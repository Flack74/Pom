package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"pom/config"
)

var (
	workMinutes  int
	breakMinutes int
	numSessions  int
	saveConfig   bool
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Pomodoro timer",
	Long: `Start the Pomodoro timer with custom durations.
If no flags are provided, values from config file or defaults will be used.
Default values: 25 minutes work, 5 minutes break, 4 sessions.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Use flags if provided, otherwise use config values
		if !cmd.Flags().Changed("work") {
			workMinutes = cfg.WorkMinutes
		}
		if !cmd.Flags().Changed("break") {
			breakMinutes = cfg.BreakMinutes
		}
		if !cmd.Flags().Changed("sessions") {
			numSessions = cfg.NumSessions
		}

		// Save config if requested
		if saveConfig {
			newConfig := config.Config{
				WorkMinutes:  workMinutes,
				BreakMinutes: breakMinutes,
				NumSessions:  numSessions,
			}
			if err := config.SaveConfig(newConfig); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Configuration saved successfully!")
		}

		StartPomodoro(workMinutes, breakMinutes, numSessions)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Add flags with default values from config.DefaultConfig
	startCmd.Flags().IntVarP(&workMinutes, "work", "w", config.DefaultConfig.WorkMinutes, "work time in minutes")
	startCmd.Flags().IntVarP(&breakMinutes, "break", "b", config.DefaultConfig.BreakMinutes, "break time in minutes")
	startCmd.Flags().IntVarP(&numSessions, "sessions", "s", config.DefaultConfig.NumSessions, "number of sessions")
	startCmd.Flags().BoolVarP(&saveConfig, "save-config", "c", false, "save current settings as default configuration")
}
