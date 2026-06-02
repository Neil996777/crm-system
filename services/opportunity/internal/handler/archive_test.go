package handler

import (
	"net/http"
	"testing"
)

func TestOpportunityArchiveAcceptance(t *testing.T) {
	db := newOpportunityTestDB(t)
	app := NewOpportunityServer(db, Config{})

	t.Run("TEST-ARCHIVE-001/002 opportunity Sales denied and manager archives with explicit archived filter", func(t *testing.T) {
		create := postOpportunityJSON(app, "/opportunities", map[string]any{
			"customerId":        "acct_archive_opp",
			"ownerId":           "sales-1",
			"stage":             "New Opportunity",
			"expectedAmount":    "10000.00",
			"expectedCloseDate": "2026-12-30",
			"title":             "Archive Opportunity",
		}, actorHeaders("sales-1", "Sales"))
		if create.Code != http.StatusCreated {
			t.Fatalf("expected create 201, got %d body=%s", create.Code, create.Body.String())
		}
		opportunityID := decodeJSON(t, create)["id"].(string)

		salesArchive := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/archive", map[string]any{
			"expectedVersion": 1,
			"reason":          "No longer active",
		}, actorHeaders("sales-1", "Sales"))
		if salesArchive.Code != http.StatusForbidden {
			t.Fatalf("expected sales archive 403, got %d body=%s", salesArchive.Code, salesArchive.Body.String())
		}
		requireErrorCode(t, salesArchive, "PERMISSION_DENIED")

		archive := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/archive", map[string]any{
			"expectedVersion": 1,
			"reason":          "No active obligations remain",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if archive.Code != http.StatusOK {
			t.Fatalf("expected archive 200, got %d body=%s", archive.Code, archive.Body.String())
		}
		if decodeJSON(t, archive)["archived"] != true {
			t.Fatalf("expected archived true, got body=%s", archive.Body.String())
		}
		requireEvent(t, db, "OpportunityArchived", opportunityID)

		active := getOpportunityJSON(app, "/opportunities?search=Archive%20Opportunity", actorHeaders("mgr-1", "Sales Manager"))
		if active.Code != http.StatusOK || contains(active.Body.String(), opportunityID) {
			t.Fatalf("archived opportunity must be hidden from active list, got %d body=%s", active.Code, active.Body.String())
		}
		archived := getOpportunityJSON(app, "/opportunities?search=Archive%20Opportunity&includeArchived=true", actorHeaders("mgr-1", "Sales Manager"))
		if archived.Code != http.StatusOK || !contains(archived.Body.String(), opportunityID) {
			t.Fatalf("archived opportunity must be visible with explicit filter, got %d body=%s", archived.Code, archived.Body.String())
		}
	})
}
