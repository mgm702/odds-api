package historical

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func withQuotaHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Requests-Remaining", "500")
	w.Header().Set("X-Requests-Used", "100")
	w.Header().Set("X-Requests-Last", "10")
}

func TestHistoricalOdds_Decode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		withQuotaHeaders(w)
		data, _ := os.ReadFile("../../testdata/historical_odds.json")
		w.Write(data)
	}))
	defer srv.Close()

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("date", "2024-11-30T12:00:00Z")
	params.Set("regions", "us")

	resp, err := c.Get(context.Background(), "/v4/historical/sports/basketball_nba/odds", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := client.Decode[model.HistoricalResponse[[]model.OddsEvent]](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if data.Timestamp != "2024-11-30T12:00:00Z" {
		t.Errorf("expected timestamp, got %s", data.Timestamp)
	}
	if data.PreviousTimestamp == nil {
		t.Error("expected previous_timestamp")
	}
	if data.NextTimestamp == nil {
		t.Error("expected next_timestamp")
	}
	if len(data.Data) != 1 {
		t.Errorf("expected 1 event, got %d", len(data.Data))
	}
}

func TestHistoricalOdds_Params(t *testing.T) {
	var gotParams url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotParams = r.URL.Query()
		withQuotaHeaders(w)
		w.Write([]byte(`{"timestamp":"2024-01-01T00:00:00Z","data":[]}`))
	}))
	defer srv.Close()

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("date", "2024-11-30T12:00:00Z")
	params.Set("regions", "us")
	params.Set("markets", "h2h,spreads")

	resp, err := c.Get(context.Background(), "/v4/historical/sports/basketball_nba/odds", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	if gotParams.Get("date") != "2024-11-30T12:00:00Z" {
		t.Errorf("expected date param, got %s", gotParams.Get("date"))
	}
	if gotParams.Get("regions") != "us" {
		t.Errorf("expected regions=us, got %s", gotParams.Get("regions"))
	}
}
