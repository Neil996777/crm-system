package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"crm-system/services/account/internal/domain"
)

var (
	ErrDuplicateTokenInvalid = errors.New("duplicate warning token invalid")
	ErrDuplicateTokenUsed    = errors.New("duplicate warning token used")
)

type DuplicateRepo struct {
	db *sql.DB
}

func NewDuplicateRepo(db *sql.DB) *DuplicateRepo {
	return &DuplicateRepo{db: db}
}

func (r *DuplicateRepo) Check(ctx context.Context, actorID, actorRole string, candidate domain.DuplicateCandidate) (domain.DuplicateCheckResult, error) {
	signature, fields := domain.DuplicateSignature(candidate)
	result := domain.DuplicateCheckResult{Result: "NoDuplicate", NormalizedFields: fields, Signature: signature}
	if signature == "" {
		return result, nil
	}
	matches, rules, err := r.findMatches(ctx, actorID, actorRole, candidate)
	if err != nil {
		return domain.DuplicateCheckResult{}, err
	}
	if len(matches) == 0 {
		return result, nil
	}
	result.Result = "PossibleDuplicate"
	result.Matches = matches
	result.Rules = rules
	result.WarningToken = "dup_" + randomHex(16)
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO account.duplicate_warning_tokens (token, target_type, normalized_signature, actor_user_id, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, result.WarningToken, candidate.TargetType, signature, actorID, time.Now().UTC().Add(30*time.Minute))
	return result, err
}

func (r *DuplicateRepo) ConsumeToken(ctx context.Context, token, targetType, actorID, signature string) error {
	var usedAt sql.NullTime
	var expiresAt time.Time
	err := r.db.QueryRowContext(ctx, `
		SELECT used_at, expires_at
		FROM account.duplicate_warning_tokens
		WHERE token = $1 AND target_type = $2 AND actor_user_id = $3 AND normalized_signature = $4
	`, token, targetType, actorID, signature).Scan(&usedAt, &expiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrDuplicateTokenInvalid
	}
	if err != nil {
		return err
	}
	if usedAt.Valid {
		return ErrDuplicateTokenUsed
	}
	if time.Now().UTC().After(expiresAt) {
		return ErrDuplicateTokenInvalid
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE account.duplicate_warning_tokens
		SET used_at = now()
		WHERE token = $1 AND used_at IS NULL
	`, token)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrDuplicateTokenUsed
	}
	return nil
}

func (r *DuplicateRepo) findMatches(ctx context.Context, actorID, actorRole string, candidate domain.DuplicateCandidate) ([]domain.DuplicateMatch, []string, error) {
	matches := make([]domain.DuplicateMatch, 0, 3)
	rules := make([]string, 0, 3)
	if candidate.TargetType == "account" {
		if company := domain.NormalizeCompanyName(candidate.CompanyName); company != "" {
			count, err := r.countAccountsByCompany(ctx, actorID, actorRole, company)
			if err != nil {
				return nil, nil, err
			}
			if count > 0 {
				matches = append(matches, safeMatch("account", "COMPANY_NAME_MATCH"))
				rules = append(rules, "COMPANY_NAME_MATCH")
			}
		}
	}
	if candidate.TargetType == "contact" {
		if email := domain.NormalizeEmail(candidate.Email); email != "" {
			count, err := r.countContactsByEmail(ctx, actorID, actorRole, email)
			if err != nil {
				return nil, nil, err
			}
			if count > 0 {
				matches = append(matches, safeMatch("contact", "CONTACT_EMAIL_MATCH"))
				rules = append(rules, "CONTACT_EMAIL_MATCH")
			}
		}
		if phone := domain.NormalizePhone(candidate.Phone); phone != "" {
			count, err := r.countContactsByPhone(ctx, actorID, actorRole, phone)
			if err != nil {
				return nil, nil, err
			}
			if count > 0 {
				matches = append(matches, safeMatch("contact", "CONTACT_PHONE_MATCH"))
				rules = append(rules, "CONTACT_PHONE_MATCH")
			}
		}
	}
	return matches, rules, nil
}

func (r *DuplicateRepo) countAccountsByCompany(ctx context.Context, actorID, actorRole, normalizedCompany string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM account.accounts
		WHERE lower(regexp_replace(btrim(company_name), '\s+', ' ', 'g')) = $1
		  AND ($2 <> 'Sales' OR owner_id = $3)
	`, normalizedCompany, actorRole, actorID).Scan(&count)
	return count, err
}

func (r *DuplicateRepo) countContactsByEmail(ctx context.Context, actorID, actorRole, normalizedEmail string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM account.contacts c
		JOIN account.accounts a ON a.id = c.account_id
		WHERE lower(btrim(c.email)) = $1
		  AND ($2 <> 'Sales' OR a.owner_id = $3)
	`, normalizedEmail, actorRole, actorID).Scan(&count)
	return count, err
}

func (r *DuplicateRepo) countContactsByPhone(ctx context.Context, actorID, actorRole, normalizedPhone string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM account.contacts c
		JOIN account.accounts a ON a.id = c.account_id
		WHERE regexp_replace(c.phone, '\D', '', 'g') IN ($1, '86' || $1)
		  AND ($2 <> 'Sales' OR a.owner_id = $3)
	`, normalizedPhone, actorRole, actorID).Scan(&count)
	return count, err
}

func safeMatch(targetType, rule string) domain.DuplicateMatch {
	return domain.DuplicateMatch{
		Type:          targetType,
		Rule:          rule,
		MatchStrength: "High",
		SafeSummary:   "Possible matching " + targetType,
		Visible:       true,
	}
}
