package domain

import (
	"strings"
	"time"
)

type Note struct {
	ID          string
	RelatedType string
	RelatedID   string
	Content     string
	ActorID     string
	OwnerID     string
	OccurredAt  time.Time
	Version     int
	UpdatedAt   time.Time
}

func NewNote(input Note) (Note, error) {
	input.RelatedType = strings.TrimSpace(input.RelatedType)
	input.RelatedID = strings.TrimSpace(input.RelatedID)
	input.Content = strings.TrimSpace(input.Content)
	input.ActorID = strings.TrimSpace(input.ActorID)
	input.OwnerID = strings.TrimSpace(input.OwnerID)
	if !validRelatedType(input.RelatedType) || input.RelatedID == "" || input.Content == "" || input.ActorID == "" || input.OwnerID == "" {
		return Note{}, ErrValidation
	}
	input.OccurredAt = time.Now().UTC()
	input.Version = 1
	return input, nil
}
