package domain

import (
	"errors"
	"testing"
	"time"
)

func TestContractDomainAcceptance(t *testing.T) {
	quote := Quote{
		ID:            "quote_domain_001",
		OpportunityID: "opp_domain_001",
		CustomerID:    "acct_domain_001",
		Amount:        "10000.00",
		Status:        StatusAccepted,
		OwnerID:       "sales-1",
	}

	t.Run("TEST-CONTRACT-CREATE-001 allows Pending Signature without signed date", func(t *testing.T) {
		contract, err := NewContract(validContractInput(), quote)
		if err != nil {
			t.Fatalf("expected valid contract create: %v", err)
		}
		if contract.Status != ContractStatusPendingSignature {
			t.Fatalf("expected Pending Signature, got %s", contract.Status)
		}
		if contract.Version != 1 {
			t.Fatalf("expected version 1, got %d", contract.Version)
		}
	})

	t.Run("TEST-CONTRACT-CREATE-003 rejects non-Accepted quote link", func(t *testing.T) {
		draft := quote
		draft.Status = StatusDraft
		_, err := NewContract(validContractInput(), draft)
		if !errors.Is(err, ErrContractQuoteInvalid) {
			t.Fatalf("expected ErrContractQuoteInvalid, got %v", err)
		}
		expired := quote
		expired.Status = StatusExpired
		_, err = NewContract(validContractInput(), expired)
		if !errors.Is(err, ErrContractQuoteInvalid) {
			t.Fatalf("expected ErrContractQuoteInvalid for expired quote, got %v", err)
		}
	})

	t.Run("TEST-CONTRACT-AMOUNT-DIFF-001 requires amount difference reason", func(t *testing.T) {
		input := validContractInput()
		input.Amount = "12000.00"
		_, err := NewContract(input, quote)
		if !errors.Is(err, ErrAmountDifferenceReasonRequired) {
			t.Fatalf("expected ErrAmountDifferenceReasonRequired, got %v", err)
		}
		input.AmountDifferenceReason = "Scope expanded after quote acceptance"
		if _, err := NewContract(input, quote); err != nil {
			t.Fatalf("expected diff contract with reason allowed: %v", err)
		}
	})
}

func validContractInput() Contract {
	return Contract{
		QuoteID:            "quote_domain_001",
		OpportunityID:      "opp_domain_001",
		CustomerID:         "acct_domain_001",
		Amount:             "10000.00",
		Status:             ContractStatusPendingSignature,
		ContractNote:       "TASK-018 contract note",
		ExpectedSignedDate: time.Date(2027, 1, 15, 0, 0, 0, 0, time.UTC),
		OwnerID:            "sales-1",
	}
}
