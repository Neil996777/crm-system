package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"crm-system/services/account/internal/authz"
	"crm-system/services/account/internal/domain"
	"crm-system/services/account/internal/event"
	"crm-system/services/account/internal/repo"
)

const intentCreateAccountForLeadConversion = "account.create_for_lead_conversion"

type Config struct {
	ServiceID          string
	ServiceTokenSecret []byte
	WorkServiceURL     string
	HTTPClient         *http.Client
}

type AccountHandler struct {
	db         *sql.DB
	repo       *repo.AccountRepo
	contacts   *repo.ContactRepo
	duplicates *repo.DuplicateRepo
	outbox     *event.Outbox
	config     Config
}

func NewAccountServer(db *sql.DB, config Config) http.Handler {
	if config.ServiceID == "" {
		config.ServiceID = "account"
	}
	handler := &AccountHandler{db: db, repo: repo.NewAccountRepo(db), contacts: repo.NewContactRepo(db), duplicates: repo.NewDuplicateRepo(db), outbox: event.NewOutbox(db), config: config}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /accounts", handler.createAccount)
	mux.HandleFunc("POST /accounts/duplicate-checks", handler.duplicateCheck)
	mux.HandleFunc("GET /accounts/{id}/archive-eligibility", handler.archiveEligibility)
	mux.HandleFunc("POST /accounts/{id}/archive", handler.archiveAccount)
	mux.HandleFunc("POST /duplicate-checks", handler.duplicateCheck)
	mux.HandleFunc("POST /internal/accounts", handler.createAccountFromLeadConversion)
	mux.HandleFunc("PATCH /accounts/{id}", handler.updateAccount)
	mux.HandleFunc("POST /accounts/{id}/contacts", handler.createContactForAccount)
	mux.HandleFunc("GET /accounts/{id}/contacts", handler.listContactsForAccount)
	mux.HandleFunc("POST /contacts", handler.createContactWithoutAccount)
	mux.HandleFunc("GET /contacts", handler.listContacts)
	mux.HandleFunc("GET /contacts/{id}", handler.getContact)
	mux.HandleFunc("GET /accounts", handler.listAccounts)
	mux.HandleFunc("GET /accounts/{id}", handler.getAccount)
	return mux
}

func (h *AccountHandler) inTransaction(ctx context.Context, fn func(*repo.AccountRepo, *repo.ContactRepo, *event.Outbox) error) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(repo.NewAccountRepoTx(tx), repo.NewContactRepoTx(tx), event.NewOutboxTx(tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func writeOutboxFailure(w http.ResponseWriter) {
	writeError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "dependency", "The audit event could not be persisted.")
}

func (h *AccountHandler) createAccountFromLeadConversion(w http.ResponseWriter, r *http.Request) {
	if !h.verifyServiceToken(r, intentCreateAccountForLeadConversion) {
		writeError(w, http.StatusUnauthorized, "SERVICE_AUTH_FAILED", "authentication", "Service authentication failed.")
		return
	}
	h.createAccount(w, r)
}

func (h *AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		CompanyName         string `json:"companyName"`
		CustomerStatus      string `json:"customerStatus"`
		OwnerID             string `json:"ownerId"`
		ProceedWarningToken string `json:"proceedWarningToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The account input is invalid.")
		return
	}
	if !domain.CanCreateAccount(actor.ID, actor.Role, request.OwnerID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	account, err := domain.NewAccount(domain.Account{
		CompanyName:    request.CompanyName,
		CustomerStatus: request.CustomerStatus,
		OwnerID:        request.OwnerID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The account input is invalid.")
		return
	}
	if request.ProceedWarningToken != "" {
		candidate := domain.DuplicateCandidate{TargetType: "account", CompanyName: account.CompanyName}
		signature, _ := domain.DuplicateSignature(candidate)
		if err := h.duplicates.ConsumeToken(r.Context(), request.ProceedWarningToken, candidate.TargetType, actor.ID, signature); err != nil {
			writeDuplicateTokenError(w, err)
			return
		}
	}
	var created domain.Account
	if err := h.inTransaction(r.Context(), func(txAccounts *repo.AccountRepo, _ *repo.ContactRepo, txOutbox *event.Outbox) error {
		var err error
		created, err = txAccounts.Create(r.Context(), account)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.AccountCreated, created.ID, map[string]any{
			"traceability":   "TASK-010 ACC-005 PIM-005 PIM-BEH-007 PSM-003 CONTRACT-005 CONTRACT-006",
			"actorId":        actor.ID,
			"actorRole":      actor.Role,
			"actorDisplay":   actor.ID,
			"accountId":      created.ID,
			"ownerId":        created.OwnerID,
			"customerStatus": created.CustomerStatus,
		})
	}); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The account input is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, accountDTO(created))
}

func (h *AccountHandler) updateAccount(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		CompanyName     string `json:"companyName"`
		CustomerStatus  string `json:"customerStatus"`
		OwnerID         string `json:"ownerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The account input is invalid.")
		return
	}
	current, err := h.repo.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if !domain.CanEditAccount(actor.ID, actor.Role, current, request.OwnerID) {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	if current.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	updated, err := domain.UpdateAccount(current, request.CompanyName, request.CustomerStatus, request.OwnerID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The account input is invalid.")
		return
	}
	var saved domain.Account
	err = h.inTransaction(r.Context(), func(txAccounts *repo.AccountRepo, _ *repo.ContactRepo, txOutbox *event.Outbox) error {
		var err error
		saved, err = txAccounts.Update(r.Context(), current.ID, request.ExpectedVersion, updated)
		if err != nil {
			return err
		}
		if err := txOutbox.Append(r.Context(), event.AccountUpdated, saved.ID, map[string]any{
			"traceability":   "TASK-010 ACC-005 PIM-BEH-007 CONTRACT-005 CONTRACT-006 CONTRACT-020",
			"actorId":        actor.ID,
			"actorRole":      actor.Role,
			"actorDisplay":   actor.ID,
			"accountId":      saved.ID,
			"customerStatus": saved.CustomerStatus,
			"version":        saved.Version,
		}); err != nil {
			return err
		}
		if current.OwnerID != saved.OwnerID {
			return txOutbox.Append(r.Context(), event.OwnerChanged, saved.ID, map[string]any{
				"traceability": "TASK-010 ACC-005 ACC-014 PIM-BEH-007 CONTRACT-006 EVT-OWNER-CHANGED",
				"actorId":      actor.ID,
				"actorRole":    actor.Role,
				"actorDisplay": actor.ID,
				"oldOwnerId":   current.OwnerID,
				"newOwnerId":   saved.OwnerID,
			})
		}
		return nil
	})
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeOutboxFailure(w)
		return
	}
	writeJSON(w, http.StatusOK, accountDTO(saved))
}

func (h *AccountHandler) verifyServiceToken(r *http.Request, intent string) bool {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}
	claims, err := authz.VerifyServiceToken(strings.TrimPrefix(authHeader, "Bearer "), authz.VerifyOptions{
		Secret:   h.config.ServiceTokenSecret,
		Audience: h.config.ServiceID,
		Intent:   intent,
		Now:      time.Now().UTC(),
	})
	if err != nil {
		return false
	}
	return r.Header.Get("X-Service-Id") == claims.Issuer && r.Header.Get("X-Intent") == intent
}
