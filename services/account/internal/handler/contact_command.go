package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"crm-system/services/account/internal/domain"
	"crm-system/services/account/internal/event"
	"crm-system/services/account/internal/repo"
)

func (h *AccountHandler) createContactForAccount(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	account, err := h.repo.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if !domain.CanReadAccount(actor.ID, actor.Role, account) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	var request struct {
		ContactName         string `json:"contactName"`
		Email               string `json:"email"`
		Phone               string `json:"phone"`
		RoleNote            string `json:"roleNote"`
		ProceedWarningToken string `json:"proceedWarningToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contact input is invalid.")
		return
	}
	contact, err := domain.NewContact(domain.Contact{
		AccountID:   account.ID,
		ContactName: request.ContactName,
		Email:       request.Email,
		Phone:       request.Phone,
		RoleNote:    request.RoleNote,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contact input is invalid.")
		return
	}
	if request.ProceedWarningToken != "" {
		candidate := domain.DuplicateCandidate{TargetType: "contact", Email: contact.Email, Phone: contact.Phone}
		signature, _ := domain.DuplicateSignature(candidate)
		if err := h.duplicates.ConsumeToken(r.Context(), request.ProceedWarningToken, candidate.TargetType, actor.ID, signature); err != nil {
			writeDuplicateTokenError(w, err)
			return
		}
	}
	var created domain.Contact
	if err := h.inTransaction(r.Context(), func(_ *repo.AccountRepo, txContacts *repo.ContactRepo, txOutbox *event.Outbox) error {
		var err error
		created, err = txContacts.Create(r.Context(), contact)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.ContactCreated, created.ID, map[string]any{
			"traceability": "TASK-011 ACC-006 PIM-006 PIM-BEH-008 PSM-003 CONTRACT-005 CONTRACT-006",
			"actorId":      actor.ID,
			"actorRole":    actor.Role,
			"actorDisplay": actor.ID,
			"accountId":    created.AccountID,
			"contactId":    created.ID,
		})
	}); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contact input is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, contactDTO(created))
}

func (h *AccountHandler) createContactWithoutAccount(w http.ResponseWriter, _ *http.Request) {
	writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The contact input is invalid.")
}
