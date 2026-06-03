package domain

type ExportScope struct {
	ObjectType      string
	IncludeArchived bool
	Confirmed       bool
}

func EscapeDangerousCSVCell(value string) string {
	if IsDangerousCSVCell(value) {
		return "'" + value
	}
	return value
}
