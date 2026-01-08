# Scenario 06: Key Vault Secrets Exfiltration

**Status:** Planned
**Difficulty:** Intermediate
**Estimated Time:** 25-30 minutes

## Overview

Players discover an Azure Key Vault and exploit misconfigured access policies to exfiltrate stored secrets, keys, and certificates.

## Learning Objectives

- Understand Azure Key Vault architecture
- Learn Key Vault access policies vs RBAC
- Identify secrets, keys, and certificates enumeration
- Recognize Key Vault security misconfigurations

## Attack Narrative

During your assessment, you discover references to an Azure Key Vault. Using previously obtained credentials, you need to determine if you can access the vault and extract any sensitive secrets stored within.

## Starting Intel

- Discovered Key Vault: `contoso-prod-kv`
- Managed identity has GET permissions
- Vault may contain database passwords and API keys

## Resources

### Key Vault: contoso-prod-kv
- Location: eastus
- SKU: Standard
- Soft Delete: Enabled
- Purge Protection: Disabled (vulnerability)

### Access Policies
| Principal | Keys | Secrets | Certificates |
|-----------|------|---------|--------------|
| backup-automation (SP) | Get, List | Get, List | Get, List |
| admin@contoso.com | All | All | All |

### Secrets
| Name | Content Type | Value (redacted for spec) |
|------|--------------|---------------------------|
| db-connection-string | text/plain | Server=...;Password=... |
| api-key-stripe | text/plain | sk_live_... |
| storage-account-key | text/plain | Base64 key |
| ssl-certificate-password | text/plain | Certificate passphrase |

## Objectives

| Order | ID | Description | Trigger | Points |
|-------|-----|-------------|---------|--------|
| 1 | discover_vault | Discover Key Vault resources | `az keyvault list` | 75 |
| 2 | list_secrets | List secrets in the vault | `az keyvault secret list` | 100 |
| 3 | read_secret | Read a secret value | `az keyvault secret show` | 125 |
| 4 | exfil_all | Extract multiple secrets | `regex:keyvault secret show.*(db\|api\|storage)` | 150 |

**Total Points:** 450

## Hints

### discover_vault
- Level 1 (15 pts): "Azure Key Vault stores secrets - can you find any vaults?"
- Level 2 (35 pts): "az keyvault list"

### list_secrets
- Level 1 (20 pts): "Once you find a vault, enumerate what's inside"
- Level 2 (50 pts): "az keyvault secret list --vault-name contoso-prod-kv"

### read_secret
- Level 1 (25 pts): "Secrets have names - try to read their values"
- Level 2 (60 pts): "az keyvault secret show --vault-name contoso-prod-kv --name <secret-name>"

### exfil_all
- Level 1 (30 pts): "Database credentials and API keys are high-value targets"
- Level 2 (75 pts): "Extract db-connection-string, api-key-stripe, and storage-account-key"

## Solution Walkthrough

```bash
# Step 1: List Key Vaults
az keyvault list
az keyvault show --name contoso-prod-kv

# Step 2: List secrets
az keyvault secret list --vault-name contoso-prod-kv

# Step 3: Read a secret
az keyvault secret show --vault-name contoso-prod-kv --name db-connection-string

# Step 4: Extract all valuable secrets
az keyvault secret show --vault-name contoso-prod-kv --name api-key-stripe
az keyvault secret show --vault-name contoso-prod-kv --name storage-account-key
```

## Debrief Topics

- Key Vault access policies vs Azure RBAC
- Secret versioning and rotation
- Soft delete and purge protection
- Private endpoints for Key Vault
- Key Vault firewall rules
- Managed identities for Key Vault access
- Azure Monitor for Key Vault auditing

## Real-World References

- MITRE ATT&CK: T1552.001 (Unsecured Credentials: Credentials In Files)
- CIS Azure Benchmark: 8.1-8.5 (Key Vault recommendations)
- Azure Security Baseline: Key Vault security

## Implementation Notes

Requires:
- az keyvault list/show commands
- az keyvault secret list/show commands
- Key Vault resources in scenario YAML
- Mock secret values with realistic formats
