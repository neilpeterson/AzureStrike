package cli

import (
	"fmt"
	"strings"

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

AZURE CLI (mocked):
  az storage ...    Storage account operations
  az ad ...         Azure AD operations
  az vm ...         Virtual machine operations

STORAGE:
  cat <path>        View blob contents (e.g., cat backups/secrets.txt)

NETWORK:
  curl <url>        Make HTTP requests (IMDS supported)

GAME:
  objective         Show current objectives
  score             Show current score
  hint [obj_id]     Get a hint for an objective
  help [command]    Show help

SYSTEM:
  clear             Clear the terminal
  exit              Exit the game`,
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
