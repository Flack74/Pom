package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/Flack74/pom/config"
	"github.com/Flack74/pom/logs"
)

var (
	workMin      int
	breakMin     int
	numberOfSess int
	saveConfig   bool
	taskID       string
	profileName  string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "🎯 Start a pomodoro session",
	Long: `🎯 Start a Pomodoro Timer Session

Start a focused work session using the Pomodoro Technique. You can customize:
  • Work duration (default: 25 minutes)
  • Break duration (default: 5 minutes)
  • Number of sessions
  • Link to a planned task

During the session:
  • Press 'p' to pause
  • Press 'r' to resume
  • Press 'q' to quit (progress is saved)

Examples:
  pom start                     Start with default settings
  pom start -w 30 -b 10        30min work + 10min break
  pom start -s 4               Do 4 sessions
  pom start -t task-id         Link to a planned task
  pom start -c                 Save settings as default`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load profile settings if specified
		if profileName != "" {
			profile, err := config.GetProfile(profileName)
			if err != nil {
				fmt.Printf("Profile '%s' not found, using default settings\n", profileName)
			} else {
				if workMin == 25 { workMin = profile.WorkMinutes }
				if breakMin == 5 { breakMin = profile.BreakMinutes }
				if numberOfSess == 1 { numberOfSess = profile.NumSessions }
				fmt.Printf("Using profile: %s\n", profile.Name)
			}
		} else {
			// Load current profile from config
			cfg, _ := config.LoadConfig()
			if cfg.CurrentProfile != "" {
				profile, err := config.GetProfile(cfg.CurrentProfile)
				if err == nil {
					if workMin == 25 { workMin = profile.WorkMinutes }
					if breakMin == 5 { breakMin = profile.BreakMinutes }
					if numberOfSess == 1 { numberOfSess = profile.NumSessions }
				}
			}
		}

		if saveConfig {
			if err := SaveConfig(workMin, breakMin, numberOfSess); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			}
		}

		// Create a channel to handle interrupts
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Execute session start plugins
		sessionData := map[string]string{
			"DURATION": fmt.Sprintf("%d", workMin),
			"BREAK_DURATION": fmt.Sprintf("%d", breakMin),
			"SESSIONS": fmt.Sprintf("%d", numberOfSess),
			"DATE": time.Now().Format("2006-01-02T15:04:05Z"),
			"TASK_ID": taskID,
		}
		config.ExecutePlugins("session_start", sessionData)

		// Start the timer in a goroutine
		doneChan := make(chan bool)
		go func() {
			startTime := time.Now()
			isCompleted := StartPomodoro(workMin, breakMin, numberOfSess, taskID)
			doneChan <- isCompleted

			// Log the session
			endTime := time.Now()
			if err := logs.LogSession(workMin, breakMin, numberOfSess, startTime, endTime, isCompleted); err != nil {
				fmt.Fprintf(os.Stderr, "⚠️  Failed to log session: %v\n", err)
			}

			// Update task progress if a task is linked
			if taskID != "" && isCompleted {
				if err := config.UpdateTaskProgress(taskID, 1, workMin); err != nil {
					fmt.Fprintf(os.Stderr, "⚠️  Failed to update task progress: %v\n", err)
				}
			}

			// Update goals progress
			if isCompleted {
				if err := config.UpdateProgress(1, workMin); err != nil {
					fmt.Fprintf(os.Stderr, "⚠️  Failed to update goals progress: %v\n", err)
				}
			}

			// Execute session end plugins
			sessionData["COMPLETED"] = fmt.Sprintf("%t", isCompleted)
			sessionData["TOTAL_MINUTES"] = fmt.Sprintf("%d", workMin*numberOfSess)
			config.ExecutePlugins("session_end", sessionData)
		}()

		// Wait for either completion or interrupt
		select {
		case <-sigChan:
			fmt.Println("\n⚠️  Pomodoro session interrupted")
			os.Exit(1)
		case isCompleted := <-doneChan:
			if isCompleted {
				fmt.Println("🎉 Pomodoro session completed!")
				if taskID != "" {
					fmt.Println("📝 Task progress updated")
				}
			}
		}
	},
}

func init() {
	startCmd.Flags().IntVarP(&workMin, "work", "w", 25, "work minutes")
	startCmd.Flags().IntVarP(&breakMin, "break", "b", 5, "break minutes")
	startCmd.Flags().IntVarP(&numberOfSess, "sessions", "s", 1, "number of sessions")
	startCmd.Flags().BoolVarP(&saveConfig, "save-config", "c", false, "save as default configuration")
	startCmd.Flags().StringVarP(&taskID, "task", "t", "", "link session to a task ID")
	startCmd.Flags().StringVarP(&profileName, "profile", "p", "", "use specific profile")

	rootCmd.AddCommand(startCmd)
}
