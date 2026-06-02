package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReminderQueryAcceptance(t *testing.T) {
	db := newWorkTestDB(t)
	commercial := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/reminders/eligibility" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("X-Service-Id") != "work" || r.Header.Get("X-Intent") != "commercial.reminder_eligibility" {
			http.Error(w, "missing service headers", http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"rows": []map[string]any{
			{
				"id":            "contract_reminder_001",
				"sourceService": "commercial-service",
				"type":          "contract_pending_signature",
				"relatedRecord": map[string]any{"type": "contract", "id": "contract_reminder_001", "display": "contract_reminder_001"},
				"ownerDisplay":  "sales-1",
				"dueDate":       "2026-01-01",
				"status":        "Overdue",
				"priority":      "P1",
				"version":       1,
			},
		}})
	}))
	defer commercial.Close()
	app := NewWorkServer(db, Config{ServiceID: "work", ServiceTokenSecret: []byte("work-test-secret"), CommercialBaseURL: commercial.URL})

	t.Run("TEST-REMINDER-001 TEST-REMINDER-002 and TEST-REMINDER-BOUNDARY-001 aggregate task and commercial reminders", func(t *testing.T) {
		task := postWorkJSON(app, "/tasks", map[string]any{
			"relatedType": "Opportunity",
			"relatedId":   "opp_reminder_task",
			"title":       "Due today task",
			"dueDate":     "2026-06-02",
			"ownerId":     "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if task.Code != http.StatusCreated {
			t.Fatalf("expected task create 201, got %d body=%s", task.Code, task.Body.String())
		}

		rec := getWork(app, "/reminders?businessDate=2026-06-02", actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected reminders 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		for _, expected := range []string{"Asia/Shanghai", "2026-06-02", "task_due", "DueToday", "contract_pending_signature", "contract_reminder_001"} {
			if !contains(body, expected) {
				t.Fatalf("expected %q in reminders body=%s", expected, body)
			}
		}
	})

	t.Run("TEST-REMINDER-004 and TEST-REMINDER-005 suppress inactive and unauthorized tasks", func(t *testing.T) {
		hidden := postWorkJSON(app, "/tasks", map[string]any{
			"relatedType": "Lead",
			"relatedId":   "lead_hidden_reminder",
			"title":       "Other owner task",
			"dueDate":     "2026-01-01",
			"ownerId":     "sales-2",
		}, actorHeaders("sales-2", "Sales"))
		if hidden.Code != http.StatusCreated {
			t.Fatalf("expected hidden task create 201, got %d body=%s", hidden.Code, hidden.Body.String())
		}
		done := postWorkJSON(app, "/tasks", map[string]any{
			"relatedType": "Lead",
			"relatedId":   "lead_done_reminder",
			"title":       "Completed reminder task",
			"dueDate":     "2026-01-01",
			"ownerId":     "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if done.Code != http.StatusCreated {
			t.Fatalf("expected done task create 201, got %d body=%s", done.Code, done.Body.String())
		}
		doneID := decodeJSON(t, done)["id"].(string)
		complete := postWorkJSON(app, "/tasks/"+doneID+"/status", map[string]any{"toStatus": "Completed"}, actorHeaders("sales-1", "Sales"))
		if complete.Code != http.StatusOK {
			t.Fatalf("expected complete 200, got %d body=%s", complete.Code, complete.Body.String())
		}

		rec := getWork(app, "/reminders?businessDate=2026-06-02", actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected reminders 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		for _, forbidden := range []string{"Other owner task", "Completed reminder task"} {
			if contains(body, forbidden) {
				t.Fatalf("forbidden reminder %q leaked: %s", forbidden, body)
			}
		}
	})
}
