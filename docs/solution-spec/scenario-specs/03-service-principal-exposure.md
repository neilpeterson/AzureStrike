# Scenario 03: Service Principal Exposure

**Status:** Planned
**Difficulty:** Beginner
**Estimated Time:** 15-20 minutes

## Overview

Players discover exposed service principal credentials in application configuration and use them to authenticate and enumerate Azure AD resources.

## Learning Objectives

- Understand Azure service principals and app registrations
- Recognize common credential exposure patterns
- Learn how to authenticate using service principal credentials
- Understand the risks of exposed application secrets

## Attack Narrative

During reconnaissance of a compromised web server, you discover configuration files containing Azure service principal credentials. Your mission is to use these credentials to authenticate to Azure AD and determine what access the service principal has.

## Starting Intel

- Found config file with Azure credentials
- App ID: `12345678-1234-1234-1234-123456789012`
- Tenant ID available in config
- Secret may be rotated - check for multiple secrets

## Resources

### Service Principal: backup-automation
- App ID: 12345678-1234-1234-1234-123456789012
- Display Name: backup-automation
- Secrets: 2 (one expired, one active)
- Permissions: Directory.Read.All, Storage.Blob.Contributor

### Configuration File (exposed)
```
AZURE_TENANT_ID=contoso.onmicrosoft.com
AZURE_CLIENT_ID=12345678-1234-1234-1234-123456789012
AZURE_CLIENT_SECRET=<secret>
```

### Users in Directory
- admin@contoso.com (Global Admin)
- svc-backup@contoso.com (Service Account)
- developer@contoso.com (Developer)

## Objectives

| Order | ID | Description | Trigger | Points |
|-------|-----|-------------|---------|--------|
| 1 | find_creds | Discover service principal credentials in config | `regex:cat.*config\|download.*config` | 50 |
| 2 | enumerate_sp | Enumerate service principals in the directory | `az ad sp list` | 75 |
| 3 | enumerate_users | List users in Azure AD | `az ad user list` | 100 |
| 4 | check_permissions | Identify the service principal's permissions | `regex:az ad sp show\|az role assignment` | 150 |

**Total Points:** 375

## Hints

### find_creds
- Level 1 (10 pts): "Application configurations often contain cloud credentials"
- Level 2 (25 pts): "Download the config-backup.tar from the storage account"

### enumerate_sp
- Level 1 (15 pts): "Azure AD service principals can be listed via CLI"
- Level 2 (35 pts): "az ad sp list --all"

### enumerate_users
- Level 1 (20 pts): "With directory read permissions, you can enumerate users"
- Level 2 (50 pts): "az ad user list"

### check_permissions
- Level 1 (25 pts): "Service principals have specific permissions - can you find what this one can do?"
- Level 2 (50 pts): "az ad sp show --id <app-id> or az role assignment list --assignee <app-id>"

## Solution Walkthrough

```bash
# Step 1: Download config file (from previous scenario or enumeration)
az storage blob download --account-name contosostorage2024 --container-name backups --name config-backup.tar

# Step 2: List service principals
az ad sp list --all

# Step 3: Enumerate directory users
az ad user list

# Step 4: Check service principal permissions
az ad sp show --id 12345678-1234-1234-1234-123456789012
az role assignment list --assignee 12345678-1234-1234-1234-123456789012
```

## Debrief Topics

- Service principal vs managed identity
- App registration security best practices
- Secret rotation policies
- Azure Key Vault for credential management
- Conditional Access for service principals
- Monitoring service principal sign-ins

## Real-World References

- MITRE ATT&CK: T1078.004 (Valid Accounts: Cloud Accounts)
- CIS Azure Benchmark: 1.11 (Ensure service principal credentials are rotated)
- Azure AD Security: Application security best practices

## Implementation Notes

Requires:
- Enhanced az ad sp commands
- Mock authentication with service principal
- Role assignment enumeration
- Links to Scenario 01 (config file in storage)
