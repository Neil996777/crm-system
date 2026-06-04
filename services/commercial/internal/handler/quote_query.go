package handler

import (
	"errors"
	"net/http"

	"crm-system/services/commercial/internal/repo"
)

func (h *CommercialHandler) listQuotes(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	quotes, err := h.quotes.List(r.Context(), actor.ID, actor.Role, r.URL.Query().Get("search"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(quotes))
	for _, quote := range quotes {
		items = append(items, quoteDTO(quote))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *CommercialHandler) getQuote(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	quote, err := h.quotes.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if !canReadCommercialRecord(actor, quote.OwnerID) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	writeJSON(w, http.StatusOK, quoteDTO(quote))
}
