package handler

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// RequestContextMiddleware - افزودن request_id و client_ip به context
func RequestContextMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Request ID از chi middleware
			requestID := middleware.GetReqID(r.Context())

			// Client IP (first valid IP from X-Forwarded-For)
			clientIP := getClientIP(r)

			// افزودن به context
			ctx := context.WithValue(r.Context(), "request_id", requestID)
			ctx = context.WithValue(ctx, "client_ip", clientIP)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// LoggingMiddleware - لاگ ساخت‌یافته برای هر درخواست
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			logger.Info("http request",
				"request_id", middleware.GetReqID(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration_ms", duration.Milliseconds(),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}

// getClientIP - استخراج IP واقعی client (اولین IP معتبر از X-Forwarded-For)
func getClientIP(r *http.Request) string {
	// بررسی X-Forwarded-For — ممکن است چند IP داشته باشد
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		for _, part := range parts {
			ip := strings.TrimSpace(part)
			if parsed := net.ParseIP(ip); parsed != nil {
				return parsed.String()
			}
		}
	}

	// بررسی X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if parsed := net.ParseIP(strings.TrimSpace(xri)); parsed != nil {
			return parsed.String()
		}
	}

	// RemoteAddr (احتمالاً شامل port است)
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		if parsed := net.ParseIP(host); parsed != nil {
			return parsed.String()
		}
	}

	// fallback
	return r.RemoteAddr
}
