package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	WorkMinutes  int    `json:"work_minutes"`
	BreakMinutes int    `json:"break_minutes"`
	NumSessions  int    `json:"num_sessions"`
	CurrentProfile string `json:"current_profile"`
	PrivacyMode  bool   `json:"privacy_mode"`
	CloudSync    bool   `json:"cloud_sync"`
	CloudProvider string `json:"cloud_provider"`
}

type Profile struct {
	Name         string `json:"name"`
	WorkMinutes  int    `json:"work_minutes"`
	BreakMinutes int    `json:"break_minutes"`
	NumSessions  int    `json:"num_sessions"`
	Description  string `json:"description"`
}

type ProfileConfig struct {
	Profiles []Profile `json:"profiles"`
}

var DefaultConfig = Config{
	WorkMinutes:    25,
	BreakMinutes:   5,
	NumSessions:    4,
	CurrentProfile: "default",
	PrivacyMode:    false,
	CloudSync:      false,
	CloudProvider:  "",
}

var DefaultProfiles = []Profile{
	{Name: "default", WorkMinutes: 25, BreakMinutes: 5, NumSessions: 4, Description: "Standard Pomodoro"},
	{Name: "work", WorkMinutes: 45, BreakMinutes: 10, NumSessions: 3, Description: "Deep work sessions"},
	{Name: "study", WorkMinutes: 30, BreakMinutes: 5, NumSessions: 4, Description: "Study sessions"},
	{Name: "quick", WorkMinutes: 15, BreakMinutes: 3, NumSessions: 6, Description: "Quick tasks"},
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

// GetConfigDir returns the path to the configuration directory
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "pom")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// Profile management functions
func GetProfilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "profiles.json"), nil
}

func LoadProfiles() (ProfileConfig, error) {
	profilePath, err := GetProfilePath()
	if err != nil {
		return ProfileConfig{}, err
	}

	data, err := os.ReadFile(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ProfileConfig{Profiles: DefaultProfiles}, nil
		}
		return ProfileConfig{}, err
	}

	var profiles ProfileConfig
	if err := json.Unmarshal(data, &profiles); err != nil {
		return ProfileConfig{}, err
	}

	return profiles, nil
}

func SaveProfiles(profiles ProfileConfig) error {
	profilePath, err := GetProfilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(profilePath, data, 0644)
}

func GetProfile(name string) (Profile, error) {
	profiles, err := LoadProfiles()
	if err != nil {
		return Profile{}, err
	}

	for _, profile := range profiles.Profiles {
		if profile.Name == name {
			return profile, nil
		}
	}

	return Profile{}, os.ErrNotExist
}

func AddProfile(profile Profile) error {
	profiles, err := LoadProfiles()
	if err != nil {
		return err
	}

	profiles.Profiles = append(profiles.Profiles, profile)
	return SaveProfiles(profiles)
}