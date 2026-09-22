package handler

import (
	"cardtoshaba/internal/domain"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// WriteError - نوشتن خطای یکنواخت JSON
func WriteError(w http.ResponseWriter, r *http.Request, err error, logger *slog.Logger) {
	requestID := getRequestIDFromContext(r)

	var appErr *domain.AppError
	if !errors.As(err, &appErr) {
		// خطای ناشناخته - تبدیل به خطای داخلی
		logger.Error("unhandled error",
			"request_id", requestID,
			"error", err.Error(),
		)
		appErr = domain.ErrInternal
	} else {
		// Log خطاهای غیر 4xx
		if appErr.HTTPStatus >= 500 {
			logger.Error("request failed",
				"request_id", requestID,
				"code", appErr.Code,
				"error", appErr.Error(),
				"cause", appErr.Unwrap(),
			)
		}
	}

	response := ErrorResponse{
		Error: ErrorBody{
			Code:      appErr.Code,
			Message:   appErr.Message,
			RequestID: requestID,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(appErr.HTTPStatus)
	_ = json.NewEncoder(w).Encode(response)
}

// WriteJSON - نوشتن پاسخ موفق JSON
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func getRequestIDFromContext(r *http.Request) string {
	if id := r.Context().Value("request_id"); id != nil {
		if s, ok := id.(string); ok {
			return s
		}
	}
	return ""
}
