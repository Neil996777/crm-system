package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLeadQueryScopeAndIDOR(t *testing.T) {
	db := newLeadTestDB(t)
	app := NewLeadServer(db, Config{})
	owned := postLeadJSON(app, "/leads", map[string]any{
		"leadName": "Owned lead",
		"source":   "Website",
		"ownerId":  "sales-1",
	}, actorHeaders("sales-1", "Sales"))
	ownedID := decodeJSON(t, owned)["id"].(string)
	other := postLeadJSON(app, "/leads", map[string]any{
		"leadName": "Restricted lead",
		"source":   "Website",
		"ownerId":  "sales-2",
	}, actorHeaders("sales-2", "Sales"))
	otherID := decodeJSON(t, other)["id"].(string)

	t.Run("TEST-NAV-RETRIEVE-002 search/filter scoped to owned leads", func(t *testing.T) {
		rec := getLead(app, "/leads?search=lead", actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if !contains(body, ownedID) {
			t.Fatalf("expected owned lead in list: %s", body)
		}
		if contains(body, otherID) || contains(body, "Restricted lead") {
			t.Fatalf("TEST-AUTHZ-SCOPE-004 leaked non-owned lead: %s", body)
		}
	})

	t.Run("TEST-ABUSE-IDOR-001 Sales non-owned detail denied without leakage", func(t *testing.T) {
		rec := getLead(app, "/leads/"+otherID, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected safe 404, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "NOT_FOUND")
		if contains(rec.Body.String(), "Restricted lead") || contains(rec.Body.String(), otherID) {
			t.Fatalf("denial leaked restricted lead detail: %s", rec.Body.String())
		}
	})
}

func getLead(handler http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	handler.ServeHTTP(rec, req)
	return rec
}
