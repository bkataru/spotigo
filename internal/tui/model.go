// Package tui provides terminal user interface components for Spotigo.
// It includes menu navigation, command execution, and interactive displays
// with a cyberpunk-inspired nuclear green aesthetic.
package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Color scheme - Nuclear Green Cyberpunk
var (
	// Primary colors
	nuclearGreen  = lipgloss.Color("#39FF14") // Bright nuclear green
	glowGreen     = lipgloss.Color("#00FF41") // Matrix green
	darkGreen     = lipgloss.Color("#0D5E1F") // Dark green accent
	deepBlack     = lipgloss.Color("#0A0E0A") // Deep black background
	charcoal      = lipgloss.Color("#1A1F1A") // Charcoal background
	neonCyan      = lipgloss.Color("#00FFFF") // Cyan accent
	terminalGreen = lipgloss.Color("#33FF33") // Terminal green
	dimGreen      = lipgloss.Color("#1F3F1F") // Dim green for borders

	// Text colors
	brightText = lipgloss.Color("#E0FFE0") // Bright text
	dimText    = lipgloss.Color("#7FFF7F") // Dimmed text
	grayText   = lipgloss.Color("#4F5F4F") // Gray text
)

// Choice represents a menu item with display text, command identifier, and icon
type Choice struct {
	Display string
	Command string
	Icon    string
}

// Model represents the TUI state and handles user interactions.
// It manages menu navigation, command execution, and display rendering.
type Model struct {
	choices  []Choice
	cursor   int
	selected bool
	quitting bool
	width    int
	height   int
}

// InitialModel creates a new TUI model with cyberpunk theme.
func InitialModel() Model {
	return Model{
		choices: []Choice{
			{Icon: "▶", Display: "BACKUP LIBRARY", Command: "backup"},
			{Icon: "◉", Display: "AI CHAT TERMINAL", Command: "chat"},
			{Icon: "◈", Display: "SEARCH MUSIC DB", Command: "search"},
			{Icon: "◆", Display: "STATISTICS CORE", Command: "stats"},
			{Icon: "◇", Display: "AUTH STATUS", Command: "auth-status"},
			{Icon: "◐", Display: "MODELS REGISTRY", Command: "models-status"},
			{Icon: "◁", Display: "EXIT SYSTEM", Command: "exit"},
		},
		cursor:   0,
		selected: false,
		quitting: false,
		width:    80,
		height:   24,
	}
}

// Init initializes the TUI model and returns any initial commands.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles user input and state changes in the TUI.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.selected = true
			return m, tea.Quit
		case tea.KeyUp, tea.KeyCtrlK:
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.choices) - 1 // Wrap to bottom
			}
		case tea.KeyDown, tea.KeyCtrlJ:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				m.cursor = 0 // Wrap to top
			}
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlQ:
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the current TUI state as a string for display.
func (m Model) View() string {
	var s strings.Builder

	// ═══════════════════════════════════════════════════════════
	// HEADER: ASCII Art Logo + Title
	// ═══════════════════════════════════════════════════════════
	logo := []string{
		" ███████╗██████╗  ██████╗ ████████╗██╗ ██████╗  ██████╗ ",
		" ██╔════╝██╔══██╗██╔═══██╗╚══██╔══╝██║██╔════╝ ██╔═══██╗",
		" ███████╗██████╔╝██║   ██║   ██║   ██║██║  ███╗██║   ██║",
		" ╚════██║██╔═══╝ ██║   ██║   ██║   ██║██║   ██║██║   ██║",
		" ███████║██║     ╚██████╔╝   ██║   ██║╚██████╔╝╚██████╔╝",
		" ╚══════╝╚═╝      ╚═════╝    ╚═╝   ╚═╝ ╚═════╝  ╚═════╝ ",
	}

	logoStyle := lipgloss.NewStyle().
		Foreground(nuclearGreen).
		Bold(true).
		Align(lipgloss.Center)

	for _, line := range logo {
		s.WriteString(logoStyle.Render(line) + "\n")
	}

	// Subtitle with glowing effect
	subtitleStyle := lipgloss.NewStyle().
		Foreground(neonCyan).
		Bold(true).
		Italic(true).
		Align(lipgloss.Center)

	s.WriteString(subtitleStyle.Render("[ SPOTIFY INTELLIGENCE TERMINAL v2.0 ]") + "\n")

	// Separator line
	separatorStyle := lipgloss.NewStyle().
		Foreground(dimGreen).
		Align(lipgloss.Center)

	separator := strings.Repeat("═", 60)
	s.WriteString(separatorStyle.Render(separator) + "\n\n")

	// ═══════════════════════════════════════════════════════════
	// MENU: Cyberpunk styled menu items
	// ═══════════════════════════════════════════════════════════

	menuTitle := lipgloss.NewStyle().
		Foreground(glowGreen).
		Bold(true).
		Render("╔═══ MAIN MENU ═══╗")

	s.WriteString("  " + menuTitle + "\n\n")

	for i, choice := range m.choices {
		// Cursor indicator
		cursor := "  "
		if m.cursor == i {
			cursor = "▶ "
		}

		cursorStyle := lipgloss.NewStyle().
			Foreground(nuclearGreen).
			Bold(true)

		// Menu item box style
		var itemStyle lipgloss.Style
		if m.cursor == i {
			// Selected item - nuclear green glow
			itemStyle = lipgloss.NewStyle().
				Foreground(deepBlack).
				Background(nuclearGreen).
				Bold(true).
				Padding(0, 2).
				Width(50)
		} else {
			// Unselected item - subtle green
			itemStyle = lipgloss.NewStyle().
				Foreground(brightText).
				Background(charcoal).
				Padding(0, 2).
				Width(50).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(dimGreen)
		}

		// Build the menu line
		itemText := fmt.Sprintf("%s  %s", choice.Icon, choice.Display)
		renderedItem := itemStyle.Render(itemText)

		line := fmt.Sprintf("%s%s", cursorStyle.Render(cursor), renderedItem)
		s.WriteString("  " + line + "\n")
	}

	s.WriteString("\n")

	// ═══════════════════════════════════════════════════════════
	// FOOTER: Help text and status bar
	// ═══════════════════════════════════════════════════════════

	footerSep := lipgloss.NewStyle().
		Foreground(dimGreen).
		Render(strings.Repeat("─", 60))

	s.WriteString("  " + footerSep + "\n")

	// Help text with cyberpunk styling
	helpStyle := lipgloss.NewStyle().
		Foreground(dimText).
		Italic(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(glowGreen).
		Bold(true)

	help := fmt.Sprintf(
		" %s/%s navigate  •  %s select  •  %s/%s quit",
		keyStyle.Render("↑"),
		keyStyle.Render("↓"),
		keyStyle.Render("ENTER"),
		keyStyle.Render("ESC"),
		keyStyle.Render("^C"),
	)

	s.WriteString("  " + helpStyle.Render(help) + "\n")

	// Status bar with system info
	statusStyle := lipgloss.NewStyle().
		Foreground(terminalGreen).
		Background(deepBlack).
		Padding(0, 1).
		Italic(true)

	status := fmt.Sprintf(
		"[ SYSTEM READY ] • [ CURSOR: %d/%d ] • [ MODE: INTERACTIVE ]",
		m.cursor+1,
		len(m.choices),
	)

	s.WriteString("\n  " + statusStyle.Render(status) + "\n")

	// Matrix-style border effect
	borderStyle := lipgloss.NewStyle().
		Foreground(darkGreen)

	border := strings.Repeat("▓", 60)
	s.WriteString("  " + borderStyle.Render(border) + "\n")

	return s.String()
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

// GetCursor returns the current cursor position
func (m Model) GetCursor() int {
	return m.cursor
}

// SetCursor sets the cursor position (for testing)
func (m *Model) SetCursor(pos int) {
	if pos >= 0 && pos < len(m.choices) {
		m.cursor = pos
	}
}

// GetChoicesCount returns the number of menu choices
func (m Model) GetChoicesCount() int {
	return len(m.choices)
}

// RenderCompact returns a compact view without ASCII art (for small terminals)
func (m Model) RenderCompact() string {
	var s strings.Builder

	// Simple header
	headerStyle := lipgloss.NewStyle().
		Foreground(nuclearGreen).
		Bold(true).
		Background(deepBlack).
		Padding(0, 2)

	s.WriteString(headerStyle.Render("SPOTIGO v2.0 TERMINAL") + "\n\n")

	// Menu items
	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "▶ "
		}

		var style lipgloss.Style
		if m.cursor == i {
			style = lipgloss.NewStyle().
				Foreground(nuclearGreen).
				Bold(true)
		} else {
			style = lipgloss.NewStyle().
				Foreground(brightText)
		}

		line := fmt.Sprintf("%s %s %s", cursor, choice.Icon, choice.Display)
		s.WriteString(style.Render(line) + "\n")
	}

	// Simple help
	s.WriteString("\n" + lipgloss.NewStyle().
		Foreground(dimText).
		Render("↑/↓: Navigate • Enter: Select • Esc: Quit") + "\n")

	return s.String()
}
