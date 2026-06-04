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

type VerifyOptions struct {
	Secret   []byte
	Audience string
	Intent   string
	Now      time.Time
}

func SignServiceToken(claims ServiceTokenClaims, secret []byte) (string, error) {
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	return encoded + "." + sign(encoded, secret), nil
}

func VerifyServiceToken(token string, options VerifyOptions) (ServiceTokenClaims, error) {
	if len(options.Secret) == 0 || options.Audience == "" || options.Intent == "" {
		return ServiceTokenClaims{}, ErrServiceAuthFailed
	}
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return ServiceTokenClaims{}, ErrServiceAuthFailed
	}
	expected := sign(parts[0], options.Secret)
	if !hmac.Equal([]byte(expected), []byte(parts[1])) {
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
	now := options.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if claims.Audience != options.Audience || claims.Intent != options.Intent || !claims.Expires.After(now) {
		return ServiceTokenClaims{}, ErrServiceAuthFailed
	}
	if claims.Expires.After(now.Add(5*time.Minute + time.Second)) {
		return ServiceTokenClaims{}, ErrServiceAuthFailed
	}
	return claims, nil
}

func sign(payload string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
