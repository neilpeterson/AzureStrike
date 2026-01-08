package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/azurestrike/azurestrike/internal/game"
)

// Status represents the status panel component
type Status struct {
	gameState *game.State
	width     int
	height    int
}

// NewStatus creates a new status panel
func NewStatus(state *game.State) *Status {
	return &Status{
		gameState: state,
	}
}

// SetSize sets the panel dimensions
func (s *Status) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// Init initializes the status panel
func (s *Status) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (s *Status) Update(msg tea.Msg) (*Status, tea.Cmd) {
	return s, nil
}

// View renders the status panel
func (s *Status) View() string {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1).
		Width(s.width).
		Height(s.height)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("75")).
		MarginTop(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	completedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("34"))

	var content strings.Builder

	// Title
	content.WriteString(titleStyle.Render("=== STATUS ==="))
	content.WriteString("\n")

	// Scenario info
	content.WriteString(sectionStyle.Render("MISSION"))
	content.WriteString("\n")
	content.WriteString(labelStyle.Render("Name: "))
	content.WriteString(valueStyle.Render(truncate(s.gameState.Scenario.Name, s.width-10)))
	content.WriteString("\n")

	// Timer
	elapsed := s.gameState.GetElapsedTime().Round(time.Second)
	content.WriteString(labelStyle.Render("Time: "))
	content.WriteString(valueStyle.Render(elapsed.String()))
	content.WriteString("\n")

	// Score
	content.WriteString(sectionStyle.Render("SCORE"))
	content.WriteString("\n")
	content.WriteString(labelStyle.Render("Points: "))
	content.WriteString(valueStyle.Render(fmt.Sprintf("%d / %d",
		s.gameState.Score.Points,
		s.gameState.Scenario.TotalPoints())))
	content.WriteString("\n")

	// Progress bar
	progress := s.gameState.GetProgress()
	progressBar := renderProgressBar(progress, s.width-8)
	content.WriteString(labelStyle.Render("Progress: "))
	content.WriteString("\n")
	content.WriteString(progressBar)
	content.WriteString("\n")

	// Objectives summary
	pending := s.gameState.GetPendingObjectives()
	completed := s.gameState.GetCompletedObjectives()
	total := len(pending) + len(completed)

	content.WriteString(sectionStyle.Render("OBJECTIVES"))
	content.WriteString("\n")
	content.WriteString(completedStyle.Render(fmt.Sprintf("%d", len(completed))))
	content.WriteString(labelStyle.Render(fmt.Sprintf(" / %d complete", total)))
	content.WriteString("\n")

	// Hints used
	if len(s.gameState.HintsUsed) > 0 {
		content.WriteString(labelStyle.Render(fmt.Sprintf("%d hint(s) used",
			len(s.gameState.HintsUsed))))
		content.WriteString("\n")
	}

	// Quick commands
	content.WriteString(sectionStyle.Render("COMMANDS"))
	content.WriteString("\n")
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("117"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	commands := []struct{ cmd, desc string }{
		{"objective", "view objectives"},
		{"hint", "get a hint"},
		{"score", "view score"},
		{"help", "all commands"},
	}
	for _, c := range commands {
		content.WriteString(cmdStyle.Render(c.cmd))
		content.WriteString(descStyle.Render(" - " + c.desc))
		content.WriteString("\n")
	}

	return borderStyle.Render(content.String())
}

// renderProgressBar creates a visual progress bar
func renderProgressBar(progress float64, width int) string {
	if width < 5 {
		width = 5
	}

	filled := int(progress * float64(width))
	empty := width - filled

	filledStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("34")).
		Background(lipgloss.Color("34"))

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("240"))

	bar := filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))

	return fmt.Sprintf("[%s] %.0f%%", bar, progress*100)
}

// truncate shortens a string to fit within maxLen
func truncate(s string, maxLen int) string {
	if maxLen <= 3 {
		return s
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
