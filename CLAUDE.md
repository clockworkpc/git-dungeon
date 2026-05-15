# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Git Dungeon** is a terminal-based puzzle game that teaches Git by making players manipulate commit graphs. Built in Go using Bubble Tea (TUI framework) and Lip Gloss (styling).

## Tech Stack

- **Go** — primary language
- **Bubble Tea** — TUI application framework (Elm-inspired model/update/view)
- **Lip Gloss** — terminal styling
- **Bubbles** — optional component library

## Build & Run

```bash
go build ./...
go run .
go test ./...          # run all tests
go test ./pkg/foo/...  # run tests in a specific package
```

## Architecture

The app follows the Bubble Tea application model with three major subsystems:

### TUI Layer (`model`, `screen`)
- `Model` holds all application state: current screen, active level, repo state, command input, history, message
- Screens: Map, Level, Graph visualization, Command input, Results view

### Game Engine
- **Level loader** — reads YAML level definitions from the content system
- **Command parser** — parses player-entered Git commands
- **Git simulation engine** — applies commands to `RepoState` (no real Git, pure simulation)
- **Goal validator** — checks whether current `RepoState` satisfies the level's win condition

### Content System
- Levels are defined in **YAML** — data-driven, not hardcoded
- Each level specifies: `id`, `title`, `objective`, `initial` repo state (branches + commits as a DAG), `goal` shape/constraints, `allowed_commands`
- Commits are identified by short labels (A, B, C…) with `parents` arrays forming the graph

## Core Data Structures

```go
type RepoState struct {
    Branches    map[string]string   // branch name → commit label
    RemoteRefs  map[string]string
    Commits     map[string]Commit
    HeadBranch  string
    WorkingTree WorkingTreeState
    Conflicts   []Conflict
}
```

## Key Design Constraints

- **No real Git execution** — the simulation engine must replicate Git graph semantics without shelling out or touching real repos
- **No arbitrary shell execution** — the command input accepts only recognized Git commands defined in `allowed_commands` per level
- **Data-driven content** — all levels, hints, and explanations live in YAML, not Go code
- **Failure is educational** — incorrect commands should produce a meaningful (but wrong) repo state, not just an error message

## MVP Scope

10 levels: `git status/add/commit` → `push/pull/diverged branches/fast-forward merge/three-way merge` → `rebase onto main/cherry-pick`. Everything else (multiplayer, web UI, AI tutoring, real filesystem, procedural generation) is explicitly out of scope.

## World Progression

Regions map to Git skills: Village (commits) → Harbor (push/pull) → Bridge (merge) → Time Tower (rebase) → Orchard (cherry-pick) → Conflict Mines (conflict resolution) → Release Castle (production repair).
