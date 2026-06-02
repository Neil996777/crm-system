package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"crm-system/services/identity-authz/internal/domain"
)

type SessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Create(ctx context.Context, session domain.Session) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO identity_authz.sessions (
			id, user_id, authz_version_at_issue, expires_at, idle_expires_at, created_at, last_seen_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, session.ID, session.UserID, session.AuthzVersionAtIssue, session.ExpiresAt, session.IdleExpiresAt, session.CreatedAt, session.LastSeenAt)
	return err
}

func (r *SessionRepo) FindByID(ctx context.Context, id string) (domain.Session, error) {
	var session domain.Session
	var revokedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, authz_version_at_issue, expires_at, idle_expires_at, revoked_at, created_at, last_seen_at
		FROM identity_authz.sessions
		WHERE id = $1
	`, id).Scan(
		&session.ID,
		&session.UserID,
		&session.AuthzVersionAtIssue,
		&session.ExpiresAt,
		&session.IdleExpiresAt,
		&revokedAt,
		&session.CreatedAt,
		&session.LastSeenAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Session{}, ErrNotFound
	}
	if err != nil {
		return domain.Session{}, err
	}
	if revokedAt.Valid {
		session.RevokedAt = &revokedAt.Time
	}
	return session, nil
}

func (r *SessionRepo) Touch(ctx context.Context, id string, now, idleExpiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE identity_authz.sessions
		SET last_seen_at = $2, idle_expires_at = $3
		WHERE id = $1 AND revoked_at IS NULL
	`, id, now, idleExpiresAt)
	return err
}

func (r *SessionRepo) Revoke(ctx context.Context, id string, now time.Time) (string, error) {
	var userID string
	err := r.db.QueryRowContext(ctx, `
		UPDATE identity_authz.sessions
		SET revoked_at = $2
		WHERE id = $1 AND revoked_at IS NULL
		RETURNING user_id
	`, id, now).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	return userID, err
}

func (r *SessionRepo) RevokeForUser(ctx context.Context, userID string, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE identity_authz.sessions
		SET revoked_at = $2
		WHERE user_id = $1 AND revoked_at IS NULL
	`, userID, now)
	return err
}
