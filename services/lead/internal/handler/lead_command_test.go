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

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestLeadCreateAssignAndConcurrencyAcceptance(t *testing.T) {
	db := newLeadTestDB(t)
	app := NewLeadServer(db, Config{})

	t.Run("TEST-LEAD-CREATE-001 creates owned lead with required fields persisted", func(t *testing.T) {
		rec := postLeadJSON(app, "/leads", map[string]any{
			"leadName":    "Inbound ERP evaluation",
			"companyName": "Acme Manufacturing",
			"source":      "Website",
			"ownerId":     "sales-1",
			"needSummary": "Needs CRM follow-up",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "status", "Pending Qualification")
		requireJSONValue(t, body, "ownerId", "sales-1")
		requireEvent(t, db, "LeadCreated", body["id"].(string))
	})

	t.Run("TEST-LEAD-CREATE-002 missing required fields blocked with safe validation", func(t *testing.T) {
		rec := postLeadJSON(app, "/leads", map[string]any{
			"leadName": "Missing source",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "VALIDATION_FAILED")
	})

	t.Run("TEST-LEAD-CREATE-003 creates unassigned lead", func(t *testing.T) {
		rec := postLeadJSON(app, "/leads", map[string]any{
			"companyName": "Unassigned Co",
			"source":      "Referral",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "status", "Unassigned")
		requireJSONValue(t, body, "ownerId", "")
	})

	t.Run("TEST-LEAD-ASSIGN-001 and TEST-OWNER-TRANSFER-001 manager assigns owner and records history event", func(t *testing.T) {
		create := postLeadJSON(app, "/leads", map[string]any{
			"companyName": "Assign Co",
			"source":      "Trade Show",
		}, actorHeaders("mgr-1", "Sales Manager"))
		leadID := decodeJSON(t, create)["id"].(string)

		assign := postLeadJSON(app, "/leads/"+leadID+"/owner-transfer", map[string]any{
			"expectedVersion": 1,
			"newOwnerId":      "sales-2",
			"reason":          "Territory assignment",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if assign.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", assign.Code, assign.Body.String())
		}
		body := decodeJSON(t, assign)
		requireJSONValue(t, body, "ownerId", "sales-2")
		requireJSONValue(t, body, "status", "Pending Qualification")
		requireEvent(t, db, "LeadOwnerChanged", leadID)
	})

	t.Run("TEST-LEAD-ASSIGN-002 and TEST-OWNER-TRANSFER-003 Sales cannot assign or transfer", func(t *testing.T) {
		create := postLeadJSON(app, "/leads", map[string]any{
			"leadName": "Sales owned",
			"source":   "Website",
			"ownerId":  "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		leadID := decodeJSON(t, create)["id"].(string)
		assign := postLeadJSON(app, "/leads/"+leadID+"/owner-transfer", map[string]any{
			"expectedVersion": 1,
			"newOwnerId":      "sales-2",
			"reason":          "Attempted reassignment",
		}, actorHeaders("sales-1", "Sales"))
		if assign.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d body=%s", assign.Code, assign.Body.String())
		}
		requireErrorCode(t, assign, "PERMISSION_DENIED")
	})

	t.Run("CONTRACT-020 expectedVersion prevents stale owner transfer", func(t *testing.T) {
		create := postLeadJSON(app, "/leads", map[string]any{
			"leadName": "Versioned lead",
			"source":   "Website",
			"ownerId":  "sales-1",
		}, actorHeaders("mgr-1", "Sales Manager"))
		leadID := decodeJSON(t, create)["id"].(string)
		conflict := postLeadJSON(app, "/leads/"+leadID+"/owner-transfer", map[string]any{
			"expectedVersion": 99,
			"newOwnerId":      "sales-2",
			"reason":          "Stale transfer",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if conflict.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", conflict.Code, conflict.Body.String())
		}
		requireErrorCode(t, conflict, "VERSION_CONFLICT")
	})
}

func newLeadTestDB(t *testing.T) *sql.DB {
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
	db := openLeadDB(t, adminDSN)
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_leads.up.sql", "0003_lead_qualification.up.sql", "0004_duplicate_warnings.up.sql", "0005_archive.up.sql"} {
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

func openLeadDB(t *testing.T, dsn string) *sql.DB {
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

func postLeadJSON(handler http.Handler, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer
	if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
		panic(err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, &requestBody)
	req.Header.Set("Content-Type", "application/json")
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	handler.ServeHTTP(rec, req)
	return rec
}

func actorHeaders(id, role string) map[string]string {
	return map[string]string{
		"X-Actor-User-Id": id,
		"X-Actor-Role":    role,
		"X-Actor-Display": id,
	}
}
