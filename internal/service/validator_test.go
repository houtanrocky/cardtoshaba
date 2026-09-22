package service

import (
	"cardtoshaba/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardValidator_Validate(t *testing.T) {
	validator := NewCardValidator()

	tests := []struct {
		name        string
		cardNumber  string
		expectedErr error
	}{
		// ✅ Valid cards (Luhn passing)
		{
			name:        "valid card - 5022291330590744 (Pasargad)",
			cardNumber:  "5022291330590744",
			expectedErr: nil,
		},
		{
			name:        "valid card - 4111111111111111 (Visa test)",
			cardNumber:  "4111111111111111",
			expectedErr: nil,
		},
		{
			name:        "valid card - 5500000000000004 (MC test)",
			cardNumber:  "5500000000000004",
			expectedErr: nil,
		},
		{
			name:        "valid card with dashes",
			cardNumber:  "5022-2913-3059-0744",
			expectedErr: nil,
		},
		{
			name:        "valid card with spaces",
			cardNumber:  "5022 2913 3059 0744",
			expectedErr: nil,
		},
		{
			name:        "valid card with mixed separators",
			cardNumber:  "5022-2913 3059-0744",
			expectedErr: nil,
		},

		// ❌ Invalid cards
		{
			name:        "empty card number",
			cardNumber:  "",
			expectedErr: domain.ErrEmptyCard,
		},
		{
			name:        "whitespace only",
			cardNumber:  "   ",
			expectedErr: domain.ErrEmptyCard,
		},
		{
			name:        "too short",
			cardNumber:  "1234",
			expectedErr: domain.ErrInvalidCardFormat,
		},
		{
			name:        "too long",
			cardNumber:  "50222913305907441",
			expectedErr: domain.ErrInvalidCardFormat,
		},
		{
			name:        "non-numeric",
			cardNumber:  "5022abcd30590744",
			expectedErr: domain.ErrInvalidCardFormat,
		},
		{
			name:        "with letters",
			cardNumber:  "502229133059074a",
			expectedErr: domain.ErrInvalidCardFormat,
		},
		{
			name:        "Luhn fails - 1234567890123456",
			cardNumber:  "1234567890123456",
			expectedErr: domain.ErrInvalidCard,
		},
		{
			name:        "Luhn fails - all ones",
			cardNumber:  "1111111111111111",
			expectedErr: domain.ErrInvalidCard,
		},
		{
			name:        "Luhn fails - zarin-hub doc example 6037997512345678",
			cardNumber:  "6037997512345678",
			expectedErr: domain.ErrInvalidCard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.cardNumber)
			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}

func TestLuhnCheck(t *testing.T) {
	tests := []struct {
		name       string
		cardNumber string
		expected   bool
	}{
		// ✅ Valid cards (Luhn passing - standard test numbers)
		{"valid - 5022291330590744 (Pasargad)", "5022291330590744", true},
		{"valid - 4111111111111111 (Visa test)", "4111111111111111", true},
		{"valid - 5500000000000004 (MC test)", "5500000000000004", true},
		{"valid - 4012888888881881 (Visa test)", "4012888888881881", true},

		// ❌ Invalid cards (Luhn failing)
		{"invalid - 1234567890123456", "1234567890123456", false},
		{"invalid - 1111111111111111", "1111111111111111", false},
		{"invalid - 5500000000000005", "5500000000000005", false},
		{"invalid - zarin-hub doc example 6037997512345678", "6037997512345678", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, luhnCheck(tt.cardNumber))
		})
	}
}

func TestCleanCardNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no separators", "5022291330590744", "5022291330590744"},
		{"with dashes", "5022-2913-3059-0744", "5022291330590744"},
		{"with spaces", "5022 2913 3059 0744", "5022291330590744"},
		{"mixed", "5022-2913 3059-0744", "5022291330590744"},
		{"leading/trailing spaces", "  5022291330590744  ", "5022291330590744"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CleanCardNumber(tt.input))
		})
	}
}
