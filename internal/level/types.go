package level

import "github.com/clockworkpc/git-dungeon/internal/git"

type CommitDef struct {
	Parents []string `yaml:"parents"`
	Message string   `yaml:"message"`
}

type WorkingFileDef struct {
	Name   string `yaml:"name"`
	Status string `yaml:"status"`
}

type InitialState struct {
	Branches    map[string]string     `yaml:"branches"`
	RemoteRefs  map[string]string     `yaml:"remote_refs"`
	Head        string                `yaml:"head"`
	Commits     map[string]CommitDef  `yaml:"commits"`
	WorkingTree []WorkingFileDef      `yaml:"working_tree"`
	StageArea   []WorkingFileDef      `yaml:"stage_area"`
}

type GoalDef struct {
	Branch                    string `yaml:"branch"`
	HeadAt                    string `yaml:"head_at"`
	Shape                     string `yaml:"shape"`
	MustNotIncludeMergeCommit bool   `yaml:"must_not_include_merge_commit"`
	RemoteMustMatchLocal      bool   `yaml:"remote_must_match_local"`
	CommitMessageContains     string `yaml:"commit_message_contains"`
	MinNewCommits             int    `yaml:"min_new_commits"`
	MinStagedFiles            int    `yaml:"min_staged_files"`
	CommandMustBeRun          string `yaml:"command_must_be_run"`
}

type LevelDef struct {
	ID              string      `yaml:"id"`
	Title           string      `yaml:"title"`
	Region          string      `yaml:"region"`
	Objective       string      `yaml:"objective"`
	Initial         InitialState `yaml:"initial"`
	Goal            GoalDef     `yaml:"goal"`
	AllowedCommands []string    `yaml:"allowed_commands"`
	Hints           []string    `yaml:"hints"`
	Explanation     string      `yaml:"explanation"`
}

func ToRepoState(init InitialState) git.RepoState {
	repo := git.RepoState{
		Branches:   make(map[string]string),
		RemoteRefs: make(map[string]string),
		Commits:    make(map[string]git.Commit),
		StageArea:  make(map[string]git.FileStatus),
		WorkingTree: git.WorkingTreeState{
			Files: make(map[string]git.FileStatus),
		},
	}

	for k, v := range init.Branches {
		repo.Branches[k] = v
	}
	for k, v := range init.RemoteRefs {
		repo.RemoteRefs[k] = v
	}

	for id, def := range init.Commits {
		parents := make([]string, len(def.Parents))
		copy(parents, def.Parents)
		repo.Commits[id] = git.Commit{
			ID:      id,
			Parents: parents,
			Message: def.Message,
			IsMerge: len(def.Parents) == 2,
		}
	}

	if len(init.Head) > 10 && init.Head[:9] == "detached:" {
		repo.HeadCommit = init.Head[9:]
	} else {
		repo.HeadBranch = init.Head
	}

	for _, f := range init.WorkingTree {
		repo.WorkingTree.Files[f.Name] = git.FileStatus(f.Status)
	}
	for _, f := range init.StageArea {
		repo.StageArea[f.Name] = git.FileStatus(f.Status)
	}

	repo.WorkingTree.Clean = len(repo.WorkingTree.Files) == 0 && len(repo.StageArea) == 0

	return repo
}
