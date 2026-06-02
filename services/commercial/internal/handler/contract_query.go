package handler

import (
	"errors"
	"net/http"

	"crm-system/services/commercial/internal/repo"
)

func (h *CommercialHandler) listContracts(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	includeArchived := r.URL.Query().Get("includeArchived") == "true"
	contracts, err := h.contracts.List(r.Context(), actor.ID, actor.Role, r.URL.Query().Get("search"), includeArchived)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(contracts))
	for _, contract := range contracts {
		items = append(items, contractDTO(contract))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *CommercialHandler) getContract(w http.ResponseWriter, r *http.Request) {
	contract, err := h.contracts.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, contractDTO(contract))
}
