package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	StatusUnassigned           = "Unassigned"
	StatusPendingQualification = "Pending Qualification"
	StatusValid                = "Valid"
	StatusInvalid              = "Invalid"
	StatusConverted            = "Converted To Opportunity"
)

var ErrValidation = errors.New("validation failed")

type Lead struct {
	ID                       string
	LeadName                 string
	CompanyName              string
	Email                    string
	Phone                    string
	Source                   string
	Status                   string
	OwnerID                  string
	NeedSummary              string
	InvalidReason            string
	ConvertedAccountID       string
	ConvertedOpportunityID   string
	ConversionIDempotencyKey string
	ArchivedAt               time.Time
	ArchivedBy               string
	ArchiveReason            string
	Version                  int
	UpdatedAt                time.Time
}

func NewLead(input Lead) (Lead, error) {
	input.LeadName = strings.TrimSpace(input.LeadName)
	input.CompanyName = strings.TrimSpace(input.CompanyName)
	input.Email = strings.TrimSpace(input.Email)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Source = strings.TrimSpace(input.Source)
	input.OwnerID = strings.TrimSpace(input.OwnerID)
	input.NeedSummary = strings.TrimSpace(input.NeedSummary)
	if input.LeadName == "" && input.CompanyName == "" {
		return Lead{}, ErrValidation
	}
	if input.Source == "" {
		return Lead{}, ErrValidation
	}
	if input.OwnerID == "" {
		input.Status = StatusUnassigned
	} else {
		input.Status = StatusPendingQualification
	}
	input.Version = 1
	return input, nil
}

func CanArchiveLead(actorRole string) bool {
	return actorRole == "Administrator" || actorRole == "Sales Manager"
}
