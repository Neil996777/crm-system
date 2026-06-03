package domain

import "time"

func ApplyClassificationAndRetention(event Event) Event {
	if event.DiffClassification == "" {
		event.DiffClassification = "Confidential"
	}
	event.RetentionPolicy = retentionPolicy(event)
	if event.RetainUntil.IsZero() {
		event.RetainUntil = retainUntil(event.OccurredAt, event.RetentionPolicy)
	}
	return event
}

func retentionPolicy(event Event) string {
	if event.ResourceType == "auth" || event.ReasonCode == "access_denied" {
		return "operation_log_access_failure_min_3y"
	}
	if event.ResourceType == "import_run" {
		return "import_result_metadata_min_1y"
	}
	return "business_record_or_operation_log_min_7y"
}

func retainUntil(anchor time.Time, policy string) time.Time {
	if anchor.IsZero() {
		anchor = time.Now().UTC()
	}
	switch policy {
	case "operation_log_access_failure_min_3y":
		return anchor.AddDate(3, 0, 0)
	case "import_result_metadata_min_1y":
		return anchor.AddDate(1, 0, 0)
	default:
		return anchor.AddDate(7, 0, 0)
	}
}

func MaskedSummary(classification string, summary map[string]any) map[string]any {
	if classification != "Restricted" && classification != "Security Critical" {
		return summary
	}
	masked := map[string]any{}
	for key := range summary {
		masked[key] = "[masked]"
	}
	return masked
}
