//go:build windows
// +build windows

package keyboard

import (
	"automation/tools/windows"
	"errors"
	"fmt"
	"time"
)

func KeyPress(options ...KeyboardPressOption) error {
	kbpOpt := &keyboardPressOption{}
	for _, opt := range options {
		opt(kbpOpt)
	}
	if kbpOpt.KeyCode == 0 {
		return errors.New("invalid key code entered")
	}

	ret, _, err := windows.KeybdEvent.Call(uintptr(kbpOpt.KeyCode), 0, 0, 0)
	if ret == 0 {
		return fmt.Errorf("failed to send key event: %v", err)
	}

	if kbpOpt.Duration > 0 {
		time.Sleep(time.Duration(kbpOpt.Duration) * time.Millisecond)
	}

	ret, _, err = windows.KeybdEvent.Call(uintptr(kbpOpt.KeyCode), 0, 2, 0)
	if ret == 0 {
		return fmt.Errorf("failed to send key event: %v", err)
	}
	return nil
}
