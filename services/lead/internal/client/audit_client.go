package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"crm-system/services/lead/internal/authz"
)

type AuditClient struct {
	config Config
}

type AuditEventInput struct {
	EventID            string
	Action             string
	ResourceType       string
	ResourceID         string
	Result             string
	BeforeSummary      map[string]any
	AfterSummary       map[string]any
	DiffClassification string
	SafeSummary        string
	CorrelationID      string
	AcceptanceIDs      []string
}

func NewAuditClient(config Config) *AuditClient {
	if config.ServiceID == "" {
		config.ServiceID = "lead"
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: 5 * time.Second}
	}
	return &AuditClient{config: config}
}

func (c *AuditClient) AppendRecordHistory(ctx context.Context, actor Actor, input AuditEventInput) error {
	if strings.TrimSpace(c.config.AuditHistoryServiceURL) == "" {
		return ErrDownstreamUnavailable
	}
	actorDisplay := actor.DisplayName
	if actorDisplay == "" {
		actorDisplay = actor.ID
	}
	body, err := json.Marshal(map[string]any{
		"eventId":            input.EventID,
		"eventVersion":       1,
		"surfaces":           []string{"record_history"},
		"action":             input.Action,
		"resourceType":       input.ResourceType,
		"resourceId":         input.ResourceID,
		"result":             input.Result,
		"beforeSummary":      input.BeforeSummary,
		"afterSummary":       input.AfterSummary,
		"diffClassification": input.DiffClassification,
		"scopeSummary":       "record permission",
		"safeSummary":        input.SafeSummary,
		"correlationId":      input.CorrelationID,
		"acceptanceIds":      input.AcceptanceIDs,
	})
	if err != nil {
		return err
	}
	token, err := authz.SignServiceToken(authz.ServiceTokenClaims{
		Issuer:   c.config.ServiceID,
		Audience: "audit-history",
		Intent:   "audit.append",
		Expires:  time.Now().UTC().Add(2 * time.Minute),
	}, c.config.ServiceTokenSecret)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.config.AuditHistoryServiceURL, "/")+"/internal/events/append", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", c.config.ServiceID)
	req.Header.Set("X-Intent", "audit.append")
	req.Header.Set("X-Actor-User-Id", actor.ID)
	req.Header.Set("X-Actor-Role", actor.Role)
	req.Header.Set("X-Actor-Display", actorDisplay)
	if strings.TrimSpace(input.CorrelationID) != "" {
		req.Header.Set("X-Correlation-Id", input.CorrelationID)
	}
	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return ErrDownstreamUnavailable
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w: audit append status %d", ErrDownstreamUnavailable, resp.StatusCode)
	}
	return nil
}
