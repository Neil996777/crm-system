package domain

type ReminderRelatedRecord struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Display string `json:"display"`
}

type ReminderRow struct {
	ID            string                `json:"id"`
	SourceService string                `json:"sourceService"`
	Type          string                `json:"type"`
	RelatedRecord ReminderRelatedRecord `json:"relatedRecord"`
	OwnerDisplay  string                `json:"ownerDisplay"`
	DueDate       string                `json:"dueDate"`
	Status        string                `json:"status"`
	Priority      string                `json:"priority"`
	Version       int                   `json:"version"`
}
