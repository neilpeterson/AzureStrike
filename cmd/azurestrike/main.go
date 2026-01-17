package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/azurestrike/azurestrike/internal/game"
	"github.com/azurestrike/azurestrike/internal/scenario"
	"github.com/azurestrike/azurestrike/internal/tui"
)

func main() {
	scenarioFlag := flag.String("scenario", "", "Scenario to load (e.g., storage-recon)")
	listFlag := flag.Bool("list", false, "List available scenarios")
	flag.Parse()

	if *listFlag {
		scenarios, err := scenario.ListScenarios("scenarios")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing scenarios: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Available scenarios:")
		for _, s := range scenarios {
			fmt.Printf("  - %s: %s (%s)\n", s.ID, s.Name, s.Difficulty)
		}
		return
	}

	var sc *scenario.Scenario
	var err error

	if *scenarioFlag != "" {
		sc, err = scenario.LoadByID("scenarios", *scenarioFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading scenario: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Load default scenario
		sc, err = scenario.LoadByID("scenarios", "storage-breach")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading default scenario: %v\n", err)
			fmt.Fprintf(os.Stderr, "Use --list to see available scenarios\n")
			os.Exit(1)
		}
	}

	// Initialize game state with scenario
	gameState := game.NewState(sc)

	// Create and run TUI
	app := tui.NewApp(gameState)
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
