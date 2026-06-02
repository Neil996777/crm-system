package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	repo       *repo.LeadRepo
	duplicates *repo.DuplicateRepo
	outbox     *event.Outbox
	conversion *client.ConversionClient
	audit      *client.AuditClient
}

func NewLeadServer(db *sql.DB, config Config) http.Handler {
	handler := &LeadHandler{
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
	if request.ProceedWarningToken != "" {
		candidate := domain.DuplicateCandidate{TargetType: "lead", CompanyName: lead.CompanyName, Email: lead.Email, Phone: lead.Phone}
		signature, _ := domain.DuplicateSignature(candidate)
		if err := h.duplicates.ConsumeToken(r.Context(), request.ProceedWarningToken, candidate.TargetType, actor.ID, signature); err != nil {
			writeDuplicateTokenError(w, err)
			return
		}
	}
	created, err := h.repo.Create(r.Context(), lead)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The lead input is invalid.")
		return
	}
	_ = h.outbox.Append(r.Context(), event.LeadCreated, created.ID, map[string]any{
		"traceability": "TASK-007 ACC-003 PIM-BEH-004 CONTRACT-004",
		"actorId":      actor.ID,
		"leadId":       created.ID,
		"ownerId":      created.OwnerID,
		"status":       created.Status,
	})
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
	before, updated, err := h.repo.TransferOwner(r.Context(), r.PathValue("id"), request.ExpectedVersion, request.NewOwnerID)
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The owner transfer input is invalid.")
		return
	}
	_ = h.outbox.Append(r.Context(), event.LeadOwnerChanged, updated.ID, map[string]any{
		"traceability": "TASK-007 ACC-003 ACC-014 PIM-BEH-005 CONTRACT-004 EVT-OWNER-CHANGED",
		"actorId":      actor.ID,
		"oldOwnerId":   before.OwnerID,
		"newOwnerId":   updated.OwnerID,
		"reason":       request.Reason,
	})
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
	archived, err := h.repo.Archive(r.Context(), lead.ID, request.ExpectedVersion, actor.ID, request.Reason)
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	_ = h.outbox.Append(r.Context(), event.LeadArchived, archived.ID, map[string]any{
		"traceability": "TASK-032 ACC-002 ACC-014 CIM-037 CIM-PROC-020 PIM-020 PIM-SM-010 PIM-BEH-024 PSM-002 FLOW-010 TEST-ARCHIVE",
		"actorId":      actor.ID,
		"leadId":       archived.ID,
		"reason":       request.Reason,
	})
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
