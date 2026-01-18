package az

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/azurestrike/azurestrike/internal/azure/compute"
	"github.com/azurestrike/azurestrike/internal/azure/entra"
	"github.com/azurestrike/azurestrike/internal/azure/storage"
)

// TokenInfo contains information about an extracted token
type TokenInfo struct {
	Token       string
	Resource    string
	PrincipalID string
}

// TokenGetter is a function that returns the latest extracted token
type TokenGetter func() *TokenInfo

// Handler processes az CLI commands
type Handler struct {
	storageEnv  *storage.Environment
	entraEnv    *entra.Environment
	computeEnv  *compute.Environment
	tokenGetter TokenGetter
}

// Result represents command execution result
type Result struct {
	Output  string
	Success bool
	Error   error
}

// NewHandler creates a new az CLI handler
func NewHandler(storageEnv *storage.Environment, entraEnv *entra.Environment, computeEnv *compute.Environment, tokenGetter TokenGetter) *Handler {
	return &Handler{
		storageEnv:  storageEnv,
		entraEnv:    entraEnv,
		computeEnv:  computeEnv,
		tokenGetter: tokenGetter,
	}
}

// Execute runs an az CLI command
func (h *Handler) Execute(args []string) Result {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		return Result{Output: h.showHelp(), Success: true}
	}

	switch args[0] {
	case "storage":
		return h.handleStorage(args[1:])
	case "ad":
		return h.handleAD(args[1:])
	case "vm":
		return h.handleVM(args[1:])
	case "account":
		return h.handleAccount(args[1:])
	case "login":
		return Result{
			Output:  "You are already logged in as attacker@yourcompany.com",
			Success: true,
		}
	default:
		return Result{
			Output:  fmt.Sprintf("az: '%s' is not an az command. See 'az --help'.", args[0]),
			Success: false,
		}
	}
}

func (h *Handler) showHelp() string {
	return `
Group
    az : Manage Azure resources.

Subgroups:
    account        : Manage Azure subscription information.
    ad             : Manage Microsoft Entra ID (formerly Azure Active Directory).
    storage        : Manage Azure Cloud Storage resources.
    vm             : Manage Linux or Windows virtual machines.

Global Arguments
    --help -h      : Show this help message and exit.
    --output -o    : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                     Default: json.
    --query        : JMESPath query string. See http://jmespath.org/ for more information.
    --verbose      : Increase logging verbosity. Use --debug for full debug logs.

To search for commands, use: az find "search term"
For more specific help, run: az [command] --help`
}

func (h *Handler) handleStorage(args []string) Result {
	if len(args) == 0 || hasHelpFlag(args) {
		return Result{Output: h.storageHelp(), Success: true}
	}

	switch args[0] {
	case "account":
		return h.handleStorageAccount(args[1:])
	case "container":
		return h.handleStorageContainer(args[1:])
	case "blob":
		return h.handleStorageBlob(args[1:])
	default:
		return Result{
			Output:  fmt.Sprintf("az storage: '%s' is not a known command", args[0]),
			Success: false,
		}
	}
}

func (h *Handler) storageHelp() string {
	return `
Group
    az storage : Manage Azure Cloud Storage resources.

Subgroups:
    account            : Manage storage accounts.
    blob               : Manage object storage for unstructured data (blobs).
    container          : Manage blob storage containers.

Global Arguments
    --help -h          : Show this help message and exit.
    --output -o        : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                         Default: json.
    --query            : JMESPath query string. See http://jmespath.org/ for more information.
    --verbose          : Increase logging verbosity. Use --debug for full debug logs.`
}

func (h *Handler) handleStorageAccount(args []string) Result {
	if hasHelpFlag(args) {
		return Result{Output: h.storageAccountHelp(), Success: true}
	}

	// Check for unrecognized arguments
	validFlags := []string{"--name", "-n", "--resource-group", "-g"}
	if err := checkUnrecognizedArgs(args, validFlags, "az storage account"); err != nil {
		return *err
	}

	// Determine command (list or show)
	command := "list"
	for _, arg := range args {
		if arg == "list" || arg == "show" {
			command = arg
			break
		}
	}

	format := parseOutputFormat(args)

	if command == "list" || len(args) == 0 {
		accounts := h.storageEnv.ListAccounts()
		// Match Azure CLI table columns
		columns := []string{"name", "resourceGroup", "location", "sku", "kind", "accessTier", "enableHttpsTrafficOnly", "provisioningState", "statusOfPrimary"}
		return Result{Output: formatOutput(accounts, format, columns), Success: true}
	}

	if command == "show" {
		name := getArgValue(args, "--name", "-n")
		if name == "" {
			return Result{Output: "az storage account show: error: the following arguments are required: --name/-n", Success: false}
		}
		acc := h.storageEnv.GetAccount(name)
		if acc == nil {
			return Result{
				Output:  fmt.Sprintf("The storage account '%s' was not found", name),
				Success: false,
			}
		}
		columns := []string{"name", "resourceGroup", "location", "kind", "sku"}
		return Result{Output: formatOutput(acc, format, columns), Success: true}
	}

	return Result{Output: fmt.Sprintf("az storage account: '%s' is not in the 'az storage account' command group.", args[0]), Success: false}
}

func (h *Handler) storageAccountHelp() string {
	return `
Group
    az storage account : Manage storage accounts.

Commands:
    list               : List storage accounts.
    show               : Show storage account properties.

Arguments
    --name -n          : The storage account name.

Global Arguments
    --help -h          : Show this help message and exit.
    --output -o        : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                         Default: json.
    --query            : JMESPath query string. See http://jmespath.org/ for more information.
    --verbose          : Increase logging verbosity. Use --debug for full debug logs.

Examples
    List all storage accounts in the subscription:
        az storage account list

    Show properties of a storage account:
        az storage account show --name mystorageaccount`
}

func (h *Handler) handleStorageContainer(args []string) Result {
	if len(args) == 0 || hasHelpFlag(args) {
		return Result{Output: h.storageContainerHelp(), Success: true}
	}

	// Check for unrecognized arguments
	validFlags := []string{"--account-name", "--account-key", "--auth-mode"}
	if err := checkUnrecognizedArgs(args, validFlags, "az storage container"); err != nil {
		return *err
	}

	accountName := getArgValue(args, "--account-name", "")
	format := parseOutputFormat(args)

	// Find command
	command := ""
	for _, arg := range args {
		if arg == "list" {
			command = arg
			break
		}
	}

	if command == "list" {
		if accountName == "" {
			return Result{
				Output:  "az storage container list: error: the following arguments are required: --account-name",
				Success: false,
			}
		}
		acc := h.storageEnv.GetAccount(accountName)
		if acc == nil {
			return Result{
				Output:  fmt.Sprintf("The storage account '%s' was not found", accountName),
				Success: false,
			}
		}

		// Return container list
		type containerInfo struct {
			Name         string `json:"name"`
			PublicAccess string `json:"publicAccess"`
		}
		var containers []containerInfo
		for _, c := range acc.Containers {
			containers = append(containers, containerInfo{
				Name:         c.Name,
				PublicAccess: c.PublicAccess,
			})
		}
		columns := []string{"name", "publicAccess"}
		return Result{Output: formatOutput(containers, format, columns), Success: true}
	}

	return Result{Output: fmt.Sprintf("az storage container: '%s' is not in the 'az storage container' command group.", args[0]), Success: false}
}

func (h *Handler) storageContainerHelp() string {
	return `
Group
    az storage container : Manage blob storage containers.

Commands:
    list               : List containers in a storage account.

Arguments
    --account-name     : Storage account name. Required.
    --auth-mode        : The mode in which to run the command. "login" mode will directly use your
                         login credentials. "key" mode will attempt to query for an account key.

Global Arguments
    --help -h          : Show this help message and exit.
    --output -o        : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                         Default: json.
    --query            : JMESPath query string. See http://jmespath.org/ for more information.
    --verbose          : Increase logging verbosity. Use --debug for full debug logs.

Examples
    List containers in a storage account:
        az storage container list --account-name mystorageaccount`
}

func (h *Handler) handleStorageBlob(args []string) Result {
	if len(args) == 0 || hasHelpFlag(args) {
		return Result{Output: h.storageBlobHelp(), Success: true}
	}

	// Check for unrecognized arguments
	validFlags := []string{"--account-name", "--container-name", "-c", "--name", "-n", "--file", "-f", "--auth-mode"}
	if err := checkUnrecognizedArgs(args, validFlags, "az storage blob"); err != nil {
		return *err
	}

	accountName := getArgValue(args, "--account-name", "")
	containerName := getArgValue(args, "--container-name", "-c")
	format := parseOutputFormat(args)

	// Find command
	command := ""
	for _, arg := range args {
		if arg == "list" || arg == "download" {
			command = arg
			break
		}
	}

	if command == "list" {
		if accountName == "" {
			return Result{
				Output:  "az storage blob list: error: the following arguments are required: --account-name",
				Success: false,
			}
		}
		if containerName == "" {
			return Result{
				Output:  "az storage blob list: error: the following arguments are required: --container-name/-c",
				Success: false,
			}
		}

		container := h.storageEnv.GetContainer(accountName, containerName)
		if container == nil {
			return Result{
				Output:  fmt.Sprintf("Container '%s' not found in account '%s'", containerName, accountName),
				Success: false,
			}
		}

		// Check access
		if !container.IsPubliclyAccessible() {
			return Result{
				Output:  "This request is not authorized to perform this operation.\nRequestId:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\nTime:2024-01-01T00:00:00.0000000Z",
				Success: false,
			}
		}

		type blobInfo struct {
			Name        string `json:"name"`
			Size        int64  `json:"contentLength"`
			ContentType string `json:"contentType"`
		}
		var blobs []blobInfo
		for _, b := range container.Blobs {
			blobs = append(blobs, blobInfo{
				Name:        b.Name,
				Size:        b.Size,
				ContentType: b.ContentType,
			})
		}
		columns := []string{"name", "contentLength", "contentType"}
		return Result{Output: formatOutput(blobs, format, columns), Success: true}
	}

	if command == "download" {
		if accountName == "" {
			return Result{
				Output:  "az storage blob download: error: the following arguments are required: --account-name",
				Success: false,
			}
		}
		if containerName == "" {
			return Result{
				Output:  "az storage blob download: error: the following arguments are required: --container-name/-c",
				Success: false,
			}
		}
		blobName := getArgValue(args, "--name", "-n")
		if blobName == "" {
			return Result{
				Output:  "az storage blob download: error: the following arguments are required: --name/-n",
				Success: false,
			}
		}

		blob := h.storageEnv.GetBlob(accountName, containerName, blobName)
		if blob == nil {
			return Result{
				Output:  "BlobNotFound: The specified blob does not exist.\nRequestId:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\nTime:2024-01-01T00:00:00.0000000Z\nErrorCode: BlobNotFound",
				Success: false,
			}
		}

		// Check container access
		container := h.storageEnv.GetContainer(accountName, containerName)
		if container != nil && !container.IsPubliclyAccessible() {
			return Result{
				Output:  "This request is not authorized to perform this operation.\nRequestId:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\nTime:2024-01-01T00:00:00.0000000Z",
				Success: false,
			}
		}

		// Return blob content
		outputFile := getArgValue(args, "--file", "-f")
		if outputFile != "" {
			return Result{
				Output:  fmt.Sprintf("Finished[#############################################################################]  100.0000%%\n{\n  \"content\": null,\n  \"name\": \"%s\"\n}\n\nContent:\n%s", blobName, blob.Content),
				Success: true,
			}
		}
		return Result{Output: blob.Content, Success: true}
	}

	return Result{Output: fmt.Sprintf("az storage blob: '%s' is not in the 'az storage blob' command group.", args[0]), Success: false}
}

func (h *Handler) storageBlobHelp() string {
	return `
Group
    az storage blob : Manage object storage for unstructured data (blobs).

Commands:
    download           : Download a blob to a file path.
    list               : List blobs in a given container.

Arguments
    --account-name     : Storage account name. Required.
    --container-name -c: The container name. Required.
    --name -n          : The blob name. Required for download.
    --file -f          : Path of file to write out to.
    --auth-mode        : The mode in which to run the command. "login" mode will directly use your
                         login credentials. "key" mode will attempt to query for an account key.

Global Arguments
    --help -h          : Show this help message and exit.
    --output -o        : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                         Default: json.
    --query            : JMESPath query string. See http://jmespath.org/ for more information.
    --verbose          : Increase logging verbosity. Use --debug for full debug logs.

Examples
    List all blobs in a container:
        az storage blob list --account-name mystorageaccount --container-name mycontainer

    Download a blob to a file:
        az storage blob download --account-name mystorageaccount --container-name mycontainer \\
            --name myblob.txt --file /path/to/file.txt`
}

func (h *Handler) handleAD(args []string) Result {
	if len(args) == 0 || hasHelpFlag(args) {
		return Result{Output: h.adHelp(), Success: true}
	}

	switch args[0] {
	case "user":
		return h.handleADUser(args[1:])
	case "sp":
		return h.handleADServicePrincipal(args[1:])
	default:
		return Result{
			Output:  fmt.Sprintf("az ad: '%s' is not a known command", args[0]),
			Success: false,
		}
	}
}

func (h *Handler) adHelp() string {
	return `
Group
    az ad : Manage Microsoft Entra ID (formerly Azure Active Directory).

Subgroups:
    sp                 : Manage service principals for automation authentication.
    user               : Manage Azure Active Directory users.

Global Arguments
    --help -h          : Show this help message and exit.
    --output -o        : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                         Default: json.
    --query            : JMESPath query string. See http://jmespath.org/ for more information.
    --verbose          : Increase logging verbosity. Use --debug for full debug logs.`
}

func (h *Handler) handleADUser(args []string) Result {
	if hasHelpFlag(args) {
		return Result{Output: h.adUserHelp(), Success: true}
	}

	// Check for unrecognized arguments
	validFlags := []string{"--id", "--upn-or-object-id", "--filter", "--all"}
	if err := checkUnrecognizedArgs(args, validFlags, "az ad user"); err != nil {
		return *err
	}

	format := parseOutputFormat(args)

	// Find command
	command := "list"
	for _, arg := range args {
		if arg == "list" || arg == "show" {
			command = arg
			break
		}
	}

	if command == "list" {
		users := h.entraEnv.ListUsers()
		if len(users) == 0 {
			return Result{Output: "[]", Success: true}
		}
		columns := []string{"displayName", "id", "userPrincipalName"}
		return Result{Output: formatOutput(users, format, columns), Success: true}
	}

	if command == "show" {
		upn := getArgValue(args, "--id", "")
		if upn == "" {
			upn = getArgValue(args, "--upn-or-object-id", "")
		}
		if upn == "" {
			return Result{Output: "az ad user show: error: the following arguments are required: --id", Success: false}
		}
		user := h.entraEnv.GetUser(upn)
		if user == nil {
			return Result{
				Output:  fmt.Sprintf("Resource '%s' does not exist or one of its queried reference-property objects are not present.", upn),
				Success: false,
			}
		}
		columns := []string{"displayName", "id", "userPrincipalName"}
		return Result{Output: formatOutput(user, format, columns), Success: true}
	}

	return Result{Output: fmt.Sprintf("az ad user: '%s' is not in the 'az ad user' command group.", args[0]), Success: false}
}

func (h *Handler) adUserHelp() string {
	return `
Group
    az ad user : Manage Azure Active Directory users.

Commands:
    list               : List Azure Active Directory users.

Global Arguments
    --help -h          : Show this help message and exit.
    --output -o        : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                         Default: json.
    --query            : JMESPath query string. See http://jmespath.org/ for more information.
    --verbose          : Increase logging verbosity. Use --debug for full debug logs.

Examples
    List all users in the directory:
        az ad user list`
}

func (h *Handler) handleADServicePrincipal(args []string) Result {
	if hasHelpFlag(args) {
		return Result{Output: h.adSpHelp(), Success: true}
	}

	// Check for unrecognized arguments
	validFlags := []string{"--id", "--filter", "--all", "--display-name"}
	if err := checkUnrecognizedArgs(args, validFlags, "az ad sp"); err != nil {
		return *err
	}

	format := parseOutputFormat(args)

	// Find command
	command := "list"
	for _, arg := range args {
		if arg == "list" || arg == "show" {
			command = arg
			break
		}
	}

	if command == "list" {
		sps := h.entraEnv.ListServicePrincipals()
		if len(sps) == 0 {
			return Result{Output: "[]", Success: true}
		}
		columns := []string{"appId", "displayName", "id"}
		return Result{Output: formatOutput(sps, format, columns), Success: true}
	}

	if command == "show" {
		id := getArgValue(args, "--id", "")
		if id == "" {
			return Result{Output: "az ad sp show: error: the following arguments are required: --id", Success: false}
		}
		sp := h.entraEnv.GetServicePrincipal(id)
		if sp == nil {
			return Result{
				Output:  fmt.Sprintf("Service principal '%s' does not exist or one of its queried reference-property objects are not present.", id),
				Success: false,
			}
		}
		columns := []string{"appId", "displayName", "id"}
		return Result{Output: formatOutput(sp, format, columns), Success: true}
	}

	return Result{Output: fmt.Sprintf("az ad sp: '%s' is not in the 'az ad sp' command group.", args[0]), Success: false}
}

func (h *Handler) adSpHelp() string {
	return `
Group
    az ad sp : Manage service principals for automation authentication.

Commands:
    list               : List service principals.

Global Arguments
    --help -h          : Show this help message and exit.
    --output -o        : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                         Default: json.
    --query            : JMESPath query string. See http://jmespath.org/ for more information.
    --verbose          : Increase logging verbosity. Use --debug for full debug logs.

Examples
    List all service principals:
        az ad sp list`
}

func (h *Handler) handleVM(args []string) Result {
	if hasHelpFlag(args) {
		return Result{Output: h.vmHelp(), Success: true}
	}

	// Check for unrecognized arguments
	validFlags := []string{"--resource-group", "-g", "--name", "-n", "--show-details", "-d"}
	if err := checkUnrecognizedArgs(args, validFlags, "az vm"); err != nil {
		return *err
	}

	format := parseOutputFormat(args)

	// Find command
	command := "list"
	for _, arg := range args {
		if arg == "list" || arg == "show" {
			command = arg
			break
		}
	}

	if command == "list" {
		// Check for resource group filter
		rg := getArgValue(args, "--resource-group", "-g")
		var vms []compute.VirtualMachine
		if rg != "" {
			vms = h.computeEnv.ListVMsByResourceGroup(rg)
		} else {
			vms = h.computeEnv.ListVMs()
		}
		if len(vms) == 0 {
			return Result{Output: "[]", Success: true}
		}
		columns := []string{"name", "resourceGroup", "powerState", "publicIpAddress", "location"}
		return Result{Output: formatOutput(vms, format, columns), Success: true}
	}

	if command == "show" {
		name := getArgValue(args, "--name", "-n")
		if name == "" {
			return Result{Output: "az vm show: error: the following arguments are required: --name/-n", Success: false}
		}
		vm := h.computeEnv.GetVM(name)
		if vm == nil {
			return Result{
				Output:  fmt.Sprintf("(ResourceNotFound) The Resource 'Microsoft.Compute/virtualMachines/%s' under resource group 'unknown' was not found.", name),
				Success: false,
			}
		}
		columns := []string{"name", "resourceGroup", "powerState", "publicIpAddress", "location"}
		return Result{Output: formatOutput(vm, format, columns), Success: true}
	}

	return Result{Output: fmt.Sprintf("az vm: '%s' is not in the 'az vm' command group.", args[0]), Success: false}
}

func (h *Handler) vmHelp() string {
	return `
Group
    az vm : Manage Linux or Windows virtual machines.

Commands:
    list               : List details of Virtual Machines.

Global Arguments
    --help -h          : Show this help message and exit.
    --output -o        : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                         Default: json.
    --query            : JMESPath query string. See http://jmespath.org/ for more information.
    --verbose          : Increase logging verbosity. Use --debug for full debug logs.

Examples
    List all VMs in the subscription:
        az vm list

    List all VMs in a resource group:
        az vm list -g MyResourceGroup`
}

func (h *Handler) handleAccount(args []string) Result {
	if hasHelpFlag(args) {
		return Result{Output: h.accountHelp(), Success: true}
	}

	if len(args) == 0 || args[0] == "show" {
		account := map[string]interface{}{
			"name":      "Contoso Production",
			"id":        "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"tenantId":  "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy",
			"state":     "Enabled",
			"isDefault": true,
		}
		output, _ := json.MarshalIndent(account, "", "  ")
		return Result{Output: string(output), Success: true}
	}
	if args[0] == "list" {
		accounts := []map[string]interface{}{
			{
				"name":      "Contoso Production",
				"id":        "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				"isDefault": true,
			},
		}
		output, _ := json.MarshalIndent(accounts, "", "  ")
		return Result{Output: string(output), Success: true}
	}
	if args[0] == "get-access-token" {
		return h.handleGetAccessToken(args[1:])
	}
	return Result{Output: "az account: unknown command", Success: false}
}

func (h *Handler) handleGetAccessToken(args []string) Result {
	if hasHelpFlag(args) {
		return Result{Output: h.getAccessTokenHelp(), Success: true}
	}

	// Check if we have a token getter and if there's a token available
	if h.tokenGetter == nil {
		return Result{
			Output:  "ERROR: Please run 'az login' to setup account.",
			Success: false,
		}
	}

	token := h.tokenGetter()
	if token == nil {
		return Result{
			Output:  "ERROR: Please run 'az login' to setup account.",
			Success: false,
		}
	}

	// Parse resource argument if provided
	resource := getArgValue(args, "--resource", "")
	if resource == "" {
		resource = "https://management.azure.com/"
	}

	// Build response matching Azure CLI format
	now := time.Now()
	expiresOn := now.Add(time.Hour * 24)

	response := map[string]interface{}{
		"accessToken":  token.Token,
		"expiresOn":    expiresOn.Format("2006-01-02 15:04:05.000000"),
		"expires_on":   expiresOn.Unix(),
		"subscription": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		"tenant":       "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy",
		"tokenType":    "Bearer",
		"resource":     resource,
	}

	output, _ := json.MarshalIndent(response, "", "  ")
	return Result{Output: string(output), Success: true}
}

func (h *Handler) getAccessTokenHelp() string {
	return `
Command
    az account get-access-token : Get a token for utilities to access Azure.

Arguments
    --resource             : Azure resource endpoints in AAD v1.0.
    --resource-type        : Type of resource. Allowed values: aad-graph, arm, ms-graph.
    --scope                : Space-separated scopes.
    --tenant -t            : Tenant ID.

Global Arguments
    --help -h              : Show this help message and exit.
    --output -o            : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                             Default: json.
    --query                : JMESPath query string. See http://jmespath.org/ for more information.

Examples
    Get access token for the current subscription:
        az account get-access-token

    Get access token for a specific resource:
        az account get-access-token --resource https://storage.azure.com/`
}

func (h *Handler) accountHelp() string {
	return `
Group
    az account : Manage Azure subscription information.

Commands:
    get-access-token   : Get a token for utilities to access Azure.
    list               : Get a list of subscriptions for the logged in account.
    show               : Get the details of a subscription.

Global Arguments
    --help -h          : Show this help message and exit.
    --output -o        : Output format.  Allowed values: json, jsonc, none, table, tsv, yaml.
                         Default: json.
    --query            : JMESPath query string. See http://jmespath.org/ for more information.
    --verbose          : Increase logging verbosity. Use --debug for full debug logs.

Examples
    Show current subscription:
        az account show

    List all subscriptions:
        az account list

    Get access token:
        az account get-access-token`
}

// hasHelpFlag checks if args contain --help, -h, or -help as the first arg
func hasHelpFlag(args []string) bool {
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h" || args[0] == "-help") {
		return true
	}
	return false
}

// getArgValue extracts a value for a flag from args
func getArgValue(args []string, flag string, shortFlag string) string {
	for i, arg := range args {
		if (arg == flag || (shortFlag != "" && arg == shortFlag)) && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(arg, flag+"=") {
			return strings.TrimPrefix(arg, flag+"=")
		}
		if shortFlag != "" && strings.HasPrefix(arg, shortFlag+"=") {
			return strings.TrimPrefix(arg, shortFlag+"=")
		}
	}
	return ""
}

// globalFlags are flags accepted by all commands
var globalFlags = []string{
	"--help", "-h", "-help",
	"--output", "-o",
	"--query",
	"--verbose",
	"--debug",
}

// checkUnrecognizedArgs checks for unrecognized arguments and returns an error message
func checkUnrecognizedArgs(args []string, validFlags []string, commandPath string) *Result {
	allValidFlags := append(globalFlags, validFlags...)

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			continue // Not a flag, skip
		}

		// Check if it's a valid flag
		isValid := false
		for _, validFlag := range allValidFlags {
			if arg == validFlag || strings.HasPrefix(arg, validFlag+"=") {
				isValid = true
				break
			}
		}

		if !isValid {
			return &Result{
				Output:  fmt.Sprintf("%s: error: unrecognized arguments: %s", commandPath, arg),
				Success: false,
			}
		}
	}
	return nil
}
