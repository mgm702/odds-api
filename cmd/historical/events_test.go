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

func TestHistoricalEvents_Decode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		withQuotaHeaders(w)
		w.Write([]byte(`{
			"timestamp": "2024-11-30T12:00:00Z",
			"previous_timestamp": "2024-11-30T11:55:00Z",
			"next_timestamp": "2024-11-30T12:05:00Z",
			"data": [
				{"id":"evt001","sport_key":"nba","sport_title":"NBA","commence_time":"2024-12-01T00:00:00Z","home_team":"LAL","away_team":"BOS"}
			]
		}`))
	}))
	defer srv.Close()

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("date", "2024-11-30T12:00:00Z")

	resp, err := c.Get(context.Background(), "/v4/historical/sports/basketball_nba/events", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := client.Decode[model.HistoricalResponse[[]model.Event]](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(data.Data) != 1 {
		t.Errorf("expected 1 event, got %d", len(data.Data))
	}
	if data.Data[0].ID != "evt001" {
		t.Errorf("expected evt001, got %s", data.Data[0].ID)
	}
}

func TestHistoricalEvents_DateParam(t *testing.T) {
	var gotDate string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotDate = r.URL.Query().Get("date")
		withQuotaHeaders(w)
		w.Write([]byte(`{"timestamp":"2024-01-01T00:00:00Z","data":[]}`))
	}))
	defer srv.Close()

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("date", "2024-06-15T18:00:00Z")

	resp, err := c.Get(context.Background(), "/v4/historical/sports/basketball_nba/events", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	if gotDate != "2024-06-15T18:00:00Z" {
		t.Errorf("expected date=2024-06-15T18:00:00Z, got %s", gotDate)
	}
}
