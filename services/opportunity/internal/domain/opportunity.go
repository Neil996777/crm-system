package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	StageNewOpportunity      = "New Opportunity"
	StageNeedsConfirmed      = "Needs Confirmed"
	StageQuote               = "Quote"
	StageContractNegotiation = "Contract Negotiation"
	StageWon                 = "Won"
	StageLost                = "Lost"
)

var ErrValidation = errors.New("validation failed")

type Opportunity struct {
	ID                string
	CustomerID        string
	OwnerID           string
	Stage             string
	ExpectedAmount    string
	ExpectedCloseDate time.Time
	Title             string
	CloseDate         time.Time
	WonContractID     string
	LostReasonCode    string
	LostReasonDetail  string
	ClosedAt          time.Time
	ArchivedAt        time.Time
	ArchivedBy        string
	ArchiveReason     string
	Version           int
	UpdatedAt         time.Time
}

func NewOpportunity(input Opportunity) (Opportunity, error) {
	input.CustomerID = strings.TrimSpace(input.CustomerID)
	input.OwnerID = strings.TrimSpace(input.OwnerID)
	input.Stage = strings.TrimSpace(input.Stage)
	input.ExpectedAmount = normalizeAmount(input.ExpectedAmount)
	input.Title = strings.TrimSpace(input.Title)
	if input.CustomerID == "" || input.OwnerID == "" || input.Stage == "" || input.ExpectedAmount == "" || input.ExpectedCloseDate.IsZero() {
		return Opportunity{}, ErrValidation
	}
	if !IsStage(input.Stage) {
		return Opportunity{}, ErrValidation
	}
	input.Version = 1
	return input, nil
}

func UpdateOpportunity(current Opportunity, customerID, ownerID, stage, expectedAmount string, expectedCloseDate time.Time, title string) (Opportunity, error) {
	updated := current
	updated.CustomerID = strings.TrimSpace(customerID)
	updated.OwnerID = strings.TrimSpace(ownerID)
	updated.Stage = strings.TrimSpace(stage)
	updated.ExpectedAmount = normalizeAmount(expectedAmount)
	updated.ExpectedCloseDate = expectedCloseDate
	updated.Title = strings.TrimSpace(title)
	if updated.CustomerID == "" || updated.OwnerID == "" || updated.Stage == "" || updated.ExpectedAmount == "" || updated.ExpectedCloseDate.IsZero() {
		return Opportunity{}, ErrValidation
	}
	if !IsStage(updated.Stage) {
		return Opportunity{}, ErrValidation
	}
	updated.Version = current.Version + 1
	return updated, nil
}

func IsStage(stage string) bool {
	switch stage {
	case StageNewOpportunity, StageNeedsConfirmed, StageQuote, StageContractNegotiation, StageWon, StageLost:
		return true
	default:
		return false
	}
}

func CanCreateOpportunity(actorID, actorRole, ownerID string) bool {
	switch actorRole {
	case "Administrator", "Sales Manager":
		return true
	case "Sales":
		return ownerID == actorID && ownerID != ""
	default:
		return false
	}
}

func CanReadOpportunity(actorID, actorRole string, opportunity Opportunity) bool {
	switch actorRole {
	case "Administrator", "Sales Manager":
		return true
	case "Sales":
		return opportunity.OwnerID == actorID && opportunity.OwnerID != ""
	default:
		return false
	}
}

func CanEditOpportunity(actorID, actorRole string, opportunity Opportunity, newOwnerID string) bool {
	switch actorRole {
	case "Administrator", "Sales Manager":
		return true
	case "Sales":
		return opportunity.OwnerID == actorID && newOwnerID == actorID && newOwnerID != ""
	default:
		return false
	}
}

func CanArchiveOpportunity(actorRole string) bool {
	return actorRole == "Administrator" || actorRole == "Sales Manager"
}

func ParseCloseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func FormatCloseDate(value time.Time) string {
	return value.UTC().Format("2006-01-02")
}

func normalizeAmount(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil || amount <= 0 {
		return ""
	}
	return fmt.Sprintf("%.2f", amount)
}
