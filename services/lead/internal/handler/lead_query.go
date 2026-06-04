package handler

import (
	"errors"
	"net/http"

	"crm-system/services/lead/internal/domain"
	"crm-system/services/lead/internal/repo"
)

func (h *LeadHandler) listLeads(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	includeArchived := r.URL.Query().Get("includeArchived") == "true"
	leads, err := h.repo.List(r.Context(), actor.ID, actor.Role, r.URL.Query().Get("search"), includeArchived)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_FILTER", "validation", "The filter is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(leads))
	for _, lead := range leads {
		items = append(items, leadDTO(lead))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *LeadHandler) getLead(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	lead, err := h.repo.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if !domain.CanReadLead(actor.ID, actor.Role, lead) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	writeJSON(w, http.StatusOK, leadDTO(lead))
}
