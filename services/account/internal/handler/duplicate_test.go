package handler

import (
	"net/http"
	"testing"
)

func TestAccountDuplicateWarningsAcceptance(t *testing.T) {
	db := newAccountTestDB(t)
	app := NewAccountServer(db, Config{})

	t.Run("TEST-DUPLICATE-WARN-001/005 and TEST-ABUSE-DUPENUM-001 company name warning and proceed creates a second record only", func(t *testing.T) {
		first := postAccountJSON(app, "/accounts", map[string]any{
			"companyName":    "Acme Duplicate",
			"customerStatus": "Prospect",
			"ownerId":        "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if first.Code != http.StatusCreated {
			t.Fatalf("expected initial create 201, got %d body=%s", first.Code, first.Body.String())
		}
		firstID := decodeJSON(t, first)["id"].(string)

		warning := postAccountJSON(app, "/duplicate-checks", map[string]any{
			"targetType": "account",
			"candidate":  map[string]any{"companyName": "  acme   duplicate "},
		}, actorHeaders("sales-1", "Sales"))
		if warning.Code != http.StatusOK {
			t.Fatalf("expected duplicate warning 200, got %d body=%s", warning.Code, warning.Body.String())
		}
		body := decodeJSON(t, warning)
		requireJSONValue(t, body, "result", "PossibleDuplicate")
		token := body["warningToken"].(string)
		if token == "" || contains(warning.Body.String(), firstID) {
			t.Fatalf("warning token missing or leaked matched id: %s", warning.Body.String())
		}

		duplicate := postAccountJSON(app, "/accounts", map[string]any{
			"companyName":         "ACME Duplicate",
			"customerStatus":      "Prospect",
			"ownerId":             "sales-1",
			"proceedWarningToken": token,
		}, actorHeaders("sales-1", "Sales"))
		if duplicate.Code != http.StatusCreated {
			t.Fatalf("expected proceed create 201, got %d body=%s", duplicate.Code, duplicate.Body.String())
		}
		secondID := decodeJSON(t, duplicate)["id"].(string)
		if secondID == firstID {
			t.Fatalf("proceed must create a new record, got same id %s", secondID)
		}
		reuse := postAccountJSON(app, "/accounts", map[string]any{
			"companyName":         "ACME Duplicate",
			"customerStatus":      "Prospect",
			"ownerId":             "sales-1",
			"proceedWarningToken": token,
		}, actorHeaders("sales-1", "Sales"))
		if reuse.Code != http.StatusConflict {
			t.Fatalf("expected single-use token conflict, got %d body=%s", reuse.Code, reuse.Body.String())
		}
		requireErrorCode(t, reuse, "DUPLICATE_WARNING_TOKEN_USED")
	})

	t.Run("TEST-DUPLICATE-WARN-002/003 contact phone and email normalized warning", func(t *testing.T) {
		accountID := createAccountForContact(t, app, "Duplicate Contact Parent", "sales-1")
		create := postAccountJSON(app, "/accounts/"+accountID+"/contacts", map[string]any{
			"contactName": "Buyer One",
			"email":       "Buyer@Example.COM",
			"phone":       "+86 138-0000-0000",
		}, actorHeaders("sales-1", "Sales"))
		if create.Code != http.StatusCreated {
			t.Fatalf("expected contact create 201, got %d body=%s", create.Code, create.Body.String())
		}
		emailWarning := postAccountJSON(app, "/duplicate-checks", map[string]any{
			"targetType": "contact",
			"candidate":  map[string]any{"email": "buyer@example.com"},
		}, actorHeaders("sales-1", "Sales"))
		if emailWarning.Code != http.StatusOK || !contains(emailWarning.Body.String(), "PossibleDuplicate") {
			t.Fatalf("expected email duplicate warning, got %d body=%s", emailWarning.Code, emailWarning.Body.String())
		}
		phoneWarning := postAccountJSON(app, "/duplicate-checks", map[string]any{
			"targetType": "contact",
			"candidate":  map[string]any{"phone": "13800000000"},
		}, actorHeaders("sales-1", "Sales"))
		if phoneWarning.Code != http.StatusOK || !contains(phoneWarning.Body.String(), "PossibleDuplicate") {
			t.Fatalf("expected phone duplicate warning, got %d body=%s", phoneWarning.Code, phoneWarning.Body.String())
		}
	})

	t.Run("TEST-DUPLICATE-WARN-006 unique account has no warning", func(t *testing.T) {
		unique := postAccountJSON(app, "/duplicate-checks", map[string]any{
			"targetType": "account",
			"candidate":  map[string]any{"companyName": "Unique Account Without Duplicate"},
		}, actorHeaders("sales-1", "Sales"))
		if unique.Code != http.StatusOK {
			t.Fatalf("expected unique check 200, got %d body=%s", unique.Code, unique.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, unique), "result", "NoDuplicate")
	})

	t.Run("TEST-DUPLICATE-WARN-TX-001 outbox failure rolls back warning token", func(t *testing.T) {
		first := postAccountJSON(app, "/accounts", map[string]any{
			"companyName":    "Transactional Duplicate",
			"customerStatus": "Prospect",
			"ownerId":        "sales-1",
		}, actorHeaders("sales-1", "Sales"))
		if first.Code != http.StatusCreated {
			t.Fatalf("expected initial create 201, got %d body=%s", first.Code, first.Body.String())
		}
		if _, err := db.Exec(`
			CREATE OR REPLACE FUNCTION account.fail_duplicate_outbox_insert() RETURNS trigger AS $$
			BEGIN
				RAISE EXCEPTION 'forced duplicate outbox failure';
			END;
			$$ LANGUAGE plpgsql;
			CREATE TRIGGER fail_duplicate_outbox_insert
			BEFORE INSERT ON account.outbox_events
			FOR EACH ROW EXECUTE FUNCTION account.fail_duplicate_outbox_insert();
		`); err != nil {
			t.Fatalf("install failing outbox trigger: %v", err)
		}
		t.Cleanup(func() {
			_, _ = db.Exec(`DROP TRIGGER IF EXISTS fail_duplicate_outbox_insert ON account.outbox_events; DROP FUNCTION IF EXISTS account.fail_duplicate_outbox_insert();`)
		})
		var beforeCount int
		if err := db.QueryRow(`SELECT count(*) FROM account.duplicate_warning_tokens WHERE actor_user_id = $1`, "sales-1").Scan(&beforeCount); err != nil {
			t.Fatalf("count duplicate warning tokens before failure: %v", err)
		}

		warning := postAccountJSON(app, "/duplicate-checks", map[string]any{
			"targetType": "account",
			"candidate":  map[string]any{"companyName": "Transactional Duplicate"},
		}, actorHeaders("sales-1", "Sales"))
		if warning.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected outbox failure 503, got %d body=%s", warning.Code, warning.Body.String())
		}
		var count int
		if err := db.QueryRow(`SELECT count(*) FROM account.duplicate_warning_tokens WHERE actor_user_id = $1`, "sales-1").Scan(&count); err != nil {
			t.Fatalf("count duplicate warning tokens: %v", err)
		}
		if count != beforeCount {
			t.Fatalf("expected duplicate warning token insert rolled back with outbox failure, before=%d after=%d", beforeCount, count)
		}
	})
}
