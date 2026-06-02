package domain

import "strings"

type PermissionRequest struct {
	ActorID  string
	Action   string
	Resource ResourceRef
	Context  PermissionContext
}

type ResourceRef struct {
	Type string
	ID   string
}

type PermissionContext struct {
	OwnerID    string
	AssigneeID string
	TeamID     string
}

type PermissionDecision struct {
	Allowed        bool
	Scope          string
	DenialCategory string
}

func DecidePermission(actor User, request PermissionRequest) PermissionDecision {
	if !actor.Active() {
		return PermissionDecision{DenialCategory: "inactive_actor"}
	}
	if isHardDelete(request.Action) {
		return PermissionDecision{DenialCategory: "hard_delete_forbidden"}
	}
	if isUserAdminAction(request.Action) {
		if actor.Role == RoleAdministrator {
			return PermissionDecision{Allowed: true, Scope: "all"}
		}
		return PermissionDecision{DenialCategory: "permission_denied"}
	}
	if request.Action == "operation_log.read" && actor.Role != RoleAdministrator {
		return PermissionDecision{DenialCategory: "permission_denied"}
	}
	switch actor.Role {
	case RoleAdministrator:
		return PermissionDecision{Allowed: true, Scope: "all"}
	case RoleSalesManager:
		if request.Context.TeamID == "" || request.Context.TeamID == "single-team" {
			return PermissionDecision{Allowed: true, Scope: "team"}
		}
		return PermissionDecision{DenialCategory: "scope_denied"}
	case RoleSales:
		if request.Context.OwnerID == actor.ID || request.Context.AssigneeID == actor.ID {
			return PermissionDecision{Allowed: true, Scope: "owned"}
		}
		return PermissionDecision{DenialCategory: "scope_denied"}
	default:
		return PermissionDecision{DenialCategory: "permission_denied"}
	}
}

func isHardDelete(action string) bool {
	return strings.Contains(action, "hard_delete") || strings.HasSuffix(action, ".delete")
}

func isUserAdminAction(action string) bool {
	return strings.HasPrefix(action, "user.") || strings.HasPrefix(action, "role.")
}
