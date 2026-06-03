package handler

import (
	"time"

	"crm-system/services/import-export/internal/repo"
)

func exportRunDTO(run repo.ExportRun, content string) map[string]any {
	return map[string]any{
		"runId":              run.RunID,
		"objectType":         run.ObjectType,
		"filename":           run.Filename,
		"status":             run.Status,
		"exportedCount":      run.ExportedCount,
		"archivedIncluded":   run.IncludeArchived,
		"content":            content,
		"operationLogStatus": run.OperationLogStatus,
		"cleanupStatus":      run.CleanupStatus,
		"retainedUntil":      run.RetainedUntil.Format(time.RFC3339),
		"fileSafety":         "dangerous_cells_prefixed",
	}
}
