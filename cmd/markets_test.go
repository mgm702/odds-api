package cmd

import (
	"context"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestMarketsCommand_Decode(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/markets.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/events/evt001/markets", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	em, err := client.Decode[model.EventMarkets](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if em.ID != "evt001" {
		t.Errorf("expected evt001, got %s", em.ID)
	}
	if len(em.Bookmakers) != 1 {
		t.Fatalf("expected 1 bookmaker, got %d", len(em.Bookmakers))
	}
	if len(em.Bookmakers[0].Markets) != 4 {
		t.Errorf("expected 4 markets, got %d", len(em.Bookmakers[0].Markets))
	}
}
