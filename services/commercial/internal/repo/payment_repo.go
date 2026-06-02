package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"crm-system/services/commercial/internal/domain"
)

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) CreatePlan(ctx context.Context, plan domain.PaymentPlan) (domain.PaymentPlan, error) {
	plan.ID = "plan_" + randomHex(16)
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO commercial.payment_plans (id, contract_id, due_amount, due_date, currency, status, version)
		VALUES ($1, $2, $3::numeric, $4, $5, $6, 1)
		RETURNING updated_at
	`, plan.ID, plan.ContractID, plan.DueAmount, plan.DueDate, plan.Currency, plan.Status).Scan(&plan.UpdatedAt)
	if err != nil {
		return domain.PaymentPlan{}, err
	}
	return plan, nil
}

func (r *PaymentRepo) FindPlan(ctx context.Context, id string) (domain.PaymentPlan, error) {
	var plan domain.PaymentPlan
	err := scanPaymentPlan(r.db.QueryRowContext(ctx, `
		SELECT id, contract_id, to_char(due_amount, 'FM999999999999990.00'), due_date, currency, status,
		       archived_at, archived_by, archive_reason, version, updated_at
		FROM commercial.payment_plans
		WHERE id = $1
	`, id), &plan)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.PaymentPlan{}, ErrNotFound
	}
	return plan, err
}

func (r *PaymentRepo) ReminderRows(ctx context.Context, actorID, actorRole string, businessDate time.Time) ([]domain.ReminderRow, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id, p.contract_id, c.opportunity_id, p.due_date, p.status, p.version, c.owner_id
		FROM commercial.payment_plans p
		JOIN commercial.contracts c ON c.id = p.contract_id
		WHERE p.status <> $1
		  AND p.archived_at IS NULL
		  AND c.archived_at IS NULL
		  AND p.due_date <= $2
		  AND c.status NOT IN ($3, $4)
		  AND ($5 <> 'Sales' OR c.owner_id = $6)
		ORDER BY p.due_date ASC, p.id ASC
	`, domain.PaymentStatusPaid, businessDate, domain.ContractStatusCompleted, domain.ContractStatusTerminated, actorRole, actorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reminders []domain.ReminderRow
	for rows.Next() {
		var planID string
		var contractID string
		var opportunityID string
		var dueDate time.Time
		var status string
		var version int
		var ownerID string
		if err := rows.Scan(&planID, &contractID, &opportunityID, &dueDate, &status, &version, &ownerID); err != nil {
			return nil, err
		}
		reminderType := "payment_due"
		reminderStatus := "DueToday"
		if dueDate.Before(businessDate) {
			reminderType = "payment_overdue"
			reminderStatus = "Overdue"
		}
		_ = status
		reminders = append(reminders, domain.ReminderRow{
			ID:            planID,
			SourceService: "commercial-service",
			Type:          reminderType,
			RelatedRecord: domain.ReminderRelatedRecord{Type: "contract", ID: contractID, Display: opportunityID},
			OwnerDisplay:  ownerID,
			DueDate:       domain.FormatDate(dueDate),
			Status:        reminderStatus,
			Priority:      "P1",
			Version:       version,
		})
	}
	return reminders, rows.Err()
}

func (r *PaymentRepo) PaidTotal(ctx context.Context, contractID string) (string, error) {
	var total string
	err := r.db.QueryRowContext(ctx, `
		SELECT to_char(COALESCE(sum(amount), 0), 'FM999999999999990.00')
		FROM commercial.actual_payments
		WHERE contract_id = $1
	`, contractID).Scan(&total)
	return total, err
}

func (r *PaymentRepo) FindPaymentByKey(ctx context.Context, contractID, idempotencyKey string) (domain.ActualPayment, error) {
	var payment domain.ActualPayment
	err := r.db.QueryRowContext(ctx, `
		SELECT id, contract_id, idempotency_key, to_char(amount, 'FM999999999999990.00'), payment_date, note,
		       currency, payment_status, to_char(remaining_amount, 'FM999999999999990.00'), version, updated_at
		FROM commercial.actual_payments
		WHERE contract_id = $1 AND idempotency_key = $2
	`, contractID, idempotencyKey).Scan(&payment.ID, &payment.ContractID, &payment.IdempotencyKey, &payment.Amount,
		&payment.PaymentDate, &payment.Note, &payment.Currency, &payment.PaymentStatus, &payment.RemainingAmount,
		&payment.Version, &payment.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ActualPayment{}, ErrNotFound
	}
	return payment, err
}

func (r *PaymentRepo) CreatePayment(ctx context.Context, payment domain.ActualPayment) (domain.ActualPayment, error) {
	payment.ID = "payment_" + randomHex(16)
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO commercial.actual_payments
			(id, contract_id, idempotency_key, amount, payment_date, note, currency, payment_status, remaining_amount, version)
		VALUES ($1, $2, $3, $4::numeric, $5, $6, $7, $8, $9::numeric, 1)
		RETURNING updated_at
	`, payment.ID, payment.ContractID, payment.IdempotencyKey, payment.Amount, payment.PaymentDate, payment.Note,
		payment.Currency, payment.PaymentStatus, payment.RemainingAmount).Scan(&payment.UpdatedAt)
	if err != nil {
		return domain.ActualPayment{}, err
	}
	_, _ = r.db.ExecContext(ctx, `
		UPDATE commercial.payment_plans
		SET status = $2,
		    version = version + 1,
		    updated_at = now()
		WHERE contract_id = $1
	`, payment.ContractID, payment.PaymentStatus)
	return payment, nil
}

func (r *PaymentRepo) ArchivePlan(ctx context.Context, id string, expectedVersion int, actorID, reason string) (domain.PaymentPlan, error) {
	var plan domain.PaymentPlan
	err := scanPaymentPlan(r.db.QueryRowContext(ctx, `
		UPDATE commercial.payment_plans
		SET archived_at = now(),
		    archived_by = $2,
		    archive_reason = $3,
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $4 AND archived_at IS NULL
		RETURNING id, contract_id, to_char(due_amount, 'FM999999999999990.00'), due_date, currency, status,
		          archived_at, archived_by, archive_reason, version, updated_at
	`, id, actorID, reason, expectedVersion), &plan)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.PaymentPlan{}, ErrVersionConflict
	}
	return plan, err
}

type paymentPlanRowScanner interface {
	Scan(dest ...any) error
}

func scanPaymentPlan(scanner paymentPlanRowScanner, plan *domain.PaymentPlan) error {
	var archivedAt sql.NullTime
	err := scanner.Scan(
		&plan.ID,
		&plan.ContractID,
		&plan.DueAmount,
		&plan.DueDate,
		&plan.Currency,
		&plan.Status,
		&archivedAt,
		&plan.ArchivedBy,
		&plan.ArchiveReason,
		&plan.Version,
		&plan.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if archivedAt.Valid {
		plan.ArchivedAt = archivedAt.Time
	}
	return nil
}
