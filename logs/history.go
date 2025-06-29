package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type SessionLog struct {
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	WorkMinutes  int       `json:"work_minutes"`
	BreakMinutes int       `json:"break_minutes"`
	NumSessions  int       `json:"num_sessions"`
	Completed    bool      `json:"completed"`
}

// getLogPath returns the path to the log file
func getLogPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	logDir := filepath.Join(homeDir, ".pom")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(logDir, "pomodoro_history.log"), nil
}

// LogSession logs a completed Pomodoro session
func LogSession(workMin, breakMin, numSessions int, startTime, endTime time.Time, completed bool) error {
	logPath, err := getLogPath()
	if err != nil {
		return fmt.Errorf("failed to get log path: %v", err)
	}

	session := SessionLog{
		StartTime:    startTime,
		EndTime:      endTime,
		WorkMinutes:  workMin,
		BreakMinutes: breakMin,
		NumSessions:  numSessions,
		Completed:    completed,
	}

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %v", err)
	}

	// Append the log entry with a newline
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write to log file: %v", err)
	}

	return nil
}

// GetSessionStats returns statistics about completed Pomodoro sessions
func GetSessionStats() (totalSessions int, totalFocusMinutes float64, avgSessionsPerDay float64, err error) {
	logPath, err := getLogPath()
	if err != nil {
		return 0, 0, 0, err
	}

	// Check if log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return 0, 0, 0, nil
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		return 0, 0, 0, err
	}

	var sessions []SessionLog
	var firstDay, lastDay time.Time

	// Parse each line as a separate JSON object
	lines := string(data)
	for _, line := range filepath.SplitList(lines) {
		if line == "" {
			continue
		}

		var session SessionLog
		if err := json.Unmarshal([]byte(line), &session); err != nil {
			continue // Skip invalid entries
		}

		if session.Completed {
			sessions = append(sessions, session)
			totalFocusMinutes += float64(session.WorkMinutes * session.NumSessions)

			// Track first and last day for average calculation
			if firstDay.IsZero() || session.StartTime.Before(firstDay) {
				firstDay = session.StartTime
			}
			if session.StartTime.After(lastDay) {
				lastDay = session.StartTime
			}
		}
	}

	totalSessions = len(sessions)
	if totalSessions > 0 {
		daysDiff := lastDay.Sub(firstDay).Hours() / 24
		if daysDiff < 1 {
			daysDiff = 1
		}
		avgSessionsPerDay = float64(totalSessions) / daysDiff
	}

	return totalSessions, totalFocusMinutes, avgSessionsPerDay, nil
}
