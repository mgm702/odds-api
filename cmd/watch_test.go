package cmd

import (
	"context"
	"net/http"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/tui"
)

func TestBuildFetchFunc_Odds(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/odds.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	fetch := buildFetchFunc(c, "odds", "basketball_nba", "us", "h2h", "", "")
	data, err := fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data.OddsEvents) != 1 {
		t.Errorf("expected 1 odds event, got %d", len(data.OddsEvents))
	}
}

func TestBuildFetchFunc_Scores(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/scores.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	fetch := buildFetchFunc(c, "scores", "basketball_nba", "", "", "", "")
	data, err := fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data.ScoreEvents) != 2 {
		t.Errorf("expected 2 score events, got %d", len(data.ScoreEvents))
	}
}

func TestBuildFetchFunc_Error(t *testing.T) {
	srv, _ := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`server error`))
	})

	c := client.New("test-key")
	c.BaseURL = srv.URL

	fetch := buildFetchFunc(c, "odds", "basketball_nba", "us", "", "", "")
	_, err := fetch(context.Background())
	if err == nil {
		t.Error("expected error from fetch")
	}
}

func TestBuildFetchFunc_QuotaTracking(t *testing.T) {
	srv, _ := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Requests-Remaining", "450")
		w.Header().Set("X-Requests-Used", "50")
		w.Header().Set("X-Requests-Last", "2")
		w.Write([]byte(`[]`))
	})

	c := client.New("test-key")
	c.BaseURL = srv.URL

	fetch := buildFetchFunc(c, "odds", "basketball_nba", "us", "", "", "")
	data, err := fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Quota.Remaining != 450 {
		t.Errorf("expected remaining=450, got %d", data.Quota.Remaining)
	}
}

func TestNewWatchModel_OddsMode(t *testing.T) {
	m := tui.NewWatchModel("odds", 60, nil, "")
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestNewWatchModel_ScoresMode(t *testing.T) {
	m := tui.NewWatchModel("scores", 60, nil, "")
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}
