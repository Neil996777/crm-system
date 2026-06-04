package handler

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	auditclient "crm-system/services/import-export/internal/client"
	"crm-system/services/import-export/internal/domain"
	"crm-system/services/import-export/internal/repo"
)

type Config struct {
	LeadServiceURL         string
	AuditHistoryServiceURL string
	ServiceID              string
	ServiceTokenSecret     []byte
	HTTPClient             *http.Client
}

type ImportExportHandler struct {
	repo         *repo.RunRepo
	leadBaseURL  string
	targetClient *http.Client
	audit        auditclient.AuditClient
}

func NewImportExportServer(db *sql.DB, config Config) http.Handler {
	client := config.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	handler := &ImportExportHandler{
		repo:         repo.NewRunRepo(db),
		leadBaseURL:  strings.TrimRight(config.LeadServiceURL, "/"),
		targetClient: client,
		audit:        auditclient.NewAuditClient(config.AuditHistoryServiceURL, config.ServiceID, config.ServiceTokenSecret, client),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /imports", handler.startImport)
	mux.HandleFunc("GET /imports/{id}", handler.getImportRun)
	mux.HandleFunc("POST /exports", handler.startExport)
	return mux
}

func (h *ImportExportHandler) startImport(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	if actor.Role == "" || actor.Role == "Sales" {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	var request struct {
		ObjectType string `json:"objectType"`
		Filename   string `json:"filename"`
		Content    string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The import input is invalid.")
		return
	}
	objectType := normalizeObjectType(request.ObjectType)
	if objectType != "lead" {
		writeError(w, http.StatusBadRequest, "UNSUPPORTED_OBJECT_TYPE", "validation", "The object type is not supported for import.")
		return
	}
	if strings.ToLower(filepath.Ext(request.Filename)) != ".csv" {
		writeError(w, http.StatusBadRequest, "UNSUPPORTED_FORMAT", "validation", "Only CSV import is supported.")
		return
	}
	rows, err := parseCSVRows(request.Content)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The CSV content is invalid.")
		return
	}
	run := repo.ImportRun{
		RunID:              newRunID(),
		ObjectType:         objectType,
		Filename:           request.Filename,
		ActorID:            actor.ID,
		ActorRole:          actor.Role,
		TeamID:             actor.TeamIDOrDefault(),
		TotalRows:          len(rows),
		OperationLogStatus: "not_configured",
		CleanupStatus:      "pending",
		RetainedUntil:      time.Now().UTC().Add(24 * time.Hour),
	}
	for _, row := range rows {
		rowErrors := domain.ValidateLeadRow(row.number, row.values)
		if len(rowErrors) > 0 {
			for _, rowError := range rowErrors {
				run.RowErrors = append(run.RowErrors, repo.ImportRowResult{
					RowNumber:   rowError.RowNumber,
					Success:     false,
					Field:       rowError.Field,
					Code:        rowError.Code,
					SafeMessage: rowError.SafeMessage,
				})
				run.FailureCount++
			}
			continue
		}
		targetID, err := h.createLead(r, row.values, actor)
		if err != nil {
			run.RowErrors = append(run.RowErrors, repo.ImportRowResult{
				RowNumber:   row.number,
				Success:     false,
				Field:       "row",
				Code:        "TARGET_COMMAND_FAILED",
				SafeMessage: "Row could not be imported.",
			})
			run.FailureCount++
			continue
		}
		run.SuccessCount++
		run.RowErrors = append(run.RowErrors, repo.ImportRowResult{RowNumber: row.number, Success: true, TargetRecordID: targetID})
	}
	run.Status = "Completed"
	if run.FailureCount > 0 && run.SuccessCount > 0 {
		run.Status = "CompletedWithErrors"
	}
	if run.FailureCount > 0 && run.SuccessCount == 0 {
		run.Status = "Failed"
	}
	completedAt := time.Now().UTC()
	run.CompletedAt = &completedAt
	run.OperationLogStatus = h.appendImportOperationLog(r, run, actor)
	if err := h.repo.SaveImportRun(r.Context(), run); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The import run could not be saved.")
		return
	}
	writeJSON(w, http.StatusCreated, importRunDTO(run))
}

func (h *ImportExportHandler) appendImportOperationLog(r *http.Request, run repo.ImportRun, actor actorContext) string {
	if err := h.audit.AppendImportRun(r.Context(), auditclient.ImportRunLogInput{
		RunID:         run.RunID,
		ActorID:       actor.ID,
		ActorRole:     actor.Role,
		ObjectType:    run.ObjectType,
		TotalRows:     run.TotalRows,
		SuccessCount:  run.SuccessCount,
		FailureCount:  run.FailureCount,
		Result:        run.Status,
		CorrelationID: r.Header.Get("X-Correlation-Id"),
	}); err != nil {
		if errors.Is(err, auditclient.ErrServiceAuthFailed) {
			return "failed"
		}
		return "failed"
	}
	if h.auditConfigured() {
		return "logged"
	}
	return "not_configured"
}

func (h *ImportExportHandler) auditConfigured() bool {
	return h.audit.BaseURLConfigured()
}

func (h *ImportExportHandler) getImportRun(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	run, err := h.repo.FindImportRun(r.Context(), r.PathValue("id"))
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if !canReadImportRun(actor, run) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	writeJSON(w, http.StatusOK, importRunDTO(run))
}

func canReadImportRun(actor actorContext, run repo.ImportRun) bool {
	if actor.ID == "" || actor.Role == "" {
		return false
	}
	if actor.Role == "Administrator" {
		return true
	}
	if actor.Role == "Sales Manager" {
		return actor.ID == run.ActorID || actor.TeamIDOrDefault() == run.TeamID
	}
	return false
}

func (h *ImportExportHandler) createLead(r *http.Request, row map[string]string, actor actorContext) (string, error) {
	if h.leadBaseURL == "" {
		return "", errors.New("lead service not configured")
	}
	body, err := json.Marshal(map[string]any{
		"companyName": strings.TrimSpace(row["companyName"]),
		"leadName":    strings.TrimSpace(row["leadName"]),
		"email":       strings.TrimSpace(row["email"]),
		"phone":       strings.TrimSpace(row["phone"]),
		"source":      strings.TrimSpace(row["source"]),
		"ownerId":     strings.TrimSpace(row["ownerId"]),
		"needSummary": strings.TrimSpace(row["needSummary"]),
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, h.leadBaseURL+"/leads", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Actor-User-Id", actor.ID)
	req.Header.Set("X-Actor-Role", actor.Role)
	req.Header.Set("X-Actor-Team-Id", actor.TeamIDOrDefault())
	if correlationID := strings.TrimSpace(r.Header.Get("X-Correlation-Id")); correlationID != "" {
		req.Header.Set("X-Correlation-Id", correlationID)
	}
	resp, err := h.targetClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("target command failed")
	}
	var response map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	if id, ok := response["id"].(string); ok {
		return id, nil
	}
	return "", nil
}

type csvRow struct {
	number int
	values map[string]string
}

func parseCSVRows(content string) ([]csvRow, error) {
	reader := csv.NewReader(strings.NewReader(content))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
	}
	var rows []csvRow
	line := 1
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		line++
		values := map[string]string{}
		for i, header := range headers {
			if i < len(record) {
				values[header] = strings.TrimSpace(record[i])
			} else {
				values[header] = ""
			}
		}
		rows = append(rows, csvRow{number: line, values: values})
	}
	return rows, nil
}

func normalizeObjectType(objectType string) string {
	return strings.ToLower(strings.TrimSpace(objectType))
}

func newRunID() string {
	var data [16]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "import-" + hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}
	return "import-" + hex.EncodeToString(data[:])
}

func importRunDTO(run repo.ImportRun) map[string]any {
	rowErrors := []map[string]any{}
	for _, row := range run.RowErrors {
		if row.Success {
			continue
		}
		rowErrors = append(rowErrors, map[string]any{
			"rowNumber":   row.RowNumber,
			"field":       row.Field,
			"code":        row.Code,
			"safeMessage": row.SafeMessage,
		})
	}
	return map[string]any{
		"runId":              run.RunID,
		"objectType":         run.ObjectType,
		"filename":           run.Filename,
		"status":             run.Status,
		"totalRows":          run.TotalRows,
		"successCount":       run.SuccessCount,
		"failureCount":       run.FailureCount,
		"rowErrors":          rowErrors,
		"operationLogStatus": run.OperationLogStatus,
		"cleanupStatus":      run.CleanupStatus,
		"retainedUntil":      run.RetainedUntil.Format(time.RFC3339),
	}
}
