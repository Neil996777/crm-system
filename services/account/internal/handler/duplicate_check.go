package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"crm-system/services/account/internal/domain"
	"crm-system/services/account/internal/event"
	"crm-system/services/account/internal/repo"
)

func (h *AccountHandler) duplicateCheck(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		TargetType string `json:"targetType"`
		Candidate  struct {
			CompanyName string `json:"companyName"`
			Email       string `json:"email"`
			Phone       string `json:"phone"`
		} `json:"candidate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The duplicate check input is invalid.")
		return
	}
	if request.TargetType != "account" && request.TargetType != "contact" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The duplicate check input is invalid.")
		return
	}
	candidate := domain.DuplicateCandidate{
		TargetType:  request.TargetType,
		CompanyName: request.Candidate.CompanyName,
		Email:       request.Candidate.Email,
		Phone:       request.Candidate.Phone,
	}
	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The duplicate check input is invalid.")
		return
	}
	defer tx.Rollback()
	txDuplicates := repo.NewDuplicateRepoTx(tx)
	txOutbox := event.NewOutboxTx(tx)
	result, err := txDuplicates.Check(r.Context(), actor.ID, actor.Role, candidate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The duplicate check input is invalid.")
		return
	}
	if result.Result == "PossibleDuplicate" {
		if err := txOutbox.Append(r.Context(), event.DuplicateWarningRaised, result.WarningToken, map[string]any{
			"traceability":     "TASK-031 ACC-019 CIM-040 CIM-PROC-005 PIM-021 PIM-BEH-025 PSM-003 CONTRACT-005 FLOW-012 TEST-DUPLICATE-WARN",
			"actorId":          actor.ID,
			"actorRole":        actor.Role,
			"actorDisplay":     actor.ID,
			"targetType":       request.TargetType,
			"normalizedFields": result.NormalizedFields,
			"rules":            result.Rules,
		}); err != nil {
			writeOutboxFailure(w)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		writeOutboxFailure(w)
		return
	}
	writeJSON(w, http.StatusOK, duplicateDTO(result))
}

func duplicateDTO(result domain.DuplicateCheckResult) map[string]any {
	body := map[string]any{
		"result":           result.Result,
		"normalizedFields": result.NormalizedFields,
		"matches":          duplicateMatchDTOs(result.Matches),
		"rules":            result.Rules,
	}
	if result.WarningToken != "" {
		body["warningToken"] = result.WarningToken
	}
	return body
}

func duplicateMatchDTOs(matches []domain.DuplicateMatch) []map[string]any {
	items := make([]map[string]any, 0, len(matches))
	for _, match := range matches {
		items = append(items, map[string]any{
			"type":          match.Type,
			"matchStrength": match.MatchStrength,
			"safeSummary":   match.SafeSummary,
			"visible":       match.Visible,
			"rule":          match.Rule,
		})
	}
	return items
}

func writeDuplicateTokenError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repo.ErrDuplicateTokenUsed):
		writeError(w, http.StatusConflict, "DUPLICATE_WARNING_TOKEN_USED", "conflict", "The duplicate warning confirmation was already used.")
	case errors.Is(err, repo.ErrDuplicateTokenInvalid):
		writeError(w, http.StatusConflict, "DUPLICATE_WARNING_TOKEN_INVALID", "conflict", "The duplicate warning confirmation is invalid.")
	default:
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
	}
}
