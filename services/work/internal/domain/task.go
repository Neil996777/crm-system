package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	TaskStatusOpen      = "Open"
	TaskStatusCompleted = "Completed"
	TaskStatusCancelled = "Cancelled"
	TaskStatusOverdue   = "Overdue"
)

var ErrInvalidTaskTransition = errors.New("invalid task transition")

type Task struct {
	ID                 string
	RelatedType        string
	RelatedID          string
	Title              string
	DueDate            time.Time
	Status             string
	ActorID            string
	OwnerID            string
	CompletedAt        time.Time
	CancelledAt        time.Time
	CancellationReason string
	Version            int
	UpdatedAt          time.Time
}

func NewTask(input Task) (Task, error) {
	input.RelatedType = strings.TrimSpace(input.RelatedType)
	input.RelatedID = strings.TrimSpace(input.RelatedID)
	input.Title = strings.TrimSpace(input.Title)
	input.ActorID = strings.TrimSpace(input.ActorID)
	input.OwnerID = strings.TrimSpace(input.OwnerID)
	if !validRelatedType(input.RelatedType) || input.RelatedID == "" || input.Title == "" || input.DueDate.IsZero() || input.ActorID == "" || input.OwnerID == "" {
		return Task{}, ErrValidation
	}
	input.Status = TaskStatusOpen
	input.Version = 1
	return input, nil
}

func ApplyTaskStatus(current Task, toStatus, reason string) (Task, error) {
	toStatus = strings.TrimSpace(toStatus)
	if current.Status != TaskStatusOpen && current.Status != TaskStatusOverdue {
		return Task{}, ErrInvalidTaskTransition
	}
	current.CancellationReason = strings.TrimSpace(reason)
	switch toStatus {
	case TaskStatusCompleted:
		current.Status = TaskStatusCompleted
		current.CompletedAt = time.Now().UTC()
	case TaskStatusCancelled:
		current.Status = TaskStatusCancelled
		current.CancelledAt = time.Now().UTC()
	default:
		return Task{}, ErrInvalidTaskTransition
	}
	return current, nil
}

func EffectiveTaskStatus(task Task, businessDate time.Time) string {
	if task.Status != TaskStatusOpen {
		return task.Status
	}
	if !businessDate.IsZero() && businessDate.After(task.DueDate) {
		return TaskStatusOverdue
	}
	return TaskStatusOpen
}

func ParseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func FormatDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format("2006-01-02")
}
