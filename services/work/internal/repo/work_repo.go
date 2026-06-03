package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"

	"crm-system/services/work/internal/domain"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrVersionConflict = errors.New("version conflict")
)

type WorkRepo struct {
	db *sql.DB
	q  sqlRunner
}

func NewWorkRepo(db *sql.DB) *WorkRepo {
	return &WorkRepo{db: db, q: db}
}

func NewWorkRepoTx(tx *sql.Tx) *WorkRepo {
	return &WorkRepo{q: tx}
}

type sqlRunner interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *WorkRepo) CreateActivity(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
	activity.ID = "activity_" + randomHex(16)
	err := r.q.QueryRowContext(ctx, `
		INSERT INTO work.activities (id, related_type, related_id, activity_type, content, actor_id, owner_id, occurred_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 1)
		RETURNING updated_at
	`, activity.ID, activity.RelatedType, activity.RelatedID, activity.ActivityType, activity.Content, activity.ActorID, activity.OwnerID, activity.OccurredAt).Scan(&activity.UpdatedAt)
	return activity, err
}

func (r *WorkRepo) CreateNote(ctx context.Context, note domain.Note) (domain.Note, error) {
	note.ID = "note_" + randomHex(16)
	err := r.q.QueryRowContext(ctx, `
		INSERT INTO work.notes (id, related_type, related_id, content, actor_id, owner_id, occurred_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 1)
		RETURNING updated_at
	`, note.ID, note.RelatedType, note.RelatedID, note.Content, note.ActorID, note.OwnerID, note.OccurredAt).Scan(&note.UpdatedAt)
	return note, err
}

func (r *WorkRepo) ListActivities(ctx context.Context, actorID, actorRole, relatedType, relatedID string) ([]domain.Activity, error) {
	rows, err := r.q.QueryContext(ctx, `
		SELECT id, related_type, related_id, activity_type, content, actor_id, owner_id, occurred_at, version, updated_at
		FROM work.activities
		WHERE ($1 <> 'Sales' OR owner_id = $2)
		  AND ($3 = '' OR related_type = $3)
		  AND ($4 = '' OR related_id = $4)
		ORDER BY occurred_at DESC, id ASC
	`, actorRole, actorID, relatedType, relatedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activities []domain.Activity
	for rows.Next() {
		var activity domain.Activity
		if err := rows.Scan(&activity.ID, &activity.RelatedType, &activity.RelatedID, &activity.ActivityType, &activity.Content, &activity.ActorID, &activity.OwnerID, &activity.OccurredAt, &activity.Version, &activity.UpdatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, rows.Err()
}

func (r *WorkRepo) ListNotes(ctx context.Context, actorID, actorRole, relatedType, relatedID string) ([]domain.Note, error) {
	rows, err := r.q.QueryContext(ctx, `
		SELECT id, related_type, related_id, content, actor_id, owner_id, occurred_at, version, updated_at
		FROM work.notes
		WHERE ($1 <> 'Sales' OR owner_id = $2)
		  AND ($3 = '' OR related_type = $3)
		  AND ($4 = '' OR related_id = $4)
		ORDER BY occurred_at DESC, id ASC
	`, actorRole, actorID, relatedType, relatedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []domain.Note
	for rows.Next() {
		var note domain.Note
		if err := rows.Scan(&note.ID, &note.RelatedType, &note.RelatedID, &note.Content, &note.ActorID, &note.OwnerID, &note.OccurredAt, &note.Version, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, rows.Err()
}

func (r *WorkRepo) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	task.ID = "task_" + randomHex(16)
	err := r.q.QueryRowContext(ctx, `
		INSERT INTO work.tasks (id, related_type, related_id, title, due_date, status, actor_id, owner_id, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 1)
		RETURNING updated_at
	`, task.ID, task.RelatedType, task.RelatedID, task.Title, task.DueDate, task.Status, task.ActorID, task.OwnerID).Scan(&task.UpdatedAt)
	return task, err
}

func (r *WorkRepo) FindTask(ctx context.Context, id string) (domain.Task, error) {
	var task domain.Task
	var completedAt sql.NullTime
	var cancelledAt sql.NullTime
	var cancellationReason sql.NullString
	err := r.q.QueryRowContext(ctx, `
		SELECT id, related_type, related_id, title, due_date, status, actor_id, owner_id, completed_at, cancelled_at, cancellation_reason, version, updated_at
		FROM work.tasks
		WHERE id = $1
	`, id).Scan(&task.ID, &task.RelatedType, &task.RelatedID, &task.Title, &task.DueDate, &task.Status, &task.ActorID, &task.OwnerID,
		&completedAt, &cancelledAt, &cancellationReason, &task.Version, &task.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Task{}, ErrNotFound
	}
	if completedAt.Valid {
		task.CompletedAt = completedAt.Time
	}
	if cancelledAt.Valid {
		task.CancelledAt = cancelledAt.Time
	}
	task.CancellationReason = cancellationReason.String
	return task, err
}

func (r *WorkRepo) UpdateTaskStatus(ctx context.Context, task domain.Task) (domain.Task, error) {
	var completedAt sql.NullTime
	if !task.CompletedAt.IsZero() {
		completedAt = sql.NullTime{Time: task.CompletedAt, Valid: true}
	}
	var cancelledAt sql.NullTime
	if !task.CancelledAt.IsZero() {
		cancelledAt = sql.NullTime{Time: task.CancelledAt, Valid: true}
	}
	err := r.q.QueryRowContext(ctx, `
		UPDATE work.tasks
		SET status = $2,
		    completed_at = $3,
		    cancelled_at = $4,
		    cancellation_reason = NULLIF($5, ''),
		    version = version + 1,
		    updated_at = now()
		WHERE id = $1 AND version = $6
		RETURNING version, updated_at
	`, task.ID, task.Status, completedAt, cancelledAt, task.CancellationReason, task.Version).Scan(&task.Version, &task.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Task{}, ErrVersionConflict
	}
	return task, err
}

func (r *WorkRepo) ListTasks(ctx context.Context, actorID, actorRole, relatedType, relatedID string, activeOnly bool) ([]domain.Task, error) {
	rows, err := r.q.QueryContext(ctx, `
		SELECT id, related_type, related_id, title, due_date, status, actor_id, owner_id, completed_at, cancelled_at, cancellation_reason, version, updated_at
		FROM work.tasks
		WHERE ($1 <> 'Sales' OR owner_id = $2)
		  AND ($3 = '' OR related_type = $3)
		  AND ($4 = '' OR related_id = $4)
		  AND ($5 = false OR status = 'Open')
		ORDER BY due_date ASC, updated_at DESC
	`, actorRole, actorID, relatedType, relatedID, activeOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		var completedAt sql.NullTime
		var cancelledAt sql.NullTime
		var cancellationReason sql.NullString
		if err := rows.Scan(&task.ID, &task.RelatedType, &task.RelatedID, &task.Title, &task.DueDate, &task.Status, &task.ActorID, &task.OwnerID,
			&completedAt, &cancelledAt, &cancellationReason, &task.Version, &task.UpdatedAt); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			task.CompletedAt = completedAt.Time
		}
		if cancelledAt.Valid {
			task.CancelledAt = cancelledAt.Time
		}
		task.CancellationReason = cancellationReason.String
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r *WorkRepo) TransferOpenWork(ctx context.Context, relatedType, relatedID, fromOwnerID, toOwnerID string) (int64, error) {
	result, err := r.q.ExecContext(ctx, `
		UPDATE work.tasks
		SET owner_id = $4,
		    version = version + 1,
		    updated_at = now()
		WHERE related_type = $1 AND related_id = $2 AND owner_id = $3 AND status = 'Open'
	`, relatedType, relatedID, fromOwnerID, toOwnerID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *WorkRepo) ActiveObligations(ctx context.Context, relatedType, relatedID string) ([]domain.Task, error) {
	rows, err := r.q.QueryContext(ctx, `
		SELECT id, related_type, related_id, title, due_date, status, actor_id, owner_id, completed_at, cancelled_at, cancellation_reason, version, updated_at
		FROM work.tasks
		WHERE related_type = $1 AND related_id = $2 AND status = 'Open'
		ORDER BY due_date ASC, id ASC
	`, relatedType, relatedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		var completedAt, cancelledAt sql.NullTime
		var cancellationReason sql.NullString
		if err := rows.Scan(&task.ID, &task.RelatedType, &task.RelatedID, &task.Title, &task.DueDate, &task.Status, &task.ActorID, &task.OwnerID, &completedAt, &cancelledAt, &cancellationReason, &task.Version, &task.UpdatedAt); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			task.CompletedAt = completedAt.Time
		}
		if cancelledAt.Valid {
			task.CancelledAt = cancelledAt.Time
		}
		task.CancellationReason = cancellationReason.String
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
