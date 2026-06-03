package handler

import (
	"net/http"
	"strings"
	"testing"
)

func TestAuditRetentionAndClassificationAcceptance(t *testing.T) {
	db, _ := newAuditTestDB(t)
	app := NewAuditServer(db, Config{ServiceTokenSecret: []byte("audit-secret")})
	token := signTestServiceToken(t, []byte("audit-secret"), "lead", "audit-history", "audit.append")

	t.Run("TEST-RETENTION-001 append-only event carries minimum retention metadata", func(t *testing.T) {
		rec := appendEvent(t, app, token, map[string]any{
			"eventId":            "EVT-RETENTION-001",
			"eventVersion":       1,
			"surfaces":           []string{"record_history"},
			"action":             "archive",
			"resourceType":       "lead",
			"resourceId":         "lead-retention-1",
			"result":             "success",
			"beforeSummary":      map[string]any{"status": "Pending Qualification"},
			"afterSummary":       map[string]any{"status": "Archived"},
			"diffClassification": "Confidential",
			"safeSummary":        "Lead archived.",
			"acceptanceIds":      []string{"ACC-014", "TEST-RETENTION-001"},
		}, actorHeaders("mgr-1", "Sales Manager", "Manager One"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected append 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		if body["classification"] != "Confidential" || body["retentionPolicy"] == "" || body["retainUntil"] == "" {
			t.Fatalf("expected classification and retention metadata, got %#v", body)
		}
		deleted := httptestDelete(app, "/records/lead/lead-retention-1/history/"+body["eventUid"].(string), nil)
		if deleted.Code != http.StatusMethodNotAllowed && deleted.Code != http.StatusNotFound {
			t.Fatalf("expected no hard-delete path, got %d body=%s", deleted.Code, deleted.Body.String())
		}
	})

	t.Run("TEST-RETENTION-002 restricted before/after values masked in safe response", func(t *testing.T) {
		rec := appendEvent(t, app, token, map[string]any{
			"eventId":            "EVT-RETENTION-002",
			"eventVersion":       1,
			"surfaces":           []string{"operation_log"},
			"action":             "payment_recorded",
			"resourceType":       "payment",
			"resourceId":         "payment-retention-1",
			"result":             "success",
			"beforeSummary":      map[string]any{"amount": "99999.00", "note": "restricted note"},
			"afterSummary":       map[string]any{"amount": "99999.00", "note": "restricted note updated"},
			"diffClassification": "Restricted",
			"safeSummary":        "Payment recorded.",
			"acceptanceIds":      []string{"ACC-014", "TEST-RETENTION-002"},
		}, actorHeaders("admin-1", "Administrator", "Admin One"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected append 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		if bodyContains(rec.Body.String(), "99999.00") || bodyContains(rec.Body.String(), "restricted note") {
			t.Fatalf("restricted values leaked in safe append response: %s", rec.Body.String())
		}
	})
}

func bodyContains(body string, needle string) bool {
	return len(needle) > 0 && len(body) >= len(needle) && strings.Contains(body, needle)
}
