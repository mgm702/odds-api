package historical

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestHistoricalEventOdds_Decode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		withQuotaHeaders(w)
		w.Write([]byte(`{
			"timestamp": "2024-11-30T12:00:00Z",
			"previous_timestamp": "2024-11-30T11:55:00Z",
			"next_timestamp": null,
			"data": {
				"id": "evt001",
				"sport_key": "nba",
				"sport_title": "NBA",
				"commence_time": "2024-12-01T00:00:00Z",
				"home_team": "LAL",
				"away_team": "BOS",
				"bookmakers": []
			}
		}`))
	}))
	defer srv.Close()

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("date", "2024-11-30T12:00:00Z")
	params.Set("regions", "us")
	params.Set("markets", "h2h")

	resp, err := c.Get(context.Background(), "/v4/historical/sports/nba/events/evt001/odds", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := client.Decode[model.HistoricalResponse[model.OddsEvent]](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if data.Data.ID != "evt001" {
		t.Errorf("expected evt001, got %s", data.Data.ID)
	}
	if data.NextTimestamp != nil {
		t.Errorf("expected nil next_timestamp")
	}
}

func TestHistoricalEventOdds_Path(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		withQuotaHeaders(w)
		w.Write([]byte(`{"timestamp":"2024-01-01T00:00:00Z","data":{"id":"e1","bookmakers":[]}}`))
	}))
	defer srv.Close()

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("date", "2024-11-30T12:00:00Z")
	params.Set("regions", "us")
	params.Set("markets", "h2h")

	resp, err := c.Get(context.Background(), "/v4/historical/sports/basketball_nba/events/evt001/odds", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	if gotPath != "/v4/historical/sports/basketball_nba/events/evt001/odds" {
		t.Errorf("unexpected path: %s", gotPath)
	}
}
