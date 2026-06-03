package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"crm-system/services/identity-authz/internal/authz"
	"crm-system/services/identity-authz/internal/domain"
	"crm-system/services/identity-authz/internal/event"
	"crm-system/services/identity-authz/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionCookieName = "crm_session"
	safeAuthMessage   = "Authentication failed."
)

var dummyPasswordHash = []byte("$2a$04$Zvn.Hww6XRH/FSQyBPwbpOqCu6E0uPqfiNGiM5dO9S0UDYtBTbTsO")

type Config struct {
	CookieSecure           bool
	SessionTTL             time.Duration
	IdleSessionTTL         time.Duration
	AuditHistoryServiceURL string
	ServiceID              string
	ServiceTokenSecret     []byte
}

type AuthHandler struct {
	users    *repo.UserRepo
	sessions *repo.SessionRepo
	outbox   *event.Outbox
	audit    authz.AuditClient
	config   Config
}

func NewAuthServer(db *sql.DB, config Config) http.Handler {
	if config.SessionTTL == 0 {
		config.SessionTTL = 12 * time.Hour
	}
	if config.IdleSessionTTL == 0 {
		config.IdleSessionTTL = 30 * time.Minute
	}
	if config.ServiceID == "" {
		config.ServiceID = "identity-authz"
	}
	handler := &AuthHandler{
		users:    repo.NewUserRepo(db),
		sessions: repo.NewSessionRepo(db),
		outbox:   event.NewOutbox(db),
		audit:    authz.NewAuditClient(config.AuditHistoryServiceURL, config.ServiceID, config.ServiceTokenSecret, nil),
		config:   config,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/sign-in", handler.signIn)
	mux.HandleFunc("POST /auth/sign-out", handler.signOut)
	mux.HandleFunc("GET /auth/current", handler.currentUser)
	mux.HandleFunc("GET /internal/sessions/check", handler.currentUser)
	mux.HandleFunc("POST /internal/permissions/check", handler.permissionCheck)
	mux.HandleFunc("GET /admin/users", handler.listUsers)
	mux.HandleFunc("POST /admin/users", handler.createUser)
	mux.HandleFunc("PATCH /admin/users/{id}/role", handler.changeUserRole)
	mux.HandleFunc("PATCH /admin/users/{id}/status", handler.changeUserStatus)
	return mux
}

func (h *AuthHandler) verifyServiceToken(r *http.Request, intent string) bool {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}
	claims, err := authz.VerifyServiceToken(strings.TrimPrefix(authHeader, "Bearer "), authz.VerifyOptions{
		Secret:   h.config.ServiceTokenSecret,
		Audience: h.config.ServiceID,
		Intent:   intent,
		Now:      time.Now().UTC(),
	})
	if err != nil {
		return false
	}
	serviceID := r.Header.Get("X-Service-Id")
	requestIntent := r.Header.Get("X-Intent")
	return serviceID != "" && serviceID == claims.Issuer && requestIntent == intent
}

func (h *AuthHandler) signIn(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusUnauthorized, safeAuthMessage)
		return
	}
	ctx := r.Context()
	user, err := h.users.FindByEmail(ctx, strings.TrimSpace(request.Email))
	if errors.Is(err, repo.ErrNotFound) {
		_ = bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(request.Password))
		h.appendAccessDenied(ctx, "", "login_failed")
		writeError(w, http.StatusUnauthorized, safeAuthMessage)
		return
	}
	if err != nil {
		writeError(w, http.StatusUnauthorized, safeAuthMessage)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil || !user.Active() {
		h.appendAccessDenied(ctx, user.ID, "login_failed")
		writeError(w, http.StatusUnauthorized, safeAuthMessage)
		return
	}

	now := time.Now().UTC()
	session := domain.Session{
		ID:                  "ses_" + randomHex(32),
		UserID:              user.ID,
		AuthzVersionAtIssue: user.AuthzVersion,
		ExpiresAt:           now.Add(h.config.SessionTTL),
		IdleExpiresAt:       now.Add(h.config.IdleSessionTTL),
		CreatedAt:           now,
		LastSeenAt:          now,
	}
	if err := h.sessions.Create(ctx, session); err != nil {
		log.Printf("create session: %v", err)
		writeError(w, http.StatusUnauthorized, safeAuthMessage)
		return
	}
	h.setSessionCookie(w, session.ID, int(h.config.SessionTTL.Seconds()))
	if err := h.outbox.Append(ctx, event.UserSignedIn, user.ID, map[string]any{
		"actorId": user.ID,
		"role":    string(user.Role),
		"result":  "success",
	}); err != nil {
		log.Printf("append sign-in event: %v", err)
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": userDTO(user)})
}

func (h *AuthHandler) signOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil && cookie.Value != "" {
		now := time.Now().UTC()
		if userID, revokeErr := h.sessions.Revoke(ctx, cookie.Value, now); revokeErr == nil {
			if err := h.outbox.Append(ctx, event.UserSignedOut, userID, map[string]any{
				"actorId": userID,
				"result":  "success",
			}); err != nil {
				log.Printf("append sign-out event: %v", err)
			}
		}
	}
	h.clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) currentUser(w http.ResponseWriter, r *http.Request) {
	user, sessionID, errorCode, ok := h.authenticate(r.Context(), r)
	if !ok {
		if errorCode == "" {
			errorCode = "AUTHENTICATION_FAILED"
		}
		writeErrorCode(w, http.StatusUnauthorized, errorCode, "authentication", safeAuthMessage)
		return
	}
	now := time.Now().UTC()
	if err := h.sessions.Touch(r.Context(), sessionID, now, now.Add(h.config.IdleSessionTTL)); err != nil {
		log.Printf("touch session: %v", err)
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": userDTO(user)})
}

func (h *AuthHandler) authenticate(ctx context.Context, r *http.Request) (domain.User, string, string, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		h.appendAccessDenied(ctx, "", "unauthenticated")
		return domain.User{}, "", "AUTHENTICATION_FAILED", false
	}
	session, err := h.sessions.FindByID(ctx, cookie.Value)
	if err != nil || !session.Active(time.Now().UTC()) {
		h.appendAccessDenied(ctx, "", "invalid_session")
		return domain.User{}, "", "AUTHENTICATION_FAILED", false
	}
	user, err := h.users.FindByID(ctx, session.UserID)
	if err != nil || !user.Active() {
		now := time.Now().UTC()
		_, _ = h.sessions.Revoke(ctx, session.ID, now)
		h.appendAccessDenied(ctx, session.UserID, "inactive_user")
		return domain.User{}, "", "AUTHENTICATION_FAILED", false
	}
	if user.AuthzVersion != session.AuthzVersionAtIssue {
		now := time.Now().UTC()
		_, _ = h.sessions.Revoke(ctx, session.ID, now)
		h.appendAccessDenied(ctx, session.UserID, "authz_version_stale")
		return domain.User{}, "", "AUTHZ_VERSION_STALE", false
	}
	return user, session.ID, "", true
}

func (h *AuthHandler) appendAccessDenied(ctx context.Context, userID, reason string) {
	if err := h.outbox.Append(ctx, event.UserAccessDenied, userID, map[string]any{
		"actorId": userID,
		"reason":  reason,
		"result":  "denied",
	}); err != nil {
		log.Printf("append access-denied event: %v", err)
	}
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	})
}

func (h *AuthHandler) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func userDTO(user domain.User) map[string]any {
	return map[string]any{
		"id":          user.ID,
		"email":       user.Email,
		"displayName": user.DisplayName,
		"role":        string(user.Role),
		"status":      string(user.Status),
	}
}

func writeError(w http.ResponseWriter, status int, safeMessage string) {
	writeErrorCode(w, status, "AUTHENTICATION_FAILED", "authentication", safeMessage)
}

func writeErrorCode(w http.ResponseWriter, status int, code, category, safeMessage string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":        code,
			"category":    category,
			"safeMessage": safeMessage,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write json: %v", err)
	}
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
