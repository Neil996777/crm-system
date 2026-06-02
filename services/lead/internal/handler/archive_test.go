package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLeadArchiveAcceptance(t *testing.T) {
	db := newLeadTestDB(t)
	app := NewLeadServer(db, Config{})

	t.Run("TEST-ARCHIVE-001/002 lead Sales denied and manager archives with explicit archived filter", func(t *testing.T) {
		create := postLeadJSON(app, "/leads", map[string]any{
			"leadName":    "Archive lead",
			"companyName": "Archive Lead Co",
			"source":      "Website",
			"ownerId":     "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if create.Code != http.StatusCreated {
			t.Fatalf("expected create 201, got %d body=%s", create.Code, create.Body.String())
		}
		leadID := decodeJSON(t, create)["id"].(string)

		salesArchive := postLeadJSON(app, "/leads/"+leadID+"/archive", map[string]any{
			"expectedVersion": 1,
			"reason":          "No longer active",
		}, actorHeaders("sales-1", "Sales"))
		if salesArchive.Code != http.StatusForbidden {
			t.Fatalf("expected sales archive 403, got %d body=%s", salesArchive.Code, salesArchive.Body.String())
		}
		requireErrorCode(t, salesArchive, "PERMISSION_DENIED")

		archive := postLeadJSON(app, "/leads/"+leadID+"/archive", map[string]any{
			"expectedVersion": 1,
			"reason":          "No active obligations remain",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if archive.Code != http.StatusOK {
			t.Fatalf("expected archive 200, got %d body=%s", archive.Code, archive.Body.String())
		}
		if decodeJSON(t, archive)["archived"] != true {
			t.Fatalf("expected archived true, got body=%s", archive.Body.String())
		}
		requireEvent(t, db, "LeadArchived", leadID)

		active := getLeadJSON(app, "/leads?search=Archive%20Lead", actorHeaders("mgr-1", "Sales Manager"))
		if active.Code != http.StatusOK || contains(active.Body.String(), leadID) {
			t.Fatalf("archived lead must be hidden from active list, got %d body=%s", active.Code, active.Body.String())
		}
		archived := getLeadJSON(app, "/leads?search=Archive%20Lead&includeArchived=true", actorHeaders("mgr-1", "Sales Manager"))
		if archived.Code != http.StatusOK || !contains(archived.Body.String(), leadID) {
			t.Fatalf("archived lead must be visible with explicit filter, got %d body=%s", archived.Code, archived.Body.String())
		}
	})
}

func getLeadJSON(handler http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	handler.ServeHTTP(rec, req)
	return rec
}
