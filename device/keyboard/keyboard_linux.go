//go:build linux
// +build linux

package keyboard

import (
	"errors"
	"time"

	"automation/tools/linux"
)

func KeyPress(options ...KeyboardPressOption) error {
	kbpOpt := &keyboardPressOption{}
	for _, opt := range options {
		opt(kbpOpt)
	}
	if kbpOpt.KeyCode == 0 {
		return errors.New("invalid key code entered")
	}

	keySym := linux.KeyCodeToKeySym(kbpOpt.KeyCode)
	err := linux.ExecuteXdotoolKeyDown(keySym)
	if err != nil {
		return err
	}

	if kbpOpt.Duration > 0 {
		time.Sleep(time.Duration(kbpOpt.Duration) * time.Millisecond)
	}

	err = linux.ExecuteXdotoolKeyUp(keySym)
	if err != nil {
		return err
	}
	return nil
}
