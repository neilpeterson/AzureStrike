# Project Specification

## Overview

**AzureStrike** is an educational terminal-based game that simulates Azure security reconnaissance scenarios. It trains security professionals to identify and exploit common Azure misconfigurations in a safe, simulated environment.

### Core Principles

1. **Realistic Simulation** - Mocked Azure CLI returns authentic responses matching real Azure output
2. **Scenario-Driven** - All game content defined in YAML, not hardcoded
3. **Educational Focus** - Narrative briefings, progressive hints, and post-game debriefs
4. **Gamified Learning** - Points, achievements, and objectives create engagement

### Target Audience

- Security professionals learning Azure attack techniques
- Blue team members understanding attacker perspectives
- Students studying cloud security concepts
- CTF participants practicing Azure scenarios

## Scope

### Included

- Terminal-based game with Bubble Tea TUI
- Mocked Azure CLI commands (az storage, az ad, az vm, az network, etc.)
- YAML-defined scenarios with objectives, hints, and scoring
- 10 planned scenarios covering beginner to advanced Azure security topics
- Achievement system and scoring mechanics
- Educational debriefs with remediation guidance

### Excluded

- Real Azure API calls (all responses are mocked)
- Multiplayer or networked gameplay
- Web-based interface
- Mobile support

## Requirements

### Functional Requirements

#### Game Flow
1. **Briefing** - Narrative introduction with mission context
2. **Gameplay** - Execute commands to complete objectives
3. **Fireworks** - Celebration animation on completion
4. **Debrief** - Educational summary with remediation guidance

#### Objective System
- Objectives tracked by command execution
- Trigger matching: substring, regex (`regex:` prefix), or exact match
- Each objective completes only once
- Hidden objectives revealed when completed
- Help commands (`-h`, `--help`) do not trigger objectives

#### Scoring System
- Points awarded for objective completion (defined in YAML)
- Hint penalties subtract from objective points
- Completion bonus (+100 points when all objectives complete)
- Points never go below 0

#### Hint System
- Progressive hints per objective at increasing point costs
- Level 1: Vague direction (10-15 pts)
- Level 2: Moderate guidance (25-35 pts)
- Level 3: Near-solution (50+ pts)

#### Achievement System
- Predefined achievements: First Blood, Speed Demon, Solo Operator, Perfect Score, etc.
- Triggered by game logic or CLI handlers
- Displayed in score command output

#### Command System
- `az` - Mocked Azure CLI (storage, ad, vm, account, network, keyvault, functionapp)
- `curl` - IMDS endpoint and blob storage URLs
- `scan` - Storage account enumeration with wordlists
- `cat` - Direct blob content viewing
- `help`, `objective`, `score`, `hint`, `clear` - Game commands

### Non-Functional Requirements

#### Terminal Requirements
- Minimum terminal size: 80x24 (recommended: 120x40)
- Alt-screen mode (full terminal takeover)
- 256-color support for styling
- Mouse support for scrolling (optional)

#### Platform Support
- Primary: macOS, Linux
- Secondary: Windows (via Windows Terminal)
- Go 1.21+ required

#### Performance
- Responsive command execution (< 100ms)
- Smooth fireworks animation (30fps)
- No network latency (all mocked locally)

### Constraints and Assumptions

- Single game session per app instance
- No save/load (planned for future phase)
- State exists only in memory
- Exit terminates all progress
- All Azure responses are pre-defined mocks

## Architecture

### Component Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                                  App                                    │
│                           (Bubble Tea Model)                            │
└─────────────────────────────────────────────────────────────────────────┘
        │                         │                         │
        ▼                         ▼                         ▼
┌───────────────┐       ┌───────────────┐       ┌───────────────┐
│   Terminal    │       │    Status     │       │   Fireworks   │
│    Panel      │       │    Panel      │       │    Overlay    │
└───────────────┘       └───────────────┘       └───────────────┘
        │
        ▼
┌───────────────┐
│    Parser     │
│  (CLI Router) │
└───────────────┘
        │
        ├──────────────┬──────────────┬──────────────┐
        ▼              ▼              ▼              ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  az Handler │  │curl Handler │  │scan Handler │  │ cat Handler │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
        │              │              │              │
        └──────────────┴──────────────┴──────────────┘
                              │
                              ▼
                    ┌───────────────┐
                    │  Game State   │
                    └───────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│    Storage    │   │    Entra      │   │    Compute    │
│  Environment  │   │  Environment  │   │  Environment  │
└───────────────┘   └───────────────┘   └───────────────┘
```

### Directory Structure

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

### Data Flow

```
User Input → Terminal.Update() → CommandMsg → App.Update() → Parser.Execute()
    → Command Handler → GameState.RecordCommand() → Check Objectives
    → Return Result → Terminal.AddOutput() → Status.Update()
```

### Bubble Tea Pattern

All TUI components implement the Model interface:
- `Init()` - Initialize subscriptions
- `Update(tea.Msg)` - Handle messages
- `View()` - Render to string

Message types: `tea.KeyMsg`, `tea.WindowSizeMsg`, `CommandMsg`, `FireworksTickMsg`

## Schema

### Scenario YAML Structure

```yaml
id: scenario-id
name: "Display Name"
difficulty: beginner|intermediate|advanced
briefing: |
  Multi-line narrative introduction...

debrief: |
  Educational summary with remediation...

resources:
  storage_accounts:
    - name: contosostorage2024
      resource_group: production-rg
      location: eastus
      public_access: true
      containers:
        - name: backups
          public_access: container
          blobs:
            - name: secrets.txt
              content: |
                Sensitive data here...
              content_type: text/plain
              size: 512

  users:
    - id: "aaaa-bbbb-cccc-dddd"
      display_name: "Admin User"
      user_principal_name: "admin@contoso.com"

  virtual_machines:
    - name: web-server-01
      resource_group: production-rg
      power_state: "VM running"
      public_ip: "20.x.x.x"

objectives:
  - id: discover_account
    description: "Discover the storage account"
    trigger: "az storage account"
    points: 50
    order: 1

  - id: download_secrets
    description: "Download sensitive files"
    trigger: "regex:curl.*secrets\\.txt"
    points: 150
    order: 4
    hidden: false

hints:
  - objective_id: discover_account
    level: 1
    text: "Try exploring 'az storage' commands"
    point_cost: 10
  - objective_id: discover_account
    level: 2
    text: "Use 'az storage account list'"
    point_cost: 25
```

### Resource Types

| Resource | YAML Key | Mock Handler |
|----------|----------|--------------|
| Storage Accounts | `storage_accounts` | `az storage account/container/blob` |
| Users | `users` | `az ad user` |
| Service Principals | `service_principals` | `az ad sp` |
| Groups | `groups` | `az ad group` |
| Virtual Machines | `virtual_machines` | `az vm` |
| Network Security Groups | `network_security_groups` | `az network nsg` |

### Game State Structure

```go
type State struct {
    Scenario            *Scenario
    StartTime           time.Time
    CompletedObjectives map[string]time.Time
    CommandHistory      []CommandRecord
    HintsUsed           map[string]int
    Score               Score
    Status              GameStatus
    StorageEnv          *storage.Environment
    EntraEnv            *entra.Environment
    ComputeEnv          *compute.Environment
}
```

## Implementation Details

### UI Layout

```
┌──────────────────────────────────────────┬────────────────────────────────┐
│                                          │        === STATUS ===          │
│              TERMINAL                    │                                │
│                                          │  MISSION                       │
│  $ az storage account list               │    Storage Breach              │
│  [                                       │    Time: 05:23                 │
│    {                                     │                                │
│      "name": "contoso2024",              │  SCORE                         │
│      ...                                 │    125 / 375                   │
│    }                                     │    ████████░░░░░░░░░░░ 33%     │
│  ]                                       │                                │
│                                          │  OBJECTIVES                    │
│  [+] OBJECTIVE COMPLETE: Discover...     │    2 / 4 complete              │
│                                          │                                │
│  $ _                                     │  COMMANDS                      │
│                                          │    objective  score            │
│                                          │    hint       help             │
└──────────────────────────────────────────┴────────────────────────────────┘
```

Panel dimensions: Terminal (width - 37), Status (32 chars fixed), Gap/borders (5 chars)

### Keyboard Controls

| Key | Action |
|-----|--------|
| Enter | Execute command |
| Up/Down | Navigate command history |
| PgUp/PgDown | Scroll output |
| Shift+Up/Down | Scroll output (alternative) |
| Home/End | Jump to top/bottom of output |
| Mouse wheel | Scroll output |
| Ctrl+L | Clear terminal |
| Ctrl+C | Quit game |

### Output Formats

Controlled via `--output` / `-o` flag:
- `json` - Pretty-printed JSON (default)
- `table` - ASCII table with headers
- `tsv` - Tab-separated values
- `none` - No output

### Fireworks Animation

- Duration: ~4.5 seconds (135 frames @ 30fps)
- Initial burst: 8 explosions
- Particles per explosion: 25-44
- Physics: Gravity +0.06 velocity/frame, fade when life < 8 frames

### Color Scheme

| Element | Color | Usage |
|---------|-------|-------|
| Title | Magenta (205) | Panel headers |
| Section headers | Cyan (75) | "MISSION", "SCORE", etc. |
| Labels | Gray (241) | Dim descriptive text |
| Values | White (252) | Data values |
| Progress filled | Green (34) | Progress bar fill |
| Borders | Cyan (63) | Panel borders |
| Errors | Red | Error messages |
| Success | Green (34) | Completion notifications |

## Dependencies

### Go Modules

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `gopkg.in/yaml.v3` - YAML parsing

### Build Requirements

- Go 1.21+
- Standard Go modules (`go mod tidy`)

### Development Tools

- `golangci-lint` - Linting
- `testify` - Test assertions
