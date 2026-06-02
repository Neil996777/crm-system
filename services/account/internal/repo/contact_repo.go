package repo

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"crm-system/services/account/internal/domain"
)

type ContactRepo struct {
	db *sql.DB
}

func NewContactRepo(db *sql.DB) *ContactRepo {
	return &ContactRepo{db: db}
}

func (r *ContactRepo) Create(ctx context.Context, contact domain.Contact) (domain.Contact, error) {
	contact.ID = "ctc_" + randomHex(16)
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO account.contacts (id, account_id, contact_name, email, phone, role_note, version)
		VALUES ($1, $2, $3, $4, $5, $6, 1)
		RETURNING updated_at
	`, contact.ID, contact.AccountID, contact.ContactName, contact.Email, contact.Phone, contact.RoleNote).Scan(&contact.UpdatedAt)
	if err != nil {
		return domain.Contact{}, err
	}
	return contact, nil
}

func (r *ContactRepo) ListByAccount(ctx context.Context, accountID string) ([]domain.Contact, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, account_id, contact_name, email, phone, role_note, version, updated_at
		FROM account.contacts
		WHERE account_id = $1
		ORDER BY updated_at DESC, id ASC
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contacts []domain.Contact
	for rows.Next() {
		var contact domain.Contact
		if err := rows.Scan(&contact.ID, &contact.AccountID, &contact.ContactName, &contact.Email, &contact.Phone, &contact.RoleNote, &contact.Version, &contact.UpdatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, rows.Err()
}

func (r *ContactRepo) List(ctx context.Context, actorID, actorRole, search string) ([]domain.Contact, error) {
	search = strings.TrimSpace(search)
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.account_id, a.company_name, c.contact_name, c.email, c.phone, c.role_note, c.version, c.updated_at
		FROM account.contacts c
		JOIN account.accounts a ON a.id = c.account_id
		WHERE ($1 <> 'Sales' OR a.owner_id = $2)
		  AND ($3 = '' OR c.contact_name ILIKE '%' || $3 || '%' OR c.email ILIKE '%' || $3 || '%' OR c.phone ILIKE '%' || $3 || '%' OR a.company_name ILIKE '%' || $3 || '%')
		ORDER BY c.updated_at DESC, c.id ASC
	`, actorRole, actorID, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contacts []domain.Contact
	for rows.Next() {
		var contact domain.Contact
		if err := rows.Scan(&contact.ID, &contact.AccountID, &contact.AccountName, &contact.ContactName, &contact.Email, &contact.Phone, &contact.RoleNote, &contact.Version, &contact.UpdatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, rows.Err()
}

func (r *ContactRepo) FindAuthorized(ctx context.Context, id, actorID, actorRole string) (domain.Contact, error) {
	var contact domain.Contact
	err := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.account_id, a.company_name, c.contact_name, c.email, c.phone, c.role_note, c.version, c.updated_at
		FROM account.contacts c
		JOIN account.accounts a ON a.id = c.account_id
		WHERE c.id = $1 AND ($2 <> 'Sales' OR a.owner_id = $3)
	`, id, actorRole, actorID).Scan(&contact.ID, &contact.AccountID, &contact.AccountName, &contact.ContactName, &contact.Email, &contact.Phone, &contact.RoleNote, &contact.Version, &contact.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Contact{}, ErrNotFound
	}
	return contact, err
}

func IsForeignKeyError(err error) bool {
	return err != nil && !errors.Is(err, sql.ErrNoRows)
}
