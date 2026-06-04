package handler

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	auditclient "crm-system/services/import-export/internal/client"
	"crm-system/services/import-export/internal/domain"
	"crm-system/services/import-export/internal/repo"
)

func (h *ImportExportHandler) startExport(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	if actor.Role == "" || actor.Role == "Sales" {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	var request struct {
		ObjectType      string `json:"objectType"`
		Confirmed       bool   `json:"confirmed"`
		IncludeArchived bool   `json:"includeArchived"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The export input is invalid.")
		return
	}
	objectType := normalizeObjectType(request.ObjectType)
	if objectType != "lead" {
		writeError(w, http.StatusBadRequest, "UNSUPPORTED_OBJECT_TYPE", "validation", "The object type is not supported for export.")
		return
	}
	if !request.Confirmed {
		writeError(w, http.StatusBadRequest, "EXPORT_CONFIRMATION_REQUIRED", "validation", "Export confirmation is required.")
		return
	}
	items, err := h.queryLeadsForExport(r, actor, request.IncludeArchived)
	if err != nil {
		writeError(w, http.StatusBadRequest, "EXPORT_FAILED", "operation", "The export could not be completed.")
		return
	}
	content, err := leadCSV(items)
	if err != nil {
		writeError(w, http.StatusBadRequest, "EXPORT_FAILED", "operation", "The export could not be completed.")
		return
	}
	run := repo.ExportRun{
		RunID:              newRunID(),
		ObjectType:         objectType,
		Filename:           "lead-export.csv",
		Status:             "Completed",
		ActorID:            actor.ID,
		ActorRole:          actor.Role,
		TeamID:             actor.TeamIDOrDefault(),
		IncludeArchived:    request.IncludeArchived,
		ExportedCount:      len(items),
		OperationLogStatus: "not_configured",
		CleanupStatus:      "pending",
		RetainedUntil:      time.Now().UTC().Add(24 * time.Hour),
	}
	completedAt := time.Now().UTC()
	run.CompletedAt = &completedAt
	run.OperationLogStatus = h.appendExportOperationLog(r, run, actor)
	if err := h.repo.SaveExportRun(r.Context(), run); err != nil {
		writeError(w, http.StatusBadRequest, "EXPORT_FAILED", "operation", "The export could not be completed.")
		return
	}
	writeJSON(w, http.StatusCreated, exportRunDTO(run, content))
}

func (h *ImportExportHandler) queryLeadsForExport(r *http.Request, actor actorContext, includeArchived bool) ([]map[string]any, error) {
	if h.leadBaseURL == "" {
		return nil, errors.New("lead service not configured")
	}
	url := h.leadBaseURL + "/leads"
	if includeArchived {
		url += "?includeArchived=true"
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Actor-User-Id", actor.ID)
	req.Header.Set("X-Actor-Role", actor.Role)
	req.Header.Set("X-Actor-Team-Id", actor.TeamIDOrDefault())
	resp, err := h.targetClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("target query failed")
	}
	var body struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Items == nil {
		return []map[string]any{}, nil
	}
	return body.Items, nil
}

func leadCSV(items []map[string]any) (string, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	headers := []string{"id", "companyName", "leadName", "source", "ownerId", "status"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}
	for _, item := range items {
		row := []string{}
		for _, header := range headers {
			row = append(row, domain.EscapeDangerousCSVCell(stringValue(item[header])))
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(value.(string))
}

func (h *ImportExportHandler) appendExportOperationLog(r *http.Request, run repo.ExportRun, actor actorContext) string {
	if err := h.audit.AppendExportRun(r.Context(), auditclient.ExportRunLogInput{
		RunID:           run.RunID,
		ActorID:         actor.ID,
		ActorRole:       actor.Role,
		ObjectType:      run.ObjectType,
		IncludeArchived: run.IncludeArchived,
		ExportedCount:   run.ExportedCount,
		Result:          run.Status,
		CorrelationID:   r.Header.Get("X-Correlation-Id"),
	}); err != nil {
		return "failed"
	}
	if h.auditConfigured() {
		return "logged"
	}
	return "not_configured"
}
