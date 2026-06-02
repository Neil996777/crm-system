package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func gatewayGet(handler http.Handler, path, correlationID string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("X-Correlation-Id", correlationID)
	req.AddCookie(&http.Cookie{Name: "crm_session", Value: "test-session"})
	handler.ServeHTTP(rec, req)
	return rec
}

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

func requireNestedJSONValue(t *testing.T, body map[string]any, parent, child, expected string) {
	t.Helper()
	nested, ok := body[parent].(map[string]any)
	if !ok {
		t.Fatalf("missing nested json object %s in %#v", parent, body)
	}
	value, ok := nested[child]
	if !ok {
		t.Fatalf("missing json path %s.%s in %#v", parent, child, body)
	}
	if value != expected {
		t.Fatalf("expected %s.%s=%q, got %#v in %#v", parent, child, expected, value, body)
	}
}
