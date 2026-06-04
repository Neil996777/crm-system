package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"crm-system/services/identity-authz/internal/domain"
	"crm-system/services/identity-authz/internal/event"
	"crm-system/services/identity-authz/internal/repo"
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
	ctx := r.Context()
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("begin create user tx: %v", err)
		writeDependencyUnavailable(w)
		return
	}
	users := repo.NewUserRepoTx(tx)
	outbox := event.NewOutboxTx(tx)
	user, err := users.Create(ctx, request.Email, request.DisplayName, request.Password, role)
	if err != nil {
		_ = tx.Rollback()
		log.Printf("create user: %v", err)
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if err := h.appendUserChangeOutbox(ctx, outbox, actor, user.ID, "create_user", "", string(user.Role), "success", r.Header.Get("X-Correlation-Id")); err != nil {
		_ = tx.Rollback()
		log.Printf("append create user event: %v", err)
		writeDependencyUnavailable(w)
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("commit create user tx: %v", err)
		writeDependencyUnavailable(w)
		return
	}
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
	ctx := r.Context()
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("begin change role tx: %v", err)
		writeDependencyUnavailable(w)
		return
	}
	users := repo.NewUserRepoTx(tx)
	outbox := event.NewOutboxTx(tx)
	target, err := users.FindByID(ctx, targetID)
	if err != nil {
		_ = tx.Rollback()
		writeErrorCode(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	blocked, guardErr := h.lastAdminGuardEvent(ctx, users, outbox, actor, target, newRole, target.Status, r.Header.Get("X-Correlation-Id"))
	if guardErr != nil {
		_ = tx.Rollback()
		log.Printf("append last admin blocked event: %v", guardErr)
		writeDependencyUnavailable(w)
		return
	}
	if blocked {
		if err := tx.Commit(); err != nil {
			log.Printf("commit last admin blocked tx: %v", err)
			writeDependencyUnavailable(w)
			return
		}
		writeErrorCode(w, http.StatusConflict, "LAST_ADMIN_BLOCKED", "conflict", "The last active Administrator cannot be disabled or downgraded.")
		return
	}
	updated, err := users.UpdateRole(ctx, target.ID, newRole)
	if err != nil {
		_ = tx.Rollback()
		log.Printf("change role: %v", err)
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if err := h.appendUserChangeOutbox(ctx, outbox, actor, updated.ID, "change_role", string(target.Role), string(updated.Role), "success", r.Header.Get("X-Correlation-Id")); err != nil {
		_ = tx.Rollback()
		log.Printf("append change role event: %v", err)
		writeDependencyUnavailable(w)
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("commit change role tx: %v", err)
		writeDependencyUnavailable(w)
		return
	}
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
	ctx := r.Context()
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("begin change status tx: %v", err)
		writeDependencyUnavailable(w)
		return
	}
	users := repo.NewUserRepoTx(tx)
	sessions := repo.NewSessionRepoTx(tx)
	outbox := event.NewOutboxTx(tx)
	target, err := users.FindByID(ctx, targetID)
	if err != nil {
		_ = tx.Rollback()
		writeErrorCode(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	blocked, guardErr := h.lastAdminGuardEvent(ctx, users, outbox, actor, target, target.Role, newStatus, r.Header.Get("X-Correlation-Id"))
	if guardErr != nil {
		_ = tx.Rollback()
		log.Printf("append last admin blocked event: %v", guardErr)
		writeDependencyUnavailable(w)
		return
	}
	if blocked {
		if err := tx.Commit(); err != nil {
			log.Printf("commit last admin blocked tx: %v", err)
			writeDependencyUnavailable(w)
			return
		}
		writeErrorCode(w, http.StatusConflict, "LAST_ADMIN_BLOCKED", "conflict", "The last active Administrator cannot be disabled or downgraded.")
		return
	}
	updated, err := users.UpdateStatus(ctx, target.ID, newStatus)
	if err != nil {
		_ = tx.Rollback()
		log.Printf("change status: %v", err)
		writeErrorCode(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if newStatus == domain.UserStatusDisabled {
		if err := sessions.RevokeForUser(ctx, target.ID, time.Now().UTC()); err != nil {
			_ = tx.Rollback()
			log.Printf("revoke disabled user sessions: %v", err)
			writeDependencyUnavailable(w)
			return
		}
	}
	if err := h.appendUserChangeOutbox(ctx, outbox, actor, updated.ID, "change_status", string(target.Status), string(updated.Status), "success", r.Header.Get("X-Correlation-Id")); err != nil {
		_ = tx.Rollback()
		log.Printf("append change status event: %v", err)
		writeDependencyUnavailable(w)
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("commit change status tx: %v", err)
		writeDependencyUnavailable(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": userDTO(updated)})
}

func (h *AuthHandler) requireAdministrator(w http.ResponseWriter, r *http.Request) (domain.User, string, bool) {
	actor, sessionID, errorCode, ok := h.authenticate(r.Context(), r)
	if !ok {
		if errorCode == "" {
			errorCode = "AUTHENTICATION_FAILED"
		}
		writeErrorCode(w, http.StatusUnauthorized, errorCode, "authentication", safeAuthMessage)
		return domain.User{}, "", false
	}
	if actor.Role != domain.RoleAdministrator {
		h.appendAccessDenied(r.Context(), actor.ID, "user_admin_denied")
		writeErrorCode(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return domain.User{}, "", false
	}
	return actor, sessionID, true
}

func (h *AuthHandler) lastAdminGuardEvent(ctx context.Context, users *repo.UserRepo, outbox *event.Outbox, actor domain.User, target domain.User, newRole domain.Role, newStatus domain.UserStatus, correlationID string) (bool, error) {
	count, err := users.CountActiveAdministrators(ctx)
	if err != nil {
		return false, err
	}
	if !domain.WouldRemoveLastActiveAdministrator(target, newRole, newStatus, count) {
		return false, nil
	}
	return true, outbox.Append(ctx, event.LastAdministratorBlocked, target.ID, map[string]any{
		"actorId":       actor.ID,
		"actorRole":     string(actor.Role),
		"actorDisplay":  actor.DisplayName,
		"targetId":      target.ID,
		"result":        "blocked",
		"reason":        "last_active_administrator",
		"correlationId": correlationID,
	})
}

func (h *AuthHandler) appendUserChangeOutbox(ctx context.Context, outbox *event.Outbox, actor domain.User, targetID, action, before, after, result, correlationID string) error {
	return outbox.Append(ctx, event.UserRoleStatusChanged, targetID, map[string]any{
		"actorId":       actor.ID,
		"actorRole":     string(actor.Role),
		"actorDisplay":  actor.DisplayName,
		"targetId":      targetID,
		"action":        action,
		"before":        before,
		"after":         after,
		"result":        result,
		"correlationId": correlationID,
	})
}

func writeDependencyUnavailable(w http.ResponseWriter) {
	writeErrorCode(w, http.StatusServiceUnavailable, "DEPENDENCY_UNAVAILABLE", "dependency", "A required dependency is unavailable.")
}
