package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"crm-system/services/commercial/internal/domain"
	"crm-system/services/commercial/internal/repo"
)

type actorContext struct {
	ID   string
	Role string
}

func actorFromRequest(r *http.Request) actorContext {
	return actorContext{
		ID:   r.Header.Get("X-Actor-User-Id"),
		Role: r.Header.Get("X-Actor-Role"),
	}
}

func canReadCommercialRecord(actor actorContext, ownerID string) bool {
	return actor.Role != "Sales" || (actor.ID != "" && actor.ID == ownerID)
}

func quoteDTO(quote domain.Quote) map[string]any {
	return map[string]any{
		"id":            quote.ID,
		"opportunityId": quote.OpportunityID,
		"customerId":    quote.CustomerID,
		"amount":        quote.Amount,
		"status":        quote.Status,
		"validityEnd":   domain.FormatDate(quote.ValidityEnd),
		"ownerId":       quote.OwnerID,
		"version":       quote.Version,
		"updatedAt":     quote.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func contractDTO(contract domain.Contract) map[string]any {
	body := map[string]any{
		"id":                 contract.ID,
		"quoteId":            contract.QuoteID,
		"opportunityId":      contract.OpportunityID,
		"customerId":         contract.CustomerID,
		"amount":             contract.Amount,
		"status":             contract.Status,
		"contractNote":       contract.ContractNote,
		"expectedSignedDate": domain.FormatDate(contract.ExpectedSignedDate),
		"ownerId":            contract.OwnerID,
		"version":            contract.Version,
		"updatedAt":          contract.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if contract.AmountDifferenceReason != "" {
		body["amountDifferenceReason"] = contract.AmountDifferenceReason
	}
	if !contract.SignedEffectiveDate.IsZero() {
		body["signedEffectiveDate"] = domain.FormatDate(contract.SignedEffectiveDate)
	}
	if !contract.ArchivedAt.IsZero() {
		body["archived"] = true
		body["archivedAt"] = contract.ArchivedAt.UTC().Format(time.RFC3339)
		body["archivedBy"] = contract.ArchivedBy
		body["archiveReason"] = contract.ArchiveReason
	} else {
		body["archived"] = false
	}
	return body
}

func paymentPlanDTO(plan domain.PaymentPlan) map[string]any {
	body := map[string]any{
		"id":         plan.ID,
		"contractId": plan.ContractID,
		"dueAmount":  plan.DueAmount,
		"dueDate":    domain.FormatDate(plan.DueDate),
		"currency":   plan.Currency,
		"status":     plan.Status,
		"version":    plan.Version,
		"updatedAt":  plan.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if !plan.ArchivedAt.IsZero() {
		body["archived"] = true
		body["archivedAt"] = plan.ArchivedAt.UTC().Format(time.RFC3339)
		body["archivedBy"] = plan.ArchivedBy
		body["archiveReason"] = plan.ArchiveReason
	} else {
		body["archived"] = false
	}
	return body
}

func actualPaymentDTO(payment domain.ActualPayment) map[string]any {
	return map[string]any{
		"paymentId":       payment.ID,
		"contractId":      payment.ContractID,
		"amount":          payment.Amount,
		"paymentDate":     domain.FormatDate(payment.PaymentDate),
		"paymentStatus":   payment.PaymentStatus,
		"remainingAmount": payment.RemainingAmount,
		"version":         payment.Version,
		"updatedAt":       payment.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func writeError(w http.ResponseWriter, status int, code, category, safeMessage string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":        code,
			"category":    category,
			"safeMessage": safeMessage,
		},
	})
}

var errPermissionDenied = errors.New("permission denied")

func writeCommercialLookupError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repo.ErrNotFound):
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
	case errors.Is(err, errPermissionDenied):
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
	default:
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write response: %v", err)
	}
}
