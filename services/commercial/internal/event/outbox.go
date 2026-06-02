package event

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"time"
)

const (
	QuoteCreated          = "QuoteCreated"
	QuoteStatusChanged    = "QuoteStatusChanged"
	QuoteAccepted         = "QuoteAccepted"
	ContractCreated       = "ContractCreated"
	ContractStatusChanged = "ContractStatusChanged"
	ContractArchived      = "ContractArchived"
	PaymentRecorded       = "PaymentRecorded"
	PaymentPlanArchived   = "PaymentPlanArchived"
)

type Outbox struct {
	db *sql.DB
}

func NewOutbox(db *sql.DB) *Outbox {
	return &Outbox{db: db}
}

func (o *Outbox) Append(ctx context.Context, eventType, aggregateID string, payload map[string]any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = o.db.ExecContext(ctx, `
		INSERT INTO commercial.outbox_events (id, event_type, aggregate_id, payload, occurred_at)
		VALUES ($1, $2, $3, $4, $5)
	`, "evt_"+randomHex(16), eventType, aggregateID, payloadBytes, time.Now().UTC())
	return err
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
