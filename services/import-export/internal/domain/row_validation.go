package domain

import "strings"

type RowError struct {
	RowNumber   int    `json:"rowNumber"`
	Field       string `json:"field"`
	Code        string `json:"code"`
	SafeMessage string `json:"safeMessage"`
}

func ValidateLeadRow(rowNumber int, row map[string]string) []RowError {
	var errors []RowError
	if strings.TrimSpace(row["source"]) == "" {
		errors = append(errors, rowError(rowNumber, "source", "REQUIRED_FIELD_MISSING", "Source is required."))
	}
	if strings.TrimSpace(row["companyName"]) == "" && strings.TrimSpace(row["leadName"]) == "" {
		errors = append(errors, rowError(rowNumber, "companyName", "REQUIRED_FIELD_MISSING", "Company name or lead name is required."))
	}
	for field, value := range row {
		if IsDangerousCSVCell(value) {
			errors = append(errors, rowError(rowNumber, field, "DANGEROUS_CSV_CELL", "CSV cell is not safe to import."))
		}
	}
	return errors
}

func rowError(rowNumber int, field, code, safeMessage string) RowError {
	return RowError{RowNumber: rowNumber, Field: field, Code: code, SafeMessage: safeMessage}
}
