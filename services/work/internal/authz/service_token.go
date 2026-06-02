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

type ServiceClaims struct {
	Issuer   string    `json:"iss"`
	Audience string    `json:"aud"`
	Intent   string    `json:"intent"`
	Expires  time.Time `json:"exp"`
}

func SignServiceToken(issuer, audience, intent string, secret []byte) string {
	claims := ServiceClaims{Issuer: issuer, Audience: audience, Intent: intent, Expires: time.Now().UTC().Add(2 * time.Minute)}
	payload, _ := json.Marshal(claims)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(encodedPayload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return encodedPayload + "." + signature
}

func VerifyServiceToken(token, audience, intent string, secret []byte, now time.Time) (ServiceClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return ServiceClaims{}, errors.New("invalid token")
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(parts[0]))
	expected := mac.Sum(nil)
	actual, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || !hmac.Equal(expected, actual) {
		return ServiceClaims{}, errors.New("invalid signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return ServiceClaims{}, err
	}
	var claims ServiceClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ServiceClaims{}, err
	}
	if claims.Audience != audience || claims.Intent != intent || !claims.Expires.After(now) {
		return ServiceClaims{}, errors.New("invalid claims")
	}
	return claims, nil
}
