package zarinhub

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockGetTokenServer - سرور Mock برای GetToken
func mockGetTokenServer(t *testing.T, responses []tokenResponse, statusCode int) *httptest.Server {
	var callCount int32
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&callCount, 1) - 1
		if int(current) >= len(responses) {
			current = int32(len(responses) - 1)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(responses[current])
	}))
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestTokenManager_GetToken_Success(t *testing.T) {
	server := mockGetTokenServer(t, []tokenResponse{
		{
			Meta: tokenMeta{IsSuccess: true, Code: 0},
			Data: tokenData{
				AccessToken:  "access-token-1",
				RefreshToken: "refresh-token-1",
				ExpiresIn:    28800,
				TokenType:    "Bearer",
			},
		},
	}, http.StatusOK)
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	tm := NewTokenManager(server.URL, "user", "pass", "app", httpClient, newTestLogger())

	token, err := tm.GetToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "access-token-1", token)

	// توکن دوم باید از cache بیاید
	token2, err := tm.GetToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "access-token-1", token2)
}

func TestTokenManager_GetToken_Unauthorized(t *testing.T) {
	server := mockGetTokenServer(t, []tokenResponse{
		{
			Meta: tokenMeta{
				IsSuccess:    false,
				Code:         6,
				ErrorMessage: stringPtr("خطای احراز هویت"),
			},
		},
	}, http.StatusUnauthorized)
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	tm := NewTokenManager(server.URL, "user", "wrong", "app", httpClient, newTestLogger())

	token, err := tm.GetToken(context.Background())
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "خطای احراز هویت")
}

func TestTokenManager_RefreshToken_UsedFirst(t *testing.T) {
	// دو پاسخ: اول برای refreshToken، دوم برای احتیاط
	server := mockGetTokenServer(t, []tokenResponse{
		{
			Meta: tokenMeta{IsSuccess: true, Code: 0},
			Data: tokenData{
				AccessToken:  "refreshed-access",
				RefreshToken: "refreshed-refresh",
				ExpiresIn:    28800,
			},
		},
	}, http.StatusOK)
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	tm := NewTokenManager(server.URL, "user", "pass", "app", httpClient, newTestLogger())

	// پیش‌تنظیم refreshToken
	tm.refreshToken = "existing-refresh-token"

	token, err := tm.GetToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "refreshed-access", token)
	assert.Equal(t, "refreshed-refresh", tm.refreshToken)
}

func TestTokenManager_Invalidate(t *testing.T) {
	server := mockGetTokenServer(t, []tokenResponse{
		{
			Meta: tokenMeta{IsSuccess: true, Code: 0},
			Data: tokenData{
				AccessToken:  "token-1",
				RefreshToken: "refresh-1",
				ExpiresIn:    28800,
			},
		},
		{
			Meta: tokenMeta{IsSuccess: true, Code: 0},
			Data: tokenData{
				AccessToken:  "token-2",
				RefreshToken: "refresh-2",
				ExpiresIn:    28800,
			},
		},
	}, http.StatusOK)
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	tm := NewTokenManager(server.URL, "user", "pass", "app", httpClient, newTestLogger())

	// اولین توکن
	token1, err := tm.GetToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "token-1", token1)

	// باطل کردن
	tm.invalidate()

	// توکن جدید باید گرفته شود
	token2, err := tm.GetToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "token-2", token2)
}

func TestTokenManager_ExpiredToken_Refreshes(t *testing.T) {
	server := mockGetTokenServer(t, []tokenResponse{
		{
			Meta: tokenMeta{IsSuccess: true, Code: 0},
			Data: tokenData{
				AccessToken:  "token-1",
				RefreshToken: "refresh-1",
				ExpiresIn:    28800,
			},
		},
		{
			Meta: tokenMeta{IsSuccess: true, Code: 0},
			Data: tokenData{
				AccessToken:  "token-2",
				RefreshToken: "refresh-2",
				ExpiresIn:    28800,
			},
		},
	}, http.StatusOK)
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	tm := NewTokenManager(server.URL, "user", "pass", "app", httpClient, newTestLogger())

	// توکن اول
	token1, _ := tm.GetToken(context.Background())
	assert.Equal(t, "token-1", token1)

	// شبیه‌سازی انقضا
	tm.mu.Lock()
	tm.expiresAt = time.Now().Add(-1 * time.Hour)
	tm.mu.Unlock()

	// توکن دوم باید گرفته شود
	token2, _ := tm.GetToken(context.Background())
	assert.Equal(t, "token-2", token2)
}

// helper
func stringPtr(s string) *string {
	return &s
}
