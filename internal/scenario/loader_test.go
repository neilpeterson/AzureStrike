package scenario

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Create a temporary scenario file
	dir := t.TempDir()
	scenarioPath := filepath.Join(dir, "test-scenario.yaml")

	scenarioContent := `
id: test-scenario
name: Test Scenario
difficulty: beginner
briefing: |
  This is a test briefing.
objectives:
  - id: obj1
    description: First objective
    trigger: "test command"
    points: 100
debrief: |
  Test debrief.
`
	err := os.WriteFile(scenarioPath, []byte(scenarioContent), 0644)
	assert.NoError(t, err)

	t.Run("loads valid scenario", func(t *testing.T) {
		sc, err := Load(scenarioPath)
		assert.NoError(t, err)
		assert.NotNil(t, sc)
		assert.Equal(t, "test-scenario", sc.ID)
		assert.Equal(t, "Test Scenario", sc.Name)
		assert.Equal(t, "beginner", sc.Difficulty)
		assert.Len(t, sc.Objectives, 1)
	})

	t.Run("returns error for nonexistent file", func(t *testing.T) {
		sc, err := Load("/nonexistent/path.yaml")
		assert.Error(t, err)
		assert.Nil(t, sc)
		assert.Contains(t, err.Error(), "failed to read scenario file")
	})

	t.Run("returns error for invalid YAML", func(t *testing.T) {
		invalidPath := filepath.Join(dir, "invalid.yaml")
		err := os.WriteFile(invalidPath, []byte("invalid: yaml: content: ["), 0644)
		assert.NoError(t, err)

		sc, err := Load(invalidPath)
		assert.Error(t, err)
		assert.Nil(t, sc)
		assert.Contains(t, err.Error(), "failed to parse scenario YAML")
	})
}

func TestValidate(t *testing.T) {
	t.Run("valid scenario", func(t *testing.T) {
		sc := &Scenario{
			ID:   "valid-scenario",
			Name: "Valid Scenario",
			Objectives: []Objective{
				{ID: "obj1", Description: "Test", Points: 100},
			},
		}
		err := sc.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing ID", func(t *testing.T) {
		sc := &Scenario{
			Name: "No ID",
			Objectives: []Objective{
				{ID: "obj1"},
			},
		}
		err := sc.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have an ID")
	})

	t.Run("missing name", func(t *testing.T) {
		sc := &Scenario{
			ID: "no-name",
			Objectives: []Objective{
				{ID: "obj1"},
			},
		}
		err := sc.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have a name")
	})

	t.Run("no objectives", func(t *testing.T) {
		sc := &Scenario{
			ID:         "no-objectives",
			Name:       "No Objectives",
			Objectives: []Objective{},
		}
		err := sc.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have at least one objective")
	})

	t.Run("objective without ID", func(t *testing.T) {
		sc := &Scenario{
			ID:   "bad-objective",
			Name: "Bad Objective",
			Objectives: []Objective{
				{Description: "No ID"},
			},
		}
		err := sc.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have an ID")
	})

	t.Run("duplicate objective IDs", func(t *testing.T) {
		sc := &Scenario{
			ID:   "duplicate-objectives",
			Name: "Duplicate Objectives",
			Objectives: []Objective{
				{ID: "same-id", Description: "First"},
				{ID: "same-id", Description: "Second"},
			},
		}
		err := sc.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate objective ID")
	})
}

func TestGetObjective(t *testing.T) {
	sc := &Scenario{
		Objectives: []Objective{
			{ID: "obj1", Description: "First", Points: 100},
			{ID: "obj2", Description: "Second", Points: 200},
			{ID: "obj3", Description: "Third", Points: 150},
		},
	}

	t.Run("existing objective", func(t *testing.T) {
		obj := sc.GetObjective("obj2")
		assert.NotNil(t, obj)
		assert.Equal(t, "Second", obj.Description)
		assert.Equal(t, 200, obj.Points)
	})

	t.Run("first objective", func(t *testing.T) {
		obj := sc.GetObjective("obj1")
		assert.NotNil(t, obj)
		assert.Equal(t, "First", obj.Description)
	})

	t.Run("last objective", func(t *testing.T) {
		obj := sc.GetObjective("obj3")
		assert.NotNil(t, obj)
		assert.Equal(t, "Third", obj.Description)
	})

	t.Run("nonexistent objective", func(t *testing.T) {
		obj := sc.GetObjective("nonexistent")
		assert.Nil(t, obj)
	})
}

func TestTotalPoints(t *testing.T) {
	t.Run("single objective", func(t *testing.T) {
		sc := &Scenario{
			Objectives: []Objective{
				{ID: "obj1", Points: 100},
			},
		}
		assert.Equal(t, 100, sc.TotalPoints())
	})

	t.Run("multiple objectives", func(t *testing.T) {
		sc := &Scenario{
			Objectives: []Objective{
				{ID: "obj1", Points: 100},
				{ID: "obj2", Points: 200},
				{ID: "obj3", Points: 50},
			},
		}
		assert.Equal(t, 350, sc.TotalPoints())
	})

	t.Run("no objectives", func(t *testing.T) {
		sc := &Scenario{
			Objectives: []Objective{},
		}
		assert.Equal(t, 0, sc.TotalPoints())
	})

	t.Run("zero point objectives", func(t *testing.T) {
		sc := &Scenario{
			Objectives: []Objective{
				{ID: "obj1", Points: 0},
				{ID: "obj2", Points: 100},
			},
		}
		assert.Equal(t, 100, sc.TotalPoints())
	})
}

func TestLoadByID(t *testing.T) {
	// Create temp scenarios directory structure
	dir := t.TempDir()

	// Create scenario in directory format
	scenarioDir := filepath.Join(dir, "test-scenario")
	err := os.MkdirAll(scenarioDir, 0755)
	assert.NoError(t, err)

	scenarioContent := `
id: test-scenario
name: Test Scenario
difficulty: beginner
objectives:
  - id: obj1
    description: Test
    points: 100
`
	err = os.WriteFile(filepath.Join(scenarioDir, "scenario.yaml"), []byte(scenarioContent), 0644)
	assert.NoError(t, err)

	t.Run("finds scenario by ID in directory", func(t *testing.T) {
		sc, err := LoadByID(dir, "test-scenario")
		assert.NoError(t, err)
		assert.NotNil(t, sc)
		assert.Equal(t, "test-scenario", sc.ID)
	})

	t.Run("returns error for nonexistent scenario", func(t *testing.T) {
		sc, err := LoadByID(dir, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, sc)
		assert.Contains(t, err.Error(), "scenario not found")
	})
}

func TestListScenarios(t *testing.T) {
	dir := t.TempDir()

	// Create multiple scenarios
	for i, id := range []string{"scenario-a", "scenario-b"} {
		scenarioDir := filepath.Join(dir, id)
		err := os.MkdirAll(scenarioDir, 0755)
		assert.NoError(t, err)

		content := `
id: ` + id + `
name: Scenario ` + string(rune('A'+i)) + `
difficulty: beginner
objectives:
  - id: obj1
    description: Test
    points: 100
`
		err = os.WriteFile(filepath.Join(scenarioDir, "scenario.yaml"), []byte(content), 0644)
		assert.NoError(t, err)
	}

	t.Run("lists all scenarios", func(t *testing.T) {
		scenarios, err := ListScenarios(dir)
		assert.NoError(t, err)
		assert.Len(t, scenarios, 2)
	})

	t.Run("returns error for nonexistent directory", func(t *testing.T) {
		scenarios, err := ListScenarios("/nonexistent/path")
		assert.Error(t, err)
		assert.Nil(t, scenarios)
	})
}

func TestObjectiveFields(t *testing.T) {
	obj := Objective{
		ID:          "test-obj",
		Description: "Complete the test",
		Trigger:     "test command",
		Points:      150,
		Hidden:      true,
		Order:       2,
	}

	assert.Equal(t, "test-obj", obj.ID)
	assert.Equal(t, "Complete the test", obj.Description)
	assert.Equal(t, "test command", obj.Trigger)
	assert.Equal(t, 150, obj.Points)
	assert.True(t, obj.Hidden)
	assert.Equal(t, 2, obj.Order)
}

func TestHintFields(t *testing.T) {
	hint := Hint{
		ObjectiveID: "obj1",
		Level:       2,
		Text:        "Try using the az command",
		PointCost:   25,
	}

	assert.Equal(t, "obj1", hint.ObjectiveID)
	assert.Equal(t, 2, hint.Level)
	assert.Equal(t, "Try using the az command", hint.Text)
	assert.Equal(t, 25, hint.PointCost)
}

func TestScenarioFields(t *testing.T) {
	sc := Scenario{
		ID:         "full-scenario",
		Name:       "Full Scenario Test",
		Difficulty: "intermediate",
		Briefing:   "Welcome to the scenario",
		Debrief:    "Congratulations!",
	}

	assert.Equal(t, "full-scenario", sc.ID)
	assert.Equal(t, "Full Scenario Test", sc.Name)
	assert.Equal(t, "intermediate", sc.Difficulty)
	assert.Equal(t, "Welcome to the scenario", sc.Briefing)
	assert.Equal(t, "Congratulations!", sc.Debrief)
}
