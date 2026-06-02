package handler

import (
	"bytes"
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

	"crm-system/services/audit-history/internal/authz"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestAppendHistoryOperationLogAndHashChain(t *testing.T) {
	adminDB, serviceDB := newAuditTestDB(t)
	app := NewAuditServer(serviceDB, Config{ServiceID: "audit-history", ServiceTokenSecret: []byte("audit-secret")})
	token := signTestServiceToken(t, []byte("audit-secret"), "lead", "audit-history", "audit.append")

	first := appendEvent(t, app, token, map[string]any{
		"eventId":            "EVT-OWNER-CHANGED",
		"eventVersion":       1,
		"surfaces":           []string{"record_history", "operation_log"},
		"action":             "Owner changed",
		"resourceType":       "Lead",
		"resourceId":         "lead-1",
		"result":             "success",
		"safeSummary":        "Owner changed",
		"beforeSummary":      map[string]any{"ownerId": "sales-old"},
		"afterSummary":       map[string]any{"ownerId": "sales-1"},
		"acceptanceIds":      []string{"ACC-014", "ACC-022"},
		"actorUserId":        "payload-claimed-actor",
		"correlationId":      "corr-1",
		"causationId":        "cmd-1",
		"diffClassification": "Confidential",
	}, actorHeaders("sales-1", "Sales", "Sales One"))
	if first.Code != http.StatusCreated {
		t.Fatalf("expected append 201, got %d body=%s", first.Code, first.Body.String())
	}
	firstBody := decodeJSON(t, first)
	firstHash := firstBody["eventHash"].(string)
	if firstBody["prevHash"] != "" {
		t.Fatalf("first event prevHash should be empty, got %#v", firstBody["prevHash"])
	}

	second := appendEvent(t, app, token, map[string]any{
		"eventId":       "EVT-STAGE-CHANGED",
		"eventVersion":  1,
		"surfaces":      []string{"record_history"},
		"action":        "Stage changed",
		"resourceType":  "Opportunity",
		"resourceId":    "opp-1",
		"result":        "success",
		"safeSummary":   "Stage changed",
		"acceptanceIds": []string{"ACC-014"},
	}, actorHeaders("sales-1", "Sales", "Sales One"))
	if second.Code != http.StatusCreated {
		t.Fatalf("expected append 201, got %d body=%s", second.Code, second.Body.String())
	}
	requireJSONValue(t, decodeJSON(t, second), "prevHash", firstHash)

	var storedActor string
	if err := adminDB.QueryRow(`SELECT actor_user_id FROM audit_history.events WHERE event_uid = $1`, firstBody["eventUid"]).Scan(&storedActor); err != nil {
		t.Fatalf("read stored actor: %v", err)
	}
	if storedActor != "sales-1" {
		t.Fatalf("TEST-ABUSE-ACTAS-001 expected authenticated actor sales-1, got %q", storedActor)
	}
}

func TestAppendRejectsMissingS2SAndEventsAreAppendOnly(t *testing.T) {
	adminDB, serviceDB := newAuditTestDB(t)
	app := NewAuditServer(serviceDB, Config{ServiceID: "audit-history", ServiceTokenSecret: []byte("audit-secret")})

	rec := appendEvent(t, app, "invalid-token", map[string]any{"eventId": "EVT-OWNER-CHANGED"}, actorHeaders("sales-1", "Sales", "Sales One"))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected invalid S2S 401, got %d body=%s", rec.Code, rec.Body.String())
	}
	requireErrorCode(t, rec, "SERVICE_AUTH_FAILED")

	var updateAllowed bool
	err := adminDB.QueryRow(`
		SELECT has_table_privilege('crm_audit_history_user', 'audit_history.events', 'UPDATE')
	`).Scan(&updateAllowed)
	if err != nil {
		t.Fatalf("check update privilege: %v", err)
	}
	if updateAllowed {
		t.Fatal("TEST-HISTORY-004 expected crm_audit_history_user to lack UPDATE on audit_history.events")
	}
}

func newAuditTestDB(t *testing.T) (*sql.DB, *sql.DB) {
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
	adminDSN := fmt.Sprintf("postgres://crm_admin:crm_admin_dev_password@%s:%s/crm_system?sslmode=disable", host, port.Port())
	adminDB := openAuditDB(t, adminDSN)
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_history_oplog.up.sql"} {
		sqlBytes, err := os.ReadFile(filepath.Join("..", "..", "migrations", migration))
		if err != nil {
			t.Fatalf("read migration %s: %v", migration, err)
		}
		if _, err := adminDB.Exec(string(sqlBytes)); err != nil {
			t.Fatalf("apply migration %s: %v", migration, err)
		}
	}
	serviceDSN := fmt.Sprintf("postgres://crm_audit_history_user:crm_audit_history_dev_password@%s:%s/crm_system?sslmode=disable&search_path=audit_history", host, port.Port())
	serviceDB := openAuditDB(t, serviceDSN)
	t.Cleanup(func() {
		adminDB.Close()
		serviceDB.Close()
	})
	return adminDB, serviceDB
}

func openAuditDB(t *testing.T, dsn string) *sql.DB {
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

func appendEvent(t *testing.T, handler http.Handler, token string, body any, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var requestBody bytes.Buffer
	if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
		t.Fatalf("encode body: %v", err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/internal/events/append", &requestBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", "lead")
	req.Header.Set("X-Intent", "audit.append")
	req.Header.Set("X-Correlation-Id", "corr-1")
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	handler.ServeHTTP(rec, req)
	return rec
}

func actorHeaders(id, role, display string) map[string]string {
	return map[string]string{
		"X-Actor-User-Id": id,
		"X-Actor-Role":    role,
		"X-Actor-Display": display,
	}
}

func signTestServiceToken(t *testing.T, secret []byte, issuer, audience, intent string) string {
	t.Helper()
	token, err := authz.SignServiceToken(authz.ServiceTokenClaims{
		Issuer:   issuer,
		Audience: audience,
		Intent:   intent,
		Expires:  time.Now().Add(5 * time.Minute),
	}, secret)
	if err != nil {
		t.Fatalf("sign service token: %v", err)
	}
	return token
}

func getHistory(handler http.Handler, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	handler.ServeHTTP(rec, req)
	return rec
}

func httptestDelete(handler http.Handler, path string, _ any) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	handler.ServeHTTP(rec, req)
	return rec
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode json body %q: %v", rec.Body.String(), err)
	}
	return body
}

func requireJSONValue(t *testing.T, body map[string]any, path, expected string) {
	t.Helper()
	value, ok := body[path]
	if !ok {
		t.Fatalf("missing json path %s in %#v", path, body)
	}
	if value != expected {
		t.Fatalf("expected %s=%q, got %#v in %#v", path, expected, value, body)
	}
}

func requireErrorCode(t *testing.T, rec *httptest.ResponseRecorder, expected string) {
	t.Helper()
	body := decodeJSON(t, rec)
	errBody, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("missing error body: %#v", body)
	}
	code, ok := errBody["code"].(string)
	if !ok {
		t.Fatalf("missing error code: %#v", errBody)
	}
	if code != expected {
		t.Fatalf("expected error code %q, got %q body=%s", expected, code, rec.Body.String())
	}
}
