package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGatewayRoutesAuthorizedListDetailAndPropagatesCorrelation(t *testing.T) {
	identity := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/current" {
			t.Fatalf("unexpected identity path %s", r.URL.Path)
		}
		if r.Header.Get("X-Correlation-Id") != "corr-test-1" {
			t.Fatalf("identity missing propagated correlation id: %q", r.Header.Get("X-Correlation-Id"))
		}
		writeTestJSON(t, w, http.StatusOK, map[string]any{
			"user": map[string]any{
				"id":     "sales-1",
				"role":   "Sales",
				"status": "Active",
			},
		})
	}))
	defer identity.Close()

	lead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Correlation-Id") != "corr-test-1" {
			t.Fatalf("lead missing propagated correlation id: %q", r.Header.Get("X-Correlation-Id"))
		}
		if r.Header.Get("X-Actor-User-Id") != "sales-1" || r.Header.Get("X-Actor-Role") != "Sales" {
			t.Fatalf("lead missing actor context: user=%q role=%q", r.Header.Get("X-Actor-User-Id"), r.Header.Get("X-Actor-Role"))
		}
		switch r.URL.Path {
		case "/leads":
			writeTestJSON(t, w, http.StatusOK, map[string]any{"items": []map[string]any{{"id": "lead-1", "name": "Authorized Lead"}}})
		case "/leads/lead-1":
			writeTestJSON(t, w, http.StatusOK, map[string]any{"id": "lead-1", "name": "Authorized Lead"})
		default:
			t.Fatalf("unexpected lead path %s", r.URL.Path)
		}
	}))
	defer lead.Close()

	app := NewGatewayServer(Config{IdentityBaseURL: identity.URL, Routes: map[string]string{"leads": lead.URL}})

	list := gatewayGet(app, "/api/leads", "corr-test-1")
	if list.Code != http.StatusOK {
		t.Fatalf("TEST-NAV-RETRIEVE-001 expected list 200, got %d body=%s", list.Code, list.Body.String())
	}
	requireJSONValue(t, decodeJSON(t, list), "correlationId", "corr-test-1")

	detail := gatewayGet(app, "/api/leads/lead-1", "corr-test-1")
	if detail.Code != http.StatusOK {
		t.Fatalf("TEST-NAV-RETRIEVE-001 expected detail 200, got %d body=%s", detail.Code, detail.Body.String())
	}
	requireJSONValue(t, decodeJSON(t, detail), "correlationId", "corr-test-1")
}

func writeTestJSON(t *testing.T, w http.ResponseWriter, status int, body any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Fatalf("write json: %v", err)
	}
}
