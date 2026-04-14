package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func testOddsEvents() []model.OddsEvent {
	return []model.OddsEvent{
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
		{
			HomeTeam: "GSW", AwayTeam: "MIA", CommenceTime: "2024-01-01T02:00:00Z",
			Bookmakers: []model.Bookmaker{
				{
					Title: "DraftKings",
					Markets: []model.Market{
						{
							Key: "h2h",
							Outcomes: []model.Outcome{
								{Name: "GSW", Price: 1.80},
								{Name: "MIA", Price: 2.05},
							},
						},
					},
				},
			},
		},
	}
}

func TestOddsModel_Init(t *testing.T) {
	m := NewOddsModel(testOddsEvents(), "")
	cmd := m.Init()
	if cmd != nil {
		t.Error("expected nil cmd from Init")
	}
}

func TestOddsModel_View(t *testing.T) {
	m := NewOddsModel(testOddsEvents(), "")
	view := m.View()
	if !strings.Contains(view, "LAL vs BOS") {
		t.Error("expected matchup in view")
	}
	if !strings.Contains(view, "1/2") {
		t.Error("expected event navigation indicator")
	}
}

func TestOddsModel_Empty(t *testing.T) {
	m := NewOddsModel(nil, "")
	view := m.View()
	if !strings.Contains(view, "No odds data") {
		t.Error("expected empty state message")
	}
}

func TestOddsModel_Quit(t *testing.T) {
	m := NewOddsModel(testOddsEvents(), "")
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	model := updated.(OddsModel)
	if !model.quitting {
		t.Error("expected quitting=true")
	}
	if cmd == nil {
		t.Error("expected quit cmd")
	}
}

func TestOddsModel_Navigate(t *testing.T) {
	m := NewOddsModel(testOddsEvents(), "")

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	model := updated.(OddsModel)
	if model.current != 1 {
		t.Errorf("expected current=1 after tab, got %d", model.current)
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = updated.(OddsModel)
	if model.current != 0 {
		t.Errorf("expected current=0 after shift+tab, got %d", model.current)
	}
}

func TestOddsModel_NavigateBounds(t *testing.T) {
	m := NewOddsModel(testOddsEvents(), "")

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model := updated.(OddsModel)
	if model.current != 0 {
		t.Errorf("should not go below 0, got %d", model.current)
	}

	model.current = 1
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updated.(OddsModel)
	if model.current != 1 {
		t.Errorf("should not exceed max, got %d", model.current)
	}
}

func TestBuildOddsRows(t *testing.T) {
	point := -3.5
	e := model.OddsEvent{
		Bookmakers: []model.Bookmaker{
			{
				Title: "FD",
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
	}
	rows := buildOddsRows(e, "")
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0][4] != "-3.5" {
		t.Errorf("expected -3.5, got %s", rows[0][4])
	}
}
