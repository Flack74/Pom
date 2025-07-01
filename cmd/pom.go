package cmd

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/Flack74/pom/config"
	"github.com/Flack74/pom/logs"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
)

// Timer states
const (
	stateRunning = iota
	statePaused
	stateQuitting
)

// handleUserInput handles keyboard input for pause/resume/quit
func handleUserInput(timerState *int, pauseChan chan struct{}, resumeChan chan struct{}) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n%s⌨️  Controls: [p]ause | [r]esume | [q]uit%s\n", colorBlue, colorReset)

	for {
		char, err := reader.ReadByte()
		if err != nil {
			continue
		}

		switch char {
		case 'p', 'P':
			if *timerState == stateRunning {
				*timerState = statePaused
				pauseChan <- struct{}{}
				fmt.Printf("\n%s⏸️  Timer paused. Press 'r' to resume.%s\n", colorYellow, colorReset)
			}
		case 'r', 'R':
			if *timerState == statePaused {
				*timerState = stateRunning
				resumeChan <- struct{}{}
				fmt.Printf("\n%s▶️  Timer resumed.%s\n", colorGreen, colorReset)
			}
		case 'q', 'Q':
			if *timerState != stateQuitting {
				*timerState = stateQuitting
				fmt.Printf("\n%s⏹️  Quitting...%s\n", colorRed, colorReset)
				return
			}
		}
	}
}

// countdown displays a live countdown timer with progress bar
func countdown(duration time.Duration, label string, color string, timerState *int, pauseChan, resumeChan chan struct{}) bool {
	startTime := time.Now()
	endTime := startTime.Add(duration)
	var pausedDuration time.Duration

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	// Create a ticker for the countdown
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Get terminal width for progress bar
	width := 40 // default width
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = w - 20 // leave room for timer and label
	}

	for {
		select {
		case <-sigChan:
			fmt.Printf("\n%s⚠️  Session interrupted!%s\n", colorRed, colorReset)
			return false
		case <-pauseChan:
			pauseStart := time.Now()
			<-resumeChan
			pausedDuration += time.Since(pauseStart)
			endTime = endTime.Add(time.Since(pauseStart))
		case <-ticker.C:
			if *timerState == stateQuitting {
				return false
			}
			if *timerState == statePaused {
				continue
			}
			if time.Now().After(endTime) {
				return true
			}

			remaining := time.Until(endTime).Round(time.Second)
			elapsed := duration - remaining
			progress := float64(elapsed) / float64(duration)
			
			// Ensure progress doesn't exceed 100%
			if progress > 1.0 {
				progress = 1.0
			}

			// Calculate progress bar with proper rounding
			barWidth := int(float64(width)*progress + 0.5)
			bar := strings.Repeat("█", barWidth) + strings.Repeat("░", width-barWidth)

			// Format time remaining
			minutes := int(remaining.Minutes())
			seconds := int(remaining.Seconds()) % 60

			// Clear line and print progress
			fmt.Printf("\r%s%s %02d:%02d [%s] %.0f%%%s",
				color, label, minutes, seconds, bar, progress*100, colorReset)
		}
	}
}

// StartPomodoro starts a pomodoro session with the given parameters
func StartPomodoro(workMin, breakMin, numberOfSess int, taskID string) bool {
	// Load theme
	theme, err := config.LoadTheme()
	if err != nil {
		theme = config.DefaultTheme
	}

	// Print session info
	fmt.Printf("\n%s🎯 Starting Pomodoro Timer%s\n", theme.HighlightColor, theme.TextColor)
	fmt.Printf("%s📚 Work: %d min | Break: %d min | Sessions: %d%s\n\n",
		theme.TextColor, workMin, breakMin, numberOfSess, theme.TextColor)

	// If a task ID is provided, verify it exists
	if taskID != "" {
		tasks, err := config.LoadTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s⚠️  Error loading tasks: %v%s\n", theme.WarningColor, err, theme.TextColor)
			return false
		}

		taskFound := false
		for _, task := range tasks.Tasks {
			if task.ID == taskID {
				taskFound = true
				fmt.Printf("%s📎 Linked to task: %s%s\n\n", theme.HighlightColor, task.Title, theme.TextColor)
				break
			}
		}

		if !taskFound {
			fmt.Fprintf(os.Stderr, "%s⚠️  Task with ID %s not found%s\n", theme.WarningColor, taskID, theme.TextColor)
			return false
		}
	}

	// Set up channels for pause/resume functionality
	pauseChan := make(chan struct{})
	resumeChan := make(chan struct{})
	timerState := stateRunning

	// Start user input handler
	go handleUserInput(&timerState, pauseChan, resumeChan)

	// Track total work time
	totalWorkTime := time.Duration(0)
	startTime := time.Now()

	// Run sessions
	for sess := 1; sess <= numberOfSess; sess++ {
		work := time.Duration(workMin) * time.Minute
		breakTime := time.Duration(breakMin) * time.Minute

		// Work period
		fmt.Printf("%s📚 Session %d/%d - Focus Time%s\n", theme.HighlightColor, sess, numberOfSess, theme.TextColor)
		if !countdown(work, "Focus", theme.TimerColor, &timerState, pauseChan, resumeChan) {
			return false
		}
		totalWorkTime += work

		// Show motivational message
		message := getRandomMotivationalMessage()
		fmt.Printf("\n%s%s%s\n", theme.SuccessColor, message, theme.TextColor)

		// Play sound and show notification
		if err := logs.PlaySound("work_end"); err != nil {
			fmt.Fprintf(os.Stderr, "%s⚠️  Error playing sound: %v%s\n", theme.WarningColor, err, theme.TextColor)
		}
		if err := logs.ShowNotification("Work session complete!", "Time for a break!"); err != nil {
			fmt.Fprintf(os.Stderr, "%s⚠️  Error showing notification: %v%s\n", theme.WarningColor, err, theme.TextColor)
		}

		// Break period (skip after last session)
		if sess < numberOfSess {
			// Execute break start plugins
			breakData := map[string]string{
				"DURATION": fmt.Sprintf("%d", breakMin),
				"SESSION": fmt.Sprintf("%d", sess),
				"DATE": time.Now().Format("2006-01-02T15:04:05Z"),
			}
			config.ExecutePlugins("break_start", breakData)

			fmt.Printf("\n%s☕ Break Time%s\n", theme.HighlightColor, theme.TextColor)
			if !countdown(breakTime, "Break", theme.ProgressColor, &timerState, pauseChan, resumeChan) {
				return false
			}

			// Execute break end plugins
			config.ExecutePlugins("break_end", breakData)

			// Play sound and show notification
			if err := logs.PlaySound("break_end"); err != nil {
				fmt.Fprintf(os.Stderr, "%s⚠️  Error playing sound: %v%s\n", theme.WarningColor, err, theme.TextColor)
			}
			if err := logs.ShowNotification("Break complete!", "Time to focus!"); err != nil {
				fmt.Fprintf(os.Stderr, "%s⚠️  Error showing notification: %v%s\n", theme.WarningColor, err, theme.TextColor)
			}
			fmt.Println()
		}
	}

	// Log the session
	endTime := time.Now()
	if err := logs.LogSession(workMin, breakMin, numberOfSess, startTime, endTime, true); err != nil {
		fmt.Fprintf(os.Stderr, "%s⚠️  Failed to log session: %v%s\n", theme.WarningColor, err, theme.TextColor)
	}

	// Update task progress if a task is linked
	if taskID != "" {
		if err := config.UpdateTaskProgress(taskID, numberOfSess, workMin*numberOfSess); err != nil {
			fmt.Fprintf(os.Stderr, "%s⚠️  Failed to update task progress: %v%s\n", theme.WarningColor, err, theme.TextColor)
		}
	}

	// Update goals progress
	if err := config.UpdateProgress(numberOfSess, workMin*numberOfSess); err != nil {
		fmt.Fprintf(os.Stderr, "%s⚠️  Failed to update goals progress: %v%s\n", theme.WarningColor, err, theme.TextColor)
	}

	// Show completion message and summary
	fmt.Printf("\n%s🎉 Pomodoro complete! Great job!%s\n", theme.SuccessColor, theme.TextColor)
	fmt.Printf("%s📊 Sessions completed: %d%s\n", theme.HighlightColor, numberOfSess, theme.TextColor)
	fmt.Printf("%s⏰ Total focus time: %.0f minutes%s\n", theme.HighlightColor, totalWorkTime.Minutes(), theme.TextColor)

	// Final notification
	if err := logs.ShowNotification("Pomodoro Complete!", "Great job on completing all your sessions!"); err != nil {
		fmt.Fprintf(os.Stderr, "%s⚠️  Error showing notification: %v%s\n", theme.WarningColor, err, theme.TextColor)
	}
	if err := logs.PlaySound("work_end"); err != nil {
		fmt.Fprintf(os.Stderr, "%s⚠️  Error playing sound: %v%s\n", theme.WarningColor, err, theme.TextColor)
	}

	return true
}

// getRandomMotivationalMessage returns a random motivational message
func getRandomMotivationalMessage() string {
	messages := []string{
		"🌟 Great work! Keep up the momentum!",
		"💪 You're making excellent progress!",
		"🎯 Stay focused, you're doing great!",
		"⭐ Well done on completing another session!",
		"🚀 You're crushing it! Keep going!",
		"✨ Fantastic work! Take a well-deserved break!",
		"🌈 You're getting closer to your goals!",
		"💫 Keep up the amazing work!",
		"🔥 You're on fire! Keep that focus!",
		"🌺 Excellent focus session!",
	}
	return messages[rand.Intn(len(messages))]
}

// SaveConfig saves the current configuration
func SaveConfig(workMin, breakMin, numberOfSess int) error {
	cfg := config.Config{
		WorkMinutes:  workMin,
		BreakMinutes: breakMin,
		NumSessions:  numberOfSess,
	}
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}
	fmt.Println("✅ Configuration saved successfully!")
	return nil
}
