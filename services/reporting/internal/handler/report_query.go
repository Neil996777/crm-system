package handler

import "net/http"

func (h *ReportingHandler) salesOverview(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	if actor.Role == "Sales" || actor.Role == "" {
		if err := h.appendReportAccessDenied(r, "sales-overview", actor); err != nil {
			writeError(w, http.StatusServiceUnavailable, "AUDIT_LOG_FAILED", "system", "Audit log failed.")
			return
		}
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	scope := "team"
	teamID := actor.TeamID
	if teamID == "" {
		teamID = "single-team"
	}
	allTeams := false
	if actor.Role == "Administrator" {
		scope = "all"
		allTeams = true
	}
	data, err := h.repo.SalesOverview(r.Context(), teamID, allTeams)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"scope": scope,
		"filters": map[string]any{
			"teamId":   teamID,
			"archived": "active_default",
			"from":     r.URL.Query().Get("from"),
			"to":       r.URL.Query().Get("to"),
			"groupBy":  r.URL.Query().Get("groupBy"),
		},
		"currency":   "CNY",
		"metrics":    data.Metrics.Map(),
		"breakdowns": data.Breakdowns,
		"groups":     data.Groups,
		"emptyState": data.Metrics.Empty(),
	})
}
