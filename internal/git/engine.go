package git

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/clockworkpc/git-dungeon/internal/parser"
)

func Apply(repo RepoState, cmd parser.ParsedCommand) (RepoState, string, error) {
	next := repo.DeepCopy()
	switch cmd.Sub {
	case "status":
		return next, simulateStatus(next), nil
	case "add":
		return simulateAdd(next, cmd)
	case "commit":
		return simulateCommit(next, cmd)
	case "push":
		return simulatePush(next, cmd)
	case "pull":
		return simulatePull(next, cmd)
	case "merge":
		return simulateMerge(next, cmd)
	case "rebase":
		return simulateRebase(next, cmd)
	case "cherry-pick":
		return simulateCherryPick(next, cmd)
	case "checkout":
		return simulateCheckout(next, cmd)
	default:
		return repo, "", fmt.Errorf("unknown git subcommand: %s", cmd.Sub)
	}
}

func generateID() string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func isAncestor(commits map[string]Commit, candidate, tip string) bool {
	if candidate == tip {
		return true
	}
	visited := map[string]bool{}
	queue := []string{tip}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if visited[cur] {
			continue
		}
		visited[cur] = true
		if cur == candidate {
			return true
		}
		if c, ok := commits[cur]; ok {
			queue = append(queue, c.Parents...)
		}
	}
	return false
}

func lca(commits map[string]Commit, a, b string) string {
	ancestors := map[string]bool{}
	queue := []string{a}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if ancestors[cur] {
			continue
		}
		ancestors[cur] = true
		if c, ok := commits[cur]; ok {
			queue = append(queue, c.Parents...)
		}
	}
	queue = []string{b}
	visited := map[string]bool{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if visited[cur] {
			continue
		}
		visited[cur] = true
		if ancestors[cur] {
			return cur
		}
		if c, ok := commits[cur]; ok {
			queue = append(queue, c.Parents...)
		}
	}
	return ""
}

// topoSort returns commits reachable from tip but not from base, oldest-first.
func topoSort(commits map[string]Commit, tip, base string) []Commit {
	reachableFromBase := map[string]bool{}
	if base != "" {
		q := []string{base}
		for len(q) > 0 {
			cur := q[0]
			q = q[1:]
			if reachableFromBase[cur] {
				continue
			}
			reachableFromBase[cur] = true
			if c, ok := commits[cur]; ok {
				q = append(q, c.Parents...)
			}
		}
	}

	var result []Commit
	visited := map[string]bool{}
	var dfs func(id string)
	dfs = func(id string) {
		if visited[id] || reachableFromBase[id] {
			return
		}
		visited[id] = true
		c, ok := commits[id]
		if !ok {
			return
		}
		for _, p := range c.Parents {
			dfs(p)
		}
		result = append(result, c)
	}
	dfs(tip)
	return result
}
