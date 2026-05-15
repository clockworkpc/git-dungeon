package screen

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/clockworkpc/git-dungeon/internal/level"
	"github.com/clockworkpc/git-dungeon/internal/progress"
)

type MapState struct {
	Levels   []level.LevelDef
	Cursor   int
	Progress progress.Progress
	Width    int
	Height   int
}

func RenderMap(s MapState) string {
	var sb strings.Builder

	title := StyleTitle.Render("G I T   D U N G E O N")
	sb.WriteString(lipgloss.PlaceHorizontal(s.Width, lipgloss.Center, title))
	sb.WriteString("\n\n")

	currentRegion := ""
	for i, lv := range s.Levels {
		if lv.Region != currentRegion {
			currentRegion = lv.Region
			sb.WriteString(StyleRegion.Render("  " + strings.ToUpper(currentRegion)))
			sb.WriteString("\n")
		}

		completed := s.Progress.IsCompleted(lv.ID)
		selected := i == s.Cursor

		var status string
		if completed {
			status = StyleCompleted.Render("✓")
		} else {
			status = StyleKeyHelp.Render("·")
		}

		cursor := "  "
		var itemStyle lipgloss.Style
		if selected {
			cursor = "▶ "
			itemStyle = StyleLevelSelected
		} else {
			itemStyle = StyleLevelItem
		}

		line := fmt.Sprintf("%s● %s", cursor, itemStyle.Render(lv.Title))
		paddedLine := lipgloss.NewStyle().Width(s.Width - 6).Render(line)
		sb.WriteString(paddedLine + "  " + status + "\n")
	}

	sb.WriteString("\n")
	help := StyleKeyHelp.Render("↑/↓ navigate  Enter: play  q: quit")
	sb.WriteString(lipgloss.PlaceHorizontal(s.Width, lipgloss.Center, help))

	return StyleBorder.Width(s.Width - 2).Render(sb.String())
}
