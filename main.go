package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/clockworkpc/git-dungeon/content"
	"github.com/clockworkpc/git-dungeon/internal/level"
	"github.com/clockworkpc/git-dungeon/internal/model"
	"github.com/clockworkpc/git-dungeon/internal/progress"
)

func main() {
	levels, err := level.LoadAll(content.FS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading levels: %v\n", err)
		os.Exit(1)
	}

	prog, err := progress.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load progress: %v\n", err)
		prog, _ = progress.Load()
	}

	m := model.New(levels, prog)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
