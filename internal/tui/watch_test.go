package tui

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func mockFetch(data WatchData, err error) FetchFunc {
	return func(ctx context.Context) (WatchData, error) {
		return data, err
	}
}

func testWatchData() WatchData {
	return WatchData{
		OddsEvents: []model.OddsEvent{
			{
				ID: "evt001", HomeTeam: "LAL", AwayTeam: "BOS",
				Bookmakers: []model.Bookmaker{
					{
						Key: "fanduel", Title: "FanDuel",
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
		},
		Quota: model.QuotaInfo{Remaining: 490, Used: 10, LastCost: 1},
	}
}

func testScoreData() WatchData {
	return WatchData{
		ScoreEvents: []model.ScoreEvent{
			{
				ID: "evt001", HomeTeam: "LAL", AwayTeam: "BOS",
				Scores: []model.Score{{Name: "LAL", Score: "55"}, {Name: "BOS", Score: "52"}},
			},
		},
		Quota: model.QuotaInfo{Remaining: 498, Used: 2, LastCost: 1},
	}
}

func TestWatchModel_Init(t *testing.T) {
	m := NewWatchModel("odds", 60*time.Second, mockFetch(testWatchData(), nil), "")
	cmd := m.Init()
	if cmd == nil {
		t.Error("expected non-nil cmd from Init (batch of fetch + tick)")
	}
}

func TestWatchModel_ViewOdds(t *testing.T) {
	m := NewWatchModel("odds", 60*time.Second, mockFetch(testWatchData(), nil), "")
	view := m.View()
	if !strings.Contains(view, "Watch: Odds") {
		t.Error("expected odds title")
	}
}

func TestWatchModel_ViewScores(t *testing.T) {
	m := NewWatchModel("scores", 60*time.Second, mockFetch(testScoreData(), nil), "")
	view := m.View()
	if !strings.Contains(view, "Watch: Scores") {
		t.Error("expected scores title")
	}
}

func TestWatchModel_Quit(t *testing.T) {
	m := NewWatchModel("odds", 60*time.Second, mockFetch(testWatchData(), nil), "")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Error("expected quit cmd")
	}
}

func TestWatchModel_Pause(t *testing.T) {
	m := NewWatchModel("odds", 60*time.Second, mockFetch(testWatchData(), nil), "")

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model := updated.(WatchModel)
	if !model.paused {
		t.Error("expected paused=true")
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model = updated.(WatchModel)
	if model.paused {
		t.Error("expected paused=false")
	}
}

func TestWatchModel_FetchResult(t *testing.T) {
	m := NewWatchModel("odds", 60*time.Second, mockFetch(testWatchData(), nil), "")

	data := testWatchData()
	updated, _ := m.Update(fetchResultMsg{data: data})
	model := updated.(WatchModel)

	if model.quota.Remaining != 490 {
		t.Errorf("expected remaining=490, got %d", model.quota.Remaining)
	}
	if model.lastUpdate.IsZero() {
		t.Error("expected lastUpdate to be set")
	}
}

func TestWatchModel_FetchError(t *testing.T) {
	m := NewWatchModel("odds", 60*time.Second, mockFetch(testWatchData(), nil), "")

	updated, _ := m.Update(fetchResultMsg{err: fmt.Errorf("network error")})
	model := updated.(WatchModel)

	if model.err == nil {
		t.Error("expected error to be set")
	}
	view := model.View()
	if !strings.Contains(view, "network error") {
		t.Error("expected error in view")
	}
}

func TestWatchModel_PriceHighlight(t *testing.T) {
	m := NewWatchModel("odds", 60*time.Second, mockFetch(testWatchData(), nil), "")

	data := testWatchData()
	updated, _ := m.Update(fetchResultMsg{data: data})
	model := updated.(WatchModel)

	data.OddsEvents[0].Bookmakers[0].Markets[0].Outcomes[0].Price = 2.00
	updated, _ = model.Update(fetchResultMsg{data: data})
	model = updated.(WatchModel)

	if model.prevPrices["evt001|fanduel|h2h|LAL"] != 2.00 {
		t.Error("expected updated price tracking")
	}
}

func TestWatchModel_WindowSize(t *testing.T) {
	m := NewWatchModel("odds", 60*time.Second, mockFetch(testWatchData(), nil), "")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model := updated.(WatchModel)
	if !model.initialized {
		t.Error("expected initialized=true after window size")
	}
}

func TestWatchModel_LowCreditsWarning(t *testing.T) {
	m := NewWatchModel("odds", 60*time.Second, mockFetch(testWatchData(), nil), "")

	data := testWatchData()
	data.Quota.Remaining = 30
	updated, _ := m.Update(fetchResultMsg{data: data})
	model := updated.(WatchModel)

	view := model.View()
	if !strings.Contains(view, "LOW CREDITS") {
		t.Error("expected low credits warning")
	}
}
