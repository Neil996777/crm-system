package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"crm-system/services/commercial/internal/domain"
	"crm-system/services/commercial/internal/event"
	"crm-system/services/commercial/internal/repo"
)

type Config struct {
	ServiceID          string
	ServiceTokenSecret []byte
}

type CommercialHandler struct {
	db        *sql.DB
	quotes    *repo.QuoteRepo
	contracts *repo.ContractRepo
	payments  *repo.PaymentRepo
	outbox    *event.Outbox
	config    Config
}

func NewCommercialServer(db *sql.DB, config Config) http.Handler {
	if config.ServiceID == "" {
		config.ServiceID = "commercial"
	}
	handler := &CommercialHandler{db: db, quotes: repo.NewQuoteRepo(db), contracts: repo.NewContractRepo(db), payments: repo.NewPaymentRepo(db), outbox: event.NewOutbox(db), config: config}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /quotes", handler.listQuotes)
	mux.HandleFunc("POST /quotes", handler.createQuote)
	mux.HandleFunc("POST /quotes/{id}/status", handler.changeQuoteStatus)
	mux.HandleFunc("GET /quotes/{id}", handler.getQuote)
	mux.HandleFunc("GET /contracts", handler.listContracts)
	mux.HandleFunc("POST /contracts", handler.createContract)
	mux.HandleFunc("POST /contracts/{id}/status", handler.changeContractStatus)
	mux.HandleFunc("POST /contracts/{id}/payment-plans", handler.createPaymentPlan)
	mux.HandleFunc("POST /contracts/{id}/payments", handler.recordPayment)
	mux.HandleFunc("GET /contracts/{id}/archive-eligibility", handler.contractArchiveEligibility)
	mux.HandleFunc("POST /contracts/{id}/archive", handler.archiveContract)
	mux.HandleFunc("POST /payment-plans/{id}/archive", handler.archivePaymentPlan)
	mux.HandleFunc("GET /contracts/{id}", handler.getContract)
	mux.HandleFunc("GET /internal/contracts/{id}/signed-status", handler.getContractSignedStatus)
	mux.HandleFunc("GET /internal/reminders/eligibility", handler.getReminderEligibility)
	return mux
}

func (h *CommercialHandler) inTransaction(ctx context.Context, fn func(*repo.QuoteRepo, *repo.ContractRepo, *repo.PaymentRepo, *event.Outbox) error) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(repo.NewQuoteRepoTx(tx), repo.NewContractRepoTx(tx), repo.NewPaymentRepoTx(tx), event.NewOutboxTx(tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func writeOutboxFailure(w http.ResponseWriter) {
	writeError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "dependency", "The audit event could not be persisted.")
}

func (h *CommercialHandler) createQuote(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		OpportunityID string `json:"opportunityId"`
		CustomerID    string `json:"customerId"`
		Amount        string `json:"amount"`
		Status        string `json:"status"`
		ValidityEnd   string `json:"validityEnd"`
		OwnerID       string `json:"ownerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The quote input is invalid.")
		return
	}
	validityEnd, err := domain.ParseDate(request.ValidityEnd)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The quote input is invalid.")
		return
	}
	if actor.Role == "Sales" && request.OwnerID != actor.ID {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	exists, err := h.quotes.ExistsForOpportunity(r.Context(), request.OpportunityID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if err := domain.EnsureCanCreateQuote(exists); err != nil {
		writeError(w, http.StatusConflict, "QUOTE_ALREADY_EXISTS", "conflict", "A quote already exists for this opportunity.")
		return
	}
	quote, err := domain.NewQuote(domain.Quote{
		OpportunityID: request.OpportunityID,
		CustomerID:    request.CustomerID,
		Amount:        request.Amount,
		Status:        request.Status,
		ValidityEnd:   validityEnd,
		OwnerID:       request.OwnerID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The quote input is invalid.")
		return
	}
	var created domain.Quote
	err = h.inTransaction(r.Context(), func(txQuotes *repo.QuoteRepo, _ *repo.ContractRepo, _ *repo.PaymentRepo, txOutbox *event.Outbox) error {
		var err error
		created, err = txQuotes.Create(r.Context(), quote)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), event.QuoteCreated, created.ID, map[string]any{
			"traceability":  "TASK-017 ACC-009 PIM-008 PIM-SM-004 PSM-005 CONTRACT-009 CONTRACT-010",
			"actorId":       actor.ID,
			"actorRole":     actor.Role,
			"actorDisplay":  actor.ID,
			"quoteId":       created.ID,
			"opportunityId": created.OpportunityID,
		})
	})
	if errors.Is(err, domain.ErrQuoteAlreadyExists) {
		writeError(w, http.StatusConflict, "QUOTE_ALREADY_EXISTS", "conflict", "A quote already exists for this opportunity.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The quote input is invalid.")
		return
	}
	writeJSON(w, http.StatusCreated, quoteDTO(created))
}

func (h *CommercialHandler) changeQuoteStatus(w http.ResponseWriter, r *http.Request) {
	actor := actorFromRequest(r)
	var request struct {
		ExpectedVersion int    `json:"expectedVersion"`
		ToStatus        string `json:"toStatus"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ExpectedVersion == 0 || request.ToStatus == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The quote status input is invalid.")
		return
	}
	current, err := h.quotes.Find(r.Context(), r.PathValue("id"))
	if errors.Is(err, repo.ErrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "not_found", "The requested resource was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	if actor.Role == "Sales" && current.OwnerID != actor.ID {
		writeError(w, http.StatusForbidden, "PERMISSION_DENIED", "permission", "Permission denied.")
		return
	}
	if current.Version != request.ExpectedVersion {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err := domain.ValidateQuoteStatusTransition(current.Status, request.ToStatus); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_TRANSITION", "business_rule", "The requested quote status transition is not allowed.")
		return
	}
	eventType := event.QuoteStatusChanged
	trace := "TASK-017 ACC-009 PIM-SM-004 PIM-BEH-012 CONTRACT-009 CONTRACT-010"
	if request.ToStatus == domain.StatusAccepted {
		eventType = event.QuoteAccepted
		trace = "TASK-017 ACC-009 ACC-014 ACC-022 PIM-BEH-013 PIM-INV-012 PSM-005 CONTRACT-009 CONTRACT-010 EVT-QUOTE-ACCEPTED"
	}
	var updated domain.Quote
	err = h.inTransaction(r.Context(), func(txQuotes *repo.QuoteRepo, _ *repo.ContractRepo, _ *repo.PaymentRepo, txOutbox *event.Outbox) error {
		var err error
		updated, err = txQuotes.ChangeStatus(r.Context(), current.ID, request.ExpectedVersion, request.ToStatus)
		if err != nil {
			return err
		}
		return txOutbox.Append(r.Context(), eventType, updated.ID, map[string]any{
			"traceability": trace,
			"actorId":      actor.ID,
			"actorRole":    actor.Role,
			"actorDisplay": actor.ID,
			"quoteId":      updated.ID,
			"fromStatus":   current.Status,
			"toStatus":     updated.Status,
		})
	})
	if errors.Is(err, repo.ErrVersionConflict) {
		writeError(w, http.StatusConflict, "VERSION_CONFLICT", "conflict", "The record changed after it was opened.")
		return
	}
	if err != nil {
		writeOutboxFailure(w)
		return
	}
	writeJSON(w, http.StatusOK, quoteDTO(updated))
}
