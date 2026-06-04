package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpportunityCloseWonAcceptance(t *testing.T) {
	db := newOpportunityTestDB(t)
	contractOpportunity := map[string]string{}
	commercial := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/contracts/contract_signed/signed-status" && r.URL.Path != "/internal/contracts/contract_unsigned/signed-status" {
			t.Fatalf("unexpected commercial path %s", r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Fatalf("commercial request missing bearer service token")
		}
		if r.Header.Get("X-Service-Id") != "opportunity" || r.Header.Get("X-Intent") != "commercial.contract_signed_status" {
			t.Fatalf("commercial request missing S2S headers: service=%q intent=%q", r.Header.Get("X-Service-Id"), r.Header.Get("X-Intent"))
		}
		if r.Header.Get("X-Correlation-Id") != "corr-close-won" {
			t.Fatalf("commercial request missing correlation id: %q", r.Header.Get("X-Correlation-Id"))
		}
		contractID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/internal/contracts/"), "/signed-status")
		signed := r.URL.Path == "/internal/contracts/contract_signed/signed-status"
		writeTestJSON(t, w, http.StatusOK, map[string]any{
			"contractId":    contractID,
			"opportunityId": contractOpportunity[contractID],
			"status":        map[bool]string{true: "Signed", false: "Pending Signature"}[signed],
			"signed":        signed,
		})
	}))
	defer commercial.Close()
	app := NewOpportunityServer(db, Config{
		ServiceID:          "opportunity",
		ServiceTokenSecret: []byte("opportunity_test_secret"),
		CommercialBaseURL:  commercial.URL,
	})

	t.Run("TEST-OPP-CLOSE-002 rejects Won without Signed contract", func(t *testing.T) {
		opportunityID := createOpportunityAtContractNegotiation(t, app, "opp_close_won_unsigned", "acct_close_won_unsigned", "sales-1")
		contractOpportunity["contract_unsigned"] = opportunityID
		headers := actorHeaders("sales-1", "Sales")
		headers["X-Correlation-Id"] = "corr-close-won"
		rec := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/close-won", map[string]any{
			"expectedVersion": 4,
			"contractId":      "contract_unsigned",
			"closeDate":       "2027-03-15",
		}, headers)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected early Won 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "EARLY_WON_BLOCKED")
		requireNoEvent(t, db, "OpportunityClosedWon", opportunityID)
	})

	t.Run("TEST-OPP-CLOSE-001 and TEST-INV-WONAFTERPAY-001 closes Won with Signed contract and no payment gate", func(t *testing.T) {
		opportunityID := createOpportunityAtContractNegotiation(t, app, "opp_close_won_signed", "acct_close_won_signed", "sales-1")
		contractOpportunity["contract_signed"] = opportunityID
		headers := actorHeaders("sales-1", "Sales")
		headers["X-Correlation-Id"] = "corr-close-won"
		rec := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/close-won", map[string]any{
			"expectedVersion": 4,
			"contractId":      "contract_signed",
			"closeDate":       "2027-03-15",
		}, headers)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected close Won 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "opportunityId", opportunityID)
		requireJSONValue(t, body, "status", "Won")
		fetch := getOpportunityJSON(app, "/opportunities/"+opportunityID, actorHeaders("sales-1", "Sales"))
		fetchBody := decodeJSON(t, fetch)
		requireJSONValue(t, fetchBody, "stage", "Won")
		requireJSONValue(t, fetchBody, "wonContractId", "contract_signed")
		requireEvent(t, db, "OpportunityClosedWon", opportunityID)
	})
}

func createOpportunityAtContractNegotiation(t *testing.T, app http.Handler, opportunityIDSeed, customerID, ownerID string) string {
	t.Helper()
	opportunityID := createOpportunityForStage(t, app, customerID, ownerID)
	for version, stage := range []string{"Needs Confirmed", "Quote", "Contract Negotiation"} {
		rec := postOpportunityJSON(app, "/opportunities/"+opportunityID+"/stage", map[string]any{
			"expectedVersion": version + 1,
			"toStage":         stage,
		}, actorHeaders(ownerID, "Sales"))
		if rec.Code != http.StatusOK {
			t.Fatalf("advance %s to %s: got %d body=%s", opportunityIDSeed, stage, rec.Code, rec.Body.String())
		}
	}
	return opportunityID
}

func writeTestJSON(t *testing.T, w http.ResponseWriter, status int, body any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Fatalf("write test json: %v", err)
	}
}
