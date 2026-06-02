package domain

import (
	"strings"
	"time"
)

type Contact struct {
	ID          string
	AccountID   string
	AccountName string
	ContactName string
	Email       string
	Phone       string
	RoleNote    string
	Version     int
	UpdatedAt   time.Time
}

func NewContact(input Contact) (Contact, error) {
	input.AccountID = strings.TrimSpace(input.AccountID)
	input.ContactName = strings.TrimSpace(input.ContactName)
	input.Email = strings.TrimSpace(input.Email)
	input.Phone = strings.TrimSpace(input.Phone)
	input.RoleNote = strings.TrimSpace(input.RoleNote)
	if input.AccountID == "" || input.ContactName == "" {
		return Contact{}, ErrValidation
	}
	if input.Email == "" && input.Phone == "" && input.RoleNote == "" {
		return Contact{}, ErrValidation
	}
	input.Version = 1
	return input, nil
}
