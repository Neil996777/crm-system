package contracts

type ErrorBody struct {
	Code        string `json:"code"`
	Category    string `json:"category"`
	SafeMessage string `json:"safeMessage"`
}

type ResponseEnvelope[T any] struct {
	CorrelationID string     `json:"correlationId"`
	Data          *T         `json:"data,omitempty"`
	Error         *ErrorBody `json:"error,omitempty"`
}
