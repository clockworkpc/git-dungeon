package screen

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/clockworkpc/git-dungeon/internal/level"
)

type ResultsState struct {
	Level  level.LevelDef
	Width  int
	Height int
}

func RenderResults(s ResultsState) string {
	var sb strings.Builder

	sb.WriteString(StyleSuccess.Render("  LEVEL COMPLETE!"))
	sb.WriteString("\n\n")
	sb.WriteString(StyleObjective.Render("  \"" + s.Level.Title + "\""))
	sb.WriteString("\n\n")

	if s.Level.Explanation != "" {
		wrapped := lipgloss.NewStyle().Width(s.Width - 8).Render(s.Level.Explanation)
		sb.WriteString(StyleExplanation.Render(wrapped))
		sb.WriteString("\n\n")
	}

	sb.WriteString(StyleKeyHelp.Render("  [N] Next Level    [R] Retry    [M] Map"))

	return StyleBorder.Width(s.Width - 4).Render(sb.String())
}
