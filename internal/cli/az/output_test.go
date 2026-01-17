package az

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOutputFormat(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected OutputFormat
	}{
		{
			name:     "default is json",
			args:     []string{"storage", "account", "list"},
			expected: OutputJSON,
		},
		{
			name:     "explicit json with -o",
			args:     []string{"storage", "list", "-o", "json"},
			expected: OutputJSON,
		},
		{
			name:     "table format with -o",
			args:     []string{"storage", "list", "-o", "table"},
			expected: OutputTable,
		},
		{
			name:     "tsv format with -o",
			args:     []string{"storage", "list", "-o", "tsv"},
			expected: OutputTSV,
		},
		{
			name:     "none format with -o",
			args:     []string{"storage", "list", "-o", "none"},
			expected: OutputNone,
		},
		{
			name:     "table format with --output",
			args:     []string{"storage", "list", "--output", "table"},
			expected: OutputTable,
		},
		{
			name:     "table format with -o= syntax",
			args:     []string{"storage", "list", "-o=table"},
			expected: OutputTable,
		},
		{
			name:     "table format with --output= syntax",
			args:     []string{"storage", "list", "--output=tsv"},
			expected: OutputTSV,
		},
		{
			name:     "case insensitive",
			args:     []string{"storage", "list", "-o", "TABLE"},
			expected: OutputTable,
		},
		{
			name:     "unknown format defaults to json",
			args:     []string{"storage", "list", "-o", "unknown"},
			expected: OutputJSON,
		},
		{
			name:     "-o at end without value",
			args:     []string{"storage", "list", "-o"},
			expected: OutputJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOutputFormat(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatOutput(t *testing.T) {
	data := []map[string]interface{}{
		{"name": "account1", "location": "eastus"},
		{"name": "account2", "location": "westus"},
	}

	t.Run("json format", func(t *testing.T) {
		output := formatOutput(data, OutputJSON, nil)
		assert.Contains(t, output, "account1")
		assert.Contains(t, output, "account2")
		assert.Contains(t, output, "eastus")
	})

	t.Run("none format returns empty", func(t *testing.T) {
		output := formatOutput(data, OutputNone, nil)
		assert.Empty(t, output)
	})

	t.Run("table format with columns", func(t *testing.T) {
		output := formatOutput(data, OutputTable, []string{"name", "location"})
		assert.Contains(t, output, "account1")
		assert.Contains(t, output, "eastus")
	})

	t.Run("tsv format with columns", func(t *testing.T) {
		output := formatOutput(data, OutputTSV, []string{"name", "location"})
		assert.Contains(t, output, "account1")
		assert.Contains(t, output, "\t")
	})
}

func TestFormatTable(t *testing.T) {
	t.Run("with data and columns", func(t *testing.T) {
		data := []map[string]interface{}{
			{"name": "vm1", "location": "eastus"},
			{"name": "vm2", "location": "westus2"},
		}
		output := formatTable(data, []string{"name", "location"})

		assert.Contains(t, output, "Name")
		assert.Contains(t, output, "Location")
		assert.Contains(t, output, "vm1")
		assert.Contains(t, output, "eastus")
		assert.Contains(t, output, "---")
	})

	t.Run("empty data returns empty", func(t *testing.T) {
		output := formatTable([]map[string]interface{}{}, []string{"name"})
		assert.Empty(t, output)
	})

	t.Run("auto columns from data", func(t *testing.T) {
		data := []map[string]interface{}{
			{"name": "test"},
		}
		output := formatTable(data, nil)
		assert.NotEmpty(t, output)
	})

	t.Run("column width adjusts to content", func(t *testing.T) {
		data := []map[string]interface{}{
			{"name": "short"},
			{"name": "verylongname"},
		}
		output := formatTable(data, []string{"name"})
		assert.Contains(t, output, "verylongname")
	})
}

func TestFormatTSV(t *testing.T) {
	t.Run("with data", func(t *testing.T) {
		data := []map[string]interface{}{
			{"col1": "a", "col2": "b"},
			{"col1": "c", "col2": "d"},
		}
		output := formatTSV(data, []string{"col1", "col2"})

		assert.Contains(t, output, "a\tb")
		assert.Contains(t, output, "c\td")
	})

	t.Run("empty data returns empty", func(t *testing.T) {
		output := formatTSV([]map[string]interface{}{}, []string{"col1"})
		assert.Empty(t, output)
	})

	t.Run("auto columns from data", func(t *testing.T) {
		data := []map[string]interface{}{
			{"name": "test", "value": "123"},
		}
		output := formatTSV(data, nil)
		assert.NotEmpty(t, output)
	})
}

func TestFormatColumnHeader(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"name", "Name"},
		{"location", "Location"},
		{"resourceGroup", "ResourceGroup"},
		{"resource_group", "ResourceGroup"},
		{"displayName", "DisplayName"},
		{"display_name", "DisplayName"},
		{"userPrincipalName", "UserPrincipalName"},
		{"appId", "AppId"},
		{"app_id", "AppId"},
		{"powerState", "PowerState"},
		{"power_state", "PowerState"},
		{"publicIpAddress", "PublicIps"},
		{"privateIpAddress", "PrivateIps"},
		{"vmSize", "VmSize"},
		{"unknown_column", "unknown_column"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatColumnHeader(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToMapSlice(t *testing.T) {
	t.Run("slice of structs", func(t *testing.T) {
		type Item struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}
		data := []Item{
			{Name: "item1", Value: 1},
			{Name: "item2", Value: 2},
		}
		result := toMapSlice(data)

		assert.Len(t, result, 2)
		assert.Equal(t, "item1", result[0]["name"])
		assert.Equal(t, float64(1), result[0]["value"])
	})

	t.Run("slice of maps", func(t *testing.T) {
		data := []map[string]interface{}{
			{"key": "value1"},
			{"key": "value2"},
		}
		result := toMapSlice(data)

		assert.Len(t, result, 2)
	})

	t.Run("single struct", func(t *testing.T) {
		type Item struct {
			Name string `json:"name"`
		}
		data := Item{Name: "single"}
		result := toMapSlice(data)

		assert.Len(t, result, 1)
		assert.Equal(t, "single", result[0]["name"])
	})

	t.Run("single map", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		result := toMapSlice(data)

		assert.Len(t, result, 1)
	})

	t.Run("pointer to slice", func(t *testing.T) {
		data := &[]map[string]interface{}{
			{"key": "value"},
		}
		result := toMapSlice(data)

		assert.Len(t, result, 1)
	})
}

func TestToMap(t *testing.T) {
	t.Run("struct to map", func(t *testing.T) {
		type TestStruct struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}
		data := TestStruct{Name: "test", Value: 42}
		result := toMap(data)

		assert.NotNil(t, result)
		assert.Equal(t, "test", result["name"])
		assert.Equal(t, float64(42), result["value"])
	})

	t.Run("map to map", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		result := toMap(data)

		assert.NotNil(t, result)
		assert.Equal(t, "value", result["key"])
	})
}
