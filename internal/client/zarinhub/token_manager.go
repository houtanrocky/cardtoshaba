package zarinhub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

const (
	getTokenEndpoint = "/api/v5/Authentication/GetToken"
)

type tokenResponse struct {
	Meta tokenMeta `json:"meta"`
	Data tokenData `json:"data"`
}

type tokenMeta struct {
	TrackID      *string `json:"trackId"`
	Status       string  `json:"status"`
	IsSuccess    bool    `json:"isSuccess"`
	Code         int     `json:"code"`
	Message      string  `json:"message"`
	ErrorMessage *string `json:"errorMessage"`
	ErrorType    *string `json:"errorType"`
}

type tokenData struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresTn"`
	RefreshToken string `json:"refreshToken"`
}

// TokenManager - Managing zarinhub token with refresh
type TokenManager struct {
	baseURL    string
	username   string
	password   string
	appName    string
	httpClient *http.Client
	logger     *slog.Logger

	mu           sync.RWMutex
	accessToken  string
	refreshToken string
	expiresAt    time.Time
}

func NewTokenManager(
	baseURL, username, password, appName string,
	httpClient *http.Client,
	logger *slog.Logger,
) *TokenManager {
	return &TokenManager{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		appName:    appName,
		httpClient: httpClient,
		logger:     logger,
	}
}

// GetToken - Getting valid token (automated with refresh)
func (tm *TokenManager) GetToken(ctx context.Context) (string, error) {
	tm.mu.RLock()
	token := tm.accessToken
	expiresAt := tm.expiresAt
	tm.mu.RUnlock()

	// 5 mins threshold before expiring
	if token != "" && time.Now().Before(expiresAt.Add(-5*time.Minute)) {
		return token, nil
	}

	return tm.renewToken(ctx)
}

func (tm *TokenManager) renewToken(ctx context.Context) (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Double-check
	if tm.accessToken != "" && time.Now().Before(tm.expiresAt.Add(-5*time.Minute)) {
		return tm.accessToken, nil
	}

	// First try with refreshToken
	if tm.refreshToken != "" {
		reqBody := map[string]string{"refreshToken": tm.refreshToken}
		resp, err := tm.callGetToken(ctx, reqBody)
		if err == nil {
			tm.updateTokens(resp)
			tm.logger.Info("token refreshed via refreshToken",
				"expires_in", resp.Data.ExpiresIn,
			)
			return tm.accessToken, nil
		}
		tm.logger.Warn("refreshToken failed, falling back to credentials",
			"error", err.Error(),
		)
		tm.refreshToken = ""
	}

	// If not try with username/password
	reqBody := map[string]string{
		"username": tm.username,
		"password": tm.password,
	}
	resp, err := tm.callGetToken(ctx, reqBody)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	tm.updateTokens(resp)
	tm.logger.Info("token obtained via credentials",
		"expires_in", resp.Data.ExpiresIn,
	)
	return tm.accessToken, nil
}

func (tm *TokenManager) callGetToken(ctx context.Context, reqBody map[string]string) (*tokenResponse, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		tm.baseURL+getTokenEndpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if tm.appName != "" {
		req.Header.Set("X-ZarinHub-Application", tm.appName)
	}

	resp, err := tm.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	if resp.StatusCode != http.StatusOK || !tokenResp.Meta.IsSuccess {
		errMsg := "unknown"
		if tokenResp.Meta.ErrorMessage != nil {
			errMsg = *tokenResp.Meta.ErrorMessage
		} else if tokenResp.Meta.Message != "" {
			errMsg = tokenResp.Meta.Message
		}
		return nil, fmt.Errorf("get token failed (http %d, code %d): %s",
			resp.StatusCode, tokenResp.Meta.Code, errMsg)
	}

	if tokenResp.Data.AccessToken == "" {
		return nil, fmt.Errorf("access token empty")
	}

	return &tokenResp, nil
}

func (tm *TokenManager) updateTokens(resp *tokenResponse) {
	tm.accessToken = resp.Data.AccessToken
	if resp.Data.RefreshToken != "" {
		tm.refreshToken = resp.Data.RefreshToken
	}
	tm.expiresAt = time.Now().Add(time.Duration(resp.Data.ExpiresIn) * time.Second)
}

// invalidate - invalidating the current token (after 401 error for retry)
func (tm *TokenManager) invalidate() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.accessToken = ""
	tm.expiresAt = time.Time{}
}
