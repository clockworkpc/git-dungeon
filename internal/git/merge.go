package git

import (
	"fmt"

	"github.com/clockworkpc/git-dungeon/internal/parser"
)

func simulateMerge(repo RepoState, cmd parser.ParsedCommand) (RepoState, string, error) {
	if len(cmd.Args) == 0 {
		return repo, "", fmt.Errorf("error: missing branch argument for merge")
	}
	target := cmd.Args[0]

	targetTip, ok := repo.Branches[target]
	if !ok {
		if targetTip, ok = repo.RemoteRefs[target]; !ok {
			return repo, "", fmt.Errorf("error: branch '%s' not found", target)
		}
	}

	localTip := repo.HEAD()
	if localTip == targetTip {
		return repo, "Already up to date.", nil
	}

	if isAncestor(repo.Commits, localTip, targetTip) {
		repo.Branches[repo.HeadBranch] = targetTip
		return repo, fmt.Sprintf("Fast-forward\n(was %s, now at %s)", localTip, targetTip), nil
	}

	if isAncestor(repo.Commits, targetTip, localTip) {
		return repo, "Already up to date.", nil
	}

	if len(repo.WorkingTree.Files) > 0 {
		var conflicts []Conflict
		for f := range repo.WorkingTree.Files {
			conflicts = append(conflicts, Conflict{File: f, Ours: localTip, Theirs: targetTip})
		}
		repo.Conflicts = conflicts
		return repo, "", fmt.Errorf("CONFLICT (content): Merge conflict in working tree files\nAutomatic merge failed; fix conflicts and then commit the result.")
	}

	mergeBase := lca(repo.Commits, localTip, targetTip)
	id := generateID()
	repo.Commits[id] = Commit{
		ID:      id,
		Parents: []string{localTip, targetTip},
		Message: fmt.Sprintf("Merge branch '%s'", target),
		IsMerge: true,
	}
	repo.Branches[repo.HeadBranch] = id

	return repo, fmt.Sprintf("Merge made by the 'ort' strategy.\n  (merge base: %s)", mergeBase), nil
}
