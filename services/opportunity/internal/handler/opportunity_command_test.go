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

	t.Run("TEST-OPP-CREATE-004 TEST-AUTHZ-SCOPE-005 and TEST-ABUSE-MUTATE-001 non-owned edit denied with no mutation and hard delete unavailable", func(t *testing.T) {
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
		detail := getOpportunityJSON(app, "/opportunities/"+opportunityID, actorHeaders("sales-1", "Sales"))
		if detail.Code != http.StatusNotFound {
			t.Fatalf("expected safe 404 for non-owned detail, got %d body=%s", detail.Code, detail.Body.String())
		}
		requireErrorCode(t, detail, "NOT_FOUND")
		if contains(detail.Body.String(), "Restricted deal") || contains(detail.Body.String(), opportunityID) {
			t.Fatalf("detail denial leaked opportunity data: %s", detail.Body.String())
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

	t.Run("TEST-HISTORY-TX-001 outbox enqueue failure rolls back opportunity create", func(t *testing.T) {
		if _, err := db.Exec(`
			CREATE OR REPLACE FUNCTION opportunity.fail_outbox_insert() RETURNS trigger AS $$
			BEGIN
				RAISE EXCEPTION 'forced outbox failure';
			END;
			$$ LANGUAGE plpgsql;
			CREATE TRIGGER fail_outbox_insert
			BEFORE INSERT ON opportunity.outbox_events
			FOR EACH ROW EXECUTE FUNCTION opportunity.fail_outbox_insert();
		`); err != nil {
			t.Fatalf("install failing outbox trigger: %v", err)
		}
		t.Cleanup(func() {
			_, _ = db.Exec(`DROP TRIGGER IF EXISTS fail_outbox_insert ON opportunity.outbox_events; DROP FUNCTION IF EXISTS opportunity.fail_outbox_insert();`)
		})
		rec := postOpportunityJSON(app, "/opportunities", map[string]any{
			"customerId":        "acct_tx_rollback",
			"ownerId":           "sales-1",
			"stage":             "New Opportunity",
			"expectedAmount":    "12000.00",
			"expectedCloseDate": "2026-10-01",
			"title":             "Must rollback",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code < 400 {
			t.Fatalf("expected outbox failure to fail request, got %d body=%s", rec.Code, rec.Body.String())
		}
		var count int
		if err := db.QueryRow(`SELECT count(*) FROM opportunity.opportunities WHERE customer_id = $1`, "acct_tx_rollback").Scan(&count); err != nil {
			t.Fatalf("count rolled back opportunity: %v", err)
		}
		if count != 0 {
			t.Fatalf("expected rollback to leave no opportunity row, got %d", count)
		}
	})
}

func TestOpportunityLeadConversionCreateIdempotency(t *testing.T) {
	db := newOpportunityTestDB(t)
	app := NewOpportunityServer(db, Config{ServiceID: "opportunity", ServiceTokenSecret: []byte("opportunity-test-secret")})
	headers := leadConversionHeaders(t, "opportunity-test-secret")
	body := map[string]any{
		"idempotencyKey":    "lead-convert-opportunity-key",
		"customerId":        "acct_from_conversion",
		"ownerId":           "sales-1",
		"stage":             "New Opportunity",
		"expectedAmount":    "50000.00",
		"expectedCloseDate": "2026-10-01",
		"title":             "Lead converted opportunity",
	}
	first := postOpportunityJSON(app, "/internal/opportunities", body, headers)
	if first.Code != http.StatusCreated {
		t.Fatalf("expected first internal create 201, got %d body=%s", first.Code, first.Body.String())
	}
	firstID := decodeJSON(t, first)["id"].(string)
	second := postOpportunityJSON(app, "/internal/opportunities", body, headers)
	if second.Code != http.StatusOK {
		t.Fatalf("expected idempotent retry 200, got %d body=%s", second.Code, second.Body.String())
	}
	if decodeJSON(t, second)["id"] != firstID {
		t.Fatalf("expected retry to return original opportunity id %s, got %s", firstID, second.Body.String())
	}
	var count int
	if err := db.QueryRow(`SELECT count(*) FROM opportunity.opportunities WHERE title = $1`, "Lead converted opportunity").Scan(&count); err != nil {
		t.Fatalf("count opportunities: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one opportunity row for idempotent conversion create, got %d", count)
	}

	t.Run("TEST-LEAD-CONVERSION-IDEMPOTENCY-005 concurrent duplicate key returns existing opportunity", func(t *testing.T) {
		key := "lead-convert-opportunity-race-key"
		lockTx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			t.Fatalf("begin lock tx: %v", err)
		}
		defer lockTx.Rollback()
		if _, err := lockTx.Exec(`LOCK TABLE opportunity.opportunities IN SHARE ROW EXCLUSIVE MODE`); err != nil {
			t.Fatalf("lock opportunity table: %v", err)
		}

		done := make(chan *httptest.ResponseRecorder, 1)
		go func() {
			done <- postOpportunityJSON(app, "/internal/opportunities", map[string]any{
				"idempotencyKey":    key,
				"customerId":        "acct_from_conversion_race",
				"ownerId":           "sales-1",
				"stage":             "New Opportunity",
				"expectedAmount":    "75000.00",
				"expectedCloseDate": "2026-10-15",
				"title":             "Concurrent lead converted opportunity",
			}, headers)
		}()
		time.Sleep(250 * time.Millisecond)
		if _, err := lockTx.Exec(`
			INSERT INTO opportunity.opportunities
				(id, customer_id, owner_id, stage, expected_amount, expected_close_date, title, version, lead_conversion_idempotency_key)
			VALUES ('opp_existing_race', 'acct_from_conversion_race', 'sales-1', 'New Opportunity', 75000.00, '2026-10-15', 'Concurrent existing opportunity', 1, $1)
		`, key); err != nil {
			t.Fatalf("insert competing opportunity: %v", err)
		}
		if err := lockTx.Commit(); err != nil {
			t.Fatalf("commit competing opportunity: %v", err)
		}
		rec := <-done
		if rec.Code != http.StatusOK {
			t.Fatalf("expected duplicate-key conversion create to return existing 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		if decodeJSON(t, rec)["id"] != "opp_existing_race" {
			t.Fatalf("expected existing opportunity from duplicate key, got %s", rec.Body.String())
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
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_opportunities.up.sql", "0003_archive.up.sql", "0004_lead_conversion_idempotency.up.sql"} {
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
	payload, err := json.Marshal(map[string]any{
		"iss":    "lead",
		"aud":    "opportunity",
		"intent": "opportunity.create_for_lead_conversion",
		"exp":    time.Now().UTC().Add(2 * time.Minute),
	})
	if err != nil {
		t.Fatalf("marshal service token: %v", err)
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(encodedPayload))
	token := encodedPayload + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	headers := actorHeaders("sales-1", "Sales")
	headers["Authorization"] = "Bearer " + token
	headers["X-Service-Id"] = "lead"
	headers["X-Intent"] = "opportunity.create_for_lead_conversion"
	return headers
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
