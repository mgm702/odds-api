package model

import (
	"encoding/json"
	"testing"
)

func TestSportUnmarshal(t *testing.T) {
	data := `{"key":"nfl","group":"American Football","title":"NFL","description":"US Football","active":true,"has_outrights":false}`
	var s Sport
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if s.Key != "nfl" || s.Title != "NFL" || !s.Active {
		t.Errorf("unexpected sport: %+v", s)
	}
}

func TestEventUnmarshal(t *testing.T) {
	data := `{"id":"e1","sport_key":"nba","sport_title":"NBA","commence_time":"2024-01-01T00:00:00Z","home_team":"LAL","away_team":"BOS"}`
	var e Event
	if err := json.Unmarshal([]byte(data), &e); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if e.ID != "e1" || e.HomeTeam != "LAL" {
		t.Errorf("unexpected event: %+v", e)
	}
}

func TestOutcomeWithPoint(t *testing.T) {
	data := `{"name":"Team A","price":1.91,"point":-3.5}`
	var o Outcome
	if err := json.Unmarshal([]byte(data), &o); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if o.Point == nil || *o.Point != -3.5 {
		t.Errorf("expected point=-3.5, got %v", o.Point)
	}
}

func TestOutcomeWithoutPoint(t *testing.T) {
	data := `{"name":"Team A","price":1.91}`
	var o Outcome
	if err := json.Unmarshal([]byte(data), &o); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if o.Point != nil {
		t.Errorf("expected nil point, got %v", o.Point)
	}
}

func TestScoreEventNullScores(t *testing.T) {
	data := `{"id":"e1","sport_key":"nba","sport_title":"NBA","commence_time":"2024-01-01T00:00:00Z","home_team":"LAL","away_team":"BOS","completed":false,"scores":null,"last_update":null}`
	var se ScoreEvent
	if err := json.Unmarshal([]byte(data), &se); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if se.Scores != nil {
		t.Errorf("expected nil scores, got %v", se.Scores)
	}
	if se.LastUpdate != nil {
		t.Errorf("expected nil last_update, got %v", se.LastUpdate)
	}
}

func TestScoreEventWithScores(t *testing.T) {
	data := `{"id":"e1","sport_key":"nba","sport_title":"NBA","commence_time":"2024-01-01T00:00:00Z","home_team":"LAL","away_team":"BOS","completed":true,"scores":[{"name":"LAL","score":"110"},{"name":"BOS","score":"105"}],"last_update":"2024-01-01T02:00:00Z"}`
	var se ScoreEvent
	if err := json.Unmarshal([]byte(data), &se); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(se.Scores) != 2 {
		t.Errorf("expected 2 scores, got %d", len(se.Scores))
	}
	if !se.Completed {
		t.Error("expected completed=true")
	}
}

func TestOddsEventEmptyBookmakers(t *testing.T) {
	data := `{"id":"e1","sport_key":"nba","sport_title":"NBA","commence_time":"2024-01-01T00:00:00Z","home_team":"LAL","away_team":"BOS","bookmakers":[]}`
	var oe OddsEvent
	if err := json.Unmarshal([]byte(data), &oe); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(oe.Bookmakers) != 0 {
		t.Errorf("expected 0 bookmakers, got %d", len(oe.Bookmakers))
	}
}

func TestHistoricalResponse(t *testing.T) {
	data := `{"timestamp":"2024-01-01T00:00:00Z","previous_timestamp":"2023-12-31T23:55:00Z","next_timestamp":"2024-01-01T00:05:00Z","data":[{"key":"nfl"}]}`
	var h HistoricalResponse[[]Sport]
	if err := json.Unmarshal([]byte(data), &h); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if h.Timestamp != "2024-01-01T00:00:00Z" {
		t.Errorf("unexpected timestamp: %s", h.Timestamp)
	}
	if h.PreviousTimestamp == nil || *h.PreviousTimestamp != "2023-12-31T23:55:00Z" {
		t.Errorf("unexpected previous_timestamp: %v", h.PreviousTimestamp)
	}
	if len(h.Data) != 1 {
		t.Errorf("expected 1 item in data, got %d", len(h.Data))
	}
}

func TestHistoricalResponseNullTimestamps(t *testing.T) {
	data := `{"timestamp":"2024-01-01T00:00:00Z","previous_timestamp":null,"next_timestamp":null,"data":[]}`
	var h HistoricalResponse[[]Sport]
	if err := json.Unmarshal([]byte(data), &h); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if h.PreviousTimestamp != nil {
		t.Errorf("expected nil previous_timestamp")
	}
	if h.NextTimestamp != nil {
		t.Errorf("expected nil next_timestamp")
	}
}

func TestQuotaInfo(t *testing.T) {
	data := `{"remaining":500,"used":100,"last_cost":2}`
	var q QuotaInfo
	if err := json.Unmarshal([]byte(data), &q); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if q.Remaining != 500 || q.Used != 100 || q.LastCost != 2 {
		t.Errorf("unexpected quota: %+v", q)
	}
}

func TestParticipantUnmarshal(t *testing.T) {
	data := `{"id":"p1","name":"Lakers","is_active":true}`
	var p Participant
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if p.ID != "p1" || !p.IsActive {
		t.Errorf("unexpected participant: %+v", p)
	}
}
