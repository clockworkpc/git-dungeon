package git

import (
	"fmt"
	"strings"

	"github.com/clockworkpc/git-dungeon/internal/parser"
)

func simulateAdd(repo RepoState, cmd parser.ParsedCommand) (RepoState, string, error) {
	pathspec := "."
	if len(cmd.Args) > 0 {
		pathspec = cmd.Args[0]
	}

	var added []string
	for f, status := range repo.WorkingTree.Files {
		if status == FileModified || status == FileUntracked || status == FileDeleted {
			if pathspec == "." || pathspec == f {
				repo.StageArea[f] = FileStaged
				delete(repo.WorkingTree.Files, f)
				added = append(added, f)
			}
		}
	}

	if len(added) == 0 {
		return repo, "", fmt.Errorf("nothing added (no matching files for '%s')", pathspec)
	}

	return repo, fmt.Sprintf("staged: %s", strings.Join(added, ", ")), nil
}
