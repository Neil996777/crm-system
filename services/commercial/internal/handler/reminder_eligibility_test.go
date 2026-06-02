package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommercialReminderEligibility(t *testing.T) {
	db := newCommercialTestDB(t)
	app := NewCommercialServer(db, Config{ServiceID: "commercial", ServiceTokenSecret: []byte("commercial_test_secret")})

	t.Run("TEST-REMINDER-002 returns pending-signature contract past expected date and suppresses signed", func(t *testing.T) {
		contractID := createPendingContractWithDate(t, app, "opp_reminder_contract_pending", "sales-1", "2026-01-01")
		signedID := createPendingContractWithDate(t, app, "opp_reminder_contract_signed", "sales-1", "2026-01-01")
		sign := postCommercialJSON(app, "/contracts/"+signedID+"/status", map[string]any{
			"expectedVersion":     1,
			"toStatus":            "Signed",
			"signedEffectiveDate": "2026-01-02",
		}, actorHeaders("sales-1", "Sales"))
		if sign.Code != http.StatusOK {
			t.Fatalf("expected sign 200, got %d body=%s", sign.Code, sign.Body.String())
		}

		rec := reminderEligibilityQuery(t, app, "2026-06-02", "sales-1", "Sales")
		if rec.Code != http.StatusOK {
			t.Fatalf("expected eligibility 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		if !contains(rec.Body.String(), contractID) || !contains(rec.Body.String(), "contract_pending_signature") {
			t.Fatalf("expected pending contract reminder, got body=%s", rec.Body.String())
		}
		if contains(rec.Body.String(), signedID) {
			t.Fatalf("signed contract must be suppressed from reminders: %s", rec.Body.String())
		}
	})

	t.Run("TEST-REMINDER-003 TEST-PAYMENT-OVERDUE-001 returns unpaid overdue payment and suppresses fully paid", func(t *testing.T) {
		overdueID := createContractForPayment(t, app, "opp_reminder_payment_overdue")
		createPaymentPlanWithDate(t, app, overdueID, "10000.00", "2026-01-01")

		paidID := createContractForPayment(t, app, "opp_reminder_payment_paid")
		createPaymentPlanWithDate(t, app, paidID, "10000.00", "2026-01-01")
		full := postCommercialJSON(app, "/contracts/"+paidID+"/payments", map[string]any{
			"idempotencyKey": "pay-reminder-full",
			"amount":         "10000.00",
			"paymentDate":    "2026-01-02",
		}, actorHeaders("sales-1", "Sales"))
		if full.Code != http.StatusCreated {
			t.Fatalf("expected payment 201, got %d body=%s", full.Code, full.Body.String())
		}

		rec := reminderEligibilityQuery(t, app, "2026-06-02", "sales-1", "Sales")
		if rec.Code != http.StatusOK {
			t.Fatalf("expected eligibility 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		if !contains(rec.Body.String(), overdueID) || !contains(rec.Body.String(), "payment_overdue") {
			t.Fatalf("expected overdue payment reminder, got body=%s", rec.Body.String())
		}
		if contains(rec.Body.String(), paidID) {
			t.Fatalf("fully paid contract payment must be suppressed: %s", rec.Body.String())
		}
	})
}

func reminderEligibilityQuery(t *testing.T, handler http.Handler, businessDate, actorID, actorRole string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/internal/reminders/eligibility?businessDate="+businessDate, nil)
	req.Header.Set("Authorization", "Bearer "+makeServiceToken(t, "work", "commercial", "commercial.reminder_eligibility", []byte("commercial_test_secret")))
	req.Header.Set("X-Service-Id", "work")
	req.Header.Set("X-Intent", "commercial.reminder_eligibility")
	req.Header.Set("X-Actor-User-Id", actorID)
	req.Header.Set("X-Actor-Role", actorRole)
	handler.ServeHTTP(rec, req)
	return rec
}

func createPendingContractWithDate(t *testing.T, app http.Handler, opportunityID, ownerID, expectedSignedDate string) string {
	t.Helper()
	quoteID := createAcceptedQuote(t, app, opportunityID, ownerID)
	body := contractCreateBody(quoteID, opportunityID, "acct_"+opportunityID, "10000.00")
	body["expectedSignedDate"] = expectedSignedDate
	rec := postCommercialJSON(app, "/contracts", body, actorHeaders(ownerID, "Sales"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected contract create 201, got %d body=%s", rec.Code, rec.Body.String())
	}
	return decodeJSON(t, rec)["id"].(string)
}

func createPaymentPlanWithDate(t *testing.T, app http.Handler, contractID, dueAmount, dueDate string) string {
	t.Helper()
	rec := postCommercialJSON(app, "/contracts/"+contractID+"/payment-plans", map[string]any{
		"dueAmount": dueAmount,
		"dueDate":   dueDate,
		"currency":  "CNY",
	}, actorHeaders("sales-1", "Sales"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected payment plan create 201, got %d body=%s", rec.Code, rec.Body.String())
	}
	return decodeJSON(t, rec)["id"].(string)
}
