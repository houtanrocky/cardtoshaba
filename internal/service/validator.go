package service

import (
	"cardtoshaba/internal/domain"
	"regexp"
	"strings"
)

var cardRegex = regexp.MustCompile(`^\d{16}$`)

type CardValidatorImpl struct{}

func NewCardValidator() domain.CardValidator {
	return &CardValidatorImpl{}
}

func (v *CardValidatorImpl) Validate(cardNumber string) error {
	cleaned := strings.ReplaceAll(cardNumber, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	if cleaned == "" {
		return domain.ErrEmptyCard
	}
	if !cardRegex.MatchString(cleaned) {
		return domain.ErrInvalidCardFormat
	}
	if !luhnCheck(cleaned) {
		return domain.ErrInvalidCard
	}
	return nil
}

func luhnCheck(cardNumber string) bool {
	sum := 0
	alt := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
		n := int(cardNumber[i] - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

// CleanCardNumber - remove unneeded characters
func CleanCardNumber(cardNumber string) string {
	cleaned := strings.ReplaceAll(cardNumber, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	return cleaned
}
