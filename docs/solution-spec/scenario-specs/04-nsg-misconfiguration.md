# Scenario 04: NSG Misconfiguration

**Status:** Planned
**Difficulty:** Intermediate
**Estimated Time:** 20-25 minutes

## Overview

Players analyze Network Security Group (NSG) configurations to identify overly permissive rules that expose internal services to the internet.

## Learning Objectives

- Understand Azure Network Security Groups
- Identify dangerous inbound rules (0.0.0.0/0, Any)
- Recognize common NSG misconfigurations
- Learn NSG security best practices

## Attack Narrative

As part of a security assessment, you need to analyze the network security posture of Contoso's Azure infrastructure. Intelligence suggests their NSGs may have overly permissive rules allowing unauthorized access to internal services.

## Starting Intel

- Target resource group: production-rg
- Multiple VMs with attached NSGs
- Known VM: web-server-01
- Suspected: Management ports exposed to internet

## Resources

### Network Security Group: web-nsg
- Resource Group: production-rg
- Associated with: web-server-01 NIC

#### Rules (Dangerous)
| Priority | Name | Direction | Access | Source | Destination | Port |
|----------|------|-----------|--------|--------|-------------|------|
| 100 | AllowSSH | Inbound | Allow | * | * | 22 |
| 110 | AllowRDP | Inbound | Allow | 0.0.0.0/0 | * | 3389 |
| 120 | AllowHTTP | Inbound | Allow | Internet | * | 80 |
| 130 | AllowHTTPS | Inbound | Allow | Internet | * | 443 |
| 200 | AllowAllOutbound | Outbound | Allow | * | * | * |

### Network Security Group: db-nsg
- Associated with: db-server-01 NIC

#### Rules
| Priority | Name | Direction | Access | Source | Destination | Port |
|----------|------|-----------|--------|--------|-------------|------|
| 100 | AllowSQL | Inbound | Allow | * | * | 1433 |
| 110 | AllowMySQL | Inbound | Allow | 0.0.0.0/0 | * | 3306 |

## Objectives

| Order | ID | Description | Trigger | Points |
|-------|-----|-------------|---------|--------|
| 1 | list_vms | Enumerate virtual machines in the environment | `az vm list` | 50 |
| 2 | find_nsgs | List Network Security Groups | `az network nsg list` | 75 |
| 3 | analyze_rules | Analyze NSG rules for misconfigurations | `az network nsg rule list` | 125 |
| 4 | identify_exposure | Identify critical services exposed to internet | `regex:nsg.*show\|rule.*list.*db` | 150 |

**Total Points:** 400

## Hints

### list_vms
- Level 1 (10 pts): "Start by understanding what VMs exist in the environment"
- Level 2 (25 pts): "az vm list -o table"

### find_nsgs
- Level 1 (15 pts): "Network Security Groups control traffic flow to VMs"
- Level 2 (35 pts): "az network nsg list --resource-group production-rg"

### analyze_rules
- Level 1 (20 pts): "Each NSG has rules - look for overly permissive sources"
- Level 2 (50 pts): "az network nsg rule list --nsg-name web-nsg --resource-group production-rg"

### identify_exposure
- Level 1 (25 pts): "Database ports exposed to the internet is a critical finding"
- Level 2 (50 pts): "Check the db-nsg for SQL and MySQL port rules"

## Solution Walkthrough

```bash
# Step 1: List VMs
az vm list -o table
az vm list --resource-group production-rg

# Step 2: List NSGs
az network nsg list --resource-group production-rg

# Step 3: Analyze web-nsg rules
az network nsg rule list --nsg-name web-nsg --resource-group production-rg

# Step 4: Check database NSG
az network nsg show --name db-nsg --resource-group production-rg
az network nsg rule list --nsg-name db-nsg --resource-group production-rg
```

## Debrief Topics

- NSG rule priority and evaluation order
- Service Tags vs IP ranges
- Application Security Groups (ASGs)
- Just-In-Time VM Access
- Azure Firewall vs NSGs
- Network Watcher for NSG diagnostics
- Azure Policy for NSG compliance

## Real-World References

- MITRE ATT&CK: T1190 (Exploit Public-Facing Application)
- CIS Azure Benchmark: 6.1-6.6 (Network Security Group recommendations)
- Azure Security Baseline: Network security

## Implementation Notes

Requires:
- az network nsg list/show commands
- az network nsg rule list command
- NSG resources in scenario YAML
- Enhanced compute package for NSG handling
