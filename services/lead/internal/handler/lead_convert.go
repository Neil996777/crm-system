package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"crm-system/services/lead/internal/client"
	"crm-system/services/lead/internal/domain"
	"crm-system/services/lead/internal/event"
	"crm-system/services/lead/internal/repo"
)

func (h *LeadHandler) convertLead(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		IDempotencyKey string `json:"idempotencyKey"`
		Target         struct {
			AccountInput     client.AccountInput     `json:"accountInput"`
			OpportunityInput client.OpportunityInput `json:"opportunityInput"`
		} `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.IDempotencyKey == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The conversion input is invalid.")
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
	if current.Status == domain.StatusConverted {
		if current.ConversionIDempotencyKey == request.IDempotencyKey {
			writeJSON(w, http.StatusOK, conversionDTO(current))
			return
		}
		writeError(w, http.StatusConflict, "LEAD_ALREADY_CONVERTED", "conflict", "The lead has already been converted.")
		return
	}
	if current.Status != domain.StatusValid {
		writeError(w, http.StatusConflict, "INVALID_LEAD_STATE", "conflict", "The lead cannot be converted in its current state.")
		return
	}
	accountInput := request.Target.AccountInput
	if accountInput.OwnerID == "" {
		accountInput.OwnerID = current.OwnerID
	}
	account, err := h.conversion.CreateAccount(r.Context(), client.Actor{ID: actor.ID, Role: actor.Role}, accountInput)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "dependency", "A required service is unavailable.")
		return
	}
	opportunityInput := request.Target.OpportunityInput
	opportunityInput.CustomerID = account.ID
	if opportunityInput.OwnerID == "" {
		opportunityInput.OwnerID = current.OwnerID
	}
	opportunity, err := h.conversion.CreateOpportunity(r.Context(), client.Actor{ID: actor.ID, Role: actor.Role}, opportunityInput)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "dependency", "A required service is unavailable.")
		return
	}
	converted, err := domain.Convert(current, request.IDempotencyKey, account.ID, opportunity.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The conversion input is invalid.")
		return
	}
	err = h.inTransaction(r.Context(), func(txLeads *repo.LeadRepo, _ *repo.DuplicateRepo, txOutbox *event.Outbox) error {
		var err error
		converted, err = txLeads.Convert(r.Context(), current.ID, current.Version, converted)
		if err != nil {
			return err
		}
		return appendLeadOutbox(r.Context(), txOutbox, event.LeadConverted, converted.ID, map[string]any{
			"traceability":  "TASK-008 ACC-004 ACC-005 ACC-007 ACC-014 FLOW-002 PIM-SM-001 PIM-INV-003 CONTRACT-003 CONTRACT-004",
			"actorId":       actor.ID,
			"leadId":        converted.ID,
			"accountId":     converted.ConvertedAccountID,
			"opportunityId": converted.ConvertedOpportunityID,
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
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The conversion input is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, conversionDTO(converted))
}

func conversionDTO(lead domain.Lead) map[string]any {
	return map[string]any{
		"leadId":        lead.ID,
		"accountId":     lead.ConvertedAccountID,
		"contactIds":    []string{},
		"opportunityId": lead.ConvertedOpportunityID,
		"status":        lead.Status,
	}
}
