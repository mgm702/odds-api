package cmd

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestOddsCommand_Decode(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/odds.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("regions", "us")

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/odds", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events, err := client.Decode[[]model.OddsEvent](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if len(events[0].Bookmakers) != 2 {
		t.Errorf("expected 2 bookmakers, got %d", len(events[0].Bookmakers))
	}
}

func TestOddsCommand_Params(t *testing.T) {
	var gotParams url.Values
	srv, _ := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotParams = r.URL.Query()
		withQuotaHeaders(w)
		_, _ = w.Write([]byte(`[]`))
	})

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("regions", "us,uk")
	params.Set("markets", "h2h,spreads")
	params.Set("oddsFormat", "american")

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/odds", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	if gotParams.Get("regions") != "us,uk" {
		t.Errorf("expected regions=us,uk, got %s", gotParams.Get("regions"))
	}
	if gotParams.Get("markets") != "h2h,spreads" {
		t.Errorf("expected markets=h2h,spreads, got %s", gotParams.Get("markets"))
	}
}

func TestOddsCommand_WithSpreads(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/odds.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("regions", "us")

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/odds", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events, err := client.Decode[[]model.OddsEvent](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	dk := events[0].Bookmakers[1]
	if len(dk.Markets) != 2 {
		t.Fatalf("expected 2 markets from DraftKings, got %d", len(dk.Markets))
	}
	spreads := dk.Markets[1]
	if spreads.Key != "spreads" {
		t.Errorf("expected spreads market, got %s", spreads.Key)
	}
	if spreads.Outcomes[0].Point == nil {
		t.Error("expected point for spread outcome")
	}
}
