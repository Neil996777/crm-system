package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuditClientPropagatesCorrelationHeader(t *testing.T) {
	var calls int
	var eventIDs []string
	audit := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Header.Get("X-Correlation-Id") != "corr-import-export" {
			t.Fatalf("missing audit correlation id: %q", r.Header.Get("X-Correlation-Id"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode audit body: %v", err)
		}
		eventIDs = append(eventIDs, body["eventId"].(string))
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(audit.Close)
	client := NewAuditClient(audit.URL, "import-export", []byte("secret"), audit.Client())
	if err := client.AppendImportRun(context.Background(), ImportRunLogInput{
		RunID:         "import-corr",
		ActorID:       "mgr-1",
		ActorRole:     "Sales Manager",
		ObjectType:    "lead",
		TotalRows:     1,
		SuccessCount:  1,
		Result:        "Completed",
		CorrelationID: "corr-import-export",
	}); err != nil {
		t.Fatalf("append import run: %v", err)
	}
	if err := client.AppendExportRun(context.Background(), ExportRunLogInput{
		RunID:         "export-corr",
		ActorID:       "mgr-1",
		ActorRole:     "Sales Manager",
		ObjectType:    "lead",
		ExportedCount: 1,
		Result:        "Completed",
		CorrelationID: "corr-import-export",
	}); err != nil {
		t.Fatalf("append export run: %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected two audit calls, got %d", calls)
	}
	if len(eventIDs) != 2 || eventIDs[0] != "EVT-IMPORT-RUN" || eventIDs[1] != "EVT-EXPORT-RUN" {
		t.Fatalf("TEST-EVT-CATALOG-IMPORTEXPORT-001 expected import/export catalog event ids, got %#v", eventIDs)
	}
}
