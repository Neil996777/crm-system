package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"crm-system/services/work/internal/domain"
)

func activityDTO(activity domain.Activity) map[string]any {
	return map[string]any{
		"id":           activity.ID,
		"relatedType":  activity.RelatedType,
		"relatedId":    activity.RelatedID,
		"activityType": activity.ActivityType,
		"content":      activity.Content,
		"actorId":      activity.ActorID,
		"ownerId":      activity.OwnerID,
		"occurredAt":   activity.OccurredAt.UTC().Format(time.RFC3339),
		"version":      activity.Version,
		"updatedAt":    activity.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func noteDTO(note domain.Note) map[string]any {
	return map[string]any{
		"id":          note.ID,
		"relatedType": note.RelatedType,
		"relatedId":   note.RelatedID,
		"content":     note.Content,
		"actorId":     note.ActorID,
		"ownerId":     note.OwnerID,
		"occurredAt":  note.OccurredAt.UTC().Format(time.RFC3339),
		"version":     note.Version,
		"updatedAt":   note.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func taskDTO(task domain.Task, businessDate time.Time) map[string]any {
	return map[string]any{
		"id":          task.ID,
		"relatedType": task.RelatedType,
		"relatedId":   task.RelatedID,
		"title":       task.Title,
		"dueDate":     domain.FormatDate(task.DueDate),
		"status":      domain.EffectiveTaskStatus(task, businessDate),
		"actorId":     task.ActorID,
		"ownerId":     task.OwnerID,
		"version":     task.Version,
		"updatedAt":   task.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func writeError(w http.ResponseWriter, status int, code, category, safeMessage string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":        code,
			"category":    category,
			"safeMessage": safeMessage,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write response: %v", err)
	}
}
