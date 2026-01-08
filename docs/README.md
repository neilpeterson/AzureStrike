# AzureStrike Documentation

## Core Documentation

- [CLI Specification](cli-spec.md) - Mocked Azure CLI command reference and behavior

## Scenario Specifications

| ID | Name | Difficulty | Status | Doc |
|----|------|------------|--------|-----|
| 01 | Storage Misconfiguration Discovery | Beginner | Complete | [Spec](scenarios/01-storage-misconfiguration.md) |
| 02 | IMDS Token Theft | Beginner | Planned | [Spec](scenarios/02-imds-token-theft.md) |
| 03 | Service Principal Exposure | Beginner | Planned | [Spec](scenarios/03-service-principal-exposure.md) |
| 04 | NSG Misconfiguration | Intermediate | Planned | [Spec](scenarios/04-nsg-misconfiguration.md) |
| 05 | Privilege Escalation via RBAC | Intermediate | Planned | [Spec](scenarios/05-privilege-escalation-rbac.md) |
| 06 | Key Vault Secrets Exfiltration | Intermediate | Planned | [Spec](scenarios/06-keyvault-secrets.md) |
| 07 | Managed Identity Abuse | Intermediate | Planned | [Spec](scenarios/07-managed-identity-abuse.md) |
| 08 | Cross-Tenant Access | Advanced | Planned | [Spec](scenarios/08-cross-tenant-access.md) |
| 09 | Azure Function Code Injection | Advanced | Planned | [Spec](scenarios/09-azure-function-injection.md) |
| 10 | Full Kill Chain | Advanced | Planned | [Spec](scenarios/10-full-kill-chain.md) |

## Scenario Spec Template

Each scenario specification includes:

- **Overview**: Brief description of the attack scenario
- **Learning Objectives**: What players will learn
- **Attack Narrative**: Story-driven context for the mission
- **Starting Intel**: Information provided to the player
- **Resources**: Mock Azure resources needed
- **Objectives**: Ordered goals with triggers and points
- **Hints**: Progressive hints with point costs
- **Solution Walkthrough**: Complete solution commands
- **Debrief Topics**: Educational content for post-scenario
- **Real-World References**: MITRE ATT&CK, CIS benchmarks, etc.
- **Implementation Notes**: Technical requirements for building

## Difficulty Levels

| Level | Target Audience | Time | Complexity |
|-------|-----------------|------|------------|
| Beginner | New to Azure security | 10-20 min | Single technique |
| Intermediate | Some Azure experience | 20-35 min | Multiple techniques |
| Advanced | Experienced practitioners | 35-60 min | Attack chains |

## MITRE ATT&CK Coverage

| Tactic | Techniques Covered |
|--------|-------------------|
| Initial Access | T1190, T1078.004 |
| Execution | T1059.007 |
| Persistence | T1098 |
| Privilege Escalation | T1098 |
| Defense Evasion | - |
| Credential Access | T1552.001, T1552.005 |
| Discovery | T1087, T1069, T1530 |
| Lateral Movement | T1550.001 |
| Collection | T1530 |
| Exfiltration | T1567 |
