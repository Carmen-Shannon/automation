//go:build linux
// +build linux

package mouse

import (
	linux "automation/tools/_linux"
	"fmt"
)

func doGetMousePosition() (int32, int32, error) {
	x, y, err := linux.ExecuteXdotoolGetMousePosition()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get mouse position: %w", err)
	}
	return x, y, nil
}

func (m *mouse) doMouseMove(x, y int32) error {
	err := linux.ExecuteXdotoolMouseMove(x, y)
	if err != nil {
		return fmt.Errorf("failed to move mouse: %w", err)
	}
	return nil
}

func (m *mouse) doMouseClick(btn int, duration int) error {
	switch btn {
	case 1:
	case 2:
	case 3:
		err := linux.ExecuteXdotoolClick(btn, duration)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("invalid button: %d", btn)
	}
	return nil
}
