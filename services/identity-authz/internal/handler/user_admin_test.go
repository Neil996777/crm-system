package handler

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestUserAdminGuardsAndPermissions(t *testing.T) {
	db := newAuthTestDB(t)
	app := NewAuthServer(db, Config{CookieSecure: true, SessionTTL: 12 * time.Hour, IdleSessionTTL: 30 * time.Minute})

	adminEmail := "admin-useradmin-" + randomSuffix(t) + "@example.com"
	adminID := insertUser(t, db, adminEmail, "pw", "Administrator", "Active")
	salesEmail := "sales-useradmin-" + randomSuffix(t) + "@example.com"
	insertUser(t, db, salesEmail, "pw", "Sales", "Active")

	adminCookie := requireSessionCookie(t, postJSON(app, "/auth/sign-in", map[string]string{"email": adminEmail, "password": "pw"}, nil))
	salesCookie := requireSessionCookie(t, postJSON(app, "/auth/sign-in", map[string]string{"email": salesEmail, "password": "pw"}, nil))

	t.Run("TEST-USER-ADMIN-001 and TEST-PERM-USERADMIN-001 admin creates user with one role", func(t *testing.T) {
		rec := postJSON(app, "/admin/users", map[string]any{
			"email":       "created-" + randomSuffix(t) + "@example.com",
			"displayName": "Created User",
			"password":    "pw",
			"role":        "Sales Manager",
		}, adminCookie)
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, rec), "user.role", "Sales Manager")
		requireEvent(t, db, "UserRoleStatusChanged", "")
	})

	t.Run("TEST-PERM-USERADMIN-002/003 and TEST-ABUSE-PRIVESC-001 non-admin cannot manage users", func(t *testing.T) {
		rec := postJSON(app, "/admin/users", map[string]any{
			"email":       "forbidden-" + randomSuffix(t) + "@example.com",
			"displayName": "Forbidden User",
			"password":    "pw",
			"role":        "Administrator",
		}, salesCookie)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "PERMISSION_DENIED")
		if strings.Contains(rec.Body.String(), "Forbidden User") {
			t.Fatal("denial leaked requested user details")
		}
	})

	t.Run("TEST-USER-ADMIN-002/003/004 role and status changes persist", func(t *testing.T) {
		targetID := insertUser(t, db, "target-"+randomSuffix(t)+"@example.com", "pw", "Sales", "Active")
		role := patchJSON(app, "/admin/users/"+targetID+"/role", map[string]string{"role": "Sales Manager"}, adminCookie)
		if role.Code != http.StatusOK {
			t.Fatalf("expected role change 200, got %d body=%s", role.Code, role.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, role), "user.role", "Sales Manager")

		disabled := patchJSON(app, "/admin/users/"+targetID+"/status", map[string]string{"status": "Disabled"}, adminCookie)
		if disabled.Code != http.StatusOK {
			t.Fatalf("expected disable 200, got %d body=%s", disabled.Code, disabled.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, disabled), "user.status", "Disabled")

		enabled := patchJSON(app, "/admin/users/"+targetID+"/status", map[string]string{"status": "Active"}, adminCookie)
		if enabled.Code != http.StatusOK {
			t.Fatalf("expected enable 200, got %d body=%s", enabled.Code, enabled.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, enabled), "user.status", "Active")
	})

	t.Run("TEST-INV-LASTADMIN-001 last active Administrator cannot be disabled or downgraded", func(t *testing.T) {
		if _, err := db.Exec(`UPDATE identity_authz.users SET status = 'Disabled' WHERE role_name = 'Administrator' AND id <> $1`, adminID); err != nil {
			t.Fatalf("prepare one active admin: %v", err)
		}
		downgrade := patchJSON(app, "/admin/users/"+adminID+"/role", map[string]string{"role": "Sales"}, adminCookie)
		if downgrade.Code != http.StatusConflict {
			t.Fatalf("expected downgrade conflict, got %d body=%s", downgrade.Code, downgrade.Body.String())
		}
		requireErrorCode(t, downgrade, "LAST_ADMIN_BLOCKED")

		disable := patchJSON(app, "/admin/users/"+adminID+"/status", map[string]string{"status": "Disabled"}, adminCookie)
		if disable.Code != http.StatusConflict {
			t.Fatalf("expected disable conflict, got %d body=%s", disable.Code, disable.Body.String())
		}
		requireErrorCode(t, disable, "LAST_ADMIN_BLOCKED")
		requireEvent(t, db, "LastAdministratorBlocked", adminID)
	})

	t.Run("TEST-INV-NODELETE-001 no hard-delete route exists", func(t *testing.T) {
		rec := httptestDelete(app, "/admin/users/"+adminID, adminCookie)
		if rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected no delete route, got %d body=%s", rec.Code, rec.Body.String())
		}
	})
}
