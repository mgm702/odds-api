package tui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mgm702/odds-api-cli/internal/model"
)

type OddsModel struct {
	tables     []eventTable
	current    int
	quitting   bool
	oddsFormat string
}

func formatOddsPrice(price float64, oddsFormat string) string {
	if strings.ToLower(oddsFormat) == "american" {
		n := int(math.Round(price))
		if n >= 0 {
			return fmt.Sprintf("+%d", n)
		}
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%.2f", price)
}

type eventTable struct {
	title string
	table table.Model
}

func NewOddsModel(events []model.OddsEvent, oddsFormat string) OddsModel {
	var tables []eventTable
	for _, e := range events {
		rows := buildOddsRows(e, oddsFormat)
		columns := []table.Column{
			{Title: "Bookmaker", Width: 20},
			{Title: "Market", Width: 12},
			{Title: "Outcome", Width: 25},
			{Title: "Price", Width: 10},
			{Title: "Point", Width: 10},
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(min(len(rows)+1, 20)),
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

		title := fmt.Sprintf("%s vs %s (%s)", e.HomeTeam, e.AwayTeam, e.CommenceTime)
		tables = append(tables, eventTable{title: title, table: t})
	}

	return OddsModel{tables: tables, oddsFormat: oddsFormat}
}

func buildOddsRows(e model.OddsEvent, oddsFormat string) []table.Row {
	var rows []table.Row
	for _, b := range e.Bookmakers {
		for _, m := range b.Markets {
			for _, o := range m.Outcomes {
				point := ""
				if o.Point != nil {
					point = fmt.Sprintf("%+.1f", *o.Point)
				}
				rows = append(rows, table.Row{
					b.Title, m.Key, o.Name, formatOddsPrice(o.Price, oddsFormat), point,
				})
			}
		}
	}
	return rows
}

func (m OddsModel) Init() tea.Cmd {
	return nil
}

func (m OddsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "tab", "right":
			if m.current < len(m.tables)-1 {
				m.current++
			}
			return m, nil
		case "shift+tab", "left":
			if m.current > 0 {
				m.current--
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		for i := range m.tables {
			m.tables[i].table.SetWidth(msg.Width)
		}
	}

	if len(m.tables) > 0 {
		var cmd tea.Cmd
		m.tables[m.current].table, cmd = m.tables[m.current].table.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m OddsModel) View() string {
	if len(m.tables) == 0 {
		return DimStyle.Render("No odds data available.")
	}

	var b strings.Builder
	et := m.tables[m.current]

	b.WriteString(TitleStyle.Render(et.title))
	b.WriteString("\n")

	if len(m.tables) > 1 {
		nav := fmt.Sprintf("Event %d/%d (Tab/Shift+Tab to navigate)", m.current+1, len(m.tables))
		b.WriteString(DimStyle.Render(nav))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(et.table.View())
	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("q: quit | tab/shift+tab: switch events | arrow keys: navigate"))

	return b.String()
}
