package handler

import (
	"net/http"

	"crm-system/services/commercial/internal/domain"
)

const intentReminderEligibility = "commercial.reminder_eligibility"

func (h *CommercialHandler) getReminderEligibility(w http.ResponseWriter, r *http.Request) {
	if !h.verifyServiceToken(r, intentReminderEligibility) {
		writeError(w, http.StatusUnauthorized, "SERVICE_AUTH_FAILED", "authentication", "Service authentication failed.")
		return
	}
	actor := actorFromRequest(r)
	businessDate, err := domain.ParseDate(r.URL.Query().Get("businessDate"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "validation", "The reminder query input is invalid.")
		return
	}
	rows, err := h.contracts.PendingSignatureReminderRows(r.Context(), actor.ID, actor.Role, businessDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	paymentRows, err := h.payments.ReminderRows(r.Context(), actor.ID, actor.Role, businessDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "validation", "The request is invalid.")
		return
	}
	rows = append(rows, paymentRows...)
	writeJSON(w, http.StatusOK, map[string]any{"rows": rows})
}
