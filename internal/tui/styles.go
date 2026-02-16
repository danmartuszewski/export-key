package tui

import "github.com/charmbracelet/lipgloss"

// Color palette aligned with hop
var (
	primaryColor   = lipgloss.Color("39")  // Cyan
	secondaryColor = lipgloss.Color("245") // Gray
	accentColor    = lipgloss.Color("212") // Pink
	successColor   = lipgloss.Color("82")  // Green
	warningColor   = lipgloss.Color("214") // Orange
	dimColor       = lipgloss.Color("240") // Dark gray
	mutedColor     = lipgloss.Color("241") // Muted gray
)

// Picker styles (TUI interactive mode)
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	versionStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	promptStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	numberStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Width(4).
			Align(lipgloss.Right)

	selectedNumberStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Width(4).
				Align(lipgloss.Right)

	envVarStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	selectedEnvVarStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	titleColumnStyle = lipgloss.NewStyle().
				Foreground(secondaryColor)

	titleSuffixStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("237"))

	selectedTitleColumnStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("255"))

	selectedTitleSuffixStyle = lipgloss.NewStyle().
					Foreground(secondaryColor)

	projectStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	selectedProjectStyle = lipgloss.NewStyle().
				Foreground(warningColor).
				Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	filterInputStyle = lipgloss.NewStyle().
				Foreground(accentColor)

	checkStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	uncheckStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	selectedCountStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true)
)

// List styles (non-interactive list command)
var (
	ListHeaderStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	ListNumberStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	ListEnvVarStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	ListTitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	ListTitleSuffixStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("237"))

	ListProjectStyle = lipgloss.NewStyle().
				Foreground(warningColor)

	ListProjectTagStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Background(lipgloss.Color("236")).
				Padding(0, 1)

	ListSeparatorStyle = lipgloss.NewStyle().
				Foreground(dimColor)
)
