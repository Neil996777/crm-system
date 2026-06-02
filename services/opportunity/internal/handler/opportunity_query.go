package handler

import (
	"errors"
	"net/http"

	"crm-system/services/opportunity/internal/domain"
	"crm-system/services/opportunity/internal/repo"
)

func (h *OpportunityHandler) listOpportunities(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	stage := r.URL.Query().Get("stage")
	if stage != "" && !domain.IsStage(stage) {
		writeError(w, http.StatusBadRequest, "INVALID_FILTER", "validation", "The filter is invalid.")
		return
	}
	includeArchived := r.URL.Query().Get("includeArchived") == "true"
	opportunities, err := h.repo.List(r.Context(), actor.ID, actor.Role, r.URL.Query().Get("search"), stage, includeArchived)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_FILTER", "validation", "The filter is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(opportunities))
	for _, opportunity := range opportunities {
		items = append(items, opportunityDTO(opportunity))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *OpportunityHandler) getOpportunity(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	opportunity, err := h.repo.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if !domain.CanReadOpportunity(actor.ID, actor.Role, opportunity) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	writeJSON(w, http.StatusOK, opportunityDTO(opportunity))
}
