package domain

type UserStatus string

const (
	UserStatusActive   UserStatus = "Active"
	UserStatusDisabled UserStatus = "Disabled"
)

type User struct {
	ID           string
	Email        string
	DisplayName  string
	PasswordHash string
	Role         Role
	Status       UserStatus
	AuthzVersion int
}

func (u User) Active() bool {
	return u.Status == UserStatusActive
}
