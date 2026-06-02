package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGatewayNormalizesSafeErrorsAndDoesNotLeakUnauthorizedDetail(t *testing.T) {
	identity := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeTestJSON(t, w, http.StatusOK, map[string]any{"user": map[string]any{"id": "sales-1", "role": "Sales", "status": "Active"}})
	}))
	defer identity.Close()

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RawQuery {
		case "status=bad-filter":
			writeTestJSON(t, w, http.StatusBadRequest, map[string]any{
				"error": map[string]any{
					"code":        "INVALID_FILTER",
					"category":    "validation",
					"safeMessage": "The filter is invalid.",
					"fieldErrors": []map[string]any{{"field": "status", "code": "UNKNOWN_VALUE", "safeMessage": "Choose a supported status."}},
				},
			})
		default:
			writeTestJSON(t, w, http.StatusForbidden, map[string]any{
				"error": map[string]any{
					"code":        "PERMISSION_DENIED",
					"category":    "permission",
					"safeMessage": "Permission denied.",
					"internal":    "restricted customer ACME amount 1000000",
				},
			})
		}
	}))
	defer target.Close()

	app := NewGatewayServer(Config{IdentityBaseURL: identity.URL, Routes: map[string]string{"leads": target.URL}})

	invalid := gatewayGet(app, "/api/leads?status=bad-filter", "corr-filter")
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("TEST-NAV-RETRIEVE-004 expected invalid filter 400, got %d body=%s", invalid.Code, invalid.Body.String())
	}
	body := decodeJSON(t, invalid)
	requireJSONValue(t, body, "correlationId", "corr-filter")
	requireNestedJSONValue(t, body, "error", "code", "INVALID_FILTER")

	denied := gatewayGet(app, "/api/leads?status=secret", "corr-denied")
	if denied.Code != http.StatusForbidden {
		t.Fatalf("TEST-NAV-RETRIEVE-005 expected permission 403, got %d body=%s", denied.Code, denied.Body.String())
	}
	if strings.Contains(denied.Body.String(), "ACME") || strings.Contains(denied.Body.String(), "1000000") {
		t.Fatalf("gateway leaked unauthorized detail: %s", denied.Body.String())
	}
	requireNestedJSONValue(t, decodeJSON(t, denied), "error", "code", "PERMISSION_DENIED")
}
