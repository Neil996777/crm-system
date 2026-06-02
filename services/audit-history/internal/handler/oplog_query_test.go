package handler

import "testing"

func TestOperationLogAdministratorOnlyAndReadOnly(t *testing.T) {
	_, serviceDB := newAuditTestDB(t)
	app := NewAuditServer(serviceDB, Config{ServiceID: "audit-history", ServiceTokenSecret: []byte("audit-secret")})
	token := signTestServiceToken(t, []byte("audit-secret"), "identity-authz", "audit-history", "audit.append")
	appendEvent(t, app, token, map[string]any{
		"eventId":       "EVT-AUTH-ACCESS-DENIED",
		"eventVersion":  1,
		"surfaces":      []string{"operation_log"},
		"action":        "Access denied",
		"resourceType":  "Auth",
		"resourceId":    "auth",
		"result":        "denied",
		"safeSummary":   "Access denied",
		"acceptanceIds": []string{"ACC-022"},
	}, actorHeaders("admin-1", "Administrator", "Admin One"))

	admin := getHistory(app, "/operation-log?actorId=admin-1&actorRole=Administrator")
	if admin.Code != 200 {
		t.Fatalf("TEST-OPLOG-001 expected admin oplog 200, got %d body=%s", admin.Code, admin.Body.String())
	}

	manager := getHistory(app, "/operation-log?actorId=mgr-1&actorRole=Sales%20Manager")
	if manager.Code != 403 {
		t.Fatalf("TEST-OPLOG-003 expected manager denial, got %d body=%s", manager.Code, manager.Body.String())
	}
	requireErrorCode(t, manager, "PERMISSION_DENIED")

	deleted := httptestDelete(app, "/operation-log/event-1", nil)
	if deleted.Code != 404 && deleted.Code != 405 {
		t.Fatalf("TEST-OPLOG-005 expected no edit/delete route, got %d body=%s", deleted.Code, deleted.Body.String())
	}
}
