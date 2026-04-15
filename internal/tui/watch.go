package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mgm702/odds-api-cli/internal/model"
)

type FetchFunc func(ctx context.Context) (WatchData, error)

type WatchData struct {
	OddsEvents  []model.OddsEvent
	ScoreEvents []model.ScoreEvent
	Quota       model.QuotaInfo
}

type WatchModel struct {
	mode         string
	interval     time.Duration
	fetch        FetchFunc
	table        table.Model
	prevPrices   map[string]float64
	prevScores   map[string]string
	lastUpdate   time.Time
	nextRefresh  time.Time
	quota        model.QuotaInfo
	creditsUsed  int
	paused       bool
	err          error
	width        int
	initialized  bool
	oddsFormat   string
}

type tickMsg time.Time
type fetchResultMsg struct {
	data WatchData
	err  error
}

func NewWatchModel(mode string, interval time.Duration, fetch FetchFunc, oddsFormat string) WatchModel {
	columns := watchColumns(mode)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(nil),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))
	t.SetStyles(s)

	return WatchModel{
		mode:       mode,
		interval:   interval,
		fetch:      fetch,
		table:      t,
		prevPrices: make(map[string]float64),
		prevScores: make(map[string]string),
		oddsFormat: oddsFormat,
	}
}

func watchColumns(mode string) []table.Column {
	if mode == "scores" {
		return []table.Column{
			{Title: "Home", Width: 22},
			{Title: "Away", Width: 22},
			{Title: "Score", Width: 12},
			{Title: "Status", Width: 10},
		}
	}
	return []table.Column{
		{Title: "Matchup", Width: 35},
		{Title: "Bookmaker", Width: 15},
		{Title: "Market", Width: 10},
		{Title: "Outcome", Width: 20},
		{Title: "Price", Width: 10},
	}
}

func (m WatchModel) Init() tea.Cmd {
	return tea.Batch(m.doFetch(), m.tick())
}

func (m WatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "p":
			m.paused = !m.paused
			return m, nil
		case "r":
			return m, m.doFetch()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height - 6)
		m.initialized = true
	case tickMsg:
		if m.paused {
			return m, m.tick()
		}
		m.nextRefresh = time.Now().Add(m.interval)
		return m, tea.Batch(m.doFetch(), m.tick())
	case fetchResultMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.err = nil
		m.lastUpdate = time.Now()
		m.quota = msg.data.Quota
		m.creditsUsed += msg.data.Quota.LastCost
		m.updateTable(msg.data)
		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *WatchModel) updateTable(data WatchData) {
	if m.mode == "scores" {
		m.table.SetRows(m.buildScoreRows(data.ScoreEvents))
	} else {
		m.table.SetRows(m.buildOddsRows(data.OddsEvents))
	}
}

func (m *WatchModel) buildOddsRows(events []model.OddsEvent) []table.Row {
	var rows []table.Row
	newPrices := make(map[string]float64)

	for _, e := range events {
		matchup := fmt.Sprintf("%s vs %s", e.HomeTeam, e.AwayTeam)
		for _, b := range e.Bookmakers {
			for _, mkt := range b.Markets {
				for _, o := range mkt.Outcomes {
					key := fmt.Sprintf("%s|%s|%s|%s", e.ID, b.Key, mkt.Key, o.Name)
					priceStr := formatOddsPrice(o.Price, m.oddsFormat)
					if prev, ok := m.prevPrices[key]; ok {
						if o.Price > prev {
							priceStr = GreenStyle.Render(formatOddsPrice(o.Price, m.oddsFormat) + " +")
						} else if o.Price < prev {
							priceStr = RedStyle.Render(formatOddsPrice(o.Price, m.oddsFormat) + " -")
						}
					}
					newPrices[key] = o.Price
					rows = append(rows, table.Row{matchup, b.Title, mkt.Key, o.Name, priceStr})
					matchup = ""
				}
			}
		}
	}
	m.prevPrices = newPrices
	return rows
}

func (m *WatchModel) buildScoreRows(events []model.ScoreEvent) []table.Row {
	var rows []table.Row
	newScores := make(map[string]string)

	for _, e := range events {
		score := "-"
		status := "upcoming"
		if e.Completed {
			status = "completed"
		} else if e.Scores != nil {
			status = "live"
		}
		if len(e.Scores) >= 2 {
			score = fmt.Sprintf("%s - %s", e.Scores[0].Score, e.Scores[1].Score)
		}

		key := e.ID
		if prev, ok := m.prevScores[key]; ok && prev != score && score != "-" {
			score = GreenStyle.Render(score)
		}
		newScores[key] = score
		rows = append(rows, table.Row{e.HomeTeam, e.AwayTeam, score, status})
	}
	m.prevScores = newScores
	return rows
}

func (m WatchModel) View() string {
	var b strings.Builder

	title := "Watch: Odds"
	if m.mode == "scores" {
		title = "Watch: Scores"
	}
	b.WriteString(TitleStyle.Render(title))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(ErrorStyle.Render(fmt.Sprintf("Error: %s", m.err)))
		b.WriteString("\n\n")
	}

	b.WriteString(m.table.View())
	b.WriteString("\n\n")

	var status []string
	if !m.lastUpdate.IsZero() {
		status = append(status, fmt.Sprintf("Updated: %s", m.lastUpdate.Format("15:04:05")))
	}
	if !m.nextRefresh.IsZero() && !m.paused {
		remaining := time.Until(m.nextRefresh).Round(time.Second)
		if remaining < 0 {
			remaining = 0
		}
		status = append(status, fmt.Sprintf("Next: %s", remaining))
	}
	status = append(status, fmt.Sprintf("Credits: %d remaining", m.quota.Remaining))
	status = append(status, fmt.Sprintf("Session cost: %d", m.creditsUsed))

	if m.paused {
		status = append(status, RedStyle.Render("PAUSED"))
	}

	if m.quota.Remaining > 0 && m.quota.Remaining < 50 {
		status = append(status, ErrorStyle.Render("LOW CREDITS"))
	}

	b.WriteString(StatusBarStyle.Render(strings.Join(status, " | ")))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("q: quit | p: pause/resume | r: refresh | arrow keys: navigate"))

	return b.String()
}

func (m WatchModel) tick() tea.Cmd {
	return tea.Tick(m.interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m WatchModel) doFetch() tea.Cmd {
	return func() tea.Msg {
		data, err := m.fetch(context.Background())
		return fetchResultMsg{data: data, err: err}
	}
}
