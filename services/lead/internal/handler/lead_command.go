package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"crm-system/services/lead/internal/client"
	"crm-system/services/lead/internal/domain"
	"crm-system/services/lead/internal/event"
	"crm-system/services/lead/internal/repo"
)

type Config struct {
	AccountServiceURL      string
	OpportunityServiceURL  string
	AuditHistoryServiceURL string
	ServiceID              string
	ServiceTokenSecret     []byte
	HTTPClient             *http.Client
}

type LeadHandler struct {
	db         *sql.DB
	repo       *repo.LeadRepo
	duplicates *repo.DuplicateRepo
	outbox     *event.Outbox
	conversion *client.ConversionClient
	audit      *client.AuditClient
}

func NewLeadServer(db *sql.DB, config Config) http.Handler {
	handler := &LeadHandler{
		db:         db,
		repo:       repo.NewLeadRepo(db),
		duplicates: repo.NewDuplicateRepo(db),
		outbox:     event.NewOutbox(db),
		conversion: client.NewConversionClient(client.Config{
			AccountServiceURL:      config.AccountServiceURL,
			OpportunityServiceURL:  config.OpportunityServiceURL,
			AuditHistoryServiceURL: config.AuditHistoryServiceURL,
			ServiceID:              config.ServiceID,
			ServiceTokenSecret:     config.ServiceTokenSecret,
			HTTPClient:             config.HTTPClient,
		}),
		audit: client.NewAuditClient(client.Config{
			AuditHistoryServiceURL: config.AuditHistoryServiceURL,
			ServiceID:              config.ServiceID,
			ServiceTokenSecret:     config.ServiceTokenSecret,
			HTTPClient:             config.HTTPClient,
		}),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /leads", handler.createLead)
	mux.HandleFunc("POST /leads/duplicate-checks", handler.duplicateCheck)
	mux.HandleFunc("POST /duplicate-checks", handler.duplicateCheck)
	mux.HandleFunc("POST /leads/{id}/owner-transfer", handler.transferOwner)
	mux.HandleFunc("POST /leads/{id}/qualify-valid", handler.qualifyValid)
	mux.HandleFunc("POST /leads/{id}/qualify-invalid", handler.qualifyInvalid)
	mux.HandleFunc("POST /leads/{id}/restore-invalid", handler.restoreInvalid)
	mux.HandleFunc("POST /leads/{id}/convert", handler.convertLead)
	mux.HandleFunc("GET /leads/{id}/archive-eligibility", handler.archiveEligibility)
	mux.HandleFunc("POST /leads/{id}/archive", handler.archiveLead)
	mux.HandleFunc("GET /leads", handler.listLeads)
	mux.HandleFunc("GET /leads/{id}", handler.getLead)
	return mux
}

var errOutboxAppendFailed = errors.New("outbox append failed")

func (h *LeadHandler) inTransaction(ctx context.Context, fn func(*repo.LeadRepo, *repo.DuplicateRepo, *event.Outbox) error) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(repo.NewLeadRepoTx(tx), repo.NewDuplicateRepoTx(tx), event.NewOutboxTx(tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func appendLeadOutbox(ctx context.Context, outbox *event.Outbox, eventType, aggregateID string, payload map[string]any) error {
	if err := outbox.Append(ctx, eventType, aggregateID, payload); err != nil {
		return fmt.Errorf("%w: %v", errOutboxAppendFailed, err)
	}
	return nil
}

func writeOutboxFailure(w http.ResponseWriter) {
	writeError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "dependency", "The audit event could not be persisted.")
}

func (h *LeadHandler) createLead(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		LeadName            string `json:"leadName"`
		CompanyName         string `json:"companyName"`
		Email               string `json:"email"`
		Phone               string `json:"phone"`
		Source              string `json:"source"`
		OwnerID             string `json:"ownerId"`
		NeedSummary         string `json:"needSummary"`
		ProceedWarningToken string `json:"proceedWarningToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The lead input is invalid.")
		return
	}
	if !domain.CanCreateLead(actor.ID, actor.Role, request.OwnerID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	lead, err := domain.NewLead(domain.Lead{
		LeadName:    request.LeadName,
		CompanyName: request.CompanyName,
		Email:       request.Email,
		Phone:       request.Phone,
		Source:      request.Source,
		OwnerID:     request.OwnerID,
		NeedSummary: request.NeedSummary,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The lead input is invalid.")
		return
	}
	var created domain.Lead
	err = h.inTransaction(r.Context(), func(txLeads *repo.LeadRepo, txDuplicates *repo.DuplicateRepo, txOutbox *event.Outbox) error {
		if request.ProceedWarningToken != "" {
			candidate := domain.DuplicateCandidate{TargetType: "lead", CompanyName: lead.CompanyName, Email: lead.Email, Phone: lead.Phone}
			signature, _ := domain.DuplicateSignature(candidate)
			if err := txDuplicates.ConsumeToken(r.Context(), request.ProceedWarningToken, candidate.TargetType, actor.ID, signature); err != nil {
				return err
			}
		}
		created, err = txLeads.Create(r.Context(), lead)
		if err != nil {
			return err
		}
		return appendLeadOutbox(r.Context(), txOutbox, event.LeadCreated, created.ID, map[string]any{
			"traceability": "TASK-007 ACC-003 PIM-BEH-004 CONTRACT-004",
			"actorId":      actor.ID,
			"leadId":       created.ID,
			"ownerId":      created.OwnerID,
			"status":       created.Status,
		})
	})
	if errors.Is(err, repo.ErrDuplicateTokenInvalid) || errors.Is(err, repo.ErrDuplicateTokenUsed) {
		writeDuplicateTokenError(w, err)
		return
	}
	if errors.Is(err, errOutboxAppendFailed) {
		writeOutboxFailure(w)
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The lead input is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, leadDTO(created))
}

func (h *LeadHandler) transferOwner(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	if !domain.CanTransferOwner(actor.Role) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		NewOwnerID      string `json:"newOwnerId"`
		Reason          string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 || request.NewOwnerID == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The owner transfer input is invalid.")
		return
	}
	var before domain.Lead
	var updated domain.Lead
	err := h.inTransaction(r.Context(), func(txLeads *repo.LeadRepo, _ *repo.DuplicateRepo, txOutbox *event.Outbox) error {
		var err error
		before, updated, err = txLeads.TransferOwner(r.Context(), r.PathValue("id"), request.ExpectedVersion, request.NewOwnerID)
		if err != nil {
			return err
		}
		return appendLeadOutbox(r.Context(), txOutbox, event.LeadOwnerChanged, updated.ID, map[string]any{
			"traceability": "TASK-007 ACC-003 ACC-014 PIM-BEH-005 CONTRACT-004 EVT-OWNER-CHANGED",
			"actorId":      actor.ID,
			"oldOwnerId":   before.OwnerID,
			"newOwnerId":   updated.OwnerID,
			"reason":       request.Reason,
		})
	})
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if errors.Is(err, errOutboxAppendFailed) {
		writeOutboxFailure(w)
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The owner transfer input is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, leadDTO(updated))
}

func (h *LeadHandler) archiveEligibility(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	lead, err := h.authorizedArchiveLead(r, actor)
	if err != nil {
		writeLeadLookupError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"resourceType":  "lead",
		"resourceId":    lead.ID,
		"canArchive":    true,
		"recordVersion": lead.Version,
		"obligations":   []any{},
	})
}

func (h *LeadHandler) archiveLead(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		Reason          string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 || strings.TrimSpace(request.Reason) == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The archive input is invalid.")
		return
	}
	lead, err := h.authorizedArchiveLead(r, actor)
	if err != nil {
		writeLeadLookupError(w, err)
		return
	}
	if lead.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	var archived domain.Lead
	err = h.inTransaction(r.Context(), func(txLeads *repo.LeadRepo, _ *repo.DuplicateRepo, txOutbox *event.Outbox) error {
		var err error
		archived, err = txLeads.Archive(r.Context(), lead.ID, request.ExpectedVersion, actor.ID, request.Reason)
		if err != nil {
			return err
		}
		return appendLeadOutbox(r.Context(), txOutbox, event.LeadArchived, archived.ID, map[string]any{
			"traceability": "TASK-032 ACC-002 ACC-014 CIM-037 CIM-PROC-020 PIM-020 PIM-SM-010 PIM-BEH-024 PSM-002 FLOW-010 TEST-ARCHIVE",
			"actorId":      actor.ID,
			"leadId":       archived.ID,
			"reason":       request.Reason,
		})
	})
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if errors.Is(err, errOutboxAppendFailed) {
		writeOutboxFailure(w)
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, leadDTO(archived))
}

func (h *LeadHandler) authorizedArchiveLead(r *http.Request, actor actorContext) (domain.Lead, error) {
	if !domain.CanArchiveLead(actor.Role) {
		return domain.Lead{}, errPermissionDenied
	}
	lead, err := h.repo.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		return domain.Lead{}, repo.ErrNotFound
	}
	if err != nil {
		return domain.Lead{}, err
	}
	if !domain.CanReadLead(actor.ID, actor.Role, lead) {
		return domain.Lead{}, errPermissionDenied
	}
	return lead, nil
}
