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
	LeadDisqualified       = "LeadDisqualified"
	LeadConverted          = "LeadConverted"
	LeadArchived           = "LeadArchived"
	DuplicateWarningRaised = "DuplicateWarningRaised"
)

type Outbox struct {
	q sqlExecer
}

type sqlExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type DispatchConfig struct {
	ServiceID              string
	ServiceTokenSecret     []byte
	AuditHistoryServiceURL string
	ReportingServiceURL    string
	HTTPClient             *http.Client
	BatchSize              int
}

type outboxEvent struct {
	ID          string
	EventType   string
	AggregateID string
	Payload     map[string]any
}

func NewOutbox(db *sql.DB) *Outbox {
	return &Outbox{q: db}
}

func NewOutboxTx(tx *sql.Tx) *Outbox {
	return &Outbox{q: tx}
}

func (o *Outbox) Append(ctx context.Context, eventType, aggregateID string, payload map[string]any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = o.q.ExecContext(ctx, `
		INSERT INTO lead.outbox_events (id, event_type, aggregate_id, payload, occurred_at)
		VALUES ($1, $2, $3, $4, $5)
	`, "evt_"+randomHex(16), eventType, aggregateID, payloadBytes, time.Now().UTC())
	return err
}

func (o *Outbox) DispatchOnce(ctx context.Context, config DispatchConfig) error {
	if strings.TrimSpace(config.AuditHistoryServiceURL) == "" || len(config.ServiceTokenSecret) == 0 {
		return errors.New("audit dispatcher is not configured")
	}
	serviceID := config.ServiceID
	if serviceID == "" {
		serviceID = "lead"
	}
	batchSize := config.BatchSize
	if batchSize <= 0 {
		batchSize = 50
	}
	rows, err := o.q.QueryContext(ctx, `
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
		if err := deliverAuditEvent(ctx, config, serviceID, item); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if strings.TrimSpace(config.ReportingServiceURL) != "" {
			if err := deliverReportingProjection(ctx, config, serviceID, item); err != nil {
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
		}
		if _, err := o.q.ExecContext(ctx, `UPDATE lead.outbox_events SET published_at = now() WHERE id = $1 AND published_at IS NULL`, item.ID); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func deliverAuditEvent(ctx context.Context, config DispatchConfig, serviceID string, item outboxEvent) error {
	body, err := json.Marshal(auditAppendBody(item))
	if err != nil {
		return err
	}
	token, err := signServiceToken(serviceID, "audit-history", "audit.append", config.ServiceTokenSecret)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(config.AuditHistoryServiceURL, "/")+"/internal/events/append", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", serviceID)
	req.Header.Set("X-Intent", "audit.append")
	req.Header.Set("X-Correlation-Id", payloadString(item.Payload, "correlationId", item.ID))
	req.Header.Set("X-Actor-User-Id", payloadString(item.Payload, "actorId", "system"))
	req.Header.Set("X-Actor-Role", payloadString(item.Payload, "actorRole", "System"))
	req.Header.Set("X-Actor-Display", payloadString(item.Payload, "actorDisplay", req.Header.Get("X-Actor-User-Id")))
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
		return fmt.Errorf("audit append failed: status %d", resp.StatusCode)
	}
	return nil
}

func deliverReportingProjection(ctx context.Context, config DispatchConfig, serviceID string, item outboxEvent) error {
	body, err := json.Marshal(reportingProjectionBody(serviceID, item))
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

func reportingProjectionBody(serviceID string, item outboxEvent) map[string]any {
	return map[string]any{
		"sourceService": serviceID,
		"recordType":    "lead",
		"recordId":      item.AggregateID,
		"ownerId":       payloadString(item.Payload, "ownerId", payloadString(item.Payload, "newOwnerId", payloadString(item.Payload, "actorId", "system"))),
		"teamId":        payloadString(item.Payload, "teamId", "single-team"),
		"status":        leadProjectionStatus(item),
		"amount":        "0.00",
	}
}

func auditAppendBody(item outboxEvent) map[string]any {
	eventID, action := auditEventContract(item.EventType)
	return map[string]any{
		"eventUid":           item.ID,
		"eventId":            eventID,
		"eventVersion":       1,
		"surfaces":           []string{"record_history", "operation_log"},
		"action":             action,
		"resourceType":       "Lead",
		"resourceId":         item.AggregateID,
		"result":             "success",
		"afterSummary":       item.Payload,
		"diffClassification": "Restricted",
		"scopeSummary":       "record permission",
		"safeSummary":        action,
		"correlationId":      payloadString(item.Payload, "correlationId", item.ID),
		"causationId":        item.ID,
		"acceptanceIds":      []string{"ACC-014", "ACC-022"},
	}
}

func auditEventContract(eventType string) (string, string) {
	switch eventType {
	case LeadOwnerChanged:
		return "EVT-OWNER-CHANGED", "Owner changed"
	case LeadQualified:
		return "EVT-LEAD-QUALIFIED", "Lead qualified"
	case LeadDisqualified:
		return "EVT-LEAD-DISQUALIFIED", "Lead disqualified"
	case LeadConverted:
		return "EVT-LEAD-CONVERTED", "Lead converted"
	case LeadArchived:
		return "EVT-RECORD-ARCHIVED", "Record archived"
	case LeadCreated:
		return "EVT-LEAD-CREATED", "Lead created"
	case DuplicateWarningRaised:
		return "EVT-DUPLICATE-WARNING-RAISED", "Duplicate warning raised"
	default:
		return "EVT-" + strings.ToUpper(strings.ReplaceAll(eventType, "_", "-")), eventType
	}
}

func leadProjectionStatus(item outboxEvent) string {
	if item.EventType == LeadQualified {
		return "Valid"
	}
	if item.EventType == LeadDisqualified {
		return "Invalid"
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
