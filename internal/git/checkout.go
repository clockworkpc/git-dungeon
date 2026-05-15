package git

import (
	"fmt"

	"github.com/clockworkpc/git-dungeon/internal/parser"
)

func simulateCheckout(repo RepoState, cmd parser.ParsedCommand) (RepoState, string, error) {
	newBranch := cmd.BoolFlags["-b"]

	if len(cmd.Args) == 0 {
		return repo, "", fmt.Errorf("error: missing argument for checkout")
	}
	target := cmd.Args[0]

	if newBranch {
		if _, exists := repo.Branches[target]; exists {
			return repo, "", fmt.Errorf("fatal: A branch named '%s' already exists.", target)
		}
		repo.Branches[target] = repo.HEAD()
		repo.HeadBranch = target
		repo.HeadCommit = ""
		return repo, fmt.Sprintf("Switched to a new branch '%s'", target), nil
	}

	if _, ok := repo.Branches[target]; ok {
		repo.HeadBranch = target
		repo.HeadCommit = ""
		return repo, fmt.Sprintf("Switched to branch '%s'", target), nil
	}

	if _, ok := repo.Commits[target]; ok {
		repo.HeadBranch = ""
		repo.HeadCommit = target
		return repo, fmt.Sprintf("HEAD is now at %s", target), nil
	}

	return repo, "", fmt.Errorf("error: pathspec '%s' did not match any known refs", target)
}
