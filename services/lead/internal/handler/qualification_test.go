package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLeadQualificationAndConversionAcceptance(t *testing.T) {
	db := newLeadTestDB(t)
	accountCalls := 0
	opportunityCalls := 0
	accountServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountCalls++
		requireS2SRequest(t, r, "account.create_for_lead_conversion")
		writeJSON(w, http.StatusCreated, map[string]any{
			"id":             "acct_from_lead",
			"companyName":    "Convert Co",
			"customerStatus": "Prospect",
			"ownerId":        "sales-1",
			"version":        1,
		})
	}))
	t.Cleanup(accountServer.Close)
	opportunityServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opportunityCalls++
		requireS2SRequest(t, r, "opportunity.create_for_lead_conversion")
		writeJSON(w, http.StatusCreated, map[string]any{
			"id":                "opp_from_lead",
			"customerId":        "acct_from_lead",
			"ownerId":           "sales-1",
			"stage":             "New Opportunity",
			"expectedAmount":    "50000.00",
			"expectedCloseDate": "2026-10-01",
			"version":           1,
		})
	}))
	t.Cleanup(opportunityServer.Close)
	auditServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requireS2SRequest(t, r, "audit.append")
		writeJSON(w, http.StatusCreated, map[string]any{"eventUid": "evt_test"})
	}))
	t.Cleanup(auditServer.Close)

	app := NewLeadServer(db, Config{
		AccountServiceURL:      accountServer.URL,
		OpportunityServiceURL:  opportunityServer.URL,
		AuditHistoryServiceURL: auditServer.URL,
		ServiceID:              "lead",
		ServiceTokenSecret:     []byte("task-008-secret"),
	})

	t.Run("TEST-LEAD-QUALIFY-001 Pending to Valid with authz", func(t *testing.T) {
		leadID := createOwnedLead(t, app, "Valid Co", "sales-1")
		rec := postLeadJSON(app, "/leads/"+leadID+"/qualify-valid", map[string]any{"expectedVersion": 1}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, rec), "status", "Valid")
		requireEvent(t, db, "LeadQualified", leadID)
	})

	t.Run("TEST-LEAD-QUALIFY-002 Pending to Invalid requires reason", func(t *testing.T) {
		leadID := createOwnedLead(t, app, "Invalid Co", "sales-1")
		missing := postLeadJSON(app, "/leads/"+leadID+"/qualify-invalid", map[string]any{"expectedVersion": 1}, actorHeaders("sales-1", "Sales"))
		if missing.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", missing.Code, missing.Body.String())
		}
		rec := postLeadJSON(app, "/leads/"+leadID+"/qualify-invalid", map[string]any{"expectedVersion": 1, "invalidReason": "No fit"}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "status", "Invalid")
		requireJSONValue(t, body, "invalidReason", "No fit")
		requireEvent(t, db, "LeadDisqualified", leadID)
	})

	t.Run("TEST-LEAD-QUALIFY-003 Valid converts through S2S and preserves lead history", func(t *testing.T) {
		leadID := createOwnedLead(t, app, "Convert Co", "sales-1")
		valid := postLeadJSON(app, "/leads/"+leadID+"/qualify-valid", map[string]any{"expectedVersion": 1}, actorHeaders("sales-1", "Sales"))
		if valid.Code != http.StatusOK {
			t.Fatalf("expected valid 200, got %d body=%s", valid.Code, valid.Body.String())
		}
		rec := postLeadJSON(app, "/leads/"+leadID+"/convert", map[string]any{
			"idempotencyKey": "convert-key-1",
			"target": map[string]any{
				"accountInput": map[string]any{
					"companyName":    "Convert Co",
					"customerStatus": "Prospect",
					"ownerId":        "sales-1",
				},
				"opportunityInput": map[string]any{
					"ownerId":           "sales-1",
					"stage":             "New Opportunity",
					"expectedAmount":    "50000.00",
					"expectedCloseDate": "2026-10-01",
					"title":             "Converted need",
				},
			},
		}, actorHeaders("sales-1", "Sales"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "status", "Converted To Opportunity")
		requireJSONValue(t, body, "accountId", "acct_from_lead")
		requireJSONValue(t, body, "opportunityId", "opp_from_lead")
		requireEvent(t, db, "LeadCreated", leadID)
		requireEvent(t, db, "LeadConverted", leadID)
	})

	t.Run("TEST-LEAD-QUALIFY-004 and TEST-ABUSE-BRBYPASS-001 Unassigned qualify or convert rejected", func(t *testing.T) {
		create := postLeadJSON(app, "/leads", map[string]any{
			"companyName": "Unassigned Convert Co",
			"source":      "Website",
		}, actorHeaders("mgr-1", "Sales Manager"))
		leadID := decodeJSON(t, create)["id"].(string)
		qualify := postLeadJSON(app, "/leads/"+leadID+"/qualify-valid", map[string]any{"expectedVersion": 1}, actorHeaders("mgr-1", "Sales Manager"))
		if qualify.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d body=%s", qualify.Code, qualify.Body.String())
		}
		convert := postLeadJSON(app, "/leads/"+leadID+"/convert", map[string]any{"idempotencyKey": "unassigned"}, actorHeaders("mgr-1", "Sales Manager"))
		if convert.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d body=%s", convert.Code, convert.Body.String())
		}
	})

	t.Run("TEST-LEAD-QUALIFY-005 Invalid convert rejected until restored", func(t *testing.T) {
		leadID := createOwnedLead(t, app, "Restore Co", "sales-1")
		invalid := postLeadJSON(app, "/leads/"+leadID+"/qualify-invalid", map[string]any{"expectedVersion": 1, "invalidReason": "No budget"}, actorHeaders("sales-1", "Sales"))
		if invalid.Code != http.StatusOK {
			t.Fatalf("expected invalid 200, got %d body=%s", invalid.Code, invalid.Body.String())
		}
		convert := postLeadJSON(app, "/leads/"+leadID+"/convert", map[string]any{"idempotencyKey": "invalid"}, actorHeaders("sales-1", "Sales"))
		if convert.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", convert.Code, convert.Body.String())
		}
		requireErrorCode(t, convert, "INVALID_LEAD_STATE")
	})

	t.Run("TEST-LEAD-QUALIFY-006 Invalid restore by admin or manager only", func(t *testing.T) {
		leadID := createOwnedLead(t, app, "Admin Restore Co", "sales-1")
		invalid := postLeadJSON(app, "/leads/"+leadID+"/qualify-invalid", map[string]any{"expectedVersion": 1, "invalidReason": "Duplicate"}, actorHeaders("sales-1", "Sales"))
		if invalid.Code != http.StatusOK {
			t.Fatalf("expected invalid 200, got %d body=%s", invalid.Code, invalid.Body.String())
		}
		salesRestore := postLeadJSON(app, "/leads/"+leadID+"/restore-invalid", map[string]any{"expectedVersion": 2}, actorHeaders("sales-1", "Sales"))
		if salesRestore.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d body=%s", salesRestore.Code, salesRestore.Body.String())
		}
		managerRestore := postLeadJSON(app, "/leads/"+leadID+"/restore-invalid", map[string]any{"expectedVersion": 2}, actorHeaders("mgr-1", "Sales Manager"))
		if managerRestore.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", managerRestore.Code, managerRestore.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, managerRestore), "status", "Pending Qualification")
	})

	t.Run("TEST-LEAD-QUALIFY-007 re-convert rejected while same idempotency key is stable", func(t *testing.T) {
		startAccountCalls := accountCalls
		startOpportunityCalls := opportunityCalls
		leadID := createOwnedLead(t, app, "Idempotent Co", "sales-1")
		_ = postLeadJSON(app, "/leads/"+leadID+"/qualify-valid", map[string]any{"expectedVersion": 1}, actorHeaders("sales-1", "Sales"))
		body := map[string]any{
			"idempotencyKey": "same-key",
			"target": map[string]any{
				"accountInput": map[string]any{
					"companyName":    "Idempotent Co",
					"customerStatus": "Prospect",
					"ownerId":        "sales-1",
				},
				"opportunityInput": map[string]any{
					"ownerId":           "sales-1",
					"stage":             "New Opportunity",
					"expectedAmount":    "70000.00",
					"expectedCloseDate": "2026-11-01",
					"title":             "Idempotent conversion",
				},
			},
		}
		first := postLeadJSON(app, "/leads/"+leadID+"/convert", body, actorHeaders("sales-1", "Sales"))
		if first.Code != http.StatusOK {
			t.Fatalf("expected first convert 200, got %d body=%s", first.Code, first.Body.String())
		}
		second := postLeadJSON(app, "/leads/"+leadID+"/convert", body, actorHeaders("sales-1", "Sales"))
		if second.Code != http.StatusOK {
			t.Fatalf("expected idempotent retry 200, got %d body=%s", second.Code, second.Body.String())
		}
		if accountCalls != startAccountCalls+1 || opportunityCalls != startOpportunityCalls+1 {
			t.Fatalf("idempotent retry must not call downstream again, account=%d opportunity=%d", accountCalls-startAccountCalls, opportunityCalls-startOpportunityCalls)
		}
		reconvert := postLeadJSON(app, "/leads/"+leadID+"/convert", map[string]any{"idempotencyKey": "different-key"}, actorHeaders("sales-1", "Sales"))
		if reconvert.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", reconvert.Code, reconvert.Body.String())
		}
		requireErrorCode(t, reconvert, "LEAD_ALREADY_CONVERTED")
	})
}

func TestLeadConversionRetriesUseDownstreamIdempotencyKeys(t *testing.T) {
	db := newLeadTestDB(t)
	accountCreates := 0
	accountByKey := map[string]string{}
	accountServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Service-Id") != "lead" || r.Header.Get("X-Intent") != "account.create_for_lead_conversion" {
			t.Fatalf("missing account conversion S2S headers: service=%q intent=%q", r.Header.Get("X-Service-Id"), r.Header.Get("X-Intent"))
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Fatalf("missing bearer service token")
		}
		if r.Header.Get("X-Actor-User-Id") == "" || r.Header.Get("X-Actor-Role") == "" {
			t.Fatalf("missing actor context")
		}
		if r.Header.Get("X-Correlation-Id") != "corr-lead-convert" {
			t.Fatalf("missing account conversion correlation id: %q", r.Header.Get("X-Correlation-Id"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode account request: %v", err)
		}
		key, _ := body["idempotencyKey"].(string)
		accountID := accountByKey[key]
		if accountID == "" {
			accountCreates++
			accountID = "acct_retry_" + string(rune('0'+accountCreates))
			if key != "" {
				accountByKey[key] = accountID
			}
		}
		writeJSON(w, http.StatusCreated, map[string]any{"id": accountID})
	}))
	t.Cleanup(accountServer.Close)
	opportunityCalls := 0
	opportunityServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opportunityCalls++
		requireS2SRequest(t, r, "opportunity.create_for_lead_conversion")
		if r.Header.Get("X-Correlation-Id") != "corr-lead-convert" {
			t.Fatalf("missing opportunity conversion correlation id: %q", r.Header.Get("X-Correlation-Id"))
		}
		if opportunityCalls == 1 {
			http.Error(w, "temporary failure", http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"id": "opp_retry_1"})
	}))
	t.Cleanup(opportunityServer.Close)
	app := NewLeadServer(db, Config{
		AccountServiceURL:     accountServer.URL,
		OpportunityServiceURL: opportunityServer.URL,
		ServiceID:             "lead",
		ServiceTokenSecret:    []byte("task-008-secret"),
	})

	leadID := createOwnedLead(t, app, "Retry Conversion Co", "sales-1")
	valid := postLeadJSON(app, "/leads/"+leadID+"/qualify-valid", map[string]any{"expectedVersion": 1}, actorHeaders("sales-1", "Sales"))
	if valid.Code != http.StatusOK {
		t.Fatalf("expected valid 200, got %d body=%s", valid.Code, valid.Body.String())
	}
	body := map[string]any{
		"idempotencyKey": "retry-key",
		"target": map[string]any{
			"accountInput": map[string]any{
				"companyName":    "Retry Conversion Co",
				"customerStatus": "Prospect",
				"ownerId":        "sales-1",
			},
			"opportunityInput": map[string]any{
				"ownerId":           "sales-1",
				"stage":             "New Opportunity",
				"expectedAmount":    "70000.00",
				"expectedCloseDate": "2026-11-01",
				"title":             "Retry conversion",
			},
		},
	}
	headers := actorHeaders("sales-1", "Sales")
	headers["X-Correlation-Id"] = "corr-lead-convert"
	first := postLeadJSON(app, "/leads/"+leadID+"/convert", body, headers)
	if first.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected first partial failure 503, got %d body=%s", first.Code, first.Body.String())
	}
	second := postLeadJSON(app, "/leads/"+leadID+"/convert", body, headers)
	if second.Code != http.StatusOK {
		t.Fatalf("expected retry convert 200, got %d body=%s", second.Code, second.Body.String())
	}
	if accountCreates != 1 {
		t.Fatalf("expected downstream account create to be idempotent across retry, got %d creates", accountCreates)
	}
	requireJSONValue(t, decodeJSON(t, second), "accountId", "acct_retry_1")
}

func createOwnedLead(t *testing.T, app http.Handler, companyName, ownerID string) string {
	t.Helper()
	create := postLeadJSON(app, "/leads", map[string]any{
		"companyName": companyName,
		"source":      "Website",
		"ownerId":     ownerID,
	}, actorHeaders(ownerID, "Sales"))
	if create.Code != http.StatusCreated {
		t.Fatalf("expected create 201, got %d body=%s", create.Code, create.Body.String())
	}
	return decodeJSON(t, create)["id"].(string)
}

func requireS2SRequest(t *testing.T, r *http.Request, intent string) {
	t.Helper()
	if r.Header.Get("X-Service-Id") != "lead" || r.Header.Get("X-Intent") != intent {
		t.Fatalf("missing S2S headers: service=%q intent=%q", r.Header.Get("X-Service-Id"), r.Header.Get("X-Intent"))
	}
	if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		t.Fatalf("missing bearer service token")
	}
	if r.Header.Get("X-Actor-User-Id") == "" || r.Header.Get("X-Actor-Role") == "" {
		t.Fatalf("missing actor context")
	}
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatalf("decode downstream request: %v", err)
	}
	if intent == "audit.append" {
		if body["eventId"] == "" || body["resourceId"] == "" || body["safeSummary"] == "" {
			t.Fatalf("audit append request missing history fields: %#v", body)
		}
		return
	}
	if body["ownerId"] == "" {
		t.Fatalf("downstream request must carry ownerId: %#v", body)
	}
}
