package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"crm-system/services/opportunity/internal/authz"
	"crm-system/services/opportunity/internal/domain"
	"crm-system/services/opportunity/internal/event"
	"crm-system/services/opportunity/internal/repo"
)

const intentCreateOpportunityForLeadConversion = "opportunity.create_for_lead_conversion"

type Config struct {
	ServiceID          string
	ServiceTokenSecret []byte
	CommercialBaseURL  string
}

type OpportunityHandler struct {
	repo       *repo.OpportunityRepo
	outbox     *event.Outbox
	commercial authz.CommercialClient
	config     Config
}

func NewOpportunityServer(db *sql.DB, config Config) http.Handler {
	if config.ServiceID == "" {
		config.ServiceID = "opportunity"
	}
	handler := &OpportunityHandler{
		repo:   repo.NewOpportunityRepo(db),
		outbox: event.NewOutbox(db),
		commercial: authz.CommercialClient{
			BaseURL:            config.CommercialBaseURL,
			ServiceID:          config.ServiceID,
			ServiceTokenSecret: config.ServiceTokenSecret,
		},
		config: config,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /opportunities", handler.createOpportunity)
	mux.HandleFunc("POST /internal/opportunities", handler.createOpportunityFromLeadConversion)
	mux.HandleFunc("POST /opportunities/{id}/stage", handler.changeStage)
	mux.HandleFunc("POST /opportunities/{id}/close-won", handler.closeWon)
	mux.HandleFunc("POST /opportunities/{id}/close-lost", handler.closeLost)
	mux.HandleFunc("GET /opportunities/{id}/archive-eligibility", handler.archiveEligibility)
	mux.HandleFunc("POST /opportunities/{id}/archive", handler.archiveOpportunity)
	mux.HandleFunc("PATCH /opportunities/{id}", handler.updateOpportunity)
	mux.HandleFunc("GET /opportunities", handler.listOpportunities)
	mux.HandleFunc("GET /opportunities/{id}", handler.getOpportunity)
	return mux
}

func (h *OpportunityHandler) createOpportunityFromLeadConversion(w http.ResponseWriter, r *http.Request) {
	if !h.verifyServiceToken(r, intentCreateOpportunityForLeadConversion) {
		writeError(w, http.StatusUnauthorized, "SERVICE_AUTH_FAILED", "authentication", "Service authentication failed.")
		return
	}
	h.createOpportunity(w, r)
}

func (h *OpportunityHandler) createOpportunity(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	request, err := decodeOpportunityInput(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The opportunity input is invalid.")
		return
	}
	if !domain.CanCreateOpportunity(actor.ID, actor.Role, request.OwnerID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	opportunity, err := domain.NewOpportunity(request.toDomain())
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The opportunity input is invalid.")
		return
	}
	created, err := h.repo.Create(r.Context(), opportunity)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The opportunity input is invalid.")
		return
	}
	_ = h.outbox.Append(r.Context(), event.OpportunityCreated, created.ID, map[string]any{
		"traceability":      "TASK-013 ACC-007 PIM-007 PIM-SM-002 PIM-BEH-009 PSM-004 CONTRACT-007 CONTRACT-008",
		"actorId":           actor.ID,
		"opportunityId":     created.ID,
		"customerId":        created.CustomerID,
		"ownerId":           created.OwnerID,
		"stage":             created.Stage,
		"expectedAmount":    created.ExpectedAmount,
		"expectedCloseDate": domain.FormatCloseDate(created.ExpectedCloseDate),
	})
	writeJSON(w, http.StatusCreated, opportunityDTO(created))
}

func (h *OpportunityHandler) updateOpportunity(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	request, err := decodeOpportunityInput(r)
	if err != nil || request.ExpectedVersion == 0 {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The opportunity input is invalid.")
		return
	}
	current, err := h.repo.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if !domain.CanEditOpportunity(actor.ID, actor.Role, current, request.OwnerID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	if err := domain.EnsureOpenForMutation(current); err != nil {
		writeError(w, http.StatusBadRequest, "TERMINAL_RECORD_READ_ONLY", "business_rule", "Terminal opportunity records cannot be edited.")
		return
	}
	if current.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	updated, err := domain.UpdateOpportunity(current, request.CustomerID, request.OwnerID, request.Stage, request.ExpectedAmount, request.closeDate.Time, request.Title)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The opportunity input is invalid.")
		return
	}
	updated, err = h.repo.Update(r.Context(), current.ID, request.ExpectedVersion, updated)
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The opportunity input is invalid.")
		return
	}
	_ = h.outbox.Append(r.Context(), event.OpportunityUpdated, updated.ID, map[string]any{
		"traceability":  "TASK-013 ACC-007 PIM-BEH-009 CONTRACT-007 CONTRACT-008 CONTRACT-020",
		"actorId":       actor.ID,
		"opportunityId": updated.ID,
		"version":       updated.Version,
		"stage":         updated.Stage,
	})
	writeJSON(w, http.StatusOK, opportunityDTO(updated))
}

func (h *OpportunityHandler) archiveEligibility(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	opportunity, err := h.authorizedArchiveOpportunity(r, actor)
	if err != nil {
		writeOpportunityLookupError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"resourceType":  "opportunity",
		"resourceId":    opportunity.ID,
		"canArchive":    true,
		"recordVersion": opportunity.Version,
		"obligations":   []any{},
	})
}

func (h *OpportunityHandler) archiveOpportunity(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		Reason          string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 || strings.TrimSpace(request.Reason) == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The archive input is invalid.")
		return
	}
	opportunity, err := h.authorizedArchiveOpportunity(r, actor)
	if err != nil {
		writeOpportunityLookupError(w, err)
		return
	}
	if opportunity.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	archived, err := h.repo.Archive(r.Context(), opportunity.ID, request.ExpectedVersion, actor.ID, request.Reason)
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	_ = h.outbox.Append(r.Context(), event.OpportunityArchived, archived.ID, map[string]any{
		"traceability":  "TASK-032 ACC-002 ACC-014 CIM-037 CIM-PROC-020 PIM-020 PIM-SM-010 PIM-BEH-024 PSM-004 FLOW-010 TEST-ARCHIVE",
		"actorId":       actor.ID,
		"opportunityId": archived.ID,
		"reason":        request.Reason,
	})
	writeJSON(w, http.StatusOK, opportunityDTO(archived))
}

func (h *OpportunityHandler) authorizedArchiveOpportunity(r *http.Request, actor actorContext) (domain.Opportunity, error) {
	if !domain.CanArchiveOpportunity(actor.Role) {
		return domain.Opportunity{}, errPermissionDenied
	}
	opportunity, err := h.repo.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		return domain.Opportunity{}, repo.ErrNotFound
	}
	if err != nil {
		return domain.Opportunity{}, err
	}
	if !domain.CanReadOpportunity(actor.ID, actor.Role, opportunity) {
		return domain.Opportunity{}, errPermissionDenied
	}
	return opportunity, nil
}

type opportunityInput struct {
	ExpectedVersion   int    `json:"expectedVersion"`
	CustomerID        string `json:"customerId"`
	OwnerID           string `json:"ownerId"`
	Stage             string `json:"stage"`
	ExpectedAmount    string `json:"expectedAmount"`
	ExpectedCloseDate string `json:"expectedCloseDate"`
	Title             string `json:"title"`
	closeDate         sql.NullTime
}

func decodeOpportunityInput(r *http.Request) (opportunityInput, error) {
	var request opportunityInput
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return opportunityInput{}, err
	}
	closeDate, err := domain.ParseCloseDate(request.ExpectedCloseDate)
	if err != nil {
		return opportunityInput{}, err
	}
	request.closeDate = sql.NullTime{Time: closeDate, Valid: true}
	return request, nil
}

func (i opportunityInput) toDomain() domain.Opportunity {
	return domain.Opportunity{
		CustomerID:        i.CustomerID,
		OwnerID:           i.OwnerID,
		Stage:             i.Stage,
		ExpectedAmount:    i.ExpectedAmount,
		ExpectedCloseDate: i.closeDate.Time,
		Title:             i.Title,
	}
}

func (h *OpportunityHandler) verifyServiceToken(r *http.Request, intent string) bool {
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
