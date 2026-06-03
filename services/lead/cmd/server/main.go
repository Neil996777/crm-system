package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"crm-system/services/lead/internal/event"
	"crm-system/services/lead/internal/handler"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const defaultServiceName = "lead"

type HealthResponse struct {
	Service   string `json:"service"`
	Process   string `json:"process"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

func healthHandler(requiresDatabase bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := HealthResponse{Service: envOrDefault("SERVICE_NAME", defaultServiceName), Process: "up", Database: "not_applicable", Timestamp: time.Now().UTC().Format(time.RFC3339)}
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
	leadServer := handler.NewLeadServer(db, handler.Config{
		AccountServiceURL:      envOrDefault("ACCOUNT_SERVICE_URL", "http://account:8080"),
		OpportunityServiceURL:  envOrDefault("OPPORTUNITY_SERVICE_URL", "http://opportunity:8080"),
		AuditHistoryServiceURL: envOrDefault("AUDIT_HISTORY_SERVICE_URL", "http://audit-history:8080"),
		ServiceID:              envOrDefault("SERVICE_ID", "lead"),
		ServiceTokenSecret:     []byte(os.Getenv("SERVICE_TOKEN_SECRET")),
	})
	startReportingDispatcher(db)
	mux.Handle("/leads", leadServer)
	mux.Handle("/leads/", leadServer)
	mux.Handle("/duplicate-checks", leadServer)
	addr := ":" + envOrDefault("PORT", "8080")
	log.Printf("%s listening on %s", defaultServiceName, addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func startReportingDispatcher(db *sql.DB) {
	reportingURL := os.Getenv("REPORTING_SERVICE_URL")
	secret := []byte(os.Getenv("SERVICE_TOKEN_SECRET"))
	if reportingURL == "" || len(secret) == 0 {
		return
	}
	outbox := event.NewOutbox(db)
	config := event.DispatchConfig{
		ServiceID:           envOrDefault("SERVICE_ID", "lead"),
		ServiceTokenSecret:  secret,
		ReportingServiceURL: reportingURL,
	}
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := outbox.DispatchOnce(ctx, config); err != nil {
				log.Printf("reporting outbox dispatch: %v", err)
			}
			cancel()
			<-ticker.C
		}
	}()
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
