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

	accountauthz "crm-system/services/account/internal/authz"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestAccountCRUDAcceptance(t *testing.T) {
	db := newAccountTestDB(t)
	app := NewAccountServer(db, Config{})

	t.Run("TEST-CUSTOMER-CRUD-001 creates account with required fields persisted", func(t *testing.T) {
		rec := postAccountJSON(app, "/accounts", map[string]any{
			"companyName":    "Acme Manufacturing",
			"customerStatus": "Prospect",
			"ownerId":        "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "companyName", "Acme Manufacturing")
		requireJSONValue(t, body, "customerStatus", "Prospect")
		requireJSONValue(t, body, "ownerId", "sales-1")
		requireEvent(t, db, "AccountCreated", body["id"].(string))

		fetch := getAccountJSON(app, "/accounts/"+body["id"].(string), actorHeaders("sales-1", "Sales"))
		if fetch.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", fetch.Code, fetch.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, fetch), "companyName", "Acme Manufacturing")
	})

	t.Run("TEST-CUSTOMER-CRUD-002 missing required fields blocked", func(t *testing.T) {
		rec := postAccountJSON(app, "/accounts", map[string]any{
			"companyName": "Missing status",
			"ownerId":     "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "VALIDATION_FAILED")
	})

	t.Run("TEST-CUSTOMER-CRUD-003 edits persist and expectedVersion is enforced", func(t *testing.T) {
		create := postAccountJSON(app, "/accounts", map[string]any{
			"companyName":    "Editable Co",
			"customerStatus": "Prospect",
			"ownerId":        "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		accountID := decodeJSON(t, create)["id"].(string)

		edit := patchAccountJSON(app, "/accounts/"+accountID, map[string]any{
			"expectedVersion": 1,
			"companyName":     "Editable Co Ltd",
			"customerStatus":  "Active",
			"ownerId":         "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if edit.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", edit.Code, edit.Body.String())
		}
		body := decodeJSON(t, edit)
		requireJSONValue(t, body, "companyName", "Editable Co Ltd")
		requireJSONValue(t, body, "customerStatus", "Active")
		if body["version"].(float64) != 2 {
			t.Fatalf("expected version 2, got %#v", body["version"])
		}
		requireEvent(t, db, "AccountUpdated", accountID)

		conflict := patchAccountJSON(app, "/accounts/"+accountID, map[string]any{
			"expectedVersion": 1,
			"companyName":     "Lost Update",
			"customerStatus":  "Active",
			"ownerId":         "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if conflict.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", conflict.Code, conflict.Body.String())
		}
		requireErrorCode(t, conflict, "VERSION_CONFLICT")
	})

	t.Run("TEST-CUSTOMER-CRUD-004 Sales unrelated denied without sensitive exposure", func(t *testing.T) {
		create := postAccountJSON(app, "/accounts", map[string]any{
			"companyName":    "Restricted Customer",
			"customerStatus": "Active",
			"ownerId":        "sales-2",
		}, actorHeaders("mgr-1", "Sales Manager"))
		accountID := decodeJSON(t, create)["id"].(string)

		fetch := getAccountJSON(app, "/accounts/"+accountID, actorHeaders("sales-1", "Sales"))
		if fetch.Code != http.StatusNotFound {
			t.Fatalf("expected safe 404, got %d body=%s", fetch.Code, fetch.Body.String())
		}
		requireErrorCode(t, fetch, "NOT_FOUND")
		if contains(fetch.Body.String(), "Restricted Customer") {
			t.Fatalf("unauthorized response leaked company name: %s", fetch.Body.String())
		}

		list := getAccountJSON(app, "/accounts?search=Restricted", actorHeaders("sales-1", "Sales"))
		if list.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", list.Code, list.Body.String())
		}
		items := decodeJSON(t, list)["items"].([]any)
		if len(items) != 0 {
			t.Fatalf("expected no unrelated records, got %#v", items)
		}
	})

	t.Run("TEST-INV-NODELETE-001 hard delete route unavailable and record persists", func(t *testing.T) {
		create := postAccountJSON(app, "/accounts", map[string]any{
			"companyName":    "No Delete Co",
			"customerStatus": "Active",
			"ownerId":        "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		accountID := decodeJSON(t, create)["id"].(string)
		del := deleteAccountJSON(app, "/accounts/"+accountID, actorHeaders("sales-1", "Sales"))
		if del.Code != http.StatusMethodNotAllowed && del.Code != http.StatusNotFound {
			t.Fatalf("expected unavailable delete route, got %d body=%s", del.Code, del.Body.String())
		}
		fetch := getAccountJSON(app, "/accounts/"+accountID, actorHeaders("sales-1", "Sales"))
		if fetch.Code != http.StatusOK {
			t.Fatalf("expected record to persist after delete attempt, got %d body=%s", fetch.Code, fetch.Body.String())
		}
	})
}

func TestAccountLeadConversionCreateIdempotency(t *testing.T) {
	db := newAccountTestDB(t)
	app := NewAccountServer(db, Config{ServiceID: "account", ServiceTokenSecret: []byte("account-test-secret")})
	headers := leadConversionHeaders(t, "account-test-secret")
	body := map[string]any{
		"idempotencyKey": "lead-convert-account-key",
		"companyName":    "Lead Converted Account",
		"customerStatus": "Prospect",
		"ownerId":        "sales-1",
	}
	first := postAccountJSON(app, "/internal/accounts", body, headers)
	if first.Code != http.StatusCreated {
		t.Fatalf("expected first internal create 201, got %d body=%s", first.Code, first.Body.String())
	}
	firstID := decodeJSON(t, first)["id"].(string)
	second := postAccountJSON(app, "/internal/accounts", body, headers)
	if second.Code != http.StatusOK {
		t.Fatalf("expected idempotent retry 200, got %d body=%s", second.Code, second.Body.String())
	}
	if decodeJSON(t, second)["id"] != firstID {
		t.Fatalf("expected retry to return original account id %s, got %s", firstID, second.Body.String())
	}
	var count int
	if err := db.QueryRow(`SELECT count(*) FROM account.accounts WHERE company_name = $1`, "Lead Converted Account").Scan(&count); err != nil {
		t.Fatalf("count accounts: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one account row for idempotent conversion create, got %d", count)
	}

	t.Run("TEST-LEAD-CONVERSION-IDEMPOTENCY-004 concurrent duplicate key returns existing account", func(t *testing.T) {
		key := "lead-convert-account-race-key"
		lockTx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			t.Fatalf("begin lock tx: %v", err)
		}
		defer lockTx.Rollback()
		if _, err := lockTx.Exec(`LOCK TABLE account.accounts IN SHARE ROW EXCLUSIVE MODE`); err != nil {
			t.Fatalf("lock account table: %v", err)
		}

		done := make(chan *httptest.ResponseRecorder, 1)
		go func() {
			done <- postAccountJSON(app, "/internal/accounts", map[string]any{
				"idempotencyKey": key,
				"companyName":    "Concurrent Lead Converted Account",
				"customerStatus": "Prospect",
				"ownerId":        "sales-1",
			}, headers)
		}()
		time.Sleep(250 * time.Millisecond)
		if _, err := lockTx.Exec(`
			INSERT INTO account.accounts (id, company_name, customer_status, owner_id, version, lead_conversion_idempotency_key)
			VALUES ('acct_existing_race', 'Concurrent Existing Account', 'Prospect', 'sales-1', 1, $1)
		`, key); err != nil {
			t.Fatalf("insert competing account: %v", err)
		}
		if err := lockTx.Commit(); err != nil {
			t.Fatalf("commit competing account: %v", err)
		}
		rec := <-done
		if rec.Code != http.StatusOK {
			t.Fatalf("expected duplicate-key conversion create to return existing 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		if decodeJSON(t, rec)["id"] != "acct_existing_race" {
			t.Fatalf("expected existing account from duplicate key, got %s", rec.Body.String())
		}
	})
}

func newAccountTestDB(t *testing.T) *sql.DB {
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
	db := openAccountDB(t, adminDSN)
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_accounts.up.sql", "0003_contacts.up.sql", "0004_duplicate_warnings.up.sql", "0005_archive.up.sql", "0006_lead_conversion_idempotency.up.sql"} {
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

func leadConversionHeaders(t *testing.T, secret string) map[string]string {
	t.Helper()
	token, err := accountauthz.SignServiceToken(accountauthz.ServiceTokenClaims{
		Issuer:   "lead",
		Audience: "account",
		Intent:   "account.create_for_lead_conversion",
		Expires:  time.Now().UTC().Add(2 * time.Minute),
	}, []byte(secret))
	if err != nil {
		t.Fatalf("sign service token: %v", err)
	}
	headers := actorHeaders("sales-1", "Sales")
	headers["Authorization"] = "Bearer " + token
	headers["X-Service-Id"] = "lead"
	headers["X-Intent"] = "account.create_for_lead_conversion"
	return headers
}

func openAccountDB(t *testing.T, dsn string) *sql.DB {
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

func postAccountJSON(handler http.Handler, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	return requestAccountJSON(handler, http.MethodPost, path, body, headers)
}

func patchAccountJSON(handler http.Handler, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	return requestAccountJSON(handler, http.MethodPatch, path, body, headers)
}

func getAccountJSON(handler http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	return requestAccountJSON(handler, http.MethodGet, path, nil, headers)
}

func deleteAccountJSON(handler http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	return requestAccountJSON(handler, http.MethodDelete, path, nil, headers)
}

func requestAccountJSON(handler http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
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
