package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"crm-system/services/audit-history/internal/domain"
)

type EventRepo struct {
	db *sql.DB
}

func NewEventRepo(db *sql.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Append(ctx context.Context, event domain.Event) (domain.Event, error) {
	event = domain.ApplyClassificationAndRetention(event)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Event{}, err
	}
	defer tx.Rollback()

	var duplicateSequenceID int64
	err = tx.QueryRowContext(ctx, `
		SELECT sequence_id
		FROM audit_history.events
		WHERE event_uid = $1
	`, event.EventUID).Scan(&duplicateSequenceID)
	if err == nil {
		return event, tx.Commit()
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return domain.Event{}, err
	}

	var prevHash sql.NullString
	err = tx.QueryRowContext(ctx, `
		SELECT event_hash
		FROM audit_history.events
		ORDER BY sequence_id DESC
		LIMIT 1
	`).Scan(&prevHash)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return domain.Event{}, err
	}
	if prevHash.Valid {
		event.PrevHash = prevHash.String
	}
	event.EventHash, err = event.ComputeHash()
	if err != nil {
		return domain.Event{}, err
	}
	beforeBytes, err := json.Marshal(event.BeforeSummary)
	if err != nil {
		return domain.Event{}, err
	}
	afterBytes, err := json.Marshal(event.AfterSummary)
	if err != nil {
		return domain.Event{}, err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO audit_history.events (
			event_uid, event_id, event_version, producer_service, surfaces,
			actor_user_id, actor_role, actor_display, action,
			resource_type, resource_id, parent_resource_type, parent_resource_id,
			result, reason_code, before_summary, after_summary, diff_classification,
			retention_policy, retain_until, scope_summary, safe_summary, correlation_id, causation_id, acceptance_ids,
			occurred_at, prev_hash, event_hash
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23,
			$24, $25, $26, $27, $28
		)
	`, event.EventUID, event.EventID, event.EventVersion, event.ProducerService, event.Surfaces,
		event.ActorUserID, event.ActorRole, event.ActorDisplay, event.Action,
		event.ResourceType, event.ResourceID, event.ParentResourceType, event.ParentResourceID,
		event.Result, event.ReasonCode, beforeBytes, afterBytes, event.DiffClassification,
		event.RetentionPolicy, event.RetainUntil, event.ScopeSummary, event.SafeSummary, event.CorrelationID, event.CausationID, event.AcceptanceIDs,
		event.OccurredAt, event.PrevHash, event.EventHash)
	if err != nil {
		return domain.Event{}, err
	}
	return event, tx.Commit()
}

func (r *EventRepo) ByRecord(ctx context.Context, resourceType, resourceID string) ([]domain.Event, error) {
	return r.query(ctx, `
		SELECT event_uid, event_id, event_version, producer_service,
			actor_user_id, actor_role, actor_display, action, resource_type, resource_id,
			result, reason_code, before_summary, after_summary, diff_classification, retention_policy, retain_until,
			safe_summary, occurred_at, prev_hash, event_hash
		FROM audit_history.events
		WHERE resource_type = $1 AND resource_id = $2 AND 'record_history' = ANY(surfaces)
		ORDER BY occurred_at ASC, sequence_id ASC
	`, resourceType, resourceID)
}

func (r *EventRepo) OperationLog(ctx context.Context) ([]domain.Event, error) {
	return r.query(ctx, `
		SELECT event_uid, event_id, event_version, producer_service,
			actor_user_id, actor_role, actor_display, action, resource_type, resource_id,
			result, reason_code, before_summary, after_summary, diff_classification, retention_policy, retain_until,
			safe_summary, occurred_at, prev_hash, event_hash
		FROM audit_history.events
		WHERE 'operation_log' = ANY(surfaces)
		ORDER BY occurred_at ASC, sequence_id ASC
	`)
}

func (r *EventRepo) query(ctx context.Context, query string, args ...any) ([]domain.Event, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		var beforeBytes []byte
		var afterBytes []byte
		if err := rows.Scan(
			&event.EventUID,
			&event.EventID,
			&event.EventVersion,
			&event.ProducerService,
			&event.ActorUserID,
			&event.ActorRole,
			&event.ActorDisplay,
			&event.Action,
			&event.ResourceType,
			&event.ResourceID,
			&event.Result,
			&event.ReasonCode,
			&beforeBytes,
			&afterBytes,
			&event.DiffClassification,
			&event.RetentionPolicy,
			&event.RetainUntil,
			&event.SafeSummary,
			&event.OccurredAt,
			&event.PrevHash,
			&event.EventHash,
		); err != nil {
			return nil, err
		}
		if len(beforeBytes) > 0 {
			if err := json.Unmarshal(beforeBytes, &event.BeforeSummary); err != nil {
				return nil, err
			}
		}
		if len(afterBytes) > 0 {
			if err := json.Unmarshal(afterBytes, &event.AfterSummary); err != nil {
				return nil, err
			}
		}
		events = append(events, event)
	}
	return events, rows.Err()
}
