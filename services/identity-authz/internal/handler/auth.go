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
	db       *sql.DB
	users    *repo.UserRepo
	sessions *repo.SessionRepo
	outbox   *event.Outbox
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
		db:       db,
		users:    repo.NewUserRepo(db),
		sessions: repo.NewSessionRepo(db),
		outbox:   event.NewOutbox(db),
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
		writeErrorCode(w, http.StatusUnauthorized, "AUTHENTICATION_FAILED", "authentication", safeAuthMessage)
		return
	}
	ctx := r.Context()
	user, err := h.users.FindByEmail(ctx, strings.TrimSpace(request.Email))
	if errors.Is(err, repo.ErrNotFound) {
		_ = bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(request.Password))
		if err := h.appendAccessDenied(ctx, auditActorAnonymous(), auditAccessDeniedOptions{ReasonCode: "login_failed"}); err != nil {
			writeDependencyUnavailable(w)
			return
		}
		writeErrorCode(w, http.StatusUnauthorized, "AUTHENTICATION_FAILED", "authentication", safeAuthMessage)
		return
	}
	if err != nil {
		writeErrorCode(w, http.StatusUnauthorized, "AUTHENTICATION_FAILED", "authentication", safeAuthMessage)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil || !user.Active() {
		if err := h.appendAccessDenied(ctx, auditActorFromUser(user), auditAccessDeniedOptions{ReasonCode: "login_failed"}); err != nil {
			writeDependencyUnavailable(w)
			return
		}
		writeErrorCode(w, http.StatusUnauthorized, "AUTHENTICATION_FAILED", "authentication", safeAuthMessage)
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
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("begin sign-in tx: %v", err)
		writeErrorCode(w, http.StatusServiceUnavailable, "DEPENDENCY_UNAVAILABLE", "dependency", "A required dependency is unavailable.")
		return
	}
	txSessions := repo.NewSessionRepoTx(tx)
	txOutbox := event.NewOutboxTx(tx)
	if err := txSessions.Create(ctx, session); err != nil {
		_ = tx.Rollback()
		log.Printf("create session: %v", err)
		writeErrorCode(w, http.StatusUnauthorized, "AUTHENTICATION_FAILED", "authentication", safeAuthMessage)
		return
	}
	if err := txOutbox.Append(ctx, event.UserSignedIn, user.ID, map[string]any{
		"actorId":      user.ID,
		"actorRole":    string(user.Role),
		"actorDisplay": user.DisplayName,
		"role":         string(user.Role),
		"result":       "success",
	}); err != nil {
		_ = tx.Rollback()
		log.Printf("append sign-in event: %v", err)
		writeErrorCode(w, http.StatusServiceUnavailable, "DEPENDENCY_UNAVAILABLE", "dependency", "A required dependency is unavailable.")
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("commit sign-in tx: %v", err)
		writeErrorCode(w, http.StatusServiceUnavailable, "DEPENDENCY_UNAVAILABLE", "dependency", "A required dependency is unavailable.")
		return
	}
	h.setSessionCookie(w, session.ID, int(h.config.SessionTTL.Seconds()))
	writeJSON(w, http.StatusOK, map[string]any{"user": userDTO(user)})
}

func (h *AuthHandler) signOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil && cookie.Value != "" {
		now := time.Now().UTC()
		tx, beginErr := h.db.BeginTx(ctx, nil)
		if beginErr != nil {
			log.Printf("begin sign-out tx: %v", beginErr)
			writeErrorCode(w, http.StatusServiceUnavailable, "DEPENDENCY_UNAVAILABLE", "dependency", "A required dependency is unavailable.")
			return
		}
		txSessions := repo.NewSessionRepoTx(tx)
		txUsers := repo.NewUserRepoTx(tx)
		txOutbox := event.NewOutboxTx(tx)
		if userID, revokeErr := txSessions.Revoke(ctx, cookie.Value, now); revokeErr == nil {
			actor := auditActorFromUserID(ctx, txUsers, userID)
			if err := txOutbox.Append(ctx, event.UserSignedOut, userID, signOutPayload(actor)); err != nil {
				_ = tx.Rollback()
				log.Printf("append sign-out event: %v", err)
				writeErrorCode(w, http.StatusServiceUnavailable, "DEPENDENCY_UNAVAILABLE", "dependency", "A required dependency is unavailable.")
				return
			}
		}
		if err := tx.Commit(); err != nil {
			log.Printf("commit sign-out tx: %v", err)
			writeErrorCode(w, http.StatusServiceUnavailable, "DEPENDENCY_UNAVAILABLE", "dependency", "A required dependency is unavailable.")
			return
		}
	}
	h.clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) currentUser(w http.ResponseWriter, r *http.Request) {
	user, sessionID, errorCode, ok := h.authenticate(r.Context(), r)
	if !ok {
		if errorCode == "DEPENDENCY_UNAVAILABLE" {
			writeDependencyUnavailable(w)
			return
		}
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
		if err := h.appendAccessDenied(ctx, auditActorAnonymous(), auditAccessDeniedOptions{ReasonCode: "unauthenticated"}); err != nil {
			return domain.User{}, "", "DEPENDENCY_UNAVAILABLE", false
		}
		return domain.User{}, "", "AUTHENTICATION_FAILED", false
	}
	session, err := h.sessions.FindByID(ctx, cookie.Value)
	if err != nil || !session.Active(time.Now().UTC()) {
		if err := h.appendAccessDenied(ctx, auditActorAnonymous(), auditAccessDeniedOptions{ReasonCode: "invalid_session"}); err != nil {
			return domain.User{}, "", "DEPENDENCY_UNAVAILABLE", false
		}
		return domain.User{}, "", "AUTHENTICATION_FAILED", false
	}
	user, err := h.users.FindByID(ctx, session.UserID)
	if err != nil || !user.Active() {
		now := time.Now().UTC()
		_, _ = h.sessions.Revoke(ctx, session.ID, now)
		actor := auditActorFromUser(user)
		if err != nil {
			actor = auditActorFromUserID(ctx, h.users, session.UserID)
		}
		if err := h.appendAccessDenied(ctx, actor, auditAccessDeniedOptions{ReasonCode: "inactive_user"}); err != nil {
			return domain.User{}, "", "DEPENDENCY_UNAVAILABLE", false
		}
		return domain.User{}, "", "AUTHENTICATION_FAILED", false
	}
	if user.AuthzVersion != session.AuthzVersionAtIssue {
		now := time.Now().UTC()
		_, _ = h.sessions.Revoke(ctx, session.ID, now)
		if err := h.appendAccessDenied(ctx, auditActorFromUser(user), auditAccessDeniedOptions{ReasonCode: "authz_version_stale"}); err != nil {
			return domain.User{}, "", "DEPENDENCY_UNAVAILABLE", false
		}
		return domain.User{}, "", "AUTHZ_VERSION_STALE", false
	}
	return user, session.ID, "", true
}

func (h *AuthHandler) appendAccessDenied(ctx context.Context, actor auditActor, options auditAccessDeniedOptions) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("begin access-denied audit tx: %v", err)
		return err
	}
	txOutbox := event.NewOutboxTx(tx)
	payload := accessDeniedPayload(actor, options)
	resourceID, _ := payload["resourceId"].(string)
	aggregateID := actor.ID
	if aggregateID == auditActorAnonymousID || aggregateID == auditActorUnknownID {
		aggregateID = resourceID
	}
	aggregateID = firstAuditNonEmpty(aggregateID, resourceID, auditAuthResourceID)
	if err := txOutbox.Append(ctx, event.UserAccessDenied, aggregateID, payload); err != nil {
		_ = tx.Rollback()
		log.Printf("append access-denied event: %v", err)
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Printf("commit access-denied audit tx: %v", err)
		return err
	}
	return nil
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
