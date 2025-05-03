//go:build linux
// +build linux

package keyboard

import (
	"errors"
	"slices"
	"strings"
	"time"

	"automation/tools/linux"
)

func KeyPress(options ...KeyboardPressOption) error {
	kbpOpt := &keyboardPressOption{}
	for _, opt := range options {
		opt(kbpOpt)
	}
	if slices.Contains(kbpOpt.KeyCodes, 0) {
		return errors.New("invalid key code entered")
	}

	action := []string{}
	for _, keyCode := range kbpOpt.KeyCodes {
		keySym := linux.XKeysymToString(uint32(keyCode))
		action = append(action, keySym)
	}

	actionStr := strings.Join(action, "+")
	err := linux.ExecuteXdotoolKeyDown(actionStr)
	if err != nil {
		return err
	}

	if kbpOpt.Duration > 0 {
		time.Sleep(time.Duration(kbpOpt.Duration) * time.Millisecond)
	}

	err = linux.ExecuteXdotoolKeyUp(actionStr)
	if err != nil {
		return err
	}
	return nil
}
