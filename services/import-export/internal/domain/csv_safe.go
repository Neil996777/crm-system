package domain

import "strings"

func IsDangerousCSVCell(value string) bool {
	if value == "" {
		return false
	}
	trimmed := strings.TrimLeft(value, " ")
	if trimmed == "" {
		return false
	}
	switch trimmed[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return true
	default:
		return false
	}
}
