package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"crm-system/services/commercial/internal/event"
	"crm-system/services/commercial/internal/handler"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const defaultServiceName = "commercial"

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
	mux := http.NewServeMux()
	mux.Handle("/health", healthHandler(true))
	db, err := openDatabase()
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()
	commercialServer := handler.NewCommercialServer(db, handler.Config{
		ServiceID:          envOrDefault("SERVICE_ID", "commercial"),
		ServiceTokenSecret: []byte(os.Getenv("SERVICE_TOKEN_SECRET")),
	})
	startAuditDispatcher(db)
	mux.Handle("/quotes", commercialServer)
	mux.Handle("/quotes/", commercialServer)
	mux.Handle("/contracts", commercialServer)
	mux.Handle("/contracts/", commercialServer)
	mux.Handle("/internal/", commercialServer)
	addr := ":" + envOrDefault("PORT", "8080")
	log.Printf("%s listening on %s", defaultServiceName, addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func startAuditDispatcher(db *sql.DB) {
	auditURL := os.Getenv("AUDIT_HISTORY_SERVICE_URL")
	secret := []byte(os.Getenv("SERVICE_TOKEN_SECRET"))
	if auditURL == "" || len(secret) == 0 {
		return
	}
	outbox := event.NewOutbox(db)
	config := event.DispatchConfig{
		ServiceID:              envOrDefault("SERVICE_ID", "commercial"),
		ServiceTokenSecret:     secret,
		AuditHistoryServiceURL: auditURL,
		ReportingServiceURL:    os.Getenv("REPORTING_SERVICE_URL"),
	}
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := outbox.DispatchOnce(ctx, config); err != nil {
				log.Printf("audit outbox dispatch: %v", err)
			}
			cancel()
			<-ticker.C
		}
	}()
}

func pingDatabase() error {
	db, err := openDatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return db.PingContext(ctx)
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
	return db, nil
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
