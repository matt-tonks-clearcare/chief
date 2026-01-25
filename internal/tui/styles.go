package tui

import "github.com/charmbracelet/lipgloss"

// Basic styles for the dashboard.
// Note: US-009 will expand this into a full styling system.
var (
	// Colors (basic palette - will be refined in US-009)
	primaryColor = lipgloss.Color("#00D7FF")
	successColor = lipgloss.Color("#5AF78E")
	warningColor = lipgloss.Color("#F3F99D")
	errorColor   = lipgloss.Color("#FF5C57")
	mutedColor   = lipgloss.Color("#6C7086")
	borderColor  = lipgloss.Color("#45475A")

	// Header style
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1)

	// Footer style
	footerStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(0, 1)

	// Panel styles
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1)

	// Selected item style
	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#313244")).
			Foreground(lipgloss.Color("#CDD6F4"))

	// Story status styles
	statusPassedStyle     = lipgloss.NewStyle().Foreground(successColor)
	statusInProgressStyle = lipgloss.NewStyle().Foreground(primaryColor)
	statusPendingStyle    = lipgloss.NewStyle().Foreground(mutedColor)
	statusFailedStyle     = lipgloss.NewStyle().Foreground(errorColor)

	// Title style for details panel
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#CDD6F4"))

	// Label style
	labelStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Progress bar styles
	progressBarFillStyle  = lipgloss.NewStyle().Foreground(successColor)
	progressBarEmptyStyle = lipgloss.NewStyle().Foreground(mutedColor)
)

// Status icons
const (
	iconPassed     = "✓"
	iconInProgress = "●"
	iconPending    = "○"
	iconFailed     = "✗"
)

// GetStatusIcon returns the appropriate icon for a story's status.
func GetStatusIcon(passed, inProgress bool) string {
	if passed {
		return statusPassedStyle.Render(iconPassed)
	}
	if inProgress {
		return statusInProgressStyle.Render(iconInProgress)
	}
	return statusPendingStyle.Render(iconPending)
}
