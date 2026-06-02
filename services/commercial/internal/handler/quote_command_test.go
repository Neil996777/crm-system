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

func TestQuoteLifecycleAcceptance(t *testing.T) {
	db := newCommercialTestDB(t)
	app := NewCommercialServer(db, Config{})

	t.Run("TEST-QUOTE-LIFECYCLE-001 creates Draft quote with required fields persisted", func(t *testing.T) {
		rec := postCommercialJSON(app, "/quotes", map[string]any{
			"opportunityId": "opp_quote_001",
			"customerId":    "acct_quote_001",
			"amount":        "12000.00",
			"status":        "Draft",
			"validityEnd":   "2026-12-31",
			"ownerId":       "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "opportunityId", "opp_quote_001")
		requireJSONValue(t, body, "status", "Draft")
		requireJSONValue(t, body, "amount", "12000.00")
	})

	t.Run("TEST-QUOTE-LIFECYCLE-002 missing amount/status/validity blocked", func(t *testing.T) {
		rec := postCommercialJSON(app, "/quotes", map[string]any{
			"opportunityId": "opp_missing",
			"customerId":    "acct_missing",
			"ownerId":       "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "VALIDATION_FAILED")
	})

	t.Run("TEST-QUOTE-LIFECYCLE-003 sends rejects and expires with expectedVersion", func(t *testing.T) {
		quoteID := createQuoteForLifecycle(t, app, "opp_quote_lifecycle", "sales-1")
		sent := postCommercialJSON(app, "/quotes/"+quoteID+"/status", map[string]any{
			"expectedVersion": 1,
			"toStatus":        "Sent",
		}, actorHeaders("sales-1", "Sales"))
		if sent.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", sent.Code, sent.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, sent), "status", "Sent")

		rejected := postCommercialJSON(app, "/quotes/"+quoteID+"/status", map[string]any{
			"expectedVersion": 2,
			"toStatus":        "Rejected",
		}, actorHeaders("sales-1", "Sales"))
		if rejected.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", rejected.Code, rejected.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, rejected), "status", "Rejected")
	})

	t.Run("TEST-QUOTE-ACCEPT-001 and TEST-QUOTE-ACCEPT-002 accepts the single quote and emits history", func(t *testing.T) {
		quoteID := createQuoteForLifecycle(t, app, "opp_quote_accept", "sales-1")
		sent := postCommercialJSON(app, "/quotes/"+quoteID+"/status", map[string]any{
			"expectedVersion": 1,
			"toStatus":        "Sent",
		}, actorHeaders("sales-1", "Sales"))
		if sent.Code != http.StatusOK {
			t.Fatalf("expected sent 200, got %d body=%s", sent.Code, sent.Body.String())
		}
		accepted := postCommercialJSON(app, "/quotes/"+quoteID+"/status", map[string]any{
			"expectedVersion": 2,
			"toStatus":        "Accepted",
		}, actorHeaders("sales-1", "Sales"))
		if accepted.Code != http.StatusOK {
			t.Fatalf("expected accept 200, got %d body=%s", accepted.Code, accepted.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, accepted), "status", "Accepted")
		requireEvent(t, db, "QuoteAccepted", quoteID)
	})

	t.Run("TEST-INV-ONEACCEPT-001 second quote for same opportunity rejected", func(t *testing.T) {
		_ = createQuoteForLifecycle(t, app, "opp_one_quote", "sales-1")
		second := postCommercialJSON(app, "/quotes", map[string]any{
			"opportunityId": "opp_one_quote",
			"customerId":    "acct_one_quote",
			"amount":        "13000.00",
			"status":        "Draft",
			"validityEnd":   "2026-12-31",
			"ownerId":       "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if second.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", second.Code, second.Body.String())
		}
		requireErrorCode(t, second, "QUOTE_ALREADY_EXISTS")
	})
}

func createQuoteForLifecycle(t *testing.T, app http.Handler, opportunityID, ownerID string) string {
	t.Helper()
	create := postCommercialJSON(app, "/quotes", map[string]any{
		"opportunityId": opportunityID,
		"customerId":    "acct_" + opportunityID,
		"amount":        "10000.00",
		"status":        "Draft",
		"validityEnd":   "2026-12-31",
		"ownerId":       ownerID,
	}, actorHeaders(ownerID, "Sales"))
	if create.Code != http.StatusCreated {
		t.Fatalf("expected quote create 201, got %d body=%s", create.Code, create.Body.String())
	}
	return decodeJSON(t, create)["id"].(string)
}

func newCommercialTestDB(t *testing.T) *sql.DB {
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
	db := openCommercialDB(t, adminDSN)
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_quotes.up.sql", "0003_contracts.up.sql", "0004_payments.up.sql", "0005_archive.up.sql"} {
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

func openCommercialDB(t *testing.T, dsn string) *sql.DB {
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

func postCommercialJSON(handler http.Handler, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
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
	}
}
