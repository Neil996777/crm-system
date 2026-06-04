package authz

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestVerifyServiceTokenRejectsLongLifetime(t *testing.T) {
	now := time.Date(2026, 6, 3, 20, 45, 0, 0, time.UTC)
	secret := []byte("test-secret")
	token := signClaims(t, ServiceTokenClaims{
		Issuer:   "identity-authz",
		Audience: "work",
		Intent:   "work.command",
		Expires:  now.Add(10 * time.Minute),
	}, secret)

	if _, err := VerifyServiceToken(token, "work", "work.command", secret, now); !errors.Is(err, ErrServiceAuthFailed) {
		t.Fatalf("TEST-SVC-TOKEN-LIFETIME-001 expected ErrServiceAuthFailed for token lifetime over 5 minutes, got %v", err)
	}
}

func signClaims(t *testing.T, claims ServiceTokenClaims, secret []byte) string {
	t.Helper()
	payload, err := json.Marshal(claims)
	if err != nil {
		t.Fatal(err)
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(encodedPayload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return encodedPayload + "." + signature
}
