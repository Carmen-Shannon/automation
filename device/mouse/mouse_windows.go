//go:build windows
// +build windows

package mouse

import (
	"automation/device/display"
	"automation/tools/windows"
	"errors"
	"fmt"
	"unsafe"
)

func Init() *mouse {
	var m mouse

	ret, _, err := windows.GetCursorPos.Call(uintptr(unsafe.Pointer(&m)))
	if ret == 0 {
		fmt.Println("failed to get the current mouse position: ", err.Error())
		return &mouse{x: 0, y: 0}
	}

	return &m
}

// Move moves the mouse to the specified coordinates on the given displays.
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

	for _, d := range displays {
		// Calculate the absolute position relative to the display
		absoluteX := d.X + int32(x)
		absoluteY := d.Y + int32(y)

		// Validate the coordinates against the virtual screen bounds
		if (absoluteX < vs.Left || absoluteX > vs.Right) ||
			(absoluteY > vs.Top || absoluteY < vs.Bottom) {
			fmt.Printf("absoluteX: %d, absoluteY: %d, vs.Left: %d, vs.Right: %d, vs.Top: %d, vs.Bottom: %d\n", absoluteX, absoluteY, vs.Left, vs.Right, vs.Top, vs.Bottom)
			return errors.New("coordinates are outside the virtual screen bounds for display")
		}

		// Call SetCursorPos to move the mouse
		ret, _, err := windows.SetCursorPos.Call(uintptr(absoluteX), uintptr(absoluteY))
		if ret == 0 {
			return errors.New("failed to move the mouse: " + err.Error())
		}

		// Update the internal mouse position
		m.x = absoluteX
		m.y = absoluteY
	}

	return nil
}
