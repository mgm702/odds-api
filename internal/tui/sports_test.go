package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestSportsModel_Init(t *testing.T) {
	sports := []model.Sport{
		{Key: "nfl", Title: "NFL", Group: "American Football", Active: true},
		{Key: "nba", Title: "NBA", Group: "Basketball", Active: true},
	}
	m := NewSportsModel(sports)
	cmd := m.Init()
	if cmd != nil {
		t.Error("expected nil cmd from Init")
	}
}

func TestSportsModel_View(t *testing.T) {
	sports := []model.Sport{
		{Key: "nfl", Title: "NFL", Group: "American Football", Active: true},
	}
	m := NewSportsModel(sports)
	view := m.View()
	if !strings.Contains(view, "Sports") {
		t.Error("expected Sports title in view")
	}
}

func TestSportsModel_Quit(t *testing.T) {
	sports := []model.Sport{
		{Key: "nfl", Title: "NFL", Group: "American Football", Active: true},
	}
	m := NewSportsModel(sports)
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	model := updated.(SportsModel)
	if !model.quitting {
		t.Error("expected quitting=true")
	}
	if cmd == nil {
		t.Error("expected quit cmd")
	}
}

func TestSportsModel_Select(t *testing.T) {
	sports := []model.Sport{
		{Key: "nfl", Title: "NFL", Group: "American Football", Active: true},
	}
	m := NewSportsModel(sports)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model := updated.(SportsModel)
	if model.Selected == nil {
		t.Error("expected selected sport")
	}
	if model.Selected.Key != "nfl" {
		t.Errorf("expected nfl, got %s", model.Selected.Key)
	}
}

func TestSportsModel_WindowSize(t *testing.T) {
	sports := []model.Sport{
		{Key: "nfl", Title: "NFL", Group: "American Football", Active: true},
	}
	m := NewSportsModel(sports)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	_ = updated.(SportsModel)
}

func TestSportItem_FilterValue(t *testing.T) {
	item := sportItem{sport: model.Sport{Key: "nfl", Title: "NFL"}}
	fv := item.FilterValue()
	if !strings.Contains(fv, "NFL") || !strings.Contains(fv, "nfl") {
		t.Errorf("expected filter value to contain key and title, got %s", fv)
	}
}
