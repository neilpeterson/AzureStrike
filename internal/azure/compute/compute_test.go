package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvironment(t *testing.T) {
	vms := []VirtualMachine{
		{ID: "vm1", Name: "web-vm"},
	}
	nsgs := []NetworkSecurityGroup{
		{ID: "nsg1", Name: "web-nsg"},
	}

	env := NewEnvironment(vms, nsgs)

	assert.NotNil(t, env)
	assert.Len(t, env.VirtualMachines, 1)
	assert.Len(t, env.NetworkSecurityGroups, 1)
}

func TestNewEnvironmentWithNil(t *testing.T) {
	env := NewEnvironment(nil, nil)

	assert.NotNil(t, env)
	assert.Nil(t, env.VirtualMachines)
	assert.Nil(t, env.NetworkSecurityGroups)
}

func TestListVMs(t *testing.T) {
	vms := []VirtualMachine{
		{ID: "1", Name: "vm1"},
		{ID: "2", Name: "vm2"},
		{ID: "3", Name: "vm3"},
	}
	env := NewEnvironment(vms, nil)

	result := env.ListVMs()

	assert.Len(t, result, 3)
	assert.Equal(t, "vm1", result[0].Name)
}

func TestListVMsEmpty(t *testing.T) {
	env := NewEnvironment(nil, nil)

	result := env.ListVMs()

	assert.Empty(t, result)
}

func TestListVMsNilEnvironment(t *testing.T) {
	var env *Environment

	result := env.ListVMs()

	assert.Empty(t, result)
}

func TestGetVM(t *testing.T) {
	vms := []VirtualMachine{
		{
			ID:            "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/web-vm",
			Name:          "web-vm",
			ResourceGroup: "production-rg",
			Location:      "eastus",
		},
	}
	env := NewEnvironment(vms, nil)

	t.Run("by name", func(t *testing.T) {
		vm := env.GetVM("web-vm")
		assert.NotNil(t, vm)
		assert.Equal(t, "production-rg", vm.ResourceGroup)
	})

	t.Run("by full ID", func(t *testing.T) {
		vm := env.GetVM("/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/web-vm")
		assert.NotNil(t, vm)
		assert.Equal(t, "web-vm", vm.Name)
	})

	t.Run("not found", func(t *testing.T) {
		vm := env.GetVM("nonexistent")
		assert.Nil(t, vm)
	})

	t.Run("nil environment", func(t *testing.T) {
		var nilEnv *Environment
		vm := nilEnv.GetVM("web-vm")
		assert.Nil(t, vm)
	})
}

func TestListVMsByResourceGroup(t *testing.T) {
	vms := []VirtualMachine{
		{Name: "vm1", ResourceGroup: "prod-rg"},
		{Name: "vm2", ResourceGroup: "prod-rg"},
		{Name: "vm3", ResourceGroup: "dev-rg"},
	}
	env := NewEnvironment(vms, nil)

	t.Run("matching resource group", func(t *testing.T) {
		result := env.ListVMsByResourceGroup("prod-rg")
		assert.Len(t, result, 2)
	})

	t.Run("single match", func(t *testing.T) {
		result := env.ListVMsByResourceGroup("dev-rg")
		assert.Len(t, result, 1)
		assert.Equal(t, "vm3", result[0].Name)
	})

	t.Run("no matches", func(t *testing.T) {
		result := env.ListVMsByResourceGroup("nonexistent-rg")
		assert.Empty(t, result)
	})

	t.Run("nil environment", func(t *testing.T) {
		var nilEnv *Environment
		result := nilEnv.ListVMsByResourceGroup("prod-rg")
		assert.Empty(t, result)
	})
}

func TestListNSGs(t *testing.T) {
	nsgs := []NetworkSecurityGroup{
		{ID: "1", Name: "nsg1"},
		{ID: "2", Name: "nsg2"},
	}
	env := NewEnvironment(nil, nsgs)

	result := env.ListNSGs()

	assert.Len(t, result, 2)
}

func TestListNSGsEmpty(t *testing.T) {
	env := NewEnvironment(nil, nil)

	result := env.ListNSGs()

	assert.Empty(t, result)
}

func TestListNSGsNilEnvironment(t *testing.T) {
	var env *Environment

	result := env.ListNSGs()

	assert.Empty(t, result)
}

func TestGetNSG(t *testing.T) {
	nsgs := []NetworkSecurityGroup{
		{
			ID:            "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/web-nsg",
			Name:          "web-nsg",
			ResourceGroup: "network-rg",
			Location:      "eastus",
		},
	}
	env := NewEnvironment(nil, nsgs)

	t.Run("by name", func(t *testing.T) {
		nsg := env.GetNSG("web-nsg")
		assert.NotNil(t, nsg)
		assert.Equal(t, "network-rg", nsg.ResourceGroup)
	})

	t.Run("by full ID", func(t *testing.T) {
		nsg := env.GetNSG("/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/web-nsg")
		assert.NotNil(t, nsg)
		assert.Equal(t, "web-nsg", nsg.Name)
	})

	t.Run("not found", func(t *testing.T) {
		nsg := env.GetNSG("nonexistent")
		assert.Nil(t, nsg)
	})

	t.Run("nil environment", func(t *testing.T) {
		var nilEnv *Environment
		nsg := nilEnv.GetNSG("web-nsg")
		assert.Nil(t, nsg)
	})
}

func TestVirtualMachineFields(t *testing.T) {
	vm := VirtualMachine{
		ID:               "vm-id",
		Name:             "production-vm",
		ResourceGroup:    "prod-rg",
		Location:         "westus2",
		VMSize:           "Standard_D2s_v3",
		PowerState:       "running",
		PublicIPAddress:  "52.1.2.3",
		PrivateIPAddress: "10.0.0.4",
		OSType:           "Linux",
		AdminUsername:    "azureuser",
		Tags:             map[string]string{"environment": "production"},
		NetworkInterfaces: []string{
			"/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/nic1",
		},
		Identity: &VMIdentity{
			Type:        "SystemAssigned",
			PrincipalID: "principal-uuid",
			TenantID:    "tenant-uuid",
		},
	}

	assert.Equal(t, "vm-id", vm.ID)
	assert.Equal(t, "production-vm", vm.Name)
	assert.Equal(t, "prod-rg", vm.ResourceGroup)
	assert.Equal(t, "westus2", vm.Location)
	assert.Equal(t, "Standard_D2s_v3", vm.VMSize)
	assert.Equal(t, "running", vm.PowerState)
	assert.Equal(t, "52.1.2.3", vm.PublicIPAddress)
	assert.Equal(t, "10.0.0.4", vm.PrivateIPAddress)
	assert.Equal(t, "Linux", vm.OSType)
	assert.Equal(t, "azureuser", vm.AdminUsername)
	assert.Equal(t, "production", vm.Tags["environment"])
	assert.Len(t, vm.NetworkInterfaces, 1)
	assert.NotNil(t, vm.Identity)
	assert.Equal(t, "SystemAssigned", vm.Identity.Type)
}

func TestNetworkSecurityGroupFields(t *testing.T) {
	nsg := NetworkSecurityGroup{
		ID:            "nsg-id",
		Name:          "web-nsg",
		ResourceGroup: "network-rg",
		Location:      "eastus",
		Rules: []NSGRule{
			{
				Name:                     "AllowHTTPS",
				Priority:                 100,
				Direction:                "Inbound",
				Access:                   "Allow",
				Protocol:                 "Tcp",
				SourcePortRange:          "*",
				DestinationPortRange:     "443",
				SourceAddressPrefix:      "*",
				DestinationAddressPrefix: "*",
			},
		},
	}

	assert.Equal(t, "nsg-id", nsg.ID)
	assert.Equal(t, "web-nsg", nsg.Name)
	assert.Equal(t, "network-rg", nsg.ResourceGroup)
	assert.Equal(t, "eastus", nsg.Location)
	assert.Len(t, nsg.Rules, 1)

	rule := nsg.Rules[0]
	assert.Equal(t, "AllowHTTPS", rule.Name)
	assert.Equal(t, 100, rule.Priority)
	assert.Equal(t, "Inbound", rule.Direction)
	assert.Equal(t, "Allow", rule.Access)
	assert.Equal(t, "Tcp", rule.Protocol)
	assert.Equal(t, "*", rule.SourcePortRange)
	assert.Equal(t, "443", rule.DestinationPortRange)
}

func TestVMIdentityFields(t *testing.T) {
	identity := VMIdentity{
		Type:        "UserAssigned",
		PrincipalID: "principal-123",
		TenantID:    "tenant-456",
	}

	assert.Equal(t, "UserAssigned", identity.Type)
	assert.Equal(t, "principal-123", identity.PrincipalID)
	assert.Equal(t, "tenant-456", identity.TenantID)
}

func TestNSGRuleFields(t *testing.T) {
	rule := NSGRule{
		Name:                     "DenyAll",
		Priority:                 4096,
		Direction:                "Outbound",
		Access:                   "Deny",
		Protocol:                 "*",
		SourcePortRange:          "*",
		DestinationPortRange:     "*",
		SourceAddressPrefix:      "10.0.0.0/8",
		DestinationAddressPrefix: "Internet",
	}

	assert.Equal(t, "DenyAll", rule.Name)
	assert.Equal(t, 4096, rule.Priority)
	assert.Equal(t, "Outbound", rule.Direction)
	assert.Equal(t, "Deny", rule.Access)
	assert.Equal(t, "*", rule.Protocol)
	assert.Equal(t, "*", rule.SourcePortRange)
	assert.Equal(t, "*", rule.DestinationPortRange)
	assert.Equal(t, "10.0.0.0/8", rule.SourceAddressPrefix)
	assert.Equal(t, "Internet", rule.DestinationAddressPrefix)
}
