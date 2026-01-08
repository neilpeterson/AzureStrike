package game

import (
	"time"
)

// Score tracks the player's score and achievements
type Score struct {
	Points       int
	PointHistory []PointEntry
	Bonuses      map[string]Bonus
	Achievements []Achievement
}

// PointEntry records a score change
type PointEntry struct {
	Amount    int
	Reason    string
	Timestamp time.Time
}

// Bonus represents a one-time bonus
type Bonus struct {
	ID          string
	Name        string
	Points      int
	Description string
	Awarded     time.Time
}

// Achievement represents an unlockable achievement
type Achievement struct {
	ID          string
	Name        string
	Description string
	Icon        string
	Unlocked    bool
	UnlockedAt  time.Time
}

// NewScore creates a new scoring tracker
func NewScore() *Score {
	return &Score{
		Points:       0,
		PointHistory: []PointEntry{},
		Bonuses:      make(map[string]Bonus),
		Achievements: []Achievement{},
	}
}

// AddPoints adds points to the score
func (s *Score) AddPoints(amount int, reason string) {
	s.Points += amount
	s.PointHistory = append(s.PointHistory, PointEntry{
		Amount:    amount,
		Reason:    reason,
		Timestamp: time.Now(),
	})
}

// DeductPoints removes points from the score
func (s *Score) DeductPoints(amount int, reason string) {
	s.Points -= amount
	if s.Points < 0 {
		s.Points = 0
	}
	s.PointHistory = append(s.PointHistory, PointEntry{
		Amount:    -amount,
		Reason:    reason,
		Timestamp: time.Now(),
	})
}

// AddBonus awards a one-time bonus (only once per ID)
func (s *Score) AddBonus(id string, points int, description string) bool {
	if _, exists := s.Bonuses[id]; exists {
		return false // Already awarded
	}

	s.Bonuses[id] = Bonus{
		ID:          id,
		Points:      points,
		Description: description,
		Awarded:     time.Now(),
	}
	s.AddPoints(points, "Bonus: "+description)
	return true
}

// UnlockAchievement unlocks an achievement
func (s *Score) UnlockAchievement(id, name, description, icon string) bool {
	// Check if already unlocked
	for i := range s.Achievements {
		if s.Achievements[i].ID == id {
			if s.Achievements[i].Unlocked {
				return false
			}
			s.Achievements[i].Unlocked = true
			s.Achievements[i].UnlockedAt = time.Now()
			return true
		}
	}

	// Add new achievement
	s.Achievements = append(s.Achievements, Achievement{
		ID:          id,
		Name:        name,
		Description: description,
		Icon:        icon,
		Unlocked:    true,
		UnlockedAt:  time.Now(),
	})
	return true
}

// GetUnlockedAchievements returns all unlocked achievements
func (s *Score) GetUnlockedAchievements() []Achievement {
	var unlocked []Achievement
	for _, a := range s.Achievements {
		if a.Unlocked {
			unlocked = append(unlocked, a)
		}
	}
	return unlocked
}

// PredefinedAchievements that can be earned across scenarios
var PredefinedAchievements = map[string]Achievement{
	"first_blood": {
		ID:          "first_blood",
		Name:        "First Blood",
		Description: "Complete your first objective",
		Icon:        "[!]",
	},
	"speed_demon": {
		ID:          "speed_demon",
		Name:        "Speed Demon",
		Description: "Complete a scenario in under 5 minutes",
		Icon:        "[>]",
	},
	"no_hints": {
		ID:          "no_hints",
		Name:        "Solo Operator",
		Description: "Complete a scenario without using any hints",
		Icon:        "[*]",
	},
	"perfect_score": {
		ID:          "perfect_score",
		Name:        "Perfect Score",
		Description: "Achieve maximum points on a scenario",
		Icon:        "[+]",
	},
	"blob_hunter": {
		ID:          "blob_hunter",
		Name:        "Blob Hunter",
		Description: "Download sensitive data from a misconfigured storage account",
		Icon:        "[B]",
	},
	"token_thief": {
		ID:          "token_thief",
		Name:        "Token Thief",
		Description: "Extract credentials from IMDS",
		Icon:        "[T]",
	},
	"lateral_mover": {
		ID:          "lateral_mover",
		Name:        "Lateral Mover",
		Description: "Move between resources using discovered credentials",
		Icon:        "[L]",
	},
}
