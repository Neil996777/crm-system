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

func (h *CommercialHandler) createContract(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		QuoteID                string `json:"quoteId"`
		OpportunityID          string `json:"opportunityId"`
		CustomerID             string `json:"customerId"`
		Amount                 string `json:"amount"`
		Status                 string `json:"status"`
		ContractNote           string `json:"contractNote"`
		ExpectedSignedDate     string `json:"expectedSignedDate"`
		AmountDifferenceReason string `json:"amountDifferenceReason"`
		OwnerID                string `json:"ownerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contract input is invalid.")
		return
	}
	if strings.TrimSpace(request.QuoteID) == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contract input is invalid.")
		return
	}
	expectedSignedDate, err := domain.ParseDate(request.ExpectedSignedDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contract input is invalid.")
		return
	}
	quote, err := h.quotes.Find(r.Context(), request.QuoteID)
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusBadRequest, "CONTRACT_QUOTE_INVALID", "business_rule", "The contract quote link is invalid.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if actor.Role == "Sales" && (request.OwnerID != actor.ID || quote.OwnerID != actor.ID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	contract, err := domain.NewContract(domain.Contract{
		QuoteID:                request.QuoteID,
		OpportunityID:          request.OpportunityID,
		CustomerID:             request.CustomerID,
		Amount:                 request.Amount,
		Status:                 request.Status,
		ContractNote:           request.ContractNote,
		ExpectedSignedDate:     expectedSignedDate,
		AmountDifferenceReason: request.AmountDifferenceReason,
		OwnerID:                request.OwnerID,
	}, quote)
	if errors.Is(err, domain.ErrContractQuoteInvalid) {
		writeError(w, http.StatusBadRequest, "CONTRACT_QUOTE_INVALID", "business_rule", "The contract quote link is invalid.")
		return
	}
	if errors.Is(err, domain.ErrAmountDifferenceReasonRequired) {
		writeError(w, http.StatusBadRequest, "AMOUNT_DIFFERENCE_REASON_REQUIRED", "business_rule", "A reason is required when contract amount differs from quote amount.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contract input is invalid.")
		return
	}
	var created domain.Contract
	err = h.inTransaction(r.Context(), func(_ *repo.QuoteRepo, txContracts *repo.ContractRepo, _ *repo.PaymentRepo, txOutbox *event.Outbox) error {
		var err error
		created, err = txContracts.Create(r.Context(), contract)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.ContractCreated, created.ID, map[string]any{
			"traceability":  "TASK-018 ACC-010 CIM-021 CIM-PROC-009 PIM-009 PIM-SM-005 PIM-INV-016 PIM-INV-018 PIM-BEH-014 PIM-BEH-016 PSM-006 CONTRACT-009 CONTRACT-010",
			"actorId":       actor.ID,
			"actorRole":     actor.Role,
			"actorDisplay":  actor.ID,
			"contractId":    created.ID,
			"quoteId":       created.QuoteID,
			"opportunityId": created.OpportunityID,
		})
	})
	if errors.Is(err, domain.ErrContractAlreadyExists) {
		writeError(w, http.StatusConflict, "CONTRACT_ALREADY_EXISTS", "conflict", "A contract already exists for this quote.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contract input is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, contractDTO(created))
}
