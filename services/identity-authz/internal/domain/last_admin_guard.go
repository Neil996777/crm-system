package domain

func WouldRemoveLastActiveAdministrator(target User, newRole Role, newStatus UserStatus, activeAdministratorCount int) bool {
	if target.Role != RoleAdministrator || target.Status != UserStatusActive {
		return false
	}
	if activeAdministratorCount > 1 {
		return false
	}
	return newRole != RoleAdministrator || newStatus != UserStatusActive
}
