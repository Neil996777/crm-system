package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"crm-system/services/lead/internal/authz"
)

var ErrDownstreamUnavailable = errors.New("downstream unavailable")

type Actor struct {
	ID          string
	Role        string
	DisplayName string
}

type Config struct {
	AccountServiceURL      string
	OpportunityServiceURL  string
	AuditHistoryServiceURL string
	ServiceID              string
	ServiceTokenSecret     []byte
	HTTPClient             *http.Client
}

type ConversionClient struct {
	config Config
}

type AccountInput struct {
	CompanyName    string `json:"companyName"`
	CustomerStatus string `json:"customerStatus"`
	OwnerID        string `json:"ownerId"`
}

type OpportunityInput struct {
	CustomerID        string `json:"customerId,omitempty"`
	OwnerID           string `json:"ownerId"`
	Stage             string `json:"stage"`
	ExpectedAmount    string `json:"expectedAmount"`
	ExpectedCloseDate string `json:"expectedCloseDate"`
	Title             string `json:"title"`
}

type AccountResult struct {
	ID string `json:"id"`
}

type OpportunityResult struct {
	ID string `json:"id"`
}

func NewConversionClient(config Config) *ConversionClient {
	if config.ServiceID == "" {
		config.ServiceID = "lead"
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: 5 * time.Second}
	}
	return &ConversionClient{config: config}
}

func (c *ConversionClient) CreateAccount(ctx context.Context, actor Actor, input AccountInput) (AccountResult, error) {
	var result AccountResult
	if err := c.postInternal(ctx, actor, c.config.AccountServiceURL, "account", "account.create_for_lead_conversion", "/internal/accounts", input, &result); err != nil {
		return AccountResult{}, err
	}
	if result.ID == "" {
		return AccountResult{}, ErrDownstreamUnavailable
	}
	return result, nil
}

func (c *ConversionClient) CreateOpportunity(ctx context.Context, actor Actor, input OpportunityInput) (OpportunityResult, error) {
	var result OpportunityResult
	if err := c.postInternal(ctx, actor, c.config.OpportunityServiceURL, "opportunity", "opportunity.create_for_lead_conversion", "/internal/opportunities", input, &result); err != nil {
		return OpportunityResult{}, err
	}
	if result.ID == "" {
		return OpportunityResult{}, ErrDownstreamUnavailable
	}
	return result, nil
}

func (c *ConversionClient) postInternal(ctx context.Context, actor Actor, baseURL, audience, intent, path string, input any, output any) error {
	if strings.TrimSpace(baseURL) == "" {
		return ErrDownstreamUnavailable
	}
	body, err := json.Marshal(input)
	if err != nil {
		return err
	}
	token, err := authz.SignServiceToken(authz.ServiceTokenClaims{
		Issuer:   c.config.ServiceID,
		Audience: audience,
		Intent:   intent,
		Expires:  time.Now().UTC().Add(2 * time.Minute),
	}, c.config.ServiceTokenSecret)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(baseURL, "/")+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", c.config.ServiceID)
	req.Header.Set("X-Intent", intent)
	req.Header.Set("X-Actor-User-Id", actor.ID)
	req.Header.Set("X-Actor-Role", actor.Role)
	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return ErrDownstreamUnavailable
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w: status %d", ErrDownstreamUnavailable, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(output)
}
