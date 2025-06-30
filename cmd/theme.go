package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Flack74/pom/config"
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "🎨 Manage themes for pom",
	Long: `🎨 Customize Your Pomodoro Experience

Change the look and feel of your Pomodoro timer with beautiful themes:

Available Themes:
  • default  - Professional and clean
  • minimal  - Distraction-free experience
  • vibrant  - Colorful and energetic

Examples:
  pom theme list         List all available themes
  pom theme set vibrant  Switch to the vibrant theme
  pom theme show        Show current theme settings`,
	Run: func(cmd *cobra.Command, args []string) {
		config.ListThemes()
	},
}

var listThemesCmd = &cobra.Command{
	Use:   "list",
	Short: "List available themes",
	Run: func(cmd *cobra.Command, args []string) {
		config.ListThemes()
	},
}

var setThemeCmd = &cobra.Command{
	Use:   "set [theme-name]",
	Short: "Set the active theme",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		themeName := args[0]
		theme, ok := config.AvailableThemes[themeName]
		if !ok {
			fmt.Printf("Theme '%s' not found. Use 'pom theme list' to see available themes.\n", themeName)
			return
		}

		if err := config.SaveTheme(theme); err != nil {
			fmt.Printf("Error saving theme: %v\n", err)
			return
		}

		fmt.Printf("Theme set to '%s'\n", themeName)
	},
}

func init() {
	themeCmd.AddCommand(listThemesCmd)
	themeCmd.AddCommand(setThemeCmd)
	rootCmd.AddCommand(themeCmd)
}
