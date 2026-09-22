package domain

import "context"

type CardValidator interface {
	Validate(cardNumber string) error
}

type ZarinHubClient interface {
	ConvertCardToSheba(ctx context.Context, cardNumber string) (*ConvertResponse, error)
}

type AuditRepository interface {
	Save(ctx context.Context, audit *AuditLog) error
	FindByRequestID(ctx context.Context, requestID string) (*AuditLog, error)
}

type ConvertService interface {
	Convert(ctx context.Context, req ConvertRequest) (*ConvertResult, error)
}
