package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"crm-system/services/identity-authz/internal/handler"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const defaultServiceName = "identity-authz"

type HealthResponse struct {
	Service   string `json:"service"`
	Process   string `json:"process"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

func healthHandler(requiresDatabase bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := HealthResponse{
			Service:   envOrDefault("SERVICE_NAME", defaultServiceName),
			Process:   "up",
			Database:  "not_applicable",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		if requiresDatabase {
			if err := pingDatabase(); err != nil {
				response.Database = "down"
				writeJSON(w, http.StatusServiceUnavailable, response)
				return
			}
			response.Database = "up"
		}
		writeJSON(w, http.StatusOK, response)
	})
}

func main() {
	db, err := openDatabase()
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.Handle("/health", healthHandler(true))
	authServer := handler.NewAuthServer(db, handler.Config{
		CookieSecure:           envBoolOrDefault("COOKIE_SECURE", true),
		SessionTTL:             12 * time.Hour,
		IdleSessionTTL:         30 * time.Minute,
		AuditHistoryServiceURL: envOrDefault("AUDIT_HISTORY_SERVICE_URL", "http://audit-history:8080"),
		ServiceID:              envOrDefault("SERVICE_ID", "identity-authz"),
		ServiceTokenSecret:     []byte(os.Getenv("SERVICE_TOKEN_SECRET")),
	})
	mux.Handle("/auth/", authServer)
	mux.Handle("/internal/sessions/check", authServer)
	mux.Handle("/internal/permissions/check", authServer)
	mux.Handle("/admin/", authServer)
	addr := ":" + envOrDefault("PORT", "8080")
	log.Printf("%s listening on %s", defaultServiceName, addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func openDatabase() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, sql.ErrConnDone
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func pingDatabase() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return sql.ErrConnDone
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return db.PingContext(ctx)
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

func envBoolOrDefault(name string, fallback bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
