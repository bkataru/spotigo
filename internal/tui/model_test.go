package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitialModel(t *testing.T) {
	model := InitialModel()

	if model.cursor != 0 {
		t.Errorf("Expected initial cursor at 0, got %d", model.cursor)
	}

	if model.selected {
		t.Error("Model should not be selected initially")
	}

	if model.quitting {
		t.Error("Model should not be quitting initially")
	}

	if len(model.choices) == 0 {
		t.Error("Model should have choices")
	}

	expectedChoices := 7 // backup, chat, search, stats, auth, models, exit
	if len(model.choices) != expectedChoices {
		t.Errorf("Expected %d choices, got %d", expectedChoices, len(model.choices))
	}
}

func TestChoiceStructure(t *testing.T) {
	model := InitialModel()

	for i, choice := range model.choices {
		if choice.Display == "" {
			t.Errorf("Choice %d has empty Display", i)
		}
		if choice.Command == "" {
			t.Errorf("Choice %d has empty Command", i)
		}
		if choice.Icon == "" {
			t.Errorf("Choice %d has empty Icon", i)
		}
	}

	// Test specific choices
	expectedCommands := map[string]bool{
		"backup":        true,
		"chat":          true,
		"search":        true,
		"stats":         true,
		"auth-status":   true,
		"models-status": true,
		"exit":          true,
	}

	foundCommands := make(map[string]bool)
	for _, choice := range model.choices {
		foundCommands[choice.Command] = true
	}

	for cmd := range expectedCommands {
		if !foundCommands[cmd] {
			t.Errorf("Missing expected command: %s", cmd)
		}
	}
}

func TestModelInit(t *testing.T) {
	model := InitialModel()
	cmd := model.Init()

	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

func TestModelUpdateKeyDown(t *testing.T) {
	model := InitialModel()
	initialCursor := model.cursor

	// Press down key
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m := updatedModel.(Model)

	if m.cursor != initialCursor+1 {
		t.Errorf("Expected cursor at %d, got %d", initialCursor+1, m.cursor)
	}
}

func TestModelUpdateKeyUp(t *testing.T) {
	model := InitialModel()
	model.cursor = 2 // Start at position 2

	// Press up key
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	m := updatedModel.(Model)

	if m.cursor != 1 {
		t.Errorf("Expected cursor at 1, got %d", m.cursor)
	}
}

func TestModelUpdateKeyUpWrap(t *testing.T) {
	model := InitialModel()
	model.cursor = 0 // Start at top

	// Press up key - should wrap to bottom
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	m := updatedModel.(Model)

	expectedCursor := len(model.choices) - 1
	if m.cursor != expectedCursor {
		t.Errorf("Expected cursor to wrap to %d, got %d", expectedCursor, m.cursor)
	}
}

func TestModelUpdateKeyDownWrap(t *testing.T) {
	model := InitialModel()
	model.cursor = len(model.choices) - 1 // Start at bottom

	// Press down key - should wrap to top
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m := updatedModel.(Model)

	if m.cursor != 0 {
		t.Errorf("Expected cursor to wrap to 0, got %d", m.cursor)
	}
}

func TestModelUpdateEnter(t *testing.T) {
	model := InitialModel()
	model.cursor = 1

	// Press enter
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updatedModel.(Model)

	if !m.selected {
		t.Error("Model should be selected after Enter")
	}

	if cmd == nil {
		t.Error("Enter should return quit command")
	}
}

func TestModelUpdateEscape(t *testing.T) {
	model := InitialModel()

	// Press escape
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m := updatedModel.(Model)

	if !m.quitting {
		t.Error("Model should be quitting after Esc")
	}

	if cmd == nil {
		t.Error("Escape should return quit command")
	}
}

func TestModelUpdateCtrlC(t *testing.T) {
	model := InitialModel()

	// Press Ctrl+C
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	m := updatedModel.(Model)

	if !m.quitting {
		t.Error("Model should be quitting after Ctrl+C")
	}

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
}

func TestModelUpdateCtrlQ(t *testing.T) {
	model := InitialModel()

	// Press Ctrl+Q
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlQ})
	m := updatedModel.(Model)

	if !m.quitting {
		t.Error("Model should be quitting after Ctrl+Q")
	}
}

func TestModelUpdateWindowSize(t *testing.T) {
	model := InitialModel()

	// Send window size message
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m := updatedModel.(Model)

	if m.width != 100 {
		t.Errorf("Expected width 100, got %d", m.width)
	}

	if m.height != 40 {
		t.Errorf("Expected height 40, got %d", m.height)
	}
}

func TestModelView(t *testing.T) {
	model := InitialModel()
	view := model.View()

	if view == "" {
		t.Error("View should not be empty")
	}

	// Check for nuclear green theme elements
	expectedStrings := []string{
		"MAIN MENU",
		"navigate",
		"select",
		"quit",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("View should contain '%s'", expected)
		}
	}

	// Check for ASCII art elements (the logo contains box drawing chars)
	if !strings.Contains(view, "SPOTIFY") && !strings.Contains(view, "═") {
		t.Error("View should contain ASCII art logo")
	}
}

func TestModelViewContainsChoices(t *testing.T) {
	model := InitialModel()
	view := model.View()

	// Each choice display should appear in the view
	for _, choice := range model.choices {
		if !strings.Contains(view, choice.Display) {
			t.Errorf("View should contain choice: %s", choice.Display)
		}
	}
}

func TestModelViewCursorIndicator(t *testing.T) {
	model := InitialModel()
	model.cursor = 2
	view := model.View()

	// Should contain cursor indicator
	if !strings.Contains(view, "▶") {
		t.Error("View should contain cursor indicator '▶'")
	}
}

func TestGetSelectedCommand(t *testing.T) {
	tests := []struct {
		name        string
		cursor      int
		selected    bool
		quitting    bool
		expectEmpty bool
	}{
		{
			name:        "selected at position 0",
			cursor:      0,
			selected:    true,
			quitting:    false,
			expectEmpty: false,
		},
		{
			name:        "not selected",
			cursor:      0,
			selected:    false,
			quitting:    false,
			expectEmpty: true,
		},
		{
			name:        "quitting",
			cursor:      0,
			selected:    false,
			quitting:    true,
			expectEmpty: true,
		},
		{
			name:        "selected at position 1",
			cursor:      1,
			selected:    true,
			quitting:    false,
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := InitialModel()
			model.cursor = tt.cursor
			model.selected = tt.selected
			model.quitting = tt.quitting

			cmd := model.GetSelectedCommand()

			isEmpty := cmd == ""
			if isEmpty != tt.expectEmpty {
				t.Errorf("Expected empty=%v, got empty=%v (cmd='%s')", tt.expectEmpty, isEmpty, cmd)
			}

			if !isEmpty && tt.cursor < len(model.choices) {
				expectedCmd := model.choices[tt.cursor].Command
				if cmd != expectedCmd {
					t.Errorf("Expected command '%s', got '%s'", expectedCmd, cmd)
				}
			}
		})
	}
}

func TestGetSelectedDisplay(t *testing.T) {
	model := InitialModel()
	model.cursor = 2

	display := model.GetSelectedDisplay()
	expected := model.choices[2].Display

	if display != expected {
		t.Errorf("Expected display '%s', got '%s'", expected, display)
	}
}

func TestGetSelectedDisplayInvalid(t *testing.T) {
	model := InitialModel()
	model.cursor = 999 // Invalid position

	display := model.GetSelectedDisplay()

	if display != "" {
		t.Errorf("Expected empty display for invalid cursor, got '%s'", display)
	}
}

func TestIsQuitting(t *testing.T) {
	model := InitialModel()

	if model.IsQuitting() {
		t.Error("Model should not be quitting initially")
	}

	model.quitting = true

	if !model.IsQuitting() {
		t.Error("Model should be quitting after setting flag")
	}
}

func TestWasSelected(t *testing.T) {
	model := InitialModel()

	if model.WasSelected() {
		t.Error("Model should not be selected initially")
	}

	model.selected = true

	if !model.WasSelected() {
		t.Error("Model should be selected after setting flag")
	}
}

func TestGetCursor(t *testing.T) {
	model := InitialModel()
	model.cursor = 3

	cursor := model.GetCursor()

	if cursor != 3 {
		t.Errorf("Expected cursor 3, got %d", cursor)
	}
}

func TestSetCursor(t *testing.T) {
	model := InitialModel()

	model.SetCursor(4)

	if model.cursor != 4 {
		t.Errorf("Expected cursor 4, got %d", model.cursor)
	}
}

func TestSetCursorInvalid(t *testing.T) {
	model := InitialModel()
	originalCursor := model.cursor

	// Try to set invalid cursor
	model.SetCursor(-1)

	if model.cursor != originalCursor {
		t.Error("Cursor should not change for invalid negative position")
	}

	model.SetCursor(999)

	if model.cursor != originalCursor {
		t.Error("Cursor should not change for invalid out-of-bounds position")
	}
}

func TestGetChoicesCount(t *testing.T) {
	model := InitialModel()
	count := model.GetChoicesCount()

	if count != len(model.choices) {
		t.Errorf("Expected count %d, got %d", len(model.choices), count)
	}

	if count != 7 {
		t.Errorf("Expected 7 choices, got %d", count)
	}
}

func TestRenderCompact(t *testing.T) {
	model := InitialModel()
	compact := model.RenderCompact()

	if compact == "" {
		t.Error("Compact view should not be empty")
	}

	// Check for essential elements
	expectedStrings := []string{
		"SPOTIGO",
		"Navigate",
		"Select",
		"Quit",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(compact, expected) {
			t.Errorf("Compact view should contain '%s'", expected)
		}
	}
}

func TestRenderCompactContainsAllChoices(t *testing.T) {
	model := InitialModel()
	compact := model.RenderCompact()

	for _, choice := range model.choices {
		if !strings.Contains(compact, choice.Display) {
			t.Errorf("Compact view should contain choice: %s", choice.Display)
		}
		if !strings.Contains(compact, choice.Icon) {
			t.Errorf("Compact view should contain icon: %s", choice.Icon)
		}
	}
}

func TestColorScheme(t *testing.T) {
	// Test that color constants are defined
	colors := []struct {
		name  string
		value string
	}{
		{"nuclearGreen", string(nuclearGreen)},
		{"glowGreen", string(glowGreen)},
		{"darkGreen", string(darkGreen)},
		{"deepBlack", string(deepBlack)},
		{"charcoal", string(charcoal)},
		{"neonCyan", string(neonCyan)},
		{"terminalGreen", string(terminalGreen)},
		{"dimGreen", string(dimGreen)},
		{"brightText", string(brightText)},
		{"dimText", string(dimText)},
		{"grayText", string(grayText)},
	}

	for _, color := range colors {
		if color.value == "" {
			t.Errorf("Color %s should not be empty", color.name)
		}
		if !strings.HasPrefix(color.value, "#") {
			t.Errorf("Color %s should start with # (hex color), got %s", color.name, color.value)
		}
	}
}

func TestNavigationWithVimKeys(t *testing.T) {
	model := InitialModel()
	model.cursor = 2

	// Test Ctrl+K (up in vim)
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlK})
	m := updatedModel.(Model)

	if m.cursor != 1 {
		t.Errorf("Ctrl+K should move cursor up, expected 1, got %d", m.cursor)
	}

	// Test Ctrl+J (down in vim)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})
	m = updatedModel.(Model)

	if m.cursor != 2 {
		t.Errorf("Ctrl+J should move cursor down, expected 2, got %d", m.cursor)
	}
}

func TestMultipleNavigationSteps(t *testing.T) {
	model := InitialModel()

	// Navigate down 3 times
	for i := 0; i < 3; i++ {
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = updatedModel.(Model)
	}

	if model.cursor != 3 {
		t.Errorf("After 3 down presses, expected cursor at 3, got %d", model.cursor)
	}

	// Navigate up 2 times
	for i := 0; i < 2; i++ {
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
		model = updatedModel.(Model)
	}

	if model.cursor != 1 {
		t.Errorf("After 2 up presses, expected cursor at 1, got %d", model.cursor)
	}
}

func TestFullNavigationCycle(t *testing.T) {
	model := InitialModel()

	choiceCount := len(model.choices)

	// Navigate through all choices
	for i := 0; i < choiceCount*2; i++ {
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = updatedModel.(Model)
	}

	// Should wrap back to start
	expectedCursor := 0
	if model.cursor != expectedCursor {
		t.Errorf("After full cycle, expected cursor at %d, got %d", expectedCursor, model.cursor)
	}
}

func TestViewContainsStatusBar(t *testing.T) {
	model := InitialModel()
	model.cursor = 2
	view := model.View()

	// Status bar should show cursor position
	if !strings.Contains(view, "CURSOR") {
		t.Error("View should contain status bar with CURSOR info")
	}

	if !strings.Contains(view, "SYSTEM READY") {
		t.Error("View should contain SYSTEM READY status")
	}

	if !strings.Contains(view, "MODE: INTERACTIVE") {
		t.Error("View should contain MODE: INTERACTIVE")
	}
}

func TestViewContainsASCIIArt(t *testing.T) {
	model := InitialModel()
	view := model.View()

	// ASCII art logo should be present
	artElements := []string{
		"███",
		"╔",
		"╗",
		"═",
	}

	for _, element := range artElements {
		if !strings.Contains(view, element) {
			t.Errorf("View should contain ASCII art element: %s", element)
		}
	}
}

func TestModelStatePersistence(t *testing.T) {
	model := InitialModel()
	model.cursor = 3
	model.width = 100
	model.height = 50

	// State should persist through non-modifying operations
	view := model.View()

	if view == "" {
		t.Error("View should not be empty")
	}

	if model.cursor != 3 {
		t.Error("Cursor should persist after View()")
	}

	if model.width != 100 {
		t.Error("Width should persist after View()")
	}

	if model.height != 50 {
		t.Error("Height should persist after View()")
	}
}
