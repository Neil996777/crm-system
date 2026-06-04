package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
)

const authFailureMessage = "Authentication failed."

func TestAuthSessionAcceptance(t *testing.T) {
	db := newAuthTestDB(t)
	app := NewAuthServer(db, Config{
		CookieSecure:   true,
		SessionTTL:     12 * time.Hour,
		IdleSessionTTL: 30 * time.Minute,
	})

	t.Run("TEST-AUTH-LOGIN-001 valid login binds role and emits login event", func(t *testing.T) {
		email := "admin-" + randomSuffix(t) + "@example.com"
		userID := insertUser(t, db, email, "correct-password", "Administrator", "Active")

		rec := postJSON(app, "/auth/sign-in", map[string]string{
			"email":    email,
			"password": "correct-password",
		}, nil)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		cookie := requireSessionCookie(t, rec)
		if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode {
			t.Fatalf("session cookie must be HttpOnly/Secure/SameSite=Lax, got %#v", cookie)
		}
		body := decodeJSON(t, rec)
		requireJSONValue(t, body, "user.id", userID)
		requireJSONValue(t, body, "user.role", "Administrator")
		requireEvent(t, db, "UserSignedIn", userID)
	})

	t.Run("TEST-AUTH-LOGIN-002 and TEST-ABUSE-ENUM-001 invalid credentials use one safe message", func(t *testing.T) {
		email := "enum-" + randomSuffix(t) + "@example.com"
		insertUser(t, db, email, "correct-password", "Sales", "Active")

		invalid := signInFailure(t, app, email, "wrong-password")
		missing := signInFailure(t, app, "missing-"+randomSuffix(t)+"@example.com", "wrong-password")

		if invalid != authFailureMessage || missing != authFailureMessage || invalid != missing {
			t.Fatalf("expected unified safe failure message, invalid=%q missing=%q", invalid, missing)
		}
		invalidBody := postJSON(app, "/auth/sign-in", map[string]string{"email": email, "password": "wrong-password"}, nil)
		requireErrorCode(t, invalidBody, "AUTHENTICATION_FAILED")
		requireErrorCategory(t, invalidBody, "authentication")
		requireEvent(t, db, "UserAccessDenied", "")
	})

	t.Run("TEST-AUTH-LOGIN-003 disabled user denied with the same safe message", func(t *testing.T) {
		email := "disabled-" + randomSuffix(t) + "@example.com"
		insertUser(t, db, email, "correct-password", "Sales Manager", "Disabled")

		message := signInFailure(t, app, email, "correct-password")

		if message != authFailureMessage {
			t.Fatalf("expected unified safe failure message, got %q", message)
		}
		requireEvent(t, db, "UserAccessDenied", "")
	})

	t.Run("TEST-AUTH-LOGIN-004 and TEST-ABUSE-UNAUTH-001 unauthenticated protected API denied and logged", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/current", nil)
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorMessage(t, rec, authFailureMessage)
		requireEvent(t, db, "UserAccessDenied", "")
	})

	t.Run("TEST-AUTH-LOGIN-005 session persists across refresh and handler restart; sign-out revokes it", func(t *testing.T) {
		email := "persist-" + randomSuffix(t) + "@example.com"
		userID := insertUser(t, db, email, "correct-password", "Sales", "Active")
		login := postJSON(app, "/auth/sign-in", map[string]string{
			"email":    email,
			"password": "correct-password",
		}, nil)
		cookie := requireSessionCookie(t, login)

		current := get(app, "/auth/current", cookie)
		if current.Code != http.StatusOK {
			t.Fatalf("expected session to survive refresh, got %d body=%s", current.Code, current.Body.String())
		}
		requireJSONValue(t, decodeJSON(t, current), "user.id", userID)

		restartedApp := NewAuthServer(db, Config{CookieSecure: true, SessionTTL: 12 * time.Hour, IdleSessionTTL: 30 * time.Minute})
		currentAfterRestart := get(restartedApp, "/auth/current", cookie)
		if currentAfterRestart.Code != http.StatusOK {
			t.Fatalf("expected persisted session after handler restart, got %d body=%s", currentAfterRestart.Code, currentAfterRestart.Body.String())
		}

		logout := postJSON(app, "/auth/sign-out", nil, cookie)
		if logout.Code != http.StatusNoContent {
			t.Fatalf("expected sign-out 204, got %d body=%s", logout.Code, logout.Body.String())
		}
		cleared := requireSessionCookie(t, logout)
		if cleared.MaxAge >= 0 {
			t.Fatalf("expected sign-out to clear cookie, got MaxAge=%d", cleared.MaxAge)
		}
		afterLogout := get(app, "/auth/current", cookie)
		if afterLogout.Code != http.StatusUnauthorized {
			t.Fatalf("expected old cookie to be revoked server-side, got %d body=%s", afterLogout.Code, afterLogout.Body.String())
		}
		requireEvent(t, db, "UserSignedOut", userID)
	})

	t.Run("TEST-AUTH-LOGIN-006 stale session re-evaluates role and disabled status with AUTHZ_VERSION_STALE", func(t *testing.T) {
		email := "stale-" + randomSuffix(t) + "@example.com"
		userID := insertUser(t, db, email, "correct-password", "Sales", "Active")
		login := postJSON(app, "/auth/sign-in", map[string]string{
			"email":    email,
			"password": "correct-password",
		}, nil)
		cookie := requireSessionCookie(t, login)

		if _, err := db.Exec(`UPDATE identity_authz.users SET role_name = 'Sales Manager', authz_version = authz_version + 1 WHERE id = $1`, userID); err != nil {
			t.Fatalf("change role: %v", err)
		}
		current := get(app, "/auth/current", cookie)
		if current.Code != http.StatusUnauthorized {
			t.Fatalf("expected stale authz version to be denied, got %d body=%s", current.Code, current.Body.String())
		}
		requireErrorCode(t, current, "AUTHZ_VERSION_STALE")

		if _, err := db.Exec(`UPDATE identity_authz.users SET status = 'Disabled', authz_version = authz_version + 1 WHERE id = $1`, userID); err != nil {
			t.Fatalf("disable user: %v", err)
		}
		denied := get(app, "/auth/current", cookie)
		if denied.Code != http.StatusUnauthorized {
			t.Fatalf("expected disabled stale session to be denied, got %d body=%s", denied.Code, denied.Body.String())
		}
		requireErrorMessage(t, denied, authFailureMessage)
		requireEvent(t, db, "UserAccessDenied", userID)
	})
}

func newAuthTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "crm_system",
				"POSTGRES_USER":     "crm_admin",
				"POSTGRES_PASSWORD": "crm_admin_dev_password",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("start postgres testcontainer: %v", err)
	}
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Fatalf("terminate postgres testcontainer: %v", err)
		}
	})
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("container port: %v", err)
	}
	dsn := fmt.Sprintf("postgres://crm_admin:crm_admin_dev_password@%s:%s/crm_system?sslmode=disable", host, port.Port())
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close db: %v", err)
		}
	})
	for _, migration := range []string{"0001_init_schema.up.sql", "0002_users_sessions.up.sql", "0003_permission_policy.up.sql"} {
		sqlBytes, err := os.ReadFile(filepath.Join("..", "..", "migrations", migration))
		if err != nil {
			t.Fatalf("read migration %s: %v", migration, err)
		}
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			t.Fatalf("apply migration %s: %v", migration, err)
		}
	}
	return db
}

func insertUser(t *testing.T, db *sql.DB, email, password, role, status string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	id := "usr_" + randomSuffix(t)
	if _, err := db.Exec(`
		INSERT INTO identity_authz.users (id, email, display_name, password_hash, role_name, status, authz_version)
		VALUES ($1, $2, $3, $4, $5, $6, 1)
	`, id, email, email, string(hash), role, status); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	return id
}

func postJSON(handler http.Handler, path string, body any, cookie *http.Cookie) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
			panic(err)
		}
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, &requestBody)
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	handler.ServeHTTP(rec, req)
	return rec
}

func get(handler http.Handler, path string, cookie *http.Cookie) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	handler.ServeHTTP(rec, req)
	return rec
}

func patchJSON(handler http.Handler, path string, body any, cookie *http.Cookie) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
			panic(err)
		}
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, path, &requestBody)
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	handler.ServeHTTP(rec, req)
	return rec
}

func httptestDelete(handler http.Handler, path string, cookie *http.Cookie) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	handler.ServeHTTP(rec, req)
	return rec
}

func permissionCheck(handler http.Handler, token string, body any) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
			panic(err)
		}
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/internal/permissions/check", &requestBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", "lead")
	req.Header.Set("X-Intent", "permission.check")
	handler.ServeHTTP(rec, req)
	return rec
}

func signInFailure(t *testing.T, handler http.Handler, email, password string) string {
	t.Helper()
	rec := postJSON(handler, "/auth/sign-in", map[string]string{"email": email, "password": password}, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
	}
	return errorMessage(t, rec)
}

func requireSessionCookie(t *testing.T, rec *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "crm_session" {
			return cookie
		}
	}
	t.Fatalf("missing crm_session cookie in %v", rec.Result().Cookies())
	return nil
}

func requireErrorMessage(t *testing.T, rec *httptest.ResponseRecorder, expected string) {
	t.Helper()
	if got := errorMessage(t, rec); got != expected {
		t.Fatalf("expected error message %q, got %q body=%s", expected, got, rec.Body.String())
	}
}

func requireErrorCode(t *testing.T, rec *httptest.ResponseRecorder, expected string) {
	t.Helper()
	body := decodeJSON(t, rec)
	errBody, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("missing error body: %#v", body)
	}
	code, ok := errBody["code"].(string)
	if !ok {
		t.Fatalf("missing error code: %#v", errBody)
	}
	if code != expected {
		t.Fatalf("expected error code %q, got %q body=%s", expected, code, rec.Body.String())
	}
}

func requireErrorCategory(t *testing.T, rec *httptest.ResponseRecorder, expected string) {
	t.Helper()
	body := decodeJSON(t, rec)
	errBody, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("missing error body: %#v", body)
	}
	category, ok := errBody["category"].(string)
	if !ok {
		t.Fatalf("missing error category: %#v", errBody)
	}
	if category != expected {
		t.Fatalf("expected error category %q, got %q body=%s", expected, category, rec.Body.String())
	}
}

func requirePermission(t *testing.T, rec *httptest.ResponseRecorder, allowed bool, scope string, denial any) {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Fatalf("expected permission check 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	body := decodeJSON(t, rec)
	if body["allowed"] != allowed {
		t.Fatalf("expected allowed=%v, got %#v body=%#v", allowed, body["allowed"], body)
	}
	if body["scope"] != scope {
		t.Fatalf("expected scope=%q, got %#v body=%#v", scope, body["scope"], body)
	}
	if denial == nil {
		if body["denialCategory"] != nil {
			t.Fatalf("expected nil denial, got %#v body=%#v", body["denialCategory"], body)
		}
		return
	}
	if body["denialCategory"] != denial {
		t.Fatalf("expected denial=%#v, got %#v body=%#v", denial, body["denialCategory"], body)
	}
}

func errorMessage(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	body := decodeJSON(t, rec)
	errBody, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("missing error body: %#v", body)
	}
	message, ok := errBody["safeMessage"].(string)
	if !ok {
		t.Fatalf("missing safeMessage: %#v", errBody)
	}
	return message
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode json body %q: %v", rec.Body.String(), err)
	}
	return body
}

func requireJSONValue(t *testing.T, body map[string]any, path, expected string) {
	t.Helper()
	current := any(body)
	for _, part := range []string{path[:4], path[5:]} {
		next, ok := current.(map[string]any)[part]
		if !ok {
			t.Fatalf("missing json path %s in %#v", path, body)
		}
		current = next
	}
	if current != expected {
		t.Fatalf("expected %s=%q, got %#v in %#v", path, expected, current, body)
	}
}

func requireEvent(t *testing.T, db *sql.DB, eventType, aggregateID string) {
	t.Helper()
	var count int
	query := `SELECT count(*) FROM identity_authz.outbox_events WHERE event_type = $1`
	args := []any{eventType}
	if aggregateID != "" {
		query += ` AND aggregate_id = $2`
		args = append(args, aggregateID)
	}
	if err := db.QueryRow(query, args...).Scan(&count); err != nil {
		t.Fatalf("count event %s: %v", eventType, err)
	}
	if count == 0 {
		t.Fatalf("expected outbox event %s aggregate=%s", eventType, aggregateID)
	}
}

func randomSuffix(t *testing.T) string {
	t.Helper()
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		t.Fatalf("random suffix: %v", err)
	}
	return hex.EncodeToString(bytes[:])
}
