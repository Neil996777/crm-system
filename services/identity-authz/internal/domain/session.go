package domain

import "time"

type Session struct {
	ID                  string
	UserID              string
	AuthzVersionAtIssue int
	ExpiresAt           time.Time
	IdleExpiresAt       time.Time
	RevokedAt           *time.Time
	CreatedAt           time.Time
	LastSeenAt          time.Time
}

func (s Session) Active(now time.Time) bool {
	if s.RevokedAt != nil {
		return false
	}
	if !s.ExpiresAt.After(now) {
		return false
	}
	return s.IdleExpiresAt.After(now)
}
