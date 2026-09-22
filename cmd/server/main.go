package main

import (
	"cardtoshaba/internal/client/zarinhub"
	"cardtoshaba/internal/config"
	"cardtoshaba/internal/handler"
	"cardtoshaba/internal/repository"
	"cardtoshaba/internal/service"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// ۱. بارگذاری تنظیمات
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// ۲. راه‌اندازی Logger
	logger := setupLogger(cfg.LogLevel)
	logger.Info("starting server", "environment", cfg.Environment)

	// ۳. اتصال به پایگاه داده
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("failed to close database", "error", closeErr)
		}
	}()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	logger.Info("database connected")

	// ۴. راه‌اندازی وابستگی‌ها (Dependency Injection)
	validator := service.NewCardValidator()

	zarinClient := zarinhub.NewClient(
		cfg.ZarinHubBaseURL,
		cfg.ZarinHubAPIUsername,
		cfg.ZarinHubAPIPassword,
		cfg.ZarinHubAppName,
		cfg.ZarinHubTimeout,
		logger,
	)

	auditRepo := repository.NewPostgresAuditRepository(db, logger)

	convertService := service.NewConvertService(
		validator,
		zarinClient,
		auditRepo,
		logger,
	)

	httpHandler := handler.NewConvertHandler(convertService, logger)

	// ۵. راه‌اندازی HTTP Server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      httpHandler.Routes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ۶. اجرای Server در Goroutine
	go func() {
		logger.Info("server listening", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// ۷. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server exited")
}

func setupLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})
	return slog.New(logHandler)
}
