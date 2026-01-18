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
	Events              []GameEvent                // Event log for state-based triggers
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

// GameEvent represents a game state change that can trigger objectives
type GameEvent struct {
	Type      string            // Event type (e.g., "token_extracted", "blob_read")
	Target    string            // Target of the event (e.g., blob name, resource URL)
	Timestamp time.Time         // When the event occurred
	Data      map[string]string // Additional event metadata
}

// Event type constants
const (
	EventTokenExtracted       = "token_extracted"
	EventMetadataRetrieved    = "metadata_retrieved"
	EventBlobRead             = "blob_read"
	EventBlobListed           = "blob_listed"
	EventContainerListed      = "container_listed"
	EventStorageAuthenticated = "storage_authenticated"
	EventSecretDiscovered     = "secret_discovered"
	EventIMDSProbed           = "imds_probed"
)

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
		Events:              []GameEvent{},
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

	// Check for command-based objective completion
	return s.checkCommandObjectives(cmd)
}

// EmitEvent records a game event and checks for state-based objective completion
func (s *State) EmitEvent(eventType, target string, data map[string]string) []string {
	event := GameEvent{
		Type:      eventType,
		Target:    target,
		Timestamp: time.Now(),
		Data:      data,
	}
	s.Events = append(s.Events, event)

	// Check for state-based objective completion
	return s.checkStateObjectives(event)
}

// checkCommandObjectives tests if a command triggers any command-based objectives
func (s *State) checkCommandObjectives(cmd string) []string {
	var completed []string

	for _, obj := range s.Scenario.Objectives {
		// Skip already completed objectives
		if _, done := s.CompletedObjectives[obj.ID]; done {
			continue
		}

		// Skip state-based triggers (they're checked by EmitEvent)
		if strings.HasPrefix(obj.Trigger, "state:") {
			continue
		}

		// Skip if prerequisites not met
		if !s.prerequisitesMet(obj) {
			continue
		}

		if s.matchesCommandTrigger(cmd, obj.Trigger) {
			s.CompleteObjective(obj.ID)
			completed = append(completed, obj.ID)
		}
	}

	return completed
}

// checkStateObjectives tests if an event triggers any state-based objectives
func (s *State) checkStateObjectives(event GameEvent) []string {
	var completed []string

	for _, obj := range s.Scenario.Objectives {
		// Skip already completed objectives
		if _, done := s.CompletedObjectives[obj.ID]; done {
			continue
		}

		// Only check state-based triggers
		if !strings.HasPrefix(obj.Trigger, "state:") {
			continue
		}

		// Skip if prerequisites not met
		if !s.prerequisitesMet(obj) {
			continue
		}

		if s.matchesStateTrigger(event, obj.Trigger) {
			s.CompleteObjective(obj.ID)
			completed = append(completed, obj.ID)
		}
	}

	return completed
}

// prerequisitesMet checks if all required objectives are completed
func (s *State) prerequisitesMet(obj scenario.Objective) bool {
	for _, reqID := range obj.Requires {
		if _, done := s.CompletedObjectives[reqID]; !done {
			return false
		}
	}
	return true
}

// matchesCommandTrigger checks if a command matches a command-based trigger pattern
func (s *State) matchesCommandTrigger(cmd, trigger string) bool {
	// Triggers can be:
	// - Simple substring match: "blob list"
	// - Regex pattern: "regex:az storage blob (list|download)"

	if strings.HasPrefix(trigger, "regex:") {
		pattern := strings.TrimPrefix(trigger, "regex:")
		matched, _ := regexp.MatchString(pattern, cmd)
		return matched
	}

	// Simple substring match (case-insensitive)
	return strings.Contains(strings.ToLower(cmd), strings.ToLower(trigger))
}

// matchesStateTrigger checks if an event matches a state-based trigger
func (s *State) matchesStateTrigger(event GameEvent, trigger string) bool {
	// State triggers have format: "state:event_type" or "state:event_type:target"
	// Examples:
	//   "state:token_extracted" - matches any token extraction
	//   "state:blob_read:secrets.txt" - matches reading specific blob
	//   "state:blob_read:*" - matches reading any blob

	parts := strings.SplitN(strings.TrimPrefix(trigger, "state:"), ":", 2)
	if len(parts) == 0 {
		return false
	}

	eventType := parts[0]

	// Check event type matches
	if event.Type != eventType {
		return false
	}

	// If no target specified, any event of this type matches
	if len(parts) == 1 {
		return true
	}

	targetPattern := parts[1]

	// Wildcard matches any target
	if targetPattern == "*" {
		return true
	}

	// Exact target match (case-insensitive)
	return strings.EqualFold(event.Target, targetPattern)
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
