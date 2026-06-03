package repo

import (
	"context"
	"database/sql"
	"time"
)

type ImportRun struct {
	RunID              string
	ObjectType         string
	Filename           string
	Status             string
	ActorID            string
	ActorRole          string
	TeamID             string
	TotalRows          int
	SuccessCount       int
	FailureCount       int
	OperationLogStatus string
	CleanupStatus      string
	RetainedUntil      time.Time
	CompletedAt        *time.Time
	RowErrors          []ImportRowResult
}

type ImportRowResult struct {
	RowNumber      int
	Success        bool
	Field          string
	Code           string
	SafeMessage    string
	TargetRecordID string
}

type RunRepo struct {
	db *sql.DB
}

func NewRunRepo(db *sql.DB) *RunRepo {
	return &RunRepo{db: db}
}

func (r *RunRepo) SaveImportRun(ctx context.Context, run ImportRun) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO import_export.import_runs
			(run_id, object_type, filename, status, actor_id, actor_role, team_id, total_rows,
			 success_count, failure_count, operation_log_status, cleanup_status, retained_until, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, run.RunID, run.ObjectType, run.Filename, run.Status, run.ActorID, run.ActorRole, run.TeamID, run.TotalRows,
		run.SuccessCount, run.FailureCount, run.OperationLogStatus, run.CleanupStatus, run.RetainedUntil, run.CompletedAt)
	if err != nil {
		return err
	}
	for _, row := range run.RowErrors {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO import_export.import_row_results
				(run_id, row_number, success, field, code, safe_message, target_record_id)
			VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, ''))
		`, run.RunID, row.RowNumber, row.Success, row.Field, row.Code, row.SafeMessage, row.TargetRecordID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *RunRepo) FindImportRun(ctx context.Context, runID string) (ImportRun, error) {
	var run ImportRun
	var completedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT run_id, object_type, filename, status, actor_id, actor_role, team_id, total_rows,
		       success_count, failure_count, operation_log_status, cleanup_status, retained_until, completed_at
		FROM import_export.import_runs
		WHERE run_id = $1
	`, runID).Scan(&run.RunID, &run.ObjectType, &run.Filename, &run.Status, &run.ActorID, &run.ActorRole, &run.TeamID,
		&run.TotalRows, &run.SuccessCount, &run.FailureCount, &run.OperationLogStatus, &run.CleanupStatus, &run.RetainedUntil, &completedAt)
	if err != nil {
		return ImportRun{}, err
	}
	if completedAt.Valid {
		run.CompletedAt = &completedAt.Time
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT row_number, success, COALESCE(field, ''), COALESCE(code, ''), COALESCE(safe_message, ''), COALESCE(target_record_id, '')
		FROM import_export.import_row_results
		WHERE run_id = $1
		ORDER BY row_number, id
	`, runID)
	if err != nil {
		return ImportRun{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var row ImportRowResult
		if err := rows.Scan(&row.RowNumber, &row.Success, &row.Field, &row.Code, &row.SafeMessage, &row.TargetRecordID); err != nil {
			return ImportRun{}, err
		}
		run.RowErrors = append(run.RowErrors, row)
	}
	if err := rows.Err(); err != nil {
		return ImportRun{}, err
	}
	return run, nil
}

func (r *RunRepo) MarkExpiredRunsDeleted(ctx context.Context, now time.Time) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE import_export.import_runs
		SET cleanup_status = 'deleted'
		WHERE retained_until <= $1 AND cleanup_status <> 'deleted'
	`, now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
