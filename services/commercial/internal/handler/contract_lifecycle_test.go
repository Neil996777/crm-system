package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestContractLifecycleAcceptance(t *testing.T) {
	db := newCommercialTestDB(t)
	app := NewCommercialServer(db, Config{ServiceID: "commercial", ServiceTokenSecret: []byte("commercial_test_secret")})

	t.Run("TEST-CONTRACT-LIFECYCLE-001 signs activates and completes with signed effective date", func(t *testing.T) {
		contractID := createPendingContract(t, app, "opp_contract_lifecycle_001", "sales-1")
		missingDate := postCommercialJSON(app, "/contracts/"+contractID+"/status", map[string]any{
			"expectedVersion": 1,
			"toStatus":        "Signed",
		}, actorHeaders("sales-1", "Sales"))
		if missingDate.Code != http.StatusBadRequest {
			t.Fatalf("expected missing signed date 400, got %d body=%s", missingDate.Code, missingDate.Body.String())
		}
		requireErrorCode(t, missingDate, "SIGNED_EFFECTIVE_DATE_REQUIRED")

		signed := postCommercialJSON(app, "/contracts/"+contractID+"/status", map[string]any{
			"expectedVersion":     1,
			"toStatus":            "Signed",
			"signedEffectiveDate": "2027-01-20",
		}, actorHeaders("sales-1", "Sales"))
		if signed.Code != http.StatusOK {
			t.Fatalf("expected sign 200, got %d body=%s", signed.Code, signed.Body.String())
		}
		body := decodeJSON(t, signed)
		requireJSONValue(t, body, "status", "Signed")
		requireJSONValue(t, body, "signedEffectiveDate", "2027-01-20")
		requireEvent(t, db, "ContractStatusChanged", contractID)

		active := postCommercialJSON(app, "/contracts/"+contractID+"/status", map[string]any{
			"expectedVersion": 2,
			"toStatus":        "Active",
		}, actorHeaders("sales-1", "Sales"))
		if active.Code != http.StatusOK {
			t.Fatalf("expected activate 200, got %d body=%s", active.Code, active.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, active), "status", "Active")

		completed := postCommercialJSON(app, "/contracts/"+contractID+"/status", map[string]any{
			"expectedVersion": 3,
			"toStatus":        "Completed",
		}, actorHeaders("sales-1", "Sales"))
		if completed.Code != http.StatusOK {
			t.Fatalf("expected complete 200, got %d body=%s", completed.Code, completed.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, completed), "status", "Completed")
	})

	t.Run("TEST-CONTRACT-LIFECYCLE-002 and TEST-INV-CONTRACTDATE-001 reject signed states without signed date", func(t *testing.T) {
		contractID := createPendingContract(t, app, "opp_contract_lifecycle_002", "sales-1")
		for _, toStatus := range []string{"Signed", "Active", "Completed"} {
			rec := postCommercialJSON(app, "/contracts/"+contractID+"/status", map[string]any{
				"expectedVersion": 1,
				"toStatus":        toStatus,
			}, actorHeaders("sales-1", "Sales"))
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected %s without date 400, got %d body=%s", toStatus, rec.Code, rec.Body.String())
			}
			requireErrorCode(t, rec, "SIGNED_EFFECTIVE_DATE_REQUIRED")
		}
	})

	t.Run("TEST-CONTRACT-LIFECYCLE-003 terminates pre-signature without date and post-signature with persisted date", func(t *testing.T) {
		preSigID := createPendingContract(t, app, "opp_contract_presig_terminate", "sales-1")
		preSig := postCommercialJSON(app, "/contracts/"+preSigID+"/status", map[string]any{
			"expectedVersion": 1,
			"toStatus":        "Terminated",
		}, actorHeaders("sales-1", "Sales"))
		if preSig.Code != http.StatusOK {
			t.Fatalf("expected pre-signature terminate 200, got %d body=%s", preSig.Code, preSig.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, preSig), "status", "Terminated")

		postSigID := createPendingContract(t, app, "opp_contract_postsig_terminate", "sales-1")
		signed := postCommercialJSON(app, "/contracts/"+postSigID+"/status", map[string]any{
			"expectedVersion":     1,
			"toStatus":            "Signed",
			"signedEffectiveDate": "2027-02-01",
		}, actorHeaders("sales-1", "Sales"))
		if signed.Code != http.StatusOK {
			t.Fatalf("expected sign 200, got %d body=%s", signed.Code, signed.Body.String())
		}
		postSig := postCommercialJSON(app, "/contracts/"+postSigID+"/status", map[string]any{
			"expectedVersion": 2,
			"toStatus":        "Terminated",
		}, actorHeaders("sales-1", "Sales"))
		if postSig.Code != http.StatusOK {
			t.Fatalf("expected post-signature terminate 200, got %d body=%s", postSig.Code, postSig.Body.String())
		}
		postSigBody := decodeJSON(t, postSig)
		requireJSONValue(t, postSigBody, "status", "Terminated")
		requireJSONValue(t, postSigBody, "signedEffectiveDate", "2027-02-01")
	})

	t.Run("Signed status is queryable by opportunity-service with S2S token", func(t *testing.T) {
		contractID := createPendingContract(t, app, "opp_contract_signed_query", "sales-1")
		signed := postCommercialJSON(app, "/contracts/"+contractID+"/status", map[string]any{
			"expectedVersion":     1,
			"toStatus":            "Signed",
			"signedEffectiveDate": "2027-02-10",
		}, actorHeaders("sales-1", "Sales"))
		if signed.Code != http.StatusOK {
			t.Fatalf("expected sign 200, got %d body=%s", signed.Code, signed.Body.String())
		}

		query := signedStatusQuery(app, contractID, makeServiceToken(t, "opportunity", "commercial", "commercial.contract_signed_status", []byte("commercial_test_secret")))
		if query.Code != http.StatusOK {
			t.Fatalf("expected S2S signed query 200, got %d body=%s", query.Code, query.Body.String())
		}
		body := decodeJSON(t, query)
		requireJSONValue(t, body, "contractId", contractID)
		requireJSONValue(t, body, "status", "Signed")
		if body["signed"] != true {
			t.Fatalf("expected signed=true, got %#v", body)
		}

		denied := signedStatusQuery(app, contractID, "invalid")
		if denied.Code != http.StatusUnauthorized {
			t.Fatalf("expected invalid S2S token 401, got %d body=%s", denied.Code, denied.Body.String())
		}
		requireErrorCode(t, denied, "SERVICE_AUTH_FAILED")
	})
}

func createPendingContract(t *testing.T, app http.Handler, opportunityID, ownerID string) string {
	t.Helper()
	quoteID := createAcceptedQuote(t, app, opportunityID, ownerID)
	rec := postCommercialJSON(app, "/contracts", contractCreateBody(quoteID, opportunityID, "acct_"+opportunityID, "10000.00"), actorHeaders(ownerID, "Sales"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected contract create 201, got %d body=%s", rec.Code, rec.Body.String())
	}
	return decodeJSON(t, rec)["id"].(string)
}

func signedStatusQuery(handler http.Handler, contractID, token string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/internal/contracts/"+contractID+"/signed-status", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", "opportunity")
	req.Header.Set("X-Intent", "commercial.contract_signed_status")
	handler.ServeHTTP(rec, req)
	return rec
}

func makeServiceToken(t *testing.T, issuer, audience, intent string, secret []byte) string {
	t.Helper()
	payload, err := json.Marshal(map[string]any{
		"iss":    issuer,
		"aud":    audience,
		"intent": intent,
		"exp":    time.Now().UTC().Add(5 * time.Minute),
	})
	if err != nil {
		t.Fatalf("marshal token payload: %v", err)
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(encodedPayload))
	return encodedPayload + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
