package screen

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/clockworkpc/git-dungeon/internal/git"
	"github.com/clockworkpc/git-dungeon/internal/level"
)

type CommandResult struct {
	Input  string
	Output string
	IsErr  bool
}

type LevelState struct {
	Level     level.LevelDef
	LevelNum  int
	Total     int
	Repo      git.RepoState
	History   []CommandResult
	Input     textinput.Model
	GraphVP   viewport.Model
	HistoryVP viewport.Model
	Hint      string
	Width     int
	Height    int
}

func RenderLevel(s LevelState) string {
	// objective bar
	objText := strings.TrimSpace(s.Level.Objective)
	levelInfo := StylePanelTitle.Render(fmt.Sprintf("Level %d/%d — %s", s.LevelNum, s.Total, s.Level.Title))
	hintText := ""
	if s.Hint != "" {
		hintText = "  " + StyleHint.Render("[?] "+s.Hint)
	}
	objLine := StyleObjective.Render(objText) + "  " + levelInfo + hintText
	objBar := StyleBorder.Width(s.Width - 4).Render(objLine)

	// status bar
	branch := s.Repo.CurrentBranch()
	statusBar := StyleStatusBar.Width(s.Width).Render("  branch: " + branch)

	// panels
	panelH := s.Height - lipgloss.Height(objBar) - lipgloss.Height(statusBar) - 5
	halfW := (s.Width - 6) / 2

	graphContent := s.GraphVP.View()
	histContent := renderHistory(s.History, halfW)

	graphPanel := StyleBorder.Width(halfW).Height(panelH).Render(
		StylePanelTitle.Render("GRAPH") + "\n" + graphContent,
	)
	histPanel := StyleBorder.Width(halfW).Height(panelH).Render(
		StylePanelTitle.Render("HISTORY") + "\n" + histContent,
	)

	panels := lipgloss.JoinHorizontal(lipgloss.Top, graphPanel, "  ", histPanel)

	// input
	inputBar := StyleBorder.Width(s.Width - 4).Render(
		StylePrompt.Render("$ ") + s.Input.View(),
	)

	helpBar := StyleKeyHelp.Render("Enter: run  ?: hint  r: restart  Esc: map")

	return lipgloss.JoinVertical(lipgloss.Left,
		objBar,
		statusBar,
		panels,
		inputBar,
		helpBar,
	)
}

func RenderFiles(repo git.RepoState) string {
	var sb strings.Builder

	for f := range repo.StageArea {
		sb.WriteString(StyleSuccess.Render("  staged:    ") + f + "\n")
	}

	for f, status := range repo.WorkingTree.Files {
		switch status {
		case git.FileModified:
			sb.WriteString(StyleError.Render("  modified:  ") + f + "\n")
		case git.FileUntracked:
			sb.WriteString(StyleOutput.Render("  untracked: ") + f + "\n")
		case git.FileDeleted:
			sb.WriteString(StyleError.Render("  deleted:   ") + f + "\n")
		}
	}

	return sb.String()
}

func renderHistory(history []CommandResult, width int) string {
	var sb strings.Builder
	for _, r := range history {
		sb.WriteString(StylePrompt.Render("$ ") + r.Input + "\n")
		if r.IsErr {
			sb.WriteString(StyleError.Render(r.Output) + "\n")
		} else {
			sb.WriteString(StyleOutput.Render(r.Output) + "\n")
		}
	}
	return lipgloss.NewStyle().Width(width).Render(sb.String())
}
