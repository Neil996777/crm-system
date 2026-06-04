package handler

import (
	"net/http"
	"testing"
)

func TestCommercialSingleRecordReadAuthorization(t *testing.T) {
	db := newCommercialTestDB(t)
	app := NewCommercialServer(db, Config{})

	t.Run("TEST-AUTHZ-SCOPE-005 quote by id denies non-owner Sales without leaking record", func(t *testing.T) {
		quoteID := createQuoteForLifecycle(t, app, "opp_quote_read_scope", "sales-1")

		denied := getCommercialJSON(app, "/quotes/"+quoteID, actorHeaders("sales-2", "Sales"))
		if denied.Code != http.StatusNotFound {
			t.Fatalf("expected safe not-found denial, got %d body=%s", denied.Code, denied.Body.String())
		}
		requireErrorCode(t, denied, "NOT_FOUND")
		if contains(denied.Body.String(), quoteID) || contains(denied.Body.String(), "opp_quote_read_scope") {
			t.Fatalf("denial leaked quote record data: %s", denied.Body.String())
		}

		owner := getCommercialJSON(app, "/quotes/"+quoteID, actorHeaders("sales-1", "Sales"))
		if owner.Code != http.StatusOK {
			t.Fatalf("expected owner read 200, got %d body=%s", owner.Code, owner.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, owner), "id", quoteID)

		manager := getCommercialJSON(app, "/quotes/"+quoteID, actorHeaders("mgr-1", "Sales Manager"))
		if manager.Code != http.StatusOK {
			t.Fatalf("expected manager read 200, got %d body=%s", manager.Code, manager.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, manager), "id", quoteID)
	})

	t.Run("TEST-AUTHZ-SCOPE-005 contract by id denies non-owner Sales without leaking record", func(t *testing.T) {
		contractID := createPendingContract(t, app, "opp_contract_read_scope", "sales-1")

		denied := getCommercialJSON(app, "/contracts/"+contractID, actorHeaders("sales-2", "Sales"))
		if denied.Code != http.StatusNotFound {
			t.Fatalf("expected safe not-found denial, got %d body=%s", denied.Code, denied.Body.String())
		}
		requireErrorCode(t, denied, "NOT_FOUND")
		if contains(denied.Body.String(), contractID) || contains(denied.Body.String(), "opp_contract_read_scope") {
			t.Fatalf("denial leaked contract record data: %s", denied.Body.String())
		}

		owner := getCommercialJSON(app, "/contracts/"+contractID, actorHeaders("sales-1", "Sales"))
		if owner.Code != http.StatusOK {
			t.Fatalf("expected owner read 200, got %d body=%s", owner.Code, owner.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, owner), "id", contractID)

		admin := getCommercialJSON(app, "/contracts/"+contractID, actorHeaders("admin-1", "Administrator"))
		if admin.Code != http.StatusOK {
			t.Fatalf("expected admin read 200, got %d body=%s", admin.Code, admin.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, admin), "id", contractID)
	})
}
