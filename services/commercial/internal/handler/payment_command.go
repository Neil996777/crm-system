package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"crm-system/services/commercial/internal/domain"
	"crm-system/services/commercial/internal/event"
	"crm-system/services/commercial/internal/repo"
)

func (h *CommercialHandler) createPaymentPlan(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	contract, err := h.contracts.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if actor.Role == "Sales" && contract.OwnerID != actor.ID {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	var request struct {
		DueAmount string `json:"dueAmount"`
		DueDate   string `json:"dueDate"`
		Currency  string `json:"currency"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The payment plan input is invalid.")
		return
	}
	plan, err := domain.NewPaymentPlan(contract.ID, request.DueAmount, request.DueDate, request.Currency)
	if errors.Is(err, domain.ErrSingleCurrencyRequired) {
		writeError(w, http.StatusBadRequest, "SINGLE_CURRENCY_REQUIRED", "business_rule", "Payments use the committed single currency.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_AMOUNT", "business_rule", "Payment amount must be greater than zero.")
		return
	}
	created, err := h.payments.CreatePlan(r.Context(), plan)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The payment plan input is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, paymentPlanDTO(created))
}

func (h *CommercialHandler) recordPayment(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	contract, err := h.contracts.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if actor.Role == "Sales" && contract.OwnerID != actor.ID {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	var request struct {
		IdempotencyKey string `json:"idempotencyKey"`
		Amount         string `json:"amount"`
		PaymentDate    string `json:"paymentDate"`
		Note           string `json:"note"`
		Currency       string `json:"currency"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The payment input is invalid.")
		return
	}
	existing, err := h.payments.FindPaymentByKey(r.Context(), contract.ID, request.IdempotencyKey)
	if err == nil {
		writeJSON(w, http.StatusOK, actualPaymentDTO(existing))
		return
	}
	if err != nil && !errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	paidBefore, err := h.payments.PaidTotal(r.Context(), contract.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	payment, err := domain.NewActualPayment(contract, request.IdempotencyKey, request.Amount, request.PaymentDate, request.Note, request.Currency, paidBefore)
	if errors.Is(err, domain.ErrSingleCurrencyRequired) {
		writeError(w, http.StatusBadRequest, "SINGLE_CURRENCY_REQUIRED", "business_rule", "Payments use the committed single currency.")
		return
	}
	if errors.Is(err, domain.ErrOverpaymentBlocked) {
		writeError(w, http.StatusBadRequest, "OVERPAYMENT_BLOCKED", "business_rule", "Payment exceeds the remaining contract amount.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_AMOUNT", "business_rule", "Payment amount must be greater than zero.")
		return
	}
	var created domain.ActualPayment
	if err := h.inTransaction(r.Context(), func(_ *repo.QuoteRepo, _ *repo.ContractRepo, txPayments *repo.PaymentRepo, txOutbox *event.Outbox) error {
		var err error
		created, err = txPayments.CreatePayment(r.Context(), payment)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.PaymentRecorded, created.ID, map[string]any{
			"traceability":    "TASK-020 ACC-011 CIM-026 CIM-PROC-010 PIM-010 PIM-011 PIM-SM-006 PIM-INV-022 PIM-INV-023 PIM-BEH-017 PIM-BEH-018 PSM-007 CONTRACT-009 CONTRACT-010 DEC-019",
			"actorId":         actor.ID,
			"actorRole":       actor.Role,
			"actorDisplay":    actor.ID,
			"contractId":      created.ContractID,
			"paymentId":       created.ID,
			"ownerId":         contract.OwnerID,
			"amount":          created.Amount,
			"paymentStatus":   created.PaymentStatus,
			"remainingAmount": created.RemainingAmount,
		})
	}); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, actualPaymentDTO(created))
}
