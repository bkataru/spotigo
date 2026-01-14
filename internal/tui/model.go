package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	choices []string
	cursor  int
	ready   bool
}

func InitialModel() Model {
	return Model{
		choices: []string{
			"🎵 Backup Library",
			"💬 AI Chat",
			"🔍 Search Music",
			"📊 Statistics",
			"🔑 Auth Status",
			"🤖 Models Status",
			"❌ Exit",
		},
		cursor: 0,
		ready:  false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 2).
		Bold(true).
		Render("Spotigo 2.0 - TUI Mode")

	choices := make([]string, len(m.choices))
	for i, choice := range m.choices {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(0, 1)

		if m.cursor == i {
			style = style.Background(lipgloss.Color("#7D56F4")).Bold(true)
		} else {
			style = style.Background(lipgloss.Color("#1A1A1A"))
		}

		choices[i] = style.Render(choice)
	}

	choicesText := lipgloss.JoinVertical(lipgloss.Left, choices...)

	return fmt.Sprintf(
		"\n\n%s\n\n%s\n\n%s",
		title,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).Render("Use ↑/↓ to navigate, Enter to select, Ctrl+C to quit"),
		choicesText,
	)
}

func GetSelectedChoice(m Model) string {
	if m.cursor >= 0 && m.cursor < len(m.choices) {
		return m.choices[m.cursor]
	}
	return ""
}
