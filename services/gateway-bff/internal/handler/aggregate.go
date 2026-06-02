package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"crm-system/services/gateway-bff/internal/middleware"
)

func writeEnvelopeError(w http.ResponseWriter, status int, correlationID string, err middleware.ErrorEnvelope) {
	writeJSON(w, status, map[string]any{
		"correlationId": correlationID,
		"error": map[string]any{
			"code":        err.Code,
			"category":    err.Category,
			"safeMessage": err.SafeMessage,
			"fieldErrors": err.FieldErrors,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write response: %v", err)
	}
}
