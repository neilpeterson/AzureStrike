package scenario

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/azurestrike/azurestrike/internal/azure/compute"
	"github.com/azurestrike/azurestrike/internal/azure/entra"
	"github.com/azurestrike/azurestrike/internal/azure/storage"
)

// Scenario represents a complete game scenario
type Scenario struct {
	ID         string      `yaml:"id"`
	Name       string      `yaml:"name"`
	Difficulty string      `yaml:"difficulty"` // beginner, intermediate, advanced
	Briefing   string      `yaml:"briefing"`
	Resources  Resources   `yaml:"resources"`
	Objectives []Objective `yaml:"objectives"`
	Debrief    string      `yaml:"debrief"`
	Hints      []Hint      `yaml:"hints"`
}

// Resources contains all mocked Azure resources for the scenario
type Resources struct {
	// Storage
	StorageAccounts []storage.Account `yaml:"storage_accounts"`

	// Entra ID (Azure AD)
	Users             []entra.User             `yaml:"users"`
	ServicePrincipals []entra.ServicePrincipal `yaml:"service_principals"`
	Groups            []entra.Group            `yaml:"groups"`

	// Compute
	VirtualMachines       []compute.VirtualMachine       `yaml:"virtual_machines"`
	NetworkSecurityGroups []compute.NetworkSecurityGroup `yaml:"network_security_groups"`
}

// Objective represents a goal the player must achieve
type Objective struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
	Trigger     string `yaml:"trigger"` // Action pattern that completes this objective
	Points      int    `yaml:"points"`
	Hidden      bool   `yaml:"hidden"` // Hidden objectives are revealed when completed
	Order       int    `yaml:"order"`  // Suggested order (0 = any order)
}

// Hint provides progressive assistance to the player
type Hint struct {
	ObjectiveID string `yaml:"objective_id"`
	Level       int    `yaml:"level"` // 1 = vague, 2 = moderate, 3 = specific
	Text        string `yaml:"text"`
	PointCost   int    `yaml:"point_cost"` // Points deducted for using hint
}

// Load reads and parses a scenario from a YAML file
func Load(path string) (*Scenario, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read scenario file: %w", err)
	}

	var sc Scenario
	if err := yaml.Unmarshal(data, &sc); err != nil {
		return nil, fmt.Errorf("failed to parse scenario YAML: %w", err)
	}

	if err := sc.Validate(); err != nil {
		return nil, fmt.Errorf("invalid scenario: %w", err)
	}

	return &sc, nil
}

// LoadByID finds and loads a scenario by its ID from a scenarios directory
func LoadByID(scenariosDir, id string) (*Scenario, error) {
	// Try common patterns
	patterns := []string{
		filepath.Join(scenariosDir, id, "scenario.yaml"),
		filepath.Join(scenariosDir, "*-"+id, "scenario.yaml"),
		filepath.Join(scenariosDir, id+".yaml"),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		if len(matches) > 0 {
			return Load(matches[0])
		}
	}

	// Try to find by scanning directories
	entries, err := os.ReadDir(scenariosDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read scenarios directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		scenarioPath := filepath.Join(scenariosDir, entry.Name(), "scenario.yaml")
		if _, err := os.Stat(scenarioPath); err != nil {
			continue
		}
		sc, err := Load(scenarioPath)
		if err != nil {
			continue
		}
		if sc.ID == id {
			return sc, nil
		}
	}

	return nil, fmt.Errorf("scenario not found: %s", id)
}

// ListScenarios returns metadata for all available scenarios
func ListScenarios(scenariosDir string) ([]Scenario, error) {
	var scenarios []Scenario

	entries, err := os.ReadDir(scenariosDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read scenarios directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		scenarioPath := filepath.Join(scenariosDir, entry.Name(), "scenario.yaml")
		sc, err := Load(scenarioPath)
		if err != nil {
			continue
		}
		scenarios = append(scenarios, *sc)
	}

	return scenarios, nil
}

// Validate checks that the scenario is properly configured
func (s *Scenario) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("scenario must have an ID")
	}
	if s.Name == "" {
		return fmt.Errorf("scenario must have a name")
	}
	if len(s.Objectives) == 0 {
		return fmt.Errorf("scenario must have at least one objective")
	}

	// Validate objectives have unique IDs
	seen := make(map[string]bool)
	for _, obj := range s.Objectives {
		if obj.ID == "" {
			return fmt.Errorf("all objectives must have an ID")
		}
		if seen[obj.ID] {
			return fmt.Errorf("duplicate objective ID: %s", obj.ID)
		}
		seen[obj.ID] = true
	}

	return nil
}

// GetObjective returns an objective by ID
func (s *Scenario) GetObjective(id string) *Objective {
	for i := range s.Objectives {
		if s.Objectives[i].ID == id {
			return &s.Objectives[i]
		}
	}
	return nil
}

// TotalPoints returns the maximum possible score for the scenario
func (s *Scenario) TotalPoints() int {
	total := 0
	for _, obj := range s.Objectives {
		total += obj.Points
	}
	return total
}

// GetNextScenario returns the next scenario after the current one, or nil if none
func GetNextScenario(scenariosDir, currentID string) (*Scenario, error) {
	scenarios, err := ListScenarios(scenariosDir)
	if err != nil {
		return nil, err
	}

	// Scenarios are returned in directory order (alphabetical by folder name)
	// Folder names like "01-storage-breach", "02-imds-token-theft" sort correctly
	foundCurrent := false
	for _, sc := range scenarios {
		if foundCurrent {
			// Return the next scenario after the current one
			return LoadByID(scenariosDir, sc.ID)
		}
		if sc.ID == currentID {
			foundCurrent = true
		}
	}

	return nil, nil // No next scenario (completed all scenarios)
}
