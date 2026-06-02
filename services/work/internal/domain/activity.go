package domain

import (
	"errors"
	"strings"
	"time"
)

var ErrValidation = errors.New("validation failed")

type Activity struct {
	ID           string
	RelatedType  string
	RelatedID    string
	ActivityType string
	Content      string
	ActorID      string
	OwnerID      string
	OccurredAt   time.Time
	Version      int
	UpdatedAt    time.Time
}

func NewActivity(input Activity) (Activity, error) {
	input.RelatedType = strings.TrimSpace(input.RelatedType)
	input.RelatedID = strings.TrimSpace(input.RelatedID)
	input.ActivityType = strings.TrimSpace(input.ActivityType)
	input.Content = strings.TrimSpace(input.Content)
	input.ActorID = strings.TrimSpace(input.ActorID)
	input.OwnerID = strings.TrimSpace(input.OwnerID)
	if !validRelatedType(input.RelatedType) || input.RelatedID == "" || input.ActivityType == "" || input.Content == "" || input.ActorID == "" || input.OwnerID == "" {
		return Activity{}, ErrValidation
	}
	input.OccurredAt = time.Now().UTC()
	input.Version = 1
	return input, nil
}

func validRelatedType(value string) bool {
	switch value {
	case "Lead", "Customer", "Contact", "Opportunity", "Quote", "Contract", "Payment":
		return true
	default:
		return false
	}
}
