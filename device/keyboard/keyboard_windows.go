//go:build windows
// +build windows

package keyboard

import (
	"automation/tools/windows"
	"errors"
	"fmt"
	"slices"
	"time"
)

func KeyPress(options ...KeyboardPressOption) error {
	kbpOpt := &keyboardPressOption{}
	for _, opt := range options {
		opt(kbpOpt)
	}
	if slices.Contains(kbpOpt.KeyCodes, 0) {
		return errors.New("invalid key code entered")
	}

	for _, keyCode := range kbpOpt.KeyCodes {
		ret, _, err := windows.KeybdEvent.Call(uintptr(keyCode), 0, 0, 0)
		if ret == 0 {
			return fmt.Errorf("failed to send key event: %v", err)
		}
	}

	if kbpOpt.Duration > 0 {
		time.Sleep(time.Duration(kbpOpt.Duration) * time.Millisecond)
	}

	for _, keyCode := range kbpOpt.KeyCodes {
		ret, _, err := windows.KeybdEvent.Call(uintptr(keyCode), 0, 2, 0)
		if ret == 0 {
			return fmt.Errorf("failed to send key event: %v", err)
		}
	}

	return nil
}
