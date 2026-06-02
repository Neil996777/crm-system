package handler

import (
	"net/http"
	"testing"
)

func TestPaymentAcceptance(t *testing.T) {
	db := newCommercialTestDB(t)
	app := NewCommercialServer(db, Config{})

	t.Run("TEST-PAYMENT-RECORD-001 creates payment plan as Unpaid", func(t *testing.T) {
		contractID := createContractForPayment(t, app, "opp_payment_plan_001")
		rec := postCommercialJSON(app, "/contracts/"+contractID+"/payment-plans", map[string]any{
			"dueAmount": "10000.00",
			"dueDate":   "2027-08-01",
			"currency":  "CNY",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected plan create 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "contractId", contractID)
		requireJSONValue(t, body, "status", "Unpaid")
		requireJSONValue(t, body, "dueAmount", "10000.00")
	})

	t.Run("TEST-PAYMENT-RECORD-002 and TEST-PAYMENT-RECORD-003 partial then full payment updates status", func(t *testing.T) {
		contractID := createContractForPayment(t, app, "opp_payment_record_001")
		createPaymentPlan(t, app, contractID, "10000.00")

		partial := postCommercialJSON(app, "/contracts/"+contractID+"/payments", map[string]any{
			"idempotencyKey": "pay-partial-001",
			"amount":         "4000.00",
			"paymentDate":    "2027-08-05",
			"note":           "Partial collection",
		}, actorHeaders("sales-1", "Sales"))
		if partial.Code != http.StatusCreated {
			t.Fatalf("expected partial payment 201, got %d body=%s", partial.Code, partial.Body.String())
		}
		partialBody := decodeJSON(t, partial)
		requireJSONValue(t, partialBody, "paymentStatus", "PartiallyPaid")
		requireJSONValue(t, partialBody, "remainingAmount", "6000.00")

		full := postCommercialJSON(app, "/contracts/"+contractID+"/payments", map[string]any{
			"idempotencyKey": "pay-full-001",
			"amount":         "6000.00",
			"paymentDate":    "2027-08-10",
		}, actorHeaders("sales-1", "Sales"))
		if full.Code != http.StatusCreated {
			t.Fatalf("expected full payment 201, got %d body=%s", full.Code, full.Body.String())
		}
		fullBody := decodeJSON(t, full)
		requireJSONValue(t, fullBody, "paymentStatus", "Paid")
		requireJSONValue(t, fullBody, "remainingAmount", "0.00")
		requireEvent(t, db, "PaymentRecorded", fullBody["paymentId"].(string))
	})

	t.Run("TEST-PAYMENT-GUARD-001 TEST-PAYMENT-GUARD-002 and TEST-INV-PAYAMOUNT-001 reject zero negative", func(t *testing.T) {
		contractID := createContractForPayment(t, app, "opp_payment_guard_amount")
		createPaymentPlan(t, app, contractID, "10000.00")
		for name, amount := range map[string]string{"zero": "0.00", "negative": "-1.00"} {
			t.Run(name, func(t *testing.T) {
				rec := postCommercialJSON(app, "/contracts/"+contractID+"/payments", map[string]any{
					"idempotencyKey": "pay-invalid-" + name,
					"amount":         amount,
					"paymentDate":    "2027-08-05",
				}, actorHeaders("sales-1", "Sales"))
				if rec.Code != http.StatusBadRequest {
					t.Fatalf("expected invalid amount 400, got %d body=%s", rec.Code, rec.Body.String())
				}
				requireErrorCode(t, rec, "INVALID_AMOUNT")
			})
		}
	})

	t.Run("TEST-PAYMENT-GUARD-003 and TEST-INV-OVERPAY-001 reject contract overpayment across plans", func(t *testing.T) {
		contractID := createContractForPayment(t, app, "opp_payment_overpay")
		createPaymentPlan(t, app, contractID, "7000.00")
		createPaymentPlan(t, app, contractID, "3000.00")
		ok := postCommercialJSON(app, "/contracts/"+contractID+"/payments", map[string]any{
			"idempotencyKey": "pay-ok-overpay-test",
			"amount":         "9000.00",
			"paymentDate":    "2027-08-05",
		}, actorHeaders("sales-1", "Sales"))
		if ok.Code != http.StatusCreated {
			t.Fatalf("expected initial payment 201, got %d body=%s", ok.Code, ok.Body.String())
		}
		over := postCommercialJSON(app, "/contracts/"+contractID+"/payments", map[string]any{
			"idempotencyKey": "pay-overpay-test",
			"amount":         "1000.01",
			"paymentDate":    "2027-08-06",
		}, actorHeaders("sales-1", "Sales"))
		if over.Code != http.StatusBadRequest {
			t.Fatalf("expected overpayment 400, got %d body=%s", over.Code, over.Body.String())
		}
		requireErrorCode(t, over, "OVERPAYMENT_BLOCKED")
	})

	t.Run("TEST-PAYMENT-GUARD-004 and TEST-INV-CURRENCY-001 enforce single currency", func(t *testing.T) {
		contractID := createContractForPayment(t, app, "opp_payment_currency")
		rec := postCommercialJSON(app, "/contracts/"+contractID+"/payment-plans", map[string]any{
			"dueAmount": "10000.00",
			"dueDate":   "2027-08-01",
			"currency":  "USD",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected invalid currency 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "SINGLE_CURRENCY_REQUIRED")
	})

	t.Run("TEST-PAYMENT-RECORD-002 idempotent record-payment repeats same result", func(t *testing.T) {
		contractID := createContractForPayment(t, app, "opp_payment_idempotent")
		createPaymentPlan(t, app, contractID, "10000.00")
		body := map[string]any{
			"idempotencyKey": "pay-idempotent-001",
			"amount":         "2500.00",
			"paymentDate":    "2027-08-05",
		}
		first := postCommercialJSON(app, "/contracts/"+contractID+"/payments", body, actorHeaders("sales-1", "Sales"))
		if first.Code != http.StatusCreated {
			t.Fatalf("expected first payment 201, got %d body=%s", first.Code, first.Body.String())
		}
		second := postCommercialJSON(app, "/contracts/"+contractID+"/payments", body, actorHeaders("sales-1", "Sales"))
		if second.Code != http.StatusOK {
			t.Fatalf("expected idempotent replay 200, got %d body=%s", second.Code, second.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, second), "paymentId", decodeJSON(t, first)["paymentId"].(string))
	})
}

func createContractForPayment(t *testing.T, app http.Handler, opportunityID string) string {
	t.Helper()
	quoteID := createAcceptedQuote(t, app, opportunityID, "sales-1")
	rec := postCommercialJSON(app, "/contracts", contractCreateBody(quoteID, opportunityID, "acct_"+opportunityID, "10000.00"), actorHeaders("sales-1", "Sales"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected contract create 201, got %d body=%s", rec.Code, rec.Body.String())
	}
	return decodeJSON(t, rec)["id"].(string)
}

func createPaymentPlan(t *testing.T, app http.Handler, contractID, dueAmount string) string {
	t.Helper()
	rec := postCommercialJSON(app, "/contracts/"+contractID+"/payment-plans", map[string]any{
		"dueAmount": dueAmount,
		"dueDate":   "2027-08-01",
		"currency":  "CNY",
	}, actorHeaders("sales-1", "Sales"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected payment plan create 201, got %d body=%s", rec.Code, rec.Body.String())
	}
	return decodeJSON(t, rec)["id"].(string)
}
