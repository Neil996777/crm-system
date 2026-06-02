package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpointRequiresOwnedDatabaseConfiguration(t *testing.T) {
	t.Setenv("SERVICE_NAME", "reporting")
	t.Setenv("DATABASE_URL", "")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	healthHandler(true).ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503 when database is not configured, got %d body=%s", rec.Code, rec.Body.String())
	}
}
