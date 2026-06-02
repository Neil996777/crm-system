package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"

	"crm-system/services/lead/internal/domain"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrVersionConflict  = errors.New("version conflict")
	ErrAlreadyConverted = errors.New("already converted")
)

type LeadRepo struct {
	db *sql.DB
}

func NewLeadRepo(db *sql.DB) *LeadRepo {
	return &LeadRepo{db: db}
}

func (r *LeadRepo) Create(ctx context.Context, lead domain.Lead) (domain.Lead, error) {
	lead.ID = "lead_" + randomHex(16)
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO lead.leads (id, lead_name, company_name, email, phone, source, status, owner_id, need_summary, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1)
		RETURNING updated_at
	`, lead.ID, lead.LeadName, lead.CompanyName, lead.Email, lead.Phone, lead.Source, lead.Status, lead.OwnerID, lead.NeedSummary).Scan(&lead.UpdatedAt)
	if err != nil {
		return domain.Lead{}, err
	}
	return lead, nil
}

func (r *LeadRepo) Find(ctx context.Context, id string) (domain.Lead, error) {
	var lead domain.Lead
	err := scanLead(r.db.QueryRowContext(ctx, `
		SELECT id, lead_name, company_name, email, phone, source, status, owner_id, need_summary,
		       invalid_reason, converted_account_id, converted_opportunity_id, conversion_idempotency_key,
		       archived_at, archived_by, archive_reason, version, updated_at
		FROM lead.leads
		WHERE id = $1
	`, id), &lead)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Lead{}, ErrNotFound
	}
	return lead, err
}

func (r *LeadRepo) List(ctx context.Context, actorID, actorRole, search string, includeArchived bool) ([]domain.Lead, error) {
	search = strings.TrimSpace(search)
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, lead_name, company_name, email, phone, source, status, owner_id, need_summary,
		       invalid_reason, converted_account_id, converted_opportunity_id, conversion_idempotency_key,
		       archived_at, archived_by, archive_reason, version, updated_at
		FROM lead.leads
		WHERE ($1 <> 'Sales' OR owner_id = $2)
		  AND ($3 = '' OR lead_name ILIKE '%' || $3 || '%' OR company_name ILIKE '%' || $3 || '%' OR email ILIKE '%' || $3 || '%' OR phone ILIKE '%' || $3 || '%')
		  AND ($4 = true OR archived_at IS NULL)
		ORDER BY updated_at DESC, id ASC
	`, actorRole, actorID, search, includeArchived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var leads []domain.Lead
	for rows.Next() {
		var lead domain.Lead
		if err := scanLead(rows, &lead); err != nil {
			return nil, err
		}
		leads = append(leads, lead)
	}
	return leads, rows.Err()
}

func (r *LeadRepo) Archive(ctx context.Context, id string, expectedVersion int, actorID, reason string) (domain.Lead, error) {
	var lead domain.Lead
	err := scanLead(r.db.QueryRowContext(ctx, `
		UPDATE lead.leads
		SET archived_at = now(),
		    archived_by = $2,
		    archive_reason = $3,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $4 AND archived_at IS NULL
		RETURNING id, lead_name, company_name, email, phone, source, status, owner_id, need_summary,
		          invalid_reason, converted_account_id, converted_opportunity_id, conversion_idempotency_key,
		          archived_at, archived_by, archive_reason, version, updated_at
	`, id, actorID, strings.TrimSpace(reason), expectedVersion), &lead)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Lead{}, ErrVersionConflict
	}
	return lead, err
}

func (r *LeadRepo) UpdateQualification(ctx context.Context, id string, expectedVersion int, updated domain.Lead) (domain.Lead, error) {
	err := r.db.QueryRowContext(ctx, `
		UPDATE lead.leads
		SET status = $2,
		    invalid_reason = $3,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $4
		RETURNING version, updated_at
	`, id, updated.Status, updated.InvalidReason, expectedVersion).Scan(&updated.Version, &updated.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Lead{}, ErrVersionConflict
	}
	return updated, err
}

func (r *LeadRepo) Convert(ctx context.Context, id string, expectedVersion int, updated domain.Lead) (domain.Lead, error) {
	err := r.db.QueryRowContext(ctx, `
		UPDATE lead.leads
		SET status = $2,
		    converted_account_id = $3,
		    converted_opportunity_id = $4,
		    conversion_idempotency_key = $5,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $6
		RETURNING version, updated_at
	`, id, updated.Status, updated.ConvertedAccountID, updated.ConvertedOpportunityID, updated.ConversionIDempotencyKey, expectedVersion).Scan(&updated.Version, &updated.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Lead{}, ErrVersionConflict
	}
	return updated, err
}

func (r *LeadRepo) TransferOwner(ctx context.Context, id string, expectedVersion int, newOwnerID string) (domain.Lead, domain.Lead, error) {
	current, err := r.Find(ctx, id)
	if err != nil {
		return domain.Lead{}, domain.Lead{}, err
	}
	if current.Version != expectedVersion {
		return domain.Lead{}, domain.Lead{}, ErrVersionConflict
	}
	updated := current
	updated.OwnerID = newOwnerID
	updated.Status = domain.StatusPendingQualification
	updated.Version = current.Version + 1
	err = r.db.QueryRowContext(ctx, `
		UPDATE lead.leads
		SET owner_id = $2, status = $3, version = version + 1, updated_at = now()
		WHERE id = $1 AND version = $4
		RETURNING updated_at
	`, id, updated.OwnerID, updated.Status, expectedVersion).Scan(&updated.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Lead{}, domain.Lead{}, ErrVersionConflict
	}
	return current, updated, err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanLead(scanner rowScanner, lead *domain.Lead) error {
	var archivedAt sql.NullTime
	err := scanner.Scan(
		&lead.ID,
		&lead.LeadName,
		&lead.CompanyName,
		&lead.Email,
		&lead.Phone,
		&lead.Source,
		&lead.Status,
		&lead.OwnerID,
		&lead.NeedSummary,
		&lead.InvalidReason,
		&lead.ConvertedAccountID,
		&lead.ConvertedOpportunityID,
		&lead.ConversionIDempotencyKey,
		&archivedAt,
		&lead.ArchivedBy,
		&lead.ArchiveReason,
		&lead.Version,
		&lead.UpdatedAt,
	)
	if archivedAt.Valid {
		lead.ArchivedAt = archivedAt.Time
	}
	return err
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
