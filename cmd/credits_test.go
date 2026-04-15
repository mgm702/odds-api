package cmd

import (
	"context"
	"net/http"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
)

func TestCreditsCommand_QuotaOnly(t *testing.T) {
	srv, _ := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Requests-Remaining", "358")
		w.Header().Set("X-Requests-Used", "142")
		w.Header().Set("X-Requests-Last", "0")
		_, _ = w.Write([]byte(`[]`))
	})

	c := client.New("test-key")
	c.BaseURL = srv.URL

	quota, err := c.GetQuotaOnly(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if quota.Used != 142 {
		t.Errorf("expected used=142, got %d", quota.Used)
	}
	if quota.Remaining != 358 {
		t.Errorf("expected remaining=358, got %d", quota.Remaining)
	}
	if quota.LastCost != 0 {
		t.Errorf("expected last_cost=0, got %d", quota.LastCost)
	}
}

func TestCreditsCommand_PathIsCorrect(t *testing.T) {
	var gotPath string
	srv, _ := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("X-Requests-Remaining", "500")
		w.Header().Set("X-Requests-Used", "0")
		w.Header().Set("X-Requests-Last", "0")
		_, _ = w.Write([]byte(`[]`))
	})

	c := client.New("test-key")
	c.BaseURL = srv.URL

	_, err := c.GetQuotaOnly(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotPath != "/v4/sports" {
		t.Errorf("expected /v4/sports, got %s", gotPath)
	}
}
