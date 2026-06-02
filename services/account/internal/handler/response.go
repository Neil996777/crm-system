package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"crm-system/services/account/internal/domain"
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

func accountDTO(account domain.Account) map[string]any {
	archivedAt := ""
	if !account.ArchivedAt.IsZero() {
		archivedAt = account.ArchivedAt.UTC().Format(time.RFC3339)
	}
	return map[string]any{
		"id":             account.ID,
		"companyName":    account.CompanyName,
		"customerStatus": account.CustomerStatus,
		"ownerId":        account.OwnerID,
		"archived":       !account.ArchivedAt.IsZero(),
		"archivedAt":     archivedAt,
		"archivedBy":     account.ArchivedBy,
		"archiveReason":  account.ArchiveReason,
		"version":        account.Version,
		"updatedAt":      account.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func contactDTO(contact domain.Contact) map[string]any {
	return map[string]any{
		"id":          contact.ID,
		"accountId":   contact.AccountID,
		"accountName": contact.AccountName,
		"contactName": contact.ContactName,
		"email":       contact.Email,
		"phone":       contact.Phone,
		"roleNote":    contact.RoleNote,
		"version":     contact.Version,
		"updatedAt":   contact.UpdatedAt.UTC().Format(time.RFC3339),
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
