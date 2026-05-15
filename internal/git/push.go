package git

import (
	"fmt"

	"github.com/clockworkpc/git-dungeon/internal/parser"
)

func simulatePush(repo RepoState, cmd parser.ParsedCommand) (RepoState, string, error) {
	branch := repo.HeadBranch
	if branch == "" {
		return repo, "", fmt.Errorf("error: HEAD detached — checkout a branch first")
	}
	if len(cmd.Args) >= 2 {
		branch = cmd.Args[1]
	}

	localTip, ok := repo.Branches[branch]
	if !ok {
		return repo, "", fmt.Errorf("error: branch '%s' not found", branch)
	}

	remoteKey := "origin/" + branch
	remoteTip, hasRemote := repo.RemoteRefs[remoteKey]

	if !hasRemote || remoteTip == "" {
		repo.RemoteRefs[remoteKey] = localTip
		return repo, fmt.Sprintf("Branch '%s' set up to track remote branch '%s' from 'origin'.", branch, branch), nil
	}

	if remoteTip == localTip {
		return repo, "Everything up-to-date", nil
	}

	if isAncestor(repo.Commits, remoteTip, localTip) {
		repo.RemoteRefs[remoteKey] = localTip
		return repo, fmt.Sprintf("   %s..%s  %s -> origin/%s", remoteTip[:min(len(remoteTip), 6)], localTip[:min(len(localTip), 6)], branch, branch), nil
	}

	if isAncestor(repo.Commits, localTip, remoteTip) {
		return repo, "", fmt.Errorf("! [rejected] %s -> %s (non-fast-forward)\nerror: Updates were rejected because the remote contains work that you do\nnot have locally.", branch, branch)
	}

	return repo, "", fmt.Errorf("! [rejected] %s -> %s (non-fast-forward)\nerror: Updates were rejected. The remote contains work you do not have locally.", branch, branch)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
