package client

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

type ImportRunLogInput struct {
	RunID        string
	ActorID      string
	ActorRole    string
	ObjectType   string
	TotalRows    int
	SuccessCount int
	FailureCount int
	Result       string
}

func NewAuditClient(baseURL, serviceID string, secret []byte, client *http.Client) AuditClient {
	if serviceID == "" {
		serviceID = "import-export"
	}
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return AuditClient{baseURL: strings.TrimRight(baseURL, "/"), serviceID: serviceID, secret: secret, client: client}
}

func (c AuditClient) BaseURLConfigured() bool {
	return c.baseURL != ""
}

func (c AuditClient) AppendImportRun(ctx context.Context, input ImportRunLogInput) error {
	if c.baseURL == "" {
		return nil
	}
	body, err := json.Marshal(map[string]any{
		"eventId":            "EVT-IMPORT-RUN",
		"eventVersion":       1,
		"surfaces":           []string{"operation_log"},
		"action":             "csv_import",
		"resourceType":       "import_run",
		"resourceId":         input.RunID,
		"result":             input.Result,
		"beforeSummary":      map[string]any{},
		"afterSummary":       map[string]any{"objectType": input.ObjectType, "totalRows": input.TotalRows, "successCount": input.SuccessCount, "failureCount": input.FailureCount},
		"diffClassification": "Restricted",
		"scopeSummary":       "import scope",
		"safeSummary":        fmt.Sprintf("CSV import completed for %s with %d successful and %d failed rows.", input.ObjectType, input.SuccessCount, input.FailureCount),
		"correlationId":      input.RunID,
		"acceptanceIds":      []string{"ACC-020", "ACC-022", "TEST-CSV-IMPORT", "TEST-OPLOG-001"},
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/internal/events/append", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", c.serviceID)
	req.Header.Set("X-Intent", "audit.append")
	req.Header.Set("X-Actor-User-Id", input.ActorID)
	req.Header.Set("X-Actor-Role", input.ActorRole)
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
