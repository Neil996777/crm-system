package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const HeaderCorrelationID = "X-Correlation-Id"

func CorrelationID(r *http.Request) string {
	if value := r.Header.Get(HeaderCorrelationID); value != "" {
		return value
	}
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		panic(err)
	}
	return "corr_" + hex.EncodeToString(bytes[:])
}
