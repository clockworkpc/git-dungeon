package screen

import "github.com/charmbracelet/lipgloss"

var (
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F1C40F")).
			Padding(0, 1)

	StyleRegion = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#95A5A6")).
			MarginTop(1)

	StyleLevelItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BDC3C7"))

	StyleLevelSelected = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#F1C40F"))

	StyleCompleted = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2ECC71"))

	StyleObjective = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ECF0F1")).
			Bold(true)

	StyleHint = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F39C12")).
			Italic(true)

	StylePrompt = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2ECC71")).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E74C3C"))

	StyleSuccess = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2ECC71")).
			Bold(true)

	StyleOutput = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BDC3C7"))

	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#34495E"))

	StylePanelTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7F8C8D"))

	StyleStatusBar = lipgloss.NewStyle().
			Background(lipgloss.Color("#2C3E50")).
			Foreground(lipgloss.Color("#ECF0F1")).
			Padding(0, 1)

	StyleKeyHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7F8C8D"))

	StyleExplanation = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#BDC3C7")).
				Padding(0, 2)
)
