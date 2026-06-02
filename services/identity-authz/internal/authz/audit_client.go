package authz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AuditClient struct {
	baseURL   string
	serviceID string
	secret    []byte
	client    *http.Client
}

type OperationLogInput struct {
	ActorID       string
	ActorRole     string
	ActorDisplay  string
	Action        string
	ResourceType  string
	ResourceID    string
	Result        string
	BeforeSummary map[string]any
	AfterSummary  map[string]any
	SafeSummary   string
	CorrelationID string
}

func NewAuditClient(baseURL, serviceID string, secret []byte, client *http.Client) AuditClient {
	if serviceID == "" {
		serviceID = "identity-authz"
	}
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return AuditClient{baseURL: baseURL, serviceID: serviceID, secret: secret, client: client}
}

func (c AuditClient) AppendOperationLog(ctx context.Context, input OperationLogInput) error {
	if strings.TrimSpace(c.baseURL) == "" {
		return nil
	}
	actorDisplay := input.ActorDisplay
	if actorDisplay == "" {
		actorDisplay = input.ActorID
	}
	body, err := json.Marshal(map[string]any{
		"eventId":            "EVT-USER-ADMIN-CHANGED",
		"eventVersion":       1,
		"surfaces":           []string{"operation_log"},
		"action":             input.Action,
		"resourceType":       input.ResourceType,
		"resourceId":         input.ResourceID,
		"result":             input.Result,
		"beforeSummary":      input.BeforeSummary,
		"afterSummary":       input.AfterSummary,
		"diffClassification": "Security Critical",
		"scopeSummary":       "administrator only",
		"safeSummary":        input.SafeSummary,
		"correlationId":      input.CorrelationID,
		"acceptanceIds":      []string{"ACC-022", "TEST-OPLOG-001", "TEST-OPLOG-002", "TEST-OPLOG-005"},
	})
	if err != nil {
		return err
	}
	token, err := SignServiceToken(ServiceTokenClaims{
		Issuer:   c.serviceID,
		Audience: "audit-history",
		Intent:   "audit.append",
		Expires:  time.Now().UTC().Add(2 * time.Minute),
	}, c.secret)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.baseURL, "/")+"/internal/events/append", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", c.serviceID)
	req.Header.Set("X-Intent", "audit.append")
	req.Header.Set("X-Actor-User-Id", input.ActorID)
	req.Header.Set("X-Actor-Role", input.ActorRole)
	req.Header.Set("X-Actor-Display", actorDisplay)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("audit append status %d", resp.StatusCode)
	}
	return nil
}
