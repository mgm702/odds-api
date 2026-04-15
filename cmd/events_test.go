package cmd

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestEventsCommand_Decode(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/events.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/events", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events, err := client.Decode[[]model.Event](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
	if events[0].HomeTeam != "Los Angeles Lakers" {
		t.Errorf("expected Lakers, got %s", events[0].HomeTeam)
	}
}

func TestEventsCommand_Filters(t *testing.T) {
	var gotParams url.Values
	srv, _ := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotParams = r.URL.Query()
		withQuotaHeaders(w)
		_, _ = w.Write([]byte(`[]`))
	})

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("eventIds", "evt001")
	params.Set("commenceTimeFrom", "2024-01-01T00:00:00Z")
	params.Set("commenceTimeTo", "2024-12-31T23:59:59Z")

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/events", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	if gotParams.Get("eventIds") != "evt001" {
		t.Errorf("expected eventIds=evt001, got %s", gotParams.Get("eventIds"))
	}
	if gotParams.Get("commenceTimeFrom") != "2024-01-01T00:00:00Z" {
		t.Error("expected commenceTimeFrom to be set")
	}
}
