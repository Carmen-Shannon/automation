package keyboard

import "automation/device/keyboard/key_codes"

type keyboardPressOption struct {
	KeyCode  key_codes.KeyCode
	Duration int
}

type KeyboardPressOption func(*keyboardPressOption)

func KeyCodeOpt(keyCode key_codes.KeyCode) KeyboardPressOption {
	return func(opt *keyboardPressOption) {
		opt.KeyCode = keyCode
	}
}

func DurationOpt(duration int) KeyboardPressOption {
	return func(opt *keyboardPressOption) {
		opt.Duration = duration
	}
}
