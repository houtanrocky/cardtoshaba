package service

import (
	"cardtoshaba/internal/domain"
	"context"
	"errors"
	"log/slog"
	"time"
)

type convertService struct {
	validator domain.CardValidator
	client    domain.ZarinHubClient
	repo      domain.AuditRepository
	logger    *slog.Logger
}

func NewConvertService(
	v domain.CardValidator,
	c domain.ZarinHubClient,
	r domain.AuditRepository,
	l *slog.Logger,
) domain.ConvertService {
	return &convertService{
		validator: v,
		client:    c,
		repo:      r,
		logger:    l,
	}
}

func (s *convertService) Convert(ctx context.Context, req domain.ConvertRequest) (*domain.ConvertResult, error) {
	startTime := time.Now()
	requestID := getRequestID(ctx)

	// Validation
	if err := s.validator.Validate(req.CardNumber); err != nil {
		s.logger.Warn("card validation failed",
			"request_id", requestID,
			"error", err.Error(),
		)
		// Audit log for error requests
		s.saveAudit(ctx, &domain.AuditLog{
			RequestID:        requestID,
			CardNumberMasked: domain.MaskCardNumber(req.CardNumber),
			Status:           domain.AuditStatusFailed,
			HTTPStatus:       httpStatusFromError(err),
			DurationMs:       time.Since(startTime).Milliseconds(),
			ErrorCode:        errorCodeFromError(err),
			ClientIP:         getClientIP(ctx),
		})
		return nil, err
	}

	cleanedCard := CleanCardNumber(req.CardNumber)

	// Call zarin hub service with timeout
	clientCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.client.ConvertCardToSheba(clientCtx, cleanedCard)
	if err != nil {
		s.logger.Error("zarinhub call failed",
			"request_id", requestID,
			"error", err.Error(),
		)

		audit := &domain.AuditLog{
			RequestID:        requestID,
			CardNumberMasked: domain.MaskCardNumber(cleanedCard),
			Status:           auditStatusFromError(err),
			HTTPStatus:       httpStatusFromError(err),
			DurationMs:       time.Since(startTime).Milliseconds(),
			ErrorCode:        errorCodeFromError(err),
			ClientIP:         getClientIP(ctx),
		}
		s.saveAudit(ctx, audit)
		return nil, err
	}

	// Save audit log for successful request
	audit := &domain.AuditLog{
		RequestID:        requestID,
		CardNumberMasked: domain.MaskCardNumber(cleanedCard),
		ShebaNumber:      resp.Sheba,
		Status:           domain.AuditStatusSuccess,
		HTTPStatus:       200,
		ZarinHubResponse: resp.RawResponse,
		DurationMs:       time.Since(startTime).Milliseconds(),
		ClientIP:         getClientIP(ctx),
	}
	s.saveAudit(ctx, audit)

	// return results
	return &domain.ConvertResult{
		Sheba:        resp.Sheba,
		BankName:     resp.BankName,
		AccountOwner: resp.AccountOwner,
		Deposit:      resp.Deposit,
		RequestID:    requestID,
	}, nil
}

// save Audit - Save Audit Log without breaking the main flow
func (s *convertService) saveAudit(ctx context.Context, audit *domain.AuditLog) {
	// Using a separate context to do saving even if the request was cancelled
	saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.repo.Save(saveCtx, audit); err != nil {
		s.logger.Error("failed to save audit log",
			"request_id", audit.RequestID,
			"error", err.Error(),
		)
		// Persistence error should not stop the applications main flow
	}
}

// Helper functions
func getRequestID(ctx context.Context) string {
	if id, ok := ctx.Value("request_id").(string); ok {
		return id
	}
	return "unknown"
}

func getClientIP(ctx context.Context) string {
	if ip, ok := ctx.Value("client_ip").(string); ok {
		return ip
	}
	return ""
}

func httpStatusFromError(err error) int {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatus
	}
	return 500
}

func errorCodeFromError(err error) string {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return "INTERNAL_ERROR"
}

func auditStatusFromError(err error) domain.AuditStatus {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		if appErr.Code == "ZARINHUB_TIMEOUT" {
			return domain.AuditStatusTimeout
		}
	}
	return domain.AuditStatusFailed
}
