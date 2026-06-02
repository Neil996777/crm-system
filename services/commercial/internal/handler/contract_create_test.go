package handler

import (
	"net/http"
	"testing"
)

func TestContractCreateAcceptance(t *testing.T) {
	db := newCommercialTestDB(t)
	app := NewCommercialServer(db, Config{})

	t.Run("TEST-CONTRACT-CREATE-001 creates Pending Signature from Accepted quote without signed date", func(t *testing.T) {
		quoteID := createAcceptedQuote(t, app, "opp_contract_create_001", "sales-1")
		rec := postCommercialJSON(app, "/contracts", contractCreateBody(quoteID, "opp_contract_create_001", "acct_opp_contract_create_001", "10000.00"), actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "quoteId", quoteID)
		requireJSONValue(t, body, "opportunityId", "opp_contract_create_001")
		requireJSONValue(t, body, "customerId", "acct_opp_contract_create_001")
		requireJSONValue(t, body, "status", "Pending Signature")
		requireJSONValue(t, body, "amount", "10000.00")
		requireJSONValue(t, body, "contractNote", "TASK-018 contract note")
		requireJSONValue(t, body, "expectedSignedDate", "2027-01-15")
		if _, ok := body["signedEffectiveDate"]; ok {
			t.Fatalf("pending signature contract must not require or expose signedEffectiveDate on create: %#v", body)
		}
		requireEvent(t, db, "ContractCreated", body["id"].(string))
	})

	t.Run("TEST-CONTRACT-CREATE-002 rejects missing note link amount or expected signed date", func(t *testing.T) {
		quoteID := createAcceptedQuote(t, app, "opp_contract_create_002", "sales-1")
		for name, body := range map[string]map[string]any{
			"note": {
				"quoteId":            quoteID,
				"opportunityId":      "opp_contract_create_002",
				"customerId":         "acct_opp_contract_create_002",
				"amount":             "10000.00",
				"status":             "Pending Signature",
				"expectedSignedDate": "2027-01-15",
				"ownerId":            "sales-1",
			},
			"link": {
				"opportunityId":      "opp_contract_create_002",
				"customerId":         "acct_opp_contract_create_002",
				"amount":             "10000.00",
				"status":             "Pending Signature",
				"contractNote":       "TASK-018 contract note",
				"expectedSignedDate": "2027-01-15",
				"ownerId":            "sales-1",
			},
			"amount": {
				"quoteId":            quoteID,
				"opportunityId":      "opp_contract_create_002",
				"customerId":         "acct_opp_contract_create_002",
				"status":             "Pending Signature",
				"contractNote":       "TASK-018 contract note",
				"expectedSignedDate": "2027-01-15",
				"ownerId":            "sales-1",
			},
			"expectedSignedDate": {
				"quoteId":       quoteID,
				"opportunityId": "opp_contract_create_002",
				"customerId":    "acct_opp_contract_create_002",
				"amount":        "10000.00",
				"status":        "Pending Signature",
				"contractNote":  "TASK-018 contract note",
				"ownerId":       "sales-1",
			},
		} {
			t.Run(name, func(t *testing.T) {
				rec := postCommercialJSON(app, "/contracts", body, actorHeaders("sales-1", "Sales"))
				if rec.Code != http.StatusBadRequest {
					t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
				}
				requireErrorCode(t, rec, "VALIDATION_FAILED")
			})
		}
	})

	t.Run("TEST-CONTRACT-CREATE-003 and TEST-INV-CONTRACTQUOTE-001 reject expired or non-Accepted quote link", func(t *testing.T) {
		draftQuoteID := createQuoteForLifecycle(t, app, "opp_contract_draft_quote", "sales-1")
		draft := postCommercialJSON(app, "/contracts", contractCreateBody(draftQuoteID, "opp_contract_draft_quote", "acct_opp_contract_draft_quote", "10000.00"), actorHeaders("sales-1", "Sales"))
		if draft.Code != http.StatusBadRequest {
			t.Fatalf("expected draft quote rejected 400, got %d body=%s", draft.Code, draft.Body.String())
		}
		requireErrorCode(t, draft, "CONTRACT_QUOTE_INVALID")

		expiredQuoteID := createQuoteForLifecycle(t, app, "opp_contract_expired_quote", "sales-1")
		expiredStatus := postCommercialJSON(app, "/quotes/"+expiredQuoteID+"/status", map[string]any{
			"expectedVersion": 1,
			"toStatus":        "Expired",
		}, actorHeaders("sales-1", "Sales"))
		if expiredStatus.Code != http.StatusOK {
			t.Fatalf("expected expire 200, got %d body=%s", expiredStatus.Code, expiredStatus.Body.String())
		}
		expired := postCommercialJSON(app, "/contracts", contractCreateBody(expiredQuoteID, "opp_contract_expired_quote", "acct_opp_contract_expired_quote", "10000.00"), actorHeaders("sales-1", "Sales"))
		if expired.Code != http.StatusBadRequest {
			t.Fatalf("expected expired quote rejected 400, got %d body=%s", expired.Code, expired.Body.String())
		}
		requireErrorCode(t, expired, "CONTRACT_QUOTE_INVALID")
	})

	t.Run("TEST-CONTRACT-AMOUNT-DIFF-001 requires and persists amount difference reason", func(t *testing.T) {
		quoteID := createAcceptedQuote(t, app, "opp_contract_amount_diff", "sales-1")
		withoutReason := contractCreateBody(quoteID, "opp_contract_amount_diff", "acct_opp_contract_amount_diff", "12000.00")
		rec := postCommercialJSON(app, "/contracts", withoutReason, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected diff without reason 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "AMOUNT_DIFFERENCE_REASON_REQUIRED")

		withReason := contractCreateBody(quoteID, "opp_contract_amount_diff", "acct_opp_contract_amount_diff", "12000.00")
		withReason["amountDifferenceReason"] = "Scope expanded after quote acceptance"
		created := postCommercialJSON(app, "/contracts", withReason, actorHeaders("sales-1", "Sales"))
		if created.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", created.Code, created.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, created), "amountDifferenceReason", "Scope expanded after quote acceptance")
	})

	t.Run("TEST-INV-NOAPPROVAL-001 creates without approval e-sign or template fields", func(t *testing.T) {
		quoteID := createAcceptedQuote(t, app, "opp_contract_noapproval", "sales-1")
		rec := postCommercialJSON(app, "/contracts", contractCreateBody(quoteID, "opp_contract_noapproval", "acct_opp_contract_noapproval", "10000.00"), actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		for _, forbidden := range []string{"approvalStatus", "approvalRequired", "eSignatureStatus", "templateId"} {
			if _, ok := body[forbidden]; ok {
				t.Fatalf("contract create must not expose approval/e-sign/template field %s: %#v", forbidden, body)
			}
		}
	})
}

func createAcceptedQuote(t *testing.T, app http.Handler, opportunityID, ownerID string) string {
	t.Helper()
	quoteID := createQuoteForLifecycle(t, app, opportunityID, ownerID)
	sent := postCommercialJSON(app, "/quotes/"+quoteID+"/status", map[string]any{
		"expectedVersion": 1,
		"toStatus":        "Sent",
	}, actorHeaders(ownerID, "Sales"))
	if sent.Code != http.StatusOK {
		t.Fatalf("expected send 200, got %d body=%s", sent.Code, sent.Body.String())
	}
	accepted := postCommercialJSON(app, "/quotes/"+quoteID+"/status", map[string]any{
		"expectedVersion": 2,
		"toStatus":        "Accepted",
	}, actorHeaders(ownerID, "Sales"))
	if accepted.Code != http.StatusOK {
		t.Fatalf("expected accept 200, got %d body=%s", accepted.Code, accepted.Body.String())
	}
	return quoteID
}

func contractCreateBody(quoteID, opportunityID, customerID, amount string) map[string]any {
	return map[string]any{
		"quoteId":            quoteID,
		"opportunityId":      opportunityID,
		"customerId":         customerID,
		"amount":             amount,
		"status":             "Pending Signature",
		"contractNote":       "TASK-018 contract note",
		"expectedSignedDate": "2027-01-15",
		"ownerId":            "sales-1",
	}
}
