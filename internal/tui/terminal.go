package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CommandMsg is sent when a command is executed
type CommandMsg string

// Terminal represents the command terminal component
type Terminal struct {
	input      textinput.Model
	viewport   viewport.Model
	history    []string
	historyIdx int
	output     strings.Builder
	width      int
	height     int
	ready      bool
}

// NewTerminal creates a new terminal component
func NewTerminal() *Terminal {
	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60
	ti.Prompt = "$ "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	return &Terminal{
		input:      ti,
		history:    []string{},
		historyIdx: -1,
	}
}

// Init initializes the terminal
func (t *Terminal) Init() tea.Cmd {
	return textinput.Blink
}

// SetSize sets the terminal dimensions
func (t *Terminal) SetSize(width, height int) {
	t.width = width
	t.height = height

	// Create viewport for output if not exists
	if !t.ready {
		t.viewport = viewport.New(width-2, height-5)
		t.viewport.Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
		t.viewport.MouseWheelEnabled = true
		t.ready = true
	} else {
		t.viewport.Width = width - 2
		t.viewport.Height = height - 5
	}

	t.input.Width = width - 6
}

// Update handles messages
func (t *Terminal) Update(msg tea.Msg) (*Terminal, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.MouseMsg:
		// Handle mouse wheel scrolling
		switch msg.Type {
		case tea.MouseWheelUp:
			t.viewport.LineUp(3)
			return t, nil
		case tea.MouseWheelDown:
			t.viewport.LineDown(3)
			return t, nil
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmd := strings.TrimSpace(t.input.Value())
			if cmd != "" {
				t.history = append(t.history, cmd)
				t.historyIdx = len(t.history)
				t.input.SetValue("")
				return t, func() tea.Msg { return CommandMsg(cmd) }
			}

		case "up":
			// Navigate history backward
			if len(t.history) > 0 && t.historyIdx > 0 {
				t.historyIdx--
				t.input.SetValue(t.history[t.historyIdx])
				t.input.CursorEnd()
			}
			return t, nil

		case "down":
			// Navigate history forward
			if t.historyIdx < len(t.history)-1 {
				t.historyIdx++
				t.input.SetValue(t.history[t.historyIdx])
				t.input.CursorEnd()
			} else if t.historyIdx == len(t.history)-1 {
				t.historyIdx = len(t.history)
				t.input.SetValue("")
			}
			return t, nil

		case "pgup", "shift+up":
			// Scroll viewport up
			t.viewport.LineUp(5)
			return t, nil

		case "pgdown", "shift+down":
			// Scroll viewport down
			t.viewport.LineDown(5)
			return t, nil

		case "home", "shift+home":
			// Scroll to top
			t.viewport.GotoTop()
			return t, nil

		case "end", "shift+end":
			// Scroll to bottom
			t.viewport.GotoBottom()
			return t, nil

		case "ctrl+l":
			// Clear terminal
			t.output.Reset()
			t.viewport.SetContent("")
			return t, nil
		}
	}

	// Update text input
	var cmd tea.Cmd
	t.input, cmd = t.input.Update(msg)
	cmds = append(cmds, cmd)

	// Update viewport
	t.viewport, cmd = t.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return t, tea.Batch(cmds...)
}

// AddOutput adds command output to the terminal
func (t *Terminal) AddOutput(cmd, output string) {
	// Add command prompt
	t.output.WriteString(fmt.Sprintf("\n$ %s\n", cmd))

	// Add output
	if output != "" {
		t.output.WriteString(output)
		t.output.WriteString("\n")
	}

	// Update viewport and scroll to bottom
	content := t.output.String()
	t.viewport.SetContent(content)

	// Calculate total lines and scroll to bottom
	lines := strings.Count(content, "\n")
	if lines > t.viewport.Height {
		t.viewport.SetYOffset(lines - t.viewport.Height + 1)
	}
}

// Clear clears the terminal output
func (t *Terminal) Clear() {
	t.output.Reset()
	t.viewport.SetContent("")
}

// View renders the terminal
func (t *Terminal) View() string {
	if !t.ready {
		return "Loading..."
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1).
		Width(t.width).
		Height(t.height)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1)

	title := titleStyle.Render("Terminal")

	// Separator line
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		Render(strings.Repeat("─", t.width-4))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		separator,
		t.viewport.View(),
		separator,
		t.input.View(),
	)

	return borderStyle.Render(content)
}
