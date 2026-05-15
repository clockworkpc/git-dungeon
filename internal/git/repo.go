package git

import "fmt"

type FileStatus string

const (
	FileUntracked FileStatus = "untracked"
	FileModified  FileStatus = "modified"
	FileStaged    FileStatus = "staged"
	FileDeleted   FileStatus = "deleted"
)

type Commit struct {
	ID      string
	Parents []string
	Message string
	IsMerge bool
}

type WorkingTreeState struct {
	Files map[string]FileStatus
	Clean bool
}

type Conflict struct {
	File   string
	Ours   string
	Theirs string
}

type RepoState struct {
	Branches    map[string]string
	RemoteRefs  map[string]string
	Commits     map[string]Commit
	HeadBranch  string
	HeadCommit  string
	WorkingTree WorkingTreeState
	StageArea   map[string]FileStatus
	Conflicts   []Conflict
}

func (r *RepoState) HEAD() string {
	if r.HeadBranch != "" {
		return r.Branches[r.HeadBranch]
	}
	return r.HeadCommit
}

func (r RepoState) CurrentBranch() string {
	if r.HeadBranch != "" {
		return r.HeadBranch
	}
	return fmt.Sprintf("(HEAD detached at %s)", r.HeadCommit)
}

func (r RepoState) DeepCopy() RepoState {
	c := RepoState{
		HeadBranch: r.HeadBranch,
		HeadCommit: r.HeadCommit,
		Branches:   make(map[string]string, len(r.Branches)),
		RemoteRefs: make(map[string]string, len(r.RemoteRefs)),
		Commits:    make(map[string]Commit, len(r.Commits)),
		StageArea:  make(map[string]FileStatus, len(r.StageArea)),
		WorkingTree: WorkingTreeState{
			Files: make(map[string]FileStatus, len(r.WorkingTree.Files)),
			Clean: r.WorkingTree.Clean,
		},
		Conflicts: make([]Conflict, len(r.Conflicts)),
	}
	for k, v := range r.Branches {
		c.Branches[k] = v
	}
	for k, v := range r.RemoteRefs {
		c.RemoteRefs[k] = v
	}
	for k, v := range r.Commits {
		cp := Commit{ID: v.ID, Message: v.Message, IsMerge: v.IsMerge, Parents: make([]string, len(v.Parents))}
		copy(cp.Parents, v.Parents)
		c.Commits[k] = cp
	}
	for k, v := range r.StageArea {
		c.StageArea[k] = v
	}
	for k, v := range r.WorkingTree.Files {
		c.WorkingTree.Files[k] = v
	}
	copy(c.Conflicts, r.Conflicts)
	return c
}
