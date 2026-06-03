package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Event struct {
	EventUID           string
	EventID            string
	EventVersion       int
	ProducerService    string
	Surfaces           []string
	ActorUserID        string
	ActorRole          string
	ActorDisplay       string
	Action             string
	ResourceType       string
	ResourceID         string
	ParentResourceType string
	ParentResourceID   string
	Result             string
	ReasonCode         string
	BeforeSummary      map[string]any
	AfterSummary       map[string]any
	DiffClassification string
	RetentionPolicy    string
	RetainUntil        time.Time
	ScopeSummary       string
	SafeSummary        string
	CorrelationID      string
	CausationID        string
	AcceptanceIDs      []string
	OccurredAt         time.Time
	PrevHash           string
	EventHash          string
}

func NewEvent() Event {
	return Event{
		EventUID:   "evt_" + randomHex(16),
		OccurredAt: time.Now().UTC(),
	}
}

func (e Event) ComputeHash() (string, error) {
	payload := map[string]any{
		"eventUid":           e.EventUID,
		"eventId":            e.EventID,
		"eventVersion":       e.EventVersion,
		"producerService":    e.ProducerService,
		"surfaces":           e.Surfaces,
		"actorUserId":        e.ActorUserID,
		"actorRole":          e.ActorRole,
		"actorDisplay":       e.ActorDisplay,
		"action":             e.Action,
		"resourceType":       e.ResourceType,
		"resourceId":         e.ResourceID,
		"parentResourceType": e.ParentResourceType,
		"parentResourceId":   e.ParentResourceID,
		"result":             e.Result,
		"reasonCode":         e.ReasonCode,
		"beforeSummary":      e.BeforeSummary,
		"afterSummary":       e.AfterSummary,
		"diffClassification": e.DiffClassification,
		"retentionPolicy":    e.RetentionPolicy,
		"retainUntil":        e.RetainUntil.Format(time.RFC3339),
		"scopeSummary":       e.ScopeSummary,
		"safeSummary":        e.SafeSummary,
		"correlationId":      e.CorrelationID,
		"causationId":        e.CausationID,
		"acceptanceIds":      e.AcceptanceIDs,
		"occurredAt":         e.OccurredAt.Format(time.RFC3339Nano),
		"prevHash":           e.PrevHash,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(bytes)
	return hex.EncodeToString(sum[:]), nil
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
