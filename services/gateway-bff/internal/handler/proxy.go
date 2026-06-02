package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"crm-system/services/gateway-bff/internal/authz"
	"crm-system/services/gateway-bff/internal/middleware"
)

func (g *Gateway) proxy(w http.ResponseWriter, r *http.Request) {
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
	baseURL, ok := g.config.Routes[resource]
	if !ok {
		writeEnvelopeError(w, http.StatusNotFound, correlationID, middleware.ErrorEnvelope{
			Code:        "NOT_FOUND",
			Category:    "not_found",
			SafeMessage: "The requested resource was not found.",
		})
		return
	}
	targetURL, err := buildTargetURL(baseURL, r)
	if err != nil {
		writeEnvelopeError(w, http.StatusBadRequest, correlationID, middleware.ErrorEnvelope{
			Code:        "INVALID_REQUEST",
			Category:    "validation",
			SafeMessage: "The request is invalid.",
		})
		return
	}
	resp, err := g.forward(r, targetURL, actor, correlationID)
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

func (g *Gateway) authProxy(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.CorrelationID(r)
	targetURL := strings.TrimRight(g.config.IdentityBaseURL, "/") + r.URL.Path
	req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, r.Body)
	if err != nil {
		writeEnvelopeError(w, http.StatusBadRequest, correlationID, middleware.ErrorEnvelope{Code: "INVALID_REQUEST", Category: "validation", SafeMessage: "The request is invalid."})
		return
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	req.Header.Set("X-Correlation-Id", correlationID)
	for _, cookie := range r.Cookies() {
		req.AddCookie(cookie)
	}
	client := g.config.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		writeEnvelopeError(w, http.StatusBadGateway, correlationID, middleware.ErrorEnvelope{Code: "DEPENDENCY_UNAVAILABLE", Category: "dependency", SafeMessage: "A required service is unavailable."})
		return
	}
	defer resp.Body.Close()
	for _, cookie := range resp.Cookies() {
		http.SetCookie(w, cookie)
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("copy auth proxy response: %v", err)
	}
}

func (g *Gateway) forward(source *http.Request, targetURL string, actor authz.ActorContext, correlationID string) (*http.Response, error) {
	client := g.config.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(source.Context(), source.Method, targetURL, source.Body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", source.Header.Get("Content-Type"))
	req.Header.Set("X-Correlation-Id", correlationID)
	req.Header.Set("X-Actor-User-Id", actor.ID)
	req.Header.Set("X-Actor-Role", actor.Role)
	req.Header.Set("X-Actor-Status", actor.Status)
	req.Header.Set("X-Actor-Display", actor.DisplayName)
	return client.Do(req)
}

func buildTargetURL(baseURL string, r *http.Request) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	targetPath := strings.TrimPrefix(r.URL.Path, "/api")
	parsed.Path = strings.TrimRight(parsed.Path, "/") + targetPath
	parsed.RawQuery = r.URL.RawQuery
	return parsed.String(), nil
}
