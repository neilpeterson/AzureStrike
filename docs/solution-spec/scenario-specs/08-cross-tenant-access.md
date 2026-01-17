# Scenario 08: Cross-Tenant Access

**Status:** Planned
**Difficulty:** Advanced
**Estimated Time:** 35-40 minutes

## Overview

Players discover and exploit misconfigured cross-tenant access settings to gain unauthorized access to resources in a partner organization's Azure tenant.

## Learning Objectives

- Understand Azure AD B2B and cross-tenant access
- Learn cross-tenant access policies
- Identify misconfigured external collaboration settings
- Recognize multi-tenant application risks

## Attack Narrative

Contoso has a business partnership with Fabrikam and has configured cross-tenant access. You've discovered that the cross-tenant settings are overly permissive, potentially allowing unauthorized access from compromised Fabrikam accounts to Contoso resources.

## Starting Intel

- Primary tenant: Contoso (target)
- Partner tenant: Fabrikam (compromised)
- B2B relationship exists
- Cross-tenant sync may be enabled

## Resources

### Contoso Tenant Configuration
- External collaboration: Allow invitations to any domain
- Cross-tenant access: Default allow all from Fabrikam
- Guest user access: Same as members (dangerous)

### Fabrikam Resources (attacker controlled)
- Compromised user: attacker@fabrikam.com
- Guest invitation to: Contoso tenant

### Cross-Tenant Access Policy
```yaml
fabrikam.com:
  b2b_collaboration:
    inbound:
      users: AllUsers
      applications: AllApplications
    outbound:
      users: AllUsers
      applications: AllApplications
  cross_tenant_sync:
    enabled: true  # Dangerous
```

### Accessible Resources (via guest access)
- SharePoint sites
- Teams channels
- Azure subscriptions (if guest has role assignment)

## Objectives

| Order | ID | Description | Trigger | Points |
|-------|-----|-------------|---------|--------|
| 1 | enum_tenants | Enumerate tenant information and relationships | `az account tenant list` | 75 |
| 2 | check_guest_access | Identify guest user permissions | `az ad user show` | 100 |
| 3 | cross_tenant_enum | Enumerate resources accessible via guest access | `az resource list` | 150 |
| 4 | access_partner_data | Access sensitive data in partner tenant | `regex:storage\|keyvault\|graph` | 200 |

**Total Points:** 525

## Hints

### enum_tenants
- Level 1 (15 pts): "Azure AD can have relationships with other tenants"
- Level 2 (35 pts): "az account tenant list; check for multiple tenants"

### check_guest_access
- Level 1 (20 pts): "Guest users have a userType property in Azure AD"
- Level 2 (50 pts): "az ad user show --id attacker@fabrikam.com"

### cross_tenant_enum
- Level 1 (30 pts): "Guests may have access to Azure resources if role assigned"
- Level 2 (75 pts): "az resource list --subscription <contoso-sub-id>"

### access_partner_data
- Level 1 (40 pts): "With resource access, what sensitive data can you reach?"
- Level 2 (100 pts): "Try storage accounts, Key Vaults, or Graph API queries"

## Solution Walkthrough

```bash
# Step 1: List accessible tenants
az account tenant list
az account list

# Step 2: Check guest user details
az ad user show --id attacker@fabrikam.com
az ad user list --filter "userType eq 'Guest'"

# Step 3: Enumerate accessible resources in Contoso
az account set --subscription <contoso-subscription>
az resource list
az role assignment list --assignee attacker@fabrikam.com

# Step 4: Access partner data
az storage account list
az keyvault list
# Or Graph API calls for SharePoint/Teams
```

## Debrief Topics

- Azure AD B2B collaboration security
- Cross-tenant access policies
- External identities governance
- Guest user permissions restrictions
- Conditional Access for guests
- External collaboration settings
- Tenant restrictions

## Real-World References

- MITRE ATT&CK: T1199 (Trusted Relationship)
- Microsoft Entra: Cross-tenant access overview
- Azure Security: External collaboration security

## Implementation Notes

Requires:
- Multi-tenant mock support
- az account tenant list command
- Guest user representation
- Cross-tenant resource access simulation
- Graph API mock for collaboration resources
