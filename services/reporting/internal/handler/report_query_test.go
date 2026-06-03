package handler

import (
	"database/sql"
	"net/http"
	"testing"
)

func TestBasicSalesReportAcceptance(t *testing.T) {
	db := newReportingTestDB(t)
	app := NewReportingServer(db, Config{})

	insertProjection(t, db, map[string]any{"recordType": "lead", "recordId": "lead_new_1", "ownerId": "sales-1", "teamId": "single-team", "status": "New"})
	insertProjection(t, db, map[string]any{"recordType": "lead", "recordId": "lead_new_2", "ownerId": "sales-1", "teamId": "single-team", "status": "New"})
	insertProjection(t, db, map[string]any{"recordType": "lead", "recordId": "lead_qualified", "ownerId": "sales-1", "teamId": "single-team", "status": "Qualified"})
	insertProjection(t, db, map[string]any{"recordType": "opportunity", "recordId": "opp_prospecting", "ownerId": "sales-1", "teamId": "single-team", "stage": "Prospecting", "amount": "2000.00"})
	insertProjection(t, db, map[string]any{"recordType": "opportunity", "recordId": "opp_won", "ownerId": "sales-1", "teamId": "single-team", "stage": "Won", "amount": "8000.00"})
	insertProjection(t, db, map[string]any{"recordType": "quote", "recordId": "quote_draft", "ownerId": "sales-1", "teamId": "single-team", "status": "Draft", "amount": "2500.00"})
	insertProjection(t, db, map[string]any{"recordType": "quote", "recordId": "quote_accepted", "ownerId": "sales-1", "teamId": "single-team", "status": "Accepted", "amount": "8000.00"})
	insertProjection(t, db, map[string]any{"recordType": "contract", "recordId": "contract_signed", "ownerId": "sales-1", "teamId": "single-team", "status": "Signed", "amount": "8000.00"})
	insertProjection(t, db, map[string]any{"recordType": "payment", "recordId": "payment_due", "ownerId": "sales-1", "teamId": "single-team", "status": "Due", "amount": "4000.00"})
	insertProjection(t, db, map[string]any{"recordType": "payment", "recordId": "payment_paid", "ownerId": "sales-1", "teamId": "single-team", "status": "Paid", "amount": "1500.00"})
	insertProjection(t, db, map[string]any{"recordType": "lead", "recordId": "lead_other_team", "ownerId": "sales-9", "teamId": "other-team", "status": "New"})
	insertArchivedProjection(t, db, map[string]any{"recordType": "lead", "recordId": "lead_archived", "ownerId": "sales-1", "teamId": "single-team", "status": "New"})
	insertArchivedProjection(t, db, map[string]any{"recordType": "opportunity", "recordId": "opp_archived", "ownerId": "sales-1", "teamId": "single-team", "stage": "Won", "amount": "9000.00"})

	t.Run("TEST-BASIC-REPORT-001 grouping traceability", func(t *testing.T) {
		rec := getReportingJSON(app, "/reports/sales-overview", actorHeaders("mgr-1", "Sales Manager", "single-team"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected manager sales report 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		metrics := body["metrics"].(map[string]any)
		if metrics["leadCount"].(float64) != 3 || metrics["opportunityCount"].(float64) != 2 {
			t.Fatalf("expected active team counts, got %#v", metrics)
		}
		if metrics["quoteAmount"] != "10500.00" || metrics["contractAmount"] != "8000.00" || metrics["paidAmount"] != "1500.00" || metrics["receivableAmount"] != "6500.00" {
			t.Fatalf("expected active team amounts, got %#v", metrics)
		}
		breakdowns := body["breakdowns"].(map[string]any)
		requireGroup(t, breakdowns["leadsByStatus"].([]any), "New", 2, "0.00")
		requireGroup(t, breakdowns["leadsByStatus"].([]any), "Qualified", 1, "0.00")
		requireGroup(t, breakdowns["opportunitiesByStage"].([]any), "Prospecting", 1, "2000.00")
		requireGroup(t, breakdowns["opportunitiesByStage"].([]any), "Won", 1, "8000.00")
		requireGroup(t, breakdowns["quotesByStatus"].([]any), "Accepted", 1, "8000.00")
		requireGroup(t, breakdowns["contractsByStatus"].([]any), "Signed", 1, "8000.00")
		requireGroup(t, breakdowns["paymentsByStatus"].([]any), "Due", 1, "4000.00")
		requireGroup(t, breakdowns["paymentsByStatus"].([]any), "Paid", 1, "1500.00")
	})

	t.Run("TEST-BASIC-REPORT-002 empty data returns zero and empty breakdowns", func(t *testing.T) {
		rec := getReportingJSON(app, "/reports/sales-overview", actorHeaders("mgr-empty", "Sales Manager", "empty-team"))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected empty report 200, got %d body=%s", rec.Code, rec.Body.String())
		}
		body := decodeJSON(t, rec)
		metrics := body["metrics"].(map[string]any)
		if body["emptyState"] != true || metrics["leadCount"].(float64) != 0 || metrics["quoteAmount"] != "0.00" {
			t.Fatalf("expected zero empty state, got %s", rec.Body.String())
		}
		breakdowns := body["breakdowns"].(map[string]any)
		if len(breakdowns["leadsByStatus"].([]any)) != 0 || len(breakdowns["paymentsByStatus"].([]any)) != 0 {
			t.Fatalf("expected empty breakdown arrays, got %#v", breakdowns)
		}
	})

	t.Run("TEST-BASIC-REPORT-003 Sales denied", func(t *testing.T) {
		rec := getReportingJSON(app, "/reports/sales-overview", actorHeaders("sales-1", "Sales", "single-team"))
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected sales denied 403, got %d body=%s", rec.Code, rec.Body.String())
		}
		requireErrorCode(t, rec, "PERMISSION_DENIED")
	})

	t.Run("TEST-BASIC-REPORT-004 unauthorized records excluded before aggregate", func(t *testing.T) {
		rec := getReportingJSON(app, "/reports/sales-overview", actorHeaders("mgr-1", "Sales Manager", "single-team"))
		body := decodeJSON(t, rec)
		metrics := body["metrics"].(map[string]any)
		if metrics["leadCount"].(float64) != 3 || contains(rec.Body.String(), "lead_other_team") {
			t.Fatalf("expected other-team records excluded, got %s", rec.Body.String())
		}
	})

	t.Run("TEST-BASIC-REPORT-005 archived excluded by default", func(t *testing.T) {
		rec := getReportingJSON(app, "/reports/sales-overview", actorHeaders("mgr-1", "Sales Manager", "single-team"))
		body := decodeJSON(t, rec)
		metrics := body["metrics"].(map[string]any)
		if metrics["leadCount"].(float64) != 3 || metrics["wonCount"].(float64) != 1 || metrics["receivableAmount"] != "6500.00" {
			t.Fatalf("expected archived records excluded, got %#v body=%s", metrics, rec.Body.String())
		}
	})
}

func insertArchivedProjection(t *testing.T, db *sql.DB, values map[string]any) {
	t.Helper()
	_, err := db.Exec(`
		INSERT INTO reporting.record_projections
			(source_service, record_type, record_id, owner_id, team_id, status, stage, amount, archived_at, updated_at)
		VALUES ('test-source', $1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), $7::numeric, now(), now())
	`, values["recordType"], values["recordId"], values["ownerId"], values["teamId"], stringValue(values["status"]), stringValue(values["stage"]), amountValue(values["amount"]))
	if err != nil {
		t.Fatalf("insert archived projection: %v", err)
	}
}

func requireGroup(t *testing.T, rows []any, key string, expectedCount int, expectedAmount string) {
	t.Helper()
	for _, raw := range rows {
		row := raw.(map[string]any)
		if row["key"] == key {
			if row["count"].(float64) != float64(expectedCount) || row["amount"] != expectedAmount {
				t.Fatalf("expected group %s count=%d amount=%s, got %#v", key, expectedCount, expectedAmount, row)
			}
			return
		}
	}
	t.Fatalf("missing group %s in %#v", key, rows)
}
