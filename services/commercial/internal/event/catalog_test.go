package event

import "testing"

func TestCommercialAuditCatalogEventIDs(t *testing.T) {
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
}
