# Card to Sheba Converter

A full-stack service that validates Iranian card numbers, converts them to IBAN (Sheba) via [ZarinHub](https://zarin-hub.com) API, persists every request to an audit log, and exposes a clean HTTP API with a modern Persian RTL web interface.

**Live:**
- Frontend: https://cardtoshaba.konkurcast.workers.dev
- Backend: https://cardtoshaba.fly.dev

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Requirements](#requirements)
- [Configuration](#configuration)
- [Running Locally](#running-locally)
- [Database Migration](#database-migration)
- [API Reference](#api-reference)
- [Error Codes](#error-codes)
- [Frontend](#frontend)
- [Testing](#testing)
- [Deployment](#deployment)
- [Design Decisions](#design-decisions)
- [Assumptions](#assumptions)
- [Limitations](#limitations)
- [Security Notes](#security-notes)

---

## Overview

This service accepts a 16-digit Iranian card number, validates it (format + Luhn), then calls ZarinHub's `CardToIban` endpoint to retrieve the corresponding Sheba (IBAN). Every attempt — successful or failed — is recorded in a PostgreSQL audit log with masked card data.

**Stack**

| Layer | Technology |
|-------|------------|
| Backend | Go 1.23 |
| Router | `go-chi/chi` + `go-chi/cors` |
| Database | PostgreSQL 16 (Neon) |
| External API | ZarinHub v5 (OAuth2 Bearer) |
| Logging | `log/slog` (JSON) |
| Frontend | Nuxt 4 + Vue 3 + Tailwind CSS |
| Hosting (Backend) | Fly.io |
| Hosting (Frontend) | Cloudflare Workers |
| CI/CD | GitHub Actions |

---

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                    HTTP Layer                        │
│      handler/  (routing, JSON, error mapping, CORS)  │
└──────────────────────┬───────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────┐
│              Application Layer                       │
│   service/  (validation, orchestration, audit)       │
└──────────┬──────────────────────────┬────────────────┘
           │                          │
┌──────────▼───────────┐   ┌──────────▼────────────────┐
│  ZarinHub Client     │   │  Audit Repository         │
│  client/zarinhub/    │   │  repository/              │
│  (token manager +    │   │  (PostgreSQL)             │
│   HTTP calls)        │   │                           │
└──────────────────────┘   └───────────────────────────┘
```

**Layers**

- `handler` — HTTP concerns only. No business logic.
- `service` — validation, workflow, audit orchestration.
- `client/zarinhub` — external API calls + automatic token refresh.
- `repository` — PostgreSQL persistence.

**Interfaces** (in `internal/domain/interfaces.go`) decouple each layer for testability.

---

## Requirements

- Go **1.23+**
- Node.js **22+** (for frontend)
- PostgreSQL **14+** (or a free Neon account)
- ZarinHub account with API credentials

---

## Configuration

Copy the example file and fill in your credentials:

```bash
cp .env.example .env
```

| Variable | Required | Description |
|----------|----------|-------------|
| `SERVER_PORT` | no | HTTP port (default `8080`) |
| `ENVIRONMENT` | no | `development` \| `production` |
| `LOG_LEVEL` | no | `debug` \| `info` \| `warn` \| `error` |
| `ZARINHUB_BASE_URL` | **yes** | e.g. `https://zarin-hub.com` |
| `ZARINHUB_API_USERNAME` | **yes** | ZarinHub API username (GUID) |
| `ZARINHUB_API_PASSWORD` | **yes** | ZarinHub API password |
| `ZARINHUB_APPLICATION_NAME` | **yes** | Registered application name |
| `ZARINHUB_TIMEOUT` | no | HTTP timeout (default `15s`) |
| `DATABASE_URL` | **yes** | PostgreSQL DSN |
| `CORS_ALLOWED_ORIGINS` | **yes** (prod) | Comma-separated list of allowed origins |

**Frontend `.env` (in `frontend/`):**

| Variable | Required | Description |
|----------|----------|-------------|
| `NUXT_PUBLIC_API_BASE_URL` | **yes** | Backend base URL (e.g. `https://cardtoshaba.fly.dev`) |

**Never commit `.env`.** It is excluded via `.gitignore`.

---

## Running Locally

### Backend

```bash
# 1. Install dependencies
go mod download

# 2. Start PostgreSQL
docker-compose up -d postgres

# 3. Apply migration (see next section)

# 4. Run
go run ./cmd/server
```

Expected output:

```json
{"level":"INFO","msg":"starting server","environment":"development"}
{"level":"INFO","msg":"database connected"}
{"level":"INFO","msg":"server listening","port":"8080"}
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Expected: `http://localhost:3000`

---

## Database Migration

Apply the schema:

```bash
psql "$DATABASE_URL" -f migrations/001_create_audit_logs.sql
```

Or, inline:

```sql
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
```

---

## API Reference

### `POST /api/v1/convert`

Convert a card number to Sheba.

**Request**

```http
POST /api/v1/convert HTTP/1.1
Content-Type: application/json

{
  "card_number": "5022291330590744"
}
```

**Success Response — `200 OK`**

```json
{
  "sheba": "IR650570310180012181077101",
  "bank_name": "بانک پاسارگاد",
  "account_owner": "هوتن عباسی راکی",
  "deposit": "3101.8000.12181077.1",
  "request_id": "d8d3e76b7d0e78/C3jgaWDeu1-000006"
}
```

**Error Response**

```json
{
  "error": {
    "code": "INVALID_CARD_FORMAT",
    "message": "شماره کارت باید ۱۶ رقم باشد",
    "request_id": "d8d3e76b7d0e78/C3jgaWDeu1-000004"
  }
}
```

### `GET /healthz`

Liveness probe.

```json
{ "status": "ok" }
```

---

## Error Codes

| Code | HTTP | Meaning |
|------|------|---------|
| `EMPTY_CARD_NUMBER` | 400 | Card number is empty |
| `INVALID_CARD_FORMAT` | 400 | Not 16 digits |
| `INVALID_CARD_NUMBER` | 400 | Luhn check failed |
| `INVALID_JSON` | 400 | Malformed request body |
| `UNAUTHORIZED` | 401 | ZarinHub auth failed |
| `ACCOUNT_DISABLED` | 403 | ZarinHub account inactive |
| `ZARINHUB_VALIDATION_ERROR` | 400 | ZarinHub rejected input |
| `ZARINHUB_REJECTED` | 422 | Business rejection by ZarinHub |
| `ZARINHUB_TIMEOUT` | 504 | Upstream timeout |
| `ZARINHUB_UNAVAILABLE` | 503 | Network failure |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

All error responses follow the same JSON shape. Sensitive fields are never leaked.

---

## Frontend

The frontend is a Nuxt 4 SPA with Persian RTL support, hosted on Cloudflare Workers.

### Features

- **Live card preview** — displays a realistic Iranian bank card matching the detected bank (BIN-based)
- **Bank detection** — 30+ Iranian banks with official brand colors and logos (`@iran-utils/iranian-banks-vue-icons`)
- **Persian digit support** — accepts `۰۱۲۳` and converts to `0123`
- **Luhn validation** — client-side, shows green checkmark when valid
- **Paste smart handling** — auto-cleans pasted card numbers
- **Toast notifications** — success/error feedback
- **Sheba + owner + bank name** displayed on the card after successful conversion
- **Copy to clipboard** for Sheba, deposit number
- **Web Share API** for mobile sharing

### Structure

```
frontend/
├── app/
│   ├── components/
│   │   ├── CardInput.vue
│   │   ├── CardPreview.vue
│   │   ├── ErrorAlert.vue
│   │   ├── ResultCard.vue
│   │   └── ToastContainer.vue
│   ├── composables/
│   │   ├── useBankDetector.ts
│   │   ├── useConvert.ts
│   │   └── useToast.ts
│   ├── layouts/
│   │   └── default.vue
│   └── pages/
│       └── index.vue
├── nuxt.config.ts
└── package.json
```

---

## Testing

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Coverage**

| Package | Statements |
|---------|------------|
| `internal/service` | ~100% |
| `internal/handler` | ~85% |
| `internal/client/zarinhub` | ~38% |

All tests use mocks for the ZarinHub client and audit repository. **No real network or credentials required.**

---

## Deployment

### Backend — Fly.io

```bash
# One-time setup
flyctl launch --no-deploy

# Configure secrets
flyctl secrets set \
  DATABASE_URL="postgres://..." \
  ZARINHUB_API_USERNAME="..." \
  ZARINHUB_API_PASSWORD="..." \
  ZARINHUB_APPLICATION_NAME="cardtoshaba" \
  ZARINHUB_BASE_URL="https://zarin-hub.com" \
  SERVER_PORT="8080" \
  ENVIRONMENT="production" \
  LOG_LEVEL="info" \
  CORS_ALLOWED_ORIGINS="https://cardtoshaba.konkurcast.workers.dev"

# Deploy
flyctl deploy
```

**Important:** Add Fly.io's outgoing IP to ZarinHub's whitelist.

```bash
flyctl ssh console -a cardtoshaba -C "wget -qO- https://api.ipify.org"
```

### Frontend — Cloudflare Workers

1. Connect GitHub repository to Cloudflare
2. Build settings:
    - **Build command:** `cd frontend && npm ci && npm run build`
    - **Deploy command:** `npx wrangler --cwd frontend/.output deploy`
3. Environment variables:
    - `NODE_VERSION` = `22`
    - `NUXT_PUBLIC_API_BASE_URL` = `https://cardtoshaba.fly.dev`

### CI/CD — GitHub Actions

`.github/workflows/backend.yml` runs tests and deploys to Fly.io on every push to `main`/`master`.

**Required GitHub Secrets:**
- `FLY_API_TOKEN` — obtain with `flyctl auth token`

---

## Design Decisions

### Layered architecture with interfaces

Handlers only decode/encode HTTP. Business logic lives in the service layer. The ZarinHub client and repository are behind interfaces — enabling fast, isolated unit tests and future swaps (e.g., a different KYC provider).

### Automatic token management

ZarinHub uses short-lived Bearer tokens (8 hours). A `TokenManager` caches the token in memory, refreshes it 5 minutes before expiry, and falls back to re-authentication when refresh fails. This keeps the service resilient without hammering the auth endpoint.

### Audit-first persistence

Every attempt is logged, including validation failures and upstream errors. The audit write is intentionally decoupled from the request context (`context.Background()`) so it succeeds even if the client disconnects. Persistence failures never break the main flow — they are logged and swallowed.

### Card masking

The full card number is never persisted or logged. Only the last 4 digits are stored (`****-****-****-0744`). ZarinHub itself returns a masked card in its response, which we store as-is.

### Idempotent errors, stable codes

Application-level error codes (`AppError.Code`) are stable strings — safe for frontends to switch on. HTTP statuses are derived from these codes centrally, keeping handlers thin.

### CORS with environment-driven origins

Allowed origins are configured via `CORS_ALLOWED_ORIGINS`, split by comma. Local development origins are always allowed. This avoids hardcoding production URLs and supports multiple environments.

---

## Assumptions

1. Card numbers are always 16 digits (Iranian standard).
2. Luhn checksum is the correct validation for Iranian cards.
3. ZarinHub is the sole KYC provider.
4. `request_id` uniqueness is per-instance, not global.
5. Audit retention is outside the scope of this service.
6. IPv4/IPv6 addresses are valid; malformed ones are ignored.
7. No rate limiting is required for the assessment.
8. Fly.io's outbound IP is whitelisted in ZarinHub.

---

## Limitations

- **No retry policy** for transient ZarinHub failures (only one retry on `401`).
- **No metrics** endpoint (Prometheus, etc.).
- **No OpenAPI spec**.
- **No graceful handling of token revocation** — a full re-login happens only on `401`.
- **Fly.io trial** shuts down machines after 5 minutes of inactivity — a credit card is required for production use.
- **Fly.io uses shared IPs** — may change; a dedicated IPv4 (~$2/month) is recommended for stability.

---

## Security Notes

- **Credentials** are read from environment variables only. `.env` is git-ignored.
- **Full card numbers** are never logged or persisted; only the last 4 digits.
- **Raw ZarinHub responses** are stored as `JSONB` for auditing — they contain no credentials.
- **Error responses** never include stack traces, raw upstream bodies, or credential values.
- **TLS verification** is enforced by Go's default `http.Client`.
- **Auth tokens** are stored in memory only; they are never written to disk.
- **CORS** is configured with an explicit allowlist; wildcard origins are not used.

**Rotate your ZarinHub API password immediately if it is ever exposed in logs, chat, or source control.**

---