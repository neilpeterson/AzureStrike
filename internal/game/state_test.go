package game

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/azurestrike/azurestrike/internal/scenario"
)

func createTestScenario() *scenario.Scenario {
	return &scenario.Scenario{
		ID:         "test-scenario",
		Name:       "Test Scenario",
		Difficulty: "beginner",
		Objectives: []scenario.Objective{
			{ID: "obj1", Description: "First objective", Trigger: "az storage", Points: 100},
			{ID: "obj2", Description: "Second objective", Trigger: "regex:curl.*metadata", Points: 150},
			{ID: "obj3", Description: "Hidden objective", Trigger: "secret command", Points: 50, Hidden: true},
		},
		Hints: []scenario.Hint{
			{ObjectiveID: "obj1", Level: 1, Text: "Try az storage", PointCost: 10},
			{ObjectiveID: "obj1", Level: 2, Text: "Use az storage account list", PointCost: 25},
			{ObjectiveID: "obj2", Level: 1, Text: "Try curl", PointCost: 15},
		},
	}
}

func TestNewState(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	assert.NotNil(t, state)
	assert.Equal(t, sc, state.Scenario)
	assert.NotNil(t, state.CompletedObjectives)
	assert.Empty(t, state.CompletedObjectives)
	assert.NotNil(t, state.CommandHistory)
	assert.Empty(t, state.CommandHistory)
	assert.NotNil(t, state.HintsUsed)
	assert.Empty(t, state.HintsUsed)
	assert.NotNil(t, state.StorageEnv)
	assert.NotNil(t, state.EntraEnv)
	assert.NotNil(t, state.ComputeEnv)
	assert.NotNil(t, state.Score)
	assert.Equal(t, StatusPlaying, state.Status)
}

func TestRecordCommand(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	completed := state.RecordCommand("ls", "output", true)

	assert.Len(t, state.CommandHistory, 1)
	assert.Equal(t, "ls", state.CommandHistory[0].Command)
	assert.Equal(t, "output", state.CommandHistory[0].Output)
	assert.True(t, state.CommandHistory[0].Success)
	assert.Empty(t, completed)
}

func TestRecordCommandTriggersObjective(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	completed := state.RecordCommand("az storage account list", "output", true)

	assert.Contains(t, completed, "obj1")
	assert.Contains(t, state.CompletedObjectives, "obj1")
	assert.Equal(t, 100, state.Score.Points)
}

func TestRecordCommandWithRegexTrigger(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	completed := state.RecordCommand("curl http://169.254.169.254/metadata/instance", "response", true)

	assert.Contains(t, completed, "obj2")
	assert.Contains(t, state.CompletedObjectives, "obj2")
}

func TestMatchesTriggerSubstring(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	t.Run("matches substring", func(t *testing.T) {
		assert.True(t, state.matchesTrigger("az storage account list", "az storage"))
	})

	t.Run("case insensitive", func(t *testing.T) {
		assert.True(t, state.matchesTrigger("AZ STORAGE account list", "az storage"))
	})

	t.Run("no match", func(t *testing.T) {
		assert.False(t, state.matchesTrigger("az vm list", "az storage"))
	})
}

func TestMatchesTriggerRegex(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	t.Run("matches regex", func(t *testing.T) {
		assert.True(t, state.matchesTrigger("curl http://169.254.169.254/metadata/instance", "regex:curl.*metadata"))
	})

	t.Run("regex no match", func(t *testing.T) {
		assert.False(t, state.matchesTrigger("wget http://example.com", "regex:curl.*metadata"))
	})

	t.Run("complex regex", func(t *testing.T) {
		trigger := "regex:az storage blob (list|download)"
		assert.True(t, state.matchesTrigger("az storage blob list", trigger))
		assert.True(t, state.matchesTrigger("az storage blob download", trigger))
		assert.False(t, state.matchesTrigger("az storage blob upload", trigger))
	})
}

func TestCompleteObjective(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	state.CompleteObjective("obj1")

	assert.Contains(t, state.CompletedObjectives, "obj1")
	assert.Equal(t, 100, state.Score.Points)
}

func TestCompleteObjectiveOnlyOnce(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	state.CompleteObjective("obj1")
	state.CompleteObjective("obj1")

	assert.Equal(t, 100, state.Score.Points)
}

func TestCompleteAllObjectivesAwardsBonus(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	state.CompleteObjective("obj1")
	state.CompleteObjective("obj2")
	state.CompleteObjective("obj3")

	assert.Equal(t, StatusCompleted, state.Status)
	assert.Contains(t, state.Score.Bonuses, "scenario_complete")
	// Points: 100 + 150 + 50 + 100 (bonus) = 400
	assert.Equal(t, 400, state.Score.Points)
}

func TestUseHint(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	t.Run("returns hint", func(t *testing.T) {
		hint := state.UseHint("obj1", 1)
		assert.NotNil(t, hint)
		assert.Equal(t, "Try az storage", hint.Text)
		assert.Equal(t, 10, hint.PointCost)
	})

	t.Run("records hint usage", func(t *testing.T) {
		assert.Equal(t, 1, state.HintsUsed["obj1"])
	})

	t.Run("returns higher level hint", func(t *testing.T) {
		hint := state.UseHint("obj1", 2)
		assert.NotNil(t, hint)
		assert.Equal(t, "Use az storage account list", hint.Text)
		assert.Equal(t, 2, state.HintsUsed["obj1"])
	})

	t.Run("returns nil for nonexistent hint", func(t *testing.T) {
		hint := state.UseHint("obj1", 5)
		assert.Nil(t, hint)
	})

	t.Run("returns nil for nonexistent objective", func(t *testing.T) {
		hint := state.UseHint("nonexistent", 1)
		assert.Nil(t, hint)
	})
}

func TestUseHintReducesPoints(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	state.UseHint("obj1", 1) // 10 points cost
	state.CompleteObjective("obj1")

	// 100 points for objective - 10 for hint = 90
	assert.Equal(t, 90, state.Score.Points)
}

func TestUseMultipleHintsReducesPoints(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	state.UseHint("obj1", 1) // 10 points
	state.UseHint("obj1", 2) // 25 points
	state.CompleteObjective("obj1")

	// 100 - 10 - 25 = 65
	assert.Equal(t, 65, state.Score.Points)
}

func TestAllObjectivesComplete(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	assert.False(t, state.AllObjectivesComplete())

	state.CompleteObjective("obj1")
	assert.False(t, state.AllObjectivesComplete())

	state.CompleteObjective("obj2")
	assert.False(t, state.AllObjectivesComplete())

	state.CompleteObjective("obj3")
	assert.True(t, state.AllObjectivesComplete())
}

func TestGetProgress(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	assert.Equal(t, 0.0, state.GetProgress())

	state.CompleteObjective("obj1")
	assert.InDelta(t, 0.333, state.GetProgress(), 0.01)

	state.CompleteObjective("obj2")
	assert.InDelta(t, 0.666, state.GetProgress(), 0.01)

	state.CompleteObjective("obj3")
	assert.Equal(t, 1.0, state.GetProgress())
}

func TestGetProgressEmptyObjectives(t *testing.T) {
	sc := &scenario.Scenario{
		ID:         "empty",
		Name:       "Empty",
		Objectives: []scenario.Objective{},
	}
	state := NewState(sc)

	assert.Equal(t, 1.0, state.GetProgress())
}

func TestGetPendingObjectives(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	pending := state.GetPendingObjectives()

	// Should not include hidden objectives
	assert.Len(t, pending, 2)

	state.CompleteObjective("obj1")
	pending = state.GetPendingObjectives()
	assert.Len(t, pending, 1)
}

func TestGetPendingObjectivesExcludesHidden(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	pending := state.GetPendingObjectives()

	for _, obj := range pending {
		assert.False(t, obj.Hidden, "Hidden objectives should not appear in pending list")
	}
}

func TestGetCompletedObjectives(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	completed := state.GetCompletedObjectives()
	assert.Empty(t, completed)

	state.CompleteObjective("obj1")
	completed = state.GetCompletedObjectives()
	assert.Len(t, completed, 1)
	assert.Equal(t, "obj1", completed[0].ID)

	state.CompleteObjective("obj2")
	completed = state.GetCompletedObjectives()
	assert.Len(t, completed, 2)
}

func TestGetElapsedTime(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	elapsed := state.GetElapsedTime()

	// Should be very small but positive
	assert.GreaterOrEqual(t, elapsed.Nanoseconds(), int64(0))
}

func TestGameStatus(t *testing.T) {
	assert.Equal(t, GameStatus(0), StatusPlaying)
	assert.Equal(t, GameStatus(1), StatusCompleted)
	assert.Equal(t, GameStatus(2), StatusFailed)
}
