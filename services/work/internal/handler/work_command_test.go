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

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWorkAcceptance(t *testing.T) {
	db := newWorkTestDB(t)
	app := NewWorkServer(db, Config{ServiceTokenSecret: []byte("work-test-secret")})

	t.Run("TEST-ACTIVITY-NOTE-001 and TEST-TASK-LIFECYCLE-001 persist related activity note and open task", func(t *testing.T) {
		activity := postWorkJSON(app, "/activities", map[string]any{
			"relatedType":  "Lead",
			"relatedId":    "lead_work_001",
			"activityType": "Call",
			"content":      "Discovery call completed",
			"ownerId":      "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if activity.Code != http.StatusCreated {
			t.Fatalf("expected activity create 201, got %d body=%s", activity.Code, activity.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, activity), "relatedId", "lead_work_001")

		note := postWorkJSON(app, "/notes", map[string]any{
			"relatedType": "Lead",
			"relatedId":   "lead_work_001",
			"content":     "Budget confirmed",
			"ownerId":     "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if note.Code != http.StatusCreated {
			t.Fatalf("expected note create 201, got %d body=%s", note.Code, note.Body.String())
		}

		task := postWorkJSON(app, "/tasks", map[string]any{
			"relatedType": "Lead",
			"relatedId":   "lead_work_001",
			"title":       "Send proposal",
			"dueDate":     "2027-03-01",
			"ownerId":     "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if task.Code != http.StatusCreated {
			t.Fatalf("expected task create 201, got %d body=%s", task.Code, task.Body.String())
		}
		taskBody := decodeJSON(t, task)
		requireJSONValue(t, taskBody, "status", "Open")
		requireEvent(t, db, "WorkItemCreated", taskBody["id"].(string))
	})

	t.Run("TEST-ACTIVITY-NOTE-002 rejects missing related record fields", func(t *testing.T) {
		rec := postWorkJSON(app, "/notes", map[string]any{
			"content": "Missing relation",
			"ownerId": "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected validation 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "VALIDATION_FAILED")
	})

	t.Run("TEST-TASK-LIFECYCLE-002 TEST-TASK-LIFECYCLE-003 and TEST-INV-TASKREMINDER-001 complete overdue task removes active reminder", func(t *testing.T) {
		task := postWorkJSON(app, "/tasks", map[string]any{
			"relatedType": "Opportunity",
			"relatedId":   "opp_work_overdue",
			"title":       "Follow up overdue deal",
			"dueDate":     "2026-01-01",
			"ownerId":     "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if task.Code != http.StatusCreated {
			t.Fatalf("expected task create 201, got %d body=%s", task.Code, task.Body.String())
		}
		taskID := decodeJSON(t, task)["id"].(string)

		overdue := getWork(app, "/tasks?businessDate=2026-06-02&activeOnly=true", actorHeaders("sales-1", "Sales"))
		if overdue.Code != http.StatusOK || !contains(overdue.Body.String(), taskID) || !contains(overdue.Body.String(), "Overdue") {
			t.Fatalf("expected overdue active task in list, got %d body=%s", overdue.Code, overdue.Body.String())
		}

		done := postWorkJSON(app, "/tasks/"+taskID+"/status", map[string]any{
			"toStatus": "Completed",
		}, actorHeaders("sales-1", "Sales"))
		if done.Code != http.StatusOK {
			t.Fatalf("expected complete 200, got %d body=%s", done.Code, done.Body.String())
		}

		active := getWork(app, "/reminders?businessDate=2026-06-02", actorHeaders("sales-1", "Sales"))
		if active.Code != http.StatusOK {
			t.Fatalf("expected reminders 200, got %d body=%s", active.Code, active.Body.String())
		}
		if contains(active.Body.String(), taskID) {
			t.Fatalf("completed task must not be active reminder: %s", active.Body.String())
		}
		requireEvent(t, db, "TaskStatusChanged", taskID)
	})

	t.Run("TEST-ACTIVITY-NOTE-003 and TEST-ABUSE-MUTATE-001 denies Sales creating work for another owner", func(t *testing.T) {
		rec := postWorkJSON(app, "/tasks", map[string]any{
			"relatedType": "Lead",
			"relatedId":   "lead_work_denied",
			"title":       "Unauthorized owner",
			"dueDate":     "2027-03-01",
			"ownerId":     "sales-2",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected permission 403, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "PERMISSION_DENIED")
	})
}

func TestOwnerTransferOpenWorkCascade(t *testing.T) {
	t.Log("TEST-OWNER-TRANSFER-004 open-work cascade transfers active tasks to the new owner")
	db := newWorkTestDB(t)
	app := NewWorkServer(db, Config{ServiceID: "work", ServiceTokenSecret: []byte("work-test-secret")})

	task := postWorkJSON(app, "/tasks", map[string]any{
		"relatedType": "Customer",
		"relatedId":   "acct_transfer_001",
		"title":       "Transfer with parent",
		"dueDate":     "2027-03-01",
		"ownerId":     "sales-1",
	}, actorHeaders("sales-1", "Sales"))
	if task.Code != http.StatusCreated {
		t.Fatalf("expected task create 201, got %d body=%s", task.Code, task.Body.String())
	}
	taskID := decodeJSON(t, task)["id"].(string)

	token := SignServiceToken("account", "work", "work.owner_transfer", []byte("work-test-secret"))
	transfer := postWorkJSON(app, "/internal/owner-transfer", map[string]any{
		"relatedType": "Customer",
		"relatedId":   "acct_transfer_001",
		"fromOwnerId": "sales-1",
		"toOwnerId":   "sales-2",
	}, map[string]string{
		"Authorization": "Bearer " + token,
		"X-Service-Id":  "account",
		"X-Intent":      "work.owner_transfer",
	})
	if transfer.Code != http.StatusOK {
		t.Fatalf("expected owner transfer 200, got %d body=%s", transfer.Code, transfer.Body.String())
	}

	list := getWork(app, "/tasks?businessDate=2026-06-02&activeOnly=true", actorHeaders("sales-2", "Sales"))
	if list.Code != http.StatusOK || !contains(list.Body.String(), taskID) {
		t.Fatalf("expected transferred task visible to new owner, got %d body=%s", list.Code, list.Body.String())
	}
}

func postWorkJSON(app http.Handler, path string, body map[string]any, headers map[string]string) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	app.ServeHTTP(rec, req)
	return rec
}

func getWork(app http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	app.ServeHTTP(rec, req)
	return rec
}

func actorHeaders(id, role string) map[string]string {
	return map[string]string{
		"X-Actor-User-Id": id,
		"X-Actor-Role":    role,
	}
}

func newWorkTestDB(t *testing.T) *sql.DB {
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
			WaitingFor: wait.ForListeningPort("5432/tcp"),
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
	adminDB, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatalf("open admin db: %v", err)
	}
	t.Cleanup(func() { _ = adminDB.Close() })
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_work.up.sql"} {
		sqlBytes, err := os.ReadFile(filepath.Join("..", "..", "migrations", migration))
		if err != nil {
			t.Fatalf("read migration %s: %v", migration, err)
		}
		if _, err := adminDB.Exec(string(sqlBytes)); err != nil {
			t.Fatalf("apply migration %s: %v", migration, err)
		}
	}
	serviceDSN := fmt.Sprintf("postgres://crm_work_user:crm_work_dev_password@%s:%s/crm_system?sslmode=disable&search_path=work", host, port.Port())
	db, err := sql.Open("pgx", serviceDSN)
	if err != nil {
		t.Fatalf("open service db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}
