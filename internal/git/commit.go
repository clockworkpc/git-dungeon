package git

import (
	"fmt"

	"github.com/clockworkpc/git-dungeon/internal/parser"
)

func simulateCommit(repo RepoState, cmd parser.ParsedCommand) (RepoState, string, error) {
	if repo.HeadBranch == "" {
		return repo, "", fmt.Errorf("HEAD detached — checkout a branch first")
	}
	msg, ok := cmd.Flags["-m"]
	if !ok || msg == "" {
		return repo, "", fmt.Errorf("error: switch `-m' requires a value")
	}
	if len(repo.StageArea) == 0 {
		return repo, "", fmt.Errorf("nothing to commit, working tree clean")
	}

	id := generateID()
	parent := repo.HEAD()
	var parents []string
	if parent != "" {
		parents = []string{parent}
	}
	repo.Commits[id] = Commit{ID: id, Parents: parents, Message: msg}
	repo.Branches[repo.HeadBranch] = id
	repo.StageArea = make(map[string]FileStatus)

	return repo, fmt.Sprintf("[%s %s] %s", repo.HeadBranch, id, msg), nil
}
