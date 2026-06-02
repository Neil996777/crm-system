package handler

import (
	"net/http"
	"testing"
)

func TestOpportunityTerminalLockAcceptance(t *testing.T) {
	db := newOpportunityTestDB(t)
	app := NewOpportunityServer(db, Config{})

	t.Run("TEST-OPP-CLOSE-005 and TEST-INV-TERMINAL-001 reject reopen rollback re-close", func(t *testing.T) {
		opportunityID := createOpportunityForStage(t, app, "acct_terminal_lost", "sales-1")
		closeLost := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/close-lost", map[string]any{
			"expectedVersion": 1,
			"closeDate":       "2027-04-01",
			"lostReason": map[string]any{
				"code": "NO_DECISION",
			},
		}, actorHeaders("sales-1", "Sales"))
		if closeLost.Code != http.StatusOK {
			t.Fatalf("expected close Lost 200, got %d body=%s", closeLost.Code, closeLost.Body.String())
		}

		stage := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/stage", map[string]any{
			"expectedVersion": 2,
			"toStage":         "New Opportunity",
		}, actorHeaders("sales-1", "Sales"))
		if stage.Code != http.StatusBadRequest {
			t.Fatalf("expected terminal stage change 400, got %d body=%s", stage.Code, stage.Body.String())
		}
		requireErrorCode(t, stage, "TERMINAL_RECORD_READ_ONLY")

		edit := patchOpportunityJSON(app, "/opportunities/"+opportunityID, map[string]any{
			"expectedVersion":   2,
			"customerId":        "acct_terminal_lost",
			"ownerId":           "sales-1",
			"stage":             "Needs Confirmed",
			"expectedAmount":    "11000.00",
			"expectedCloseDate": "2027-04-30",
			"title":             "Terminal edit attempt",
		}, actorHeaders("sales-1", "Sales"))
		if edit.Code != http.StatusBadRequest {
			t.Fatalf("expected terminal edit 400, got %d body=%s", edit.Code, edit.Body.String())
		}
		requireErrorCode(t, edit, "TERMINAL_RECORD_READ_ONLY")

		reclose := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/close-lost", map[string]any{
			"expectedVersion": 2,
			"closeDate":       "2027-04-02",
			"lostReason": map[string]any{
				"code": "OTHER",
			},
		}, actorHeaders("sales-1", "Sales"))
		if reclose.Code != http.StatusBadRequest {
			t.Fatalf("expected terminal re-close 400, got %d body=%s", reclose.Code, reclose.Body.String())
		}
		requireErrorCode(t, reclose, "TERMINAL_RECORD_READ_ONLY")
	})

	t.Run("TEST-OPP-CLOSE-006 rejects post-close stage edit and preserves record", func(t *testing.T) {
		opportunityID := createOpportunityForStage(t, app, "acct_terminal_preserve", "sales-1")
		closeLost := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/close-lost", map[string]any{
			"expectedVersion": 1,
			"closeDate":       "2027-04-01",
			"lostReason": map[string]any{
				"code": "COMPETITOR",
			},
		}, actorHeaders("sales-1", "Sales"))
		if closeLost.Code != http.StatusOK {
			t.Fatalf("expected close Lost 200, got %d body=%s", closeLost.Code, closeLost.Body.String())
		}
		fetch := getOpportunityJSON(app, "/opportunities/"+opportunityID, actorHeaders("sales-1", "Sales"))
		body := decodeJSON(t, fetch)
		requireJSONValue(t, body, "stage", "Lost")
		requireJSONValue(t, body, "lostReasonCode", "COMPETITOR")
	})
}
