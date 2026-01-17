# Project Plan

## Status

Phase 3 complete. Testing framework and CI/CD pipelines are implemented.

## Phases

### Phase 1: Foundation (Complete)

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
- [x] Fix `--help` flag routing (`az storage -h` was showing main help)
- [x] Fix `clear` command breaking TUI (was sending raw ANSI codes)
- [x] Fix `exit`/`quit` commands to actually exit the game
- [x] Fix fireworks - now full-screen celebration overlay
- [x] Fix fireworks not showing on completion (debrief was rendered first)
- [x] Refactor az ad handlers to use modular entra environment
- [x] Refactor az vm handlers to use modular compute environment
- [x] Create internal/azure/entra package for Users, ServicePrincipals, Groups
- [x] Create internal/azure/compute package for VMs, NSGs
- [x] Update scenario loader to support all resource types
- [x] Simplify status panel - show summary + quick commands
- [x] Add `o` shortcut for objectives command
- [x] Terminal scroll support (PgUp/PgDown, Shift+Up/Down, Home/End, Mouse wheel)

### Phase 2: Storage Misconfiguration (Complete)

- [x] Storage account enumeration via az CLI
- [x] Container listing with public access detection
- [x] Blob download from public containers
- [x] curl support for Azure blob endpoints with realistic XML responses
- [x] scan command for account name enumeration with wordlists
- [x] cat command for direct blob content viewing
- [x] External reconnaissance approach (no pre-auth required)

### Phase 3: Testing & CI/CD (Complete)

- [x] Test framework setup (Go testing + testify)
- [x] CLI parser unit tests
- [x] Azure handler unit tests (storage, entra, compute)
- [x] Game state unit tests (objectives, scoring, achievements)
- [x] Scenario loader unit tests
- [x] Code coverage reporting
- [x] CI workflow (build, test, lint on push/PR)
- [x] Go linting (golangci-lint)
- [x] Code coverage badge
- [x] Multi-platform build matrix (linux, darwin, windows)
- [x] Automated PR checks
- [x] Release workflow (tag-triggered)
- [x] Cross-platform binary builds (GOOS/GOARCH)
- [x] GitHub Releases with artifacts

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
