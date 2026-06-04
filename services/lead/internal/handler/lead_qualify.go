package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"crm-system/services/lead/internal/domain"
	"crm-system/services/lead/internal/event"
	"crm-system/services/lead/internal/repo"
)

func (h *LeadHandler) qualifyValid(w http.ResponseWriter, r *http.Request) {
	h.updateQualification(w, r, func(current domain.Lead, _ string) (domain.Lead, error) {
		return domain.QualifyValid(current)
	})
}

func (h *LeadHandler) qualifyInvalid(w http.ResponseWriter, r *http.Request) {
	h.updateQualification(w, r, func(current domain.Lead, reason string) (domain.Lead, error) {
		return domain.QualifyInvalid(current, reason)
	})
}

func (h *LeadHandler) restoreInvalid(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	if !domain.CanRestoreLead(actor.Role) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	h.updateQualification(w, r, func(current domain.Lead, _ string) (domain.Lead, error) {
		return domain.RestoreInvalid(current)
	})
}

func (h *LeadHandler) updateQualification(w http.ResponseWriter, r *http.Request, transition func(domain.Lead, string) (domain.Lead, error)) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		InvalidReason   string `json:"invalidReason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The qualification input is invalid.")
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
	if !domain.CanQualifyLead(actor.ID, actor.Role, current) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	if current.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	updated, err := transition(current, request.InvalidReason)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The qualification input is invalid.")
		return
	}
	eventType := leadQualificationEventType(updated)
	err = h.inTransaction(r.Context(), func(txLeads *repo.LeadRepo, _ *repo.DuplicateRepo, txOutbox *event.Outbox) error {
		var err error
		updated, err = txLeads.UpdateQualification(r.Context(), current.ID, request.ExpectedVersion, updated)
		if err != nil {
			return err
		}
		return appendLeadOutbox(r.Context(), txOutbox, eventType, updated.ID, map[string]any{
			"traceability":  "TASK-008 ACC-004 PIM-SM-001 PIM-BEH-006 CONTRACT-004",
			"actorId":       actor.ID,
			"actorRole":     actor.Role,
			"actorDisplay":  actor.DisplayName,
			"correlationId": r.Header.Get("X-Correlation-Id"),
			"leadId":        updated.ID,
			"beforeStatus":  current.Status,
			"status":        updated.Status,
			"invalidReason": updated.InvalidReason,
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
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The qualification input is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, leadDTO(updated))
}

func leadQualificationEventType(updated domain.Lead) string {
	if updated.Status == domain.StatusInvalid {
		return event.LeadDisqualified
	}
	return event.LeadQualified
}
