package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Session represents a completed Pomodoro session
type Session struct {
	WorkMinutes  int       `json:"work_minutes"`
	BreakMinutes int       `json:"break_minutes"`
	NumSessions  int       `json:"num_sessions"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	IsCompleted  bool      `json:"is_completed"`
}

// getLogFilePath returns the path to the session log file
func getLogFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	logDir := filepath.Join(homeDir, ".config", "pom", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(logDir, "sessions.json"), nil
}

// LogSession logs a completed Pomodoro session
func LogSession(workMin, breakMin, numSessions int, startTime, endTime time.Time, isCompleted bool) error {
	logPath, err := getLogFilePath()
	if err != nil {
		return fmt.Errorf("failed to get log path: %v", err)
	}

	// Create new session entry
	session := Session{
		WorkMinutes:  workMin,
		BreakMinutes: breakMin,
		NumSessions:  numSessions,
		StartTime:    startTime,
		EndTime:      endTime,
		IsCompleted:  isCompleted,
	}

	// Read existing sessions
	var sessions []Session
	data, err := os.ReadFile(logPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read log file: %v", err)
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &sessions); err != nil {
			return fmt.Errorf("failed to parse log file: %v", err)
		}
	}

	// Add new session
	sessions = append(sessions, session)

	// Save updated sessions
	data, err = json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal log data: %v", err)
	}

	if err := os.WriteFile(logPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write log file: %v", err)
	}

	return nil
}

// GetSessionStats returns statistics about completed Pomodoro sessions
func GetSessionStats() (totalSessions int, totalFocusMinutes float64, avgSessionsPerDay float64, err error) {
	logPath, err := getLogFilePath()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get log path: %v", err)
	}

	// Load sessions
	var sessions []Session
	data, err := os.ReadFile(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, 0, 0, nil
		}
		return 0, 0, 0, fmt.Errorf("failed to read log file: %v", err)
	}

	if err := json.Unmarshal(data, &sessions); err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse log file: %v", err)
	}

	// Calculate statistics
	if len(sessions) == 0 {
		return 0, 0, 0, nil
	}

	var completedSessions []Session
	for _, session := range sessions {
		if session.IsCompleted {
			completedSessions = append(completedSessions, session)
			totalSessions += session.NumSessions
			totalFocusMinutes += float64(session.WorkMinutes * session.NumSessions)
		}
	}

	if len(completedSessions) == 0 {
		return 0, 0, 0, nil
	}

	// Calculate average sessions per day
	firstSession := completedSessions[0].StartTime
	lastSession := completedSessions[len(completedSessions)-1].EndTime
	daysDiff := lastSession.Sub(firstSession).Hours() / 24
	if daysDiff < 1 {
		daysDiff = 1
	}
	avgSessionsPerDay = float64(totalSessions) / daysDiff

	return totalSessions, totalFocusMinutes, avgSessionsPerDay, nil
}

// GetDailyStats returns statistics for the current day
func GetDailyStats() (sessions int, minutes int, err error) {
	logPath, err := getLogFilePath()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get log path: %v", err)
	}

	// Load sessions
	var allSessions []Session
	data, err := os.ReadFile(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, 0, nil
		}
		return 0, 0, fmt.Errorf("failed to read log file: %v", err)
	}

	if err := json.Unmarshal(data, &allSessions); err != nil {
		return 0, 0, fmt.Errorf("failed to parse log file: %v", err)
	}

	// Get today's date
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	// Calculate today's statistics
	for _, session := range allSessions {
		if session.IsCompleted && session.EndTime.After(startOfDay) {
			sessions += session.NumSessions
			minutes += session.WorkMinutes * session.NumSessions
		}
	}

	return sessions, minutes, nil
}
