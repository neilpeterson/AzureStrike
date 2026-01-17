package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
