package handler

import (
	"encoding/json"
	"net/http"

	"crm-system/services/identity-authz/internal/domain"
	"crm-system/services/identity-authz/internal/event"
)

func (h *AuthHandler) permissionCheck(w http.ResponseWriter, r *http.Request) {
	if !h.verifyServiceToken(r, "permission.check") {
		writeErrorCode(w, http.StatusUnauthorized, "SERVICE_AUTH_FAILED", "authentication", "Service authentication failed.")
		return
	}
	var request struct {
		ActorID  string `json:"actorId"`
		Action   string `json:"action"`
		Resource struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"resource"`
		Context struct {
			OwnerID    string `json:"ownerId"`
			AssigneeID string `json:"assigneeId"`
			TeamID     string `json:"teamId"`
		} `json:"context"`
		CorrelationID string `json:"correlationId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	actor, err := h.users.FindByID(r.Context(), request.ActorID)
	if err != nil {
		if err := h.appendAccessDenied(r.Context(), request.ActorID, "actor_not_found"); err != nil {
			writeDependencyUnavailable(w)
			return
		}
		writeJSON(w, http.StatusOK, permissionResponse(domain.PermissionDecision{DenialCategory: "permission_denied"}))
		return
	}
	decision := domain.DecidePermission(actor, domain.PermissionRequest{
		ActorID: request.ActorID,
		Action:  request.Action,
		Resource: domain.ResourceRef{
			Type: request.Resource.Type,
			ID:   request.Resource.ID,
		},
		Context: domain.PermissionContext{
			OwnerID:    request.Context.OwnerID,
			AssigneeID: request.Context.AssigneeID,
			TeamID:     request.Context.TeamID,
		},
	})
	if !decision.Allowed {
		tx, err := h.db.BeginTx(r.Context(), nil)
		if err != nil {
			writeDependencyUnavailable(w)
			return
		}
		txOutbox := event.NewOutboxTx(tx)
		if err := txOutbox.Append(r.Context(), event.UserAccessDenied, actor.ID, map[string]any{
			"actorId":       request.ActorID,
			"action":        request.Action,
			"resource":      map[string]any{"type": request.Resource.Type},
			"result":        "denied",
			"reason":        decision.DenialCategory,
			"correlationId": request.CorrelationID,
		}); err != nil {
			_ = tx.Rollback()
			writeDependencyUnavailable(w)
			return
		}
		if err := tx.Commit(); err != nil {
			writeDependencyUnavailable(w)
			return
		}
	}
	writeJSON(w, http.StatusOK, permissionResponse(decision))
}

func permissionResponse(decision domain.PermissionDecision) map[string]any {
	if decision.Allowed {
		return map[string]any{
			"allowed":        true,
			"scope":          decision.Scope,
			"denialCategory": nil,
		}
	}
	return map[string]any{
		"allowed":        false,
		"scope":          "",
		"denialCategory": decision.DenialCategory,
	}
}
