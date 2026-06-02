package domain

import "strings"

func QualifyValid(current Lead) (Lead, error) {
	if current.Status != StatusPendingQualification || current.OwnerID == "" {
		return Lead{}, ErrValidation
	}
	updated := current
	updated.Status = StatusValid
	updated.InvalidReason = ""
	updated.Version = current.Version + 1
	return updated, nil
}

func QualifyInvalid(current Lead, reason string) (Lead, error) {
	reason = strings.TrimSpace(reason)
	if current.Status != StatusPendingQualification || current.OwnerID == "" || reason == "" {
		return Lead{}, ErrValidation
	}
	updated := current
	updated.Status = StatusInvalid
	updated.InvalidReason = reason
	updated.Version = current.Version + 1
	return updated, nil
}

func RestoreInvalid(current Lead) (Lead, error) {
	if current.Status != StatusInvalid || current.OwnerID == "" {
		return Lead{}, ErrValidation
	}
	updated := current
	updated.Status = StatusPendingQualification
	updated.InvalidReason = ""
	updated.Version = current.Version + 1
	return updated, nil
}

func CanQualifyLead(actorID, actorRole string, lead Lead) bool {
	return CanReadLead(actorID, actorRole, lead) && lead.OwnerID != ""
}

func CanRestoreLead(actorRole string) bool {
	return actorRole == "Administrator" || actorRole == "Sales Manager"
}
