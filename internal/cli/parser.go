package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/azurestrike/azurestrike/internal/azure/storage"
	"github.com/azurestrike/azurestrike/internal/cli/az"
	"github.com/azurestrike/azurestrike/internal/game"
)

// Parser handles command parsing and routing
type Parser struct {
	gameState  *game.State
	azHandler  *az.Handler
}

// Result represents the result of executing a command
type Result struct {
	Output  string
	Success bool
	Error   error
}

// NewParser creates a new command parser
func NewParser(state *game.State) *Parser {
	return &Parser{
		gameState: state,
		azHandler: az.NewHandler(state.StorageEnv, state.EntraEnv, state.ComputeEnv),
	}
}

// Execute parses and executes a command
func (p *Parser) Execute(input string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Output: "", Success: true}
	}

	parts := parseCommand(input)
	if len(parts) == 0 {
		return Result{Output: "", Success: true}
	}

	cmd := parts[0]
	args := parts[1:]

	var result Result

	switch cmd {
	case "az":
		azResult := p.azHandler.Execute(args)
		result = Result{Output: azResult.Output, Success: azResult.Success, Error: azResult.Error}
	case "curl":
		result = p.handleCurl(args)
	case "scan":
		result = p.handleScan(args)
	case "cat":
		result = p.handleCat(args)
	case "help":
		result = p.handleHelp(args)
	case "clear":
		result = Result{Output: "\033[H\033[2J", Success: true}
	case "objectives", "objective", "obj":
		result = p.handleObjectives()
	case "score":
		result = p.handleScore()
	case "hint":
		result = p.handleHint(args)
	case "exit", "quit":
		result = Result{Output: "Use Ctrl+C to exit", Success: true}
	default:
		result = Result{
			Output:  fmt.Sprintf("Command not found: %s\nType 'help' for available commands.", cmd),
			Success: false,
		}
	}

	// Don't record help commands for objective completion
	if isHelpCommand(input) {
		return result
	}

	// Record the command and check for objective completion
	completedObjectives := p.gameState.RecordCommand(input, result.Output, result.Success)

	// Add notifications for completed objectives
	if len(completedObjectives) > 0 {
		for _, objID := range completedObjectives {
			obj := p.gameState.Scenario.GetObjective(objID)
			if obj != nil {
				notification := fmt.Sprintf("\n\n[+] OBJECTIVE COMPLETE: %s (+%d points)", obj.Description, obj.Points)
				result.Output += notification
			}
		}
	}

	return result
}

// isHelpCommand checks if a command is asking for help
func isHelpCommand(input string) bool {
	lower := strings.ToLower(input)
	return strings.Contains(lower, "--help") ||
		strings.Contains(lower, "-help") ||
		strings.HasSuffix(lower, " -h") ||
		strings.HasSuffix(lower, "\t-h") ||
		lower == "-h" ||
		lower == "help"
}

// parseCommand splits a command string respecting quotes
func parseCommand(input string) []string {
	var parts []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range input {
		switch {
		case r == '"' || r == '\'':
			if !inQuote {
				inQuote = true
				quoteChar = r
			} else if r == quoteChar {
				inQuote = false
				quoteChar = 0
			} else {
				current.WriteRune(r)
			}
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func (p *Parser) handleCurl(args []string) Result {
	if len(args) == 0 {
		return Result{Output: "curl: try 'curl --help' for more information", Success: false}
	}

	url := ""
	headers := make(map[string]string)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-H" && i+1 < len(args):
			i++
			headerParts := strings.SplitN(args[i], ":", 2)
			if len(headerParts) == 2 {
				headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
			}
		case arg == "-s", arg == "-S", arg == "-f":
			// Ignore common curl flags
		case !strings.HasPrefix(arg, "-"):
			url = arg
		}
	}

	// Handle IMDS requests
	if strings.Contains(url, "169.254.169.254") {
		return p.handleIMDS(url, headers)
	}

	// Handle Azure Blob Storage requests
	if strings.Contains(url, ".blob.core.windows.net") {
		return p.handleBlobStorage(url, headers)
	}

	return Result{
		Output:  fmt.Sprintf("curl: (6) Could not resolve host: %s", extractHost(url)),
		Success: false,
	}
}

func (p *Parser) handleIMDS(url string, headers map[string]string) Result {
	// Check for required Metadata header
	if headers["Metadata"] != "true" {
		return Result{
			Output:  "400 Bad Request: Required metadata header not specified",
			Success: false,
		}
	}

	// Mock IMDS responses
	if strings.Contains(url, "/metadata/identity/oauth2/token") {
		return Result{
			Output: `{
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik1uQ19WWmNBVG...",
  "client_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "expires_in": "28799",
  "expires_on": "1234567890",
  "ext_expires_in": "28799",
  "not_before": "1234567890",
  "resource": "https://management.azure.com/",
  "token_type": "Bearer"
}`,
			Success: true,
		}
	}

	if strings.Contains(url, "/metadata/instance") {
		return Result{
			Output: `{
  "compute": {
    "name": "victim-vm",
    "resourceGroupName": "production-rg",
    "subscriptionId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "vmId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "location": "eastus",
    "tags": "environment:production;owner:admin"
  }
}`,
			Success: true,
		}
	}

	return Result{
		Output:  `{"error": "Not found"}`,
		Success: false,
	}
}

func extractHost(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return url
}

// handleBlobStorage handles curl requests to Azure Blob Storage endpoints
func (p *Parser) handleBlobStorage(url string, headers map[string]string) Result {
	// Parse the URL: https://<account>.blob.core.windows.net/<container>/<blob>?params
	urlClean := strings.TrimPrefix(url, "http://")
	urlClean = strings.TrimPrefix(urlClean, "https://")

	// Split off query parameters
	urlPath := urlClean
	queryParams := ""
	if idx := strings.Index(urlClean, "?"); idx != -1 {
		urlPath = urlClean[:idx]
		queryParams = urlClean[idx+1:]
	}

	// Extract account name from host
	parts := strings.Split(urlPath, "/")
	if len(parts) == 0 {
		return Result{
			Output:  "curl: (6) Could not resolve host",
			Success: false,
		}
	}

	host := parts[0]
	accountName := strings.TrimSuffix(host, ".blob.core.windows.net")

	// Check if storage account exists in our environment
	account := p.gameState.StorageEnv.GetAccount(accountName)
	if account == nil {
		// Return realistic "account not found" response
		return Result{
			Output: fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Error>
  <Code>AccountNotFound</Code>
  <Message>The specified account %s does not exist.
RequestId:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Time:%s</Message>
</Error>`, accountName, time.Now().Format(time.RFC3339)),
			Success: false,
		}
	}

	// Just the account root - return account info
	if len(parts) == 1 || (len(parts) == 2 && parts[1] == "") {
		// Check for container list request: ?comp=list
		if strings.Contains(queryParams, "comp=list") {
			return p.listContainersXML(account)
		}
		// Just probing the account - return 400 (missing required params)
		return Result{
			Output: fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Error>
  <Code>InvalidQueryParameterValue</Code>
  <Message>Value for one of the query parameters specified in the request URI is invalid.
RequestId:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Time:%s</Message>
</Error>`, time.Now().Format(time.RFC3339)),
			Success: true, // Account exists, just missing params
		}
	}

	containerName := parts[1]
	container := p.gameState.StorageEnv.GetContainer(accountName, containerName)

	// Container doesn't exist
	if container == nil {
		return Result{
			Output: fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Error>
  <Code>ContainerNotFound</Code>
  <Message>The specified container does not exist.
RequestId:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Time:%s</Message>
</Error>`, time.Now().Format(time.RFC3339)),
			Success: false,
		}
	}

	// Container exists but is private
	if !container.IsPubliclyAccessible() {
		return Result{
			Output: fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Error>
  <Code>AuthenticationFailed</Code>
  <Message>Server failed to authenticate the request. Make sure the value of Authorization header is formed correctly.
RequestId:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Time:%s</Message>
</Error>`, time.Now().Format(time.RFC3339)),
			Success: false,
		}
	}

	// Check if listing blobs in container: ?restype=container&comp=list
	if strings.Contains(queryParams, "comp=list") {
		return p.listBlobsXML(account, container)
	}

	// Accessing a specific blob
	if len(parts) >= 3 {
		blobName := strings.Join(parts[2:], "/")
		blob := p.gameState.StorageEnv.GetBlob(accountName, containerName, blobName)
		if blob == nil {
			return Result{
				Output: fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Error>
  <Code>BlobNotFound</Code>
  <Message>The specified blob does not exist.
RequestId:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Time:%s</Message>
</Error>`, time.Now().Format(time.RFC3339)),
				Success: false,
			}
		}
		// Return blob content
		return Result{
			Output:  blob.Content,
			Success: true,
		}
	}

	// Just container path without query - return blob list
	return p.listBlobsXML(account, container)
}

// listContainersXML returns containers in Azure XML format
func (p *Parser) listContainersXML(account *storage.Account) Result {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://` + account.Name + `.blob.core.windows.net/">
  <Containers>
`)
	for _, c := range account.Containers {
		sb.WriteString(fmt.Sprintf(`    <Container>
      <Name>%s</Name>
      <Properties>
        <PublicAccess>%s</PublicAccess>
      </Properties>
    </Container>
`, c.Name, c.PublicAccess))
	}
	sb.WriteString(`  </Containers>
</EnumerationResults>`)

	return Result{Output: sb.String(), Success: true}
}

// listBlobsXML returns blobs in Azure XML format
func (p *Parser) listBlobsXML(account *storage.Account, container *storage.Container) Result {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://` + account.Name + `.blob.core.windows.net/" ContainerName="` + container.Name + `">
  <Blobs>
`)
	for _, b := range container.Blobs {
		sb.WriteString(fmt.Sprintf(`    <Blob>
      <Name>%s</Name>
      <Properties>
        <Content-Type>%s</Content-Type>
        <Content-Length>%d</Content-Length>
      </Properties>
    </Blob>
`, b.Name, b.ContentType, b.Size))
	}
	sb.WriteString(`  </Blobs>
</EnumerationResults>`)

	return Result{Output: sb.String(), Success: true}
}

// handleScan handles the scan command for storage account enumeration
func (p *Parser) handleScan(args []string) Result {
	if len(args) == 0 {
		return Result{
			Output: `scan - Azure Storage Account Scanner

Usage:
  scan storage <account-name>              Check if a specific account exists
  scan storage --wordlist <type>           Scan using a built-in wordlist
  scan storage --accounts <name1,name2>    Scan multiple account names

Wordlist types:
  common        Common storage account names (backup, storage, data, etc.)
  company       Company-related patterns (requires --company flag)

Examples:
  scan storage contoso2024
  scan storage --wordlist common
  scan storage --wordlist company --company contoso
  scan storage --accounts backup2024,storage2024,data2024

After finding an account, probe for containers:
  curl https://<account>.blob.core.windows.net/?comp=list

Then check for public containers:
  curl https://<account>.blob.core.windows.net/<container>?comp=list`,
			Success: true,
		}
	}

	if args[0] != "storage" {
		return Result{
			Output:  fmt.Sprintf("scan: unknown scan type '%s'. Currently supported: storage", args[0]),
			Success: false,
		}
	}

	// Parse scan arguments
	var accountsToScan []string
	var companyName string

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--wordlist":
			if i+1 >= len(args) {
				return Result{Output: "scan: --wordlist requires a type (common, company)", Success: false}
			}
			i++
			wordlistType := args[i]
			switch wordlistType {
			case "common":
				accountsToScan = append(accountsToScan, getCommonStorageNames()...)
			case "company":
				if companyName == "" {
					return Result{Output: "scan: --wordlist company requires --company <name>", Success: false}
				}
				accountsToScan = append(accountsToScan, getCompanyStorageNames(companyName)...)
			default:
				return Result{Output: fmt.Sprintf("scan: unknown wordlist type '%s'", wordlistType), Success: false}
			}
		case "--company":
			if i+1 >= len(args) {
				return Result{Output: "scan: --company requires a company name", Success: false}
			}
			i++
			companyName = strings.ToLower(args[i])
			// If we already queued company wordlist, regenerate
			accountsToScan = append(accountsToScan, getCompanyStorageNames(companyName)...)
		case "--accounts":
			if i+1 >= len(args) {
				return Result{Output: "scan: --accounts requires a comma-separated list", Success: false}
			}
			i++
			names := strings.Split(args[i], ",")
			for _, n := range names {
				accountsToScan = append(accountsToScan, strings.TrimSpace(n))
			}
		default:
			// Single account name to scan
			if !strings.HasPrefix(args[i], "--") {
				accountsToScan = append(accountsToScan, args[i])
			}
		}
	}

	if len(accountsToScan) == 0 {
		return Result{Output: "scan: no accounts specified. Use --wordlist, --accounts, or provide an account name", Success: false}
	}

	// Perform the scan
	var sb strings.Builder
	sb.WriteString("=== Azure Storage Account Scanner ===\n\n")
	sb.WriteString(fmt.Sprintf("Scanning %d account name(s)...\n\n", len(accountsToScan)))

	foundCount := 0
	for _, name := range accountsToScan {
		account := p.gameState.StorageEnv.GetAccount(name)
		if account != nil {
			foundCount++
			sb.WriteString(fmt.Sprintf("[+] FOUND: %s.blob.core.windows.net\n", name))
			sb.WriteString(fmt.Sprintf("    └─ Location: %s | Kind: %s\n", account.Location, account.Kind))
		} else {
			sb.WriteString(fmt.Sprintf("[-] Not found: %s\n", name))
		}
	}

	sb.WriteString(fmt.Sprintf("\nScan complete: %d/%d accounts found\n", foundCount, len(accountsToScan)))

	if foundCount > 0 {
		sb.WriteString("\nNext steps:\n")
		sb.WriteString("  1. List containers: curl https://<account>.blob.core.windows.net/?comp=list\n")
		sb.WriteString("  2. Check container access: curl https://<account>.blob.core.windows.net/<container>?comp=list\n")
		sb.WriteString("  3. Download public blobs: curl https://<account>.blob.core.windows.net/<container>/<blob>\n")
	}

	return Result{Output: sb.String(), Success: true}
}

// getCommonStorageNames returns common storage account name patterns
func getCommonStorageNames() []string {
	return []string{
		"backup", "backups", "storage", "data", "files",
		"uploads", "assets", "static", "media", "public",
		"private", "logs", "archive", "temp", "dev",
		"prod", "production", "staging", "test", "cdn",
	}
}

// getCompanyStorageNames generates storage names based on company name
func getCompanyStorageNames(company string) []string {
	suffixes := []string{
		"", "storage", "data", "backup", "backups",
		"files", "assets", "blob", "store", "2024", "2023",
		"prod", "dev", "staging",
	}

	years := []string{"", "2024", "2023", "2022"}

	var names []string
	for _, suffix := range suffixes {
		for _, year := range years {
			name := company
			if suffix != "" {
				name = company + suffix
			}
			if year != "" && suffix != year {
				name = name + year
			}
			// Azure storage accounts must be 3-24 chars, lowercase alphanumeric
			if len(name) >= 3 && len(name) <= 24 {
				names = append(names, name)
			}
		}
	}

	return names
}

func (p *Parser) handleCat(args []string) Result {
	if len(args) == 0 {
		return Result{Output: "cat: missing file operand", Success: false}
	}

	path := args[0]

	// Parse path: account/container/blob or just container/blob (uses first account)
	parts := strings.Split(path, "/")

	var accountName, containerName, blobName string

	if len(parts) == 3 {
		// Full path: account/container/blob
		accountName = parts[0]
		containerName = parts[1]
		blobName = parts[2]
	} else if len(parts) == 2 {
		// Short path: container/blob (find first account with this container)
		containerName = parts[0]
		blobName = parts[1]
		// Find account with this container
		for _, acc := range p.gameState.StorageEnv.ListAccounts() {
			for _, c := range acc.Containers {
				if c.Name == containerName {
					accountName = acc.Name
					break
				}
			}
			if accountName != "" {
				break
			}
		}
		if accountName == "" {
			return Result{
				Output:  fmt.Sprintf("cat: %s: No such file or directory", path),
				Success: false,
			}
		}
	} else {
		return Result{
			Output:  fmt.Sprintf("cat: %s: No such file or directory\nUsage: cat <account>/<container>/<blob> or cat <container>/<blob>", path),
			Success: false,
		}
	}

	// Check container access
	container := p.gameState.StorageEnv.GetContainer(accountName, containerName)
	if container == nil {
		return Result{
			Output:  fmt.Sprintf("cat: %s: No such file or directory", path),
			Success: false,
		}
	}

	if !container.IsPubliclyAccessible() {
		return Result{
			Output:  fmt.Sprintf("cat: %s: Permission denied", path),
			Success: false,
		}
	}

	// Get blob content
	blob := p.gameState.StorageEnv.GetBlob(accountName, containerName, blobName)
	if blob == nil {
		return Result{
			Output:  fmt.Sprintf("cat: %s: No such file or directory", path),
			Success: false,
		}
	}

	return Result{Output: blob.Content, Success: true}
}

func (p *Parser) handleHelp(args []string) Result {
	if len(args) > 0 {
		switch args[0] {
		case "az":
			return Result{
				Output: `Azure CLI Commands (mocked):
  az storage account list         List storage accounts
  az storage container list       List containers in a storage account
  az storage blob list            List blobs in a container
  az storage blob download        Download a blob
  az ad user list                 List Azure AD users
  az ad sp list                   List service principals`,
				Success: true,
			}
		}
	}

	return Result{
		Output: `AzureStrike - Available Commands:

RECONNAISSANCE:
  scan storage ...   Scan for Azure storage accounts
  curl <url>         Make HTTP requests (blob storage, IMDS)

AZURE CLI (mocked):
  az storage ...     Storage account operations
  az ad ...          Azure AD operations
  az vm ...          Virtual machine operations

STORAGE ACCESS:
  cat <path>         View blob contents (e.g., cat backups/secrets.txt)
  curl https://<account>.blob.core.windows.net/<container>/<blob>

GAME:
  objective          Show current objectives
  score              Show current score
  hint [obj_id]      Get a hint for an objective
  help [command]     Show help

SYSTEM:
  clear              Clear the terminal
  exit               Exit the game

Try 'scan' or 'scan storage' for scanner help.`,
		Success: true,
	}
}

func (p *Parser) handleObjectives() Result {
	var sb strings.Builder
	sb.WriteString("=== OBJECTIVES ===\n\n")

	pending := p.gameState.GetPendingObjectives()
	completed := p.gameState.GetCompletedObjectives()

	if len(pending) > 0 {
		sb.WriteString("PENDING:\n")
		for i, obj := range pending {
			sb.WriteString(fmt.Sprintf("  %d. [ ] %s (%d pts)\n", i+1, obj.Description, obj.Points))
		}
	}

	if len(completed) > 0 {
		sb.WriteString("\nCOMPLETED:\n")
		for _, obj := range completed {
			sb.WriteString(fmt.Sprintf("  [x] %s (+%d pts)\n", obj.Description, obj.Points))
		}
	}

	progress := p.gameState.GetProgress() * 100
	sb.WriteString(fmt.Sprintf("\nProgress: %.0f%%", progress))

	return Result{Output: sb.String(), Success: true}
}

func (p *Parser) handleScore() Result {
	score := p.gameState.Score
	var sb strings.Builder

	sb.WriteString("=== SCORE ===\n\n")
	sb.WriteString(fmt.Sprintf("Total Points: %d\n", score.Points))
	sb.WriteString(fmt.Sprintf("Max Possible: %d\n", p.gameState.Scenario.TotalPoints()))

	if len(score.Bonuses) > 0 {
		sb.WriteString("\nBonuses:\n")
		for _, bonus := range score.Bonuses {
			sb.WriteString(fmt.Sprintf("  [+] %s: +%d\n", bonus.Description, bonus.Points))
		}
	}

	achievements := score.GetUnlockedAchievements()
	if len(achievements) > 0 {
		sb.WriteString("\nAchievements:\n")
		for _, a := range achievements {
			sb.WriteString(fmt.Sprintf("  %s %s\n", a.Icon, a.Name))
		}
	}

	return Result{Output: sb.String(), Success: true}
}

func (p *Parser) handleHint(args []string) Result {
	if len(args) == 0 {
		// Show available hints
		pending := p.gameState.GetPendingObjectives()
		if len(pending) == 0 {
			return Result{Output: "No pending objectives!", Success: true}
		}

		var sb strings.Builder
		sb.WriteString("Available hints for objectives:\n")
		for _, obj := range pending {
			sb.WriteString(fmt.Sprintf("  hint %s\n", obj.ID))
		}
		return Result{Output: sb.String(), Success: true}
	}

	objID := args[0]
	currentLevel := p.gameState.HintsUsed[objID]
	nextLevel := currentLevel + 1

	hint := p.gameState.UseHint(objID, nextLevel)
	if hint == nil {
		return Result{
			Output:  fmt.Sprintf("No hint available for '%s' at level %d", objID, nextLevel),
			Success: false,
		}
	}

	return Result{
		Output:  fmt.Sprintf("HINT (Level %d, -%d pts):\n%s", hint.Level, hint.PointCost, hint.Text),
		Success: true,
	}
}
