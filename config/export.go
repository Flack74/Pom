package config

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type ExportData struct {
	Sessions []SessionData `json:"sessions"`
	Tasks    []Task        `json:"tasks"`
	Goal     Goal          `json:"goal"`
	Progress GoalProgress  `json:"progress"`
	Config   Config        `json:"config"`
	Profiles []Profile     `json:"profiles"`
}

type SessionData struct {
	Date         time.Time `json:"date"`
	WorkMinutes  int       `json:"work_minutes"`
	BreakMinutes int       `json:"break_minutes"`
	Sessions     int       `json:"sessions"`
	Completed    bool      `json:"completed"`
	Profile      string    `json:"profile"`
}

func ExportToJSON(filepath string) error {
	// Load all data
	tasks, _ := LoadTasks()
	goal, _ := LoadGoal()
	progress, _ := LoadProgress()
	config, _ := LoadConfig()
	profiles, _ := LoadProfiles()
	sessions, _ := loadSessionHistory()

	exportData := ExportData{
		Sessions: sessions,
		Tasks:    tasks.Tasks,
		Goal:     goal,
		Progress: progress,
		Config:   config,
		Profiles: profiles.Profiles,
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

func ExportToCSV(filepath string) error {
	sessions, err := loadSessionHistory()
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Date", "Profile", "Work Minutes", "Break Minutes", "Sessions", "Completed"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, session := range sessions {
		record := []string{
			session.Date.Format("2006-01-02 15:04:05"),
			session.Profile,
			strconv.Itoa(session.WorkMinutes),
			strconv.Itoa(session.BreakMinutes),
			strconv.Itoa(session.Sessions),
			strconv.FormatBool(session.Completed),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func ImportFromJSON(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var exportData ExportData
	if err := json.Unmarshal(data, &exportData); err != nil {
		return err
	}

	// Import tasks
	if len(exportData.Tasks) > 0 {
		taskList := TaskList{Tasks: exportData.Tasks}
		if err := SaveTasks(taskList); err != nil {
			return fmt.Errorf("failed to import tasks: %v", err)
		}
	}

	// Import profiles
	if len(exportData.Profiles) > 0 {
		profiles := ProfileConfig{Profiles: exportData.Profiles}
		if err := SaveProfiles(profiles); err != nil {
			return fmt.Errorf("failed to import profiles: %v", err)
		}
	}

	// Import config
	if err := SaveConfig(exportData.Config); err != nil {
		return fmt.Errorf("failed to import config: %v", err)
	}

	return nil
}

func loadSessionHistory() ([]SessionData, error) {
	// This would load from logs/session.go data
	// For now, return empty slice
	return []SessionData{}, nil
}