package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"

	"crm-system/services/commercial/internal/domain"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrVersionConflict = errors.New("version conflict")
)

type QuoteRepo struct {
	db *sql.DB
}

func NewQuoteRepo(db *sql.DB) *QuoteRepo {
	return &QuoteRepo{db: db}
}

func (r *QuoteRepo) ExistsForOpportunity(ctx context.Context, opportunityID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM commercial.quotes WHERE opportunity_id = $1)`, opportunityID).Scan(&exists)
	return exists, err
}

func (r *QuoteRepo) Create(ctx context.Context, quote domain.Quote) (domain.Quote, error) {
	quote.ID = "quote_" + randomHex(16)
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO commercial.quotes
			(id, opportunity_id, customer_id, amount, status, validity_end, owner_id, version)
		VALUES ($1, $2, $3, $4::numeric, $5, $6, $7, 1)
		RETURNING updated_at
	`, quote.ID, quote.OpportunityID, quote.CustomerID, quote.Amount, quote.Status, quote.ValidityEnd, quote.OwnerID).Scan(&quote.UpdatedAt)
	if err != nil && strings.Contains(err.Error(), "quotes_opportunity_unique") {
		return domain.Quote{}, domain.ErrQuoteAlreadyExists
	}
	if err != nil {
		return domain.Quote{}, err
	}
	return quote, nil
}

func (r *QuoteRepo) Find(ctx context.Context, id string) (domain.Quote, error) {
	var quote domain.Quote
	err := r.db.QueryRowContext(ctx, `
		SELECT id, opportunity_id, customer_id, to_char(amount, 'FM999999999999990.00'), status, validity_end, owner_id, version, updated_at
		FROM commercial.quotes
		WHERE id = $1
	`, id).Scan(&quote.ID, &quote.OpportunityID, &quote.CustomerID, &quote.Amount, &quote.Status, &quote.ValidityEnd, &quote.OwnerID, &quote.Version, &quote.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Quote{}, ErrNotFound
	}
	return quote, err
}

func (r *QuoteRepo) List(ctx context.Context, actorID, actorRole, search string) ([]domain.Quote, error) {
	search = strings.TrimSpace(search)
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, opportunity_id, customer_id, to_char(amount, 'FM999999999999990.00'), status, validity_end, owner_id, version, updated_at
		FROM commercial.quotes
		WHERE ($1 <> 'Sales' OR owner_id = $2)
		  AND ($3 = '' OR opportunity_id ILIKE '%' || $3 || '%' OR customer_id ILIKE '%' || $3 || '%')
		ORDER BY updated_at DESC, id ASC
	`, actorRole, actorID, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var quotes []domain.Quote
	for rows.Next() {
		var quote domain.Quote
		if err := rows.Scan(&quote.ID, &quote.OpportunityID, &quote.CustomerID, &quote.Amount, &quote.Status, &quote.ValidityEnd, &quote.OwnerID, &quote.Version, &quote.UpdatedAt); err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	return quotes, rows.Err()
}

func (r *QuoteRepo) ChangeStatus(ctx context.Context, id string, expectedVersion int, toStatus string) (domain.Quote, error) {
	var quote domain.Quote
	err := r.db.QueryRowContext(ctx, `
		UPDATE commercial.quotes
		SET status = $2,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $3
		RETURNING id, opportunity_id, customer_id, to_char(amount, 'FM999999999999990.00'), status, validity_end, owner_id, version, updated_at
	`, id, toStatus, expectedVersion).Scan(&quote.ID, &quote.OpportunityID, &quote.CustomerID, &quote.Amount, &quote.Status, &quote.ValidityEnd, &quote.OwnerID, &quote.Version, &quote.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Quote{}, ErrVersionConflict
	}
	return quote, err
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
