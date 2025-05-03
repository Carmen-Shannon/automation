//go:build linux
// +build linux

package mouse

import (
	"automation/device/display"
	"automation/tools/linux"
	"errors"
	"fmt"
)

func Init() *mouse {
	var m mouse

	x, y, err := linux.ExecuteXdotoolGetMousePosition()
	if err != nil {
		fmt.Println("failed to get the current mouse position: ", err.Error())
		return &mouse{x: 0, y: 0}
	}

	// Initialize the mouse struct with the current position
	m.x = x
	m.y = y
	return &m
}

func (m *mouse) Move(options ...MouseMoveOption) error {
	moveOptions := &mouseMoveOption{ToX: 0, ToY: 0}
	// If no displays are provided, use the primary display
	for _, opt := range options {
		opt(moveOptions)
	}
	// Get the virtual screen bounds
	if vs == nil {
		vs = display.Init()
	}
	// default to primary display if no options are provided
	if moveOptions.Display == nil {
		if pd == nil {
			d, err := vs.GetPrimaryDisplay()
			if err != nil {
				return err
			}
			pd = &d
		}
		moveOptions.Display = pd
	}

	absoluteX := moveOptions.Display.X + int32(moveOptions.ToX)
	absoluteY := moveOptions.Display.Y + int32(moveOptions.ToY)

	// Validate the coordinates against the virtual screen bounds
	if absoluteX < vs.GetLeft() || absoluteX > vs.GetRight() ||
		absoluteY < vs.GetTop() || absoluteY > vs.GetBottom() {
		return errors.New("coordinates are outside the virtual screen bounds for display")
	}

	// Execute xdotool to move the mouse
	err := linux.ExecuteXdotoolMouseMove(absoluteX, absoluteY)
	if err != nil {
		return fmt.Errorf("failed to move mouse: %w", err)
	}

	// Update the internal mouse position
	m.x = absoluteX
	m.y = absoluteY

	return nil
}

func (m *mouse) Click(options ...MouseClickOption) error {
	clickOptions := &mouseClickOption{}
	for _, opt := range options {
		opt(clickOptions)
	}
	// default to left click if no options are provided
	if !clickOptions.Left && !clickOptions.Right && !clickOptions.Middle {
		clickOptions.Left = true
	}

	// Perform the click(s) based on the options
	if clickOptions.Left {
		err := linux.ExecuteXdotoolClick(1, clickOptions.Duration)
		if err != nil {
			return fmt.Errorf("failed to perform left click: %w", err)
		}
	}

	if clickOptions.Right {
		err := linux.ExecuteXdotoolClick(3, clickOptions.Duration)
		if err != nil {
			return fmt.Errorf("failed to perform right click: %w", err)
		}
	}

	if clickOptions.Middle {
		err := linux.ExecuteXdotoolClick(2, clickOptions.Duration)
		if err != nil {
			return fmt.Errorf("failed to perform middle click: %w", err)
		}
	}

	return nil
}
