package keyboard

import "github.com/Carmen-Shannon/automation/device/keyboard/key_codes"

type keyboardPressOption struct {
	KeyCodes []key_codes.KeyCode
	Duration int
}

type KeyboardPressOption func(*keyboardPressOption)

// KeyCodeOpt is the option to specify the key codes for the keyboard press event.
// This works with modifiers for both windows and linux, simply add the modifier before the key code you want to modify in the slice.
//
// Parameters:
//   - keyCodes: A slice of key codes to press. This can include multiple key codes for simultaneous key presses.
//   	Example: []key_codes.KeyCode{key_codes.KeyCodeLeftShift, key_codes.KeyCodeX} will press the left shift key and the 'X' key simultaneously.
func KeyCodeOpt(keyCodes []key_codes.KeyCode) KeyboardPressOption {
	return func(opt *keyboardPressOption) {
		opt.KeyCodes = keyCodes
	}
}

// DurationOpt is the option to specify the duration for the key press event.
// This is the time in milliseconds that the key will be held down before being released.
//
// Parameters:
//   - duration: The duration to hold the key down in milliseconds. If 0, it will be an instant key press.
//   	Example: 1000 will hold the key down for 1 second before releasing it.
func DurationOpt(duration int) KeyboardPressOption {
	return func(opt *keyboardPressOption) {
		opt.Duration = duration
	}
}
