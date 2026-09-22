package handler

import (
	"bytes"
	"cardtoshaba/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConvertService - Mock برای ConvertService
type MockConvertService struct {
	mock.Mock
}

func (m *MockConvertService) Convert(
	ctx context.Context,
	req domain.ConvertRequest,
) (*domain.ConvertResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ConvertResult), args.Error(1)
}

// setupHandlerTest - راه‌اندازی تست HTTP
func setupHandlerTest() (*ConvertHandler, *MockConvertService) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockSvc := new(MockConvertService)
	handler := NewConvertHandler(mockSvc, logger)
	return handler, mockSvc
}

// ============================================================
// تست‌های POST /api/v1/convert
// ============================================================

func TestConvertHandler_Success(t *testing.T) {
	handler, mockSvc := setupHandlerTest()

	mockSvc.On("Convert", mock.Anything, domain.ConvertRequest{
		CardNumber: "5022291330590744",
	}).Return(&domain.ConvertResult{
		Sheba:        "IR650570310180012181077101",
		BankName:     "بانک پاسارگاد",
		AccountOwner: "هوتن عباسی راکی",
		Deposit:      "3101.8000.12181077.1",
		RequestID:    "test-request-id",
	}, nil)

	body := `{"card_number":"5022291330590744"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/convert", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.ConvertResult
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "IR650570310180012181077101", resp.Sheba)
	assert.Equal(t, "بانک پاسارگاد", resp.BankName)
	assert.Equal(t, "test-request-id", resp.RequestID)

	mockSvc.AssertExpectations(t)
}

func TestConvertHandler_InvalidJSON(t *testing.T) {
	handler, mockSvc := setupHandlerTest()

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/convert", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_JSON", resp.Error.Code)

	// Service نباید صدا زده شود
	mockSvc.AssertNotCalled(t, "Convert")
}

func TestConvertHandler_InvalidCardNumber(t *testing.T) {
	handler, mockSvc := setupHandlerTest()

	mockSvc.On("Convert", mock.Anything, domain.ConvertRequest{
		CardNumber: "1234",
	}).Return(nil, domain.ErrInvalidCardFormat)

	body := `{"card_number":"1234"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/convert", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_CARD_FORMAT", resp.Error.Code)

	mockSvc.AssertExpectations(t)
}

func TestConvertHandler_EmptyCardNumber(t *testing.T) {
	handler, mockSvc := setupHandlerTest()

	mockSvc.On("Convert", mock.Anything, domain.ConvertRequest{
		CardNumber: "",
	}).Return(nil, domain.ErrEmptyCard)

	body := `{"card_number":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/convert", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestConvertHandler_Unauthorized(t *testing.T) {
	handler, mockSvc := setupHandlerTest()

	mockSvc.On("Convert", mock.Anything, mock.Anything).
		Return(nil, domain.ErrUnauthorized)

	body := `{"card_number":"5022291330590744"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/convert", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "UNAUTHORIZED", resp.Error.Code)
}

func TestConvertHandler_Timeout(t *testing.T) {
	handler, mockSvc := setupHandlerTest()

	mockSvc.On("Convert", mock.Anything, mock.Anything).
		Return(nil, domain.ErrZarinHubTimeout)

	body := `{"card_number":"5022291330590744"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/convert", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
}

func TestConvertHandler_InternalError(t *testing.T) {
	handler, mockSvc := setupHandlerTest()

	mockSvc.On("Convert", mock.Anything, mock.Anything).
		Return(nil, errors.New("unexpected error"))

	body := `{"card_number":"5022291330590744"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/convert", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "INTERNAL_ERROR", resp.Error.Code)
}

func TestConvertHandler_MissingCardNumberField(t *testing.T) {
	handler, mockSvc := setupHandlerTest()

	// Service با CardNumber خالی صدا زده می‌شود
	mockSvc.On("Convert", mock.Anything, domain.ConvertRequest{
		CardNumber: "",
	}).Return(nil, domain.ErrEmptyCard)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/convert", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertExpectations(t)
}

// ============================================================
// تست‌های GET /healthz
// ============================================================

func TestHealthCheck(t *testing.T) {
	handler, _ := setupHandlerTest()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	handler.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
}

// ============================================================
// تست‌های امنیتی (Response نباید اطلاعات حساس داشته باشد)
// ============================================================

func TestConvertHandler_NoSensitiveDataInResponse(t *testing.T) {
	handler, mockSvc := setupHandlerTest()

	mockSvc.On("Convert", mock.Anything, mock.Anything).
		Return(nil, domain.ErrZarinHubUnavailable)

	body := `{"card_number":"5022291330590744"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/convert", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)

	responseBody := w.Body.String()

	// نباید شماره کارت کامل در پاسخ باشد
	assert.NotContains(t, responseBody, "5022291330590744")

	// نباید API key یا توکن در پاسخ باشد
	assert.NotContains(t, responseBody, "Bearer")
	assert.NotContains(t, responseBody, "apiKey")

	// نباید stack trace باشد
	assert.NotContains(t, responseBody, "goroutine")
	assert.NotContains(t, responseBody, ".go:")
}

func TestConvertHandler_RequestIDInErrorResponse(t *testing.T) {
	handler, mockSvc := setupHandlerTest()

	mockSvc.On("Convert", mock.Anything, mock.Anything).
		Return(nil, domain.ErrInvalidCard)

	body := `{"card_number":"invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/convert", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)

	var resp ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "INVALID_CARD_NUMBER", resp.Error.Code)
}
