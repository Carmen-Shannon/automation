//go:build linux
// +build linux

package linux

import (
	"fmt"
	"os/exec"
	"time"
)

func ExecuteXrandr() ([]byte, error) {
	return exec.Command("xrandr", "--query").Output()
}

func ExecuteXdotoolMouseMove(x, y int32) error {
	err := exec.Command("xdotool", "mousemove", fmt.Sprintf("%d", x), fmt.Sprintf("%d", y)).Run()
	if err != nil {
		return err
	}
	return nil
}

func ExecuteXdotoolGetMousePosition() (int32, int32, error) {
	cmd := exec.Command("xdotool", "getmouselocation")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get mouse position: %w", err)
	}

	var x, y int32
	_, err = fmt.Sscanf(string(output), "x:%d y:%d", &x, &y)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse mouse position: %w", err)
	}

	return x, y, nil
}

func ExecuteXdotoolClick(button int, duration int) error {
	// Simulate the button press
	err := exec.Command("xdotool", "mousedown", fmt.Sprintf("%d", button)).Run()
	if err != nil {
		return fmt.Errorf("failed to press mouse button %d: %w", button, err)
	}

	// Add delay if duration is specified
	if duration > 0 {
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}

	// Simulate the button release
	err = exec.Command("xdotool", "mouseup", fmt.Sprintf("%d", button)).Run()
	if err != nil {
		return fmt.Errorf("failed to release mouse button %d: %w", button, err)
	}

	return nil
}
