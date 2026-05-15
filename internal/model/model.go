package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clockworkpc/git-dungeon/internal/git"
	"github.com/clockworkpc/git-dungeon/internal/graph"
	"github.com/clockworkpc/git-dungeon/internal/level"
	"github.com/clockworkpc/git-dungeon/internal/messages"
	"github.com/clockworkpc/git-dungeon/internal/parser"
	"github.com/clockworkpc/git-dungeon/internal/progress"
	"github.com/clockworkpc/git-dungeon/internal/screen"
)

type Screen int

const (
	ScreenMap Screen = iota
	ScreenLevel
	ScreenResults
)

type Model struct {
	currentScreen Screen
	width         int
	height        int

	levels     []level.LevelDef
	levelIndex int
	repo       git.RepoState
	initialRepo git.RepoState

	input        textinput.Model
	history      []screen.CommandResult
	inputHistory []string
	inputHistIdx int

	graphVP   viewport.Model
	historyVP viewport.Model

	mapCursor  int
	hintIndex  int
	lastCmdSub string

	prog progress.Progress
}

func New(levels []level.LevelDef, prog progress.Progress) Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 60

	return Model{
		currentScreen: ScreenMap,
		levels:        levels,
		input:         ti,
		prog:          prog,
		graphVP:       viewport.New(40, 20),
		historyVP:     viewport.New(40, 20),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		halfW := (m.width - 6) / 2
		panelH := m.height - 12
		if panelH < 5 {
			panelH = 5
		}
		m.graphVP.Width = halfW - 4
		m.graphVP.Height = panelH - 2
		m.historyVP.Width = halfW - 4
		m.historyVP.Height = panelH - 2
		m.input.Width = m.width - 10

	case tea.KeyMsg:
		return m.handleKey(msg)

	case messages.LevelSelectedMsg:
		return m.loadLevel(msg.Index)

	case messages.LevelPassedMsg:
		m.prog.MarkCompleted(m.levels[m.levelIndex].ID)
		_ = m.prog.Save()
		m.currentScreen = ScreenResults

	case messages.RestartLevelMsg:
		return m.restartLevel()

	case messages.NextLevelMsg:
		next := m.levelIndex + 1
		if next >= len(m.levels) {
			m.currentScreen = ScreenMap
			return m, nil
		}
		return m.loadLevel(next)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.currentScreen {
	case ScreenMap:
		return m.handleMapKey(msg)
	case ScreenLevel:
		return m.handleLevelKey(msg)
	case ScreenResults:
		return m.handleResultsKey(msg)
	}
	return m, nil
}

func (m Model) handleMapKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.mapCursor > 0 {
			m.mapCursor--
		}
	case "down", "j":
		if m.mapCursor < len(m.levels)-1 {
			m.mapCursor++
		}
	case "enter":
		return m.loadLevel(m.mapCursor)
	}
	return m, nil
}

func (m Model) handleLevelKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.currentScreen = ScreenMap
		return m, nil
	case "r":
		return m.restartLevel()
	case "?":
		lv := m.levels[m.levelIndex]
		if len(lv.Hints) > 0 {
			m.hintIndex = (m.hintIndex + 1) % len(lv.Hints)
		}
		return m, nil
	case "up":
		if len(m.inputHistory) > 0 {
			if m.inputHistIdx > 0 {
				m.inputHistIdx--
			}
			m.input.SetValue(m.inputHistory[m.inputHistIdx])
			return m, nil
		}
	case "down":
		if m.inputHistIdx < len(m.inputHistory)-1 {
			m.inputHistIdx++
			m.input.SetValue(m.inputHistory[m.inputHistIdx])
		} else {
			m.input.SetValue("")
		}
		return m, nil
	case "enter":
		raw := strings.TrimSpace(m.input.Value())
		m.input.SetValue("")
		if raw == "" {
			return m, nil
		}
		m.inputHistory = append(m.inputHistory, raw)
		m.inputHistIdx = len(m.inputHistory)
		return m.submitCommand(raw)
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) handleResultsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "n", "enter":
		return m.Update(messages.NextLevelMsg{})
	case "r":
		return m.restartLevel()
	case "m", "esc":
		m.currentScreen = ScreenMap
	}
	return m, nil
}

func (m Model) loadLevel(index int) (tea.Model, tea.Cmd) {
	if index < 0 || index >= len(m.levels) {
		return m, nil
	}
	m.levelIndex = index
	m.mapCursor = index
	lv := m.levels[index]
	repo := level.ToRepoState(lv.Initial)
	m.repo = repo
	m.initialRepo = repo.DeepCopy()
	m.history = nil
	m.inputHistory = nil
	m.inputHistIdx = 0
	m.hintIndex = 0
	m.lastCmdSub = ""
	m.input.SetValue("")
	m.input.Focus()
	m.updateGraph()
	m.currentScreen = ScreenLevel
	return m, textinput.Blink
}

func (m Model) restartLevel() (tea.Model, tea.Cmd) {
	return m.loadLevel(m.levelIndex)
}

func (m *Model) updateGraph() {
	content := graph.Render(m.repo)
	files := screen.RenderFiles(m.repo)
	if files != "" {
		content += "\n\n" + screen.StylePanelTitle.Render("FILES") + "\n" + files
	}
	m.graphVP.SetContent(content)
}

func (m Model) submitCommand(raw string) (tea.Model, tea.Cmd) {
	lv := m.levels[m.levelIndex]

	cmd, err := parser.Parse(raw)
	if err != nil {
		m.history = append(m.history, screen.CommandResult{
			Input:  raw,
			Output: err.Error(),
			IsErr:  true,
		})
		return m, nil
	}

	if !parser.IsAllowed(cmd, lv.AllowedCommands) {
		m.history = append(m.history, screen.CommandResult{
			Input:  raw,
			Output: "Command not allowed in this level.",
			IsErr:  true,
		})
		return m, nil
	}

	newRepo, output, applyErr := git.Apply(m.repo, cmd)
	isErr := applyErr != nil
	if isErr {
		output = applyErr.Error()
	} else {
		m.repo = newRepo
		m.updateGraph()
	}

	m.lastCmdSub = cmd.Sub
	m.history = append(m.history, screen.CommandResult{
		Input:  raw,
		Output: output,
		IsErr:  isErr,
	})

	if !isErr {
		result := level.Validate(lv.Goal, m.initialRepo, m.repo, m.lastCmdSub)
		if result.Passed {
			return m, func() tea.Msg { return messages.LevelPassedMsg{} }
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}
	switch m.currentScreen {
	case ScreenMap:
		return screen.RenderMap(screen.MapState{
			Levels:   m.levels,
			Cursor:   m.mapCursor,
			Progress: m.prog,
			Width:    m.width,
			Height:   m.height,
		})
	case ScreenLevel:
		lv := m.levels[m.levelIndex]
		hint := ""
		if len(lv.Hints) > 0 {
			hint = lv.Hints[m.hintIndex%len(lv.Hints)]
		}
		return screen.RenderLevel(screen.LevelState{
			Level:     lv,
			LevelNum:  m.levelIndex + 1,
			Total:     len(m.levels),
			Repo:      m.repo,
			History:   m.history,
			Input:     m.input,
			GraphVP:   m.graphVP,
			HistoryVP: m.historyVP,
			Hint:      hint,
			Width:     m.width,
			Height:    m.height,
		})
	case ScreenResults:
		lv := m.levels[m.levelIndex]
		return screen.RenderResults(screen.ResultsState{
			Level:  lv,
			Width:  m.width,
			Height: m.height,
		})
	}
	return ""
}
