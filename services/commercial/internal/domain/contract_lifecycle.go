package domain

import (
	"strings"
	"time"
)

func ValidateContractStatusTransition(current Contract, toStatus string, requestedSignedEffectiveDate time.Time) (time.Time, error) {
	toStatus = strings.TrimSpace(toStatus)
	if requiresSignedEffectiveDate(current.Status, toStatus) && current.SignedEffectiveDate.IsZero() && requestedSignedEffectiveDate.IsZero() {
		return time.Time{}, ErrSignedEffectiveDateRequired
	}
	if !canTransitionContract(current.Status, toStatus) {
		return time.Time{}, ErrInvalidContractTransition
	}
	if requestedSignedEffectiveDate.IsZero() {
		return current.SignedEffectiveDate, nil
	}
	return requestedSignedEffectiveDate, nil
}

func canTransitionContract(fromStatus, toStatus string) bool {
	switch fromStatus {
	case ContractStatusPendingSignature:
		return toStatus == ContractStatusSigned || toStatus == ContractStatusTerminated
	case ContractStatusSigned:
		return toStatus == ContractStatusActive || toStatus == ContractStatusTerminated
	case ContractStatusActive:
		return toStatus == ContractStatusCompleted || toStatus == ContractStatusTerminated
	}
	return false
}

func requiresSignedEffectiveDate(fromStatus, toStatus string) bool {
	if toStatus == ContractStatusSigned || toStatus == ContractStatusActive || toStatus == ContractStatusCompleted {
		return true
	}
	return toStatus == ContractStatusTerminated && (fromStatus == ContractStatusSigned || fromStatus == ContractStatusActive)
}
