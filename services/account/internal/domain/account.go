package domain

import (
	"errors"
	"strings"
	"time"
)

var ErrValidation = errors.New("validation failed")

type Account struct {
	ID             string
	CompanyName    string
	CustomerStatus string
	OwnerID        string
	ArchivedAt     time.Time
	ArchivedBy     string
	ArchiveReason  string
	Version        int
	UpdatedAt      time.Time
}

func NewAccount(input Account) (Account, error) {
	input.CompanyName = strings.TrimSpace(input.CompanyName)
	input.CustomerStatus = strings.TrimSpace(input.CustomerStatus)
	input.OwnerID = strings.TrimSpace(input.OwnerID)
	if input.CompanyName == "" || input.CustomerStatus == "" || input.OwnerID == "" {
		return Account{}, ErrValidation
	}
	input.Version = 1
	return input, nil
}

func UpdateAccount(current Account, companyName, customerStatus, ownerID string) (Account, error) {
	updated := current
	updated.CompanyName = strings.TrimSpace(companyName)
	updated.CustomerStatus = strings.TrimSpace(customerStatus)
	updated.OwnerID = strings.TrimSpace(ownerID)
	if updated.CompanyName == "" || updated.CustomerStatus == "" || updated.OwnerID == "" {
		return Account{}, ErrValidation
	}
	updated.Version = current.Version + 1
	return updated, nil
}

func CanCreateAccount(actorID, actorRole, ownerID string) bool {
	switch actorRole {
	case "Administrator", "Sales Manager":
		return true
	case "Sales":
		return ownerID == actorID && ownerID != ""
	default:
		return false
	}
}

func CanReadAccount(actorID, actorRole string, account Account) bool {
	switch actorRole {
	case "Administrator", "Sales Manager":
		return true
	case "Sales":
		return account.OwnerID == actorID && account.OwnerID != ""
	default:
		return false
	}
}

func CanEditAccount(actorID, actorRole string, account Account, newOwnerID string) bool {
	switch actorRole {
	case "Administrator", "Sales Manager":
		return true
	case "Sales":
		return account.OwnerID == actorID && newOwnerID == actorID && newOwnerID != ""
	default:
		return false
	}
}

func CanArchiveAccount(actorRole string) bool {
	return actorRole == "Administrator" || actorRole == "Sales Manager"
}
