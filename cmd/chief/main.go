package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/minicodemonkey/chief/internal/tui"
)

func main() {
	// For now, use a default PRD path (will be configurable via CLI flags in US-022)
	prdPath := ".chief/prds/main/prd.json"

	// Check for command-line argument for PRD path
	if len(os.Args) > 1 {
		prdPath = os.Args[1]
	}

	app, err := tui.NewApp(prdPath)
	if err != nil {
		fmt.Printf("Error loading PRD: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
