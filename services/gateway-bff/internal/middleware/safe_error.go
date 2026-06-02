package middleware

import "encoding/json"

type ErrorEnvelope struct {
	Code        string       `json:"code"`
	Category    string       `json:"category"`
	SafeMessage string       `json:"safeMessage"`
	FieldErrors []FieldError `json:"fieldErrors,omitempty"`
}

type FieldError struct {
	Field       string `json:"field"`
	Code        string `json:"code"`
	SafeMessage string `json:"safeMessage"`
}

func NormalizeError(body []byte, fallbackCode, fallbackCategory, fallbackMessage string) ErrorEnvelope {
	var parsed struct {
		Error ErrorEnvelope `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err == nil && parsed.Error.Code != "" {
		return parsed.Error
	}
	return ErrorEnvelope{Code: fallbackCode, Category: fallbackCategory, SafeMessage: fallbackMessage}
}
