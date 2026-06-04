package handler

import (
	"net/http"
	"testing"
)

func TestContactLinkAcceptance(t *testing.T) {
	db := newAccountTestDB(t)
	app := NewAccountServer(db, Config{})

	t.Run("TEST-CONTACT-LINK-001 creates contact with company and method or role note", func(t *testing.T) {
		accountID := createAccountForContact(t, app, "Contact Parent Co", "sales-1")
		rec := postAccountJSON(app, "/accounts/"+accountID+"/contacts", map[string]any{
			"contactName": "Ada Buyer",
			"email":       "ada@example.com",
			"roleNote":    "Decision maker",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "accountId", accountID)
		requireJSONValue(t, body, "contactName", "Ada Buyer")
		requireEvent(t, db, "ContactCreated", body["id"].(string))
	})

	t.Run("TEST-CONTACT-LINK-002 save without company blocked", func(t *testing.T) {
		rec := postAccountJSON(app, "/contacts", map[string]any{
			"contactName": "No Company",
			"email":       "missing-company@example.com",
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "VALIDATION_FAILED")
	})

	t.Run("TEST-CONTACT-LINK-003 multiple contacts visible in company context", func(t *testing.T) {
		accountID := createAccountForContact(t, app, "Multiple Contact Co", "sales-1")
		for _, name := range []string{"Primary Contact", "Technical Contact"} {
			rec := postAccountJSON(app, "/accounts/"+accountID+"/contacts", map[string]any{
				"contactName": name,
				"phone":       "13800000000",
			}, actorHeaders("sales-1", "Sales"))
			if rec.Code != http.StatusCreated {
				t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
			}
		}
		list := getAccountJSON(app, "/accounts/"+accountID+"/contacts", actorHeaders("sales-1", "Sales"))
		if list.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", list.Code, list.Body.String())
		}
		items := decodeJSON(t, list)["items"].([]any)
		if len(items) != 2 {
			t.Fatalf("expected 2 contacts, got %#v", items)
		}
	})

	t.Run("TEST-CONTACT-LINK-004 unrelated Sales denied", func(t *testing.T) {
		accountID := createAccountForContact(t, app, "Restricted Contact Co", "sales-2")
		create := postAccountJSON(app, "/accounts/"+accountID+"/contacts", map[string]any{
			"contactName": "Hidden Contact",
			"email":       "hidden@example.com",
		}, actorHeaders("sales-1", "Sales"))
		if create.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d body=%s", create.Code, create.Body.String())
		}
		requireErrorCode(t, create, "PERMISSION_DENIED")
		if contains(create.Body.String(), "Hidden Contact") {
			t.Fatalf("unauthorized response leaked contact details: %s", create.Body.String())
		}
	})

	t.Run("TEST-DENIAL-CONTACTS-001 unrelated Sales list gets safe 404", func(t *testing.T) {
		accountID := createAccountForContact(t, app, "Restricted Contact List Co", "sales-2")
		contact := postAccountJSON(app, "/accounts/"+accountID+"/contacts", map[string]any{
			"contactName": "Hidden List Contact",
			"email":       "hidden-list@example.com",
		}, actorHeaders("mgr-1", "Sales Manager"))
		if contact.Code != http.StatusCreated {
			t.Fatalf("expected manager contact create 201, got %d body=%s", contact.Code, contact.Body.String())
		}

		list := getAccountJSON(app, "/accounts/"+accountID+"/contacts", actorHeaders("sales-1", "Sales"))

		if list.Code != http.StatusNotFound {
			t.Fatalf("expected safe 404 for unreadable account contacts, got %d body=%s", list.Code, list.Body.String())
		}
		requireErrorCode(t, list, "NOT_FOUND")
		if contains(list.Body.String(), "Hidden List Contact") {
			t.Fatalf("unauthorized response leaked contact details: %s", list.Body.String())
		}
	})
}

func createAccountForContact(t *testing.T, app http.Handler, companyName, ownerID string) string {
	t.Helper()
	create := postAccountJSON(app, "/accounts", map[string]any{
		"companyName":    companyName,
		"customerStatus": "Prospect",
		"ownerId":        ownerID,
	}, actorHeaders("mgr-1", "Sales Manager"))
	if create.Code != http.StatusCreated {
		t.Fatalf("expected account create 201, got %d body=%s", create.Code, create.Body.String())
	}
	return decodeJSON(t, create)["id"].(string)
}
