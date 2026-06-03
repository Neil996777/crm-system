package contracts

import "testing"

func TestDataClassificationContract(t *testing.T) {
	if ClassificationSecurityCritical.String() != "Security Critical" {
		t.Fatalf("expected Security Critical classification")
	}
	if !IsCommittedClassification(ClassificationConfidential) || !IsCommittedClassification(ClassificationRestricted) {
		t.Fatalf("expected committed classifications to validate")
	}
	if IsCommittedClassification(DataClassification("Public")) {
		t.Fatalf("unexpected public classification accepted")
	}
}
