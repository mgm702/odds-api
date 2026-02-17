package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mgm702/odds-api-cli/internal/model"
)

type sportItem struct {
	sport model.Sport
}

func (s sportItem) Title() string       { return s.sport.Title }
func (s sportItem) Description() string { return fmt.Sprintf("%s | %s", s.sport.Key, s.sport.Group) }
func (s sportItem) FilterValue() string { return s.sport.Title + " " + s.sport.Key }

type SportsModel struct {
	list     list.Model
	Selected *model.Sport
	quitting bool
}

func NewSportsModel(sports []model.Sport) SportsModel {
	items := make([]list.Item, len(sports))
	for i, s := range sports {
		items[i] = sportItem{sport: s}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("170")).
		BorderForeground(lipgloss.Color("170"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("243")).
		BorderForeground(lipgloss.Color("170"))

	l := list.New(items, delegate, 60, 20)
	l.Title = "Sports"
	l.Styles.Title = TitleStyle
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(true)

	return SportsModel{list: l}
}

func (m SportsModel) Init() tea.Cmd {
	return nil
}

func (m SportsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(sportItem); ok {
				s := item.sport
				m.Selected = &s
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SportsModel) View() string {
	return m.list.View()
}

func (m SportsModel) WriteSelected(w io.Writer) {
	if m.Selected != nil {
		fmt.Fprintln(w, m.Selected.Key)
	}
}
