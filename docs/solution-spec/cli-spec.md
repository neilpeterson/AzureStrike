# AzureStrike CLI Specification

This document specifies how the mocked Azure CLI and other commands behave in AzureStrike.

## Design Principles

1. **Realistic Output**: Commands should return JSON output matching real Azure CLI format
2. **Scenario-Driven**: All data comes from the scenario YAML, not hardcoded values
3. **Error Handling**: Return realistic Azure error messages for invalid operations
4. **Access Control**: Respect container/resource access levels defined in scenarios

## Command Structure

```
az <group> <subgroup> <command> [options]
```

## Supported Command Groups

### az account

Manage Azure subscription context.

| Command | Description | Output |
|---------|-------------|--------|
| `az account show` | Show current subscription | Subscription JSON |
| `az account list` | List all subscriptions | Array of subscriptions |

**Options:** `--output`, `--query`

**Example Output:**
```json
{
  "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "name": "Contoso Production",
  "tenantId": "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy",
  "state": "Enabled",
  "isDefault": true
}
```

---

### az storage account

Manage storage accounts.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az storage account list` | List storage accounts | - | Array of accounts |
| `az storage account show` | Show account details | `--name` | Account JSON |

**Data Source:** `scenario.resources.storage_accounts[]`

**Example Output:**
```json
{
  "name": "contosostorage2024",
  "resourceGroup": "production-rg",
  "location": "eastus",
  "kind": "StorageV2",
  "sku": "Standard_LRS",
  "publicAccess": true,
  "httpsOnly": false,
  "minimumTlsVersion": "TLS1_0"
}
```

---

### az storage container

Manage blob containers.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az storage container list` | List containers | `--account-name` | Array of containers |

**Data Source:** `scenario.resources.storage_accounts[].containers[]`

**Access Control:** Always returns container list (metadata is not protected)

**Example Output:**
```json
[
  {
    "name": "backups",
    "publicAccess": "container"
  },
  {
    "name": "private-data",
    "publicAccess": "none"
  }
]
```

---

### az storage blob

Manage blobs.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az storage blob list` | List blobs in container | `--account-name`, `--container-name` | Array of blobs |
| `az storage blob download` | Download blob content | `--account-name`, `--container-name`, `--name` | Blob content |

**Options:**
- `--file`, `-f`: Output file path (simulated)
- `--auth-mode`: `key` (default) or `login`

**Access Control:**
| Container Access | blob list | blob download |
|------------------|-----------|---------------|
| `container` | Allowed | Allowed |
| `blob` | Denied (403) | Allowed (if URL known) |
| `none` | Denied (403) | Denied (403) |

**Error Response (403):**
```
This request is not authorized to perform this operation.
RequestId:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Time:2024-01-01T00:00:00.0000000Z
```

**Example Output (list):**
```json
[
  {
    "name": "secrets.txt",
    "contentLength": 512,
    "contentType": "text/plain"
  }
]
```

---

### az ad user

Manage Azure AD users.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az ad user list` | List all users | - | Array of users |
| `az ad user show` | Show user details | `--id` or `--upn-or-object-id` | User JSON |

**Data Source:** `scenario.resources.users[]`

**Example Output:**
```json
{
  "id": "aaaa-bbbb-cccc-dddd",
  "displayName": "Admin User",
  "userPrincipalName": "admin@contoso.com",
  "mail": "admin@contoso.com",
  "jobTitle": "IT Administrator",
  "department": "IT"
}
```

---

### az ad sp

Manage service principals.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az ad sp list` | List service principals | - | Array of SPs |
| `az ad sp show` | Show SP details | `--id` | SP JSON |

**Data Source:** `scenario.resources.service_principals[]`

**Note:** Secret values are never returned in output (matches real Azure behavior)

**Example Output:**
```json
{
  "id": "1111-2222-3333-4444",
  "appId": "xxxx-yyyy-zzzz-0000",
  "displayName": "backup-automation",
  "permissions": ["Directory.Read.All"]
}
```

---

### az ad group

Manage Azure AD groups.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az ad group list` | List all groups | - | Array of groups |
| `az ad group show` | Show group details | `--group` | Group JSON |
| `az ad group member list` | List group members | `--group` | Array of members |

**Data Source:** `scenario.resources.groups[]`

---

### az vm

Manage virtual machines.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az vm list` | List VMs | - | Array of VMs |
| `az vm show` | Show VM details | `--name` | VM JSON |

**Options:**
- `--resource-group`, `-g`: Filter by resource group

**Data Source:** `scenario.resources.virtual_machines[]`

**Example Output:**
```json
{
  "id": "/subscriptions/.../virtualMachines/web-server-01",
  "name": "web-server-01",
  "resourceGroup": "production-rg",
  "location": "eastus",
  "vmSize": "Standard_D2s_v3",
  "powerState": "VM running",
  "publicIpAddress": "20.xxx.xxx.xxx",
  "privateIpAddress": "10.0.1.4",
  "osType": "Linux",
  "identity": {
    "type": "SystemAssigned",
    "principalId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }
}
```

---

### az network nsg

Manage network security groups.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az network nsg list` | List NSGs | - | Array of NSGs |
| `az network nsg show` | Show NSG details | `--name` | NSG JSON |
| `az network nsg rule list` | List NSG rules | `--nsg-name` | Array of rules |

**Options:**
- `--resource-group`, `-g`: Filter/scope by resource group

**Data Source:** `scenario.resources.network_security_groups[]`

**Example Output (rule):**
```json
{
  "name": "AllowSSH",
  "priority": 100,
  "direction": "Inbound",
  "access": "Allow",
  "protocol": "Tcp",
  "sourcePortRange": "*",
  "destinationPortRange": "22",
  "sourceAddressPrefix": "*",
  "destinationAddressPrefix": "*"
}
```

---

### az keyvault (Planned)

Manage Key Vaults.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az keyvault list` | List Key Vaults | - | Array of vaults |
| `az keyvault show` | Show vault details | `--name` | Vault JSON |
| `az keyvault secret list` | List secrets | `--vault-name` | Array of secret metadata |
| `az keyvault secret show` | Get secret value | `--vault-name`, `--name` | Secret JSON with value |

**Access Control:** Based on scenario-defined access policies

---

### az role (Planned)

Manage RBAC.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az role assignment list` | List role assignments | - | Array of assignments |
| `az role assignment create` | Create assignment | `--assignee`, `--role` | Assignment JSON |
| `az role definition list` | List role definitions | - | Array of definitions |

**Options:**
- `--assignee`: Filter by user/SP
- `--scope`: Filter by resource scope
- `--all`: Include inherited assignments

---

### az functionapp (Planned)

Manage Azure Functions.

| Command | Description | Required Args | Output |
|---------|-------------|---------------|--------|
| `az functionapp list` | List function apps | - | Array of apps |
| `az functionapp show` | Show app details | `--name` | App JSON |
| `az functionapp config appsettings list` | List app settings | `--name` | Settings array |

---

## cat

View blob contents directly from storage containers.

### Syntax

```
cat <path>
```

### Path Formats

| Format | Description | Example |
|--------|-------------|---------|
| `account/container/blob` | Full path with account | `cat contosostorage2024/backups/secrets.txt` |
| `container/blob` | Short path (auto-discovers account) | `cat backups/secrets.txt` |

### Access Control

- Returns content only from containers with `public_access: blob` or `public_access: container`
- Private containers return: `cat: <path>: Permission denied`

### Error Responses

| Error | Cause |
|-------|-------|
| `cat: missing file operand` | No path provided |
| `cat: <path>: No such file or directory` | Blob/container not found |
| `cat: <path>: Permission denied` | Container is private |

### Example

```bash
$ cat backups/secrets.txt
=== EMERGENCY ACCESS CREDENTIALS ===
Last Updated: 2024-01-15

Azure Portal Admin: admin@contoso.com / AzureAdmin#2024!
...
```

---

## curl

HTTP client for IMDS and web requests.

### Supported Endpoints

#### Azure IMDS (169.254.169.254)

| Endpoint | Description | Required Header |
|----------|-------------|-----------------|
| `/metadata/instance` | VM instance metadata | `Metadata: true` |
| `/metadata/identity/oauth2/token` | Get managed identity token | `Metadata: true` |

**Query Parameters (token endpoint):**
- `api-version`: Required (e.g., `2018-02-01`)
- `resource`: Target resource (e.g., `https://management.azure.com/`)

**Error Handling:**
- Missing `Metadata: true` header: `400 Bad Request`
- Invalid endpoint: `404 Not Found`

**Example Token Response:**
```json
{
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIs...",
  "client_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "expires_in": "28799",
  "expires_on": "1234567890",
  "resource": "https://management.azure.com/",
  "token_type": "Bearer"
}
```

---

## Global Options

All `az` commands support:

| Option | Short | Alternate | Description |
|--------|-------|-----------|-------------|
| `--help` | `-h` | `-help` | Show command help |
| `--output` | `-o` | - | Output format: `json` (default), `table`, `tsv` |
| `--query` | - | - | JMESPath query (not implemented) |
| `--verbose` | - | - | Increase logging verbosity (no-op) |
| `--debug` | - | - | Show debug output (no-op) |

### Output Formats

**JSON (default):**
```bash
$ az storage account list
[
  {
    "name": "contosostorage2024",
    "resourceGroup": "production-rg",
    ...
  }
]
```

**Table:**
```bash
$ az storage account list -o table
Name                  ResourceGroup    Location    Kind
--------------------  ---------------  ----------  ---------
contosostorage2024    production-rg    eastus      StorageV2
```

Table format displays columns appropriate for each command type with aligned headers.

---

## Error Handling

### Standard Error Format

```
az <command>: error: <message>
```

### Common Errors

| Error | Cause | Message |
|-------|-------|---------|
| Missing required arg | Argument not provided | `az storage blob list: --account-name is required` |
| Unrecognized argument | Invalid flag used | `az storage account list: unrecognized arguments: -f` |
| Resource not found | Invalid name/ID | `The storage account 'xxx' was not found` |
| Access denied | Insufficient permissions | `This request is not authorized to perform this operation.` |
| Unknown command | Invalid command | `az: 'xxx' is not an az command. See 'az --help'.` |

### Argument Validation

Commands validate that all provided arguments are recognized. Unknown flags return an error matching Azure CLI behavior:

```
$ az storage account list -f
az storage account list: unrecognized arguments: -f
```

---

## Scenario YAML Schema

### Storage Account

```yaml
storage_accounts:
  - name: string              # Account name (required)
    resource_group: string    # Resource group (required)
    location: string          # Azure region
    kind: string              # StorageV2, BlobStorage, etc.
    sku: string               # Standard_LRS, etc.
    public_access: boolean    # Allow public blob access
    https_only: boolean       # Require HTTPS
    min_tls_version: string   # TLS1_0, TLS1_2
    tags: map[string]string   # Resource tags
    containers:
      - name: string          # Container name
        public_access: string # none, blob, container
        blobs:
          - name: string      # Blob name
            content: string   # Blob content (returned on download)
            content_type: string
            size: integer
```

### User

```yaml
users:
  - id: string                    # Object ID
    display_name: string          # Display name
    user_principal_name: string   # UPN (email format)
    mail: string                  # Email
    job_title: string
    department: string
    roles: [string]               # Assigned roles
    groups: [string]              # Group memberships
```

### Service Principal

```yaml
service_principals:
  - id: string                    # Object ID
    app_id: string                # Application ID
    display_name: string
    permissions: [string]         # API permissions
    secrets:                      # Credentials (hidden from output)
      - key_id: string
        display_name: string
        end_date_time: string
        secret_text: string       # The actual secret
```

### Virtual Machine

```yaml
virtual_machines:
  - id: string                    # Resource ID
    name: string                  # VM name
    resource_group: string
    location: string
    vm_size: string               # Standard_D2s_v3, etc.
    power_state: string           # VM running, VM stopped
    public_ip: string
    private_ip: string
    os_type: string               # Linux, Windows
    admin_username: string
    tags: map[string]string
    identity:
      type: string                # SystemAssigned, UserAssigned
      principal_id: string
      tenant_id: string
```

### Network Security Group

```yaml
network_security_groups:
  - id: string
    name: string
    resource_group: string
    location: string
    rules:
      - name: string
        priority: integer         # 100-4096
        direction: string         # Inbound, Outbound
        access: string            # Allow, Deny
        protocol: string          # Tcp, Udp, *
        source_port_range: string
        destination_port_range: string
        source_address_prefix: string      # *, 0.0.0.0/0, IP, tag
        destination_address_prefix: string
```

---

## Implementation Status

| Command Group | Status | Notes |
|---------------|--------|-------|
| az account | Implemented | Basic show/list |
| az storage account | Implemented | Full support with -o table |
| az storage container | Implemented | Full support with -o table |
| az storage blob | Implemented | list/download with access control, -o table |
| az ad user | Implemented | list/show with -o table |
| az ad sp | Implemented | list/show with -o table |
| az ad group | Planned | Need group commands |
| az vm | Implemented | list/show with RG filter, -o table |
| az network nsg | Planned | Need NSG commands |
| az keyvault | Planned | Scenario 06 requirement |
| az role | Planned | Scenario 05 requirement |
| az functionapp | Planned | Scenario 09 requirement |
| curl (IMDS) | Implemented | Basic token/instance |
| cat | Implemented | View blob contents directly |
