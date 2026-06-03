package event

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	LeadCreated            = "LeadCreated"
	LeadOwnerChanged       = "LeadOwnerChanged"
	LeadQualified          = "LeadQualified"
	LeadConverted          = "LeadConverted"
	LeadArchived           = "LeadArchived"
	DuplicateWarningRaised = "DuplicateWarningRaised"
)

type Outbox struct {
	db *sql.DB
}

type DispatchConfig struct {
	ServiceID           string
	ServiceTokenSecret  []byte
	ReportingServiceURL string
	HTTPClient          *http.Client
	BatchSize           int
}

type outboxEvent struct {
	ID          string
	EventType   string
	AggregateID string
	Payload     map[string]any
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
		INSERT INTO lead.outbox_events (id, event_type, aggregate_id, payload, occurred_at)
		VALUES ($1, $2, $3, $4, $5)
	`, "evt_"+randomHex(16), eventType, aggregateID, payloadBytes, time.Now().UTC())
	return err
}

func (o *Outbox) DispatchOnce(ctx context.Context, config DispatchConfig) error {
	if strings.TrimSpace(config.ReportingServiceURL) == "" || len(config.ServiceTokenSecret) == 0 {
		return errors.New("reporting dispatcher is not configured")
	}
	serviceID := config.ServiceID
	if serviceID == "" {
		serviceID = "lead"
	}
	batchSize := config.BatchSize
	if batchSize <= 0 {
		batchSize = 50
	}
	rows, err := o.db.QueryContext(ctx, `
		SELECT id, event_type, aggregate_id, payload
		FROM lead.outbox_events
		WHERE published_at IS NULL
		ORDER BY occurred_at ASC, id ASC
		LIMIT $1
	`, batchSize)
	if err != nil {
		return err
	}
	defer rows.Close()
	var events []outboxEvent
	for rows.Next() {
		var item outboxEvent
		var payloadBytes []byte
		if err := rows.Scan(&item.ID, &item.EventType, &item.AggregateID, &payloadBytes); err != nil {
			return err
		}
		if err := json.Unmarshal(payloadBytes, &item.Payload); err != nil {
			return err
		}
		events = append(events, item)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	var firstErr error
	for _, item := range events {
		if err := deliverReportingProjection(ctx, config, serviceID, item); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if _, err := o.db.ExecContext(ctx, `UPDATE lead.outbox_events SET published_at = now() WHERE id = $1 AND published_at IS NULL`, item.ID); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func deliverReportingProjection(ctx context.Context, config DispatchConfig, serviceID string, item outboxEvent) error {
	body, err := json.Marshal(map[string]any{
		"sourceService": serviceID,
		"recordType":    "lead",
		"recordId":      item.AggregateID,
		"ownerId":       payloadString(item.Payload, "ownerId", payloadString(item.Payload, "newOwnerId", payloadString(item.Payload, "actorId", "system"))),
		"teamId":        payloadString(item.Payload, "teamId", "single-team"),
		"status":        leadProjectionStatus(item),
		"amount":        "0.00",
	})
	if err != nil {
		return err
	}
	token, err := signServiceToken(serviceID, "reporting", "reporting.projection_ingest", config.ServiceTokenSecret)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(config.ReportingServiceURL, "/")+"/internal/projections", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", serviceID)
	req.Header.Set("X-Intent", "reporting.projection_ingest")
	req.Header.Set("X-Correlation-Id", payloadString(item.Payload, "correlationId", item.ID))
	client := config.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("reporting projection failed: status %d", resp.StatusCode)
	}
	return nil
}

func leadProjectionStatus(item outboxEvent) string {
	if item.EventType == LeadQualified {
		return "Valid"
	}
	if item.EventType == LeadConverted {
		return "Converted"
	}
	if item.EventType == LeadArchived {
		return "Archived"
	}
	return payloadString(item.Payload, "status", "Pending Qualification")
}

func payloadString(payload map[string]any, key, fallback string) string {
	if value, ok := payload[key].(string); ok && strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func signServiceToken(issuer, audience, intent string, secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("missing service token secret")
	}
	payload, err := json.Marshal(map[string]any{"iss": issuer, "aud": audience, "intent": intent, "exp": time.Now().UTC().Add(2 * time.Minute)})
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(encodedPayload))
	return encodedPayload + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
