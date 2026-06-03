package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

var ErrServiceAuthFailed = errors.New("SERVICE_AUTH_FAILED")

type ServiceTokenClaims struct {
	Issuer   string    `json:"iss"`
	Audience string    `json:"aud"`
	Intent   string    `json:"intent"`
	Expires  time.Time `json:"exp"`
}

func SignServiceToken(claims ServiceTokenClaims, secret []byte) (string, error) {
	if len(secret) == 0 || claims.Issuer == "" || claims.Audience == "" || claims.Intent == "" || claims.Expires.IsZero() {
		return "", ErrServiceAuthFailed
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(encodedPayload))
	return encodedPayload + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}
