package handler

import (
	"context"
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
	"github.com/jackc/pgx/v5/pgconn"
)

const intentCreateOpportunityForLeadConversion = "opportunity.create_for_lead_conversion"

type Config struct {
	ServiceID          string
	ServiceTokenSecret []byte
	CommercialBaseURL  string
}

type OpportunityHandler struct {
	db         *sql.DB
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
		db:     db,
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

func (h *OpportunityHandler) inTransaction(ctx context.Context, fn func(*repo.OpportunityRepo, *event.Outbox) error) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(repo.NewOpportunityRepoTx(tx), event.NewOutboxTx(tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func writeOutboxFailure(w http.ResponseWriter) {
	writeError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "dependency", "The audit event could not be persisted.")
}

func (h *OpportunityHandler) createOpportunityFromLeadConversion(w http.ResponseWriter, r *http.Request) {
	if !h.verifyServiceToken(r, intentCreateOpportunityForLeadConversion) {
		writeError(w, http.StatusUnauthorized, "SERVICE_AUTH_FAILED", "authentication", "Service authentication failed.")
		return
	}
	h.createOpportunityCommand(w, r, true)
}

func (h *OpportunityHandler) createOpportunity(w http.ResponseWriter, r *http.Request) {
	h.createOpportunityCommand(w, r, false)
}

func (h *OpportunityHandler) createOpportunityCommand(w http.ResponseWriter, r *http.Request, requireIDempotency bool) {
	actor := actorFromRequest(r)
	request, err := decodeOpportunityInput(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The opportunity input is invalid.")
		return
	}
	if requireIDempotency {
		request.IDempotencyKey = strings.TrimSpace(request.IDempotencyKey)
		if request.IDempotencyKey == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The opportunity input is invalid.")
			return
		}
		existing, err := h.repo.FindByLeadConversionIDempotencyKey(r.Context(), request.IDempotencyKey)
		if err == nil {
			writeJSON(w, http.StatusOK, opportunityDTO(existing))
			return
		}
		if !errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
			return
		}
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
	var created domain.Opportunity
	if err := h.inTransaction(r.Context(), func(txRepo *repo.OpportunityRepo, txOutbox *event.Outbox) error {
		var err error
		created, err = txRepo.CreateForLeadConversion(r.Context(), opportunity, request.IDempotencyKey)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.OpportunityCreated, created.ID, map[string]any{
			"traceability":      "TASK-013 ACC-007 PIM-007 PIM-SM-002 PIM-BEH-009 PSM-004 CONTRACT-007 CONTRACT-008",
			"actorId":           actor.ID,
			"actorRole":         actor.Role,
			"actorDisplay":      actor.ID,
			"opportunityId":     created.ID,
			"customerId":        created.CustomerID,
			"ownerId":           created.OwnerID,
			"stage":             created.Stage,
			"expectedAmount":    created.ExpectedAmount,
			"expectedCloseDate": domain.FormatCloseDate(created.ExpectedCloseDate),
		})
	}); err != nil {
		if requireIDempotency && isLeadConversionIDempotencyConflict(err) {
			existing, findErr := h.repo.FindByLeadConversionIDempotencyKey(r.Context(), request.IDempotencyKey)
			if findErr == nil {
				writeJSON(w, http.StatusOK, opportunityDTO(existing))
				return
			}
		}
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The opportunity input is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, opportunityDTO(created))
}

func isLeadConversionIDempotencyConflict(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "lead_conversion_idempotency_key")
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
	var saved domain.Opportunity
	err = h.inTransaction(r.Context(), func(txRepo *repo.OpportunityRepo, txOutbox *event.Outbox) error {
		var err error
		saved, err = txRepo.Update(r.Context(), current.ID, request.ExpectedVersion, updated)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.OpportunityUpdated, saved.ID, map[string]any{
			"traceability":  "TASK-013 ACC-007 PIM-BEH-009 CONTRACT-007 CONTRACT-008 CONTRACT-020",
			"actorId":       actor.ID,
			"actorRole":     actor.Role,
			"actorDisplay":  actor.ID,
			"opportunityId": saved.ID,
			"version":       saved.Version,
			"stage":         saved.Stage,
		})
	})
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeOutboxFailure(w)
		return
	}
	writeJSON(w, http.StatusOK, opportunityDTO(saved))
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
	var archived domain.Opportunity
	err = h.inTransaction(r.Context(), func(txRepo *repo.OpportunityRepo, txOutbox *event.Outbox) error {
		var err error
		archived, err = txRepo.Archive(r.Context(), opportunity.ID, request.ExpectedVersion, actor.ID, request.Reason)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.OpportunityArchived, archived.ID, map[string]any{
			"traceability":  "TASK-032 ACC-002 ACC-014 CIM-037 CIM-PROC-020 PIM-020 PIM-SM-010 PIM-BEH-024 PSM-004 FLOW-010 TEST-ARCHIVE",
			"actorId":       actor.ID,
			"actorRole":     actor.Role,
			"actorDisplay":  actor.ID,
			"opportunityId": archived.ID,
			"reason":        request.Reason,
		})
	})
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeOutboxFailure(w)
		return
	}
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
	IDempotencyKey    string `json:"idempotencyKey"`
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
