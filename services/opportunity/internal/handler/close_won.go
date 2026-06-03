package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"crm-system/services/opportunity/internal/authz"
	"crm-system/services/opportunity/internal/domain"
	"crm-system/services/opportunity/internal/event"
	"crm-system/services/opportunity/internal/repo"
)

func (h *OpportunityHandler) closeWon(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		ContractID      string `json:"contractId"`
		CloseDate       string `json:"closeDate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The close-won input is invalid.")
		return
	}
	closeDate, err := domain.ParseCloseDate(request.CloseDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The close-won input is invalid.")
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
	if !domain.CanEditOpportunity(actor.ID, actor.Role, current, current.OwnerID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	if err := domain.EnsureOpenForMutation(current); err != nil {
		writeError(w, http.StatusBadRequest, "TERMINAL_RECORD_READ_ONLY", "business_rule", "Terminal opportunity records cannot be closed again.")
		return
	}
	if current.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	contract, err := h.commercial.ContractSignedStatus(r.Context(), request.ContractID)
	if err != nil || contract.OpportunityID != current.ID || !contract.Signed {
		if errors.Is(err, authz.ErrCommercialUnavailable) {
			writeError(w, http.StatusBadRequest, "EARLY_WON_BLOCKED", "business_rule", "Won requires a Signed related contract.")
			return
		}
		writeError(w, http.StatusBadRequest, "EARLY_WON_BLOCKED", "business_rule", "Won requires a Signed related contract.")
		return
	}
	closed, err := domain.CloseWon(current, contract.Signed, request.ContractID)
	if errors.Is(err, domain.ErrEarlyWonBlocked) {
		writeError(w, http.StatusBadRequest, "EARLY_WON_BLOCKED", "business_rule", "Won requires a Signed related contract.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	closed = domain.ApplyClosureDates(closed, closeDate, time.Now().UTC())
	var updated domain.Opportunity
	err = h.inTransaction(r.Context(), func(txRepo *repo.OpportunityRepo, txOutbox *event.Outbox) error {
		var err error
		updated, err = txRepo.Close(r.Context(), current.ID, request.ExpectedVersion, closed)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.OpportunityClosedWon, updated.ID, map[string]any{
			"traceability":  "TASK-015 ACC-013 CIM-017 CIM-PROC-011 PIM-SM-009 PIM-INV-035 PIM-BEH-011 PSM-004 CONTRACT-007 CONTRACT-008 FLOW-004 DEC-017 DEC-019",
			"actorId":       actor.ID,
			"actorRole":     actor.Role,
			"actorDisplay":  actor.ID,
			"opportunityId": updated.ID,
			"contractId":    updated.WonContractID,
			"stage":         updated.Stage,
			"closeDate":     domain.FormatCloseDate(updated.CloseDate),
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
	writeJSON(w, http.StatusOK, closeResultDTO(updated))
}

func closeResultDTO(opportunity domain.Opportunity) map[string]any {
	return map[string]any{
		"opportunityId": opportunity.ID,
		"status":        opportunity.Stage,
		"closedAt":      opportunity.ClosedAt.UTC().Format(time.RFC3339),
		"version":       opportunity.Version,
	}
}
