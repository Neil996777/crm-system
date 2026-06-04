package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"crm-system/services/reporting/internal/authz"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestReportingAccessDeniedDispatchesAuditCatalogEventAndRetries(t *testing.T) {
	db := newEventTestDB(t)
	outbox := NewOutbox(db)
	if err := outbox.AppendReportAccessDenied(context.Background(), ReportAccessDeniedInput{
		ActorID:       "sales-denied",
		ActorRole:     "Sales",
		ActorDisplay:  "sales-denied",
		ReportType:    "team-overview",
		CorrelationID: "corr-report-denied",
	}); err != nil {
		t.Fatalf("append access denied: %v", err)
	}

	var attempts int
	var delivered map[string]any
	audit := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			http.Error(w, "temporary outage", http.StatusServiceUnavailable)
			return
		}
		if r.URL.Path != "/internal/events/append" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("X-Service-Id") != "reporting" || r.Header.Get("X-Intent") != "audit.append" {
			t.Fatalf("missing S2S headers: %#v", r.Header)
		}
		if r.Header.Get("X-Correlation-Id") != "corr-report-denied" {
			t.Fatalf("missing correlation id: %#v", r.Header)
		}
		if _, err := authz.VerifyServiceToken(
			r.Header.Get("Authorization")[len("Bearer "):],
			authz.VerifyOptions{Secret: []byte("secret"), Audience: "audit-history", Intent: "audit.append", Now: time.Now().UTC()},
		); err != nil {
			t.Fatalf("verify token: %v", err)
		}
		if err := json.NewDecoder(r.Body).Decode(&delivered); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(audit.Close)

	err := outbox.DispatchOnce(context.Background(), DispatchConfig{
		ServiceID:              "reporting",
		ServiceTokenSecret:     []byte("secret"),
		AuditHistoryServiceURL: audit.URL,
		BatchSize:              10,
	})
	if err == nil {
		t.Fatal("expected first dispatch failure")
	}
	requireUnpublishedCount(t, db, 1)

	if err := outbox.DispatchOnce(context.Background(), DispatchConfig{
		ServiceID:              "reporting",
		ServiceTokenSecret:     []byte("secret"),
		AuditHistoryServiceURL: audit.URL,
		BatchSize:              10,
	}); err != nil {
		t.Fatalf("dispatch retry: %v", err)
	}
	requireUnpublishedCount(t, db, 0)
	if delivered["eventId"] != "EVT-REPORT-ACCESS-DENIED" || delivered["resourceType"] != "Report" || delivered["resourceId"] != "team-overview" {
		t.Fatalf("unexpected delivered body: %#v", delivered)
	}
}

func newEventTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "crm_system",
				"POSTGRES_USER":     "crm_admin",
				"POSTGRES_PASSWORD": "crm_admin_dev_password",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("start postgres testcontainer: %v", err)
	}
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Fatalf("terminate postgres testcontainer: %v", err)
		}
	})
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("container port: %v", err)
	}
	db, err := sql.Open("pgx", fmt.Sprintf("postgres://crm_admin:crm_admin_dev_password@%s:%s/crm_system?sslmode=disable", host, port.Port()))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_reporting_projections.up.sql", "0003_reporting_outbox.up.sql"} {
		sqlBytes, err := os.ReadFile(filepath.Join("..", "..", "migrations", migration))
		if err != nil {
			t.Fatalf("read migration %s: %v", migration, err)
		}
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			t.Fatalf("apply migration %s: %v", migration, err)
		}
	}
	return db
}

func requireUnpublishedCount(t *testing.T, db *sql.DB, expected int) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT count(*) FROM reporting.outbox_events WHERE published_at IS NULL`).Scan(&count); err != nil {
		t.Fatalf("count unpublished: %v", err)
	}
	if count != expected {
		t.Fatalf("expected %d unpublished events, got %d", expected, count)
	}
}
