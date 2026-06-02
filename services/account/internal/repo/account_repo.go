package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"

	"crm-system/services/account/internal/domain"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrVersionConflict = errors.New("version conflict")
)

type AccountRepo struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	account.ID = "acct_" + randomHex(16)
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO account.accounts (id, company_name, customer_status, owner_id, version)
		VALUES ($1, $2, $3, $4, 1)
		RETURNING updated_at
	`, account.ID, account.CompanyName, account.CustomerStatus, account.OwnerID).Scan(&account.UpdatedAt)
	if err != nil {
		return domain.Account{}, err
	}
	return account, nil
}

func (r *AccountRepo) Find(ctx context.Context, id string) (domain.Account, error) {
	var account domain.Account
	var archivedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, company_name, customer_status, owner_id, archived_at, archived_by, archive_reason, version, updated_at
		FROM account.accounts
		WHERE id = $1
	`, id).Scan(&account.ID, &account.CompanyName, &account.CustomerStatus, &account.OwnerID, &archivedAt, &account.ArchivedBy, &account.ArchiveReason, &account.Version, &account.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Account{}, ErrNotFound
	}
	if archivedAt.Valid {
		account.ArchivedAt = archivedAt.Time
	}
	return account, err
}

func (r *AccountRepo) List(ctx context.Context, actorID, actorRole, search, customerStatus string, includeArchived bool) ([]domain.Account, error) {
	search = strings.TrimSpace(search)
	customerStatus = strings.TrimSpace(customerStatus)
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, company_name, customer_status, owner_id, archived_at, archived_by, archive_reason, version, updated_at
		FROM account.accounts
		WHERE ($1 <> 'Sales' OR owner_id = $2)
		  AND ($3 = '' OR company_name ILIKE '%' || $3 || '%')
		  AND ($4 = '' OR customer_status = $4)
		  AND ($5 = true OR archived_at IS NULL)
		ORDER BY updated_at DESC, id ASC
	`, actorRole, actorID, search, customerStatus, includeArchived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		var archivedAt sql.NullTime
		if err := rows.Scan(&account.ID, &account.CompanyName, &account.CustomerStatus, &account.OwnerID, &archivedAt, &account.ArchivedBy, &account.ArchiveReason, &account.Version, &account.UpdatedAt); err != nil {
			return nil, err
		}
		if archivedAt.Valid {
			account.ArchivedAt = archivedAt.Time
		}
		accounts = append(accounts, account)
	}
	return accounts, rows.Err()
}

func (r *AccountRepo) Update(ctx context.Context, id string, expectedVersion int, updated domain.Account) (domain.Account, error) {
	err := r.db.QueryRowContext(ctx, `
		UPDATE account.accounts
		SET company_name = $2, customer_status = $3, owner_id = $4, version = version + 1, updated_at = now()
		WHERE id = $1 AND version = $5
		RETURNING version, updated_at
	`, id, updated.CompanyName, updated.CustomerStatus, updated.OwnerID, expectedVersion).Scan(&updated.Version, &updated.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Account{}, ErrVersionConflict
	}
	if err != nil {
		return domain.Account{}, err
	}
	updated.ID = id
	return updated, nil
}

func (r *AccountRepo) Archive(ctx context.Context, id string, expectedVersion int, actorID, reason string) (domain.Account, error) {
	var account domain.Account
	var archivedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		UPDATE account.accounts
		SET archived_at = now(),
		    archived_by = $2,
		    archive_reason = $3,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $4 AND archived_at IS NULL
		RETURNING id, company_name, customer_status, owner_id, archived_at, archived_by, archive_reason, version, updated_at
	`, id, actorID, strings.TrimSpace(reason), expectedVersion).Scan(&account.ID, &account.CompanyName, &account.CustomerStatus, &account.OwnerID, &archivedAt, &account.ArchivedBy, &account.ArchiveReason, &account.Version, &account.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Account{}, ErrVersionConflict
	}
	if archivedAt.Valid {
		account.ArchivedAt = archivedAt.Time
	}
	return account, err
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
