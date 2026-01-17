# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**AzureStrike** - A terminal-based game in Go that simulates Azure security scenarios for training security professionals. Features narrative-driven missions with a hybrid interface (CLI commands + TUI status panels) and a scoring/achievement system.

## Technology Stack

- **Language**: Go
- **TUI Framework**: Bubble Tea + Lip Gloss (Charm ecosystem)
- **Scenario Definitions**: YAML files
- **Build**: Standard Go modules

## Build Commands

```bash
# Build the binary
go build -o azurestrike ./cmd/azurestrike

# Run in development
go run ./cmd/azurestrike

# Run specific scenario
./azurestrike --scenario storage-recon

# Run tests
go test ./...

# Run a single test
go test -run TestName ./internal/package/

# Run tests with coverage
go test -cover ./...
```

## Architecture

```
cmd/azurestrike/       Entry point
internal/
  game/                Game engine (state, objectives, scoring)
  azure/               Mocked Azure environment
    entra/             Users, groups, roles, tokens
    storage/           Storage accounts, blobs, SAS tokens
    compute/           VMs, networking, NSGs
    arm/               Mock ARM API responses
  cli/                 Fake CLI implementations (az, kubectl, curl)
  tui/                 Terminal UI (Bubble Tea components)
  scenario/            YAML scenario loader
scenarios/             YAML scenario definitions
```

## Key Design Patterns

- **Bubble Tea Model-View-Update**: All TUI components follow the Bubble Tea pattern with `Init()`, `Update()`, and `View()` methods
- **Scenario-driven**: Game content defined in YAML, not hardcoded
- **Mock responses**: CLI commands return realistic Azure-formatted responses from the mocked environment
- **Objective triggers**: Player actions are tracked and matched against scenario objectives

## Scenario YAML Structure

```yaml
id: scenario-id
name: "Display Name"
difficulty: beginner|intermediate|advanced
briefing: |
  Narrative introduction...
resources:
  storage_accounts: [...]
  users: [...]
objectives:
  - id: objective_id
    description: "What the player must do"
    trigger: "action that completes this"
    points: 100
debrief: |
  Educational summary after completion...
```

## Mocked Azure Commands

The game intercepts and mocks these command patterns:
- `az storage blob list/download`
- `az ad user/sp list`
- `az vm list`
- `az network nsg rule list`
- `curl http://169.254.169.254/metadata/...` (IMDS)

## SPEC.md Reference

The `docs/project-management/SPEC.md` file is the authoritative technical and design specification for AzureStrike. Consult it for:
- Game design (objectives, scoring, hints, achievements)
- UI layout and component architecture
- Command system and routing
- Scenario YAML structure details
- Extension points for adding new features

## PLAN.md Guidelines

The `docs/project-management/PLAN.md` file tracks project progress like a product/sprint document. **Update it after every change.**

### Required Updates
- After completing ANY task, add it to the Current Phase's completed items
- When starting work, update "Current Work" section
- When fixing bugs, add to "Bug Fixes" in the current phase
- When adding features, add to "Features" in the current phase

### Document Structure
- **Current Work**: What's actively being worked on
- **Current Phase**: Active phase with Features, Bug Fixes, Improvements sections
- **Backlog**: Planned work organized by future phases
- **Completed Phases**: Historical record of past phases

### Rules
- Do NOT add documentation (build commands, usage examples, how-to sections)
- Keep entries concise (one line per item)
- Use checkboxes for trackable items
