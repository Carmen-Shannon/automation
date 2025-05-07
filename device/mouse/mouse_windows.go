//go:build windows
// +build windows

package mouse

import (
	"errors"
	"time"
	"unsafe"

	windows "github.com/Carmen-Shannon/automation/tools/_windows"
)

// doGetMousePosition retrieves the current mouse position on the screen.
// It returns the x and y coordinates of the mouse cursor.
// If the retrieval fails, it returns an error.
// This function is specific to Windows OS and uses the Windows API to get the mouse position.
//
// Returns:
//   - (int32, int32): The x and y coordinates of the mouse cursor.
//   - error: An error if the retrieval fails, otherwise nil.
func doGetMousePosition() (int32, int32, error) {
	var p struct {
		x int32
		y int32
	}
	ret, _, err := windows.GetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	if ret == 0 {
		return 0, 0, err
	}

	return p.x, p.y, nil
}

// doMouseClick performs a mouse click at the current mouse position.
// It accepts the button to click (1 for left, 2 for middle, 3 for right) and an optional duration for the click.
// The function uses the Windows API to simulate the mouse click event.
// It first simulates a mouse button down event, waits for the specified duration (if any), and then simulates a mouse button up event.
//
// Parameters:
//   - btn: The button to click (1 for left, 2 for middle, 3 for right).
//   - duration: The duration to hold the button down in milliseconds. If 0, it will be an instant click.
//
// Returns:
//   - error: An error if the click operation fails, otherwise nil.
func (m *mouse) doMouseClick(btn int, duration int) error {
	var downFlags, upFlags uintptr
	if btn == 1 {
		downFlags |= windows.MOUSEEVENTF_LEFTDOWN
		upFlags |= windows.MOUSEEVENTF_LEFTUP
	}
	if btn == 3 {
		downFlags |= windows.MOUSEEVENTF_RIGHTDOWN
		upFlags |= windows.MOUSEEVENTF_RIGHTUP
	}
	if btn == 2 {
		downFlags |= windows.MOUSEEVENTF_MIDDLEDOWN
		upFlags |= windows.MOUSEEVENTF_MIDDLEUP
	}

	windows.MouseEvent.Call(downFlags, 0, 0, 0, 0)

	if duration > 0 {
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}

	windows.MouseEvent.Call(upFlags, 0, 0, 0, 0)
	return nil
}

// doMouseMove moves the mouse cursor to the specified x and y coordinates on the screen.
// It uses the Windows API to set the cursor position. The coordinates are relative to the screen, not the window.
//
// Parameters:
//   - x: The x-coordinate to move the mouse to.
//   - y: The y-coordinate to move the mouse to.
func (m *mouse) doMouseMove(x, y int32) error {
	ret, _, err := windows.SetCursorPos.Call(uintptr(x), uintptr(y))
	if ret == 0 {
		return errors.New("failed to move the mouse: " + err.Error())
	}
	return nil
}
