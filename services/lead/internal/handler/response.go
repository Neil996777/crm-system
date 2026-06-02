package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"crm-system/services/lead/internal/domain"
	"crm-system/services/lead/internal/repo"
)

type actorContext struct {
	ID          string
	Role        string
	DisplayName string
}

func actorFromRequest(r *http.Request) actorContext {
	return actorContext{
		ID:          r.Header.Get("X-Actor-User-Id"),
		Role:        r.Header.Get("X-Actor-Role"),
		DisplayName: r.Header.Get("X-Actor-Display"),
	}
}

func leadDTO(lead domain.Lead) map[string]any {
	body := map[string]any{
		"id":                     lead.ID,
		"leadName":               lead.LeadName,
		"companyName":            lead.CompanyName,
		"email":                  lead.Email,
		"phone":                  lead.Phone,
		"source":                 lead.Source,
		"status":                 lead.Status,
		"ownerId":                lead.OwnerID,
		"needSummary":            lead.NeedSummary,
		"invalidReason":          lead.InvalidReason,
		"convertedAccountId":     lead.ConvertedAccountID,
		"convertedOpportunityId": lead.ConvertedOpportunityID,
		"version":                lead.Version,
		"updatedAt":              lead.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if !lead.ArchivedAt.IsZero() {
		body["archived"] = true
		body["archivedAt"] = lead.ArchivedAt.UTC().Format(time.RFC3339)
		body["archivedBy"] = lead.ArchivedBy
		body["archiveReason"] = lead.ArchiveReason
	} else {
		body["archived"] = false
	}
	return body
}

var errPermissionDenied = errors.New("permission denied")

func writeLeadLookupError(w http.ResponseWriter, err error) {
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
