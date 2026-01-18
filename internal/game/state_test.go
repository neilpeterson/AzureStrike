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

func TestMatchesCommandTriggerSubstring(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	t.Run("matches substring", func(t *testing.T) {
		assert.True(t, state.matchesCommandTrigger("az storage account list", "az storage"))
	})

	t.Run("case insensitive", func(t *testing.T) {
		assert.True(t, state.matchesCommandTrigger("AZ STORAGE account list", "az storage"))
	})

	t.Run("no match", func(t *testing.T) {
		assert.False(t, state.matchesCommandTrigger("az vm list", "az storage"))
	})
}

func TestMatchesCommandTriggerRegex(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	t.Run("matches regex", func(t *testing.T) {
		assert.True(t, state.matchesCommandTrigger("curl http://169.254.169.254/metadata/instance", "regex:curl.*metadata"))
	})

	t.Run("regex no match", func(t *testing.T) {
		assert.False(t, state.matchesCommandTrigger("wget http://example.com", "regex:curl.*metadata"))
	})

	t.Run("complex regex", func(t *testing.T) {
		trigger := "regex:az storage blob (list|download)"
		assert.True(t, state.matchesCommandTrigger("az storage blob list", trigger))
		assert.True(t, state.matchesCommandTrigger("az storage blob download", trigger))
		assert.False(t, state.matchesCommandTrigger("az storage blob upload", trigger))
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

func TestStoreToken(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	state.StoreToken("test-token-123", "https://storage.azure.com/", "principal-id-123")

	assert.Contains(t, state.ExtractedTokens, "test-token-123")
	assert.Equal(t, "https://storage.azure.com/", state.ExtractedTokens["test-token-123"].Resource)
	assert.Equal(t, "principal-id-123", state.ExtractedTokens["test-token-123"].PrincipalID)
}

func TestValidateToken(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	// Store a token scoped to storage
	state.StoreToken("storage-token", "https://storage.azure.com/", "principal-1")
	// Store a token scoped to management
	state.StoreToken("management-token", "https://management.azure.com/", "principal-2")

	t.Run("valid storage token for storage resource", func(t *testing.T) {
		assert.True(t, state.ValidateToken("storage-token", "storage"))
	})

	t.Run("management token works for storage", func(t *testing.T) {
		// ARM tokens can access storage
		assert.True(t, state.ValidateToken("management-token", "storage"))
	})

	t.Run("management token for management resource", func(t *testing.T) {
		assert.True(t, state.ValidateToken("management-token", "management"))
	})

	t.Run("storage token invalid for management", func(t *testing.T) {
		assert.False(t, state.ValidateToken("storage-token", "management"))
	})

	t.Run("invalid token", func(t *testing.T) {
		assert.False(t, state.ValidateToken("nonexistent-token", "storage"))
	})
}

func TestGetTokenPrincipal(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	state.StoreToken("test-token", "https://storage.azure.com/", "principal-abc")

	t.Run("returns principal for valid token", func(t *testing.T) {
		principal := state.GetTokenPrincipal("test-token")
		assert.Equal(t, "principal-abc", principal)
	})

	t.Run("returns empty for invalid token", func(t *testing.T) {
		principal := state.GetTokenPrincipal("invalid-token")
		assert.Empty(t, principal)
	})
}

func TestExtractedTokensInitialized(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	assert.NotNil(t, state.ExtractedTokens)
	assert.Empty(t, state.ExtractedTokens)
}

// Tests for event-based objective system

func createStateEventScenario() *scenario.Scenario {
	return &scenario.Scenario{
		ID:         "event-test-scenario",
		Name:       "Event Test Scenario",
		Difficulty: "beginner",
		Objectives: []scenario.Objective{
			{ID: "cmd_obj", Description: "Command-based objective", Trigger: "az storage", Points: 100},
			{ID: "event_obj", Description: "Event-based objective", Trigger: "state:token_extracted", Points: 150},
			{ID: "event_target", Description: "Event with target", Trigger: "state:blob_read:secrets.txt", Points: 200},
			{ID: "event_wildcard", Description: "Event with wildcard", Trigger: "state:blob_listed:*", Points: 75},
		},
	}
}

func TestEventsInitialized(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	assert.NotNil(t, state.Events)
	assert.Empty(t, state.Events)
}

func TestEmitEvent(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	completed := state.EmitEvent(EventTokenExtracted, "https://storage.azure.com/", map[string]string{
		"principal_id": "test-principal",
	})

	assert.Len(t, state.Events, 1)
	assert.Equal(t, EventTokenExtracted, state.Events[0].Type)
	assert.Equal(t, "https://storage.azure.com/", state.Events[0].Target)
	assert.Equal(t, "test-principal", state.Events[0].Data["principal_id"])
	assert.Empty(t, completed) // No state-based objectives in test scenario
}

func TestEmitEventTriggersStateObjective(t *testing.T) {
	sc := createStateEventScenario()
	state := NewState(sc)

	completed := state.EmitEvent(EventTokenExtracted, "https://storage.azure.com/", nil)

	assert.Contains(t, completed, "event_obj")
	assert.Contains(t, state.CompletedObjectives, "event_obj")
	assert.Equal(t, 150, state.Score.Points)
}

func TestEmitEventWithTargetMatch(t *testing.T) {
	sc := createStateEventScenario()
	state := NewState(sc)

	// Should not match - wrong target
	completed := state.EmitEvent(EventBlobRead, "other-file.txt", nil)
	assert.Empty(t, completed)

	// Should match - correct target
	completed = state.EmitEvent(EventBlobRead, "secrets.txt", nil)
	assert.Contains(t, completed, "event_target")
	assert.Equal(t, 200, state.Score.Points)
}

func TestEmitEventWithWildcardMatch(t *testing.T) {
	sc := createStateEventScenario()
	state := NewState(sc)

	// Any blob_listed event should match the wildcard trigger
	completed := state.EmitEvent(EventBlobListed, "any-container", nil)

	assert.Contains(t, completed, "event_wildcard")
	assert.Equal(t, 75, state.Score.Points)
}

func TestMatchesStateTrigger(t *testing.T) {
	sc := createTestScenario()
	state := NewState(sc)

	t.Run("simple event type match", func(t *testing.T) {
		event := GameEvent{Type: EventTokenExtracted, Target: "https://storage.azure.com/"}
		assert.True(t, state.matchesStateTrigger(event, "state:token_extracted"))
	})

	t.Run("event type with target match", func(t *testing.T) {
		event := GameEvent{Type: EventBlobRead, Target: "secrets.txt"}
		assert.True(t, state.matchesStateTrigger(event, "state:blob_read:secrets.txt"))
	})

	t.Run("event type with wrong target", func(t *testing.T) {
		event := GameEvent{Type: EventBlobRead, Target: "other.txt"}
		assert.False(t, state.matchesStateTrigger(event, "state:blob_read:secrets.txt"))
	})

	t.Run("event type with wildcard target", func(t *testing.T) {
		event := GameEvent{Type: EventBlobRead, Target: "any-file.txt"}
		assert.True(t, state.matchesStateTrigger(event, "state:blob_read:*"))
	})

	t.Run("wrong event type", func(t *testing.T) {
		event := GameEvent{Type: EventBlobListed, Target: "container"}
		assert.False(t, state.matchesStateTrigger(event, "state:token_extracted"))
	})

	t.Run("case insensitive target match", func(t *testing.T) {
		event := GameEvent{Type: EventBlobRead, Target: "SECRETS.TXT"}
		assert.True(t, state.matchesStateTrigger(event, "state:blob_read:secrets.txt"))
	})
}

func TestCommandObjectivesSkipStateTriggersAndViceVersa(t *testing.T) {
	sc := createStateEventScenario()
	state := NewState(sc)

	// Command should only trigger command-based objectives
	completed := state.RecordCommand("az storage account list", "output", true)
	assert.Contains(t, completed, "cmd_obj")
	assert.NotContains(t, completed, "event_obj")

	// Event should only trigger state-based objectives
	completed = state.EmitEvent(EventTokenExtracted, "https://storage.azure.com/", nil)
	assert.Contains(t, completed, "event_obj")
	assert.NotContains(t, completed, "cmd_obj") // Already completed anyway
}

func TestMultipleEventsAndObjectives(t *testing.T) {
	sc := createStateEventScenario()
	state := NewState(sc)

	// First event
	state.EmitEvent(EventTokenExtracted, "https://storage.azure.com/", nil)
	assert.Len(t, state.CompletedObjectives, 1)

	// Second event
	state.EmitEvent(EventBlobListed, "container", nil)
	assert.Len(t, state.CompletedObjectives, 2)

	// Third event
	state.EmitEvent(EventBlobRead, "secrets.txt", nil)
	assert.Len(t, state.CompletedObjectives, 3)

	// Command
	state.RecordCommand("az storage account list", "output", true)
	assert.Len(t, state.CompletedObjectives, 4)

	// All objectives complete
	assert.True(t, state.AllObjectivesComplete())
}
