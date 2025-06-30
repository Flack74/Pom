package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Theme represents a color scheme for the application
type Theme struct {
	Name            string `json:"name"`
	TimerColor      string `json:"timer_color"`      // Color for timer display
	ProgressColor   string `json:"progress_color"`   // Color for progress bar
	SuccessColor    string `json:"success_color"`    // Color for success messages
	WarningColor    string `json:"warning_color"`    // Color for warning messages
	TextColor       string `json:"text_color"`       // Color for regular text
	HighlightColor  string `json:"highlight_color"`  // Color for highlighted text
	BackgroundColor string `json:"background_color"` // Color for background (if supported)
}

var (
	// DefaultTheme is the default color scheme
	DefaultTheme = Theme{
		Name:           "default",
		TimerColor:     "\033[1;32m", // Bold Green
		ProgressColor:  "\033[1;34m", // Bold Blue
		SuccessColor:   "\033[1;32m", // Bold Green
		WarningColor:   "\033[1;33m", // Bold Yellow
		TextColor:      "\033[0m",    // Reset
		HighlightColor: "\033[1;36m", // Bold Cyan
	}

	// MinimalTheme is a minimal color scheme
	MinimalTheme = Theme{
		Name:           "minimal",
		TimerColor:     "\033[0m", // Reset
		ProgressColor:  "\033[0m", // Reset
		SuccessColor:   "\033[0m", // Reset
		WarningColor:   "\033[0m", // Reset
		TextColor:      "\033[0m", // Reset
		HighlightColor: "\033[1m", // Bold
	}

	// VibrantTheme is a colorful theme
	VibrantTheme = Theme{
		Name:           "vibrant",
		TimerColor:     "\033[1;35m", // Bold Magenta
		ProgressColor:  "\033[1;36m", // Bold Cyan
		SuccessColor:   "\033[1;32m", // Bold Green
		WarningColor:   "\033[1;31m", // Bold Red
		TextColor:      "\033[1;37m", // Bold White
		HighlightColor: "\033[1;33m", // Bold Yellow
	}

	// Available themes
	AvailableThemes = map[string]Theme{
		"default": DefaultTheme,
		"minimal": MinimalTheme,
		"vibrant": VibrantTheme,
	}
)

// GetThemeFilePath returns the path to the theme configuration file
func GetThemeFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "theme.json"), nil
}

// SaveTheme saves the current theme to the configuration file
func SaveTheme(theme Theme) error {
	themePath, err := GetThemeFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(theme, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(themePath, data, 0644)
}

// LoadTheme loads the theme from the configuration file
func LoadTheme() (Theme, error) {
	themePath, err := GetThemeFilePath()
	if err != nil {
		return DefaultTheme, err
	}

	data, err := os.ReadFile(themePath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultTheme, nil
		}
		return DefaultTheme, err
	}

	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return DefaultTheme, err
	}

	return theme, nil
}

// ListThemes prints all available themes
func ListThemes() {
	fmt.Println("Available themes:")
	for name, theme := range AvailableThemes {
		fmt.Printf("  %s%s%s\n", theme.HighlightColor, name, "\033[0m")
	}
}
