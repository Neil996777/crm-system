package authz

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ErrServiceAuthFailed = errors.New("SERVICE_AUTH_FAILED")

type ServiceTokenClaims struct {
	Issuer   string    `json:"iss"`
	Audience string    `json:"aud"`
	Intent   string    `json:"intent"`
	Expires  time.Time `json:"exp"`
}

func SignServiceToken(issuer, audience, intent string, secret []byte) string {
	claims := ServiceTokenClaims{Issuer: issuer, Audience: audience, Intent: intent, Expires: time.Now().UTC().Add(2 * time.Minute)}
	payload, _ := json.Marshal(claims)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(encodedPayload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return encodedPayload + "." + signature
}

func VerifyServiceToken(token, audience, intent string, secret []byte, now time.Time) (ServiceTokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return ServiceTokenClaims{}, ErrServiceAuthFailed
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(parts[0]))
	expected := mac.Sum(nil)
	actual, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || !hmac.Equal(expected, actual) {
		return ServiceTokenClaims{}, ErrServiceAuthFailed
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return ServiceTokenClaims{}, ErrServiceAuthFailed
	}
	var claims ServiceTokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ServiceTokenClaims{}, ErrServiceAuthFailed
	}
	if claims.Audience != audience || claims.Intent != intent || !claims.Expires.After(now) {
		return ServiceTokenClaims{}, ErrServiceAuthFailed
	}
	if claims.Expires.After(now.Add(5*time.Minute + time.Second)) {
		return ServiceTokenClaims{}, ErrServiceAuthFailed
	}
	return claims, nil
}
