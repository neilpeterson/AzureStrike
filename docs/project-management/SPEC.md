# AzureStrike - Project Specification

This document is the authoritative technical and design reference for AzureStrike.

## 1. Vision & Goals

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

---

## 2. Game Design

### 2.1 Game Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Briefing   │ ──▶ │  Gameplay   │ ──▶ │  Fireworks  │ ──▶ │   Debrief   │
│  (Modal)    │     │  (CLI+TUI)  │     │  (Overlay)  │     │   (Modal)   │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
     Enter              Commands            ~4.5 sec            Enter
```

1. **Briefing** - Narrative introduction with mission context
2. **Gameplay** - Execute commands to complete objectives
3. **Fireworks** - Celebration animation on completion
4. **Debrief** - Educational summary with remediation guidance

### 2.2 Objective System

Objectives are goals the player must achieve, tracked by command execution.

**Trigger Matching** (checked after every command):

| Type | Syntax | Example |
|------|--------|---------|
| Substring | Plain text | `"blob list"` matches any command containing these words |
| Regex | `regex:` prefix | `"regex:curl.*contoso.*blob"` for complex patterns |
| Exact | Full command | `"az storage blob download --name secrets.txt"` |

**Completion Rules:**
- Each objective completes only once (duplicate commands ignored)
- Multiple objectives can complete from a single command
- Help commands (`-h`, `--help`) do not trigger objectives
- Hidden objectives are revealed when completed

**Ordering:**
- `order: 0` - Can be completed at any time
- `order: N` - Suggests sequence (not enforced, for display purposes)

### 2.3 Scoring System

**Point Mechanics:**

| Source | Calculation |
|--------|-------------|
| Objective completion | Base points from YAML |
| Hint penalty | Subtract cumulative cost of hints used |
| Completion bonus | +100 points when all objectives complete |

**Invariants:**
- Points never go below 0 (clamped)
- Points never exceed maximum possible
- Point history logged with timestamp and reason

**Example:**
```
Objective "download_secrets": 150 points
Hint Level 1 used: -10 points
Hint Level 2 used: -25 points
Final award: 150 - 10 - 25 = 115 points
```

### 2.4 Hint System

Progressive hints available per objective at increasing point costs.

| Level | Purpose | Typical Cost |
|-------|---------|--------------|
| 1 | Vague direction | 10-15 pts |
| 2 | Moderate guidance | 25-35 pts |
| 3 | Near-solution | 50+ pts |

**Mechanics:**
- `hint <objective_id>` shows next unused hint level
- Hint costs accumulate (level 1 + level 2 = both costs deducted)
- Hints tracked per objective (using level 2 doesn't require level 1)

### 2.5 Achievement System

Achievements are optional milestones beyond core objectives.

**Predefined Achievements:**

| ID | Name | Trigger |
|----|------|---------|
| `first_blood` | First Blood | Complete first objective |
| `speed_demon` | Speed Demon | Finish scenario in < 5 minutes |
| `no_hints` | Solo Operator | Complete without using hints |
| `perfect_score` | Perfect Score | Achieve maximum points |
| `blob_hunter` | Blob Hunter | Download sensitive storage data |
| `token_thief` | Token Thief | Extract IMDS credentials |
| `lateral_mover` | Lateral Mover | Move between resources using credentials |

**Unlock Mechanics:**
- Triggered explicitly by game logic or CLI handlers
- Not automatic from objectives (enables complex conditions)
- Unlock attempts on already-unlocked achievements ignored

---

## 3. User Interface Design

### 3.1 Layout Architecture

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

**Panel Dimensions:**
- Terminal: `width - 37` (remaining space after status panel)
- Status: 32 characters fixed width
- Gap/borders: 5 characters

### 3.2 Terminal Panel

**Features:**
- Command input with `$ ` prompt
- Scrollable output viewport (command history + results)
- Command history navigation (up/down arrows)

**Keyboard Controls:**

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

### 3.3 Status Panel

**Sections (top to bottom):**

1. **Header** - "=== STATUS ===" title
2. **Mission** - Scenario name + elapsed time
3. **Score** - Points with visual progress bar
4. **Objectives** - Completion count + hints used
5. **Commands** - Quick reference (objective, hint, score, help)

**Progress Bar:**
```
████████████░░░░░░░░░░░░ 52%
 (filled)    (empty)
```

### 3.4 Modal Screens

**Briefing (shown first):**
- Magenta double border
- Scenario name and difficulty
- Full narrative text
- Dismiss: Enter or Space

**Debrief (shown after completion):**
- Green double border
- "MISSION COMPLETE" header
- Final score in yellow
- Educational content
- Dismiss: Enter or Space

### 3.5 Fireworks Animation

Triggered on scenario completion before debrief.

**Parameters:**
- Duration: ~4.5 seconds (135 frames @ 30fps)
- Initial burst: 8 explosions
- Continuous spawning: 1-2 explosions every 8 frames
- Particles per explosion: 25-44
- Particle life: 20-44 frames

**Physics:**
- Gravity: +0.06 velocity per frame
- Horizontal spread: 3x base velocity
- Vertical spread: 1.8x base velocity
- Fade effect: dim when life < 8 frames

---

## 4. Architecture

### 4.1 Component Diagram

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
│    Parser     │ ─────────────────────────────────────────────────────┐
│  (CLI Router) │                                                      │
└───────────────┘                                                      │
        │                                                              │
        ├──────────────┬──────────────┬──────────────┐                 │
        ▼              ▼              ▼              ▼                 │
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  az Handler │  │curl Handler │  │scan Handler │  │ cat Handler │     │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘     │
        │              │              │              │                 │
        └──────────────┴──────────────┴──────────────┘                 │
                              │                                        │
                              ▼                                        │
                    ┌───────────────┐                                  │
                    │  Game State   │ ◀────────────────────────────────┘
                    └───────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│    Storage    │   │    Entra      │   │    Compute    │
│  Environment  │   │  Environment  │   │  Environment  │
└───────────────┘   └───────────────┘   └───────────────┘
```

### 4.2 Data Flow

```
User Input
    │
    ▼
Terminal.Update() ─────▶ CommandMsg
                              │
                              ▼
                    App.Update() ─────▶ Parser.Execute()
                                              │
                                              ▼
                                      Command Handler
                                              │
                                              ▼
                                      GameState.RecordCommand()
                                              │
                                              ▼
                                      Check Objectives ─────▶ Complete if matched
                                              │
                                              ▼
                                      Return Result (output + notifications)
                                              │
                                              ▼
                                      Terminal.AddOutput()
                                              │
                                              ▼
                                      Status.Update()
```

### 4.3 Game State

**Core Structure:**
```go
type State struct {
    Scenario           *Scenario           // Loaded scenario definition
    StartTime          time.Time           // Session start
    CompletedObjectives map[string]time.Time // ID → completion time
    CommandHistory     []CommandRecord     // Full command log
    HintsUsed          map[string]int      // Objective ID → highest level used
    Score              Score               // Points + history + achievements
    Status             GameStatus          // Playing, Completed, Failed

    // Mocked environments (immutable during gameplay)
    StorageEnv         *storage.Environment
    EntraEnv           *entra.Environment
    ComputeEnv         *compute.Environment
}
```

**Invariants:**
- Environments are immutable during gameplay
- Command history only appends (no undo)
- Objectives complete once (checked via map)
- State changes only through defined methods

### 4.4 Bubble Tea Pattern

All TUI components implement the Model interface:

```go
type Model interface {
    Init() tea.Cmd           // Initialize subscriptions
    Update(tea.Msg) (Model, tea.Cmd)  // Handle messages
    View() string            // Render to string
}
```

**Message Types:**
- `tea.KeyMsg` - Keyboard input
- `tea.WindowSizeMsg` - Terminal resize
- `CommandMsg` - Command execution request
- `FireworksTickMsg` - Animation frame

---

## 5. Command System

### 5.1 Command Routing

```
Input: "az storage blob list --account-name contoso"
       │
       ▼
┌──────────────────────────────────────────────────┐
│                     Parser                       │
│  Split: ["az", "storage", "blob", "list", ...]   │
└──────────────────────────────────────────────────┘
       │
       ▼ (switch on args[0])
       │
       ├── "az"        → azHandler.Execute()
       ├── "curl"      → handleCurl()
       ├── "scan"      → handleScan()
       ├── "cat"       → handleCat()
       ├── "help"      → handleHelp()
       ├── "objective" → handleObjectives()
       ├── "score"     → handleScore()
       ├── "hint"      → handleHint()
       ├── "clear"     → handleClear()
       └── default     → "Unknown command"
```

### 5.2 Azure CLI Handler Hierarchy

```
az
├── storage
│   ├── account
│   │   ├── list     → List all storage accounts
│   │   └── show     → Show specific account (--name required)
│   ├── container
│   │   └── list     → List containers (--account-name required)
│   └── blob
│       ├── list     → List blobs (--account-name, --container-name required)
│       └── download → Download blob (--name required)
├── ad
│   ├── user
│   │   ├── list     → List Entra ID users
│   │   └── show     → Show user (--id or --upn-or-object-id required)
│   └── sp
│       ├── list     → List service principals
│       └── show     → Show SP (--id required)
├── vm
│   ├── list         → List VMs (-g for resource group filter)
│   └── show         → Show VM (--name required)
├── account
│   ├── list         → List subscriptions
│   └── show         → Show current subscription
└── login            → (Informational only in game context)
```

### 5.3 Output Formats

Controlled via `--output` / `-o` flag:

| Format | Description | Example |
|--------|-------------|---------|
| `json` | Pretty-printed JSON (default) | `[{"name": "value", ...}]` |
| `table` | ASCII table with headers | `Name    Location\n------  --------\ncontoso eastus` |
| `tsv` | Tab-separated values | `contoso\teastus` |
| `none` | No output | (empty) |

### 5.4 Error Messages

Errors match real Azure CLI format:

```
az storage blob list: --account-name is required
```

```
This request is not authorized to perform this operation.
RequestId:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Time:2024-01-15T10:30:00.0000000Z
```

### 5.5 Special Commands

**curl (IMDS):**
```bash
curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/"
```

**curl (Blob Storage):**
```bash
curl https://contoso.blob.core.windows.net/backups?comp=list    # List blobs
curl https://contoso.blob.core.windows.net/backups/secrets.txt  # Download blob
```

**scan (Storage Enumeration):**
```bash
scan storage --company contoso              # Company-based wordlist
scan storage --wordlist common              # Common names wordlist
scan storage --accounts name1,name2,name3   # Specific accounts
```

**cat (Blob Content):**
```bash
cat backups/secrets.txt                     # Short form (auto-discover account)
cat contoso2024/backups/secrets.txt         # Full form (explicit account)
```

---

## 6. Scenario System

### 6.1 YAML Structure

```yaml
id: storage-breach
name: "Storage Misconfiguration Discovery"
difficulty: beginner
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

### 6.2 Resource Types

| Resource | YAML Key | Mock Handler |
|----------|----------|--------------|
| Storage Accounts | `storage_accounts` | `az storage account/container/blob` |
| Users | `users` | `az ad user` |
| Service Principals | `service_principals` | `az ad sp` |
| Groups | `groups` | `az ad group` |
| Virtual Machines | `virtual_machines` | `az vm` |
| Network Security Groups | `network_security_groups` | `az network nsg` |

### 6.3 Access Control Simulation

**Container Public Access Levels:**

| Level | `blob list` | `blob download` |
|-------|-------------|-----------------|
| `container` | Allowed | Allowed |
| `blob` | Denied (403) | Allowed (if URL known) |
| `none` | Denied (403) | Denied (403) |

### 6.4 Scenario Loading

**Search Order:**
1. `scenarios/<id>/scenario.yaml`
2. `scenarios/*-<id>/scenario.yaml` (prefix wildcard)
3. `scenarios/<id>.yaml` (flat file)
4. Scan all directories, match by `id` field

**Validation:**
- Required: id, name, at least 1 objective
- Objectives require: id, description, trigger, points
- Duplicate objective IDs are rejected

---

## 7. Extension Points

### 7.1 Adding New Commands

1. Add case in `Parser.Execute()` switch
2. Implement handler returning `Result{Output, Success}`
3. Parser auto-records and checks objectives

### 7.2 Adding New Azure Resources

1. Define struct in `internal/azure/<service>/`
2. Add YAML fields to scenario schema
3. Create Environment with lookup methods
4. Add CLI handlers in `internal/cli/az/`
5. Wire environment to game state

### 7.3 Adding New Scenarios

1. Create `scenarios/<id>/scenario.yaml`
2. Define resources, objectives, hints
3. Test with `--scenario <id>`
4. Document in `docs/solution-docs/scenarios/`

### 7.4 Adding New Achievements

1. Add achievement ID to predefined list in game engine
2. Add unlock logic in appropriate handler
3. Achievement displays in `score` command output

---

## 8. Technical Constraints

### 8.1 Session Model

- Single game session per app instance
- No save/load (planned for Phase 2)
- State exists only in memory
- Exit terminates all progress

### 8.2 Terminal Requirements

- Minimum terminal size: 80x24 (recommended: 120x40)
- Alt-screen mode (full terminal takeover)
- 256-color support for styling
- Mouse support for scrolling (optional)

### 8.3 Platform Support

- Primary: macOS, Linux
- Secondary: Windows (via Windows Terminal)
- Go 1.21+ required

---

## Appendix A: Color Scheme

| Element | Color Code | Usage |
|---------|------------|-------|
| Title | Magenta (205) | Panel headers |
| Section headers | Cyan (75) | "MISSION", "SCORE", etc. |
| Labels | Gray (241) | Dim descriptive text |
| Values | White (252) | Data values |
| Progress filled | Green (34) | Progress bar fill |
| Progress empty | Dark gray (236) | Progress bar background |
| Borders | Cyan (63) | Panel borders |
| Errors | Red | Error messages |
| Success | Green (34) | Completion notifications |

---

## Appendix B: Message Format Reference

**Objective Completion:**
```
[+] OBJECTIVE COMPLETE: <description> (+<points>)
```

**Hint Display:**
```
HINT (Level <N>, -<cost> pts):
<hint text>
```

**Score Display:**
```
=== SCORE ===

Total Points: <current>
Max Possible: <max>

Bonuses:
  [+] <bonus description>: +<points>

Achievements:
  [*] <achievement name>
```
