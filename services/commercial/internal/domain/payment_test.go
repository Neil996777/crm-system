package domain

import (
	"errors"
	"testing"
)

func TestPaymentDomainAcceptance(t *testing.T) {
	contract := Contract{ID: "contract_payment_domain", Amount: "10000.00"}

	t.Run("TEST-PAYMENT-GUARD-001 and TEST-PAYMENT-GUARD-002 reject zero negative", func(t *testing.T) {
		for _, amount := range []string{"0.00", "-1.00"} {
			_, err := NewActualPayment(contract, "key-"+amount, amount, "2027-08-01", "", "CNY", "0.00")
			if !errors.Is(err, ErrInvalidAmount) {
				t.Fatalf("expected ErrInvalidAmount for %s, got %v", amount, err)
			}
		}
	})

	t.Run("TEST-PAYMENT-GUARD-003 rejects overpayment", func(t *testing.T) {
		_, err := NewActualPayment(contract, "key-over", "0.01", "2027-08-01", "", "CNY", "10000.00")
		if !errors.Is(err, ErrOverpaymentBlocked) {
			t.Fatalf("expected ErrOverpaymentBlocked, got %v", err)
		}
	})

	t.Run("TEST-INV-CURRENCY-001 rejects non-CNY currency", func(t *testing.T) {
		if _, err := NewPaymentPlan(contract.ID, "1000.00", "2027-08-01", "USD"); !errors.Is(err, ErrSingleCurrencyRequired) {
			t.Fatalf("expected ErrSingleCurrencyRequired, got %v", err)
		}
	})
}
