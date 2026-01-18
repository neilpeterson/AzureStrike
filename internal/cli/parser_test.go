package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/azurestrike/azurestrike/internal/azure/compute"
	"github.com/azurestrike/azurestrike/internal/azure/storage"
	"github.com/azurestrike/azurestrike/internal/game"
	"github.com/azurestrike/azurestrike/internal/scenario"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple command",
			input:    "az storage list",
			expected: []string{"az", "storage", "list"},
		},
		{
			name:     "command with double quotes",
			input:    `az storage blob download --name "file name.txt"`,
			expected: []string{"az", "storage", "blob", "download", "--name", "file name.txt"},
		},
		{
			name:     "command with single quotes",
			input:    `az storage blob download --name 'file name.txt'`,
			expected: []string{"az", "storage", "blob", "download", "--name", "file name.txt"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: nil,
		},
		{
			name:     "multiple spaces between args",
			input:    "az    storage   list",
			expected: []string{"az", "storage", "list"},
		},
		{
			name:     "quotes containing quotes",
			input:    `echo "it's working"`,
			expected: []string{"echo", "it's working"},
		},
		{
			name:     "nested quotes",
			input:    `az query --query "[?name=='test']"`,
			expected: []string{"az", "query", "--query", "[?name=='test']"},
		},
		{
			name:     "equals sign in argument",
			input:    "curl -H Authorization=Bearer",
			expected: []string{"curl", "-H", "Authorization=Bearer"},
		},
		{
			name:     "url argument",
			input:    "curl https://example.com/path?query=value",
			expected: []string{"curl", "https://example.com/path?query=value"},
		},
		{
			name:     "empty quotes",
			input:    `az test --arg ""`,
			expected: []string{"az", "test", "--arg"},
		},
		{
			name:     "single word",
			input:    "help",
			expected: []string{"help"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommand(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsHelpCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "help keyword",
			input:    "help",
			expected: true,
		},
		{
			name:     "help with uppercase",
			input:    "HELP",
			expected: true,
		},
		{
			name:     "--help flag",
			input:    "az storage --help",
			expected: true,
		},
		{
			name:     "-help flag",
			input:    "az storage -help",
			expected: true,
		},
		{
			name:     "-h at end",
			input:    "az storage -h",
			expected: true,
		},
		{
			name:     "-h alone",
			input:    "-h",
			expected: true,
		},
		{
			name:     "-h with tab",
			input:    "az\t-h",
			expected: true,
		},
		{
			name:     "normal command",
			input:    "az storage account list",
			expected: false,
		},
		{
			name:     "-h in middle of command",
			input:    "az -h storage list",
			expected: false, // isHelpCommand only checks HasSuffix for -h
		},
		{
			name:     "helpful as part of word",
			input:    "echo helpful",
			expected: false, // "helpful" doesn't contain "--help"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHelpCommand(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseCommandEdgeCases(t *testing.T) {
	t.Run("unclosed quote at end", func(t *testing.T) {
		// Current implementation treats unclosed quote content as part of last token
		result := parseCommand(`az test "unclosed`)
		assert.Contains(t, result, "az")
		assert.Contains(t, result, "test")
	})

	t.Run("mixed quote types", func(t *testing.T) {
		result := parseCommand(`az test "double" 'single'`)
		assert.Equal(t, []string{"az", "test", "double", "single"}, result)
	})

	t.Run("backslash content in quotes", func(t *testing.T) {
		// Backslashes in raw string literals are literal backslashes
		result := parseCommand(`echo "path/to/file"`)
		assert.Equal(t, []string{"echo", "path/to/file"}, result)
	})
}

// Helper function to create a test scenario with VMs and managed identities
func createIMDSTestScenario() *scenario.Scenario {
	return &scenario.Scenario{
		ID:         "imds-test",
		Name:       "IMDS Test Scenario",
		Difficulty: "intermediate",
		Resources: scenario.Resources{
			VirtualMachines: []compute.VirtualMachine{
				{
					Name:          "test-vm",
					ResourceGroup: "test-rg",
					Location:      "eastus",
					PowerState:    "VM running",
					Tags: map[string]string{
						"environment": "test",
						"owner":       "test-user",
					},
					Identity: &compute.VMIdentity{
						Type:        "SystemAssigned",
						PrincipalID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
						TenantID:    "12345678-abcd-ef12-3456-7890abcdef12",
					},
				},
			},
			StorageAccounts: []storage.Account{
				{
					Name:          "testbackups",
					ResourceGroup: "test-rg",
					Location:      "eastus",
					Kind:          "StorageV2",
					Containers: []storage.Container{
						{
							Name:         "private-data",
							PublicAccess: "none",
							Blobs: []storage.Blob{
								{
									Name:        "secret.txt",
									Content:     "This is secret data!",
									ContentType: "text/plain",
									Size:        20,
								},
							},
						},
					},
				},
			},
		},
		Objectives: []scenario.Objective{
			{ID: "test-obj", Description: "Test objective", Trigger: "test", Points: 100},
		},
	}
}

func TestIMDSHandlerNoIdentity(t *testing.T) {
	// Create a scenario with no managed identity
	sc := &scenario.Scenario{
		ID:         "no-identity",
		Name:       "No Identity Test",
		Difficulty: "beginner",
		Resources: scenario.Resources{
			VirtualMachines: []compute.VirtualMachine{
				{
					Name:          "vm-no-identity",
					ResourceGroup: "test-rg",
					Location:      "eastus",
					Identity:      nil,
				},
			},
		},
		Objectives: []scenario.Objective{
			{ID: "obj1", Description: "Test", Trigger: "test", Points: 100},
		},
	}

	state := game.NewState(sc)
	parser := NewParser(state)

	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/"`)

	assert.False(t, result.Success)
	assert.Contains(t, result.Output, "IdentityNotFound")
}

func TestIMDSHandlerRequiresMetadataHeader(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	result := parser.Execute(`curl "http://169.254.169.254/metadata/instance?api-version=2021-02-01"`)

	assert.False(t, result.Success)
	// Real Azure returns JSON error format
	assert.Contains(t, result.Output, `"error"`)
	assert.Contains(t, result.Output, "Required metadata header not specified")
}

func TestIMDSHandlerInstanceMetadata(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/instance?api-version=2021-02-01"`)

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, `"name": "test-vm"`)
	assert.Contains(t, result.Output, `"resourceGroupName": "test-rg"`)
	assert.Contains(t, result.Output, `"location": "eastus"`)
	assert.Contains(t, result.Output, "environment:test")
}

func TestIMDSHandlerTokenExtraction(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://storage.azure.com/"`)

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "access_token")
	assert.Contains(t, result.Output, "Bearer")
	assert.Contains(t, result.Output, "https://storage.azure.com/")

	// Token should be stored in game state
	assert.NotEmpty(t, state.ExtractedTokens)
}

func TestIMDSHandlerTokenUsageForStorage(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	// First, extract a token
	tokenResult := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://storage.azure.com/"`)
	assert.True(t, tokenResult.Success)

	// Extract the token from the response
	var token string
	for t := range state.ExtractedTokens {
		token = t
		break
	}
	assert.NotEmpty(t, token, "Token should be stored in state")

	// Try to access private storage without token - should fail
	failResult := parser.Execute(`curl "https://testbackups.blob.core.windows.net/private-data/secret.txt"`)
	assert.False(t, failResult.Success)
	assert.Contains(t, failResult.Output, "AuthenticationFailed")

	// Try with token but without x-ms-version - should fail
	noVersionResult := parser.Execute(`curl -H "Authorization: Bearer ` + token + `" "https://testbackups.blob.core.windows.net/private-data/secret.txt"`)
	assert.False(t, noVersionResult.Success)
	assert.Contains(t, noVersionResult.Output, "MissingRequiredHeader")
	assert.Contains(t, noVersionResult.Output, "x-ms-version")

	// Try with both token and x-ms-version - should succeed
	successResult := parser.Execute(`curl -H "Authorization: Bearer ` + token + `" -H "x-ms-version: 2020-04-08" "https://testbackups.blob.core.windows.net/private-data/secret.txt"`)
	assert.True(t, successResult.Success)
	assert.Contains(t, successResult.Output, "This is secret data!")
}

func TestIMDSHandlerListCategories(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/"`)

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "instance")
	assert.Contains(t, result.Output, "identity")
}

func TestIMDSHandlerIdentityInfo(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/info?api-version=2021-02-01"`)

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "tenantId")
	assert.Contains(t, result.Output, "12345678-abcd-ef12-3456-7890abcdef12")
}

func TestIMDSHandlerScheduledEvents(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/scheduledevents?api-version=2020-07-01"`)

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "DocumentIncarnation")
	assert.Contains(t, result.Output, "Events")
}

func TestIMDSHandlerNotFound(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/nonexistent?api-version=2021-02-01"`)

	assert.False(t, result.Success)
	assert.Contains(t, result.Output, "NotFound")
}

func TestIMDSHandlerRequiresApiVersion(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	// Endpoints other than /metadata root require api-version
	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/instance"`)

	assert.False(t, result.Success)
	assert.Contains(t, result.Output, "api-version")
	assert.Contains(t, result.Output, "newest-versions")
}

func TestIMDSHandlerRequiresResourceParam(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	// Token endpoint requires resource parameter
	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01"`)

	assert.False(t, result.Success)
	assert.Contains(t, result.Output, "resource")
	assert.Contains(t, result.Output, "missing")
}

func TestGenerateMockAccessToken(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	// Use reflection or just verify through the result
	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/"`)

	assert.True(t, result.Success)

	// Token should have JWT format (header.payload.signature)
	var token string
	for t := range state.ExtractedTokens {
		token = t
		break
	}

	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3, "JWT should have 3 parts separated by dots")
}

func TestIMDSHandlerTokenResponseFormat(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://storage.azure.com/"`)

	assert.True(t, result.Success)
	// Check for all required fields in real Azure token response
	assert.Contains(t, result.Output, "access_token")
	assert.Contains(t, result.Output, "refresh_token")
	assert.Contains(t, result.Output, "expires_in")
	assert.Contains(t, result.Output, "expires_on")
	assert.Contains(t, result.Output, "not_before")
	assert.Contains(t, result.Output, "resource")
	assert.Contains(t, result.Output, "token_type")
	assert.Contains(t, result.Output, "Bearer")
}

func TestIMDSHandlerInstanceMetadataRichFields(t *testing.T) {
	sc := createIMDSTestScenario()
	state := game.NewState(sc)
	parser := NewParser(state)

	result := parser.Execute(`curl -H "Metadata: true" "http://169.254.169.254/metadata/instance?api-version=2021-02-01"`)

	assert.True(t, result.Success)
	// Check for rich compute fields
	assert.Contains(t, result.Output, "vmId")
	assert.Contains(t, result.Output, "resourceId")
	assert.Contains(t, result.Output, "osProfile")
	assert.Contains(t, result.Output, "securityProfile")
	assert.Contains(t, result.Output, "storageProfile")
	assert.Contains(t, result.Output, "platformFaultDomain")
	assert.Contains(t, result.Output, "publisher")
	assert.Contains(t, result.Output, "offer")
	assert.Contains(t, result.Output, "sku")
	// Check for network section
	assert.Contains(t, result.Output, "network")
	assert.Contains(t, result.Output, "ipAddress")
}
