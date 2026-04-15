package cmd

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestEventOddsCommand_Decode(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/event_odds.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("regions", "us")
	params.Set("markets", "player_points")

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/events/evt001/odds", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	event, err := client.Decode[model.OddsEvent](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if event.ID != "evt001" {
		t.Errorf("expected evt001, got %s", event.ID)
	}
	if len(event.Bookmakers) != 1 {
		t.Fatalf("expected 1 bookmaker, got %d", len(event.Bookmakers))
	}
	if event.Bookmakers[0].Markets[0].Key != "player_points" {
		t.Errorf("expected player_points market, got %s", event.Bookmakers[0].Markets[0].Key)
	}
}

func TestEventOddsCommand_Params(t *testing.T) {
	var gotPath string
	var gotParams url.Values
	srv, _ := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotParams = r.URL.Query()
		withQuotaHeaders(w)
		_, _ = w.Write([]byte(`{"id":"evt001","bookmakers":[]}`))
	})

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("regions", "us")
	params.Set("markets", "h2h")
	params.Set("oddsFormat", "american")

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/events/evt001/odds", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	if gotPath != "/v4/sports/basketball_nba/events/evt001/odds" {
		t.Errorf("unexpected path: %s", gotPath)
	}
	if gotParams.Get("oddsFormat") != "american" {
		t.Errorf("expected oddsFormat=american, got %s", gotParams.Get("oddsFormat"))
	}
}
