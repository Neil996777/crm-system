package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"crm-system/services/work/internal/authz"
	"crm-system/services/work/internal/domain"
	"crm-system/services/work/internal/event"
	"crm-system/services/work/internal/repo"
)

const (
	intentOwnerTransfer     = "work.owner_transfer"
	intentActiveObligations = "work.active_obligations"
)

type Config struct {
	ServiceID          string
	ServiceTokenSecret []byte
	CommercialBaseURL  string
	HTTPClient         *http.Client
}

type WorkHandler struct {
	db     *sql.DB
	repo   *repo.WorkRepo
	outbox *event.Outbox
	config Config
}

type actorContext struct {
	ID   string
	Role string
}

func NewWorkServer(db *sql.DB, config Config) http.Handler {
	if config.ServiceID == "" {
		config.ServiceID = "work"
	}
	handler := &WorkHandler{db: db, repo: repo.NewWorkRepo(db), outbox: event.NewOutbox(db), config: config}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /activities", handler.createActivity)
	mux.HandleFunc("GET /activities", handler.listActivities)
	mux.HandleFunc("POST /notes", handler.createNote)
	mux.HandleFunc("GET /notes", handler.listNotes)
	mux.HandleFunc("POST /tasks", handler.createTask)
	mux.HandleFunc("GET /tasks", handler.listTasks)
	mux.HandleFunc("GET /reminders", handler.listReminders)
	mux.HandleFunc("POST /tasks/{id}/status", handler.changeTaskStatus)
	mux.HandleFunc("POST /internal/owner-transfer", handler.transferOwner)
	mux.HandleFunc("GET /internal/active-obligations", handler.activeObligations)
	return mux
}

func (h *WorkHandler) inTransaction(ctx context.Context, fn func(*repo.WorkRepo, *event.Outbox) error) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(repo.NewWorkRepoTx(tx), event.NewOutboxTx(tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func writeOutboxFailure(w http.ResponseWriter) {
	writeError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "dependency", "The audit event could not be persisted.")
}

func SignServiceToken(issuer, audience, intent string, secret []byte) string {
	return authz.SignServiceToken(issuer, audience, intent, secret)
}

func (h *WorkHandler) createActivity(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		RelatedType  string `json:"relatedType"`
		RelatedID    string `json:"relatedId"`
		ActivityType string `json:"activityType"`
		Content      string `json:"content"`
		OwnerID      string `json:"ownerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The work item input is invalid.")
		return
	}
	if !canActForOwner(actor, request.OwnerID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	activity, err := domain.NewActivity(domain.Activity{
		RelatedType:  request.RelatedType,
		RelatedID:    request.RelatedID,
		ActivityType: request.ActivityType,
		Content:      request.Content,
		ActorID:      actor.ID,
		OwnerID:      request.OwnerID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The work item input is invalid.")
		return
	}
	var created domain.Activity
	if err := h.inTransaction(r.Context(), func(txRepo *repo.WorkRepo, txOutbox *event.Outbox) error {
		var err error
		created, err = txRepo.CreateActivity(r.Context(), activity)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.WorkItemCreated, created.ID, map[string]any{
			"traceability": "TASK-024 ACC-012 CIM-030 CIM-PROC-012 PIM-013 PIM-BEH-020 PSM-008 CONTRACT-011 CONTRACT-012",
			"actorId":      actor.ID,
			"actorRole":    actor.Role,
			"actorDisplay": actor.ID,
			"workItemId":   created.ID,
			"relatedType":  created.RelatedType,
			"relatedId":    created.RelatedID,
		})
	}); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The work item input is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, activityDTO(created))
}

func (h *WorkHandler) createNote(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		RelatedType string `json:"relatedType"`
		RelatedID   string `json:"relatedId"`
		Content     string `json:"content"`
		OwnerID     string `json:"ownerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The work item input is invalid.")
		return
	}
	if !canActForOwner(actor, request.OwnerID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	note, err := domain.NewNote(domain.Note{
		RelatedType: request.RelatedType,
		RelatedID:   request.RelatedID,
		Content:     request.Content,
		ActorID:     actor.ID,
		OwnerID:     request.OwnerID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The work item input is invalid.")
		return
	}
	var created domain.Note
	if err := h.inTransaction(r.Context(), func(txRepo *repo.WorkRepo, txOutbox *event.Outbox) error {
		var err error
		created, err = txRepo.CreateNote(r.Context(), note)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.WorkItemCreated, created.ID, map[string]any{
			"traceability": "TASK-024 ACC-012 CIM-031 CIM-PROC-012 PIM-014 PIM-BEH-020 PSM-008 CONTRACT-011 CONTRACT-012",
			"actorId":      actor.ID,
			"actorRole":    actor.Role,
			"actorDisplay": actor.ID,
			"workItemId":   created.ID,
			"relatedType":  created.RelatedType,
			"relatedId":    created.RelatedID,
		})
	}); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The work item input is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, noteDTO(created))
}

func (h *WorkHandler) createTask(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		RelatedType string `json:"relatedType"`
		RelatedID   string `json:"relatedId"`
		Title       string `json:"title"`
		DueDate     string `json:"dueDate"`
		OwnerID     string `json:"ownerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The task input is invalid.")
		return
	}
	if !canActForOwner(actor, request.OwnerID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	dueDate, err := domain.ParseDate(request.DueDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The task input is invalid.")
		return
	}
	task, err := domain.NewTask(domain.Task{
		RelatedType: request.RelatedType,
		RelatedID:   request.RelatedID,
		Title:       request.Title,
		DueDate:     dueDate,
		ActorID:     actor.ID,
		OwnerID:     request.OwnerID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The task input is invalid.")
		return
	}
	var created domain.Task
	if err := h.inTransaction(r.Context(), func(txRepo *repo.WorkRepo, txOutbox *event.Outbox) error {
		var err error
		created, err = txRepo.CreateTask(r.Context(), task)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.WorkItemCreated, created.ID, map[string]any{
			"traceability": "TASK-024 ACC-012 ACC-021 CIM-032 CIM-PROC-012 PIM-015 PIM-SM-007 PIM-BEH-021 PSM-008 CONTRACT-011 CONTRACT-012",
			"actorId":      actor.ID,
			"actorRole":    actor.Role,
			"actorDisplay": actor.ID,
			"taskId":       created.ID,
			"relatedType":  created.RelatedType,
			"relatedId":    created.RelatedID,
		})
	}); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The task input is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, taskDTO(created, time.Time{}))
}

func (h *WorkHandler) changeTaskStatus(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ToStatus           string `json:"toStatus"`
		CancellationReason string `json:"cancellationReason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The task status input is invalid.")
		return
	}
	current, err := h.repo.FindTask(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if !canActForOwner(actor, current.OwnerID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	updated, err := domain.ApplyTaskStatus(current, request.ToStatus, request.CancellationReason)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_TRANSITION", "business_rule", "The requested task status transition is not allowed.")
		return
	}
	var saved domain.Task
	if err := h.inTransaction(r.Context(), func(txRepo *repo.WorkRepo, txOutbox *event.Outbox) error {
		var err error
		saved, err = txRepo.UpdateTaskStatus(r.Context(), updated)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.TaskStatusChanged, saved.ID, map[string]any{
			"traceability": "TASK-024 ACC-012 ACC-014 ACC-021 CIM-032 PIM-SM-007 PIM-BEH-021 PSM-008 CONTRACT-011 CONTRACT-012",
			"actorId":      actor.ID,
			"actorRole":    actor.Role,
			"actorDisplay": actor.ID,
			"taskId":       saved.ID,
			"toStatus":     saved.Status,
		})
	}); err != nil {
		writeOutboxFailure(w)
		return
	}
	writeJSON(w, http.StatusOK, taskDTO(saved, time.Time{}))
}

func (h *WorkHandler) transferOwner(w http.ResponseWriter, r *http.Request) {
	if !h.verifyServiceToken(r, intentOwnerTransfer) {
		writeError(w, http.StatusUnauthorized, "SERVICE_AUTH_FAILED", "authentication", "Service authentication failed.")
		return
	}
	var request struct {
		RelatedType string `json:"relatedType"`
		RelatedID   string `json:"relatedId"`
		FromOwnerID string `json:"fromOwnerId"`
		ToOwnerID   string `json:"toOwnerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || strings.TrimSpace(request.ToOwnerID) == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The owner transfer input is invalid.")
		return
	}
	var count int64
	if err := h.inTransaction(r.Context(), func(txRepo *repo.WorkRepo, txOutbox *event.Outbox) error {
		var err error
		count, err = txRepo.TransferOpenWork(r.Context(), request.RelatedType, request.RelatedID, request.FromOwnerID, request.ToOwnerID)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.OpenWorkTransferred, request.RelatedID, map[string]any{
			"traceability": "TASK-024 ACC-012 ACC-014 CIM-032 PIM-INV-030 PIM-INV-033 PSM-008 CONTRACT-011 CONTRACT-012 EDGE-024",
			"actorId":      r.Header.Get("X-Actor-User-Id"),
			"actorRole":    r.Header.Get("X-Actor-Role"),
			"actorDisplay": r.Header.Get("X-Actor-User-Id"),
			"relatedType":  request.RelatedType,
			"relatedId":    request.RelatedID,
			"fromOwnerId":  request.FromOwnerID,
			"toOwnerId":    request.ToOwnerID,
			"count":        count,
		})
	}); err != nil {
		writeOutboxFailure(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"transferred": count})
}

func (h *WorkHandler) activeObligations(w http.ResponseWriter, r *http.Request) {
	if !h.verifyServiceToken(r, intentActiveObligations) {
		writeError(w, http.StatusUnauthorized, "SERVICE_AUTH_FAILED", "authentication", "Service authentication failed.")
		return
	}
	relatedType := r.URL.Query().Get("relatedType")
	relatedID := r.URL.Query().Get("relatedId")
	if relatedType == "" || relatedID == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The obligation query input is invalid.")
		return
	}
	tasks, err := h.repo.ActiveObligations(r.Context(), relatedType, relatedID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	obligations := make([]map[string]any, 0, len(tasks))
	for _, task := range tasks {
		obligations = append(obligations, map[string]any{
			"type":         "open_task",
			"id":           task.ID,
			"service":      "work-service",
			"status":       task.Status,
			"dueDate":      domain.FormatDate(task.DueDate),
			"ownerDisplay": task.OwnerID,
			"blocking":     true,
			"safeMessage":  task.Title,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"obligations": obligations})
}

func actorFromRequest(r *http.Request) actorContext {
	return actorContext{ID: r.Header.Get("X-Actor-User-Id"), Role: r.Header.Get("X-Actor-Role")}
}

func canActForOwner(actor actorContext, ownerID string) bool {
	return actor.Role != "Sales" || actor.ID == ownerID
}

func (h *WorkHandler) verifyServiceToken(r *http.Request, intent string) bool {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") || len(h.config.ServiceTokenSecret) == 0 {
		return false
	}
	claims, err := authz.VerifyServiceToken(strings.TrimPrefix(authHeader, "Bearer "), h.config.ServiceID, intent, h.config.ServiceTokenSecret, time.Now().UTC())
	if err != nil {
		return false
	}
	return r.Header.Get("X-Service-Id") == claims.Issuer && r.Header.Get("X-Intent") == intent
}
