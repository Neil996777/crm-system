package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"crm-system/services/commercial/internal/authz"
	"crm-system/services/commercial/internal/domain"
	"crm-system/services/commercial/internal/event"
	"crm-system/services/commercial/internal/repo"
)

const intentContractSignedStatus = "commercial.contract_signed_status"

func (h *CommercialHandler) changeContractStatus(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion     int    `json:"expectedVersion"`
		ToStatus            string `json:"toStatus"`
		SignedEffectiveDate string `json:"signedEffectiveDate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 || strings.TrimSpace(request.ToStatus) == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contract status input is invalid.")
		return
	}
	var requestedSignedDate time.Time
	if strings.TrimSpace(request.SignedEffectiveDate) != "" {
		parsed, err := domain.ParseDate(request.SignedEffectiveDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contract status input is invalid.")
			return
		}
		requestedSignedDate = parsed
	}
	current, err := h.contracts.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if actor.Role == "Sales" && current.OwnerID != actor.ID {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	if current.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	resolvedSignedDate, err := domain.ValidateContractStatusTransition(current, request.ToStatus, requestedSignedDate)
	if errors.Is(err, domain.ErrSignedEffectiveDateRequired) {
		writeError(w, http.StatusBadRequest, "SIGNED_EFFECTIVE_DATE_REQUIRED", "business_rule", "Signed or effective date is required for this contract status.")
		return
	}
	if errors.Is(err, domain.ErrInvalidContractTransition) {
		writeError(w, http.StatusBadRequest, "INVALID_TRANSITION", "business_rule", "The requested contract status transition is not allowed.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contract status input is invalid.")
		return
	}
	updated, err := h.contracts.ChangeStatus(r.Context(), current.ID, request.ExpectedVersion, request.ToStatus, resolvedSignedDate)
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	_ = h.outbox.Append(r.Context(), event.ContractStatusChanged, updated.ID, map[string]any{
		"traceability":        "TASK-019 ACC-010 ACC-014 ACC-022 CIM-022 CIM-024 PIM-SM-005 PIM-INV-017 PIM-INV-021 PIM-BEH-015 PSM-006 CONTRACT-009 CONTRACT-010 FLOW-004",
		"actorId":             actor.ID,
		"contractId":          updated.ID,
		"opportunityId":       updated.OpportunityID,
		"fromStatus":          current.Status,
		"toStatus":            updated.Status,
		"signedEffectiveDate": optionalDate(updated.SignedEffectiveDate),
	})
	writeJSON(w, http.StatusOK, contractDTO(updated))
}

func (h *CommercialHandler) getContractSignedStatus(w http.ResponseWriter, r *http.Request) {
	if !h.verifyServiceToken(r, intentContractSignedStatus) {
		writeError(w, http.StatusUnauthorized, "SERVICE_AUTH_FAILED", "authentication", "Service authentication failed.")
		return
	}
	contract, err := h.contracts.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"contractId":          contract.ID,
		"opportunityId":       contract.OpportunityID,
		"status":              contract.Status,
		"signed":              contract.Status == domain.ContractStatusSigned,
		"signedEffectiveDate": optionalDate(contract.SignedEffectiveDate),
		"version":             contract.Version,
	})
}

func (h *CommercialHandler) verifyServiceToken(r *http.Request, intent string) bool {
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

func optionalDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return domain.FormatDate(value)
}
