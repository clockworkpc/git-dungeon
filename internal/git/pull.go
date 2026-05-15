package git

import (
	"fmt"

	"github.com/clockworkpc/git-dungeon/internal/parser"
)

func simulatePull(repo RepoState, cmd parser.ParsedCommand) (RepoState, string, error) {
	branch := repo.HeadBranch
	if branch == "" {
		return repo, "", fmt.Errorf("error: HEAD detached")
	}
	if len(cmd.Args) >= 2 {
		branch = cmd.Args[1]
	}

	remoteKey := "origin/" + branch
	remoteTip, ok := repo.RemoteRefs[remoteKey]
	if !ok || remoteTip == "" {
		return repo, "Already up to date.", nil
	}

	localTip := repo.Branches[branch]

	if localTip == remoteTip {
		return repo, "Already up to date.", nil
	}

	if isAncestor(repo.Commits, localTip, remoteTip) {
		repo.Branches[branch] = remoteTip
		return repo, fmt.Sprintf("Fast-forward\n %s -> %s", branch, remoteTip), nil
	}

	if isAncestor(repo.Commits, remoteTip, localTip) {
		return repo, "Already up to date.", nil
	}

	mergeBase := lca(repo.Commits, localTip, remoteTip)
	id := generateID()
	repo.Commits[id] = Commit{
		ID:      id,
		Parents: []string{localTip, remoteTip},
		Message: fmt.Sprintf("Merge branch 'origin/%s'", branch),
		IsMerge: true,
	}
	repo.Branches[branch] = id

	return repo, fmt.Sprintf("Merge made by the 'ort' strategy.\nMerge base: %s", mergeBase), nil
}
