package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"

	"crm-system/services/identity-authz/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("not found")

type UserRepo struct {
	db *sql.DB
	q  sqlRunner
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db, q: db}
}

func NewUserRepoTx(tx *sql.Tx) *UserRepo {
	return &UserRepo{q: tx}
}

type sqlRunner interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.find(ctx, `
		SELECT id, email, display_name, password_hash, role_name, status, authz_version
		FROM identity_authz.users
		WHERE lower(email) = lower($1)
	`, email)
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (domain.User, error) {
	return r.find(ctx, `
		SELECT id, email, display_name, password_hash, role_name, status, authz_version
		FROM identity_authz.users
		WHERE id = $1
	`, id)
}

func (r *UserRepo) List(ctx context.Context) ([]domain.User, error) {
	rows, err := r.q.QueryContext(ctx, `
		SELECT id, email, display_name, password_hash, role_name, status, authz_version
		FROM identity_authz.users
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []domain.User
	for rows.Next() {
		var user domain.User
		var role string
		var status string
		if err := rows.Scan(&user.ID, &user.Email, &user.DisplayName, &user.PasswordHash, &role, &status, &user.AuthzVersion); err != nil {
			return nil, err
		}
		user.Role = domain.Role(role)
		user.Status = domain.UserStatus(status)
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepo) Create(ctx context.Context, email, displayName, password string, role domain.Role) (domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		ID:           "usr_" + randomHex(16),
		Email:        strings.TrimSpace(email),
		DisplayName:  strings.TrimSpace(displayName),
		PasswordHash: string(hash),
		Role:         role,
		Status:       domain.UserStatusActive,
		AuthzVersion: 1,
	}
	var roleName string
	var status string
	err = r.q.QueryRowContext(ctx, `
		INSERT INTO identity_authz.users (id, email, display_name, password_hash, role_name, status, authz_version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, email, display_name, password_hash, role_name, status, authz_version
	`, user.ID, user.Email, user.DisplayName, user.PasswordHash, string(user.Role), string(user.Status), user.AuthzVersion).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&roleName,
		&status,
		&user.AuthzVersion,
	)
	user.Role = domain.Role(roleName)
	user.Status = domain.UserStatus(status)
	return user, err
}

func (r *UserRepo) UpdateRole(ctx context.Context, id string, role domain.Role) (domain.User, error) {
	_, err := r.q.ExecContext(ctx, `
		UPDATE identity_authz.users
		SET role_name = $2, authz_version = authz_version + 1, updated_at = now()
		WHERE id = $1
	`, id, string(role))
	if err != nil {
		return domain.User{}, err
	}
	return r.FindByID(ctx, id)
}

func (r *UserRepo) UpdateStatus(ctx context.Context, id string, status domain.UserStatus) (domain.User, error) {
	_, err := r.q.ExecContext(ctx, `
		UPDATE identity_authz.users
		SET status = $2, authz_version = authz_version + 1, updated_at = now()
		WHERE id = $1
	`, id, string(status))
	if err != nil {
		return domain.User{}, err
	}
	return r.FindByID(ctx, id)
}

func (r *UserRepo) CountActiveAdministrators(ctx context.Context) (int, error) {
	var count int
	err := r.q.QueryRowContext(ctx, `
		SELECT count(*)
		FROM identity_authz.users
		WHERE role_name = 'Administrator' AND status = 'Active'
	`).Scan(&count)
	return count, err
}

func (r *UserRepo) find(ctx context.Context, query string, arg any) (domain.User, error) {
	var user domain.User
	var role string
	var status string
	err := r.q.QueryRowContext(ctx, query, arg).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&role,
		&status,
		&user.AuthzVersion,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, ErrNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	user.Role = domain.Role(role)
	user.Status = domain.UserStatus(status)
	return user, nil
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
