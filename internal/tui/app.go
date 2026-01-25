package tui

import (
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/minicodemonkey/chief/internal/prd"
)

// AppState represents the current state of the application.
type AppState int

const (
	StateReady AppState = iota
	StateRunning
	StatePaused
	StateStopped
	StateComplete
	StateError
)

func (s AppState) String() string {
	switch s {
	case StateReady:
		return "Ready"
	case StateRunning:
		return "Running"
	case StatePaused:
		return "Paused"
	case StateStopped:
		return "Stopped"
	case StateComplete:
		return "Complete"
	case StateError:
		return "Error"
	default:
		return "Unknown"
	}
}

// App is the main Bubble Tea model for the Chief TUI.
type App struct {
	prd           *prd.PRD
	prdPath       string
	prdName       string
	state         AppState
	iteration     int
	startTime     time.Time
	selectedIndex int
	width         int
	height        int
	err           error
}

// NewApp creates a new App with the given PRD.
func NewApp(prdPath string) (*App, error) {
	p, err := prd.LoadPRD(prdPath)
	if err != nil {
		return nil, err
	}

	// Extract PRD name from path (directory name or filename without extension)
	prdName := filepath.Base(filepath.Dir(prdPath))
	if prdName == "." || prdName == "/" {
		prdName = filepath.Base(prdPath)
	}

	return &App{
		prd:           p,
		prdPath:       prdPath,
		prdName:       prdName,
		state:         StateReady,
		iteration:     0,
		selectedIndex: 0,
	}, nil
}

// Init initializes the App.
func (a App) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
	)
}

// Update handles messages and updates the model.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit

		// Navigation
		case "up", "k":
			if a.selectedIndex > 0 {
				a.selectedIndex--
			}
		case "down", "j":
			if a.selectedIndex < len(a.prd.UserStories)-1 {
				a.selectedIndex++
			}
		}
	}

	return a, nil
}

// View renders the TUI.
func (a App) View() string {
	return a.renderDashboard()
}

// GetPRD returns the current PRD.
func (a *App) GetPRD() *prd.PRD {
	return a.prd
}

// GetSelectedStory returns the currently selected story.
func (a *App) GetSelectedStory() *prd.UserStory {
	if a.selectedIndex >= 0 && a.selectedIndex < len(a.prd.UserStories) {
		return &a.prd.UserStories[a.selectedIndex]
	}
	return nil
}

// GetState returns the current app state.
func (a *App) GetState() AppState {
	return a.state
}

// GetIteration returns the current iteration count.
func (a *App) GetIteration() int {
	return a.iteration
}

// GetElapsedTime returns the elapsed time since the loop started.
func (a *App) GetElapsedTime() time.Duration {
	if a.startTime.IsZero() {
		return 0
	}
	return time.Since(a.startTime)
}

// GetCompletionPercentage returns the percentage of completed stories.
func (a *App) GetCompletionPercentage() float64 {
	if len(a.prd.UserStories) == 0 {
		return 100.0
	}
	var completed int
	for _, s := range a.prd.UserStories {
		if s.Passes {
			completed++
		}
	}
	return float64(completed) / float64(len(a.prd.UserStories)) * 100.0
}
