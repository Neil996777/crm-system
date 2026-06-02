package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"crm-system/services/gateway-bff/internal/authz"
	"crm-system/services/gateway-bff/internal/middleware"
)

type recordContext struct {
	OwnerID    string `json:"ownerId"`
	AssigneeID string `json:"assigneeId"`
	TeamID     string `json:"teamId"`
}

// ACC-014 / PSM-009 / CONTRACT-013: BFF exposes record-local history only after
// the owning Query API proves the actor may read the related record.
func (g *Gateway) recordHistory(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.CorrelationID(r)
	actor, err := g.authz.CurrentUser(r.Context(), r, correlationID)
	if err != nil {
		writeEnvelopeError(w, http.StatusUnauthorized, correlationID, middleware.ErrorEnvelope{
			Code:        "AUTHENTICATION_FAILED",
			Category:    "authentication",
			SafeMessage: "Authentication failed.",
		})
		return
	}

	resource := r.PathValue("resource")
	recordID := r.PathValue("id")
	baseURL, ok := g.config.Routes[resource]
	if !ok {
		writeEnvelopeError(w, http.StatusNotFound, correlationID, middleware.ErrorEnvelope{
			Code:        "NOT_FOUND",
			Category:    "not_found",
			SafeMessage: "The requested resource was not found.",
		})
		return
	}

	record, ok := g.authorizedRecordContext(w, r, actor, correlationID, baseURL, resource, recordID)
	if !ok {
		return
	}

	auditBaseURL, ok := g.config.Routes["history"]
	if !ok || strings.TrimSpace(auditBaseURL) == "" {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{
			Code:        "DEPENDENCY_UNAVAILABLE",
			Category:    "dependency",
			SafeMessage: "A required service is unavailable.",
		})
		return
	}
	historyURL, err := buildHistoryURL(auditBaseURL, resourceHistoryType(resource), recordID, actor, record)
	if err != nil {
		writeEnvelopeError(w, http.StatusBadRequest, correlationID, middleware.ErrorEnvelope{
			Code:        "INVALID_REQUEST",
			Category:    "validation",
			SafeMessage: "The request is invalid.",
		})
		return
	}
	resp, err := g.forward(r, historyURL, actor, correlationID)
	if err != nil {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{
			Code:        "DEPENDENCY_UNAVAILABLE",
			Category:    "dependency",
			SafeMessage: "A required service is unavailable.",
		})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{
			Code:        "DEPENDENCY_UNAVAILABLE",
			Category:    "dependency",
			SafeMessage: "A required service is unavailable.",
		})
		return
	}
	if resp.StatusCode >= 400 {
		writeEnvelopeError(w, resp.StatusCode, correlationID, middleware.NormalizeError(body, "DEPENDENCY_ERROR", "dependency", "A required service returned an error."))
		return
	}
	var data any
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&data); err != nil {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{
			Code:        "DEPENDENCY_ERROR",
			Category:    "dependency",
			SafeMessage: "A required service returned an invalid response.",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"correlationId": correlationID, "data": data})
}

func (g *Gateway) authorizedRecordContext(w http.ResponseWriter, r *http.Request, actor authz.ActorContext, correlationID, baseURL, resource, recordID string) (recordContext, bool) {
	recordURL := strings.TrimRight(baseURL, "/") + "/" + resource + "/" + url.PathEscape(recordID)
	resp, err := g.forward(r, recordURL, actor, correlationID)
	if err != nil {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{Code: "DEPENDENCY_UNAVAILABLE", Category: "dependency", SafeMessage: "A required service is unavailable."})
		return recordContext{}, false
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{Code: "DEPENDENCY_UNAVAILABLE", Category: "dependency", SafeMessage: "A required service is unavailable."})
		return recordContext{}, false
	}
	if resp.StatusCode >= 400 {
		writeEnvelopeError(w, resp.StatusCode, correlationID, middleware.NormalizeError(body, "DEPENDENCY_ERROR", "dependency", "A required service returned an error."))
		return recordContext{}, false
	}
	var record recordContext
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&record); err != nil {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{Code: "DEPENDENCY_ERROR", Category: "dependency", SafeMessage: "A required service returned an invalid response."})
		return recordContext{}, false
	}
	return record, true
}

func buildHistoryURL(baseURL, resourceType, recordID string, actor authz.ActorContext, record recordContext) (string, error) {
	parsed, err := url.Parse(strings.TrimRight(baseURL, "/") + fmt.Sprintf("/records/%s/%s/history", url.PathEscape(resourceType), url.PathEscape(recordID)))
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	query.Set("actorId", actor.ID)
	query.Set("actorRole", actor.Role)
	query.Set("ownerId", record.OwnerID)
	query.Set("assigneeId", record.AssigneeID)
	query.Set("teamId", record.TeamID)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func resourceHistoryType(resource string) string {
	switch resource {
	case "leads":
		return "Lead"
	case "accounts":
		return "Account"
	case "opportunities":
		return "Opportunity"
	case "quotes":
		return "Quote"
	case "contracts":
		return "Contract"
	case "payments":
		return "Payment"
	case "tasks":
		return "Task"
	default:
		return strings.TrimSuffix(resource, "s")
	}
}
