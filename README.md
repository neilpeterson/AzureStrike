# AzureStrike

[![CI](https://github.com/neilpeterson/AzureStrike/actions/workflows/ci.yml/badge.svg)](https://github.com/neilpeterson/AzureStrike/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/neilpeterson/AzureStrike/branch/main/graph/badge.svg)](https://codecov.io/gh/neilpeterson/AzureStrike)
[![Go Report Card](https://goreportcard.com/badge/github.com/neilpeterson/AzureStrike)](https://goreportcard.com/report/github.com/neilpeterson/AzureStrike)

A terminal-based game that simulates Azure security scenarios for training security professionals. Features narrative-driven missions with a hybrid interface (CLI commands + TUI status panels) and a scoring/achievement system.

## Overview

AzureStrike provides an interactive learning environment where security professionals can practice identifying and exploiting common Azure misconfigurations in a safe, simulated environment. Each scenario presents a realistic attack narrative with objectives to complete using mocked Azure CLI commands.

## Features

- **Realistic Azure CLI**: Mocked `az` commands that return Azure-formatted responses
- **Scenario-Based Learning**: YAML-defined scenarios with progressive objectives
- **Interactive TUI**: Status panels showing objectives, score, and progress
- **Hint System**: Progressive hints available at a point cost
- **Scoring System**: Points for completing objectives with bonuses and achievements
- **Educational Debriefs**: Post-scenario summaries with remediation recommendations

## Installation

```bash
# Clone the repository
git clone https://github.com/azurestrike/azurestrike.git
cd azurestrike

# Build the binary
make build

# Or using go directly
go build -o azurestrike ./cmd/azurestrike
```

## Usage

```bash
# List available scenarios
./azurestrike --list

# Start a specific scenario
./azurestrike --scenario storage-breach

# Run in development mode
make run
```

## Supported Commands

AzureStrike intercepts and mocks the following command patterns:

### Azure CLI (az)
- `az storage account list/show` - List and show storage accounts
- `az storage container list` - List containers in a storage account
- `az storage blob list/download` - List and download blobs
- `az ad user list/show` - List and show Azure AD users
- `az ad sp list/show` - List and show service principals
- `az vm list/show` - List and show virtual machines
- `az account show/list` - Show subscription information

### Network
- `curl http://169.254.169.254/metadata/...` - Azure IMDS requests

### Game Commands
- `objectives` / `obj` - Show current objectives
- `score` - Show current score
- `hint [objective_id]` - Get a hint for an objective
- `help [command]` - Show help

## Scenario Modules

### Available Scenarios

| ID | Name | Difficulty | Description |
|----|------|------------|-------------|
| 01 | Storage Misconfiguration Discovery | Beginner | Identify publicly accessible storage containers exposing sensitive data |

### Planned Scenarios

| ID | Name | Difficulty | Description |
|----|------|------------|-------------|
| 02 | IMDS Token Theft | Beginner | Exploit Azure IMDS from a compromised VM to steal managed identity tokens |
| 03 | Service Principal Exposure | Beginner | Discover and abuse exposed service principal credentials |
| 04 | NSG Misconfiguration | Intermediate | Identify overly permissive network security group rules |
| 05 | Privilege Escalation via RBAC | Intermediate | Discover and exploit excessive RBAC permissions |
| 06 | Key Vault Secrets Exfiltration | Intermediate | Access Key Vault secrets using compromised credentials |
| 07 | Managed Identity Abuse | Intermediate | Abuse system-assigned managed identity for lateral movement |
| 08 | Cross-Tenant Access | Advanced | Exploit misconfigured cross-tenant access policies |
| 09 | Azure Function Code Injection | Advanced | Discover and exploit vulnerable Azure Functions |
| 10 | Full Kill Chain | Advanced | Multi-stage attack combining multiple techniques |

## Architecture

```
cmd/azurestrike/       Entry point
internal/
  game/                Game engine (state, objectives, scoring)
  azure/               Mocked Azure environment
    entra/             Users, groups, service principals
    storage/           Storage accounts, blobs, SAS tokens
    compute/           VMs, networking, NSGs
  cli/                 Command parser and handlers
    az/                Azure CLI mock handlers
  tui/                 Terminal UI (Bubble Tea components)
  scenario/            YAML scenario loader
scenarios/             YAML scenario definitions
```

## Creating Scenarios

Scenarios are defined in YAML files under the `scenarios/` directory. Each scenario includes:

- **Briefing**: Narrative introduction and initial intel
- **Resources**: Mocked Azure resources (storage, VMs, users, etc.)
- **Objectives**: Goals with trigger patterns and point values
- **Hints**: Progressive hints with point costs
- **Debrief**: Educational summary and remediation steps

See `scenarios/01-storage-breach/scenario.yaml` for a complete example.

### Resource Types

Scenarios can define the following resource types:

```yaml
resources:
  storage_accounts:     # Azure Storage accounts with containers and blobs
  users:               # Entra ID (Azure AD) users
  service_principals:  # Entra ID service principals
  groups:              # Entra ID groups
  virtual_machines:    # Azure VMs with identity configuration
  network_security_groups:  # NSGs with security rules
```

### Objective Triggers

Objectives are completed when the player executes commands matching trigger patterns:

- **Simple substring**: `"az storage account"` - matches any command containing this substring
- **Regex pattern**: `"regex:az storage blob (list|download)"` - matches regex pattern

## Development

```bash
# Run tests
make test

# Run tests with coverage
make test-cover

# Run tests with race detector
make test-race

# Format code
make fmt

# Run go vet
make vet

# Run golangci-lint
make lint

# Run all CI checks locally
make ci
```

## Technology Stack

- **Language**: Go
- **TUI Framework**: Bubble Tea + Lip Gloss (Charm ecosystem)
- **Scenario Definitions**: YAML files
- **Build**: Standard Go modules

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit pull requests with new scenarios, bug fixes, or feature improvements.
