package service

import (
	"cardtoshaba/internal/domain"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockZarinHubClient - Mock برای ZarinHubClient
type MockZarinHubClient struct {
	mock.Mock
}

func (m *MockZarinHubClient) ConvertCardToSheba(
	ctx context.Context,
	cardNumber string,
) (*domain.ConvertResponse, error) {
	args := m.Called(ctx, cardNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ConvertResponse), args.Error(1)
}

// MockAuditRepository - Mock برای AuditRepository
type MockAuditRepository struct {
	mock.Mock
}

func (m *MockAuditRepository) Save(ctx context.Context, audit *domain.AuditLog) error {
	args := m.Called(ctx, audit)
	return args.Error(0)
}

func (m *MockAuditRepository) FindByRequestID(
	ctx context.Context,
	requestID string,
) (*domain.AuditLog, error) {
	args := m.Called(ctx, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuditLog), args.Error(1)
}

// setupTest - راه‌اندازی مشترک تست
func setupTest() (
	*convertService,
	*MockZarinHubClient,
	*MockAuditRepository,
) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	validator := NewCardValidator()
	mockClient := new(MockZarinHubClient)
	mockRepo := new(MockAuditRepository)

	svc := NewConvertService(validator, mockClient, mockRepo, logger)
	return svc.(*convertService), mockClient, mockRepo
}

// contextWithRequestID - ساخت context با request_id
func contextWithRequestID(requestID string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", requestID)
	ctx = context.WithValue(ctx, "client_ip", "127.0.0.1")
	return ctx
}

// ============================================================
// تست‌های Convert
// ============================================================

func TestConvertService_Success(t *testing.T) {
	svc, mockClient, mockRepo := setupTest()

	cardNumber := "5022291330590744"
	expectedSheba := "IR650570310180012181077101"

	mockClient.On("ConvertCardToSheba", mock.Anything, cardNumber).
		Return(&domain.ConvertResponse{
			Sheba:        expectedSheba,
			BankName:     "بانک پاسارگاد",
			AccountOwner: "هوتن عباسی راکی",
			Deposit:      "3101.8000.12181077.1",
			TrackID:      "3b2f0890-d77f-43ab-ba56-7408f6e1c8ef",
			RawResponse:  []byte(`{"data":{"iban":"IR650570310180012181077101"}}`),
		}, nil)

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	ctx := contextWithRequestID("test-request-1")
	result, err := svc.Convert(ctx, domain.ConvertRequest{
		CardNumber: cardNumber,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedSheba, result.Sheba)
	assert.Equal(t, "بانک پاسارگاد", result.BankName)
	assert.Equal(t, "هوتن عباسی راکی", result.AccountOwner)
	assert.Equal(t, "3101.8000.12181077.1", result.Deposit)
	assert.Equal(t, "test-request-1", result.RequestID)

	mockClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestConvertService_InvalidCard_Empty(t *testing.T) {
	svc, mockClient, mockRepo := setupTest()

	// برای درخواست نامعتبر، فقط Audit Log ذخیره می‌شود
	mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(a *domain.AuditLog) bool {
		return a.Status == domain.AuditStatusFailed &&
			a.ErrorCode == "EMPTY_CARD_NUMBER"
	})).Return(nil)

	ctx := contextWithRequestID("test-invalid-1")
	result, err := svc.Convert(ctx, domain.ConvertRequest{
		CardNumber: "",
	})

	assert.ErrorIs(t, err, domain.ErrEmptyCard)
	assert.Nil(t, result)

	// مطمئن شو که به زرین‌هاب درخواست نرفته
	mockClient.AssertNotCalled(t, "ConvertCardToSheba")
	mockRepo.AssertExpectations(t)
}

func TestConvertService_InvalidCard_Format(t *testing.T) {
	svc, mockClient, mockRepo := setupTest()

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	ctx := contextWithRequestID("test-invalid-2")
	result, err := svc.Convert(ctx, domain.ConvertRequest{
		CardNumber: "1234",
	})

	assert.ErrorIs(t, err, domain.ErrInvalidCardFormat)
	assert.Nil(t, result)
	mockClient.AssertNotCalled(t, "ConvertCardToSheba")
}

func TestConvertService_InvalidCard_Luhn(t *testing.T) {
	svc, mockClient, mockRepo := setupTest()

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	ctx := contextWithRequestID("test-invalid-3")
	result, err := svc.Convert(ctx, domain.ConvertRequest{
		CardNumber: "1234567890123456",
	})

	assert.ErrorIs(t, err, domain.ErrInvalidCard)
	assert.Nil(t, result)
	mockClient.AssertNotCalled(t, "ConvertCardToSheba")
}

func TestConvertService_ZarinHubUnauthorized(t *testing.T) {
	svc, mockClient, mockRepo := setupTest()

	cardNumber := "5022291330590744"

	mockClient.On("ConvertCardToSheba", mock.Anything, cardNumber).
		Return(nil, domain.ErrUnauthorized)

	mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(a *domain.AuditLog) bool {
		return a.Status == domain.AuditStatusFailed &&
			a.ErrorCode == "UNAUTHORIZED"
	})).Return(nil)

	ctx := contextWithRequestID("test-unauth")
	result, err := svc.Convert(ctx, domain.ConvertRequest{
		CardNumber: cardNumber,
	})

	assert.ErrorIs(t, err, domain.ErrUnauthorized)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestConvertService_ZarinHubTimeout(t *testing.T) {
	svc, mockClient, mockRepo := setupTest()

	cardNumber := "5022291330590744"

	mockClient.On("ConvertCardToSheba", mock.Anything, cardNumber).
		Return(nil, domain.ErrZarinHubTimeout)

	mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(a *domain.AuditLog) bool {
		return a.Status == domain.AuditStatusTimeout &&
			a.ErrorCode == "ZARINHUB_TIMEOUT"
	})).Return(nil)

	ctx := contextWithRequestID("test-timeout")
	result, err := svc.Convert(ctx, domain.ConvertRequest{
		CardNumber: cardNumber,
	})

	assert.ErrorIs(t, err, domain.ErrZarinHubTimeout)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestConvertService_ZarinHubRejected(t *testing.T) {
	svc, mockClient, mockRepo := setupTest()

	cardNumber := "5022291330590744"

	rejectedErr := domain.NewAppError(
		"ZARINHUB_REJECTED",
		"درخواست توسط سرویس تبدیل رد شد",
		422,
		errors.New("zarinhub rejected"),
	)

	mockClient.On("ConvertCardToSheba", mock.Anything, cardNumber).
		Return(nil, rejectedErr)

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	ctx := contextWithRequestID("test-rejected")
	result, err := svc.Convert(ctx, domain.ConvertRequest{
		CardNumber: cardNumber,
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "ZARINHUB_REJECTED", err.(*domain.AppError).Code)
	mockClient.AssertExpectations(t)
}

func TestConvertService_AuditLogFails_DoesNotBlockRequest(t *testing.T) {
	svc, mockClient, mockRepo := setupTest()

	cardNumber := "5022291330590744"
	expectedSheba := "IR650570310180012181077101"

	mockClient.On("ConvertCardToSheba", mock.Anything, cardNumber).
		Return(&domain.ConvertResponse{
			Sheba:    expectedSheba,
			BankName: "بانک پاسارگاد",
		}, nil)

	// Audit Log شکست می‌خورد
	mockRepo.On("Save", mock.Anything, mock.Anything).
		Return(errors.New("database connection lost"))

	ctx := contextWithRequestID("test-audit-fail")
	result, err := svc.Convert(ctx, domain.ConvertRequest{
		CardNumber: cardNumber,
	})

	// درخواست باید موفق باشد حتی اگر Audit Log شکست بخورد
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedSheba, result.Sheba)
	mockClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestConvertService_CleansCardNumber(t *testing.T) {
	svc, mockClient, mockRepo := setupTest()

	// ورودی با خط تیره
	inputCard := "5022-2913-3059-0744"
	// انتظار می‌رود که به زرین‌هاب با فرمت پاک ارسال شود
	expectedCleanCard := "5022291330590744"

	mockClient.On("ConvertCardToSheba", mock.Anything, expectedCleanCard).
		Return(&domain.ConvertResponse{
			Sheba:    "IR650570310180012181077101",
			BankName: "بانک پاسارگاد",
		}, nil)

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	ctx := contextWithRequestID("test-clean")
	result, err := svc.Convert(ctx, domain.ConvertRequest{
		CardNumber: inputCard,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockClient.AssertExpectations(t) // مطمئن می‌شویم که با شماره پاک صدا زده شده
}

// ============================================================
// تست‌های Helper Functions
// ============================================================

func TestHttpStatusFromError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"AppError 400", domain.ErrInvalidCard, 400},
		{"AppError 401", domain.ErrUnauthorized, 401},
		{"AppError 504", domain.ErrZarinHubTimeout, 504},
		{"plain error", errors.New("unknown"), 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, httpStatusFromError(tt.err))
		})
	}
}

func TestErrorCodeFromError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"AppError INVALID_CARD_NUMBER", domain.ErrInvalidCard, "INVALID_CARD_NUMBER"},
		{"AppError UNAUTHORIZED", domain.ErrUnauthorized, "UNAUTHORIZED"},
		{"plain error", errors.New("unknown"), "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, errorCodeFromError(tt.err))
		})
	}
}

func TestAuditStatusFromError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected domain.AuditStatus
	}{
		{"timeout", domain.ErrZarinHubTimeout, domain.AuditStatusTimeout},
		{"unauthorized", domain.ErrUnauthorized, domain.AuditStatusFailed},
		{"invalid card", domain.ErrInvalidCard, domain.AuditStatusFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, auditStatusFromError(tt.err))
		})
	}
}

// ============================================================
// تست Context Helpers
// ============================================================

func TestGetRequestID_FromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "request_id", "my-request-id")
	assert.Equal(t, "my-request-id", getRequestID(ctx))
}

func TestGetRequestID_NotInContext(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, "unknown", getRequestID(ctx))
}

func TestGetClientIP_FromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "client_ip", "192.168.1.1")
	assert.Equal(t, "192.168.1.1", getClientIP(ctx))
}

func TestGetClientIP_NotInContext(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, "", getClientIP(ctx))
}

// ============================================================
// تست Timeout در Service
// ============================================================

func TestConvertService_ContextTimeout(t *testing.T) {
	svc, mockClient, mockRepo := setupTest()

	cardNumber := "5022291330590744"

	// Mock کلاینت که با تأخیر پاسخ می‌دهد
	mockClient.On("ConvertCardToSheba", mock.Anything, cardNumber).
		Run(func(args mock.Arguments) {
			time.Sleep(100 * time.Millisecond)
		}).
		Return(nil, domain.ErrZarinHubTimeout)

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	ctx := contextWithRequestID("test-ctx-timeout")
	result, err := svc.Convert(ctx, domain.ConvertRequest{
		CardNumber: cardNumber,
	})

	assert.ErrorIs(t, err, domain.ErrZarinHubTimeout)
	assert.Nil(t, result)
}
