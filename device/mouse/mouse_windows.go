//go:build windows
// +build windows

package mouse

import (
	windows "github.com/Carmen-Shannon/automation/tools/_windows"
	"errors"
	"time"
	"unsafe"
)

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

func (m *mouse) doMouseMove(x, y int32) error {
	ret, _, err := windows.SetCursorPos.Call(uintptr(x), uintptr(y))
	if ret == 0 {
		return errors.New("failed to move the mouse: " + err.Error())
	}
	return nil
}
