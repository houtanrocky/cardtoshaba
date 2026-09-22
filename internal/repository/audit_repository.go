package repository

import (
	"cardtoshaba/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type PostgresAuditRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresAuditRepository(db *sql.DB, logger *slog.Logger) domain.AuditRepository {
	return &PostgresAuditRepository{db: db, logger: logger}
}

func (r *PostgresAuditRepository) Save(ctx context.Context, audit *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			request_id, card_number_masked, sheba_number, status,
			http_status, zarinhub_response, duration_ms, client_ip,
			error_code, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
	`

	var sheba interface{}
	if audit.ShebaNumber != "" {
		sheba = audit.ShebaNumber
	}

	var zarinResp interface{}
	if len(audit.ZarinHubResponse) > 0 {
		zarinResp = audit.ZarinHubResponse
	}

	var clientIP interface{}
	if audit.ClientIP != "" {
		clientIP = audit.ClientIP
	}

	var errorCode interface{}
	if audit.ErrorCode != "" {
		errorCode = audit.ErrorCode
	}

	_, err := r.db.ExecContext(ctx, query,
		audit.RequestID,
		audit.CardNumberMasked,
		sheba,
		string(audit.Status),
		audit.HTTPStatus,
		zarinResp,
		audit.DurationMs,
		clientIP,
		errorCode,
	)
	if err != nil {
		r.logger.Error("failed to save audit log",
			"request_id", audit.RequestID,
			"error", err.Error(),
		)
		return fmt.Errorf("save audit log: %w", err)
	}

	return nil
}

func (r *PostgresAuditRepository) FindByRequestID(ctx context.Context, requestID string) (*domain.AuditLog, error) {
	query := `
		SELECT id, request_id, card_number_masked, sheba_number, status,
		       http_status, zarinhub_response, duration_ms, client_ip,
		       error_code, created_at
		FROM audit_logs
		WHERE request_id = $1
		LIMIT 1
	`

	var audit domain.AuditLog
	var sheba, clientIP, errorCode sql.NullString
	var zarinResp []byte

	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&audit.ID,
		&audit.RequestID,
		&audit.CardNumberMasked,
		&sheba,
		&audit.Status,
		&audit.HTTPStatus,
		&zarinResp,
		&audit.DurationMs,
		&clientIP,
		&errorCode,
		&audit.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find audit log: %w", err)
	}

	audit.ShebaNumber = sheba.String
	audit.ClientIP = clientIP.String
	audit.ErrorCode = errorCode.String
	audit.ZarinHubResponse = zarinResp

	return &audit, nil
}
