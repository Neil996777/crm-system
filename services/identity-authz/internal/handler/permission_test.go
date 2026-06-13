package handler

import (
	"net/http"
	"testing"
	"time"

	"crm-system/services/identity-authz/internal/authz"
)

func TestPermissionScopeAndS2SAcceptance(t *testing.T) {
	db := newAuthTestDB(t)
	secret := []byte("permission-test-secret")
	app := NewAuthServer(db, Config{
		CookieSecure:       true,
		SessionTTL:         12 * time.Hour,
		IdleSessionTTL:     30 * time.Minute,
		ServiceID:          "identity-authz",
		ServiceTokenSecret: secret,
	})

	adminID := insertUser(t, db, "perm-admin-"+randomSuffix(t)+"@example.com", "pw", "Administrator", "Active")
	managerID := insertUser(t, db, "perm-manager-"+randomSuffix(t)+"@example.com", "pw", "Sales Manager", "Active")
	salesID := insertUser(t, db, "perm-sales-"+randomSuffix(t)+"@example.com", "pw", "Sales", "Active")

	token := signTestServiceToken(t, secret, "lead", "identity-authz", "permission.check")

	t.Run("TEST-AUTHZ-SCOPE-001 and TEST-PERM-CRUD-ADMIN-001 admin governs all", func(t *testing.T) {
		rec := permissionCheck(app, token, map[string]any{
			"actorId": adminID,
			"action":  "lead.update",
			"resource": map[string]any{
				"type": "lead",
				"id":   "lead-1",
			},
			"context": map[string]any{
				"ownerId": "another-user",
				"teamId":  "single-team",
			},
		})
		requirePermission(t, rec, true, "all", nil)
	})

	t.Run("TEST-AUTHZ-SCOPE-002 and TEST-PERM-CRUD-MGR-001 manager team scope", func(t *testing.T) {
		rec := permissionCheck(app, token, map[string]any{
			"actorId": managerID,
			"action":  "opportunity.update",
			"resource": map[string]any{
				"type": "opportunity",
				"id":   "opp-1",
			},
			"context": map[string]any{
				"ownerId": "another-user",
				"teamId":  "single-team",
			},
		})
		requirePermission(t, rec, true, "team", nil)
	})

	t.Run("TEST-AUTHZ-SCOPE-003/004/005 and TEST-PERM-CRUD-SALES-001 sales owned allowed and cross-owner denied without leakage or mutation", func(t *testing.T) {
		owned := permissionCheck(app, token, map[string]any{
			"actorId": salesID,
			"action":  "lead.update",
			"resource": map[string]any{
				"type": "lead",
				"id":   "lead-owned",
			},
			"context": map[string]any{
				"ownerId": salesID,
				"teamId":  "single-team",
			},
		})
		requirePermission(t, owned, true, "owned", nil)

		denied := permissionCheck(app, token, map[string]any{
			"actorId": salesID,
			"action":  "lead.update",
			"resource": map[string]any{
				"type": "lead",
				"id":   "restricted-lead-name-must-not-leak",
			},
			"context": map[string]any{
				"ownerId": "another-user",
				"teamId":  "single-team",
			},
		})
		requirePermission(t, denied, false, "", "scope_denied")
		if denied.Body.String() == "restricted-lead-name-must-not-leak" {
			t.Fatal("permission denial leaked restricted resource id")
		}
		_, payload := latestOutboxPayload(t, db, "UserAccessDenied", salesID)
		requirePayloadString(t, payload, "actorId", salesID)
		requirePayloadString(t, payload, "actorRole", "Sales")
		requirePayloadString(t, payload, "reasonCode", "scope_denied")
		requirePayloadString(t, payload, "resourceType", "lead")
		requirePayloadString(t, payload, "resourceId", "permission-check")
		requirePayloadString(t, payload, "result", "denied")
		requirePayloadEmptySummary(t, payload, "beforeSummary")
		requirePayloadEmptySummary(t, payload, "afterSummary")
	})

	t.Run("TEST-AUTHZ-SCOPE-006 and TEST-INV-NODELETE-001 hard delete denied for every role", func(t *testing.T) {
		for _, actorID := range []string{adminID, managerID, salesID} {
			rec := permissionCheck(app, token, map[string]any{
				"actorId": actorID,
				"action":  "lead.hard_delete",
				"resource": map[string]any{
					"type": "lead",
					"id":   "lead-1",
				},
				"context": map[string]any{"ownerId": actorID, "teamId": "single-team"},
			})
			requirePermission(t, rec, false, "", "hard_delete_forbidden")
		}
	})

	t.Run("TEST-ABUSE-S2S-001 invalid service token returns SERVICE_AUTH_FAILED", func(t *testing.T) {
		rec := permissionCheck(app, "invalid-token", map[string]any{
			"actorId": adminID,
			"action":  "lead.update",
		})
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "SERVICE_AUTH_FAILED")
	})
}

func TestPermissionDeniedAuditOutboxFailureIsSurfaced(t *testing.T) {
	db := newAuthTestDB(t)
	secret := []byte("permission-test-secret")
	app := NewAuthServer(db, Config{
		CookieSecure:       true,
		SessionTTL:         12 * time.Hour,
		IdleSessionTTL:     30 * time.Minute,
		ServiceID:          "identity-authz",
		ServiceTokenSecret: secret,
	})
	salesID := insertUser(t, db, "perm-denial-audit-"+randomSuffix(t)+"@example.com", "pw", "Sales", "Active")
	token := signTestServiceToken(t, secret, "lead", "identity-authz", "permission.check")
	if _, err := db.Exec(`DROP TABLE identity_authz.outbox_events`); err != nil {
		t.Fatalf("force outbox failure: %v", err)
	}

	rec := permissionCheck(app, token, map[string]any{
		"actorId": salesID,
		"action":  "lead.update",
		"resource": map[string]any{
			"type": "lead",
			"id":   "restricted-lead",
		},
		"context": map[string]any{
			"ownerId": "another-user",
			"teamId":  "single-team",
		},
	})

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("TEST-AUDIT-DENIAL-DURABLE-002 expected audit persistence failure to surface as 503, got %d body=%s", rec.Code, rec.Body.String())
	}
	requireErrorCode(t, rec, "DEPENDENCY_UNAVAILABLE")
}

func signTestServiceToken(t *testing.T, secret []byte, issuer, audience, intent string) string {
	t.Helper()
	token, err := authz.SignServiceToken(authz.ServiceTokenClaims{
		Issuer:   issuer,
		Audience: audience,
		Intent:   intent,
		Expires:  time.Now().Add(5 * time.Minute),
	}, secret)
	if err != nil {
		t.Fatalf("sign service token: %v", err)
	}
	return token
}
