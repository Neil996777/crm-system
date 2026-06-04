package event

import (
	"context"
	"database/sql"
	"encoding/base64"
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

func TestOutboxDispatcherDeliversRetriesAndDedupesAuditHistoryEvents(t *testing.T) {
	db := newDispatcherTestDB(t)
	outbox := NewOutbox(db)
	ctx := context.Background()
	if err := outbox.Append(ctx, OwnerChanged, "acct_dispatch_1", map[string]any{
		"actorId":       "mgr-1",
		"actorRole":     "Sales Manager",
		"actorDisplay":  "Manager One",
		"correlationId": "corr-account-dispatch",
		"ownerId":       "sales-2",
	}); err != nil {
		t.Fatalf("append outbox event: %v", err)
	}
	eventUID := readOutboxID(t, db, "acct_dispatch_1")

	failures := 0
	failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failures++
		http.Error(w, "audit unavailable", http.StatusServiceUnavailable)
	}))
	defer failServer.Close()
	if err := outbox.DispatchOnce(ctx, DispatchConfig{
		ServiceID:              "account",
		ServiceTokenSecret:     []byte("dispatch-secret"),
		AuditHistoryServiceURL: failServer.URL,
	}); err == nil {
		t.Fatal("TEST-HISTORY-DISPATCH-RETRY-001 expected failed audit append to keep account event unpublished")
	}
	if failures != 1 {
		t.Fatalf("expected one failed delivery attempt, got %d", failures)
	}
	requirePublishedState(t, db, eventUID, false)

	seen := map[string]int{}
	var received map[string]any
	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requireAuditRequest(t, r, "account", "mgr-1", "Sales Manager", "corr-account-dispatch")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode audit append body: %v", err)
		}
		uid, _ := received["eventUid"].(string)
		seen[uid]++
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer okServer.Close()

	if err := outbox.DispatchOnce(ctx, DispatchConfig{
		ServiceID:              "account",
		ServiceTokenSecret:     []byte("dispatch-secret"),
		AuditHistoryServiceURL: okServer.URL,
	}); err != nil {
		t.Fatalf("dispatch retry: %v", err)
	}
	requirePublishedState(t, db, eventUID, true)
	if received["eventUid"] != eventUID || received["eventId"] != "EVT-OWNER-CHANGED" {
		t.Fatalf("TEST-HISTORY-DISPATCH-RETRY-001 expected account owner-change event uid/id, got %#v", received)
	}

	if _, err := db.Exec(`UPDATE account.outbox_events SET published_at = NULL WHERE id = $1`, eventUID); err != nil {
		t.Fatalf("reset published state for duplicate retry: %v", err)
	}
	if err := outbox.DispatchOnce(ctx, DispatchConfig{
		ServiceID:              "account",
		ServiceTokenSecret:     []byte("dispatch-secret"),
		AuditHistoryServiceURL: okServer.URL,
	}); err != nil {
		t.Fatalf("dispatch duplicate retry: %v", err)
	}
	if len(seen) != 1 || seen[eventUID] != 2 {
		t.Fatalf("TEST-HISTORY-IDEMPOTENT-001 expected consumer idempotency by one eventUid with duplicate attempts, got %#v", seen)
	}
}

func TestOutboxDispatcherDeliversReportingProjectionAndRetries(t *testing.T) {
	db := newDispatcherTestDB(t)
	outbox := NewOutbox(db)
	ctx := context.Background()
	if err := outbox.Append(ctx, AccountCreated, "acct_reporting_1", map[string]any{
		"actorId":        "mgr-1",
		"actorRole":      "Sales Manager",
		"actorDisplay":   "Manager One",
		"correlationId":  "corr-account-reporting",
		"ownerId":        "sales-2",
		"teamId":         "team-a",
		"customerStatus": "Active",
	}); err != nil {
		t.Fatalf("append outbox event: %v", err)
	}
	eventUID := readOutboxID(t, db, "acct_reporting_1")

	auditServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requireAuditRequest(t, r, "account", "mgr-1", "Sales Manager", "corr-account-reporting")
		w.WriteHeader(http.StatusCreated)
	}))
	defer auditServer.Close()

	var reportingAttempts int
	reportingFail := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reportingAttempts++
		http.Error(w, "reporting unavailable", http.StatusServiceUnavailable)
	}))
	defer reportingFail.Close()
	if err := outbox.DispatchOnce(ctx, DispatchConfig{
		ServiceID:              "account",
		ServiceTokenSecret:     []byte("dispatch-secret"),
		AuditHistoryServiceURL: auditServer.URL,
		ReportingServiceURL:    reportingFail.URL,
	}); err == nil {
		t.Fatal("TEST-REPORTING-DISPATCH-ACCOUNT-001 expected reporting failure to keep account event unpublished")
	}
	if reportingAttempts != 1 {
		t.Fatalf("expected one reporting failure attempt, got %d", reportingAttempts)
	}
	requirePublishedState(t, db, eventUID, false)

	var projection map[string]any
	reportingOK := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requireReportingRequest(t, r, "account", "corr-account-reporting")
		if err := json.NewDecoder(r.Body).Decode(&projection); err != nil {
			t.Fatalf("decode reporting projection: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer reportingOK.Close()
	if err := outbox.DispatchOnce(ctx, DispatchConfig{
		ServiceID:              "account",
		ServiceTokenSecret:     []byte("dispatch-secret"),
		AuditHistoryServiceURL: auditServer.URL,
		ReportingServiceURL:    reportingOK.URL,
	}); err != nil {
		t.Fatalf("dispatch reporting retry: %v", err)
	}
	requirePublishedState(t, db, eventUID, true)
	if projection["sourceService"] != "account" || projection["recordType"] != "account" || projection["recordId"] != "acct_reporting_1" || projection["ownerId"] != "sales-2" || projection["teamId"] != "team-a" || projection["status"] != "Active" {
		t.Fatalf("TEST-REPORTING-DISPATCH-ACCOUNT-002 unexpected account projection: %#v", projection)
	}
}

func requireAuditRequest(t *testing.T, r *http.Request, serviceID, actorID, actorRole, correlationID string) {
	t.Helper()
	if r.URL.Path != "/internal/events/append" {
		t.Fatalf("unexpected audit path %s", r.URL.Path)
	}
	if r.Header.Get("X-Service-Id") != serviceID || r.Header.Get("X-Intent") != "audit.append" {
		t.Fatalf("missing S2S identity headers: service=%q intent=%q", r.Header.Get("X-Service-Id"), r.Header.Get("X-Intent"))
	}
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == "" || token == r.Header.Get("Authorization") {
		t.Fatalf("missing bearer service token")
	}
	payload := strings.Split(token, ".")[0]
	decoded, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		t.Fatalf("decode token payload: %v", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(decoded, &claims); err != nil {
		t.Fatalf("decode token claims: %v", err)
	}
	if claims["iss"] != serviceID || claims["aud"] != "audit-history" || claims["intent"] != "audit.append" {
		t.Fatalf("unexpected service token claims: %#v", claims)
	}
	if r.Header.Get("X-Actor-User-Id") != actorID || r.Header.Get("X-Actor-Role") != actorRole {
		t.Fatalf("missing actor headers: actor=%q role=%q", r.Header.Get("X-Actor-User-Id"), r.Header.Get("X-Actor-Role"))
	}
	if r.Header.Get("X-Correlation-Id") != correlationID {
		t.Fatalf("missing correlation id: %q", r.Header.Get("X-Correlation-Id"))
	}
}

func requireReportingRequest(t *testing.T, r *http.Request, serviceID, correlationID string) {
	t.Helper()
	if r.URL.Path != "/internal/projections" {
		t.Fatalf("unexpected reporting path %s", r.URL.Path)
	}
	if r.Header.Get("X-Service-Id") != serviceID || r.Header.Get("X-Intent") != "reporting.projection_ingest" {
		t.Fatalf("missing reporting S2S headers: service=%q intent=%q", r.Header.Get("X-Service-Id"), r.Header.Get("X-Intent"))
	}
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == "" || token == r.Header.Get("Authorization") {
		t.Fatalf("missing bearer service token")
	}
	payload := strings.Split(token, ".")[0]
	decoded, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		t.Fatalf("decode token payload: %v", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(decoded, &claims); err != nil {
		t.Fatalf("decode token claims: %v", err)
	}
	if claims["iss"] != serviceID || claims["aud"] != "reporting" || claims["intent"] != "reporting.projection_ingest" {
		t.Fatalf("unexpected reporting token claims: %#v", claims)
	}
	if r.Header.Get("X-Correlation-Id") != correlationID {
		t.Fatalf("missing reporting correlation id: %q", r.Header.Get("X-Correlation-Id"))
	}
}

func readOutboxID(t *testing.T, db *sql.DB, aggregateID string) string {
	t.Helper()
	var eventUID string
	if err := db.QueryRow(`SELECT id FROM account.outbox_events WHERE aggregate_id = $1`, aggregateID).Scan(&eventUID); err != nil {
		t.Fatalf("read outbox id: %v", err)
	}
	return eventUID
}

func requirePublishedState(t *testing.T, db *sql.DB, eventUID string, published bool) {
	t.Helper()
	var isPublished bool
	if err := db.QueryRow(`SELECT published_at IS NOT NULL FROM account.outbox_events WHERE id = $1`, eventUID).Scan(&isPublished); err != nil {
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
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_accounts.up.sql", "0003_contacts.up.sql", "0004_duplicate_warnings.up.sql", "0005_archive.up.sql"} {
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
