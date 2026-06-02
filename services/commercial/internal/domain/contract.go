package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	ContractStatusPendingSignature = "Pending Signature"
	ContractStatusSigned           = "Signed"
	ContractStatusActive           = "Active"
	ContractStatusCompleted        = "Completed"
	ContractStatusTerminated       = "Terminated"
)

var (
	ErrContractQuoteInvalid           = errors.New("contract quote invalid")
	ErrAmountDifferenceReasonRequired = errors.New("amount difference reason required")
	ErrContractAlreadyExists          = errors.New("contract already exists")
	ErrInvalidContractTransition      = errors.New("invalid contract transition")
	ErrSignedEffectiveDateRequired    = errors.New("signed effective date required")
)

type Contract struct {
	ID                     string
	QuoteID                string
	OpportunityID          string
	CustomerID             string
	Amount                 string
	Status                 string
	ContractNote           string
	ExpectedSignedDate     time.Time
	SignedEffectiveDate    time.Time
	AmountDifferenceReason string
	OwnerID                string
	ArchivedAt             time.Time
	ArchivedBy             string
	ArchiveReason          string
	Version                int
	UpdatedAt              time.Time
}

func NewContract(input Contract, quote Quote) (Contract, error) {
	input.QuoteID = strings.TrimSpace(input.QuoteID)
	input.OpportunityID = strings.TrimSpace(input.OpportunityID)
	input.CustomerID = strings.TrimSpace(input.CustomerID)
	input.Amount = normalizeAmount(input.Amount)
	input.Status = strings.TrimSpace(input.Status)
	input.ContractNote = strings.TrimSpace(input.ContractNote)
	input.AmountDifferenceReason = strings.TrimSpace(input.AmountDifferenceReason)
	input.OwnerID = strings.TrimSpace(input.OwnerID)

	if input.QuoteID == "" || input.OpportunityID == "" || input.CustomerID == "" || input.Amount == "" ||
		input.Status == "" || input.ContractNote == "" || input.ExpectedSignedDate.IsZero() || input.OwnerID == "" {
		return Contract{}, ErrValidation
	}
	if input.Status != ContractStatusPendingSignature {
		return Contract{}, ErrValidation
	}
	if quote.ID != input.QuoteID || quote.OpportunityID != input.OpportunityID || quote.CustomerID != input.CustomerID || quote.Status != StatusAccepted {
		return Contract{}, ErrContractQuoteInvalid
	}
	if input.Amount != normalizeAmount(quote.Amount) && input.AmountDifferenceReason == "" {
		return Contract{}, ErrAmountDifferenceReasonRequired
	}
	input.Version = 1
	return input, nil
}

func CanArchiveCommercialRecord(actorRole string) bool {
	return actorRole == "Administrator" || actorRole == "Sales Manager"
}
