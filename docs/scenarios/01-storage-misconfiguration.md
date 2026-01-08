# Scenario 01: Storage Misconfiguration Discovery

**Status:** Complete
**Difficulty:** Beginner
**Estimated Time:** 10-15 minutes

## Overview

Players discover and exploit a misconfigured Azure Storage account with publicly accessible containers containing sensitive backup data.

## Learning Objectives

- Understand Azure Storage account enumeration using Azure CLI
- Identify public vs private container access levels
- Recognize common data exposure patterns in cloud storage
- Learn remediation steps for storage misconfigurations

## Attack Narrative

Intel indicates Contoso Corporation may have misconfigured storage accounts. The player must enumerate storage resources, identify publicly accessible containers, and extract sensitive data to document the security exposure.

## Starting Intel

- Target: Contoso Corporation
- Known storage account: `contosostorage2024`
- Suspected container: `backups`

## Resources

### Storage Account: contosostorage2024
- Location: eastus
- Resource Group: production-rg
- Public Access: Enabled
- HTTPS Only: Disabled (vulnerability)
- Min TLS: 1.0 (vulnerability)

### Containers

| Container | Access Level | Contents |
|-----------|--------------|----------|
| public-assets | Blob | logo.png, style.css |
| backups | Container (PUBLIC) | db-backup, config-backup.tar, secrets.txt |
| private-data | None | employee-records.xlsx |

## Objectives

| Order | ID | Description | Trigger | Points |
|-------|-----|-------------|---------|--------|
| 1 | discover_account | Discover the storage account configuration | `az storage account` | 50 |
| 2 | list_containers | List containers in the storage account | `az storage container list` | 75 |
| 3 | find_public_container | Identify the publicly accessible backup container | `regex:az storage blob list.*backups` | 100 |
| 4 | extract_secrets | Download and examine sensitive backup files | `regex:az storage blob download.*(secrets\|config)` | 150 |

**Total Points:** 375

## Hints

### discover_account
- Level 1 (10 pts): "Azure CLI has commands to interact with storage accounts. Try exploring 'az storage'."
- Level 2 (25 pts): "Use 'az storage account list' or 'az storage account show --name <account>'"

### list_containers
- Level 1 (15 pts): "Storage accounts contain containers. What command might list them?"
- Level 2 (35 pts): "Try: az storage container list --account-name contosostorage2024"

### find_public_container
- Level 1 (20 pts): "One of the containers has interesting contents related to backups..."
- Level 2 (50 pts): "az storage blob list --account-name contosostorage2024 --container-name backups"

### extract_secrets
- Level 1 (25 pts): "Look at the blob names - which ones might contain sensitive information?"
- Level 2 (50 pts): "Use 'az storage blob download' with --name to download specific files"

## Solution Walkthrough

```bash
# Step 1: Enumerate storage account
az storage account list
az storage account show --name contosostorage2024

# Step 2: List containers
az storage container list --account-name contosostorage2024

# Step 3: Identify public container and list blobs
az storage blob list --account-name contosostorage2024 --container-name backups

# Step 4: Download sensitive files
az storage blob download --account-name contosostorage2024 --container-name backups --name secrets.txt
az storage blob download --account-name contosostorage2024 --container-name backups --name config-backup.tar
```

## Debrief Topics

- Azure Blob Storage access levels (Private, Blob, Container)
- Risks of public blob access
- Credential rotation after exposure
- Azure Policy for preventing public access
- Microsoft Defender for Storage
- Storage account firewall rules
- Private endpoints

## Real-World References

- MITRE ATT&CK: T1530 (Data from Cloud Storage Object)
- CIS Azure Benchmark: 3.6 (Ensure public access level is set to private)
- Azure Security Baseline: Storage security recommendations

## Implementation Notes

This scenario is fully implemented in `scenarios/01-storage-breach/scenario.yaml`.
