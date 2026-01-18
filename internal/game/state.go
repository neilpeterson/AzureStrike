package game

import (
	"regexp"
	"strings"
	"time"

	"github.com/azurestrike/azurestrike/internal/azure/compute"
	"github.com/azurestrike/azurestrike/internal/azure/entra"
	"github.com/azurestrike/azurestrike/internal/azure/storage"
	"github.com/azurestrike/azurestrike/internal/scenario"
)

// State represents the current game state
type State struct {
	Scenario            *scenario.Scenario
	StartTime           time.Time
	CompletedObjectives map[string]time.Time
	CommandHistory      []CommandRecord
	HintsUsed           map[string]int // objective_id -> highest hint level used
	StorageEnv          *storage.Environment
	EntraEnv            *entra.Environment
	ComputeEnv          *compute.Environment
	Score               *Score
	Status              GameStatus
	ExtractedTokens     map[string]*ExtractedToken // token -> metadata
}

// ExtractedToken represents a token extracted via IMDS or other means
type ExtractedToken struct {
	Token       string
	Resource    string // The resource the token is scoped to (e.g., https://storage.azure.com/)
	PrincipalID string
	ExtractedAt time.Time
}

// CommandRecord stores information about executed commands
type CommandRecord struct {
	Command   string
	Output    string
	Timestamp time.Time
	Success   bool
}

// GameStatus represents the current status of the game
type GameStatus int

const (
	StatusPlaying GameStatus = iota
	StatusCompleted
	StatusFailed
)

// NewState creates a new game state from a scenario
func NewState(sc *scenario.Scenario) *State {
	state := &State{
		Scenario:            sc,
		StartTime:           time.Now(),
		CompletedObjectives: make(map[string]time.Time),
		CommandHistory:      []CommandRecord{},
		HintsUsed:           make(map[string]int),
		StorageEnv:          storage.NewEnvironment(),
		EntraEnv:            entra.NewEnvironment(sc.Resources.Users, sc.Resources.ServicePrincipals, sc.Resources.Groups),
		ComputeEnv:          compute.NewEnvironment(sc.Resources.VirtualMachines, sc.Resources.NetworkSecurityGroups),
		Score:               NewScore(),
		Status:              StatusPlaying,
		ExtractedTokens:     make(map[string]*ExtractedToken),
	}

	// Initialize storage environment from scenario resources
	for i := range sc.Resources.StorageAccounts {
		state.StorageEnv.AddAccount(&sc.Resources.StorageAccounts[i])
	}

	return state
}

// StoreToken records a token that was extracted via IMDS or other methods
func (s *State) StoreToken(token, resource, principalID string) {
	s.ExtractedTokens[token] = &ExtractedToken{
		Token:       token,
		Resource:    resource,
		PrincipalID: principalID,
		ExtractedAt: time.Now(),
	}
}

// ValidateToken checks if a token is valid for accessing a given resource
func (s *State) ValidateToken(token string, requiredResource string) bool {
	extracted, exists := s.ExtractedTokens[token]
	if !exists {
		return false
	}

	// Check if the token is scoped for the required resource
	// Storage tokens can be scoped to:
	// - https://storage.azure.com/
	// - https://management.azure.com/ (ARM can access storage)
	// - https://<account>.blob.core.windows.net/ (specific account)
	switch requiredResource {
	case "storage":
		return extracted.Resource == "https://storage.azure.com/" ||
			extracted.Resource == "https://management.azure.com/" ||
			strings.Contains(extracted.Resource, ".blob.core.windows.net")
	case "management":
		return extracted.Resource == "https://management.azure.com/"
	default:
		return extracted.Resource == requiredResource
	}
}

// GetTokenPrincipal returns the principal ID for a given token
func (s *State) GetTokenPrincipal(token string) string {
	if extracted, exists := s.ExtractedTokens[token]; exists {
		return extracted.PrincipalID
	}
	return ""
}

// GetLatestToken returns the most recently extracted token, if any
func (s *State) GetLatestToken() *ExtractedToken {
	var latest *ExtractedToken
	for _, token := range s.ExtractedTokens {
		if latest == nil || token.ExtractedAt.After(latest.ExtractedAt) {
			latest = token
		}
	}
	return latest
}

// RecordCommand adds a command to the history and checks for objective completion
func (s *State) RecordCommand(cmd, output string, success bool) []string {
	s.CommandHistory = append(s.CommandHistory, CommandRecord{
		Command:   cmd,
		Output:    output,
		Timestamp: time.Now(),
		Success:   success,
	})

	// Check for objective completion
	return s.checkObjectives(cmd)
}

// checkObjectives tests if a command triggers any objectives
func (s *State) checkObjectives(cmd string) []string {
	var completed []string

	for _, obj := range s.Scenario.Objectives {
		// Skip already completed objectives
		if _, done := s.CompletedObjectives[obj.ID]; done {
			continue
		}

		if s.matchesTrigger(cmd, obj.Trigger) {
			s.CompleteObjective(obj.ID)
			completed = append(completed, obj.ID)
		}
	}

	return completed
}

// matchesTrigger checks if a command matches an objective trigger pattern
func (s *State) matchesTrigger(cmd, trigger string) bool {
	// Triggers can be:
	// - Simple substring match: "blob list"
	// - Regex pattern: "regex:az storage blob (list|download)"
	// - Command with specific args: "az storage blob download --name secrets.txt"

	if strings.HasPrefix(trigger, "regex:") {
		pattern := strings.TrimPrefix(trigger, "regex:")
		matched, _ := regexp.MatchString(pattern, cmd)
		return matched
	}

	// Simple substring match (case-insensitive)
	return strings.Contains(strings.ToLower(cmd), strings.ToLower(trigger))
}

// CompleteObjective marks an objective as completed
func (s *State) CompleteObjective(id string) {
	if _, done := s.CompletedObjectives[id]; done {
		return
	}

	s.CompletedObjectives[id] = time.Now()

	// Award points
	if obj := s.Scenario.GetObjective(id); obj != nil {
		// Reduce points if hints were used
		points := obj.Points
		if hintLevel, ok := s.HintsUsed[id]; ok {
			for _, hint := range s.Scenario.Hints {
				if hint.ObjectiveID == id && hint.Level <= hintLevel {
					points -= hint.PointCost
				}
			}
		}
		if points < 0 {
			points = 0
		}
		s.Score.AddPoints(points, "Completed: "+obj.Description)
	}

	// Check if all objectives are complete
	if s.AllObjectivesComplete() {
		s.Status = StatusCompleted
		s.Score.AddBonus("scenario_complete", 100, "Scenario completed!")
	}
}

// UseHint records that a hint was used and deducts points
func (s *State) UseHint(objectiveID string, level int) *scenario.Hint {
	// Find the hint
	var hint *scenario.Hint
	for i := range s.Scenario.Hints {
		h := &s.Scenario.Hints[i]
		if h.ObjectiveID == objectiveID && h.Level == level {
			hint = h
			break
		}
	}

	if hint == nil {
		return nil
	}

	// Record highest hint level used
	if current, ok := s.HintsUsed[objectiveID]; !ok || level > current {
		s.HintsUsed[objectiveID] = level
	}

	return hint
}

// AllObjectivesComplete returns true if all objectives are done
func (s *State) AllObjectivesComplete() bool {
	for _, obj := range s.Scenario.Objectives {
		if _, done := s.CompletedObjectives[obj.ID]; !done {
			return false
		}
	}
	return true
}

// GetProgress returns the completion progress (0.0 to 1.0)
func (s *State) GetProgress() float64 {
	if len(s.Scenario.Objectives) == 0 {
		return 1.0
	}
	return float64(len(s.CompletedObjectives)) / float64(len(s.Scenario.Objectives))
}

// GetElapsedTime returns how long the player has been playing
func (s *State) GetElapsedTime() time.Duration {
	return time.Since(s.StartTime)
}

// GetPendingObjectives returns objectives not yet completed
func (s *State) GetPendingObjectives() []scenario.Objective {
	var pending []scenario.Objective
	for _, obj := range s.Scenario.Objectives {
		if _, done := s.CompletedObjectives[obj.ID]; !done && !obj.Hidden {
			pending = append(pending, obj)
		}
	}
	return pending
}

// GetCompletedObjectives returns completed objectives
func (s *State) GetCompletedObjectives() []scenario.Objective {
	var completed []scenario.Objective
	for _, obj := range s.Scenario.Objectives {
		if _, done := s.CompletedObjectives[obj.ID]; done {
			completed = append(completed, obj)
		}
	}
	return completed
}
