package handler

import (
	"context"
	"strings"

	"crm-system/services/identity-authz/internal/domain"
	"crm-system/services/identity-authz/internal/repo"
)

const (
	auditActorAnonymousID      = "anonymous"
	auditActorAnonymousRole    = "Unauthenticated"
	auditActorAnonymousDisplay = "Unauthenticated actor"
	auditActorUnknownID        = "unknown"
	auditActorUnknownRole      = "UnknownActor"
	auditActorUnknownDisplay   = "Unknown actor"
	auditAuthResourceID        = "auth"
	auditAdminUsersResourceID  = "admin-users"
	auditPermissionResourceID  = "permission-check"
)

type auditActor struct {
	ID      string
	Role    string
	Display string
}

type auditAccessDeniedOptions struct {
	ReasonCode    string
	Action        string
	ResourceType  string
	ResourceID    string
	Result        string
	ScopeSummary  string
	CorrelationID string
}

func auditActorFromUser(user domain.User) auditActor {
	actor := auditActor{
		ID:      strings.TrimSpace(user.ID),
		Role:    strings.TrimSpace(string(user.Role)),
		Display: strings.TrimSpace(user.DisplayName),
	}
	if actor.ID == "" {
		actor.ID = auditActorUnknownID
	}
	if actor.Role == "" {
		actor.Role = auditActorUnknownRole
	}
	if actor.Display == "" {
		actor.Display = strings.TrimSpace(user.Email)
	}
	if actor.Display == "" {
		actor.Display = actor.ID
	}
	return actor
}

func auditActorAnonymous() auditActor {
	return auditActor{
		ID:      auditActorAnonymousID,
		Role:    auditActorAnonymousRole,
		Display: auditActorAnonymousDisplay,
	}
}

func auditActorUnknown() auditActor {
	return auditActor{
		ID:      auditActorUnknownID,
		Role:    auditActorUnknownRole,
		Display: auditActorUnknownDisplay,
	}
}

func auditActorFromUserID(ctx context.Context, users *repo.UserRepo, userID string) auditActor {
	if strings.TrimSpace(userID) == "" {
		return auditActorAnonymous()
	}
	user, err := users.FindByID(ctx, userID)
	if err != nil {
		return auditActorUnknown()
	}
	return auditActorFromUser(user)
}

func accessDeniedPayload(actor auditActor, options auditAccessDeniedOptions) map[string]any {
	actor = completeAuditActor(actor)
	reasonCode := strings.TrimSpace(options.ReasonCode)
	if reasonCode == "" {
		reasonCode = "access_denied"
	}
	action := firstAuditNonEmpty(options.Action, accessDeniedAction(reasonCode))
	result := firstAuditNonEmpty(options.Result, accessDeniedResult(reasonCode))
	resourceType := firstAuditNonEmpty(options.ResourceType, accessDeniedResourceType(reasonCode))
	resourceID := firstAuditNonEmpty(options.ResourceID, accessDeniedResourceID(reasonCode))
	scopeSummary := firstAuditNonEmpty(options.ScopeSummary, accessDeniedScopeSummary(reasonCode))
	return map[string]any{
		"actorId":            actor.ID,
		"actorRole":          actor.Role,
		"actorDisplay":       actor.Display,
		"action":             action,
		"resourceType":       resourceType,
		"resourceId":         resourceID,
		"result":             result,
		"reason":             reasonCode,
		"reasonCode":         reasonCode,
		"beforeSummary":      map[string]any{},
		"afterSummary":       map[string]any{},
		"diffClassification": "Security Critical",
		"scopeSummary":       scopeSummary,
		"safeSummary":        "Identity authorization denied: " + reasonCode,
		"correlationId":      strings.TrimSpace(options.CorrelationID),
	}
}

func signOutPayload(actor auditActor) map[string]any {
	actor = completeAuditActor(actor)
	return map[string]any{
		"actorId":      actor.ID,
		"actorRole":    actor.Role,
		"actorDisplay": actor.Display,
		"resourceType": "Auth",
		"resourceId":   auditAuthResourceID,
		"result":       "success",
		"action":       "sign_out",
		"safeSummary":  "Identity authorization sign_out",
	}
}

func completeAuditActor(actor auditActor) auditActor {
	if strings.TrimSpace(actor.ID) == "" {
		actor.ID = auditActorUnknownID
	}
	if strings.TrimSpace(actor.Role) == "" {
		actor.Role = auditActorUnknownRole
	}
	if strings.TrimSpace(actor.Display) == "" {
		actor.Display = actor.ID
	}
	return actor
}

func accessDeniedAction(reasonCode string) string {
	if reasonCode == "login_failed" {
		return "login_failed"
	}
	return "access_denied"
}

func accessDeniedResult(reasonCode string) string {
	if reasonCode == "login_failed" {
		return "failed"
	}
	return "denied"
}

func accessDeniedResourceType(reasonCode string) string {
	switch reasonCode {
	case "login_failed", "unauthenticated", "invalid_session", "inactive_user", "authz_version_stale":
		return "Auth"
	case "user_admin_denied":
		return "User"
	default:
		return "Permission"
	}
}

func accessDeniedResourceID(reasonCode string) string {
	switch reasonCode {
	case "login_failed", "unauthenticated", "invalid_session", "inactive_user", "authz_version_stale":
		return auditAuthResourceID
	case "user_admin_denied":
		return auditAdminUsersResourceID
	default:
		return auditPermissionResourceID
	}
}

func accessDeniedScopeSummary(reasonCode string) string {
	switch reasonCode {
	case "login_failed", "unauthenticated", "invalid_session", "inactive_user", "authz_version_stale":
		return "authentication"
	case "user_admin_denied":
		return "administrator only"
	default:
		return "permission denied"
	}
}

func firstAuditNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
