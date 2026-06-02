package handler

import (
	"errors"
	"net/http"

	"crm-system/services/account/internal/domain"
	"crm-system/services/account/internal/repo"
)

func (h *AccountHandler) listAccounts(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	accounts, err := h.repo.List(r.Context(), actor.ID, actor.Role, r.URL.Query().Get("search"), r.URL.Query().Get("customerStatus"), r.URL.Query().Get("includeArchived") == "true")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_FILTER", "validation", "The filter is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(accounts))
	for _, account := range accounts {
		items = append(items, accountDTO(account))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *AccountHandler) getAccount(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, accountDTO(account))
}
