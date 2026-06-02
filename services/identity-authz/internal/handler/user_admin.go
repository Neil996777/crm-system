package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"crm-system/services/identity-authz/internal/authz"
	"crm-system/services/identity-authz/internal/domain"
	"crm-system/services/identity-authz/internal/event"
)

func (h *AuthHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := h.requireAdministrator(w, r); !ok {
		return
	}
	users, err := h.users.List(r.Context())
	if err != nil {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	items := make([]map[string]any, 0, len(users))
	for _, user := range users {
		items = append(items, userDTO(user))
	}
	activeAdmins, err := h.users.CountActiveAdministrators(r.Context())
	if err != nil {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "activeAdministratorCount": activeAdmins})
}

func (h *AuthHandler) createUser(w http.ResponseWriter, r *http.Request) {
	actor, _, ok := h.requireAdministrator(w, r)
	if !ok {
		return
	}
	var request struct {
		Email       string `json:"email"`
		DisplayName string `json:"displayName"`
		Password    string `json:"password"`
		Role        string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	role := domain.Role(request.Role)
	if !role.Valid() || request.Email == "" || request.Password == "" {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	user, err := h.users.Create(r.Context(), request.Email, request.DisplayName, request.Password, role)
	if err != nil {
		log.Printf("create user: %v", err)
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	h.appendUserChange(r, actor.ID, user.ID, "create_user", "", string(user.Role), "success")
	writeJSON(w, http.StatusCreated, map[string]any{"user": userDTO(user)})
}

func (h *AuthHandler) changeUserRole(w http.ResponseWriter, r *http.Request) {
	actor, _, ok := h.requireAdministrator(w, r)
	if !ok {
		return
	}
	targetID := r.PathValue("id")
	var request struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	newRole := domain.Role(request.Role)
	if !newRole.Valid() {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	target, err := h.users.FindByID(r.Context(), targetID)
	if err != nil {
		writeErrorCode(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if h.blockedByLastAdminGuard(w, r, actor.ID, target, newRole, target.Status) {
		return
	}
	updated, err := h.users.UpdateRole(r.Context(), target.ID, newRole)
	if err != nil {
		log.Printf("change role: %v", err)
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	h.appendUserChange(r, actor.ID, updated.ID, "change_role", string(target.Role), string(updated.Role), "success")
	writeJSON(w, http.StatusOK, map[string]any{"user": userDTO(updated)})
}

func (h *AuthHandler) changeUserStatus(w http.ResponseWriter, r *http.Request) {
	actor, _, ok := h.requireAdministrator(w, r)
	if !ok {
		return
	}
	targetID := r.PathValue("id")
	var request struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	newStatus := domain.UserStatus(request.Status)
	if newStatus != domain.UserStatusActive && newStatus != domain.UserStatusDisabled {
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	target, err := h.users.FindByID(r.Context(), targetID)
	if err != nil {
		writeErrorCode(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if h.blockedByLastAdminGuard(w, r, actor.ID, target, target.Role, newStatus) {
		return
	}
	updated, err := h.users.UpdateStatus(r.Context(), target.ID, newStatus)
	if err != nil {
		log.Printf("change status: %v", err)
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if newStatus == domain.UserStatusDisabled {
		if err := h.sessions.RevokeForUser(r.Context(), target.ID, time.Now().UTC()); err != nil {
			log.Printf("revoke disabled user sessions: %v", err)
		}
	}
	h.appendUserChange(r, actor.ID, updated.ID, "change_status", string(target.Status), string(updated.Status), "success")
	writeJSON(w, http.StatusOK, map[string]any{"user": userDTO(updated)})
}

func (h *AuthHandler) requireAdministrator(w http.ResponseWriter, r *http.Request) (domain.User, string, bool) {
	actor, sessionID, ok := h.authenticate(r.Context(), r)
	if !ok {
		writeError(w, http.StatusUnauthorized, safeAuthMessage)
		return domain.User{}, "", false
	}
	if actor.Role != domain.RoleAdministrator {
		h.appendAccessDenied(r.Context(), actor.ID, "user_admin_denied")
		writeErrorCode(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return domain.User{}, "", false
	}
	return actor, sessionID, true
}

func (h *AuthHandler) blockedByLastAdminGuard(w http.ResponseWriter, r *http.Request, actorID string, target domain.User, newRole domain.Role, newStatus domain.UserStatus) bool {
	count, err := h.users.CountActiveAdministrators(r.Context())
	if err != nil {
		log.Printf("count active admins: %v", err)
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return true
	}
	if !domain.WouldRemoveLastActiveAdministrator(target, newRole, newStatus, count) {
		return false
	}
	if err := h.outbox.Append(r.Context(), event.LastAdministratorBlocked, target.ID, map[string]any{
		"actorId":  actorID,
		"targetId": target.ID,
		"result":   "blocked",
		"reason":   "last_active_administrator",
	}); err != nil {
		log.Printf("append last admin blocked event: %v", err)
	}
	writeErrorCode(w, http.StatusConflict, "LAST_ADMIN_BLOCKED", "conflict", "The last active Administrator cannot be disabled or downgraded.")
	return true
}

func (h *AuthHandler) appendUserChange(r *http.Request, actorID, targetID, action, before, after, result string) {
	if err := h.outbox.Append(r.Context(), event.UserRoleStatusChanged, targetID, map[string]any{
		"actorId":  actorID,
		"targetId": targetID,
		"action":   action,
		"before":   before,
		"after":    after,
		"result":   result,
	}); err != nil {
		log.Printf("append user change event: %v", err)
	}
	actorDisplay := actorID
	if actor, err := h.users.FindByID(r.Context(), actorID); err == nil {
		actorDisplay = actor.DisplayName
	}
	// TASK-028 / ACC-022 / PSM-009 / CONTRACT-013: administrator user-management
	// actions are persisted as global operation-log events in audit-history.
	if err := h.audit.AppendOperationLog(r.Context(), authz.OperationLogInput{
		ActorID:       actorID,
		ActorRole:     "Administrator",
		ActorDisplay:  actorDisplay,
		Action:        action,
		ResourceType:  "User",
		ResourceID:    targetID,
		Result:        result,
		BeforeSummary: map[string]any{"value": before},
		AfterSummary:  map[string]any{"value": after},
		SafeSummary:   "Administrator " + action + " on user",
		CorrelationID: r.Header.Get("X-Correlation-Id"),
	}); err != nil {
		log.Printf("append operation log: %v", err)
	}
}
