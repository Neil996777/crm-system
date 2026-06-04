package domain

func EscapeDangerousCSVCell(value string) string {
	if IsDangerousCSVCell(value) {
		return "'" + value
	}
	return value
}
