package domain

import (
	"strings"
	"unicode"
)

type DuplicateCandidate struct {
	TargetType  string
	CompanyName string
	Email       string
	Phone       string
}

type DuplicateMatch struct {
	Type          string
	Rule          string
	MatchStrength string
	SafeSummary   string
	Visible       bool
}

type DuplicateCheckResult struct {
	Result           string
	WarningToken     string
	NormalizedFields []string
	Matches          []DuplicateMatch
	Rules            []string
	Signature        string
}

func NormalizeCompanyName(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), " "))
}

func NormalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func NormalizePhone(value string) string {
	var builder strings.Builder
	for _, char := range value {
		if unicode.IsDigit(char) {
			builder.WriteRune(char)
		}
	}
	digits := builder.String()
	if strings.HasPrefix(digits, "86") && len(digits) > 11 {
		return digits[2:]
	}
	return digits
}

func DuplicateSignature(candidate DuplicateCandidate) (string, []string) {
	parts := make([]string, 0, 3)
	fields := make([]string, 0, 3)
	if company := NormalizeCompanyName(candidate.CompanyName); company != "" {
		parts = append(parts, "companyName="+company)
		fields = append(fields, "companyName")
	}
	if email := NormalizeEmail(candidate.Email); email != "" {
		parts = append(parts, "email="+email)
		fields = append(fields, "email")
	}
	if phone := NormalizePhone(candidate.Phone); phone != "" {
		parts = append(parts, "phone="+phone)
		fields = append(fields, "phone")
	}
	return strings.Join(parts, "|"), fields
}
