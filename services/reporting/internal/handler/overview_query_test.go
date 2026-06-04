package handler

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"crm-system/services/reporting/internal/event"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestManagerTeamOverviewAcceptance(t *testing.T) {
	db := newReportingTestDB(t)
	app := NewReportingServer(db, Config{})

	insertProjection(t, db, map[string]any{"recordType": "lead", "recordId": "lead_team", "ownerId": "sales-1", "teamId": "single-team", "status": "Pending Qualification"})
	insertProjection(t, db, map[string]any{"recordType": "opportunity", "recordId": "opp_team", "ownerId": "sales-1", "teamId": "single-team", "stage": "Quote", "amount": "10000.00"})
	insertProjection(t, db, map[string]any{"recordType": "quote", "recordId": "quote_team", "ownerId": "sales-1", "teamId": "single-team", "status": "Accepted", "amount": "10000.00"})
	insertProjection(t, db, map[string]any{"recordType": "contract", "recordId": "contract_team", "ownerId": "sales-1", "teamId": "single-team", "status": "Signed", "amount": "10000.00"})
	insertProjection(t, db, map[string]any{"recordType": "payment", "recordId": "payment_team", "ownerId": "sales-1", "teamId": "single-team", "status": "PartiallyPaid", "amount": "3000.00"})
	insertProjection(t, db, map[string]any{"recordType": "task", "recordId": "task_team", "ownerId": "sales-1", "teamId": "single-team", "status": "Open"})
	insertProjection(t, db, map[string]any{"recordType": "lead", "recordId": "lead_other", "ownerId": "sales-9", "teamId": "other-team", "status": "Pending Qualification"})

	t.Run("TEST-TEAM-OVERVIEW-001/004 manager sees only team records and pipeline status", func(t *testing.T) {
		rec := getReportingJSON(app, "/reports/team-overview", actorHeaders("mgr-1", "Sales Manager", "single-team"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected manager overview 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		metrics := body["metrics"].(map[string]any)
		if metrics["leadCount"].(float64) != 1 || metrics["opportunityCount"].(float64) != 1 || metrics["taskCount"].(float64) != 1 {
			t.Fatalf("expected single-team counts only, got %#v", metrics)
		}
		if metrics["quoteAmount"] != "10000.00" || metrics["contractAmount"] != "10000.00" || metrics["paidAmount"] != "3000.00" {
			t.Fatalf("expected team amounts only, got %#v", metrics)
		}
		if body["scope"] != "team" || body["emptyState"] != false {
			t.Fatalf("expected team non-empty overview, got %s", rec.Body.String())
		}
		if !contains(rec.Body.String(), "Quote") || contains(rec.Body.String(), "lead_other") {
			t.Fatalf("expected pipeline status without non-team leakage, got %s", rec.Body.String())
		}
	})

	t.Run("TEST-TEAM-OVERVIEW-002 Sales denied", func(t *testing.T) {
		rec := getReportingJSON(app, "/reports/team-overview", actorHeaders("sales-1", "Sales", "single-team"))
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected sales denied 403, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "PERMISSION_DENIED")
	})

	t.Run("TEST-TEAM-OVERVIEW-003 empty authorized team returns empty state", func(t *testing.T) {
		rec := getReportingJSON(app, "/reports/team-overview", actorHeaders("mgr-empty", "Sales Manager", "empty-team"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected empty overview 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		metrics := body["metrics"].(map[string]any)
		if body["emptyState"] != true || metrics["leadCount"].(float64) != 0 {
			t.Fatalf("expected empty-state metrics, got %s", rec.Body.String())
		}
	})
}

func TestProjectionIngestRequiresS2SToken(t *testing.T) {
	db := newReportingTestDB(t)
	app := NewReportingServer(db, Config{ServiceID: "reporting", ServiceTokenSecret: []byte("reporting-secret")})
	body := map[string]any{
		"sourceService": "opportunity",
		"recordType":    "opportunity",
		"recordId":      "opp_s2s",
		"ownerId":       "sales-1",
		"teamId":        "single-team",
		"stage":         "Quote",
		"amount":        "1000.00",
	}

	for name, token := range map[string]string{
		"TEST-REPORTING-S2S-001 missing token":   "",
		"TEST-REPORTING-S2S-002 expired token":   signReportingTestToken(t, "opportunity", "reporting", "reporting.projection_ingest", []byte("reporting-secret"), time.Now().Add(-time.Minute)),
		"TEST-REPORTING-S2S-003 wrong audience":  signReportingTestToken(t, "opportunity", "account", "reporting.projection_ingest", []byte("reporting-secret"), time.Now().Add(2*time.Minute)),
		"TEST-REPORTING-S2S-004 wrong intent":    signReportingTestToken(t, "opportunity", "reporting", "audit.append", []byte("reporting-secret"), time.Now().Add(2*time.Minute)),
		"TEST-REPORTING-S2S-005 wrong signature": signReportingTestToken(t, "opportunity", "reporting", "reporting.projection_ingest", []byte("wrong-secret"), time.Now().Add(2*time.Minute)),
	} {
		t.Run(name, func(t *testing.T) {
			rec := postProjectionJSON(app, body, token, "opportunity", "reporting.projection_ingest")
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
			}
			requireErrorCode(t, rec, "SERVICE_AUTH_FAILED")
			requireProjectionCount(t, db, "opp_s2s", 0)
		})
	}

	valid := signReportingTestToken(t, "opportunity", "reporting", "reporting.projection_ingest", []byte("reporting-secret"), time.Now().Add(2*time.Minute))
	rec := postProjectionJSON(app, body, valid, "opportunity", "reporting.projection_ingest")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected valid S2S projection ingest 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	requireProjectionCount(t, db, "opp_s2s", 1)
	overview := getReportingJSON(app, "/reports/team-overview", actorHeaders("mgr-1", "Sales Manager", "single-team"))
	if overview.Code != http.StatusOK {
		t.Fatalf("TEST-REPORTING-PROJECTION-INGEST-006 expected manager query over ingested projection 200, got %d body=%s", overview.Code, overview.Body.String())
	}
	metrics := decodeJSON(t, overview)["metrics"].(map[string]any)
	pipeline := decodeJSON(t, overview)["pipeline"].([]any)
	if metrics["opportunityCount"].(float64) != 1 || len(pipeline) != 1 || pipeline[0].(map[string]any)["amount"] != "1000.00" {
		t.Fatalf("expected ingested projection in manager aggregate, got %s", overview.Body.String())
	}
	denied := getReportingJSON(app, "/reports/team-overview", actorHeaders("sales-1", "Sales", "single-team"))
	if denied.Code != http.StatusForbidden {
		t.Fatalf("expected Sales denied with no aggregate leakage, got %d body=%s", denied.Code, denied.Body.String())
	}
	requireErrorCode(t, denied, "PERMISSION_DENIED")
}

func TestReportAccessDeniedCreatesCatalogEvent(t *testing.T) {
	db := newReportingTestDB(t)
	app := NewReportingServer(db, Config{ServiceID: "reporting", ServiceTokenSecret: []byte("reporting-secret")})
	rec := getReportingJSON(app, "/reports/team-overview", actorHeaders("sales-denied", "Sales", "single-team"))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected sales denied 403, got %d body=%s", rec.Code, rec.Body.String())
	}
	var eventID, payload string
	err := db.QueryRow(`
		SELECT event_type, payload::text
		FROM reporting.outbox_events
		WHERE aggregate_id = 'team-overview'
		ORDER BY occurred_at DESC
		LIMIT 1
	`).Scan(&eventID, &payload)
	if err != nil {
		t.Fatalf("TEST-EVT-CATALOG-REPORTING-001 expected reporting access denied outbox event: %v", err)
	}
	if eventID != event.ReportAccessDenied || !contains(payload, "sales-denied") {
		t.Fatalf("expected report access denied event with actor, got event=%s payload=%s", eventID, payload)
	}
}

func newReportingTestDB(t *testing.T) *sql.DB {
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
	db := openReportingDB(t, fmt.Sprintf("postgres://crm_admin:crm_admin_dev_password@%s:%s/crm_system?sslmode=disable", host, port.Port()))
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_reporting_projections.up.sql", "0003_reporting_outbox.up.sql"} {
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

func openReportingDB(t *testing.T, dsn string) *sql.DB {
	t.Helper()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}
	return db
}

func insertProjection(t *testing.T, db *sql.DB, values map[string]any) {
	t.Helper()
	_, err := db.Exec(`
		INSERT INTO reporting.record_projections
			(source_service, record_type, record_id, owner_id, team_id, status, stage, amount, updated_at)
		VALUES ('test-source', $1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), $7::numeric, now())
	`, values["recordType"], values["recordId"], values["ownerId"], values["teamId"], stringValue(values["status"]), stringValue(values["stage"]), amountValue(values["amount"]))
	if err != nil {
		t.Fatalf("insert projection: %v", err)
	}
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	return value.(string)
}

func amountValue(value any) string {
	if value == nil {
		return "0.00"
	}
	return value.(string)
}

func getReportingJSON(handler http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	return requestReportingJSON(handler, http.MethodGet, path, nil, headers)
}

func postProjectionJSON(handler http.Handler, body any, token, serviceID, intent string) *httptest.ResponseRecorder {
	headers := map[string]string{"X-Service-Id": serviceID, "X-Intent": intent}
	if token != "" {
		headers["Authorization"] = "Bearer " + token
	}
	return requestReportingJSON(handler, http.MethodPost, "/internal/projections", body, headers)
}

func requestReportingJSON(handler http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
			panic(err)
		}
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, &requestBody)
	req.Header.Set("Content-Type", "application/json")
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	handler.ServeHTTP(rec, req)
	return rec
}

func actorHeaders(userID, role, teamID string) map[string]string {
	return map[string]string{
		"X-Actor-User-Id": userID,
		"X-Actor-Role":    role,
		"X-Actor-Team-Id": teamID,
	}
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode json body %q: %v", rec.Body.String(), err)
	}
	return body
}

func requireErrorCode(t *testing.T, rec *httptest.ResponseRecorder, expected string) {
	t.Helper()
	body := decodeJSON(t, rec)
	errBody, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("missing error body: %#v", body)
	}
	if errBody["code"] != expected {
		t.Fatalf("expected error code %q, got %#v body=%s", expected, errBody["code"], rec.Body.String())
	}
}

func contains(haystack, needle string) bool {
	return bytes.Contains([]byte(haystack), []byte(needle))
}

func requireProjectionCount(t *testing.T, db *sql.DB, recordID string, expected int) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT count(*) FROM reporting.record_projections WHERE record_id = $1`, recordID).Scan(&count); err != nil {
		t.Fatalf("count projection: %v", err)
	}
	if count != expected {
		t.Fatalf("expected projection count %d for %s, got %d", expected, recordID, count)
	}
}

func signReportingTestToken(t *testing.T, issuer, audience, intent string, secret []byte, expires time.Time) string {
	t.Helper()
	payload, err := json.Marshal(map[string]any{
		"iss":    issuer,
		"aud":    audience,
		"intent": intent,
		"exp":    expires.UTC(),
	})
	if err != nil {
		t.Fatalf("marshal token: %v", err)
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(encodedPayload))
	return encodedPayload + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
