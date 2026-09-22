package handler

import (
	"cardtoshaba/internal/domain"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type ConvertHandler struct {
	service domain.ConvertService
	logger  *slog.Logger
}

func NewConvertHandler(s domain.ConvertService, l *slog.Logger) *ConvertHandler {
	return &ConvertHandler{service: s, logger: l}
}

func (h *ConvertHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// ─── CORS ───
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
		"https://cardtoshaba.konkurcast.workers.dev", // ← HARDCODE موقت
	}

	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		for _, o := range strings.Split(origins, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				allowedOrigins = append(allowedOrigins, o)
			}
		}
	}

	h.logger.Info("CORS allowed origins", "origins", allowedOrigins)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"X-Request-ID", "Link", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
		Debug:            true,
	}))

	// ─── Middlewares ───
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(RequestContextMiddleware(h.logger))
	r.Use(LoggingMiddleware(h.logger))

	// ─── Routes ───
	r.Post("/api/v1/convert", h.Convert)
	r.Get("/healthz", h.Health)

	return r
}

func (h *ConvertHandler) Convert(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		WriteError(w, r, domain.NewAppError(
			"INVALID_REQUEST",
			"خطا در خواندن درخواست",
			http.StatusBadRequest,
			err,
		), h.logger)
		return
	}
	defer r.Body.Close()

	var req domain.ConvertRequest
	if err := json.Unmarshal(body, &req); err != nil {
		WriteError(w, r, domain.NewAppError(
			"INVALID_JSON",
			"ساختار JSON نامعتبر است",
			http.StatusBadRequest,
			err,
		), h.logger)
		return
	}

	result, err := h.service.Convert(r.Context(), req)
	if err != nil {
		WriteError(w, r, err, h.logger)
		return
	}

	WriteJSON(w, http.StatusOK, result)
}

func (h *ConvertHandler) Health(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func getRequestID(r *http.Request) string {
	if id := r.Context().Value("request_id"); id != nil {
		if s, ok := id.(string); ok {
			return s
		}
	}
	return ""
}
