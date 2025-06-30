package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"pom/logs"
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

// sendNotification sends a system notification based on the OS
func sendNotification(title, message string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("notify-send", title, message).Run()
		exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga").Run()
	case "darwin":
		exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)).Run()
		exec.Command("say", message).Run()
	default:
		fmt.Print("\a") // Default beep for other OS
	}
}

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

// countdown displays a live countdown timer
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
			minutes := int(remaining.Minutes())
			seconds := int(remaining.Seconds()) % 60
			fmt.Printf("\r%s%s time remaining: %02d:%02d%s", color, label, minutes, seconds, colorReset)
		}
	}
}

func StartPomodoro(workMin, breakMin, numberOfSess int) {
	counter := 0
	totalWorkTime := time.Duration(0)
	startTime := time.Now()

	fmt.Printf("%s🎯 Starting Pomodoro Timer%s\n\n", colorPurple, colorReset)

	// Set up signal handling for cleanup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Channels for pause/resume functionality
	pauseChan := make(chan struct{})
	resumeChan := make(chan struct{})
	timerState := stateRunning

	// Start user input handler
	go handleUserInput(&timerState, pauseChan, resumeChan)

	// Create a channel to track completion
	completed := make(chan bool, 1)
	go func() {
		for counter < numberOfSess {
			work := time.Duration(workMin) * time.Minute
			breakTime := time.Duration(breakMin) * time.Minute

			fmt.Printf("%s📚 Session %d/%d%s\n", colorBlue, counter+1, numberOfSess, colorReset)

			// Work period
			fmt.Printf("%s⏱️  Starting work session...%s\n", colorGreen, colorReset)
			if !countdown(work, "Focus", colorGreen, &timerState, pauseChan, resumeChan) {
				completed <- false
				return
			}
			sendNotification("Pomodoro", "Work session complete! Time for a break!")
			totalWorkTime += work

			if counter < numberOfSess-1 { // Skip break after last session
				// Break period
				fmt.Printf("\n%s☕ Time for a break!%s\n", colorYellow, colorReset)
				if !countdown(breakTime, "Break", colorYellow, &timerState, pauseChan, resumeChan) {
					completed <- false
					return
				}
				sendNotification("Pomodoro", "Break time is over! Let's focus!")
				fmt.Println()
			}

			counter++
		}
		completed <- true
	}()

	// Wait for either completion or interruption
	var isCompleted bool
	select {
	case <-sigChan:
		isCompleted = false
	case isCompleted = <-completed:
	}

	// Log the session
	endTime := time.Now()
	if err := logs.LogSession(workMin, breakMin, numberOfSess, startTime, endTime, isCompleted); err != nil {
		fmt.Fprintf(os.Stderr, "%s⚠️  Failed to log session: %v%s\n", colorRed, err, colorReset)
	}

	if isCompleted {
		// Session summary
		fmt.Printf("\n%s🎉 Pomodoro complete! Here's your session summary:%s\n", colorPurple, colorReset)
		fmt.Printf("%s📊 Total Sessions: %d%s\n", colorBlue, counter, colorReset)
		fmt.Printf("%s⏰ Total Focus Time: %.0f minutes%s\n", colorGreen, totalWorkTime.Minutes(), colorReset)
		fmt.Printf("%s🌟 Great job staying focused!%s\n", colorYellow, colorReset)
		sendNotification("Pomodoro", "All sessions completed! Great job!")
	} else {
		fmt.Printf("\n%s⚠️  Session interrupted! Progress has been saved.%s\n", colorRed, colorReset)
		sendNotification("Pomodoro", "Session interrupted! Progress saved.")
	}
}
