package domain

import "testing"

func TestStageMachineAcceptance(t *testing.T) {
	t.Run("TEST-OPP-STAGE-001 allows forward transition", func(t *testing.T) {
		next, err := AdvanceStage(StageNewOpportunity)
		if err != nil {
			t.Fatalf("expected transition allowed: %v", err)
		}
		if next != StageNeedsConfirmed {
			t.Fatalf("expected %q, got %q", StageNeedsConfirmed, next)
		}
	})

	t.Run("TEST-OPP-STAGE-002 rejects forbidden skip", func(t *testing.T) {
		if err := ValidateStageTransition(StageNewOpportunity, StageQuote); err != ErrInvalidTransition {
			t.Fatalf("expected invalid transition, got %v", err)
		}
	})

	t.Run("TEST-OPP-STAGE-003 rejects arbitrary rollback", func(t *testing.T) {
		if err := ValidateStageTransition(StageQuote, StageNeedsConfirmed); err != ErrInvalidTransition {
			t.Fatalf("expected rollback invalid transition, got %v", err)
		}
	})
}
