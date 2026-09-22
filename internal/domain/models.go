package domain

import "time"

// ConvertRequest - Incoming requests from frontend
type ConvertRequest struct {
	CardNumber string `json:"card_number"`
}

// ConvertResponse - Raw response from zarinhub
type ConvertResponse struct {
	Sheba        string // data.iban
	BankName     string // data.bankName
	AccountOwner string // data.depositOwners
	Deposit      string // data.deposit
	TrackID      string // meta.trackId
	RawResponse  []byte // Whole response for Audit
}

// ConvertResult - Final result sent to frontend
type ConvertResult struct {
	Sheba        string `json:"sheba"`
	BankName     string `json:"bank_name,omitempty"`
	AccountOwner string `json:"account_owner,omitempty"`
	Deposit      string `json:"deposit,omitempty"`
	RequestID    string `json:"request_id"`
}

// AuditLog - Saved record in database
type AuditLog struct {
	ID               int64       `json:"id"`
	RequestID        string      `json:"request_id"`
	CardNumberMasked string      `json:"card_number_masked"`
	ShebaNumber      string      `json:"sheba_number,omitempty"`
	Status           AuditStatus `json:"status"`
	HTTPStatus       int         `json:"http_status"`
	ZarinHubResponse []byte      `json:"zarinhub_response,omitempty"`
	DurationMs       int64       `json:"duration_ms"`
	ClientIP         string      `json:"client_ip,omitempty"`
	ErrorCode        string      `json:"error_code,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
}

type AuditStatus string

const (
	AuditStatusSuccess AuditStatus = "SUCCESS"
	AuditStatusFailed  AuditStatus = "FAILED"
	AuditStatusTimeout AuditStatus = "TIMEOUT"
)

// MaskCardNumber - Masking card number (only show the last 4 digits)
func MaskCardNumber(cardNumber string) string {
	if len(cardNumber) < 4 {
		return "****"
	}
	return "****-****-****-" + cardNumber[len(cardNumber)-4:]
}
