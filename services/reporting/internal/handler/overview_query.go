package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"crm-system/services/reporting/internal/projection"
	"crm-system/services/reporting/internal/repo"
)

type Config struct{}

type ReportingHandler struct {
	repo     *repo.ProjectionRepo
	consumer *projection.Consumer
}

func NewReportingServer(db *sql.DB, _ Config) http.Handler {
	handler := &ReportingHandler{
		repo:     repo.NewProjectionRepo(db),
		consumer: projection.NewConsumer(db),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /reports/team-overview", handler.teamOverview)
	mux.HandleFunc("POST /internal/projections", handler.upsertProjection)
	return mux
}

func (h *ReportingHandler) teamOverview(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	if actor.Role == "Sales" || actor.Role == "" {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	scope := "team"
	teamID := actor.TeamID
	if teamID == "" {
		teamID = "single-team"
	}
	if actor.Role == "Administrator" {
		scope = "all"
	}
	data, err := h.repo.TeamOverview(r.Context(), teamID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"scope":      scope,
		"filters":    map[string]any{"teamId": teamID, "archived": "active_default"},
		"currency":   "CNY",
		"metrics":    data.Metrics.Map(),
		"pipeline":   data.Pipeline,
		"emptyState": data.Metrics.Empty(),
	})
}

func (h *ReportingHandler) upsertProjection(w http.ResponseWriter, r *http.Request) {
	var input projection.RecordProjection
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The projection input is invalid.")
		return
	}
	if err := h.consumer.Upsert(r.Context(), input); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The projection input is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "updated"})
}

type actorContext struct {
	ID     string
	Role   string
	TeamID string
}

func actorFromRequest(r *http.Request) actorContext {
	return actorContext{
		ID:     r.Header.Get("X-Actor-User-Id"),
		Role:   r.Header.Get("X-Actor-Role"),
		TeamID: r.Header.Get("X-Actor-Team-Id"),
	}
}

func writeError(w http.ResponseWriter, status int, code, category, safeMessage string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":        code,
			"category":    category,
			"safeMessage": safeMessage,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write response: %v", err)
	}
}
