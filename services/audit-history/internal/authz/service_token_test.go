package authz

import (
	"errors"
	"testing"
	"time"
)

func TestVerifyServiceTokenFailsClosedOnEmptyAudienceOrIntent(t *testing.T) {
	token, err := SignServiceToken(ServiceTokenClaims{
		Issuer:   "identity-authz",
		Audience: "",
		Intent:   "",
		Expires:  time.Now().UTC().Add(2 * time.Minute),
	}, []byte("secret"))
	if err != nil {
		t.Fatalf("sign service token: %v", err)
	}
	if _, err := VerifyServiceToken(token, VerifyOptions{
		Secret:   []byte("secret"),
		Audience: "",
		Intent:   "",
		Now:      time.Now().UTC(),
	}); !errors.Is(err, ErrServiceAuthFailed) {
		t.Fatalf("TEST-SVC-TOKEN-FAILCLOSED-001 expected SERVICE_AUTH_FAILED for empty audience/intent, got %v", err)
	}
}
