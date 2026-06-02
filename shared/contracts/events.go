package contracts

import "time"

type EventEnvelope struct {
	EventID         string    `json:"eventId"`
	EventType       string    `json:"eventType"`
	EventVersion    int       `json:"eventVersion"`
	ProducerService string    `json:"producerService"`
	AggregateType   string    `json:"aggregateType"`
	AggregateID     string    `json:"aggregateId"`
	ActorID         string    `json:"actorId"`
	OccurredAt      time.Time `json:"occurredAt"`
	CorrelationID   string    `json:"correlationId"`
	CausationID     string    `json:"causationId,omitempty"`
	SafeSummary     string    `json:"safeSummary"`
}
