package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"crm-system/services/gateway-bff/internal/handler"
)

type HealthResponse struct {
	Service   string `json:"service"`
	Process   string `json:"process"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

func healthHandler(requiresDatabase bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if requiresDatabase {
			http.Error(w, "gateway-bff must not require database access", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, HealthResponse{
			Service:   envOrDefault("SERVICE_NAME", "gateway-bff"),
			Process:   "up",
			Database:  "not_applicable",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/health", healthHandler(false))
	gateway := handler.NewGatewayServer(handler.Config{
		IdentityBaseURL: envOrDefault("IDENTITY_AUTHZ_URL", "http://identity-authz:8080"),
		Routes: map[string]string{
			"leads":         envOrDefault("LEAD_SERVICE_URL", "http://lead:8080"),
			"accounts":      envOrDefault("ACCOUNT_SERVICE_URL", "http://account:8080"),
			"contacts":      envOrDefault("ACCOUNT_SERVICE_URL", "http://account:8080"),
			"opportunities": envOrDefault("OPPORTUNITY_SERVICE_URL", "http://opportunity:8080"),
			"quotes":        envOrDefault("COMMERCIAL_SERVICE_URL", "http://commercial:8080"),
			"contracts":     envOrDefault("COMMERCIAL_SERVICE_URL", "http://commercial:8080"),
			"activities":    envOrDefault("WORK_SERVICE_URL", "http://work:8080"),
			"notes":         envOrDefault("WORK_SERVICE_URL", "http://work:8080"),
			"tasks":         envOrDefault("WORK_SERVICE_URL", "http://work:8080"),
			"reminders":     envOrDefault("WORK_SERVICE_URL", "http://work:8080"),
			"history":       envOrDefault("AUDIT_HISTORY_SERVICE_URL", "http://audit-history:8080"),
		},
	})
	mux.Handle("/auth/", gateway)
	mux.Handle("/admin/", gateway)
	mux.Handle("/api/", gateway)
	addr := ":" + envOrDefault("PORT", "8080")
	log.Printf("gateway-bff listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write response: %v", err)
	}
}

func envOrDefault(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}
