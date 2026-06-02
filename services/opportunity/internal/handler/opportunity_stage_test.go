package handler

import (
	"net/http"
	"testing"
)

func TestOpportunityStageTransitionAcceptance(t *testing.T) {
	db := newOpportunityTestDB(t)
	app := NewOpportunityServer(db, Config{})

	t.Run("TEST-OPP-STAGE-001 and TEST-HISTORY-002 allowed forward transition persists and emits history", func(t *testing.T) {
		opportunityID := createOpportunityForStage(t, app, "stage-acct-1", "sales-1")
		rec := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/stage", map[string]any{
			"expectedVersion": 1,
			"toStage":         "Needs Confirmed",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "stage", "Needs Confirmed")
		if body["version"].(float64) != 2 {
			t.Fatalf("expected version 2, got %#v", body["version"])
		}
		requireEvent(t, db, "OpportunityStageChanged", opportunityID)
	})

	t.Run("TEST-OPP-STAGE-002 forbidden skip rejected without mutation", func(t *testing.T) {
		opportunityID := createOpportunityForStage(t, app, "stage-acct-2", "sales-1")
		rec := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/stage", map[string]any{
			"expectedVersion": 1,
			"toStage":         "Quote",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "INVALID_TRANSITION")
		fetch := getOpportunityJSON(app, "/opportunities/"+opportunityID, actorHeaders("sales-1", "Sales"))
		requireJSONValue(t, decodeJSON(t, fetch), "stage", "New Opportunity")
		requireNoEvent(t, db, "OpportunityStageChanged", opportunityID)
	})

	t.Run("TEST-OPP-STAGE-003 rollback rejected without mutation", func(t *testing.T) {
		opportunityID := createOpportunityForStage(t, app, "stage-acct-3", "sales-1")
		first := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/stage", map[string]any{
			"expectedVersion": 1,
			"toStage":         "Needs Confirmed",
		}, actorHeaders("sales-1", "Sales"))
		if first.Code != http.StatusOK {
			t.Fatalf("expected initial transition 200, got %d body=%s", first.Code, first.Body.String())
		}
		rollback := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/stage", map[string]any{
			"expectedVersion": 2,
			"toStage":         "New Opportunity",
		}, actorHeaders("sales-1", "Sales"))
		if rollback.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", rollback.Code, rollback.Body.String())
		}
		requireErrorCode(t, rollback, "INVALID_TRANSITION")
		fetch := getOpportunityJSON(app, "/opportunities/"+opportunityID, actorHeaders("sales-1", "Sales"))
		requireJSONValue(t, decodeJSON(t, fetch), "stage", "Needs Confirmed")
	})

	t.Run("TEST-ABUSE-BRBYPASS-001 non-owned Sales transition denied", func(t *testing.T) {
		opportunityID := createOpportunityForStage(t, app, "stage-acct-4", "sales-2")
		rec := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/stage", map[string]any{
			"expectedVersion": 1,
			"toStage":         "Needs Confirmed",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "PERMISSION_DENIED")
	})
}

func createOpportunityForStage(t *testing.T, app http.Handler, customerID, ownerID string) string {
	t.Helper()
	create := postOpportunityJSON(app, "/opportunities", map[string]any{
		"customerId":        customerID,
		"ownerId":           ownerID,
		"stage":             "New Opportunity",
		"expectedAmount":    "10000.00",
		"expectedCloseDate": "2026-10-30",
		"title":             "Stage test",
	}, actorHeaders("mgr-1", "Sales Manager"))
	if create.Code != http.StatusCreated {
		t.Fatalf("expected opportunity create 201, got %d body=%s", create.Code, create.Body.String())
	}
	return decodeJSON(t, create)["id"].(string)
}
