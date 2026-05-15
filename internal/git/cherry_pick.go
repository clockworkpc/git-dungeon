package git

import (
	"fmt"

	"github.com/clockworkpc/git-dungeon/internal/parser"
)

func simulateCherryPick(repo RepoState, cmd parser.ParsedCommand) (RepoState, string, error) {
	if repo.HeadBranch == "" {
		return repo, "", fmt.Errorf("error: HEAD detached — checkout a branch first")
	}
	if len(cmd.Args) == 0 {
		return repo, "", fmt.Errorf("error: missing commit argument for cherry-pick")
	}
	target := cmd.Args[0]

	src, ok := repo.Commits[target]
	if !ok {
		return repo, "", fmt.Errorf("error: bad revision '%s'", target)
	}

	id := generateID()
	repo.Commits[id] = Commit{
		ID:      id,
		Parents: []string{repo.HEAD()},
		Message: src.Message + fmt.Sprintf(" (cherry picked from %s)", target),
	}
	repo.Branches[repo.HeadBranch] = id

	return repo, fmt.Sprintf("[%s %s] %s", repo.HeadBranch, id, src.Message), nil
}
