package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Plugin struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Script      string   `json:"script"`
	Triggers    []string `json:"triggers"` // "session_start", "session_end", "break_start", "break_end"
	Enabled     bool     `json:"enabled"`
	Args        []string `json:"args"`
}

type PluginConfig struct {
	Plugins []Plugin `json:"plugins"`
}

func GetPluginPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "plugins.json"), nil
}

func LoadPlugins() (PluginConfig, error) {
	pluginPath, err := GetPluginPath()
	if err != nil {
		return PluginConfig{}, err
	}

	data, err := os.ReadFile(pluginPath)
	if err != nil {
		if os.IsNotExist(err) {
			return PluginConfig{Plugins: getDefaultPlugins()}, nil
		}
		return PluginConfig{}, err
	}

	var plugins PluginConfig
	if err := json.Unmarshal(data, &plugins); err != nil {
		return PluginConfig{}, err
	}

	return plugins, nil
}

func SavePlugins(plugins PluginConfig) error {
	pluginPath, err := GetPluginPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(plugins, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pluginPath, data, 0644)
}

func getDefaultPlugins() []Plugin {
	return []Plugin{
		{
			Name:        "notion-logger",
			Description: "Log sessions to Notion database",
			Script:      "curl -X POST https://api.notion.com/v1/pages -H 'Authorization: Bearer $NOTION_TOKEN' -H 'Content-Type: application/json' -d '{\"parent\":{\"database_id\":\"$NOTION_DB_ID\"},\"properties\":{\"Name\":{\"title\":[{\"text\":{\"content\":\"Pomodoro Session\"}}]},\"Duration\":{\"number\":$DURATION},\"Date\":{\"date\":{\"start\":\"$DATE\"}}}}'",
			Triggers:    []string{"session_end"},
			Enabled:     false,
			Args:        []string{},
		},
		{
			Name:        "slack-notify",
			Description: "Send Slack notification on session completion",
			Script:      "curl -X POST -H 'Content-type: application/json' --data '{\"text\":\"🍅 Completed a $DURATION minute focus session!\"}' $SLACK_WEBHOOK_URL",
			Triggers:    []string{"session_end"},
			Enabled:     false,
			Args:        []string{},
		},
		{
			Name:        "break-reminder",
			Description: "Play sound and show desktop notification",
			Script:      "notify-send 'Pomodoro' 'Time for a break!' && paplay /usr/share/sounds/alsa/Front_Left.wav",
			Triggers:    []string{"break_start"},
			Enabled:     false,
			Args:        []string{},
		},
		{
			Name:        "focus-mode",
			Description: "Block distracting websites during focus sessions",
			Script:      "echo '127.0.0.1 facebook.com twitter.com reddit.com' | sudo tee -a /etc/hosts",
			Triggers:    []string{"session_start"},
			Enabled:     false,
			Args:        []string{},
		},
	}
}

func ExecutePlugins(trigger string, sessionData map[string]string) error {
	plugins, err := LoadPlugins()
	if err != nil {
		return err
	}

	for _, plugin := range plugins.Plugins {
		if !plugin.Enabled {
			continue
		}

		// Check if plugin should run for this trigger
		shouldRun := false
		for _, t := range plugin.Triggers {
			if t == trigger {
				shouldRun = true
				break
			}
		}

		if !shouldRun {
			continue
		}

		if err := executePlugin(plugin, sessionData); err != nil {
			fmt.Printf("Plugin '%s' failed: %v\n", plugin.Name, err)
		}
	}

	return nil
}

func executePlugin(plugin Plugin, sessionData map[string]string) error {
	script := plugin.Script

	// Replace variables in script
	for key, value := range sessionData {
		script = strings.ReplaceAll(script, "$"+key, value)
	}

	// Replace environment variables
	script = os.ExpandEnv(script)

	// Execute script
	cmd := exec.Command("sh", "-c", script)
	cmd.Env = os.Environ()
	
	// Add session data as environment variables
	for key, value := range sessionData {
		cmd.Env = append(cmd.Env, fmt.Sprintf("POM_%s=%s", strings.ToUpper(key), value))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("script failed: %v, output: %s", err, string(output))
	}

	return nil
}

func AddPlugin(plugin Plugin) error {
	plugins, err := LoadPlugins()
	if err != nil {
		return err
	}

	plugins.Plugins = append(plugins.Plugins, plugin)
	return SavePlugins(plugins)
}

func EnablePlugin(name string, enabled bool) error {
	plugins, err := LoadPlugins()
	if err != nil {
		return err
	}

	for i, plugin := range plugins.Plugins {
		if plugin.Name == name {
			plugins.Plugins[i].Enabled = enabled
			return SavePlugins(plugins)
		}
	}

	return fmt.Errorf("plugin '%s' not found", name)
}