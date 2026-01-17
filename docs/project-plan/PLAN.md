# AzureStrike - Project Status

## Current Work

None

---

## Phase 1: Foundation (Complete)

### Features
- [x] Project structure and Go modules
- [x] Bubble Tea TUI shell with terminal and status panels
- [x] Command parser framework
- [x] Mock az CLI (storage, ad, vm, account)
- [x] YAML scenario loader
- [x] Game state with objective tracking
- [x] Scoring system with achievements
- [x] Implement `--help` / `-h` for all az commands
- [x] Match real Azure CLI help format
- [x] Add Makefile with standard targets
- [x] Fireworks animation on objective completion
- [x] Secret commands: `fireworks` and `konami`
- [x] Modular architecture for pluggable scenarios
- [x] README.md with project overview and scenario roadmap

### Bug Fixes
- [x] Fix `--help` flag routing (`az storage -h` was showing main help)
- [x] Fix `clear` command breaking TUI (was sending raw ANSI codes)
- [x] Fix `exit`/`quit` commands to actually exit the game
- [x] Fix fireworks - now full-screen celebration overlay (cleaner approach)
- [x] Fix fireworks not showing on completion (debrief was rendered first)

### Improvements
- [x] Reorganize STATUS.md as phase-based document
- [x] Refactor az ad handlers to use modular entra environment
- [x] Refactor az vm handlers to use modular compute environment
- [x] Create internal/azure/entra package for Users, ServicePrincipals, Groups
- [x] Create internal/azure/compute package for VMs, NSGs
- [x] Update scenario loader to support all resource types
- [x] Simplify status panel - show summary + quick commands instead of full objectives
- [x] Add `o` shortcut for objectives command (displays full objectives in terminal)
- [x] Terminal scroll support (PgUp/PgDown, Shift+Up/Down, Home/End, Mouse wheel)
- [x] docs/WALKTHROUGHS.md with step-by-step commands for each scenario
- [x] Reorganize docs into project-plan/ and solution-docs/ directories
- [x] Create comprehensive SPECIFICATION.md (game design, architecture, systems)

---

## Phase 2: Storage Misconfiguration (Complete)

- [x] Storage account enumeration via az CLI
- [x] Container listing with public access detection
- [x] Blob download from public containers
- [x] curl support for Azure blob endpoints with realistic XML responses
- [x] scan command for account name enumeration with wordlists
- [x] cat command for direct blob content viewing
- [x] External reconnaissance approach (no pre-auth required)

---

## Backlog

### Phase 3: Testing & CI/CD

#### Unit Tests
- [ ] Test framework setup (Go testing + testify)
- [ ] CLI parser unit tests
- [ ] Azure handler unit tests (storage, entra, compute)
- [ ] Game state unit tests (objectives, scoring, achievements)
- [ ] Scenario loader unit tests
- [ ] TUI component tests (terminal, status panel)
- [ ] Code coverage reporting

#### GitHub Actions
- [ ] CI workflow (build, test, lint on push/PR)
- [ ] Go linting (golangci-lint)
- [ ] Code coverage badge
- [ ] Multi-platform build matrix (linux, darwin, windows)
- [ ] Automated PR checks

#### Release Management
- [ ] Semantic versioning setup
- [ ] Release workflow (tag-triggered)
- [ ] Cross-platform binary builds (GOOS/GOARCH)
- [ ] GitHub Releases with artifacts
- [ ] Changelog generation
- [ ] Homebrew formula (optional)

---

### Phase 4: IMDS Token Theft

- [ ] curl support for IMDS endpoint (169.254.169.254)
- [ ] Instance metadata retrieval
- [ ] Managed identity token extraction
- [ ] Token usage for Azure resource access

### Phase 5: Service Principal Exposure

- [ ] Environment/config file discovery (cat .env, config.json)
- [ ] Service principal credential extraction
- [ ] az login with service principal
- [ ] Enumerate SP permissions

### Phase 6: NSG Misconfiguration

- [ ] az network nsg commands implementation
- [ ] NSG rule listing and analysis
- [ ] Identify overly permissive rules (0.0.0.0/0, Any)
- [ ] VM to NSG association discovery

### Phase 7: Privilege Escalation via RBAC

- [ ] az role assignment commands
- [ ] az role definition commands
- [ ] Current user permission enumeration
- [ ] RBAC misconfiguration exploitation

### Phase 8: Key Vault Secrets Exfiltration

- [ ] az keyvault commands implementation
- [ ] Secret listing and retrieval
- [ ] Access policy analysis
- [ ] Credential extraction from vault

### Phase 9: Managed Identity Abuse

- [ ] System-assigned identity enumeration
- [ ] User-assigned identity discovery
- [ ] Cross-resource access via identity
- [ ] Lateral movement techniques

### Phase 10: Cross-Tenant Access

- [ ] Multi-tenant scenario support
- [ ] Cross-tenant access policy discovery
- [ ] B2B collaboration exploitation
- [ ] Tenant enumeration techniques

### Phase 11: Azure Function Injection

- [ ] az functionapp commands implementation
- [ ] Function app settings enumeration
- [ ] Code injection vulnerability
- [ ] Function execution exploitation

### Phase 12: Full Kill Chain

- [ ] Multi-stage attack scenario
- [ ] Combines techniques from previous scenarios
- [ ] Initial access to data exfiltration
- [ ] Persistence and lateral movement

### Phase 13: Polish

- [ ] Save/load game state
- [ ] Graph API mock endpoints
- [ ] Scenario editor/validator tool
- [ ] Leaderboard system
- [ ] Custom scenario loading from URL

---

## Scenario Tracker

| ID | Name | Phase | Status | Difficulty |
|----|------|-------|--------|------------|
| 01 | Storage Misconfiguration Discovery | 2 | Complete | Beginner |
| 02 | IMDS Token Theft | 4 | Planned | Beginner |
| 03 | Service Principal Exposure | 5 | Planned | Beginner |
| 04 | NSG Misconfiguration | 6 | Planned | Intermediate |
| 05 | Privilege Escalation via RBAC | 7 | Planned | Intermediate |
| 06 | Key Vault Secrets Exfiltration | 8 | Planned | Intermediate |
| 07 | Managed Identity Abuse | 9 | Planned | Intermediate |
| 08 | Cross-Tenant Access | 10 | Planned | Advanced |
| 09 | Azure Function Code Injection | 11 | Planned | Advanced |
| 10 | Full Kill Chain | 12 | Planned | Advanced |
