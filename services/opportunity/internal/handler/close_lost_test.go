package handler

import (
	"net/http"
	"testing"
)

func TestOpportunityCloseLostAcceptance(t *testing.T) {
	db := newOpportunityTestDB(t)
	app := NewOpportunityServer(db, Config{})

	t.Run("TEST-OPP-CLOSE-004 rejects Lost without reason", func(t *testing.T) {
		opportunityID := createOpportunityForStage(t, app, "acct_close_lost_missing", "sales-1")
		rec := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/close-lost", map[string]any{
			"expectedVersion": 1,
			"closeDate":       "2027-04-01",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected missing lost reason 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "LOST_REASON_REQUIRED")
		requireNoEvent(t, db, "OpportunityClosedLost", opportunityID)
	})

	t.Run("TEST-OPP-CLOSE-003 closes Lost with reason", func(t *testing.T) {
		opportunityID := createOpportunityForStage(t, app, "acct_close_lost_reason", "sales-1")
		rec := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/close-lost", map[string]any{
			"expectedVersion": 1,
			"closeDate":       "2027-04-01",
			"lostReason": map[string]any{
				"code":   "PRICE",
				"detail": "Competitor pricing won",
			},
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected close Lost 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "opportunityId", opportunityID)
		requireJSONValue(t, body, "status", "Lost")
		fetch := getOpportunityJSON(app, "/opportunities/"+opportunityID, actorHeaders("sales-1", "Sales"))
		fetchBody := decodeJSON(t, fetch)
		requireJSONValue(t, fetchBody, "stage", "Lost")
		requireJSONValue(t, fetchBody, "lostReasonCode", "PRICE")
		requireEvent(t, db, "OpportunityClosedLost", opportunityID)
	})
}
