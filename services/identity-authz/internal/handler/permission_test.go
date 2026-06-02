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

	t.Run("TEST-AUTHZ-SCOPE-001 admin governs all", func(t *testing.T) {
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

	t.Run("TEST-AUTHZ-SCOPE-002 manager team scope", func(t *testing.T) {
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

	t.Run("TEST-AUTHZ-SCOPE-003 and 004 sales owned allowed and cross-owner denied without leakage", func(t *testing.T) {
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
