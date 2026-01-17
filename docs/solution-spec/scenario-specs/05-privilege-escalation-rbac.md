# Scenario 05: Privilege Escalation via RBAC

**Status:** Planned
**Difficulty:** Intermediate
**Estimated Time:** 25-30 minutes

## Overview

Players discover excessive RBAC permissions that allow a low-privilege user to escalate privileges and gain unauthorized access to sensitive resources.

## Learning Objectives

- Understand Azure Role-Based Access Control (RBAC)
- Identify dangerous role assignments
- Recognize privilege escalation paths in Azure
- Learn RBAC security best practices

## Attack Narrative

You have compromised a developer account with seemingly limited permissions. Your mission is to enumerate the account's role assignments and discover if there are any privilege escalation paths that could grant broader access.

## Starting Intel

- Compromised account: developer@contoso.com
- Initial role: Reader on resource group
- Target: Gain write access to production resources

## Resources

### Users

| User | Initial Roles | Notes |
|------|---------------|-------|
| developer@contoso.com | Reader (dev-rg), User Access Administrator (prod-rg) | Dangerous UAA role |
| admin@contoso.com | Owner (subscription) | Target |

### Role Assignments

```
developer@contoso.com:
  - Reader on /subscriptions/.../resourceGroups/dev-rg
  - User Access Administrator on /subscriptions/.../resourceGroups/prod-rg  # ESCALATION PATH
```

### Dangerous Permissions (User Access Administrator)
- Microsoft.Authorization/roleAssignments/write
- Microsoft.Authorization/roleAssignments/delete
- Can grant themselves any role including Owner

## Objectives

| Order | ID | Description | Trigger | Points |
|-------|-----|-------------|---------|--------|
| 1 | check_identity | Verify current identity and context | `az account show` | 50 |
| 2 | enumerate_roles | Enumerate role assignments for current user | `az role assignment list` | 100 |
| 3 | identify_escalation | Identify User Access Administrator role | `regex:role.*definition\|User Access` | 125 |
| 4 | escalate_privileges | Grant yourself elevated permissions | `az role assignment create` | 175 |

**Total Points:** 450

## Hints

### check_identity
- Level 1 (10 pts): "First, understand who you're logged in as"
- Level 2 (25 pts): "az account show"

### enumerate_roles
- Level 1 (15 pts): "Azure RBAC assigns roles to identities - what roles do you have?"
- Level 2 (40 pts): "az role assignment list --assignee developer@contoso.com"

### identify_escalation
- Level 1 (25 pts): "Some roles can create other role assignments..."
- Level 2 (60 pts): "User Access Administrator can assign any role - check your assignments carefully"

### escalate_privileges
- Level 1 (30 pts): "If you can create role assignments, you can grant yourself more access"
- Level 2 (75 pts): "az role assignment create --assignee developer@contoso.com --role Contributor --resource-group prod-rg"

## Solution Walkthrough

```bash
# Step 1: Check current identity
az account show
az ad signed-in-user show

# Step 2: List role assignments
az role assignment list --assignee developer@contoso.com --all
az role assignment list --include-inherited

# Step 3: Check what User Access Administrator can do
az role definition list --name "User Access Administrator"

# Step 4: Escalate by assigning Contributor role
az role assignment create --assignee developer@contoso.com --role "Contributor" --resource-group prod-rg
# Or even Owner:
az role assignment create --assignee developer@contoso.com --role "Owner" --resource-group prod-rg
```

## Debrief Topics

- Principle of least privilege
- Dangerous built-in roles (Owner, User Access Administrator)
- Custom role definitions
- Azure AD Privileged Identity Management (PIM)
- Role assignment conditions (ABAC)
- Monitoring role assignments with Azure Monitor

## Real-World References

- MITRE ATT&CK: T1098 (Account Manipulation)
- CIS Azure Benchmark: 1.22 (Ensure no custom subscription owner roles exist)
- Azure Security: RBAC best practices

## Implementation Notes

Requires:
- az role assignment list/create commands
- az role definition list command
- RBAC resources in scenario YAML
- Mock role assignment validation
