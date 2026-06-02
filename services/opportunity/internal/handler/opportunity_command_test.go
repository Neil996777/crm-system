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

func TestOpportunityCreateAcceptance(t *testing.T) {
	db := newOpportunityTestDB(t)
	app := NewOpportunityServer(db, Config{})

	t.Run("TEST-OPP-CREATE-001 creates opportunity with required fields and stage persisted", func(t *testing.T) {
		rec := postOpportunityJSON(app, "/opportunities", map[string]any{
			"customerId":        "acct_001",
			"ownerId":           "sales-1",
			"stage":             "New Opportunity",
			"expectedAmount":    "120000.50",
			"expectedCloseDate": "2026-07-15",
			"title":             "ERP expansion",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "customerId", "acct_001")
		requireJSONValue(t, body, "ownerId", "sales-1")
		requireJSONValue(t, body, "stage", "New Opportunity")
		requireJSONValue(t, body, "expectedAmount", "120000.50")
		requireEvent(t, db, "OpportunityCreated", body["id"].(string))
	})

	t.Run("TEST-OPP-CREATE-002 missing required fields blocked", func(t *testing.T) {
		rec := postOpportunityJSON(app, "/opportunities", map[string]any{
			"customerId":        "acct_missing_amount",
			"ownerId":           "sales-1",
			"stage":             "New Opportunity",
			"expectedCloseDate": "2026-07-15",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "VALIDATION_FAILED")
	})

	t.Run("TEST-OPP-CREATE-003 persists plain opportunity and exposes no Status field", func(t *testing.T) {
		rec := postOpportunityJSON(app, "/opportunities", map[string]any{
			"customerId":        "acct_no_status",
			"ownerId":           "sales-1",
			"stage":             "New Opportunity",
			"expectedAmount":    "9000.00",
			"expectedCloseDate": "2026-08-20",
			"title":             "No status deal",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		if _, ok := body["status"]; ok {
			t.Fatalf("status field must not be exposed: %#v", body)
		}
		accountID := body["id"].(string)
		fetch := getOpportunityJSON(app, "/opportunities/"+accountID, actorHeaders("sales-1", "Sales"))
		if fetch.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", fetch.Code, fetch.Body.String())
		}
		if _, ok := decodeJSON(t, fetch)["status"]; ok {
			t.Fatalf("detail response exposed retired status field: %s", fetch.Body.String())
		}
		requireNoStatusColumn(t, db)
	})

	t.Run("TEST-OPP-CREATE-004 non-owned edit denied and hard delete unavailable", func(t *testing.T) {
		create := postOpportunityJSON(app, "/opportunities", map[string]any{
			"customerId":        "acct_restricted",
			"ownerId":           "sales-2",
			"stage":             "New Opportunity",
			"expectedAmount":    "30000.00",
			"expectedCloseDate": "2026-09-01",
			"title":             "Restricted deal",
		}, actorHeaders("mgr-1", "Sales Manager"))
		opportunityID := decodeJSON(t, create)["id"].(string)

		edit := patchOpportunityJSON(app, "/opportunities/"+opportunityID, map[string]any{
			"expectedVersion":   1,
			"customerId":        "acct_restricted",
			"ownerId":           "sales-2",
			"stage":             "New Opportunity",
			"expectedAmount":    "31000.00",
			"expectedCloseDate": "2026-09-01",
			"title":             "Restricted deal updated",
		}, actorHeaders("sales-1", "Sales"))
		if edit.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d body=%s", edit.Code, edit.Body.String())
		}
		requireErrorCode(t, edit, "PERMISSION_DENIED")
		if contains(edit.Body.String(), "Restricted deal") {
			t.Fatalf("unauthorized response leaked opportunity data: %s", edit.Body.String())
		}

		del := deleteOpportunityJSON(app, "/opportunities/"+opportunityID, actorHeaders("mgr-1", "Sales Manager"))
		if del.Code != http.StatusMethodNotAllowed && del.Code != http.StatusNotFound {
			t.Fatalf("expected unavailable delete route, got %d body=%s", del.Code, del.Body.String())
		}
		fetch := getOpportunityJSON(app, "/opportunities/"+opportunityID, actorHeaders("mgr-1", "Sales Manager"))
		if fetch.Code != http.StatusOK {
			t.Fatalf("expected opportunity to persist after delete attempt, got %d body=%s", fetch.Code, fetch.Body.String())
		}
	})
}

func newOpportunityTestDB(t *testing.T) *sql.DB {
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
	db := openOpportunityDB(t, adminDSN)
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_opportunities.up.sql", "0003_archive.up.sql"} {
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

func openOpportunityDB(t *testing.T, dsn string) *sql.DB {
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

func postOpportunityJSON(handler http.Handler, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	return requestOpportunityJSON(handler, http.MethodPost, path, body, headers)
}

func patchOpportunityJSON(handler http.Handler, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	return requestOpportunityJSON(handler, http.MethodPatch, path, body, headers)
}

func getOpportunityJSON(handler http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	return requestOpportunityJSON(handler, http.MethodGet, path, nil, headers)
}

func deleteOpportunityJSON(handler http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	return requestOpportunityJSON(handler, http.MethodDelete, path, nil, headers)
}

func requestOpportunityJSON(handler http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
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

func actorHeaders(id, role string) map[string]string {
	return map[string]string{
		"X-Actor-User-Id": id,
		"X-Actor-Role":    role,
	}
}
