package git

import (
	"fmt"

	"github.com/clockworkpc/git-dungeon/internal/parser"
)

func simulateRebase(repo RepoState, cmd parser.ParsedCommand) (RepoState, string, error) {
	if repo.HeadBranch == "" {
		return repo, "", fmt.Errorf("error: HEAD detached — checkout a branch first")
	}
	if len(cmd.Args) == 0 {
		return repo, "", fmt.Errorf("error: missing upstream argument for rebase")
	}
	upstream := cmd.Args[0]

	upstreamTip, ok := repo.Branches[upstream]
	if !ok {
		return repo, "", fmt.Errorf("error: invalid upstream '%s'", upstream)
	}

	branchTip := repo.Branches[repo.HeadBranch]
	if isAncestor(repo.Commits, branchTip, upstreamTip) {
		repo.Branches[repo.HeadBranch] = upstreamTip
		return repo, fmt.Sprintf("Current branch %s is up to date.", repo.HeadBranch), nil
	}

	base := lca(repo.Commits, branchTip, upstreamTip)
	toReplay := topoSort(repo.Commits, branchTip, base)

	for _, c := range toReplay {
		if c.IsMerge {
			return repo, "", fmt.Errorf("error: cannot rebase: commit %s is a merge commit", c.ID)
		}
	}

	prev := upstreamTip
	for _, c := range toReplay {
		newID := generateID()
		repo.Commits[newID] = Commit{
			ID:      newID,
			Parents: []string{prev},
			Message: c.Message,
		}
		prev = newID
	}

	repo.Branches[repo.HeadBranch] = prev
	return repo, fmt.Sprintf("Successfully rebased and updated refs/heads/%s.", repo.HeadBranch), nil
}
