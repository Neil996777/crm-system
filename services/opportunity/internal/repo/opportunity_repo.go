package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"

	"crm-system/services/opportunity/internal/domain"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrVersionConflict = errors.New("version conflict")
)

type OpportunityRepo struct {
	db *sql.DB
}

func NewOpportunityRepo(db *sql.DB) *OpportunityRepo {
	return &OpportunityRepo{db: db}
}

func (r *OpportunityRepo) Create(ctx context.Context, opportunity domain.Opportunity) (domain.Opportunity, error) {
	opportunity.ID = "opp_" + randomHex(16)
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO opportunity.opportunities
			(id, customer_id, owner_id, stage, expected_amount, expected_close_date, title, version)
		VALUES ($1, $2, $3, $4, $5::numeric, $6, $7, 1)
		RETURNING updated_at
	`, opportunity.ID, opportunity.CustomerID, opportunity.OwnerID, opportunity.Stage, opportunity.ExpectedAmount, opportunity.ExpectedCloseDate, opportunity.Title).Scan(&opportunity.UpdatedAt)
	if err != nil {
		return domain.Opportunity{}, err
	}
	return opportunity, nil
}

func (r *OpportunityRepo) Find(ctx context.Context, id string) (domain.Opportunity, error) {
	var opportunity domain.Opportunity
	err := scanOpportunityRow(r.db.QueryRowContext(ctx, `
		SELECT id, customer_id, owner_id, stage, to_char(expected_amount, 'FM999999999999990.00'), expected_close_date, title,
		       close_date, won_contract_id, lost_reason_code, lost_reason_detail, closed_at,
		       archived_at, archived_by, archive_reason, version, updated_at
		FROM opportunity.opportunities
		WHERE id = $1
	`, id), &opportunity)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Opportunity{}, ErrNotFound
	}
	return opportunity, err
}

func (r *OpportunityRepo) List(ctx context.Context, actorID, actorRole, search, stage string, includeArchived bool) ([]domain.Opportunity, error) {
	search = strings.TrimSpace(search)
	stage = strings.TrimSpace(stage)
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, customer_id, owner_id, stage, to_char(expected_amount, 'FM999999999999990.00'), expected_close_date, title,
		       close_date, won_contract_id, lost_reason_code, lost_reason_detail, closed_at,
		       archived_at, archived_by, archive_reason, version, updated_at
		FROM opportunity.opportunities
		WHERE ($1 <> 'Sales' OR owner_id = $2)
		  AND ($3 = '' OR customer_id ILIKE '%' || $3 || '%' OR title ILIKE '%' || $3 || '%')
		  AND ($4 = '' OR stage = $4)
		  AND ($5 = true OR archived_at IS NULL)
		ORDER BY updated_at DESC, id ASC
	`, actorRole, actorID, search, stage, includeArchived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var opportunities []domain.Opportunity
	for rows.Next() {
		var opportunity domain.Opportunity
		if err := scanOpportunityRow(rows, &opportunity); err != nil {
			return nil, err
		}
		opportunities = append(opportunities, opportunity)
	}
	return opportunities, rows.Err()
}

func (r *OpportunityRepo) Update(ctx context.Context, id string, expectedVersion int, updated domain.Opportunity) (domain.Opportunity, error) {
	err := scanOpportunityRow(r.db.QueryRowContext(ctx, `
		UPDATE opportunity.opportunities
		SET customer_id = $2,
		    owner_id = $3,
		    stage = $4,
		    expected_amount = $5::numeric,
		    expected_close_date = $6,
		    title = $7,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $8
		RETURNING id, customer_id, owner_id, stage, to_char(expected_amount, 'FM999999999999990.00'), expected_close_date, title,
		          close_date, won_contract_id, lost_reason_code, lost_reason_detail, closed_at,
		          archived_at, archived_by, archive_reason, version, updated_at
	`, id, updated.CustomerID, updated.OwnerID, updated.Stage, updated.ExpectedAmount, updated.ExpectedCloseDate, updated.Title, expectedVersion), &updated)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Opportunity{}, ErrVersionConflict
	}
	if err != nil {
		return domain.Opportunity{}, err
	}
	return updated, nil
}

func (r *OpportunityRepo) ChangeStage(ctx context.Context, id string, expectedVersion int, toStage string) (domain.Opportunity, error) {
	var updated domain.Opportunity
	err := scanOpportunityRow(r.db.QueryRowContext(ctx, `
		UPDATE opportunity.opportunities
		SET stage = $2,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $3
		RETURNING id, customer_id, owner_id, stage, to_char(expected_amount, 'FM999999999999990.00'), expected_close_date, title,
		          close_date, won_contract_id, lost_reason_code, lost_reason_detail, closed_at,
		          archived_at, archived_by, archive_reason, version, updated_at
	`, id, toStage, expectedVersion), &updated)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Opportunity{}, ErrVersionConflict
	}
	return updated, err
}

func (r *OpportunityRepo) Close(ctx context.Context, id string, expectedVersion int, closed domain.Opportunity) (domain.Opportunity, error) {
	err := scanOpportunityRow(r.db.QueryRowContext(ctx, `
		UPDATE opportunity.opportunities
		SET stage = $2,
		    close_date = $3,
		    won_contract_id = NULLIF($4, ''),
		    lost_reason_code = NULLIF($5, ''),
		    lost_reason_detail = NULLIF($6, ''),
		    closed_at = $7,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $8
		RETURNING id, customer_id, owner_id, stage, to_char(expected_amount, 'FM999999999999990.00'), expected_close_date, title,
		          close_date, won_contract_id, lost_reason_code, lost_reason_detail, closed_at,
		          archived_at, archived_by, archive_reason, version, updated_at
	`, id, closed.Stage, closed.CloseDate, closed.WonContractID, closed.LostReasonCode, closed.LostReasonDetail, closed.ClosedAt, expectedVersion), &closed)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Opportunity{}, ErrVersionConflict
	}
	if err != nil {
		return domain.Opportunity{}, err
	}
	return closed, nil
}

func (r *OpportunityRepo) Archive(ctx context.Context, id string, expectedVersion int, actorID, reason string) (domain.Opportunity, error) {
	var archived domain.Opportunity
	err := scanOpportunityRow(r.db.QueryRowContext(ctx, `
		UPDATE opportunity.opportunities
		SET archived_at = now(),
		    archived_by = $2,
		    archive_reason = $3,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $4 AND archived_at IS NULL
		RETURNING id, customer_id, owner_id, stage, to_char(expected_amount, 'FM999999999999990.00'), expected_close_date, title,
		          close_date, won_contract_id, lost_reason_code, lost_reason_detail, closed_at,
		          archived_at, archived_by, archive_reason, version, updated_at
	`, id, actorID, strings.TrimSpace(reason), expectedVersion), &archived)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Opportunity{}, ErrVersionConflict
	}
	return archived, err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanOpportunityRow(scanner rowScanner, opportunity *domain.Opportunity) error {
	var closeDate sql.NullTime
	var wonContractID sql.NullString
	var lostReasonCode sql.NullString
	var lostReasonDetail sql.NullString
	var closedAt sql.NullTime
	var archivedAt sql.NullTime
	err := scanner.Scan(
		&opportunity.ID,
		&opportunity.CustomerID,
		&opportunity.OwnerID,
		&opportunity.Stage,
		&opportunity.ExpectedAmount,
		&opportunity.ExpectedCloseDate,
		&opportunity.Title,
		&closeDate,
		&wonContractID,
		&lostReasonCode,
		&lostReasonDetail,
		&closedAt,
		&archivedAt,
		&opportunity.ArchivedBy,
		&opportunity.ArchiveReason,
		&opportunity.Version,
		&opportunity.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if closeDate.Valid {
		opportunity.CloseDate = closeDate.Time
	}
	opportunity.WonContractID = wonContractID.String
	opportunity.LostReasonCode = lostReasonCode.String
	opportunity.LostReasonDetail = lostReasonDetail.String
	if closedAt.Valid {
		opportunity.ClosedAt = closedAt.Time
	}
	if archivedAt.Valid {
		opportunity.ArchivedAt = archivedAt.Time
	}
	return nil
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
