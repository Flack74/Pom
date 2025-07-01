package logs

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Common sound files on different Linux distributions

// ShowNotification displays a system notification
func ShowNotification(title, message string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		// Try notify-send first
		if err = exec.Command("notify-send", title, message).Run(); err != nil {
			// If notify-send fails, try zenity
			if zenityErr := exec.Command("zenity", "--notification", "--text", fmt.Sprintf("%s: %s", title, message)).Run(); zenityErr != nil {
				return fmt.Errorf("notification failed: notify-send error: %v, zenity error: %v", err, zenityErr)
			}
		}

	case "darwin":
		script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
		if err = exec.Command("osascript", "-e", script).Run(); err != nil {
			return fmt.Errorf("notification failed: %v", err)
		}

	default:
		// For other OS, print to stderr
		fmt.Fprintf(os.Stderr, "\a%s: %s\n", title, message)
	}

	return nil
}

// PlaySound plays a notification sound
func PlaySound(soundType string) error {
	switch runtime.GOOS {
	case "linux":
		// Try canberra-gtk-play first (best option)
		soundID := "complete"
		if soundType == "break_end" {
			soundID = "dialog-warning"
		}
		if err := exec.Command("canberra-gtk-play", "--id", soundID, "--description", "Pomodoro notification").Run(); err == nil {
			return nil
		}

		// Try paplay as fallback
		soundFile := "/usr/share/sounds/freedesktop/stereo/complete.oga"
		if err := exec.Command("paplay", soundFile).Run(); err == nil {
			return nil
		}

		// Try aplay as last resort
		if err := exec.Command("aplay", "-q", "-f", "cd", "/dev/zero", "trim", "0", "0.1").Run(); err == nil {
			return nil
		}

		// If all else fails, use terminal bell
		fmt.Print("\a")
		return nil

	case "darwin":
		var message string
		switch soundType {
		case "work_end":
			message = "Work session complete"
		case "break_end":
			message = "Break time over"
		default:
			message = "Notification"
		}

		if err := exec.Command("say", message).Run(); err != nil {
			// If speech fails, try system alert sound
			if err := exec.Command("afplay", "/System/Library/Sounds/Glass.aiff").Run(); err != nil {
				fmt.Print("\a") // Fallback to terminal bell
			}
		}
		return nil

	default:
		// For other OS, just use the terminal bell
		fmt.Print("\a")
		return nil
	}
}
