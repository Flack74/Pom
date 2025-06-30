package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	WorkMinutes  int `json:"work_minutes"`
	BreakMinutes int `json:"break_minutes"`
	NumSessions  int `json:"num_sessions"`
}

var DefaultConfig = Config{
	WorkMinutes:  25,
	BreakMinutes: 5,
	NumSessions:  4,
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".pomorc"), nil
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return DefaultConfig, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig, nil
		}
		return DefaultConfig, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultConfig, err
	}

	return config, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
