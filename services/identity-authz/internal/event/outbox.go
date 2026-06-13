package event

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"crm-system/services/identity-authz/internal/authz"
)

const (
	UserSignedIn             = "UserSignedIn"
	UserSignedOut            = "UserSignedOut"
	UserAccessDenied         = "UserAccessDenied"
	UserRoleStatusChanged    = "UserRoleStatusChanged"
	LastAdministratorBlocked = "LastAdministratorBlocked"
)

type Outbox struct {
	q sqlExecer
}

type sqlExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func NewOutbox(db *sql.DB) *Outbox {
	return &Outbox{q: db}
}

func NewOutboxTx(tx *sql.Tx) *Outbox {
	return &Outbox{q: tx}
}

type DispatchConfig struct {
	ServiceID              string
	ServiceTokenSecret     []byte
	AuditHistoryServiceURL string
	HTTPClient             *http.Client
	BatchSize              int
}

type outboxEvent struct {
	ID          string
	EventType   string
	AggregateID string
	Payload     map[string]any
}

func (o *Outbox) Append(ctx context.Context, eventType, aggregateID string, payload map[string]any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = o.q.ExecContext(ctx, `
		INSERT INTO identity_authz.outbox_events (id, event_type, aggregate_type, aggregate_id, payload, occurred_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, "evt_"+randomHex(16), eventType, "User", aggregateID, payloadBytes, time.Now().UTC())
	return err
}

func (o *Outbox) DispatchOnce(ctx context.Context, config DispatchConfig) error {
	if strings.TrimSpace(config.AuditHistoryServiceURL) == "" || len(config.ServiceTokenSecret) == 0 {
		return errors.New("audit dispatcher is not configured")
	}
	serviceID := config.ServiceID
	if serviceID == "" {
		serviceID = "identity-authz"
	}
	batchSize := config.BatchSize
	if batchSize <= 0 {
		batchSize = 50
	}
	rows, err := o.q.QueryContext(ctx, `
		SELECT id, event_type, aggregate_id, payload
		FROM identity_authz.outbox_events
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
		if _, err := o.q.ExecContext(ctx, `
			UPDATE identity_authz.outbox_events
			SET published_at = now()
			WHERE id = $1 AND published_at IS NULL
		`, item.ID); err != nil && firstErr == nil {
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
	token, err := authz.SignServiceToken(authz.ServiceTokenClaims{
		Issuer:   serviceID,
		Audience: "audit-history",
		Intent:   "audit.append",
		Expires:  time.Now().UTC().Add(2 * time.Minute),
	}, config.ServiceTokenSecret)
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

func auditAppendBody(item outboxEvent) map[string]any {
	eventID, action := auditEventContract(item.EventType, item.Payload)
	beforeSummary := map[string]any{"value": payloadString(item.Payload, "before", "")}
	if summary, ok := payloadMap(item.Payload, "beforeSummary"); ok {
		beforeSummary = summary
	}
	afterSummary := item.Payload
	if summary, ok := payloadMap(item.Payload, "afterSummary"); ok {
		afterSummary = summary
	}
	resourceID := payloadString(item.Payload, "resourceId", item.AggregateID)
	if strings.TrimSpace(resourceID) == "" {
		resourceID = "unknown"
	}
	return map[string]any{
		"eventUid":           item.ID,
		"eventId":            eventID,
		"eventVersion":       1,
		"surfaces":           []string{"operation_log"},
		"action":             action,
		"resourceType":       payloadString(item.Payload, "resourceType", "User"),
		"resourceId":         resourceID,
		"result":             payloadString(item.Payload, "result", "success"),
		"reasonCode":         payloadString(item.Payload, "reasonCode", payloadString(item.Payload, "reason", "")),
		"beforeSummary":      beforeSummary,
		"afterSummary":       afterSummary,
		"diffClassification": payloadString(item.Payload, "diffClassification", "Security Critical"),
		"scopeSummary":       payloadString(item.Payload, "scopeSummary", "administrator only"),
		"safeSummary":        payloadString(item.Payload, "safeSummary", safeSummary(action)),
		"correlationId":      payloadString(item.Payload, "correlationId", item.ID),
		"causationId":        item.ID,
		"acceptanceIds":      []string{"ACC-022", "TEST-OPLOG-001", "TEST-OPLOG-002", "TEST-OPLOG-005"},
	}
}

func auditEventContract(eventType string, payload map[string]any) (string, string) {
	switch eventType {
	case UserRoleStatusChanged:
		action := payloadString(payload, "action", "user_admin_changed")
		switch action {
		case "change_role":
			return "EVT-USER-ROLE-CHANGED", action
		case "change_status":
			return "EVT-USER-STATUS-CHANGED", action
		default:
			return "EVT-USER-STATUS-CHANGED", action
		}
	case LastAdministratorBlocked:
		return "EVT-LAST-ADMIN-BLOCKED", "last_admin_blocked"
	case UserSignedIn:
		return "EVT-AUTH-LOGIN-SUCCEEDED", "sign_in"
	case UserSignedOut:
		return "EVT-USER-SIGNED-OUT", "sign_out"
	case UserAccessDenied:
		if payloadString(payload, "reasonCode", payloadString(payload, "reason", "")) == "login_failed" {
			return "EVT-AUTH-LOGIN-FAILED", "login_failed"
		}
		return "EVT-AUTH-ACCESS-DENIED", "access_denied"
	default:
		return "EVT-" + strings.ToUpper(strings.ReplaceAll(eventType, "_", "-")), eventType
	}
}

func payloadString(payload map[string]any, key, fallback string) string {
	if value, ok := payload[key].(string); ok && strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func payloadMap(payload map[string]any, key string) (map[string]any, bool) {
	if value, ok := payload[key].(map[string]any); ok {
		return value, true
	}
	return nil, false
}

func safeSummary(action string) string {
	if action == "" {
		return "Identity authorization event"
	}
	return "Identity authorization " + action
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
