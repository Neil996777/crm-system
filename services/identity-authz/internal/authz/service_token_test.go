package authz

import (
	"testing"
	"time"
)

func TestServiceTokenSignedTokenVerification(t *testing.T) {
	secret := []byte("test-secret-with-enough-entropy")
	token, err := SignServiceToken(ServiceTokenClaims{
		Issuer:   "lead",
		Audience: "identity-authz",
		Intent:   "permission.check",
		Expires:  time.Now().Add(5 * time.Minute),
	}, secret)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	claims, err := VerifyServiceToken(token, VerifyOptions{
		Secret:   secret,
		Audience: "identity-authz",
		Intent:   "permission.check",
		Now:      time.Now(),
	})
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if claims.Issuer != "lead" || claims.Audience != "identity-authz" || claims.Intent != "permission.check" {
		t.Fatalf("unexpected claims: %#v", claims)
	}

	if _, err := VerifyServiceToken(token, VerifyOptions{Secret: secret, Audience: "lead", Intent: "permission.check", Now: time.Now()}); err == nil {
		t.Fatal("wrong audience token was accepted")
	}
	if _, err := VerifyServiceToken(token, VerifyOptions{Secret: secret, Audience: "identity-authz", Intent: "audit.append", Now: time.Now()}); err == nil {
		t.Fatal("wrong intent token was accepted")
	}
	if _, err := VerifyServiceToken(token+"tampered", VerifyOptions{Secret: secret, Audience: "identity-authz", Intent: "permission.check", Now: time.Now()}); err == nil {
		t.Fatal("tampered token was accepted")
	}
}
