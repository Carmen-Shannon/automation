//go:build windows
// +build windows

package keyboard

import (
	"errors"
	"fmt"
	"slices"
	"time"

	windows "github.com/Carmen-Shannon/automation/tools/_windows"
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
