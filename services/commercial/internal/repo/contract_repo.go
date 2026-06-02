package repo

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"crm-system/services/commercial/internal/domain"
)

type ContractRepo struct {
	db *sql.DB
}

func NewContractRepo(db *sql.DB) *ContractRepo {
	return &ContractRepo{db: db}
}

func (r *ContractRepo) Create(ctx context.Context, contract domain.Contract) (domain.Contract, error) {
	contract.ID = "contract_" + randomHex(16)
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO commercial.contracts
			(id, quote_id, opportunity_id, customer_id, amount, status, contract_note, expected_signed_date,
			 amount_difference_reason, owner_id, version)
		VALUES ($1, $2, $3, $4, $5::numeric, $6, $7, $8, NULLIF($9, ''), $10, 1)
		RETURNING updated_at
	`, contract.ID, contract.QuoteID, contract.OpportunityID, contract.CustomerID, contract.Amount,
		contract.Status, contract.ContractNote, contract.ExpectedSignedDate, contract.AmountDifferenceReason, contract.OwnerID).Scan(&contract.UpdatedAt)
	if err != nil && strings.Contains(err.Error(), "contracts_quote_unique") {
		return domain.Contract{}, domain.ErrContractAlreadyExists
	}
	if err != nil {
		return domain.Contract{}, err
	}
	return contract, nil
}

func (r *ContractRepo) Find(ctx context.Context, id string) (domain.Contract, error) {
	var contract domain.Contract
	err := scanContract(r.db.QueryRowContext(ctx, `
		SELECT id, quote_id, opportunity_id, customer_id, to_char(amount, 'FM999999999999990.00'), status,
		       contract_note, expected_signed_date, signed_effective_date, amount_difference_reason, owner_id,
		       archived_at, archived_by, archive_reason, version, updated_at
		FROM commercial.contracts
		WHERE id = $1
	`, id), &contract)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Contract{}, ErrNotFound
	}
	return contract, err
}

func (r *ContractRepo) List(ctx context.Context, actorID, actorRole, search string, includeArchived bool) ([]domain.Contract, error) {
	search = strings.TrimSpace(search)
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, quote_id, opportunity_id, customer_id, to_char(amount, 'FM999999999999990.00'), status,
		       contract_note, expected_signed_date, signed_effective_date, amount_difference_reason, owner_id,
		       archived_at, archived_by, archive_reason, version, updated_at
		FROM commercial.contracts
		WHERE ($1 <> 'Sales' OR owner_id = $2)
		  AND ($3 = '' OR id ILIKE '%' || $3 || '%' OR quote_id ILIKE '%' || $3 || '%' OR opportunity_id ILIKE '%' || $3 || '%' OR customer_id ILIKE '%' || $3 || '%')
		  AND ($4 = true OR archived_at IS NULL)
		ORDER BY updated_at DESC, id ASC
	`, actorRole, actorID, search, includeArchived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contracts []domain.Contract
	for rows.Next() {
		var contract domain.Contract
		if err := scanContract(rows, &contract); err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	return contracts, rows.Err()
}

func (r *ContractRepo) PendingSignatureReminderRows(ctx context.Context, actorID, actorRole string, businessDate time.Time) ([]domain.ReminderRow, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, opportunity_id, expected_signed_date, owner_id, version
		FROM commercial.contracts
		WHERE status = $1
		  AND archived_at IS NULL
		  AND expected_signed_date < $2
		  AND ($3 <> 'Sales' OR owner_id = $4)
		ORDER BY expected_signed_date ASC, id ASC
	`, domain.ContractStatusPendingSignature, businessDate, actorRole, actorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reminders []domain.ReminderRow
	for rows.Next() {
		var contractID string
		var opportunityID string
		var dueDate time.Time
		var ownerID string
		var version int
		if err := rows.Scan(&contractID, &opportunityID, &dueDate, &ownerID, &version); err != nil {
			return nil, err
		}
		reminders = append(reminders, domain.ReminderRow{
			ID:            contractID,
			SourceService: "commercial-service",
			Type:          "contract_pending_signature",
			RelatedRecord: domain.ReminderRelatedRecord{Type: "contract", ID: contractID, Display: opportunityID},
			OwnerDisplay:  ownerID,
			DueDate:       domain.FormatDate(dueDate),
			Status:        "Overdue",
			Priority:      "P1",
			Version:       version,
		})
	}
	return reminders, rows.Err()
}

func (r *ContractRepo) ChangeStatus(ctx context.Context, id string, expectedVersion int, toStatus string, signedEffectiveDate time.Time) (domain.Contract, error) {
	var dateParam sql.NullTime
	if !signedEffectiveDate.IsZero() {
		dateParam = sql.NullTime{Time: signedEffectiveDate, Valid: true}
	}
	var contract domain.Contract
	err := scanContract(r.db.QueryRowContext(ctx, `
		UPDATE commercial.contracts
		SET status = $2,
		    signed_effective_date = COALESCE($4, signed_effective_date),
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $3
		RETURNING id, quote_id, opportunity_id, customer_id, to_char(amount, 'FM999999999999990.00'), status,
		          contract_note, expected_signed_date, signed_effective_date, amount_difference_reason, owner_id,
		          archived_at, archived_by, archive_reason, version, updated_at
	`, id, toStatus, expectedVersion, dateParam), &contract)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Contract{}, ErrVersionConflict
	}
	return contract, nil
}

func (r *ContractRepo) Archive(ctx context.Context, id string, expectedVersion int, actorID, reason string) (domain.Contract, error) {
	var contract domain.Contract
	err := scanContract(r.db.QueryRowContext(ctx, `
		UPDATE commercial.contracts
		SET archived_at = now(),
		    archived_by = $2,
		    archive_reason = $3,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $4 AND archived_at IS NULL
		RETURNING id, quote_id, opportunity_id, customer_id, to_char(amount, 'FM999999999999990.00'), status,
		          contract_note, expected_signed_date, signed_effective_date, amount_difference_reason, owner_id,
		          archived_at, archived_by, archive_reason, version, updated_at
	`, id, actorID, strings.TrimSpace(reason), expectedVersion), &contract)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Contract{}, ErrVersionConflict
	}
	return contract, err
}

type contractRowScanner interface {
	Scan(dest ...any) error
}

func scanContract(scanner contractRowScanner, contract *domain.Contract) error {
	var amountDifferenceReason sql.NullString
	var signedEffectiveDate sql.NullTime
	var archivedAt sql.NullTime
	err := scanner.Scan(
		&contract.ID,
		&contract.QuoteID,
		&contract.OpportunityID,
		&contract.CustomerID,
		&contract.Amount,
		&contract.Status,
		&contract.ContractNote,
		&contract.ExpectedSignedDate,
		&signedEffectiveDate,
		&amountDifferenceReason,
		&contract.OwnerID,
		&archivedAt,
		&contract.ArchivedBy,
		&contract.ArchiveReason,
		&contract.Version,
		&contract.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if signedEffectiveDate.Valid {
		contract.SignedEffectiveDate = signedEffectiveDate.Time
	}
	contract.AmountDifferenceReason = amountDifferenceReason.String
	if archivedAt.Valid {
		contract.ArchivedAt = archivedAt.Time
	}
	return nil
}
