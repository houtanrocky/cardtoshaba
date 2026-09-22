package domain

import "net/http"

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	Err        error  `json:"-"` // Original error (just for log)
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(code, message string, status int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
		Err:        err,
	}
}

var (
	ErrInvalidCard = &AppError{
		Code:       "INVALID_CARD_NUMBER",
		Message:    "شماره کارت وارد شده معتبر نیست",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrEmptyCard = &AppError{
		Code:       "EMPTY_CARD_NUMBER",
		Message:    "شماره کارت نمی‌تواند خالی باشد",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrInvalidCardFormat = &AppError{
		Code:       "INVALID_CARD_FORMAT",
		Message:    "شماره کارت باید ۱۶ رقم باشد",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "خطا در احراز هویت سرویس خارجی",
		HTTPStatus: http.StatusUnauthorized,
	}
	ErrZarinHubRejected = &AppError{
		Code:       "ZARINHUB_REJECTED",
		Message:    "درخواست توسط سرویس تبدیل رد شد",
		HTTPStatus: http.StatusUnprocessableEntity,
	}
	ErrZarinHubTimeout = &AppError{
		Code:       "ZARINHUB_TIMEOUT",
		Message:    "سرویس تبدیل موقتاً در دسترس نیست",
		HTTPStatus: http.StatusGatewayTimeout,
	}
	ErrZarinHubUnavailable = &AppError{
		Code:       "ZARINHUB_UNAVAILABLE",
		Message:    "ارتباط با سرویس تبدیل برقرار نشد",
		HTTPStatus: http.StatusServiceUnavailable,
	}
	ErrInternal = &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "خطای داخلی سرور",
		HTTPStatus: http.StatusInternalServerError,
	}
	ErrPersistence = &AppError{
		Code:       "PERSISTENCE_ERROR",
		Message:    "خطا در ذخیره‌سازی اطلاعات",
		HTTPStatus: http.StatusInternalServerError,
	}
)
