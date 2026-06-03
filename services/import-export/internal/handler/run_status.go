package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type actorContext struct {
	ID     string
	Role   string
	TeamID string
}

func (a actorContext) TeamIDOrDefault() string {
	if a.TeamID == "" {
		return "single-team"
	}
	return a.TeamID
}

func actorFromRequest(r *http.Request) actorContext {
	return actorContext{
		ID:     r.Header.Get("X-Actor-User-Id"),
		Role:   r.Header.Get("X-Actor-Role"),
		TeamID: r.Header.Get("X-Actor-Team-Id"),
	}
}

func writeError(w http.ResponseWriter, status int, code, category, safeMessage string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":        code,
			"category":    category,
			"safeMessage": safeMessage,
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
