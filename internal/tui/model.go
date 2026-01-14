package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Choice represents a menu item with display text and command identifier
type Choice struct {
	Display string
	Command string
}

type Model struct {
	choices  []Choice
	cursor   int
	selected bool
	quitting bool
}

func InitialModel() Model {
	return Model{
		choices: []Choice{
			{Display: "🎵 Backup Library", Command: "backup"},
			{Display: "💬 AI Chat", Command: "chat"},
			{Display: "🔍 Search Music", Command: "search"},
			{Display: "📊 Statistics", Command: "stats"},
			{Display: "🔑 Auth Status", Command: "auth-status"},
			{Display: "🤖 Models Status", Command: "models-status"},
			{Display: "❌ Exit", Command: "exit"},
		},
		cursor:   0,
		selected: false,
		quitting: false,
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
			m.selected = true
			return m, tea.Quit
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
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

		choices[i] = style.Render(choice.Display)
	}

	choicesText := lipgloss.JoinVertical(lipgloss.Left, choices...)

	return fmt.Sprintf(
		"\n\n%s\n\n%s\n\n%s",
		title,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).Render("Use ↑/↓ to navigate, Enter to select, Ctrl+C/Esc to quit"),
		choicesText,
	)
}

// GetSelectedCommand returns the command identifier of the selected choice
func (m Model) GetSelectedCommand() string {
	if !m.selected || m.quitting {
		return ""
	}
	if m.cursor >= 0 && m.cursor < len(m.choices) {
		return m.choices[m.cursor].Command
	}
	return ""
}

// GetSelectedDisplay returns the display text of the selected choice
func (m Model) GetSelectedDisplay() string {
	if m.cursor >= 0 && m.cursor < len(m.choices) {
		return m.choices[m.cursor].Display
	}
	return ""
}

// IsQuitting returns true if the user quit without selecting
func (m Model) IsQuitting() bool {
	return m.quitting
}

// WasSelected returns true if a selection was made
func (m Model) WasSelected() bool {
	return m.selected
}
