package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"crm-system/services/opportunity/internal/domain"
	"crm-system/services/opportunity/internal/repo"
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

func opportunityDTO(opportunity domain.Opportunity) map[string]any {
	body := map[string]any{
		"id":                opportunity.ID,
		"customerId":        opportunity.CustomerID,
		"ownerId":           opportunity.OwnerID,
		"stage":             opportunity.Stage,
		"expectedAmount":    opportunity.ExpectedAmount,
		"expectedCloseDate": domain.FormatCloseDate(opportunity.ExpectedCloseDate),
		"title":             opportunity.Title,
		"version":           opportunity.Version,
		"updatedAt":         opportunity.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if !opportunity.CloseDate.IsZero() {
		body["closeDate"] = domain.FormatCloseDate(opportunity.CloseDate)
	}
	if opportunity.WonContractID != "" {
		body["wonContractId"] = opportunity.WonContractID
	}
	if opportunity.LostReasonCode != "" {
		body["lostReasonCode"] = opportunity.LostReasonCode
	}
	if opportunity.LostReasonDetail != "" {
		body["lostReasonDetail"] = opportunity.LostReasonDetail
	}
	if !opportunity.ClosedAt.IsZero() {
		body["closedAt"] = opportunity.ClosedAt.UTC().Format(time.RFC3339)
	}
	if !opportunity.ArchivedAt.IsZero() {
		body["archived"] = true
		body["archivedAt"] = opportunity.ArchivedAt.UTC().Format(time.RFC3339)
		body["archivedBy"] = opportunity.ArchivedBy
		body["archiveReason"] = opportunity.ArchiveReason
	} else {
		body["archived"] = false
	}
	return body
}

var errPermissionDenied = errors.New("permission denied")

func writeOpportunityLookupError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repo.ErrNotFound):
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
	case errors.Is(err, errPermissionDenied):
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
	default:
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
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

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write response: %v", err)
	}
}
