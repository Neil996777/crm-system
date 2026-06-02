package handler

import (
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode json body %q: %v", rec.Body.String(), err)
	}
	return body
}

func requireJSONValue(t *testing.T, body map[string]any, path, expected string) {
	t.Helper()
	value, ok := body[path]
	if !ok {
		t.Fatalf("missing json path %s in %#v", path, body)
	}
	if value != expected {
		t.Fatalf("expected %s=%q, got %#v in %#v", path, expected, value, body)
	}
}

func requireErrorCode(t *testing.T, rec *httptest.ResponseRecorder, expected string) {
	t.Helper()
	body := decodeJSON(t, rec)
	errBody, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("missing error body: %#v", body)
	}
	code, ok := errBody["code"].(string)
	if !ok || code != expected {
		t.Fatalf("expected error code %q, got %#v body=%s", expected, errBody["code"], rec.Body.String())
	}
}

func requireEvent(t *testing.T, db queryer, eventType, aggregateID string) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT count(*) FROM commercial.outbox_events WHERE event_type = $1 AND aggregate_id = $2`, eventType, aggregateID).Scan(&count); err != nil {
		t.Fatalf("count event %s: %v", eventType, err)
	}
	if count == 0 {
		t.Fatalf("expected outbox event %s aggregate=%s", eventType, aggregateID)
	}
}

type queryer interface {
	QueryRow(query string, args ...any) *sql.Row
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
