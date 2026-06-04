package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"crm-system/services/reporting/internal/authz"
	"crm-system/services/reporting/internal/event"
	"crm-system/services/reporting/internal/projection"
	"crm-system/services/reporting/internal/repo"
)

const intentProjectionIngest = "reporting.projection_ingest"

type Config struct {
	ServiceID          string
	ServiceTokenSecret []byte
}

type ReportingHandler struct {
	repo     *repo.ProjectionRepo
	consumer *projection.Consumer
	outbox   *event.Outbox
	config   Config
}

func NewReportingServer(db *sql.DB, config Config) http.Handler {
	if config.ServiceID == "" {
		config.ServiceID = "reporting"
	}
	handler := &ReportingHandler{
		repo:     repo.NewProjectionRepo(db),
		consumer: projection.NewConsumer(db),
		outbox:   event.NewOutbox(db),
		config:   config,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /reports/team-overview", handler.teamOverview)
	mux.HandleFunc("GET /reports/sales-overview", handler.salesOverview)
	mux.HandleFunc("POST /internal/projections", handler.upsertProjection)
	return mux
}

func (h *ReportingHandler) teamOverview(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	if actor.Role == "Sales" || actor.Role == "" {
		if err := h.appendReportAccessDenied(r, "team-overview", actor); err != nil {
			writeError(w, http.StatusServiceUnavailable, "AUDIT_LOG_FAILED", "system", "Audit log failed.")
			return
		}
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
	if !h.verifyServiceToken(r, intentProjectionIngest) {
		writeError(w, http.StatusUnauthorized, "SERVICE_AUTH_FAILED", "authentication", "Service authentication failed.")
		return
	}
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

func (h *ReportingHandler) verifyServiceToken(r *http.Request, intent string) bool {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}
	claims, err := authz.VerifyServiceToken(strings.TrimPrefix(authHeader, "Bearer "), authz.VerifyOptions{
		Secret:   h.config.ServiceTokenSecret,
		Audience: h.config.ServiceID,
		Intent:   intent,
		Now:      time.Now().UTC(),
	})
	if err != nil {
		return false
	}
	return r.Header.Get("X-Service-Id") == claims.Issuer && r.Header.Get("X-Intent") == intent
}

func (h *ReportingHandler) appendReportAccessDenied(r *http.Request, reportType string, actor actorContext) error {
	return h.outbox.AppendReportAccessDenied(r.Context(), event.ReportAccessDeniedInput{
		ActorID:       actor.ID,
		ActorRole:     actor.Role,
		ActorDisplay:  actor.ID,
		ReportType:    reportType,
		CorrelationID: r.Header.Get("X-Correlation-Id"),
	})
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
