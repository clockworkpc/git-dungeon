# Git Dungeon — Founding Product & Game Design Brief

## Vision

**Git Dungeon** is a terminal-based puzzle game that teaches Git through increasingly complex repository-state challenges.

The core educational insight:

> Players learn Git by manipulating commit graphs rather than memorizing commands.

The game presents realistic Git workflows in a safe, simulated environment where users can experiment, fail, and develop intuition around branching, rebasing, merging, cherry-picking, and remote synchronization.

---

# Product Thesis

Git is difficult because repository state is largely invisible.

Most tutorials teach commands linearly:

```bash
git pull
git merge
git rebase
```

But Git is fundamentally a graph manipulation system.

Git Dungeon visualizes repository state continuously and turns Git workflows into spatial, interactive puzzles.

---

# Target Audience

## Primary Audience

- Junior developers
- Bootcamp graduates
- Self-taught programmers
- Engineers uncomfortable with advanced Git workflows

## Secondary Audience

- Engineering managers onboarding new hires
- Coding schools
- Internal developer education teams
- Developer tooling enthusiasts

---

# Design Pillars

These principles guide all product and engineering decisions.

## 1. Graph-First Learning

Every Git operation visibly transforms the commit graph.

The graph is the primary educational surface.

## 2. Real Commands, Safe Environment

Players use authentic Git commands and workflows.

The game simulates Git behavior without risk of damaging real repositories.

## 3. Progressive Complexity

The player journey progresses from:

```text
status → commit → push → pull → merge → rebase → cherry-pick → conflict resolution
```

Each mechanic builds on previously learned graph concepts.

## 4. Terminal-Native Experience

The game should feel like a real developer tool.

The UI should resemble authentic terminal workflows rather than gamified web tutorials.

## 5. Failure Is Educational

Incorrect commands should create meaningful repository states.

The player learns by observing consequences rather than simply receiving “wrong answer” messages.

---

# Core Gameplay Loop

```text
Read objective
Inspect repository graph
Enter Git command
Observe graph transformation
Receive feedback
Advance or retry
```

---

# MVP Scope

The MVP should focus on the smallest fully playable educational experience.

## Included Levels

### Beginner

1. `git status`
2. `git add`
3. `git commit`

### Intermediate

4. `git push`
5. `git pull`
6. Diverged branches
7. Fast-forward merge
8. Three-way merge

### Advanced Intro

9. Rebase onto main
10. Cherry-pick a single commit

---

# Explicit Non-Goals

To prevent scope explosion, the following are excluded from the MVP.

## Excluded Features

- Real Git repository mutation
- Multiplayer
- Online accounts
- Cloud sync
- Web UI
- Arbitrary shell execution
- Full Git implementation
- AI tutoring
- Real filesystem integration
- Advanced merge conflict editors
- Procedural level generation

---

# Game Structure

## World Progression

| Region | Core Skill |
|---|---|
| Village | Commit fundamentals |
| Harbor | Push / pull / remotes |
| Bridge | Merge workflows |
| Time Tower | Rebase |
| Orchard | Cherry-pick |
| Conflict Mines | Conflict resolution |
| Release Castle | Production repair workflows |

---

# Educational Philosophy

The game teaches:

## Mental Models

- Commits are immutable snapshots
- Branches are movable pointers
- Rebasing rewrites history
- Merging preserves branch topology
- Cherry-picking copies commits
- Remote repositories are independent graphs

## Not Memorization

The game should avoid rote command memorization.

Instead, players should develop the ability to reason about repository state transitions.

---

# Technical Architecture

## Stack

- Go
- Bubble Tea
- Lip Gloss
- Optional: Bubbles components

---

# High-Level Architecture

```text
Bubble Tea TUI
  ├── Screens
  │   ├── Map screen
  │   ├── Level screen
  │   ├── Graph visualization
  │   ├── Command input
  │   └── Results view
  │
  ├── Game engine
  │   ├── Level loader
  │   ├── Command parser
  │   ├── Git simulation engine
  │   └── Goal validator
  │
  └── Content system
      ├── YAML levels
      ├── Hints
      └── Explanations
```

---

# Bubble Tea Application Model

## Core Application State

```go
type Model struct {
    screen      Screen
    level       Level
    repo        RepoState
    command     string
    history     []CommandResult
    selected    int
    message     string
}
```

---

# Repository Simulation Engine

## Repository State

```go
type RepoState struct {
    Branches      map[string]string
    RemoteRefs    map[string]string
    Commits       map[string]Commit
    HeadBranch    string
    WorkingTree   WorkingTreeState
    Conflicts     []Conflict
}
```

---

# Level Definition Format

All gameplay content should be data-driven.

## Example Level

```yaml
id: rebase-001

title: Rebase Onto Main

objective: >
  Move feature/login onto the latest main
  without creating a merge commit.

initial:
  branches:
    main: C
    feature/login: E

  head: feature/login

  commits:
    A:
      parents: []

    B:
      parents: [A]

    C:
      parents: [B]

    D:
      parents: [B]

    E:
      parents: [D]

goal:
  branch: feature/login
  shape: linear_after_main
  must_not_include_merge_commit: true

allowed_commands:
  - git checkout
  - git rebase
```

---

# Success Criteria

- Players understand merge vs rebase
- Players understand local vs remote branches
- Players can reason about commit graphs
- The game is fully playable from a terminal
