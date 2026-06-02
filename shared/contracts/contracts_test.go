package contracts

import "testing"

func TestCoreContractConstantsExist(t *testing.T) {
	if ErrorServiceAuthFailed == "" {
		t.Fatal("SERVICE_AUTH_FAILED error constant must exist")
	}
	if HeaderCorrelationID == "" {
		t.Fatal("correlation header constant must exist")
	}
	if PermissionActionLeadCreate == "" {
		t.Fatal("permission action constants must exist")
	}
}
