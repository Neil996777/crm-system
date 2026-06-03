package domain

import "testing"

func TestCSVRowValidationAcceptance(t *testing.T) {
	t.Run("TEST-CSV-IMPORT-002 required fields produce safe row errors", func(t *testing.T) {
		errors := ValidateLeadRow(2, map[string]string{"companyName": "Missing Source"})
		if len(errors) == 0 {
			t.Fatalf("expected missing source row error")
		}
		if errors[0].RowNumber != 2 || errors[0].Field != "source" || errors[0].Code != "REQUIRED_FIELD_MISSING" {
			t.Fatalf("expected safe required-field error, got %#v", errors[0])
		}
	})

	t.Run("TEST-ABUSE-CSVINJECT-001 dangerous formula cells rejected", func(t *testing.T) {
		errors := ValidateLeadRow(3, map[string]string{"companyName": "=cmd|' /C calc'!A0", "source": "Website"})
		if len(errors) == 0 {
			t.Fatalf("expected dangerous-cell row error")
		}
		if errors[0].Code != "DANGEROUS_CSV_CELL" || errors[0].SafeMessage == "" {
			t.Fatalf("expected safe dangerous-cell error, got %#v", errors[0])
		}
	})
}
