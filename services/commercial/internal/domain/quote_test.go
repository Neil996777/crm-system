package domain

import "testing"

func TestQuoteDomainAcceptance(t *testing.T) {
	t.Run("TEST-INV-ONEACCEPT-001 one quote per opportunity guard rejects duplicate", func(t *testing.T) {
		if err := EnsureCanCreateQuote(false); err != nil {
			t.Fatalf("expected first quote allowed: %v", err)
		}
		if err := EnsureCanCreateQuote(true); err != ErrQuoteAlreadyExists {
			t.Fatalf("expected quote already exists, got %v", err)
		}
	})

	t.Run("TEST-QUOTE-LIFECYCLE-003 status machine allows send/reject/expire and accept", func(t *testing.T) {
		if err := ValidateQuoteStatusTransition(StatusDraft, StatusSent); err != nil {
			t.Fatalf("draft to sent should pass: %v", err)
		}
		if err := ValidateQuoteStatusTransition(StatusSent, StatusAccepted); err != nil {
			t.Fatalf("sent to accepted should pass: %v", err)
		}
		if err := ValidateQuoteStatusTransition(StatusDraft, StatusAccepted); err != ErrInvalidQuoteTransition {
			t.Fatalf("draft to accepted must be rejected, got %v", err)
		}
	})
}
