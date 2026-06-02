package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpointReportsProcessUpWithoutDatabase(t *testing.T) {
	t.Setenv("SERVICE_NAME", "gateway-bff")
	t.Setenv("DATABASE_URL", "should-not-be-used")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	healthHandler(false).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	var body HealthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid health response: %v", err)
	}
	if body.Service != "gateway-bff" || body.Process != "up" || body.Database != "not_applicable" {
		t.Fatalf("unexpected health body: %+v", body)
	}
}
