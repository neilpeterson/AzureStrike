package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScore(t *testing.T) {
	score := NewScore()

	assert.NotNil(t, score)
	assert.Equal(t, 0, score.Points)
	assert.NotNil(t, score.PointHistory)
	assert.Empty(t, score.PointHistory)
	assert.NotNil(t, score.Bonuses)
	assert.Empty(t, score.Bonuses)
	assert.NotNil(t, score.Achievements)
	assert.Empty(t, score.Achievements)
}

func TestAddPoints(t *testing.T) {
	score := NewScore()

	score.AddPoints(100, "Test objective")
	assert.Equal(t, 100, score.Points)
	assert.Len(t, score.PointHistory, 1)
	assert.Equal(t, 100, score.PointHistory[0].Amount)
	assert.Equal(t, "Test objective", score.PointHistory[0].Reason)

	score.AddPoints(50, "Another objective")
	assert.Equal(t, 150, score.Points)
	assert.Len(t, score.PointHistory, 2)
}

func TestAddPointsMultiple(t *testing.T) {
	score := NewScore()

	score.AddPoints(10, "First")
	score.AddPoints(20, "Second")
	score.AddPoints(30, "Third")

	assert.Equal(t, 60, score.Points)
	assert.Len(t, score.PointHistory, 3)
}

func TestDeductPoints(t *testing.T) {
	score := NewScore()
	score.Points = 100

	score.DeductPoints(30, "Hint used")
	assert.Equal(t, 70, score.Points)
	assert.Len(t, score.PointHistory, 1)
	assert.Equal(t, -30, score.PointHistory[0].Amount)
	assert.Equal(t, "Hint used", score.PointHistory[0].Reason)
}

func TestDeductPointsDoesNotGoNegative(t *testing.T) {
	score := NewScore()
	score.Points = 20

	score.DeductPoints(50, "Large deduction")
	assert.Equal(t, 0, score.Points)
}

func TestDeductPointsFromZero(t *testing.T) {
	score := NewScore()

	score.DeductPoints(10, "Deduction from zero")
	assert.Equal(t, 0, score.Points)
}

func TestAddBonus(t *testing.T) {
	score := NewScore()

	awarded := score.AddBonus("first_find", 50, "First discovery")
	assert.True(t, awarded)
	assert.Equal(t, 50, score.Points)
	assert.Contains(t, score.Bonuses, "first_find")
	assert.Equal(t, 50, score.Bonuses["first_find"].Points)
	assert.Equal(t, "First discovery", score.Bonuses["first_find"].Description)
}

func TestAddBonusOnlyOnce(t *testing.T) {
	score := NewScore()

	first := score.AddBonus("unique_bonus", 100, "Unique bonus")
	assert.True(t, first)
	assert.Equal(t, 100, score.Points)

	second := score.AddBonus("unique_bonus", 100, "Unique bonus")
	assert.False(t, second)
	assert.Equal(t, 100, score.Points)
}

func TestAddMultipleBonuses(t *testing.T) {
	score := NewScore()

	score.AddBonus("bonus1", 10, "First bonus")
	score.AddBonus("bonus2", 20, "Second bonus")
	score.AddBonus("bonus3", 30, "Third bonus")

	assert.Equal(t, 60, score.Points)
	assert.Len(t, score.Bonuses, 3)
}

func TestUnlockAchievement(t *testing.T) {
	score := NewScore()

	unlocked := score.UnlockAchievement("first_blood", "First Blood", "Complete first objective", "[!]")
	assert.True(t, unlocked)
	assert.Len(t, score.Achievements, 1)
	assert.Equal(t, "first_blood", score.Achievements[0].ID)
	assert.Equal(t, "First Blood", score.Achievements[0].Name)
	assert.Equal(t, "Complete first objective", score.Achievements[0].Description)
	assert.Equal(t, "[!]", score.Achievements[0].Icon)
	assert.True(t, score.Achievements[0].Unlocked)
}

func TestUnlockAchievementOnlyOnce(t *testing.T) {
	score := NewScore()

	first := score.UnlockAchievement("test_achievement", "Test", "Test achievement", "[T]")
	assert.True(t, first)

	second := score.UnlockAchievement("test_achievement", "Test", "Test achievement", "[T]")
	assert.False(t, second)
	assert.Len(t, score.Achievements, 1)
}

func TestUnlockAchievementExistingNotUnlocked(t *testing.T) {
	score := NewScore()

	// Add an achievement that's not unlocked
	score.Achievements = append(score.Achievements, Achievement{
		ID:       "existing",
		Name:     "Existing",
		Unlocked: false,
	})

	// Unlock it
	unlocked := score.UnlockAchievement("existing", "Existing", "Desc", "[E]")
	assert.True(t, unlocked)
	assert.True(t, score.Achievements[0].Unlocked)
}

func TestGetUnlockedAchievements(t *testing.T) {
	score := NewScore()

	score.UnlockAchievement("a1", "Achievement 1", "Desc 1", "[1]")
	score.UnlockAchievement("a2", "Achievement 2", "Desc 2", "[2]")

	unlocked := score.GetUnlockedAchievements()
	assert.Len(t, unlocked, 2)
}

func TestGetUnlockedAchievementsFiltersLocked(t *testing.T) {
	score := NewScore()

	score.Achievements = append(score.Achievements, Achievement{
		ID:       "locked",
		Name:     "Locked Achievement",
		Unlocked: false,
	})
	score.UnlockAchievement("unlocked", "Unlocked Achievement", "Desc", "[U]")

	unlocked := score.GetUnlockedAchievements()
	assert.Len(t, unlocked, 1)
	assert.Equal(t, "unlocked", unlocked[0].ID)
}

func TestGetUnlockedAchievementsEmpty(t *testing.T) {
	score := NewScore()

	unlocked := score.GetUnlockedAchievements()
	assert.Nil(t, unlocked)
}

func TestPredefinedAchievementsExist(t *testing.T) {
	expectedAchievements := []string{
		"first_blood",
		"speed_demon",
		"no_hints",
		"perfect_score",
		"blob_hunter",
		"token_thief",
		"lateral_mover",
	}

	for _, id := range expectedAchievements {
		_, exists := PredefinedAchievements[id]
		assert.True(t, exists, "Expected achievement %s to exist", id)
	}
}
