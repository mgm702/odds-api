package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestWriteSports(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, false)
	tw.WriteSports([]model.Sport{
		{Key: "nfl", Title: "NFL", Group: "American Football", Active: true},
		{Key: "nba", Title: "NBA", Group: "Basketball", Active: false},
	})
	out := buf.String()
	if !strings.Contains(out, "nfl") {
		t.Error("expected nfl in output")
	}
	if !strings.Contains(out, "yes") {
		t.Error("expected yes for active sport")
	}
	if !strings.Contains(out, "no") {
		t.Error("expected no for inactive sport")
	}
	if !strings.Contains(out, "KEY") {
		t.Error("expected header row")
	}
}

func TestWriteEvents(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, false)
	tw.WriteEvents([]model.Event{
		{ID: "e1", HomeTeam: "LAL", AwayTeam: "BOS", CommenceTime: "2024-01-01T00:00:00Z"},
	})
	out := buf.String()
	if !strings.Contains(out, "LAL") || !strings.Contains(out, "BOS") {
		t.Error("expected team names in output")
	}
}

func TestWriteOdds(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, false)
	tw.WriteOdds([]model.OddsEvent{
		{
			HomeTeam: "LAL", AwayTeam: "BOS", CommenceTime: "2024-01-01T00:00:00Z",
			Bookmakers: []model.Bookmaker{
				{
					Title: "FanDuel",
					Markets: []model.Market{
						{
							Key: "h2h",
							Outcomes: []model.Outcome{
								{Name: "LAL", Price: 1.91},
								{Name: "BOS", Price: 1.95},
							},
						},
					},
				},
			},
		},
	})
	out := buf.String()
	if !strings.Contains(out, "LAL vs BOS") {
		t.Error("expected matchup header")
	}
	if !strings.Contains(out, "FanDuel") {
		t.Error("expected bookmaker name")
	}
	if !strings.Contains(out, "1.91") {
		t.Error("expected odds value")
	}
}

func TestWriteOddsWithSpreads(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, false)
	point := -3.5
	tw.WriteOdds([]model.OddsEvent{
		{
			HomeTeam: "LAL", AwayTeam: "BOS", CommenceTime: "2024-01-01T00:00:00Z",
			Bookmakers: []model.Bookmaker{
				{
					Title: "DK",
					Markets: []model.Market{
						{
							Key: "spreads",
							Outcomes: []model.Outcome{
								{Name: "LAL", Price: 1.91, Point: &point},
							},
						},
					},
				},
			},
		},
	})
	out := buf.String()
	if !strings.Contains(out, "-3.5") {
		t.Error("expected spread point in output")
	}
}

func TestWriteScores(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, false)
	updated := "2024-01-01T01:00:00Z"
	tw.WriteScores([]model.ScoreEvent{
		{
			HomeTeam: "LAL", AwayTeam: "BOS", Completed: false,
			Scores:     []model.Score{{Name: "LAL", Score: "55"}, {Name: "BOS", Score: "52"}},
			LastUpdate: &updated,
		},
		{
			HomeTeam: "CHI", AwayTeam: "NYK", Completed: false,
			Scores: nil, LastUpdate: nil,
		},
	})
	out := buf.String()
	if !strings.Contains(out, "55 - 52") {
		t.Error("expected score in output")
	}
	if !strings.Contains(out, "live") {
		t.Error("expected live status")
	}
	if !strings.Contains(out, "upcoming") {
		t.Error("expected upcoming status")
	}
}

func TestWriteCredits(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, false)
	tw.WriteCredits(model.QuotaInfo{Used: 142, Remaining: 358, LastCost: 0})
	out := buf.String()
	if !strings.Contains(out, "Credits Report") {
		t.Error("expected header")
	}
	if !strings.Contains(out, "142") {
		t.Error("expected used count")
	}
	if !strings.Contains(out, "358") {
		t.Error("expected remaining count")
	}
	if !strings.Contains(out, "28.4%") {
		t.Error("expected percentage")
	}
}

func TestWriteParticipants(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, false)
	tw.WriteParticipants([]model.Participant{
		{ID: "p1", Name: "Lakers", IsActive: true},
	})
	out := buf.String()
	if !strings.Contains(out, "Lakers") {
		t.Error("expected participant name")
	}
}

func TestWriteEventMarkets(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, false)
	tw.WriteEventMarkets(model.EventMarkets{
		ID: "e1",
		Bookmakers: []model.BookmakerMarkets{
			{Key: "fd", Title: "FanDuel", Markets: []model.MarketInfo{{Key: "h2h"}, {Key: "spreads"}}},
		},
	})
	out := buf.String()
	if !strings.Contains(out, "FanDuel") {
		t.Error("expected bookmaker")
	}
	if !strings.Contains(out, "h2h") {
		t.Error("expected market key")
	}
}
