package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"crm-system/services/audit-history/internal/authz"
	"crm-system/services/audit-history/internal/domain"
	"crm-system/services/audit-history/internal/repo"
)

type Config struct {
	ServiceID          string
	ServiceTokenSecret []byte
}

type AuditHandler struct {
	events *repo.EventRepo
	config Config
}

type appendEventRequest struct {
	EventUID           string         `json:"eventUid"`
	EventID            string         `json:"eventId"`
	EventVersion       int            `json:"eventVersion"`
	Surfaces           []string       `json:"surfaces"`
	Action             string         `json:"action"`
	ResourceType       string         `json:"resourceType"`
	ResourceID         string         `json:"resourceId"`
	ParentResourceType string         `json:"parentResourceType"`
	ParentResourceID   string         `json:"parentResourceId"`
	Result             string         `json:"result"`
	ReasonCode         string         `json:"reasonCode"`
	BeforeSummary      map[string]any `json:"beforeSummary"`
	AfterSummary       map[string]any `json:"afterSummary"`
	DiffClassification string         `json:"diffClassification"`
	ScopeSummary       string         `json:"scopeSummary"`
	SafeSummary        string         `json:"safeSummary"`
	CorrelationID      string         `json:"correlationId"`
	CausationID        string         `json:"causationId"`
	AcceptanceIDs      []string       `json:"acceptanceIds"`
}

func NewAuditServer(db *sql.DB, config Config) http.Handler {
	if config.ServiceID == "" {
		config.ServiceID = "audit-history"
	}
	handler := &AuditHandler{events: repo.NewEventRepo(db), config: config}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /internal/events/append", handler.appendEvent)
	mux.HandleFunc("GET /records/{type}/{id}/history", handler.recordHistory)
	mux.HandleFunc("GET /operation-log", handler.operationLog)
	return mux
}

func (h *AuditHandler) appendEvent(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.verifyServiceToken(r, "audit.append")
	if !ok {
		writeErrorCode(w, http.StatusUnauthorized, "SERVICE_AUTH_FAILED", "authentication", "Service authentication failed.")
		return
	}
	var request appendEventRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logAppendRejection("json_decode", claims.Issuer, appendEventRequest{}, r, nil, []string{"json"})
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if missing := missingAppendBodyFields(request); len(missing) > 0 {
		logAppendRejection("body_required_fields", claims.Issuer, request, r, missing, nil)
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	event := domain.NewEvent()
	if strings.TrimSpace(request.EventUID) != "" {
		event.EventUID = strings.TrimSpace(request.EventUID)
	}
	event.EventID = request.EventID
	event.EventVersion = request.EventVersion
	event.ProducerService = claims.Issuer
	event.Surfaces = request.Surfaces
	event.ActorUserID = r.Header.Get("X-Actor-User-Id")
	event.ActorRole = r.Header.Get("X-Actor-Role")
	event.ActorDisplay = r.Header.Get("X-Actor-Display")
	event.Action = request.Action
	event.ResourceType = request.ResourceType
	event.ResourceID = request.ResourceID
	event.ParentResourceType = request.ParentResourceType
	event.ParentResourceID = request.ParentResourceID
	event.Result = request.Result
	event.ReasonCode = request.ReasonCode
	event.BeforeSummary = request.BeforeSummary
	event.AfterSummary = request.AfterSummary
	event.DiffClassification = request.DiffClassification
	event.ScopeSummary = request.ScopeSummary
	event.SafeSummary = request.SafeSummary
	event.CorrelationID = firstNonEmpty(request.CorrelationID, r.Header.Get("X-Correlation-Id"))
	event.CausationID = request.CausationID
	event.AcceptanceIDs = request.AcceptanceIDs
	if missing := missingActorHeaderFields(event); len(missing) > 0 {
		logAppendRejection("actor_headers", claims.Issuer, request, r, missing, nil)
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	stored, err := h.events.Append(r.Context(), event)
	if err != nil {
		logAppendRejection("repo_append", claims.Issuer, request, r, nil, []string{"repo_append_failed"})
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, eventDTO(stored))
}

func (h *AuditHandler) recordHistory(w http.ResponseWriter, r *http.Request) {
	actorID := r.URL.Query().Get("actorId")
	actorRole := r.URL.Query().Get("actorRole")
	ownerID := r.URL.Query().Get("ownerId")
	assigneeID := r.URL.Query().Get("assigneeId")
	teamID := r.URL.Query().Get("teamId")
	if !canReadRecord(actorID, actorRole, ownerID, assigneeID, teamID) {
		writeErrorCode(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	events, err := h.events.ByRecord(r.Context(), r.PathValue("type"), r.PathValue("id"))
	if err != nil {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": eventsDTO(events)})
}

func (h *AuditHandler) operationLog(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("actorRole") != "Administrator" {
		writeErrorCode(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	events, err := h.events.OperationLog(r.Context())
	if err != nil {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": eventsDTO(events)})
}

func (h *AuditHandler) verifyServiceToken(r *http.Request, intent string) (authz.ServiceTokenClaims, bool) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return authz.ServiceTokenClaims{}, false
	}
	claims, err := authz.VerifyServiceToken(strings.TrimPrefix(authHeader, "Bearer "), authz.VerifyOptions{
		Secret:   h.config.ServiceTokenSecret,
		Audience: h.config.ServiceID,
		Intent:   intent,
		Now:      time.Now().UTC(),
	})
	if err != nil {
		return authz.ServiceTokenClaims{}, false
	}
	return claims, r.Header.Get("X-Service-Id") == claims.Issuer && r.Header.Get("X-Intent") == intent
}

func canReadRecord(actorID, actorRole, ownerID, assigneeID, teamID string) bool {
	switch actorRole {
	case "Administrator":
		return true
	case "Sales Manager":
		return teamID == "" || teamID == "single-team"
	case "Sales":
		return actorID != "" && (ownerID == actorID || assigneeID == actorID)
	default:
		return false
	}
}

func eventDTO(event domain.Event) map[string]any {
	return map[string]any{
		"eventUid":        event.EventUID,
		"eventId":         event.EventID,
		"eventVersion":    event.EventVersion,
		"actorUserId":     event.ActorUserID,
		"actorRole":       event.ActorRole,
		"actorDisplay":    event.ActorDisplay,
		"action":          event.Action,
		"resourceType":    event.ResourceType,
		"resourceId":      event.ResourceID,
		"result":          event.Result,
		"reasonCode":      event.ReasonCode,
		"classification":  event.DiffClassification,
		"retentionPolicy": event.RetentionPolicy,
		"retainUntil":     event.RetainUntil.Format(time.RFC3339),
		"beforeSummary":   domain.MaskedSummary(event.DiffClassification, event.BeforeSummary),
		"afterSummary":    domain.MaskedSummary(event.DiffClassification, event.AfterSummary),
		"safeSummary":     event.SafeSummary,
		"occurredAt":      event.OccurredAt.Format(time.RFC3339),
		"prevHash":        event.PrevHash,
		"eventHash":       event.EventHash,
	}
}

func eventsDTO(events []domain.Event) []map[string]any {
	items := make([]map[string]any, 0, len(events))
	for _, event := range events {
		items = append(items, eventDTO(event))
	}
	return items
}

func writeErrorCode(w http.ResponseWriter, status int, code, category, safeMessage string) {
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
	_ = json.NewEncoder(w).Encode(body)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func missingAppendBodyFields(request appendEventRequest) []string {
	var missing []string
	if strings.TrimSpace(request.EventID) == "" {
		missing = append(missing, "eventId")
	}
	if request.EventVersion == 0 {
		missing = append(missing, "eventVersion")
	}
	if strings.TrimSpace(request.Action) == "" {
		missing = append(missing, "action")
	}
	if strings.TrimSpace(request.ResourceType) == "" {
		missing = append(missing, "resourceType")
	}
	if strings.TrimSpace(request.ResourceID) == "" {
		missing = append(missing, "resourceId")
	}
	if strings.TrimSpace(request.Result) == "" {
		missing = append(missing, "result")
	}
	if strings.TrimSpace(request.SafeSummary) == "" {
		missing = append(missing, "safeSummary")
	}
	if len(request.Surfaces) == 0 {
		missing = append(missing, "surfaces")
	}
	if len(request.AcceptanceIDs) == 0 {
		missing = append(missing, "acceptanceIds")
	}
	return missing
}

func missingActorHeaderFields(event domain.Event) []string {
	var missing []string
	if strings.TrimSpace(event.ActorUserID) == "" {
		missing = append(missing, "actorUserId")
	}
	if strings.TrimSpace(event.ActorRole) == "" {
		missing = append(missing, "actorRole")
	}
	if strings.TrimSpace(event.ActorDisplay) == "" {
		missing = append(missing, "actorDisplay")
	}
	return missing
}

func logAppendRejection(stage, producerService string, request appendEventRequest, r *http.Request, missingFields, invalidFields []string) {
	log.Printf(
		"audit append rejected stage=%s producerService=%s eventUid=%s correlationId=%s missingFields=%s invalidFields=%s",
		safeDiagnosticValue(stage),
		safeDiagnosticValue(producerService),
		safeDiagnosticValue(request.EventUID),
		safeDiagnosticValue(firstNonEmpty(request.CorrelationID, r.Header.Get("X-Correlation-Id"))),
		safeDiagnosticList(missingFields),
		safeDiagnosticList(invalidFields),
	)
}

func safeDiagnosticValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return strings.NewReplacer(" ", "_", "\n", "_", "\r", "_", "\t", "_").Replace(value)
}

func safeDiagnosticList(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		cleaned = append(cleaned, safeDiagnosticValue(value))
	}
	return strings.Join(cleaned, ",")
}
