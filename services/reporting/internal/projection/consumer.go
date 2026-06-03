package projection

import (
	"context"
	"database/sql"
	"strings"
)

type RecordProjection struct {
	SourceService string
	RecordType    string
	RecordID      string
	OwnerID       string
	TeamID        string
	Status        string
	Stage         string
	Amount        string
}

type Consumer struct {
	db *sql.DB
}

func NewConsumer(db *sql.DB) *Consumer {
	return &Consumer{db: db}
}

func (c *Consumer) Upsert(ctx context.Context, projection RecordProjection) error {
	if strings.TrimSpace(projection.TeamID) == "" {
		projection.TeamID = "single-team"
	}
	if strings.TrimSpace(projection.Amount) == "" {
		projection.Amount = "0.00"
	}
	_, err := c.db.ExecContext(ctx, `
		INSERT INTO reporting.record_projections
			(source_service, record_type, record_id, owner_id, team_id, status, stage, amount, updated_at)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), NULLIF($7, ''), $8::numeric, now())
		ON CONFLICT (record_type, record_id) DO UPDATE
		SET source_service = EXCLUDED.source_service,
		    owner_id = EXCLUDED.owner_id,
		    team_id = EXCLUDED.team_id,
		    status = EXCLUDED.status,
		    stage = EXCLUDED.stage,
		    amount = EXCLUDED.amount,
		    archived_at = NULL,
		    updated_at = now()
	`, projection.SourceService, projection.RecordType, projection.RecordID, projection.OwnerID, projection.TeamID, projection.Status, projection.Stage, projection.Amount)
	return err
}
