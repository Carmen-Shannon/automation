//go:build linux
// +build linux

package linux

import (
	"fmt"
	"os/exec"
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
