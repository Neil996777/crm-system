package handler

import "testing"

func TestHistoryQueryPermission(t *testing.T) {
	_, serviceDB := newAuditTestDB(t)
	app := NewAuditServer(serviceDB, Config{ServiceID: "audit-history", ServiceTokenSecret: []byte("audit-secret")})
	token := signTestServiceToken(t, []byte("audit-secret"), "lead", "audit-history", "audit.append")
	appendEvent(t, app, token, map[string]any{
		"eventId":       "EVT-OWNER-CHANGED",
		"eventVersion":  1,
		"surfaces":      []string{"record_history"},
		"action":        "Owner changed",
		"resourceType":  "Lead",
		"resourceId":    "lead-1",
		"result":        "success",
		"safeSummary":   "Owner changed",
		"acceptanceIds": []string{"ACC-014"},
	}, actorHeaders("sales-1", "Sales", "Sales One"))

	owner := getHistory(app, "/records/Lead/lead-1/history?actorId=sales-1&actorRole=Sales&ownerId=sales-1")
	if owner.Code != 200 {
		t.Fatalf("TEST-HISTORY-001 expected owner history 200, got %d body=%s", owner.Code, owner.Body.String())
	}

	denied := getHistory(app, "/records/Lead/lead-1/history?actorId=sales-2&actorRole=Sales&ownerId=sales-1")
	if denied.Code != 403 {
		t.Fatalf("TEST-HISTORY-003 expected non-owner denial, got %d body=%s", denied.Code, denied.Body.String())
	}
	requireErrorCode(t, denied, "PERMISSION_DENIED")
}
