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

func (m *mouse) Move(x, y int, displays ...display.Display) error {
	// If no displays are provided, use the primary display
	if len(displays) == 0 {
		if pd == nil {
			d, err := display.GetPrimaryDisplay()
			if err != nil {
				return err
			}
			pd = &d
		}
		displays = append(displays, *pd)
	}

	// Get the virtual screen bounds
	if vs == nil {
		vsp, err := display.DetectVirtualScreen()
		if err != nil {
			return err
		}
		vs = &vsp
	}

	// Iterate through the provided displays
	for _, d := range displays {
		// Calculate the absolute position relative to the display
		absoluteX := d.X + int32(x)
		absoluteY := d.Y + int32(y)

		// Validate the coordinates against the virtual screen bounds
		if absoluteX < vs.Left || absoluteX > vs.Right ||
			absoluteY > vs.Top || absoluteY < vs.Bottom {
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
	}

	return nil
}
