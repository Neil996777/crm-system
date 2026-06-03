package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"crm-system/services/opportunity/internal/domain"
	"crm-system/services/opportunity/internal/event"
	"crm-system/services/opportunity/internal/repo"
)

func (h *OpportunityHandler) changeStage(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		ToStage         string `json:"toStage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 || request.ToStage == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The stage transition input is invalid.")
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
		writeError(w, http.StatusBadRequest, "TERMINAL_RECORD_READ_ONLY", "business_rule", "Terminal opportunity records cannot change stage.")
		return
	}
	if current.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err := domain.ValidateStageTransition(current.Stage, request.ToStage); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_TRANSITION", "business_rule", "The requested stage transition is not allowed.")
		return
	}
	var updated domain.Opportunity
	err = h.inTransaction(r.Context(), func(txRepo *repo.OpportunityRepo, txOutbox *event.Outbox) error {
		var err error
		updated, err = txRepo.ChangeStage(r.Context(), current.ID, request.ExpectedVersion, request.ToStage)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.OpportunityStageChanged, updated.ID, map[string]any{
			"traceability":   "TASK-014 ACC-008 PIM-SM-002 PIM-INV-006 PIM-BEH-010 PSM-004 CONTRACT-007 CONTRACT-008",
			"actorId":        actor.ID,
			"actorRole":      actor.Role,
			"actorDisplay":   actor.ID,
			"opportunityId":  updated.ID,
			"ownerId":        updated.OwnerID,
			"fromStage":      current.Stage,
			"toStage":        updated.Stage,
			"expectedAmount": updated.ExpectedAmount,
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
	writeJSON(w, http.StatusOK, opportunityDTO(updated))
}
