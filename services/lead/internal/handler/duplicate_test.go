package handler

import (
	"net/http"
	"testing"
)

func TestLeadDuplicateWarningsAcceptance(t *testing.T) {
	db := newLeadTestDB(t)
	app := NewLeadServer(db, Config{})

	t.Run("TEST-DUPLICATE-WARN-004/005 and TEST-ABUSE-DUPENUM-001 lead company/contact warning and proceed creates a second record only", func(t *testing.T) {
		first := postLeadJSON(app, "/leads", map[string]any{
			"leadName":    "ERP buyer",
			"companyName": "Lead Duplicate Co",
			"source":      "Website",
			"ownerId":     "sales-1",
			"email":       "LeadBuyer@Example.COM",
			"phone":       "+86 (139) 0000-0000",
		}, actorHeaders("sales-1", "Sales"))
		if first.Code != http.StatusCreated {
			t.Fatalf("expected initial create 201, got %d body=%s", first.Code, first.Body.String())
		}
		firstID := decodeJSON(t, first)["id"].(string)

		warning := postLeadJSON(app, "/duplicate-checks", map[string]any{
			"targetType": "lead",
			"candidate": map[string]any{
				"companyName": "  lead   duplicate   co ",
				"email":       "leadbuyer@example.com",
				"phone":       "13900000000",
			},
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

		duplicate := postLeadJSON(app, "/leads", map[string]any{
			"leadName":            "ERP buyer second",
			"companyName":         "LEAD Duplicate Co",
			"source":              "Website",
			"ownerId":             "sales-1",
			"email":               "leadbuyer@example.com",
			"phone":               "13900000000",
			"proceedWarningToken": token,
		}, actorHeaders("sales-1", "Sales"))
		if duplicate.Code != http.StatusCreated {
			t.Fatalf("expected proceed create 201, got %d body=%s", duplicate.Code, duplicate.Body.String())
		}
		secondID := decodeJSON(t, duplicate)["id"].(string)
		if secondID == firstID {
			t.Fatalf("proceed must create a new record, got same id %s", secondID)
		}

		reuse := postLeadJSON(app, "/leads", map[string]any{
			"leadName":            "ERP buyer third",
			"companyName":         "Lead Duplicate Co",
			"source":              "Website",
			"ownerId":             "sales-1",
			"email":               "leadbuyer@example.com",
			"phone":               "13900000000",
			"proceedWarningToken": token,
		}, actorHeaders("sales-1", "Sales"))
		if reuse.Code != http.StatusConflict {
			t.Fatalf("expected single-use token conflict, got %d body=%s", reuse.Code, reuse.Body.String())
		}
		requireErrorCode(t, reuse, "DUPLICATE_WARNING_TOKEN_USED")
	})

	t.Run("TEST-DUPLICATE-WARN-006 unique lead has no warning", func(t *testing.T) {
		unique := postLeadJSON(app, "/duplicate-checks", map[string]any{
			"targetType": "lead",
			"candidate":  map[string]any{"companyName": "Unique Lead Without Duplicate"},
		}, actorHeaders("sales-1", "Sales"))
		if unique.Code != http.StatusOK {
			t.Fatalf("expected unique check 200, got %d body=%s", unique.Code, unique.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, unique), "result", "NoDuplicate")
	})
}
