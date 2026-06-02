package domain

func CanReadLead(actorID, actorRole string, lead Lead) bool {
	switch actorRole {
	case "Administrator", "Sales Manager":
		return true
	case "Sales":
		return lead.OwnerID == actorID && lead.OwnerID != ""
	default:
		return false
	}
}

func CanCreateLead(actorID, actorRole string, ownerID string) bool {
	switch actorRole {
	case "Administrator", "Sales Manager":
		return true
	case "Sales":
		return ownerID == "" || ownerID == actorID
	default:
		return false
	}
}

func CanTransferOwner(actorRole string) bool {
	return actorRole == "Administrator" || actorRole == "Sales Manager"
}
