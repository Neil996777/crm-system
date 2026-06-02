package domain

import "strings"

func Convert(current Lead, idempotencyKey, accountID, opportunityID string) (Lead, error) {
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	accountID = strings.TrimSpace(accountID)
	opportunityID = strings.TrimSpace(opportunityID)
	if current.Status != StatusValid || current.OwnerID == "" || idempotencyKey == "" || accountID == "" || opportunityID == "" {
		return Lead{}, ErrValidation
	}
	updated := current
	updated.Status = StatusConverted
	updated.ConvertedAccountID = accountID
	updated.ConvertedOpportunityID = opportunityID
	updated.ConversionIDempotencyKey = idempotencyKey
	updated.Version = current.Version + 1
	return updated, nil
}
