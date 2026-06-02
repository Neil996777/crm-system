package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"crm-system/services/account/internal/authz"
	"crm-system/services/account/internal/domain"
	"crm-system/services/account/internal/event"
	"crm-system/services/account/internal/repo"
)

const intentWorkActiveObligations = "work.active_obligations"

type ArchiveObligation struct {
	Type         string `json:"type"`
	ID           string `json:"id"`
	Service      string `json:"service"`
	Status       string `json:"status"`
	DueDate      string `json:"dueDate"`
	OwnerDisplay string `json:"ownerDisplay"`
	Blocking     bool   `json:"blocking"`
	SafeMessage  string `json:"safeMessage"`
}

func (h *AccountHandler) archiveEligibility(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	account, err := h.authorizedArchiveAccount(r, actor)
	if err != nil {
		writeArchiveLookupError(w, err)
		return
	}
	obligations, err := h.accountArchiveObligations(r, account.ID)
	if err != nil {
		writeError(w, http.StatusBadGateway, "DEPENDENCY_UNAVAILABLE", "dependency", "A required service is unavailable.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"resourceType":  "account",
		"resourceId":    account.ID,
		"canArchive":    len(obligations) == 0,
		"recordVersion": account.Version,
		"obligations":   obligations,
	})
}

func (h *AccountHandler) archiveAccount(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		Reason          string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 || strings.TrimSpace(request.Reason) == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The archive input is invalid.")
		return
	}
	account, err := h.authorizedArchiveAccount(r, actor)
	if err != nil {
		writeArchiveLookupError(w, err)
		return
	}
	if account.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	obligations, err := h.accountArchiveObligations(r, account.ID)
	if err != nil {
		writeError(w, http.StatusBadGateway, "DEPENDENCY_UNAVAILABLE", "dependency", "A required service is unavailable.")
		return
	}
	if len(obligations) > 0 {
		writeJSON(w, http.StatusConflict, map[string]any{
			"error": map[string]any{
				"code":        "ARCHIVE_BLOCKED_ACTIVE_OBLIGATION",
				"category":    "business_rule",
				"safeMessage": "Active obligations must be resolved before archive.",
			},
			"obligations": obligations,
		})
		return
	}
	archived, err := h.repo.Archive(r.Context(), account.ID, request.ExpectedVersion, actor.ID, request.Reason)
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	_ = h.outbox.Append(r.Context(), event.AccountArchived, archived.ID, map[string]any{
		"traceability": "TASK-032 ACC-002 ACC-014 CIM-037 CIM-PROC-020 PIM-020 PIM-BEH-024 PSM-003 CONTRACT-005 FLOW-010 TEST-ARCHIVE",
		"actorId":      actor.ID,
		"accountId":    archived.ID,
		"reason":       request.Reason,
	})
	writeJSON(w, http.StatusOK, accountDTO(archived))
}

func (h *AccountHandler) authorizedArchiveAccount(r *http.Request, actor actorContext) (domain.Account, error) {
	if !domain.CanArchiveAccount(actor.Role) {
		return domain.Account{}, errPermissionDenied
	}
	account, err := h.repo.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		return domain.Account{}, repo.ErrNotFound
	}
	if err != nil {
		return domain.Account{}, err
	}
	if !domain.CanReadAccount(actor.ID, actor.Role, account) {
		return domain.Account{}, errPermissionDenied
	}
	return account, nil
}

func (h *AccountHandler) accountArchiveObligations(r *http.Request, accountID string) ([]ArchiveObligation, error) {
	if h.config.WorkServiceURL == "" {
		return nil, nil
	}
	token, err := authz.SignServiceToken(authz.ServiceTokenClaims{
		Issuer:   h.config.ServiceID,
		Audience: "work",
		Intent:   intentWorkActiveObligations,
		Expires:  time.Now().UTC().Add(2 * time.Minute),
	}, h.config.ServiceTokenSecret)
	if err != nil {
		return nil, err
	}
	target := strings.TrimRight(h.config.WorkServiceURL, "/") + "/internal/active-obligations?relatedType=" + url.QueryEscape("Customer") + "&relatedId=" + url.QueryEscape(accountID)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Service-Id", h.config.ServiceID)
	req.Header.Set("X-Intent", intentWorkActiveObligations)
	client := h.config.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("work obligation dependency failed")
	}
	var body struct {
		Obligations []ArchiveObligation `json:"obligations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return body.Obligations, nil
}

var errPermissionDenied = errors.New("permission denied")

func writeArchiveLookupError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repo.ErrNotFound):
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
	case errors.Is(err, errPermissionDenied):
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
	default:
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
	}
}
