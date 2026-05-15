package git

import (
	"fmt"
	"strings"
)

func simulateStatus(repo RepoState) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("On branch %s\n", repo.CurrentBranch()))

	if len(repo.StageArea) > 0 {
		sb.WriteString("\nChanges to be committed:\n")
		for f := range repo.StageArea {
			sb.WriteString(fmt.Sprintf("  staged:   %s\n", f))
		}
	}

	var unstaged, untracked []string
	for f, s := range repo.WorkingTree.Files {
		switch s {
		case FileModified, FileDeleted:
			unstaged = append(unstaged, f)
		case FileUntracked:
			untracked = append(untracked, f)
		}
	}

	if len(unstaged) > 0 {
		sb.WriteString("\nChanges not staged for commit:\n")
		for _, f := range unstaged {
			sb.WriteString(fmt.Sprintf("  modified: %s\n", f))
		}
	}
	if len(untracked) > 0 {
		sb.WriteString("\nUntracked files:\n")
		for _, f := range untracked {
			sb.WriteString(fmt.Sprintf("  %s\n", f))
		}
	}

	if len(repo.StageArea) == 0 && len(unstaged) == 0 && len(untracked) == 0 {
		sb.WriteString("\nnothing to commit, working tree clean")
	}

	return sb.String()
}
