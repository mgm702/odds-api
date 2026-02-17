package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestSportsCommand_Table(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/sports.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	resp, err := c.Get(context.Background(), "/v4/sports", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sports, err := client.Decode[[]model.Sport](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(sports) != 2 {
		t.Errorf("expected 2 sports, got %d", len(sports))
	}
	if sports[0].Key != "americanfootball_nfl" {
		t.Errorf("expected nfl, got %s", sports[0].Key)
	}
}

func TestSportsCommand_JSON(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/sports.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	resp, err := c.Get(context.Background(), "/v4/sports", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sports, err := client.Decode[[]model.Sport](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	var buf bytes.Buffer
	data, _ := json.MarshalIndent(sports, "", "  ")
	buf.Write(data)

	var result []model.Sport
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
}

func TestSportsCommand_AllFlag(t *testing.T) {
	var gotAll string
	srv, _ := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAll = r.URL.Query().Get("all")
		withQuotaHeaders(w)
		w.Write([]byte(`[]`))
	})

	c := client.New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("all", "true")
	resp, err := c.Get(context.Background(), "/v4/sports", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	if gotAll != "true" {
		t.Errorf("expected all=true, got %s", gotAll)
	}
}

func TestSportsCommand_MissingAPIKey(t *testing.T) {
	apiKey = ""
	os.Unsetenv("ODDS_API_KEY")
	_, err := getAPIKey()
	if err == nil {
		t.Error("expected error for missing API key")
	}
	if !strings.Contains(err.Error(), "API key required") {
		t.Errorf("expected API key error message, got: %s", err.Error())
	}
}
