package az

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	OutputJSON  OutputFormat = "json"
	OutputTable OutputFormat = "table"
	OutputTSV   OutputFormat = "tsv"
	OutputNone  OutputFormat = "none"
)

// parseOutputFormat extracts the output format from args
func parseOutputFormat(args []string) OutputFormat {
	for i, arg := range args {
		if (arg == "-o" || arg == "--output") && i+1 < len(args) {
			switch strings.ToLower(args[i+1]) {
			case "table":
				return OutputTable
			case "tsv":
				return OutputTSV
			case "none":
				return OutputNone
			default:
				return OutputJSON
			}
		}
		if strings.HasPrefix(arg, "-o=") {
			return OutputFormat(strings.TrimPrefix(arg, "-o="))
		}
		if strings.HasPrefix(arg, "--output=") {
			return OutputFormat(strings.TrimPrefix(arg, "--output="))
		}
	}
	return OutputJSON
}

// formatOutput formats data based on the output format
func formatOutput(data interface{}, format OutputFormat, columns []string) string {
	switch format {
	case OutputTable:
		return formatTable(data, columns)
	case OutputTSV:
		return formatTSV(data, columns)
	case OutputNone:
		return ""
	default:
		output, _ := json.MarshalIndent(data, "", "  ")
		return string(output)
	}
}

// formatTable formats data as an ASCII table
func formatTable(data interface{}, columns []string) string {
	// Convert to slice of maps for uniform handling
	rows := toMapSlice(data)
	if len(rows) == 0 {
		return ""
	}

	// If no columns specified, use all keys from first row
	if len(columns) == 0 {
		for k := range rows[0] {
			columns = append(columns, k)
		}
	}

	// Calculate column widths
	widths := make(map[string]int)
	for _, col := range columns {
		widths[col] = len(formatColumnHeader(col))
	}
	for _, row := range rows {
		for _, col := range columns {
			val := fmt.Sprintf("%v", row[col])
			if len(val) > widths[col] {
				widths[col] = len(val)
			}
		}
	}

	// Build table
	var sb strings.Builder

	// Header
	for i, col := range columns {
		header := formatColumnHeader(col)
		sb.WriteString(fmt.Sprintf("%-*s", widths[col], header))
		if i < len(columns)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")

	// Separator
	for i, col := range columns {
		sb.WriteString(strings.Repeat("-", widths[col]))
		if i < len(columns)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")

	// Rows
	for _, row := range rows {
		for i, col := range columns {
			val := fmt.Sprintf("%v", row[col])
			sb.WriteString(fmt.Sprintf("%-*s", widths[col], val))
			if i < len(columns)-1 {
				sb.WriteString("  ")
			}
		}
		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

// formatTSV formats data as tab-separated values
func formatTSV(data interface{}, columns []string) string {
	rows := toMapSlice(data)
	if len(rows) == 0 {
		return ""
	}

	if len(columns) == 0 {
		for k := range rows[0] {
			columns = append(columns, k)
		}
	}

	var sb strings.Builder
	for _, row := range rows {
		vals := make([]string, len(columns))
		for i, col := range columns {
			vals[i] = fmt.Sprintf("%v", row[col])
		}
		sb.WriteString(strings.Join(vals, "\t"))
		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

// formatColumnHeader converts camelCase/snake_case to the Azure CLI header format
func formatColumnHeader(s string) string {
	// Handle common Azure CLI column names - match exact Azure CLI output
	headerMap := map[string]string{
		"name":                   "Name",
		"resourceGroup":          "ResourceGroup",
		"resource_group":         "ResourceGroup",
		"location":               "Location",
		"id":                     "Id",
		"displayName":            "DisplayName",
		"display_name":           "DisplayName",
		"userPrincipalName":      "UserPrincipalName",
		"appId":                  "AppId",
		"app_id":                 "AppId",
		"powerState":             "PowerState",
		"power_state":            "PowerState",
		"publicIpAddress":        "PublicIps",
		"public_ip":              "PublicIps",
		"privateIpAddress":       "PrivateIps",
		"private_ip":             "PrivateIps",
		"vmSize":                 "VmSize",
		"vm_size":                "VmSize",
		"publicAccess":           "PublicAccess",
		"public_access":          "PublicAccess",
		"kind":                   "Kind",
		"sku":                    "Sku",
		"contentLength":          "ContentLength",
		"contentType":            "ContentType",
		"accessTier":             "AccessTier",
		"allowBlobPublicAccess":  "AllowBlobPublicAccess",
		"creationTime":           "CreationTime",
		"enableHttpsTrafficOnly": "EnableHttpsTrafficOnly",
		"minimumTlsVersion":      "MinimumTlsVersion",
		"provisioningState":      "ProvisioningState",
		"primaryLocation":        "PrimaryLocation",
		"statusOfPrimary":        "StatusOfPrimary",
	}

	if mapped, ok := headerMap[s]; ok {
		return mapped
	}
	return s
}

// toMapSlice converts various data types to []map[string]interface{}
func toMapSlice(data interface{}) []map[string]interface{} {
	var result []map[string]interface{}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if m := toMap(item); m != nil {
				result = append(result, m)
			}
		}
	case reflect.Struct, reflect.Map:
		if m := toMap(data); m != nil {
			result = append(result, m)
		}
	}

	return result
}

// toMap converts a struct or map to map[string]interface{}
func toMap(data interface{}) map[string]interface{} {
	// First try JSON marshaling (handles struct tags)
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}

	return m
}
