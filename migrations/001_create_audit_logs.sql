CREATE TABLE IF NOT EXISTS audit_logs (
                                          id                 BIGSERIAL PRIMARY KEY,
                                          request_id         VARCHAR(255) NOT NULL,
    card_number_masked VARCHAR(19) NOT NULL,
    sheba_number       VARCHAR(26),
    status             VARCHAR(20) NOT NULL,
    http_status        INTEGER,
    zarinhub_response  JSONB,
    duration_ms        BIGINT,
    client_ip          INET,
    error_code         VARCHAR(50),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_audit_logs_request_id ON audit_logs(request_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_status ON audit_logs(status);