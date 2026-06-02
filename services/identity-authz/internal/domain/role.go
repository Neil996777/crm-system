package domain

type Role string

const (
	RoleAdministrator Role = "Administrator"
	RoleSalesManager  Role = "Sales Manager"
	RoleSales         Role = "Sales"
)

func (r Role) Valid() bool {
	switch r {
	case RoleAdministrator, RoleSalesManager, RoleSales:
		return true
	default:
		return false
	}
}
