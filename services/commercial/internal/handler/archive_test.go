package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommercialArchiveAcceptance(t *testing.T) {
	db := newCommercialTestDB(t)
	app := NewCommercialServer(db, Config{})

	t.Run("TEST-ARCHIVE-001/002/003 contract archive blocks pending signature then succeeds after signed", func(t *testing.T) {
		contractID := createContractForPayment(t, app, "opp_archive_contract")
		blocked := postCommercialJSON(app, "/contracts/"+contractID+"/archive", map[string]any{
			"expectedVersion": 1,
			"reason":          "Archive pending signature",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if blocked.Code != http.StatusConflict {
			t.Fatalf("expected archive block 409, got %d body=%s", blocked.Code, blocked.Body.String())
		}
		requireErrorCode(t, blocked, "ARCHIVE_BLOCKED_ACTIVE_OBLIGATION")

		signed := postCommercialJSON(app, "/contracts/"+contractID+"/status", map[string]any{
			"expectedVersion":       1,
			"toStatus":              "Signed",
			"signedEffectiveDate":   "2027-01-16",
		}, actorHeaders("sales-1", "Sales"))
		if signed.Code != http.StatusOK {
			t.Fatalf("expected signed 200, got %d body=%s", signed.Code, signed.Body.String())
		}
		archive := postCommercialJSON(app, "/contracts/"+contractID+"/archive", map[string]any{
			"expectedVersion": 2,
			"reason":          "No active obligations remain",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if archive.Code != http.StatusOK {
			t.Fatalf("expected archive 200, got %d body=%s", archive.Code, archive.Body.String())
		}
		if decodeJSON(t, archive)["archived"] != true {
			t.Fatalf("expected archived true, got body=%s", archive.Body.String())
		}
		requireEvent(t, db, "ContractArchived", contractID)
	})

	t.Run("TEST-INV-ARCHIVEBLOCK-001 payment plan archive blocks unpaid and succeeds when paid", func(t *testing.T) {
		contractID := createContractForPayment(t, app, "opp_archive_payment")
		planID := createPaymentPlan(t, app, contractID, "10000.00")
		blocked := postCommercialJSON(app, "/payment-plans/"+planID+"/archive", map[string]any{
			"expectedVersion": 1,
			"reason":          "Archive unpaid payment",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if blocked.Code != http.StatusConflict {
			t.Fatalf("expected unpaid archive block 409, got %d body=%s", blocked.Code, blocked.Body.String())
		}
		requireErrorCode(t, blocked, "ARCHIVE_BLOCKED_ACTIVE_OBLIGATION")

		paid := postCommercialJSON(app, "/contracts/"+contractID+"/payments", map[string]any{
			"idempotencyKey": "archive-payment-paid",
			"amount":         "10000.00",
			"paymentDate":    "2027-08-10",
		}, actorHeaders("sales-1", "Sales"))
		if paid.Code != http.StatusCreated {
			t.Fatalf("expected payment 201, got %d body=%s", paid.Code, paid.Body.String())
		}
		archive := postCommercialJSON(app, "/payment-plans/"+planID+"/archive", map[string]any{
			"expectedVersion": 2,
			"reason":          "Payment obligation resolved",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if archive.Code != http.StatusOK {
			t.Fatalf("expected paid payment archive 200, got %d body=%s", archive.Code, archive.Body.String())
		}
		if decodeJSON(t, archive)["archived"] != true {
			t.Fatalf("expected archived true, got body=%s", archive.Body.String())
		}
	})
}

func getCommercialJSON(handler http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	handler.ServeHTTP(rec, req)
	return rec
}
