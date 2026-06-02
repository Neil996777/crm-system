package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"crm-system/services/gateway-bff/internal/middleware"
)

// ACC-022 / PSM-009 / CONTRACT-013: global operation logs are queried only
// through audit-history, with the authenticated actor role supplied by BFF.
func (g *Gateway) operationLog(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.CorrelationID(r)
	actor, err := g.authz.CurrentUser(r.Context(), r, correlationID)
	if err != nil {
		writeEnvelopeError(w, http.StatusUnauthorized, correlationID, middleware.ErrorEnvelope{Code: "AUTHENTICATION_FAILED", Category: "authentication", SafeMessage: "Authentication failed."})
		return
	}
	auditBaseURL, ok := g.config.Routes["history"]
	if !ok || strings.TrimSpace(auditBaseURL) == "" {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{Code: "DEPENDENCY_UNAVAILABLE", Category: "dependency", SafeMessage: "A required service is unavailable."})
		return
	}
	parsed, err := url.Parse(strings.TrimRight(auditBaseURL, "/") + "/operation-log")
	if err != nil {
		writeEnvelopeError(w, http.StatusBadRequest, correlationID, middleware.ErrorEnvelope{Code: "INVALID_REQUEST", Category: "validation", SafeMessage: "The request is invalid."})
		return
	}
	query := parsed.Query()
	query.Set("actorId", actor.ID)
	query.Set("actorRole", actor.Role)
	parsed.RawQuery = query.Encode()

	resp, err := g.forward(r, parsed.String(), actor, correlationID)
	if err != nil {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{Code: "DEPENDENCY_UNAVAILABLE", Category: "dependency", SafeMessage: "A required service is unavailable."})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{Code: "DEPENDENCY_UNAVAILABLE", Category: "dependency", SafeMessage: "A required service is unavailable."})
		return
	}
	if resp.StatusCode >= 400 {
		writeEnvelopeError(w, resp.StatusCode, correlationID, middleware.NormalizeError(body, "DEPENDENCY_ERROR", "dependency", "A required service returned an error."))
		return
	}
	var data any
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&data); err != nil {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{Code: "DEPENDENCY_ERROR", Category: "dependency", SafeMessage: "A required service returned an invalid response."})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"correlationId": correlationID, "data": data})
}
