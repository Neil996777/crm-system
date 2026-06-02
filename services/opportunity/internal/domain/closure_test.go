package domain

import (
	"errors"
	"testing"
)

func TestOpportunityClosureDomain(t *testing.T) {
	opportunity := Opportunity{Stage: StageContractNegotiation}

	t.Run("TEST-OPP-CLOSE-002 Won requires signed contract", func(t *testing.T) {
		if _, err := CloseWon(opportunity, false, "contract_1"); !errors.Is(err, ErrEarlyWonBlocked) {
			t.Fatalf("expected ErrEarlyWonBlocked, got %v", err)
		}
		closed, err := CloseWon(opportunity, true, "contract_1")
		if err != nil {
			t.Fatalf("expected close won allowed: %v", err)
		}
		if closed.Stage != StageWon || closed.WonContractID != "contract_1" {
			t.Fatalf("expected won closure fields, got %#v", closed)
		}
	})

	t.Run("TEST-OPP-CLOSE-004 Lost requires reason", func(t *testing.T) {
		if _, err := CloseLost(opportunity, "", ""); !errors.Is(err, ErrLostReasonRequired) {
			t.Fatalf("expected ErrLostReasonRequired, got %v", err)
		}
	})

	t.Run("TEST-INV-TERMINAL-001 terminal records are read-only", func(t *testing.T) {
		if err := EnsureOpenForMutation(Opportunity{Stage: StageWon}); !errors.Is(err, ErrTerminalRecordReadOnly) {
			t.Fatalf("expected terminal read-only for Won, got %v", err)
		}
		if err := EnsureOpenForMutation(Opportunity{Stage: StageLost}); !errors.Is(err, ErrTerminalRecordReadOnly) {
			t.Fatalf("expected terminal read-only for Lost, got %v", err)
		}
	})
}
