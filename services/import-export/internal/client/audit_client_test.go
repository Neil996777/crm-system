package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuditClientPropagatesCorrelationHeader(t *testing.T) {
	var calls int
	audit := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Header.Get("X-Correlation-Id") != "corr-import-export" {
			t.Fatalf("missing audit correlation id: %q", r.Header.Get("X-Correlation-Id"))
		}
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
}
