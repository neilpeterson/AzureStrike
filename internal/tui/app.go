package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/azurestrike/azurestrike/internal/cli"
	"github.com/azurestrike/azurestrike/internal/game"
)

// App is the main application model
type App struct {
	gameState    *game.State
	parser       *cli.Parser
	terminal     *Terminal
	status       *Status
	fireworks    *Fireworks
	width        int
	height       int
	showBriefing bool
	showDebrief  bool
}

// KeyMap defines the key bindings
type KeyMap struct {
	Quit    key.Binding
	Tab     key.Binding
	Enter   key.Binding
	Dismiss key.Binding
}

var keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "execute"),
	),
	Dismiss: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "dismiss"),
	),
}

// NewApp creates a new application instance
func NewApp(state *game.State) *App {
	return &App{
		gameState:    state,
		parser:       cli.NewParser(state),
		terminal:     NewTerminal(),
		status:       NewStatus(state),
		fireworks:    NewFireworks(),
		showBriefing: true,
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		a.terminal.Init(),
	)
}

// Update handles messages and updates state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit
		if key.Matches(msg, keys.Quit) {
			return a, tea.Quit
		}

		// Handle briefing/debrief dismissal
		if a.showBriefing || a.showDebrief {
			if key.Matches(msg, keys.Dismiss) {
				a.showBriefing = false
				a.showDebrief = false
				return a, nil
			}
			return a, nil
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		statusWidth := 32
		a.terminal.SetSize(msg.Width-statusWidth-5, msg.Height-4) // Leave room for status panel
		a.status.SetSize(statusWidth, msg.Height-4)
		a.fireworks.SetSize(msg.Width, msg.Height) // Full screen for celebration

	case FireworksTickMsg:
		var cmd tea.Cmd
		a.fireworks, cmd = a.fireworks.Update(msg)
		return a, cmd

	case FireworksDoneMsg:
		// Fireworks finished, nothing special to do
		return a, nil

	case CommandMsg:
		cmdStr := string(msg)
		// Handle special commands
		switch cmdStr {
		case "clear":
			a.terminal.Clear()
		case "exit", "quit":
			return a, tea.Quit
		case "konami", "fireworks": // Secret commands!
			a.terminal.AddOutput(cmdStr, "🎆 Fireworks mode activated!")
			return a, a.fireworks.Start()
		default:
			// Execute the command
			result := a.parser.Execute(cmdStr)
			a.terminal.AddOutput(cmdStr, result.Output)
		}

		// Check for game completion - fireworks only on full scenario completion
		if a.gameState.Status == game.StatusCompleted {
			a.showDebrief = true
			cmds = append(cmds, a.fireworks.Start())
		}

		// Update status panel
		a.status.Update(nil)
	}

	// Update terminal
	var cmd tea.Cmd
	a.terminal, cmd = a.terminal.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

// View renders the application
func (a *App) View() string {
	if a.showBriefing {
		return a.renderBriefing()
	}

	if a.showDebrief {
		return a.renderDebrief()
	}

	// Show full-screen fireworks celebration if active
	if a.fireworks.IsActive() {
		return a.fireworks.View()
	}

	// Main game view
	terminalView := a.terminal.View()
	statusView := a.status.View()

	// Layout: terminal on left, status on right
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		terminalView,
		statusView,
	)
}

func (a *App) renderBriefing() string {
	style := lipgloss.NewStyle().
		Width(a.width - 4).
		Height(a.height - 4).
		Padding(2, 4).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("63"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(2)

	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		MarginTop(2)

	content := titleStyle.Render("=== MISSION BRIEFING ===") + "\n\n"
	content += titleStyle.Render(a.gameState.Scenario.Name) + "\n"
	content += lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("Difficulty: "+a.gameState.Scenario.Difficulty) + "\n\n"
	content += contentStyle.Render(a.gameState.Scenario.Briefing) + "\n"
	content += hintStyle.Render("\n[Press ENTER to begin]")

	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center,
		lipgloss.Center,
		style.Render(content),
	)
}

func (a *App) renderDebrief() string {
	style := lipgloss.NewStyle().
		Width(a.width - 4).
		Height(a.height - 4).
		Padding(2, 4).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("34"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("34")).
		MarginBottom(2)

	content := titleStyle.Render("=== MISSION COMPLETE ===") + "\n\n"
	content += lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226")).
		Render("SCORE: ") +
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255")).
			Render(string(rune(a.gameState.Score.Points+'0'))) + " points\n\n"

	content += titleStyle.Render("Debrief:") + "\n"
	content += a.gameState.Scenario.Debrief + "\n"

	content += lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		MarginTop(2).
		Render("\n[Press ENTER to continue]")

	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center,
		lipgloss.Center,
		style.Render(content),
	)
}
