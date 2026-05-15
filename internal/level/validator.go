package level

import (
	"strings"

	"github.com/clockworkpc/git-dungeon/internal/git"
)

type ValidationResult struct {
	Passed bool
	Reason string
}

func Validate(goal GoalDef, initial git.RepoState, current git.RepoState, lastCmdSub string) ValidationResult {
	switch goal.Shape {
	case "command_run":
		if lastCmdSub == goal.CommandMustBeRun[4:] { // strip "git "
			return ValidationResult{Passed: true}
		}
		return ValidationResult{Reason: "Run the required command."}

	case "working_tree_staged":
		if len(current.StageArea) >= goal.MinStagedFiles {
			return ValidationResult{Passed: true}
		}
		return ValidationResult{Reason: "Stage the required files."}

	case "has_new_commit":
		newCount := len(current.Commits) - len(initial.Commits)
		if newCount < goal.MinNewCommits {
			newCount = 1
		}
		if len(current.Commits) > len(initial.Commits) {
			tip := current.Branches[goal.Branch]
			if goal.CommitMessageContains != "" {
				c, ok := current.Commits[tip]
				if !ok || !strings.Contains(c.Message, goal.CommitMessageContains) {
					return ValidationResult{Reason: "Commit message must contain: " + goal.CommitMessageContains}
				}
			}
			_ = newCount
			return ValidationResult{Passed: true}
		}
		return ValidationResult{Reason: "Create a new commit."}

	case "remote_up_to_date":
		local := current.Branches[goal.Branch]
		remote := current.RemoteRefs["origin/"+goal.Branch]
		if local != "" && local == remote {
			return ValidationResult{Passed: true}
		}
		return ValidationResult{Reason: "Push your commits to the remote."}

	case "fast_forwarded":
		initialTip := initial.Branches[goal.Branch]
		currentTip := current.Branches[goal.Branch]
		if currentTip == "" || currentTip == initialTip {
			return ValidationResult{Reason: "The branch has not advanced."}
		}
		if goal.HeadAt != "" && currentTip != goal.HeadAt {
			return ValidationResult{Reason: "Branch tip is not at the expected commit."}
		}
		if goal.MustNotIncludeMergeCommit {
			if hasMergeCommitBetween(current.Commits, currentTip, initialTip) {
				return ValidationResult{Reason: "History must not contain a merge commit."}
			}
		}
		return ValidationResult{Passed: true}

	case "merged_into_main":
		targetTip, ok := initial.Branches[goal.Branch]
		if !ok {
			return ValidationResult{Reason: "Target branch not found."}
		}
		mainTip := current.Branches["master"]
		if isReachable(current.Commits, targetTip, mainTip) {
			return ValidationResult{Passed: true}
		}
		return ValidationResult{Reason: "Merge the branch into master."}

	case "linear_after_main":
		mainTip := current.Branches["master"]
		branchTip := current.Branches[goal.Branch]
		if branchTip == "" {
			return ValidationResult{Reason: "Branch not found."}
		}
		if !isReachable(current.Commits, mainTip, branchTip) {
			return ValidationResult{Reason: "Branch must be based on main."}
		}
		if goal.MustNotIncludeMergeCommit {
			if hasMergeCommitBetween(current.Commits, branchTip, mainTip) {
				return ValidationResult{Reason: "History must not contain a merge commit."}
			}
		}
		return ValidationResult{Passed: true}
	}

	return ValidationResult{Reason: "Unknown goal shape: " + goal.Shape}
}

func isReachable(commits map[string]git.Commit, target, from string) bool {
	visited := map[string]bool{}
	queue := []string{from}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if visited[cur] {
			continue
		}
		visited[cur] = true
		if cur == target {
			return true
		}
		if c, ok := commits[cur]; ok {
			queue = append(queue, c.Parents...)
		}
	}
	return false
}

func hasMergeCommitBetween(commits map[string]git.Commit, tip, base string) bool {
	visited := map[string]bool{}
	queue := []string{tip}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if visited[cur] || cur == base {
			continue
		}
		visited[cur] = true
		c, ok := commits[cur]
		if !ok {
			continue
		}
		if c.IsMerge {
			return true
		}
		queue = append(queue, c.Parents...)
	}
	return false
}
