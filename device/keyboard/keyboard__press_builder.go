package keyboard

import "automation/device/keyboard/key_codes"

type keyboardPressOption struct {
	KeyCodes  []key_codes.KeyCode
	Duration int
}

type KeyboardPressOption func(*keyboardPressOption)

func KeyCodeOpt(keyCodes []key_codes.KeyCode) KeyboardPressOption {
	return func(opt *keyboardPressOption) {
		opt.KeyCodes = keyCodes
	}
}

func DurationOpt(duration int) KeyboardPressOption {
	return func(opt *keyboardPressOption) {
		opt.Duration = duration
	}
}
