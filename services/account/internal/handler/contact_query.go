package handler

import (
	"errors"
	"net/http"

	"crm-system/services/account/internal/domain"
	"crm-system/services/account/internal/repo"
)

func (h *AccountHandler) listContactsForAccount(w http.ResponseWriter, r *http.Request) {
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
	contacts, err := h.contacts.ListByAccount(r.Context(), account.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(contacts))
	for _, contact := range contacts {
		items = append(items, contactDTO(contact))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *AccountHandler) listContacts(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	contacts, err := h.contacts.List(r.Context(), actor.ID, actor.Role, r.URL.Query().Get("search"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_FILTER", "validation", "The filter is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(contacts))
	for _, contact := range contacts {
		items = append(items, contactDTO(contact))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *AccountHandler) getContact(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	contact, err := h.contacts.FindAuthorized(r.Context(), r.PathValue("id"), actor.ID, actor.Role)
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, contactDTO(contact))
}
