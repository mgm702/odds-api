package tui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("75")).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder())

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	GreenStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	RedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)
