package zarinhub

import (
	"bytes"
	"cardtoshaba/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const (
	convertEndpoint = "/api/v5/KYC/CardToIban"
)

type Client struct {
	baseURL      string
	appName      string
	tokenManager *TokenManager
	httpClient   *http.Client
	logger       *slog.Logger
}

func NewClient(
	baseURL, username, password, appName string,
	timeout time.Duration,
	logger *slog.Logger,
) domain.ZarinHubClient {
	httpClient := &http.Client{
		Timeout: timeout,
	}
	tokenManager := NewTokenManager(baseURL, username, password, appName, httpClient, logger)

	return &Client{
		baseURL:      baseURL,
		appName:      appName,
		tokenManager: tokenManager,
		httpClient:   httpClient,
		logger:       logger,
	}
}

func (c *Client) ConvertCardToSheba(
	ctx context.Context,
	cardNumber string,
) (*domain.ConvertResponse, error) {
	// Get valid token
	token, err := c.tokenManager.GetToken(ctx)
	if err != nil {
		c.logger.Error("failed to get auth token", "error", err.Error())
		return nil, domain.ErrUnauthorized
	}

	// Make Request Body (card field)
	reqBody := ConvertRequest{Card: cardNumber}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, domain.NewAppError(
			"INTERNAL_ERROR",
			"خطا در ساخت درخواست",
			http.StatusInternalServerError,
			err,
		)
	}

	// Sending request
	endpoint := c.baseURL + convertEndpoint
	resp, body, err := c.doConvert(ctx, endpoint, bodyBytes, token)
	if err != nil {
		return nil, err
	}

	// If got 401, invalidate token and retry
	if resp.StatusCode == http.StatusUnauthorized {
		c.logger.Warn("received 401, invalidating token and retrying once")
		c.tokenManager.invalidate()

		newToken, err := c.tokenManager.GetToken(ctx)
		if err != nil {
			return nil, domain.ErrUnauthorized
		}

		resp, body, err = c.doConvert(ctx, endpoint, bodyBytes, newToken)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return nil, domain.ErrUnauthorized
		}
	}

	// Unmarshal response
	var zResp Response
	if err := json.Unmarshal(body, &zResp); err != nil {
		c.logger.Error("failed to parse zarinhub response",
			"error", err.Error(),
			"body_preview", string(body[:minInt(len(body), 500)]),
		)
		return nil, domain.NewAppError(
			"ZARINHUB_INVALID_RESPONSE",
			"پاسخ سرویس تبدیل نامعتبر است",
			http.StatusBadGateway,
			err,
		)
	}

	// Checking meta
	if !zResp.Meta.IsSuccess {
		c.logger.Warn("zarinhub returned business error",
			"code", zResp.Meta.Code,
			"status", zResp.Meta.Status,
			"error_type", ptrToString(zResp.Meta.ErrorType),
			"error_message", ptrToString(zResp.Meta.ErrorMessage),
		)
		return nil, c.mapZarinHubError(&zResp.Meta)
	}

	// Validate iban exists
	if zResp.Data.Iban == "" {
		c.logger.Warn("zarinhub success but no iban returned",
			"track_id", ptrToString(zResp.Meta.TrackID),
		)
		return nil, domain.NewAppError(
			"ZARINHUB_EMPTY_IBAN",
			"شبا از سرویس تبدیل دریافت نشد",
			http.StatusBadGateway,
			nil,
		)
	}

	// Return successful result
	return &domain.ConvertResponse{
		Sheba:        zResp.Data.Iban,
		BankName:     zResp.Data.BankName,
		AccountOwner: zResp.Data.DepositOwners,
		Deposit:      zResp.Data.Deposit,
		TrackID:      ptrToString(zResp.Meta.TrackID),
		RawResponse:  body,
	}, nil
}

// doConvert - Send request with token
func (c *Client) doConvert(
	ctx context.Context,
	endpoint string,
	bodyBytes []byte,
	token string,
) (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, nil, domain.ErrInternal
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	if c.appName != "" {
		req.Header.Set("X-ZarinHub-Application", c.appName)
	}

	c.logger.Debug("sending request to zarinhub",
		"url", endpoint,
		"card_masked", domain.MaskCardNumber(string(bodyBytes)),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("zarinhub http request failed",
			"url", endpoint,
			"error", err.Error(),
		)
		return nil, nil, c.handleNetworkError(ctx, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, nil, domain.NewAppError(
			"INTERNAL_ERROR",
			"خطا در خواندن پاسخ سرویس",
			http.StatusInternalServerError,
			err,
		)
	}

	c.logger.Debug("zarinhub response received",
		"status_code", resp.StatusCode,
		"body_size", len(body),
	)

	return resp, body, nil
}

func (c *Client) mapZarinHubError(meta *Meta) error {
	if meta.ErrorType != nil {
		switch *meta.ErrorType {
		case "BadRequest", "validation_error":
			return domain.NewAppError(
				"ZARINHUB_VALIDATION_ERROR",
				"اطلاعات ارسالی نامعتبر است",
				http.StatusBadRequest,
				fmt.Errorf("zarinhub validation: %s", ptrToString(meta.ErrorMessage)),
			)
		case "UnAuthorized", "authentication_error", "authorization_error":
			return domain.ErrUnauthorized
		case "Inaccessible":
			return domain.NewAppError(
				"ACCOUNT_DISABLED",
				"حساب کاربری غیرفعال است",
				http.StatusForbidden,
				fmt.Errorf("zarinhub: account disabled"),
			)
		}
	}

	message := ptrToString(meta.ErrorMessage)
	if message == "" {
		message = meta.Message
	}
	if message == "" {
		message = "درخواست توسط سرویس تبدیل رد شد"
	}

	return domain.NewAppError(
		"ZARINHUB_REJECTED",
		message,
		http.StatusUnprocessableEntity,
		fmt.Errorf("zarinhub code=%d status=%s", meta.Code, meta.Status),
	)
}

func (c *Client) handleNetworkError(ctx context.Context, err error) error {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return domain.ErrZarinHubTimeout
	}
	if errors.Is(err, context.Canceled) {
		return domain.NewAppError(
			"REQUEST_CANCELED",
			"درخواست لغو شد",
			http.StatusRequestTimeout,
			err,
		)
	}
	return domain.ErrZarinHubUnavailable
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
