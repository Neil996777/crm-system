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

func TestLeadOutboxDispatcherDeliversAuditHistoryAndReportingWithRetry(t *testing.T) {
	db := newLeadDispatcherTestDB(t)
	outbox := NewOutbox(db)
	ctx := context.Background()
	if err := outbox.Append(ctx, LeadQualified, "lead_dispatch_1", map[string]any{
		"actorId":       "sales-1",
		"actorRole":     "Sales",
		"actorDisplay":  "Sales One",
		"correlationId": "corr-lead-dispatch-1",
		"ownerId":       "sales-1",
		"status":        "Valid",
	}); err != nil {
		t.Fatalf("append outbox event: %v", err)
	}
	var eventUID string
	if err := db.QueryRow(`SELECT id FROM lead.outbox_events WHERE aggregate_id = $1`, "lead_dispatch_1").Scan(&eventUID); err != nil {
		t.Fatalf("read outbox id: %v", err)
	}

	auditFailures := 0
	auditFailServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auditFailures++
		http.Error(w, "audit unavailable", http.StatusServiceUnavailable)
	}))
	t.Cleanup(auditFailServer.Close)
	reportingCalls := 0
	reportingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reportingCalls++
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(reportingServer.Close)

	err := outbox.DispatchOnce(ctx, DispatchConfig{
		ServiceID:              "lead",
		ServiceTokenSecret:     []byte("dispatch-secret"),
		AuditHistoryServiceURL: auditFailServer.URL,
		ReportingServiceURL:    reportingServer.URL,
	})
	if err == nil {
		t.Fatal("TEST-HISTORY-DISPATCH-RETRY-001 expected failed audit append to keep event unpublished")
	}
	if auditFailures != 1 {
		t.Fatalf("expected one failed audit attempt, got %d", auditFailures)
	}
	if reportingCalls != 0 {
		t.Fatalf("reporting must not run after failed audit append, got %d calls", reportingCalls)
	}
	requireLeadPublishedState(t, db, eventUID, false)

	var receivedAudit map[string]any
	auditOKServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/events/append" {
			t.Fatalf("unexpected audit path %s", r.URL.Path)
		}
		if r.Header.Get("X-Service-Id") != "lead" || r.Header.Get("X-Intent") != "audit.append" {
			t.Fatalf("missing audit S2S headers: service=%q intent=%q", r.Header.Get("X-Service-Id"), r.Header.Get("X-Intent"))
		}
		if r.Header.Get("X-Correlation-Id") != "corr-lead-dispatch-1" {
			t.Fatalf("expected audit correlation id, got %q", r.Header.Get("X-Correlation-Id"))
		}
		if r.Header.Get("X-Actor-User-Id") != "sales-1" || r.Header.Get("X-Actor-Role") != "Sales" {
			t.Fatalf("missing actor headers: actor=%q role=%q", r.Header.Get("X-Actor-User-Id"), r.Header.Get("X-Actor-Role"))
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Fatalf("missing audit bearer service token")
		}
		if err := json.NewDecoder(r.Body).Decode(&receivedAudit); err != nil {
			t.Fatalf("decode audit append body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(auditOKServer.Close)
	reportingFailures := 0
	reportingFailServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reportingFailures++
		http.Error(w, "reporting unavailable", http.StatusServiceUnavailable)
	}))
	t.Cleanup(reportingFailServer.Close)

	err = outbox.DispatchOnce(ctx, DispatchConfig{
		ServiceID:              "lead",
		ServiceTokenSecret:     []byte("dispatch-secret"),
		AuditHistoryServiceURL: auditOKServer.URL,
		ReportingServiceURL:    reportingFailServer.URL,
	})
	if err == nil {
		t.Fatal("TEST-REPORTING-PROJECTION-INGEST-001 expected failed reporting projection to keep event unpublished")
	}
	if reportingFailures != 1 {
		t.Fatalf("expected one failed reporting attempt, got %d", reportingFailures)
	}
	requireLeadPublishedState(t, db, eventUID, false)

	var receivedProjection map[string]any
	reportingOKServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/projections" {
			t.Fatalf("unexpected reporting path %s", r.URL.Path)
		}
		if r.Header.Get("X-Service-Id") != "lead" || r.Header.Get("X-Intent") != "reporting.projection_ingest" {
			t.Fatalf("missing reporting S2S headers: service=%q intent=%q", r.Header.Get("X-Service-Id"), r.Header.Get("X-Intent"))
		}
		if r.Header.Get("X-Correlation-Id") != "corr-lead-dispatch-1" {
			t.Fatalf("expected reporting correlation id, got %q", r.Header.Get("X-Correlation-Id"))
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
	t.Cleanup(reportingOKServer.Close)

	if err := outbox.DispatchOnce(ctx, DispatchConfig{
		ServiceID:              "lead",
		ServiceTokenSecret:     []byte("dispatch-secret"),
		AuditHistoryServiceURL: auditOKServer.URL,
		ReportingServiceURL:    reportingOKServer.URL,
	}); err != nil {
		t.Fatalf("dispatch retry: %v", err)
	}
	requireLeadPublishedState(t, db, eventUID, true)
	if receivedAudit["eventUid"] != eventUID || receivedAudit["eventId"] != "EVT-LEAD-QUALIFIED" {
		t.Fatalf("expected qualified audit event uid/id, got %#v", receivedAudit)
	}
	if receivedProjection["sourceService"] != "lead" || receivedProjection["recordType"] != "lead" || receivedProjection["recordId"] != "lead_dispatch_1" {
		t.Fatalf("expected lead reporting projection, got %#v", receivedProjection)
	}
}

func TestLeadOutboxDispatcherMapsQualificationAuditEvents(t *testing.T) {
	qualified := auditAppendBody(outboxEvent{ID: "evt_q", EventType: LeadQualified, AggregateID: "lead_q", Payload: map[string]any{"status": "Valid"}})
	disqualified := auditAppendBody(outboxEvent{ID: "evt_d", EventType: LeadDisqualified, AggregateID: "lead_d", Payload: map[string]any{"status": "Invalid", "invalidReason": "No fit"}})

	if qualified["eventId"] != "EVT-LEAD-QUALIFIED" {
		t.Fatalf("expected qualified event id, got %#v", qualified["eventId"])
	}
	if disqualified["eventId"] != "EVT-LEAD-DISQUALIFIED" {
		t.Fatalf("expected disqualified event id, got %#v", disqualified["eventId"])
	}
	if qualified["eventId"] == disqualified["eventId"] {
		t.Fatalf("TEST-HISTORY-LEAD-EVENT-ID-001 expected distinct lead qualify/disqualify event ids")
	}
	if after, ok := disqualified["afterSummary"].(map[string]any); !ok || after["invalidReason"] != "No fit" {
		t.Fatalf("expected disqualified afterSummary to include invalid reason, got %#v", disqualified["afterSummary"])
	}
}

func requireLeadPublishedState(t *testing.T, db *sql.DB, eventUID string, published bool) {
	t.Helper()
	var isPublished bool
	if err := db.QueryRow(`SELECT published_at IS NOT NULL FROM lead.outbox_events WHERE id = $1`, eventUID).Scan(&isPublished); err != nil {
		t.Fatalf("read published state: %v", err)
	}
	if isPublished != published {
		t.Fatalf("expected published=%t for %s, got %t", published, eventUID, isPublished)
	}
}

func newLeadDispatcherTestDB(t *testing.T) *sql.DB {
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
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_leads.up.sql"} {
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
