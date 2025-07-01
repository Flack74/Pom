package cmd

import (
	"fmt"

	"github.com/Flack74/pom/config"
	"github.com/spf13/cobra"
)

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "🧩 Plugin and script system",
	Long: `🧩 Plugin/Script System

Run custom scripts at different points in your Pomodoro sessions:
  • Session start/end hooks
  • Break start/end hooks
  • Integration with external tools
  • Custom notifications

Examples:
  pom plugins list              List all plugins
  pom plugins enable notion-logger  Enable a plugin
  pom plugins disable slack-notify  Disable a plugin
  pom plugins add "my-script" "echo 'Session done!'" session_end`,
}

var listPluginsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available plugins",
	Run: func(cmd *cobra.Command, args []string) {
		plugins, err := config.LoadPlugins()
		if err != nil {
			fmt.Printf("Error loading plugins: %v\n", err)
			return
		}

		fmt.Println("🧩 Available Plugins:\n")
		for _, plugin := range plugins.Plugins {
			status := "❌ Disabled"
			if plugin.Enabled {
				status = "✅ Enabled"
			}
			
			fmt.Printf("  %s %s\n", status, plugin.Name)
			fmt.Printf("    %s\n", plugin.Description)
			fmt.Printf("    Triggers: %v\n", plugin.Triggers)
			fmt.Println()
		}
	},
}

var enablePluginCmd = &cobra.Command{
	Use:   "enable [plugin-name]",
	Short: "Enable a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]
		
		if err := config.EnablePlugin(pluginName, true); err != nil {
			fmt.Printf("Error enabling plugin: %v\n", err)
			return
		}

		fmt.Printf("✅ Plugin '%s' enabled\n", pluginName)
	},
}

var disablePluginCmd = &cobra.Command{
	Use:   "disable [plugin-name]",
	Short: "Disable a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]
		
		if err := config.EnablePlugin(pluginName, false); err != nil {
			fmt.Printf("Error disabling plugin: %v\n", err)
			return
		}

		fmt.Printf("❌ Plugin '%s' disabled\n", pluginName)
	},
}

var addPluginCmd = &cobra.Command{
	Use:   "add [name] [script] [trigger]",
	Short: "Add a custom plugin",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		script := args[1]
		trigger := args[2]
		
		description, _ := cmd.Flags().GetString("description")
		if description == "" {
			description = "Custom plugin"
		}

		plugin := config.Plugin{
			Name:        name,
			Description: description,
			Script:      script,
			Triggers:    []string{trigger},
			Enabled:     false,
		}

		if err := config.AddPlugin(plugin); err != nil {
			fmt.Printf("Error adding plugin: %v\n", err)
			return
		}

		fmt.Printf("✅ Plugin '%s' added (disabled by default)\n", name)
		fmt.Printf("Enable it with: pom plugins enable %s\n", name)
	},
}

func init() {
	addPluginCmd.Flags().String("description", "", "Plugin description")
	
	pluginsCmd.AddCommand(listPluginsCmd)
	pluginsCmd.AddCommand(enablePluginCmd)
	pluginsCmd.AddCommand(disablePluginCmd)
	pluginsCmd.AddCommand(addPluginCmd)
	rootCmd.AddCommand(pluginsCmd)
}