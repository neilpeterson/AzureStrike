package compute

// VirtualMachine represents an Azure VM
type VirtualMachine struct {
	ID                string            `json:"id" yaml:"id"`
	Name              string            `json:"name" yaml:"name"`
	ResourceGroup     string            `json:"resourceGroup" yaml:"resource_group"`
	Location          string            `json:"location" yaml:"location"`
	VMSize            string            `json:"vmSize,omitempty" yaml:"vm_size,omitempty"`
	PowerState        string            `json:"powerState" yaml:"power_state"`
	PublicIPAddress   string            `json:"publicIpAddress,omitempty" yaml:"public_ip,omitempty"`
	PrivateIPAddress  string            `json:"privateIpAddress,omitempty" yaml:"private_ip,omitempty"`
	OSType            string            `json:"osType,omitempty" yaml:"os_type,omitempty"`
	AdminUsername     string            `json:"adminUsername,omitempty" yaml:"admin_username,omitempty"`
	Tags              map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
	NetworkInterfaces []string          `json:"networkInterfaces,omitempty" yaml:"network_interfaces,omitempty"`
	Identity          *VMIdentity       `json:"identity,omitempty" yaml:"identity,omitempty"`
}

// VMIdentity represents managed identity configuration
type VMIdentity struct {
	Type        string `json:"type" yaml:"type"` // SystemAssigned, UserAssigned
	PrincipalID string `json:"principalId,omitempty" yaml:"principal_id,omitempty"`
	TenantID    string `json:"tenantId,omitempty" yaml:"tenant_id,omitempty"`
}

// NetworkSecurityGroup represents an Azure NSG
type NetworkSecurityGroup struct {
	ID            string        `json:"id" yaml:"id"`
	Name          string        `json:"name" yaml:"name"`
	ResourceGroup string        `json:"resourceGroup" yaml:"resource_group"`
	Location      string        `json:"location" yaml:"location"`
	Rules         []NSGRule     `json:"securityRules" yaml:"rules"`
}

// NSGRule represents an NSG security rule
type NSGRule struct {
	Name                     string `json:"name" yaml:"name"`
	Priority                 int    `json:"priority" yaml:"priority"`
	Direction                string `json:"direction" yaml:"direction"` // Inbound, Outbound
	Access                   string `json:"access" yaml:"access"`       // Allow, Deny
	Protocol                 string `json:"protocol" yaml:"protocol"`   // Tcp, Udp, *
	SourcePortRange          string `json:"sourcePortRange" yaml:"source_port_range"`
	DestinationPortRange     string `json:"destinationPortRange" yaml:"destination_port_range"`
	SourceAddressPrefix      string `json:"sourceAddressPrefix" yaml:"source_address_prefix"`
	DestinationAddressPrefix string `json:"destinationAddressPrefix" yaml:"destination_address_prefix"`
}

// Environment holds the mocked compute state for a scenario
type Environment struct {
	VirtualMachines       []VirtualMachine
	NetworkSecurityGroups []NetworkSecurityGroup
}

// NewEnvironment creates a new compute environment from scenario resources
func NewEnvironment(vms []VirtualMachine, nsgs []NetworkSecurityGroup) *Environment {
	return &Environment{
		VirtualMachines:       vms,
		NetworkSecurityGroups: nsgs,
	}
}

// ListVMs returns all virtual machines
func (e *Environment) ListVMs() []VirtualMachine {
	if e == nil || e.VirtualMachines == nil {
		return []VirtualMachine{}
	}
	return e.VirtualMachines
}

// GetVM returns a VM by name
func (e *Environment) GetVM(name string) *VirtualMachine {
	if e == nil {
		return nil
	}
	for i, vm := range e.VirtualMachines {
		if vm.Name == name || vm.ID == name {
			return &e.VirtualMachines[i]
		}
	}
	return nil
}

// ListVMsByResourceGroup returns VMs in a specific resource group
func (e *Environment) ListVMsByResourceGroup(rg string) []VirtualMachine {
	if e == nil {
		return []VirtualMachine{}
	}
	var result []VirtualMachine
	for _, vm := range e.VirtualMachines {
		if vm.ResourceGroup == rg {
			result = append(result, vm)
		}
	}
	return result
}

// ListNSGs returns all network security groups
func (e *Environment) ListNSGs() []NetworkSecurityGroup {
	if e == nil || e.NetworkSecurityGroups == nil {
		return []NetworkSecurityGroup{}
	}
	return e.NetworkSecurityGroups
}

// GetNSG returns an NSG by name
func (e *Environment) GetNSG(name string) *NetworkSecurityGroup {
	if e == nil {
		return nil
	}
	for i, nsg := range e.NetworkSecurityGroups {
		if nsg.Name == name || nsg.ID == name {
			return &e.NetworkSecurityGroups[i]
		}
	}
	return nil
}
