package cmd

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestScoresCommand_Decode(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/scores.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/scores", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	scores, err := client.Decode[[]model.ScoreEvent](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(scores) != 2 {
		t.Fatalf("expected 2 score events, got %d", len(scores))
	}

	if scores[0].Scores == nil {
		t.Error("expected non-nil scores for first event")
	}
	if scores[0].Scores[0].Score != "55" {
		t.Errorf("expected score 55, got %s", scores[0].Scores[0].Score)
	}

	if scores[1].Scores != nil {
		t.Error("expected nil scores for upcoming event")
	}
}

func TestScoresCommand_DaysFrom(t *testing.T) {
	var gotParams url.Values
	srv, _ := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotParams = r.URL.Query()
		withQuotaHeaders(w)
		w.Write([]byte(`[]`))
	})

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("daysFrom", "2")

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/scores", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	if gotParams.Get("daysFrom") != "2" {
		t.Errorf("expected daysFrom=2, got %s", gotParams.Get("daysFrom"))
	}
}
