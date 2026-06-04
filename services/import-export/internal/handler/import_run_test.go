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
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCSVImportRunAcceptance(t *testing.T) {
	db := newImportExportTestDB(t)
	var imported []map[string]any
	var mu sync.Mutex
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/leads" {
			http.NotFound(w, r)
			return
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		mu.Lock()
		imported = append(imported, body)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "lead-" + fmt.Sprint(len(imported)), "companyName": body["companyName"], "source": body["source"]})
	}))
	t.Cleanup(target.Close)
	app := NewImportExportServer(db, Config{LeadServiceURL: target.URL, HTTPClient: target.Client()})

	t.Run("TEST-CSV-IMPORT-001/002 valid rows imported and invalid rows reported without corruption", func(t *testing.T) {
		rec := postImportJSON(app, "/imports", map[string]any{
			"objectType": "lead",
			"filename":   "leads.csv",
			"content":    "companyName,leadName,source,ownerId\nImport Good Co,Good lead,Website,sales-1\nImport Bad Co,Bad lead,,sales-1\n",
		}, actorHeaders("mgr-1", "Sales Manager", "single-team"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected import 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		if body["status"] != "CompletedWithErrors" || body["successCount"].(float64) != 1 || body["failureCount"].(float64) != 1 {
			t.Fatalf("expected partial success summary, got %#v", body)
		}
		if len(body["rowErrors"].([]any)) != 1 {
			t.Fatalf("expected one row error, got %s", rec.Body.String())
		}
		mu.Lock()
		defer mu.Unlock()
		if len(imported) != 1 || imported[0]["companyName"] != "Import Good Co" {
			t.Fatalf("expected only valid row sent to target command API, got %#v", imported)
		}
	})

	t.Run("TEST-CSV-IMPORT-003 unsupported format rejected before mutation", func(t *testing.T) {
		before := importCallCount(&mu, imported)
		rec := postImportJSON(app, "/imports", map[string]any{
			"objectType": "lead",
			"filename":   "leads.xlsx",
			"content":    "companyName,source\nUnsupported,Website\n",
		}, actorHeaders("mgr-1", "Sales Manager", "single-team"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected unsupported format 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "UNSUPPORTED_FORMAT")
		if importCallCount(&mu, imported) != before {
			t.Fatalf("unsupported format mutated target")
		}
	})

	t.Run("TEST-CSV-IMPORT-004 and TEST-ABUSE-IMPORTAUTHZ-001 Sales denied", func(t *testing.T) {
		rec := postImportJSON(app, "/imports", map[string]any{
			"objectType": "lead",
			"filename":   "leads.csv",
			"content":    "companyName,source\nDenied,Website\n",
		}, actorHeaders("sales-1", "Sales", "single-team"))
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected sales denied 403, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "PERMISSION_DENIED")
	})

	t.Run("TEST-ABUSE-CSVINJECT-001 dangerous cells fail row without mutation", func(t *testing.T) {
		before := importCallCount(&mu, imported)
		rec := postImportJSON(app, "/imports", map[string]any{
			"objectType": "lead",
			"filename":   "leads.csv",
			"content":    "companyName,source\n=cmd,Website\n",
		}, actorHeaders("mgr-1", "Sales Manager", "single-team"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected import run 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		if body["successCount"].(float64) != 0 || body["failureCount"].(float64) != 1 || importCallCount(&mu, imported) != before {
			t.Fatalf("expected dangerous row rejected without mutation, got %s", rec.Body.String())
		}
	})
}

func TestImportRunReadAuthorization(t *testing.T) {
	db := newImportExportTestDB(t)
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/leads" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "lead-authz-1"})
	}))
	t.Cleanup(target.Close)
	app := NewImportExportServer(db, Config{LeadServiceURL: target.URL, HTTPClient: target.Client()})

	create := postImportJSON(app, "/imports", map[string]any{
		"objectType": "lead",
		"filename":   "private-import-run.csv",
		"content":    "companyName,leadName,source,ownerId\nVisible Co,Visible Lead,Website,sales-1\nPrivate Bad Co,Bad lead,,sales-1\n",
	}, actorHeaders("mgr-owner", "Sales Manager", "team-a"))
	if create.Code != http.StatusCreated {
		t.Fatalf("expected import 201, got %d body=%s", create.Code, create.Body.String())
	}
	runID, ok := decodeJSON(t, create)["runId"].(string)
	if !ok || runID == "" {
		t.Fatalf("missing runId in response %s", create.Body.String())
	}

	for name, headers := range map[string]map[string]string{
		"owner manager":     actorHeaders("mgr-owner", "Sales Manager", "team-a"),
		"same-team manager": actorHeaders("mgr-peer", "Sales Manager", "team-a"),
		"administrator":     actorHeaders("admin-1", "Administrator", "other-team"),
	} {
		t.Run(name+" may read import run", func(t *testing.T) {
			rec := getImportJSON(app, "/imports/"+runID, headers)
			if rec.Code != http.StatusOK {
				t.Fatalf("expected authorized read 200, got %d body=%s", rec.Code, rec.Body.String())
			}
			body := decodeJSON(t, rec)
			if body["runId"] != runID || body["filename"] != "private-import-run.csv" {
				t.Fatalf("expected import run detail, got %#v", body)
			}
		})
	}

	for name, headers := range map[string]map[string]string{
		"cross-team manager": actorHeaders("mgr-other", "Sales Manager", "team-b"),
		"sales":              actorHeaders("sales-1", "Sales", "team-a"),
		"anonymous":          {},
	} {
		t.Run(name+" cannot infer import run existence", func(t *testing.T) {
			rec := getImportJSON(app, "/imports/"+runID, headers)
			if rec.Code != http.StatusNotFound {
				t.Fatalf("expected unauthorized read hidden as 404, got %d body=%s", rec.Code, rec.Body.String())
			}
			requireErrorCode(t, rec, "NOT_FOUND")
			body := rec.Body.String()
			for _, leaked := range []string{runID, "private-import-run.csv", "Private Bad Co"} {
				if strings.Contains(body, leaked) {
					t.Fatalf("unauthorized response leaked %q: %s", leaked, body)
				}
			}
		})
	}
}

func newImportExportTestDB(t *testing.T) *sql.DB {
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
	db := openImportExportDB(t, fmt.Sprintf("postgres://crm_admin:crm_admin_dev_password@%s:%s/crm_system?sslmode=disable", host, port.Port()))
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_import_export_runs.up.sql", "0003_export_runs.up.sql"} {
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

func openImportExportDB(t *testing.T, dsn string) *sql.DB {
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

func postImportJSON(handler http.Handler, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
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

func getImportJSON(handler http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
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

func importCallCount(mu *sync.Mutex, imported []map[string]any) int {
	mu.Lock()
	defer mu.Unlock()
	return len(imported)
}
