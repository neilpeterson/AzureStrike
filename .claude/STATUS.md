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
- [x] First scenario: Storage Breach
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

### Improvements
- [x] Reorganize STATUS.md as phase-based document
- [x] Refactor az ad handlers to use modular entra environment
- [x] Refactor az vm handlers to use modular compute environment
- [x] Create internal/azure/entra package for Users, ServicePrincipals, Groups
- [x] Create internal/azure/compute package for VMs, NSGs
- [x] Update scenario loader to support all resource types
- [x] Simplify status panel - show summary + quick commands instead of full objectives
- [x] Add `o` shortcut for objectives command (displays full objectives in terminal)

---

## Backlog

### Phase 2: Core Scenarios
- [ ] Scenario 02: IMDS Token Theft
- [ ] Scenario 03: Service Principal Exposure
- [ ] Scenario 04: NSG Misconfiguration
- [ ] Scenario 05: Privilege Escalation via RBAC
- [ ] Save/load game state

### Phase 3: Advanced Scenarios
- [ ] Scenario 06: Key Vault Secrets Exfiltration
- [ ] Scenario 07: Managed Identity Abuse
- [ ] Scenario 08: Cross-Tenant Access
- [ ] Scenario 09: Azure Function Code Injection
- [ ] Scenario 10: Full Kill Chain (multi-stage)

### Phase 4: Polish
- [ ] Graph API mock endpoints
- [ ] Scenario editor/validator tool
- [ ] Leaderboard system
- [ ] Custom scenario loading from URL

---

## Scenario Tracker

| ID | Name | Status | Difficulty |
|----|------|--------|------------|
| 01 | Storage Misconfiguration Discovery | Complete | Beginner |
| 02 | IMDS Token Theft | Planned | Beginner |
| 03 | Service Principal Exposure | Planned | Beginner |
| 04 | NSG Misconfiguration | Planned | Intermediate |
| 05 | Privilege Escalation via RBAC | Planned | Intermediate |
| 06 | Key Vault Secrets Exfiltration | Planned | Intermediate |
| 07 | Managed Identity Abuse | Planned | Intermediate |
| 08 | Cross-Tenant Access | Planned | Advanced |
| 09 | Azure Function Code Injection | Planned | Advanced |
| 10 | Full Kill Chain | Planned | Advanced |
