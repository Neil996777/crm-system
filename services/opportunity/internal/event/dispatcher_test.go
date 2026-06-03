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
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestOutboxDispatcherDeliversAndRetriesAuditHistoryEvents(t *testing.T) {
	db := newDispatcherTestDB(t)
	outbox := NewOutbox(db)
	ctx := context.Background()
	if err := outbox.Append(ctx, OpportunityStageChanged, "opp_dispatch_1", map[string]any{
		"actorId":       "sales-1",
		"actorRole":     "Sales",
		"actorDisplay":  "Sales One",
		"correlationId": "corr-dispatch-1",
		"fromStage":     "New Opportunity",
		"toStage":       "Needs Confirmed",
	}); err != nil {
		t.Fatalf("append outbox event: %v", err)
	}

	var eventUID string
	if err := db.QueryRow(`SELECT id FROM opportunity.outbox_events WHERE aggregate_id = $1`, "opp_dispatch_1").Scan(&eventUID); err != nil {
		t.Fatalf("read outbox id: %v", err)
	}

	failures := 0
	failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failures++
		http.Error(w, "audit unavailable", http.StatusServiceUnavailable)
	}))
	defer failServer.Close()

	err := outbox.DispatchOnce(ctx, DispatchConfig{
		ServiceID:              "opportunity",
		ServiceTokenSecret:     []byte("dispatch-secret"),
		AuditHistoryServiceURL: failServer.URL,
	})
	if err == nil {
		t.Fatal("TEST-HISTORY-DISPATCH-RETRY-001 expected failed audit append to keep event unpublished")
	}
	if failures != 1 {
		t.Fatalf("expected one failed delivery attempt, got %d", failures)
	}
	requirePublishedState(t, db, eventUID, false)

	var received map[string]any
	var receivedProjection map[string]any
	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Service-Id") != "opportunity" || r.Header.Get("X-Intent") != "audit.append" {
			t.Fatalf("missing S2S identity headers: service=%q intent=%q", r.Header.Get("X-Service-Id"), r.Header.Get("X-Intent"))
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Fatalf("missing bearer service token")
		}
		if r.Header.Get("X-Actor-User-Id") != "sales-1" || r.Header.Get("X-Actor-Role") != "Sales" {
			t.Fatalf("missing actor headers: actor=%q role=%q", r.Header.Get("X-Actor-User-Id"), r.Header.Get("X-Actor-Role"))
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode audit append body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer okServer.Close()
	reportingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/projections" {
			t.Fatalf("unexpected reporting path %s", r.URL.Path)
		}
		if r.Header.Get("X-Service-Id") != "opportunity" || r.Header.Get("X-Intent") != "reporting.projection_ingest" {
			t.Fatalf("missing reporting S2S headers: service=%q intent=%q", r.Header.Get("X-Service-Id"), r.Header.Get("X-Intent"))
		}
		if r.Header.Get("X-Correlation-Id") != "corr-dispatch-1" {
			t.Fatalf("TEST-REPORTING-CORRELATION-001 expected reporting correlation id, got %q", r.Header.Get("X-Correlation-Id"))
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Fatalf("missing reporting bearer service token")
		}
		if err := json.NewDecoder(r.Body).Decode(&receivedProjection); err != nil {
			t.Fatalf("decode reporting projection body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"updated"}`))
	}))
	defer reportingServer.Close()

	if err := outbox.DispatchOnce(ctx, DispatchConfig{
		ServiceID:              "opportunity",
		ServiceTokenSecret:     []byte("dispatch-secret"),
		AuditHistoryServiceURL: okServer.URL,
		ReportingServiceURL:    reportingServer.URL,
	}); err != nil {
		t.Fatalf("dispatch retry: %v", err)
	}
	requirePublishedState(t, db, eventUID, true)
	if received["eventUid"] != eventUID {
		t.Fatalf("expected audit eventUid to use outbox id %q, got %#v", eventUID, received["eventUid"])
	}
	if received["eventId"] != "EVT-STAGE-CHANGED" {
		t.Fatalf("expected EVT-STAGE-CHANGED, got %#v", received["eventId"])
	}
	if receivedProjection["sourceService"] != "opportunity" || receivedProjection["recordType"] != "opportunity" || receivedProjection["recordId"] != "opp_dispatch_1" {
		t.Fatalf("TEST-REPORTING-PROJECTION-INGEST-001 expected opportunity projection from dispatcher, got %#v", receivedProjection)
	}
}

func requirePublishedState(t *testing.T, db *sql.DB, eventUID string, published bool) {
	t.Helper()
	var isPublished bool
	if err := db.QueryRow(`SELECT published_at IS NOT NULL FROM opportunity.outbox_events WHERE id = $1`, eventUID).Scan(&isPublished); err != nil {
		t.Fatalf("read published state: %v", err)
	}
	if isPublished != published {
		t.Fatalf("expected published=%t for %s, got %t", published, eventUID, isPublished)
	}
}

func newDispatcherTestDB(t *testing.T) *sql.DB {
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
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_opportunities.up.sql"} {
		sqlBytes, err := os.ReadFile(filepath.Join("..", "..", "migrations", migration))
		if err != nil {
			t.Fatalf("read migration %s: %v", migration, err)
		}
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			t.Fatalf("apply migration %s: %v", migration, err)
		}
	}
	t.Cleanup(func() { db.Close() })
	return db
}
