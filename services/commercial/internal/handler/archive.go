package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"crm-system/services/commercial/internal/domain"
	"crm-system/services/commercial/internal/event"
	"crm-system/services/commercial/internal/repo"
)

type archiveObligation struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Service     string `json:"service"`
	Status      string `json:"status"`
	Blocking    bool   `json:"blocking"`
	SafeMessage string `json:"safeMessage"`
}

func (h *CommercialHandler) contractArchiveEligibility(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	contract, err := h.authorizedArchiveContract(r, actor)
	if err != nil {
		writeCommercialLookupError(w, err)
		return
	}
	obligations := contractArchiveObligations(contract)
	writeJSON(w, http.StatusOK, map[string]any{
		"resourceType":  "contract",
		"resourceId":    contract.ID,
		"canArchive":    len(obligations) == 0,
		"recordVersion": contract.Version,
		"obligations":   obligations,
	})
}

func (h *CommercialHandler) archiveContract(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		Reason          string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 || strings.TrimSpace(request.Reason) == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The archive input is invalid.")
		return
	}
	contract, err := h.authorizedArchiveContract(r, actor)
	if err != nil {
		writeCommercialLookupError(w, err)
		return
	}
	if contract.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	obligations := contractArchiveObligations(contract)
	if len(obligations) > 0 {
		writeArchiveBlocked(w, obligations)
		return
	}
	archived, err := h.contracts.Archive(r.Context(), contract.ID, request.ExpectedVersion, actor.ID, request.Reason)
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	_ = h.outbox.Append(r.Context(), event.ContractArchived, archived.ID, map[string]any{
		"traceability": "TASK-032 ACC-002 ACC-014 CIM-037 CIM-PROC-020 PIM-020 PIM-SM-010 PIM-BEH-024 PSM-006 FLOW-010 TEST-ARCHIVE",
		"actorId":      actor.ID,
		"contractId":   archived.ID,
		"reason":       request.Reason,
	})
	writeJSON(w, http.StatusOK, contractDTO(archived))
}

func (h *CommercialHandler) archivePaymentPlan(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	if !domain.CanArchiveCommercialRecord(actor.Role) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		Reason          string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 || strings.TrimSpace(request.Reason) == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The archive input is invalid.")
		return
	}
	plan, err := h.payments.FindPlan(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if plan.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if plan.Status != domain.PaymentStatusPaid {
		writeArchiveBlocked(w, []archiveObligation{{
			Type:        "unpaid_payment",
			ID:          plan.ID,
			Service:     "commercial-service",
			Status:      plan.Status,
			Blocking:    true,
			SafeMessage: "Unpaid payment plans must be paid before archive.",
		}})
		return
	}
	archived, err := h.payments.ArchivePlan(r.Context(), plan.ID, request.ExpectedVersion, actor.ID, request.Reason)
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	_ = h.outbox.Append(r.Context(), event.PaymentPlanArchived, archived.ID, map[string]any{
		"traceability":  "TASK-032 ACC-002 ACC-014 CIM-037 CIM-PROC-020 PIM-020 PIM-SM-010 PIM-BEH-024 PSM-007 FLOW-010 TEST-INV-ARCHIVEBLOCK-001",
		"actorId":       actor.ID,
		"paymentPlanId": archived.ID,
		"contractId":    archived.ContractID,
		"reason":        request.Reason,
	})
	writeJSON(w, http.StatusOK, paymentPlanDTO(archived))
}

func (h *CommercialHandler) authorizedArchiveContract(r *http.Request, actor actorContext) (domain.Contract, error) {
	if !domain.CanArchiveCommercialRecord(actor.Role) {
		return domain.Contract{}, errPermissionDenied
	}
	contract, err := h.contracts.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		return domain.Contract{}, repo.ErrNotFound
	}
	if err != nil {
		return domain.Contract{}, err
	}
	return contract, nil
}

func contractArchiveObligations(contract domain.Contract) []archiveObligation {
	if contract.Status != domain.ContractStatusPendingSignature {
		return []archiveObligation{}
	}
	return []archiveObligation{{
		Type:        "pending_signature_contract",
		ID:          contract.ID,
		Service:     "commercial-service",
		Status:      contract.Status,
		Blocking:    true,
		SafeMessage: "Pending signature contracts must be signed or terminated before archive.",
	}}
}

func writeArchiveBlocked(w http.ResponseWriter, obligations []archiveObligation) {
	writeJSON(w, http.StatusConflict, map[string]any{
		"error": map[string]any{
			"code":        "ARCHIVE_BLOCKED_ACTIVE_OBLIGATION",
			"category":    "business_rule",
			"safeMessage": "Active obligations must be resolved before archive.",
		},
		"obligations": obligations,
	})
}
