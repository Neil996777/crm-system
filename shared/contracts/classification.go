package contracts

type DataClassification string

const (
	ClassificationSecurityCritical DataClassification = "Security Critical"
	ClassificationConfidential     DataClassification = "Confidential"
	ClassificationRestricted       DataClassification = "Restricted"
)

func (c DataClassification) String() string {
	return string(c)
}

func IsCommittedClassification(classification DataClassification) bool {
	switch classification {
	case ClassificationSecurityCritical, ClassificationConfidential, ClassificationRestricted:
		return true
	default:
		return false
	}
}
