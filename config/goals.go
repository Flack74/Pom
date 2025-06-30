package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Goal represents a daily Pomodoro goal
type Goal struct {
	DailySessionTarget int       `json:"daily_session_target"` // Target number of sessions per day
	DailyMinutes       int       `json:"daily_minutes"`        // Target minutes of focus time per day
	StartDate          time.Time `json:"start_date"`           // When this goal was set
	EndDate            time.Time `json:"end_date,omitempty"`   // Optional end date for the goal
}

// GoalProgress tracks progress towards goals
type GoalProgress struct {
	CurrentDate    time.Time `json:"current_date"`
	SessionsToday  int       `json:"sessions_today"`
	MinutesToday   int       `json:"minutes_today"`
	CurrentStreak  int       `json:"current_streak"`   // Days in a row meeting goals
	LongestStreak  int       `json:"longest_streak"`   // Longest streak ever
	LastUpdateDate time.Time `json:"last_update_date"` // Last time progress was updated
}

// GetGoalFilePath returns the path to the goals configuration file
func GetGoalFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "goals.json"), nil
}

// SaveGoal saves the current goal to the configuration file
func SaveGoal(goal Goal) error {
	goalPath, err := GetGoalFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(goal, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(goalPath, data, 0644)
}

// LoadGoal loads the goal from the configuration file
func LoadGoal() (Goal, error) {
	goalPath, err := GetGoalFilePath()
	if err != nil {
		return Goal{}, err
	}

	data, err := os.ReadFile(goalPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Goal{}, nil
		}
		return Goal{}, err
	}

	var goal Goal
	if err := json.Unmarshal(data, &goal); err != nil {
		return Goal{}, err
	}

	return goal, nil
}

// UpdateProgress updates the goal progress for today
func UpdateProgress(sessions, minutes int) error {
	progress, err := LoadProgress()
	if err != nil {
		progress = GoalProgress{
			CurrentDate:   time.Now(),
			SessionsToday: 0,
			MinutesToday:  0,
			CurrentStreak: 0,
			LongestStreak: 0,
		}
	}

	// Check if we need to reset daily progress
	today := time.Now()
	if !isSameDay(progress.LastUpdateDate, today) {
		// It's a new day, check if yesterday's goals were met before resetting
		goal, err := LoadGoal()
		if err == nil && progress.SessionsToday >= goal.DailySessionTarget &&
			progress.MinutesToday >= goal.DailyMinutes {
			progress.CurrentStreak++
			if progress.CurrentStreak > progress.LongestStreak {
				progress.LongestStreak = progress.CurrentStreak
			}
		} else {
			progress.CurrentStreak = 0
		}

		progress.SessionsToday = 0
		progress.MinutesToday = 0
		progress.CurrentDate = today
	}

	// Update today's progress
	progress.SessionsToday += sessions
	progress.MinutesToday += minutes
	progress.LastUpdateDate = today

	return SaveProgress(progress)
}

// SaveProgress saves the current progress to the configuration file
func SaveProgress(progress GoalProgress) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	progressPath := filepath.Join(configDir, "progress.json")
	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(progressPath, data, 0644)
}

// LoadProgress loads the progress from the configuration file
func LoadProgress() (GoalProgress, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return GoalProgress{}, err
	}

	progressPath := filepath.Join(configDir, "progress.json")
	data, err := os.ReadFile(progressPath)
	if err != nil {
		if os.IsNotExist(err) {
			return GoalProgress{
				CurrentDate:    time.Now(),
				LastUpdateDate: time.Now(),
			}, nil
		}
		return GoalProgress{}, err
	}

	var progress GoalProgress
	if err := json.Unmarshal(data, &progress); err != nil {
		return GoalProgress{}, err
	}

	return progress, nil
}

// ShowProgress displays the current progress towards goals
func ShowProgress() error {
	goal, err := LoadGoal()
	if err != nil {
		return fmt.Errorf("failed to load goal: %v", err)
	}

	progress, err := LoadProgress()
	if err != nil {
		return fmt.Errorf("failed to load progress: %v", err)
	}

	fmt.Printf("\nDaily Goals Progress:\n")
	fmt.Printf("Sessions: %d/%d\n", progress.SessionsToday, goal.DailySessionTarget)
	fmt.Printf("Minutes:  %d/%d\n", progress.MinutesToday, goal.DailyMinutes)
	fmt.Printf("Current Streak: %d days\n", progress.CurrentStreak)
	fmt.Printf("Longest Streak: %d days\n", progress.LongestStreak)

	return nil
}

// Helper function to check if two times are on the same day
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
