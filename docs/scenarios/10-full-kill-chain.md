# Scenario 10: Full Kill Chain

**Status:** Planned
**Difficulty:** Advanced
**Estimated Time:** 45-60 minutes

## Overview

A comprehensive multi-stage attack scenario that combines techniques from previous scenarios into a complete attack chain from initial access to data exfiltration.

## Learning Objectives

- Execute a complete cloud attack kill chain
- Chain multiple vulnerabilities for maximum impact
- Practice realistic attack sequences
- Understand defense-in-depth importance

## Attack Narrative

You're conducting a full penetration test of Contoso's Azure environment. Starting with minimal information, you must chain together multiple vulnerabilities to achieve full compromise: from initial reconnaissance through privilege escalation to data exfiltration.

## Starting Intel

- Target: Contoso Corporation
- Entry point: Publicly accessible web application
- Goal: Exfiltrate customer database
- Rules of engagement: Full scope authorized

## Kill Chain Stages

### Stage 1: Reconnaissance
- Enumerate public-facing resources
- Discover storage misconfigurations
- Identify exposed credentials

### Stage 2: Initial Access
- Exploit public storage container
- Extract service principal credentials
- Authenticate to Azure

### Stage 3: Discovery
- Enumerate Azure AD users and groups
- Map role assignments
- Identify high-value targets

### Stage 4: Privilege Escalation
- Exploit excessive RBAC permissions
- Gain access to Key Vault
- Obtain database credentials

### Stage 5: Lateral Movement
- Access production storage
- Compromise VM via managed identity
- Move to database tier

### Stage 6: Data Exfiltration
- Access production database
- Extract customer data
- Establish persistence (optional)

## Resources

### Complete Environment
Combines resources from all previous scenarios:
- Storage accounts (public and private)
- Azure AD users, groups, service principals
- Virtual machines with managed identities
- Key Vaults with secrets
- Network Security Groups
- Azure Functions

### Attack Path
```
Public Storage → Creds → Service Principal → RBAC Enum →
Privilege Escalation → Key Vault → DB Creds → Customer Data
```

## Objectives

| Stage | ID | Description | Trigger | Points |
|-------|-----|-------------|---------|--------|
| Recon | recon_storage | Discover exposed storage account | `az storage account list` | 50 |
| Recon | recon_containers | Find public containers | `az storage container list` | 50 |
| Access | initial_creds | Extract credentials from storage | `az storage blob download` | 75 |
| Access | authenticate | Authenticate with stolen creds | `az login` | 75 |
| Discovery | enum_users | Enumerate Azure AD | `az ad user list` | 75 |
| Discovery | enum_roles | Map role assignments | `az role assignment list` | 75 |
| Escalation | find_escalation | Identify privilege escalation path | `regex:User Access Administrator` | 100 |
| Escalation | escalate | Escalate privileges | `az role assignment create` | 125 |
| Lateral | access_keyvault | Access Key Vault | `az keyvault secret show` | 100 |
| Lateral | get_db_creds | Obtain database credentials | `regex:keyvault.*db\|sql` | 100 |
| Exfil | access_database | Connect to production database | `regex:sql.*connect\|query` | 125 |
| Exfil | exfiltrate_data | Exfiltrate customer data | `regex:SELECT.*customers\|export` | 150 |

**Total Points:** 1100

## Hints

Hints are provided per-stage, with increasing specificity:

### Stage 1 Hints
- Level 1: "Start with public enumeration - what's exposed to the internet?"
- Level 2: "Check storage accounts for public access"

### Stage 2 Hints
- Level 1: "Configuration backups often contain credentials"
- Level 2: "Download config files and look for Azure credentials"

### Stage 3 Hints
- Level 1: "With valid credentials, enumerate the directory"
- Level 2: "az ad user list; az role assignment list --all"

### Stage 4 Hints
- Level 1: "Look for roles that can grant other roles"
- Level 2: "User Access Administrator can assign any role"

### Stage 5 Hints
- Level 1: "Key Vault contains production secrets"
- Level 2: "Look for database connection strings in secrets"

### Stage 6 Hints
- Level 1: "Use the database credentials to connect"
- Level 2: "Query the customers table for PII"

## Solution Walkthrough

```bash
# Stage 1: Reconnaissance
az storage account list
az storage container list --account-name contosostorage2024

# Stage 2: Initial Access
az storage blob download --account-name contosostorage2024 --container-name backups --name config-backup.tar
# Extract credentials from config

# Stage 3: Discovery
az ad user list
az ad sp list
az role assignment list --all

# Stage 4: Privilege Escalation
# Discover UAA role on current user
az role assignment list --assignee <current-user>
# Grant Contributor role
az role assignment create --assignee <current-user> --role Contributor --scope /subscriptions/<sub>

# Stage 5: Lateral Movement
az keyvault list
az keyvault secret list --vault-name contoso-prod-kv
az keyvault secret show --vault-name contoso-prod-kv --name db-connection-string

# Stage 6: Data Exfiltration
az sql db list --server contoso-sql --resource-group prod-rg
# Use credentials to query database
# SELECT * FROM customers WHERE ...
```

## Debrief Topics

- Defense in depth
- Kill chain detection opportunities
- Logging and monitoring gaps
- Security controls at each stage
- Incident response considerations
- Prevention vs detection

## Real-World References

- MITRE ATT&CK: Full cloud attack chain
- Cyber Kill Chain (Lockheed Martin)
- Azure Security Center alerts
- Microsoft Defender for Cloud attack paths

## Implementation Notes

This scenario requires:
- All previous scenario implementations complete
- Linked resources across scenarios
- State persistence between commands
- Complex trigger patterns for multi-step objectives
- Achievement system for completing full chain
- Time-based scoring bonuses
