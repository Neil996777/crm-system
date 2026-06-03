package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestCSVExportRunAcceptance(t *testing.T) {
	db := newImportExportTestDB(t)
	var queryCount int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/leads" {
			http.NotFound(w, r)
			return
		}
		atomic.AddInt32(&queryCount, 1)
		if r.URL.Query().Get("includeArchived") == "true" {
			t.Fatalf("default export requested archived records")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"items": []map[string]any{
			{"id": "lead-1", "companyName": "Export Team Co", "leadName": "Team lead", "source": "Website", "ownerId": "sales-1", "status": "Pending Qualification"},
			{"id": "lead-2", "companyName": "=Unsafe Co", "leadName": "Unsafe lead", "source": "Website", "ownerId": "sales-1", "status": "Pending Qualification"},
		}})
	}))
	t.Cleanup(target.Close)
	app := NewImportExportServer(db, Config{LeadServiceURL: target.URL, HTTPClient: target.Client()})

	t.Run("TEST-ABUSE-EXPORTCONFIRM-001 confirmation required before query", func(t *testing.T) {
		rec := postImportJSON(app, "/exports", map[string]any{"objectType": "lead", "confirmed": false}, actorHeaders("mgr-1", "Sales Manager", "single-team"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected confirmation 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "EXPORT_CONFIRMATION_REQUIRED")
		if atomic.LoadInt32(&queryCount) != 0 {
			t.Fatalf("unconfirmed export queried target service")
		}
	})

	t.Run("TEST-CSV-EXPORT-001 and TEST-ABUSE-EXPORTLEAK-001 authorized active records exported with formula safety", func(t *testing.T) {
		rec := postImportJSON(app, "/exports", map[string]any{"objectType": "lead", "confirmed": true}, actorHeaders("mgr-1", "Sales Manager", "single-team"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected export 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		if body["status"] != "Completed" || body["exportedCount"].(float64) != 2 || body["archivedIncluded"].(bool) {
			t.Fatalf("expected active export metadata, got %#v", body)
		}
		csvContent := body["content"].(string)
		if !strings.Contains(csvContent, "Export Team Co") || !strings.Contains(csvContent, "'=Unsafe Co") {
			t.Fatalf("expected escaped CSV content, got %q", csvContent)
		}
	})

	t.Run("TEST-CSV-EXPORT-002 Sales denied", func(t *testing.T) {
		rec := postImportJSON(app, "/exports", map[string]any{"objectType": "lead", "confirmed": true}, actorHeaders("sales-1", "Sales", "single-team"))
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected sales denied 403, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "PERMISSION_DENIED")
	})
}
