package event

import "testing"

func TestCommercialAuditCatalogEventIDs(t *testing.T) {
	t.Run("TEST-EVT-CATALOG-COMMERCIAL-002 quote accepted", func(t *testing.T) {
		body := auditAppendBody(outboxEvent{
			ID:          "evt_catalog_quote_accepted",
			EventType:   QuoteAccepted,
			AggregateID: "quote_1",
			Payload:     map[string]any{"result": "success"},
		})
		if body["eventId"] != "EVT-QUOTE-ACCEPTED" {
			t.Fatalf("expected EVT-QUOTE-ACCEPTED, got %#v body=%#v", body["eventId"], body)
		}
	})

	t.Run("TEST-EVT-CATALOG-COMMERCIAL-003 payment recorded", func(t *testing.T) {
		body := auditAppendBody(outboxEvent{
			ID:          "evt_catalog_payment_recorded",
			EventType:   PaymentRecorded,
			AggregateID: "payment_1",
			Payload:     map[string]any{"result": "success"},
		})
		if body["eventId"] != "EVT-PAYMENT-RECORDED" {
			t.Fatalf("expected EVT-PAYMENT-RECORDED, got %#v body=%#v", body["eventId"], body)
		}
	})

	t.Run("TEST-EVT-CATALOG-COMMERCIAL-001 contract terminated", func(t *testing.T) {
		body := auditAppendBody(outboxEvent{
			ID:          "evt_catalog_contract_terminated",
			EventType:   ContractStatusChanged,
			AggregateID: "contract_1",
			Payload:     map[string]any{"toStatus": "Terminated", "terminationReason": "customer_request"},
		})
		if body["eventId"] != "EVT-CONTRACT-TERMINATED" {
			t.Fatalf("expected EVT-CONTRACT-TERMINATED, got %#v body=%#v", body["eventId"], body)
		}
	})

	t.Run("TEST-EVT-CATALOG-COMMERCIAL-004 archived", func(t *testing.T) {
		body := auditAppendBody(outboxEvent{
			ID:          "evt_catalog_payment_plan_archived",
			EventType:   PaymentPlanArchived,
			AggregateID: "payment_plan_1",
			Payload:     map[string]any{"result": "success"},
		})
		if body["eventId"] != "EVT-RECORD-ARCHIVED" {
			t.Fatalf("expected EVT-RECORD-ARCHIVED, got %#v body=%#v", body["eventId"], body)
		}
	})
}
