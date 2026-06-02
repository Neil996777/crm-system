package domain

import "time"

const ReminderTimezone = "Asia/Shanghai"

type RelatedRecord struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Display string `json:"display"`
}

type ReminderRow struct {
	ID            string        `json:"id"`
	SourceService string        `json:"sourceService"`
	Type          string        `json:"type"`
	RelatedRecord RelatedRecord `json:"relatedRecord"`
	OwnerDisplay  string        `json:"ownerDisplay"`
	DueDate       string        `json:"dueDate"`
	Status        string        `json:"status"`
	Priority      string        `json:"priority"`
	Version       int           `json:"version"`
}

func ReminderFromTask(task Task, businessDate time.Time) (ReminderRow, bool) {
	status := EffectiveTaskStatus(task, businessDate)
	if status != TaskStatusOpen && status != TaskStatusOverdue {
		return ReminderRow{}, false
	}
	if task.DueDate.After(businessDate) {
		return ReminderRow{}, false
	}
	reminderType := "task_due"
	reminderStatus := "DueToday"
	if status == TaskStatusOverdue {
		reminderType = "task_overdue"
		reminderStatus = "Overdue"
	}
	return ReminderRow{
		ID:            task.ID,
		SourceService: "work-service",
		Type:          reminderType,
		RelatedRecord: RelatedRecord{Type: task.RelatedType, ID: task.RelatedID, Display: task.Title},
		OwnerDisplay:  task.OwnerID,
		DueDate:       FormatDate(task.DueDate),
		Status:        reminderStatus,
		Priority:      "P1",
		Version:       task.Version,
	}, true
}
