package handler

import (
	"net/http"
	"net/url"
	"testing"
)

func TestAccountArchiveAcceptance(t *testing.T) {
	db := newAccountTestDB(t)
	app := NewAccountServer(db, Config{})

	t.Run("TEST-ARCHIVE-001/002 Sales denied and manager can archive with explicit archived filter", func(t *testing.T) {
		create := postAccountJSON(app, "/accounts", map[string]any{
			"companyName":    "Archive Eligible Account",
			"customerStatus": "Active",
			"ownerId":        "sales-1",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if create.Code != http.StatusCreated {
			t.Fatalf("expected create 201, got %d body=%s", create.Code, create.Body.String())
		}
		accountID := decodeJSON(t, create)["id"].(string)

		salesArchive := postAccountJSON(app, "/accounts/"+accountID+"/archive", map[string]any{
			"expectedVersion": 1,
			"reason":          "Duplicate retired from active operations",
		}, actorHeaders("sales-1", "Sales"))
		if salesArchive.Code != http.StatusForbidden {
			t.Fatalf("expected sales archive 403, got %d body=%s", salesArchive.Code, salesArchive.Body.String())
		}
		requireErrorCode(t, salesArchive, "PERMISSION_DENIED")

		eligibility := getAccountJSON(app, "/accounts/"+accountID+"/archive-eligibility", actorHeaders("mgr-1", "Sales Manager"))
		if eligibility.Code != http.StatusOK {
			t.Fatalf("expected eligibility 200, got %d body=%s", eligibility.Code, eligibility.Body.String())
		}
		if decodeJSON(t, eligibility)["canArchive"] != true {
			t.Fatalf("expected canArchive true, got body=%s", eligibility.Body.String())
		}

		archive := postAccountJSON(app, "/accounts/"+accountID+"/archive", map[string]any{
			"expectedVersion": 1,
			"reason":          "No active obligations remain",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if archive.Code != http.StatusOK {
			t.Fatalf("expected archive 200, got %d body=%s", archive.Code, archive.Body.String())
		}
		body := decodeJSON(t, archive)
		if body["archived"] != true {
			t.Fatalf("expected archived true, got body=%s", archive.Body.String())
		}
		requireEvent(t, db, "AccountArchived", accountID)

		activeList := getAccountJSON(app, "/accounts?search="+url.QueryEscape("Archive Eligible"), actorHeaders("mgr-1", "Sales Manager"))
		if activeList.Code != http.StatusOK || contains(activeList.Body.String(), accountID) {
			t.Fatalf("archived record must be hidden from active list, got %d body=%s", activeList.Code, activeList.Body.String())
		}
		archivedList := getAccountJSON(app, "/accounts?search="+url.QueryEscape("Archive Eligible")+"&includeArchived=true", actorHeaders("mgr-1", "Sales Manager"))
		if archivedList.Code != http.StatusOK || !contains(archivedList.Body.String(), accountID) {
			t.Fatalf("archived record must be retrievable with explicit filter, got %d body=%s", archivedList.Code, archivedList.Body.String())
		}
	})
}
