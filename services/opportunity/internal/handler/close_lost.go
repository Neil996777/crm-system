package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"crm-system/services/opportunity/internal/domain"
	"crm-system/services/opportunity/internal/event"
	"crm-system/services/opportunity/internal/repo"
)

func (h *OpportunityHandler) closeLost(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		CloseDate       string `json:"closeDate"`
		LostReason      struct {
			Code   string `json:"code"`
			Detail string `json:"detail"`
		} `json:"lostReason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The close-lost input is invalid.")
		return
	}
	closeDate, err := domain.ParseCloseDate(request.CloseDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The close-lost input is invalid.")
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
	closed, err := domain.CloseLost(current, request.LostReason.Code, request.LostReason.Detail)
	if errors.Is(err, domain.ErrLostReasonRequired) {
		writeError(w, http.StatusBadRequest, "LOST_REASON_REQUIRED", "business_rule", "Lost reason is required.")
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
		return txOutbox.Append(r.Context(), event.OpportunityClosedLost, updated.ID, map[string]any{
			"traceability":   "TASK-015 ACC-013 CIM-018 CIM-PROC-011 PIM-SM-009 PIM-INV-037 PIM-BEH-011 PSM-004 CONTRACT-007 CONTRACT-008 FLOW-013",
			"actorId":        actor.ID,
			"actorRole":      actor.Role,
			"actorDisplay":   actor.ID,
			"opportunityId":  updated.ID,
			"ownerId":        updated.OwnerID,
			"stage":          updated.Stage,
			"expectedAmount": updated.ExpectedAmount,
			"closeDate":      domain.FormatCloseDate(updated.CloseDate),
			"lostReasonCode": updated.LostReasonCode,
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
