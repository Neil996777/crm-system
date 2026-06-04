package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"crm-system/services/work/internal/authz"
	"crm-system/services/work/internal/domain"
)

func (h *WorkHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	businessDate := parseBusinessDate(r.URL.Query().Get("businessDate"))
	activeOnly := r.URL.Query().Get("activeOnly") == "true"
	tasks, err := h.repo.ListTasks(r.Context(), actor.ID, actor.Role, r.URL.Query().Get("relatedType"), r.URL.Query().Get("relatedId"), activeOnly)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, taskDTO(task, businessDate))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *WorkHandler) listReminders(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	businessDateValue := normalizedBusinessDate(r.URL.Query().Get("businessDate"))
	businessDate := parseBusinessDate(businessDateValue)
	tasks, err := h.repo.ListTasks(r.Context(), actor.ID, actor.Role, "", "", true)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	rows := make([]domain.ReminderRow, 0, len(tasks))
	for _, task := range tasks {
		if row, ok := domain.ReminderFromTask(task, businessDate); ok {
			rows = append(rows, row)
		}
	}
	commercialRows, err := h.commercialReminderRows(r, actor, businessDateValue)
	if err != nil {
		writeError(w, http.StatusBadGateway, "DEPENDENCY_UNAVAILABLE", "dependency", "A required service is unavailable.")
		return
	}
	rows = append(rows, commercialRows...)
	writeJSON(w, http.StatusOK, map[string]any{
		"timezone":     domain.ReminderTimezone,
		"businessDate": businessDateValue,
		"rows":         rows,
	})
}

func (h *WorkHandler) listActivities(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	activities, err := h.repo.ListActivities(r.Context(), actor.ID, actor.Role, r.URL.Query().Get("relatedType"), r.URL.Query().Get("relatedId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(activities))
	for _, activity := range activities {
		items = append(items, activityDTO(activity))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *WorkHandler) listNotes(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	notes, err := h.repo.ListNotes(r.Context(), actor.ID, actor.Role, r.URL.Query().Get("relatedType"), r.URL.Query().Get("relatedId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(notes))
	for _, note := range notes {
		items = append(items, noteDTO(note))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func parseBusinessDate(value string) time.Time {
	if value == "" {
		return time.Now().UTC()
	}
	date, err := domain.ParseDate(value)
	if err != nil {
		return time.Now().UTC()
	}
	return date
}

func normalizedBusinessDate(value string) string {
	if strings.TrimSpace(value) == "" {
		return domain.FormatDate(time.Now().UTC())
	}
	return strings.TrimSpace(value)
}

func (h *WorkHandler) commercialReminderRows(r *http.Request, actor actorContext, businessDate string) ([]domain.ReminderRow, error) {
	if strings.TrimSpace(h.config.CommercialBaseURL) == "" {
		return nil, nil
	}
	token := authz.SignServiceToken(h.config.ServiceID, "commercial", "commercial.reminder_eligibility", h.config.ServiceTokenSecret)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, strings.TrimRight(h.config.CommercialBaseURL, "/")+"/internal/reminders/eligibility?businessDate="+businessDate, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", h.config.ServiceID)
	req.Header.Set("X-Intent", "commercial.reminder_eligibility")
	req.Header.Set("X-Actor-User-Id", actor.ID)
	req.Header.Set("X-Actor-Role", actor.Role)
	if correlationID := strings.TrimSpace(r.Header.Get("X-Correlation-Id")); correlationID != "" {
		req.Header.Set("X-Correlation-Id", correlationID)
	}
	client := h.config.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errDependencyStatus(resp.StatusCode)
	}
	var body struct {
		Rows []domain.ReminderRow `json:"rows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return body.Rows, nil
}

type errDependencyStatus int

func (e errDependencyStatus) Error() string {
	return "dependency status " + http.StatusText(int(e))
}
