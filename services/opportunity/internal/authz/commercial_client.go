package authz

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

const IntentContractSignedStatus = "commercial.contract_signed_status"

var ErrCommercialUnavailable = errors.New("commercial unavailable")

type CommercialClient struct {
	BaseURL            string
	ServiceID          string
	ServiceTokenSecret []byte
	HTTPClient         *http.Client
}

type ContractSignedStatus struct {
	ContractID    string
	OpportunityID string
	Status        string
	Signed        bool
}

func (c CommercialClient) ContractSignedStatus(ctx context.Context, contractID string) (ContractSignedStatus, error) {
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(contractID) == "" || len(c.ServiceTokenSecret) == 0 {
		return ContractSignedStatus{}, ErrCommercialUnavailable
	}
	serviceID := c.ServiceID
	if serviceID == "" {
		serviceID = "opportunity"
	}
	target := strings.TrimRight(c.BaseURL, "/") + "/internal/contracts/" + contractID + "/signed-status"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return ContractSignedStatus{}, err
	}
	req.Header.Set("Authorization", "Bearer "+createServiceToken(serviceID, "commercial", IntentContractSignedStatus, c.ServiceTokenSecret, time.Now().UTC()))
	req.Header.Set("X-Service-Id", serviceID)
	req.Header.Set("X-Intent", IntentContractSignedStatus)
	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return ContractSignedStatus{}, ErrCommercialUnavailable
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ContractSignedStatus{}, ErrCommercialUnavailable
	}
	var status ContractSignedStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return ContractSignedStatus{}, err
	}
	return status, nil
}

func createServiceToken(issuer, audience, intent string, secret []byte, now time.Time) string {
	payload, err := json.Marshal(ServiceTokenClaims{
		Issuer:   issuer,
		Audience: audience,
		Intent:   intent,
		Expires:  now.Add(5 * time.Minute),
	})
	if err != nil {
		panic(err)
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	return encodedPayload + "." + sign(encodedPayload, secret)
}
