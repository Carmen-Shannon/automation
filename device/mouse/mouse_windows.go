//go:build windows
// +build windows

package mouse

import (
	"automation/device/display"
	"automation/tools/windows"
	"errors"
	"fmt"
	"time"
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
func (m *mouse) Move(options ...MouseMoveOption) error {
	moveOptions := &mouseMoveOption{ToX: 0, ToY: 0}
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
	if (absoluteX < vs.GetLeft() || absoluteX > vs.GetRight()) ||
		(absoluteY > vs.GetTop() || absoluteY < vs.GetBottom()) {
		return errors.New("coordinates are outside the virtual screen bounds for display")
	}

	ret, _, err := windows.SetCursorPos.Call(uintptr(absoluteX), uintptr(absoluteY))
	if ret == 0 {
		return errors.New("failed to move the mouse: " + err.Error())
	}

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

	// Combine all click events if multiple options are provided
	var downFlags, upFlags uintptr
	if clickOptions.Left {
		downFlags |= windows.MOUSEEVENTF_LEFTDOWN
		upFlags |= windows.MOUSEEVENTF_LEFTUP
	}
	if clickOptions.Right {
		downFlags |= windows.MOUSEEVENTF_RIGHTDOWN
		upFlags |= windows.MOUSEEVENTF_RIGHTUP
	}
	if clickOptions.Middle {
		downFlags |= windows.MOUSEEVENTF_MIDDLEDOWN
		upFlags |= windows.MOUSEEVENTF_MIDDLEUP
	}

	// Perform the click down event
	windows.MouseEvent.Call(downFlags, 0, 0, 0, 0)

	// Add delay if DurationOpt is specified
	if clickOptions.Duration > 0 {
		time.Sleep(time.Duration(clickOptions.Duration) * time.Millisecond)
	}

	// Perform the click up event
	windows.MouseEvent.Call(upFlags, 0, 0, 0, 0)

	return nil
}
