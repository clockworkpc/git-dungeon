package graph

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/clockworkpc/git-dungeon/internal/git"
)

const (
	nodeChar  = "●"
	vertChar  = "│"
	horizChar = "─"
	branchOut = "╮"
	mergeIn   = "╯"
	junction  = "├"
)

var laneColors = []lipgloss.Color{
	"#F1C40F", "#3498DB", "#2ECC71", "#9B59B6", "#E67E22", "#1ABC9C",
}

var (
	styleHEAD    = lipgloss.NewStyle().Foreground(lipgloss.Color("#E74C3C")).Bold(true)
	styleBranch  = lipgloss.NewStyle().Foreground(lipgloss.Color("#2ECC71")).Bold(true)
	styleRemote  = lipgloss.NewStyle().Foreground(lipgloss.Color("#3498DB"))
	styleMessage = lipgloss.NewStyle().Foreground(lipgloss.Color("#BDC3C7"))
)

type commitPos struct {
	row int
	col int
}

func Render(repo git.RepoState) string {
	order := topoOrder(repo)
	if len(order) == 0 {
		return "(empty repository)"
	}

	pos := make(map[string]commitPos, len(order))
	lanes := []string{}

	for row, id := range order {
		col := -1
		for i, lane := range lanes {
			if lane == id {
				col = i
				break
			}
		}
		if col == -1 {
			col = len(lanes)
			lanes = append(lanes, id)
		}
		pos[id] = commitPos{row: row, col: col}

		c := repo.Commits[id]
		if len(c.Parents) == 1 {
			parent := c.Parents[0]
			for i, lane := range lanes {
				if lane == id {
					lanes[i] = parent
					break
				}
			}
		} else if len(c.Parents) == 0 {
			for i, lane := range lanes {
				if lane == id {
					lanes = append(lanes[:i], lanes[i+1:]...)
					break
				}
			}
		} else {
			for i, lane := range lanes {
				if lane == id {
					lanes[i] = c.Parents[0]
					break
				}
			}
			for _, p := range c.Parents[1:] {
				found := false
				for _, l := range lanes {
					if l == p {
						found = true
						break
					}
				}
				if !found {
					lanes = append(lanes, p)
				}
			}
		}
	}

	// build reversed tip→branches map
	branchAt := map[string][]string{}
	for name, id := range repo.Branches {
		branchAt[id] = append(branchAt[id], name)
	}
	remoteAt := map[string][]string{}
	for name, id := range repo.RemoteRefs {
		remoteAt[id] = append(remoteAt[id], name)
	}

	var lines []string
	maxCol := 0
	for _, p := range pos {
		if p.col > maxCol {
			maxCol = p.col
		}
	}

	for _, id := range order {
		p := pos[id]
		c := repo.Commits[id]

		// node line
		nodeLine := buildNodeLine(id, p.col, maxCol, c, repo, branchAt, remoteAt)
		lines = append(lines, nodeLine)

		// connector line
		if len(c.Parents) > 0 {
			conn := buildConnectorLine(id, p, c.Parents, pos, maxCol)
			lines = append(lines, conn)
		}
	}

	return strings.Join(lines, "\n")
}

func buildNodeLine(id string, col, maxCol int, c git.Commit, repo git.RepoState,
	branchAt, remoteAt map[string][]string) string {

	colColor := laneColors[col%len(laneColors)]
	nodeStyle := lipgloss.NewStyle().Foreground(colColor).Bold(true)

	var sb strings.Builder
	for i := 0; i <= maxCol; i++ {
		if i == col {
			sb.WriteString(nodeStyle.Render(nodeChar))
		} else {
			edgeStyle := lipgloss.NewStyle().Foreground(laneColors[i%len(laneColors)])
			sb.WriteString(edgeStyle.Render(vertChar))
		}
		if i < maxCol {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("  ")

	// labels
	var labels []string
	isHead := (repo.HeadBranch != "" && repo.Branches[repo.HeadBranch] == id) ||
		(repo.HeadBranch == "" && repo.HeadCommit == id)

	branches := branchAt[id]
	sort.Strings(branches)
	remotes := remoteAt[id]
	sort.Strings(remotes)

	if len(branches) > 0 || len(remotes) > 0 || isHead {
		var parts []string
		if isHead {
			headLabel := "HEAD"
			if repo.HeadBranch != "" {
				headLabel = fmt.Sprintf("HEAD -> %s", repo.HeadBranch)
			}
			parts = append(parts, styleHEAD.Render(headLabel))
		}
		for _, b := range branches {
			if b != repo.HeadBranch {
				parts = append(parts, styleBranch.Render(b))
			}
		}
		for _, r := range remotes {
			parts = append(parts, styleRemote.Render(r))
		}
		labels = append(labels, "("+strings.Join(parts, ", ")+")")
	}

	sb.WriteString(strings.Join(labels, " "))
	if len(labels) > 0 {
		sb.WriteString(" ")
	}
	sb.WriteString(styleMessage.Render(c.Message))

	return sb.String()
}

func buildConnectorLine(id string, p commitPos, parents []string, pos map[string]commitPos, maxCol int) string {
	var sb strings.Builder
	for i := 0; i <= maxCol; i++ {
		colColor := laneColors[i%len(laneColors)]
		edgeStyle := lipgloss.NewStyle().Foreground(colColor)

		if i == p.col {
			if len(parents) == 1 {
				sb.WriteString(edgeStyle.Render(vertChar))
			} else {
				sb.WriteString(edgeStyle.Render(junction))
			}
		} else {
			// check if any parent is at a different column
			parentHere := false
			for _, pid := range parents {
				if pp, ok := pos[pid]; ok && pp.col == i {
					parentHere = true
					break
				}
			}
			if parentHere && i != p.col {
				edgeStyle2 := lipgloss.NewStyle().Foreground(laneColors[i%len(laneColors)])
				sb.WriteString(edgeStyle2.Render(mergeIn))
			} else {
				sb.WriteString(edgeStyle.Render(vertChar))
			}
		}
		if i < maxCol {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

func topoOrder(repo git.RepoState) []string {
	tips := []string{}
	seen := map[string]bool{}
	for _, id := range repo.Branches {
		if !seen[id] {
			seen[id] = true
			tips = append(tips, id)
		}
	}
	for _, id := range repo.RemoteRefs {
		if !seen[id] {
			seen[id] = true
			tips = append(tips, id)
		}
	}
	sort.Strings(tips)

	var order []string
	visited := map[string]bool{}
	var visit func(id string)
	visit = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true
		c, ok := repo.Commits[id]
		if !ok {
			return
		}
		for _, p := range c.Parents {
			visit(p)
		}
		order = append(order, id)
	}
	for _, t := range tips {
		visit(t)
	}

	// reverse so newest is first
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return order
}
