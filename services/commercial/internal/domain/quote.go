package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	StatusDraft    = "Draft"
	StatusSent     = "Sent"
	StatusAccepted = "Accepted"
	StatusRejected = "Rejected"
	StatusExpired  = "Expired"
)

var (
	ErrValidation             = errors.New("validation failed")
	ErrQuoteAlreadyExists     = errors.New("quote already exists")
	ErrInvalidQuoteTransition = errors.New("invalid quote transition")
)

type Quote struct {
	ID            string
	OpportunityID string
	CustomerID    string
	Amount        string
	Status        string
	ValidityEnd   time.Time
	OwnerID       string
	Version       int
	UpdatedAt     time.Time
}

func NewQuote(input Quote) (Quote, error) {
	input.OpportunityID = strings.TrimSpace(input.OpportunityID)
	input.CustomerID = strings.TrimSpace(input.CustomerID)
	input.Amount = normalizeAmount(input.Amount)
	input.Status = strings.TrimSpace(input.Status)
	input.OwnerID = strings.TrimSpace(input.OwnerID)
	if input.OpportunityID == "" || input.CustomerID == "" || input.Amount == "" || input.Status == "" || input.ValidityEnd.IsZero() || input.OwnerID == "" {
		return Quote{}, ErrValidation
	}
	if input.Status != StatusDraft {
		return Quote{}, ErrValidation
	}
	input.Version = 1
	return input, nil
}

func EnsureCanCreateQuote(existing bool) error {
	if existing {
		return ErrQuoteAlreadyExists
	}
	return nil
}

func ValidateQuoteStatusTransition(fromStatus, toStatus string) error {
	switch fromStatus {
	case StatusDraft:
		if toStatus == StatusSent || toStatus == StatusExpired {
			return nil
		}
	case StatusSent:
		if toStatus == StatusAccepted || toStatus == StatusRejected || toStatus == StatusExpired {
			return nil
		}
	}
	return ErrInvalidQuoteTransition
}

func ParseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func FormatDate(value time.Time) string {
	return value.UTC().Format("2006-01-02")
}

func normalizeAmount(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil || amount <= 0 {
		return ""
	}
	return fmt.Sprintf("%.2f", amount)
}
