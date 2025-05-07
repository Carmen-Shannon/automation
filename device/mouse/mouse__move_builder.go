package mouse

import "github.com/Carmen-Shannon/automation/device/display"

type mouseMoveOption struct {
	Velocity int
	Jitter   int
	Done     chan struct{}
	Display  *display.Display
}

type MouseMoveOption func(*mouseMoveOption)

// JitterOpt is the option to control mouse movement jitter.
//
// Parameters:
//   - jitter: The amount of jitter to apply to the mouse movement. This is a random value added to the x and y coordinates of the mouse movement.
func JitterOpt(jitter int) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.Jitter = jitter
	}
}

// DisplayOpt is the option to specify the display for mouse movement.
//
// Parameters:
//   - display: The display to use for mouse movement. This is useful for multi-display setups where you want to move the mouse on a specific display.
func DisplayOpt(display *display.Display) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.Display = display
	}
}

// VelocityOpt is the option to control mouse movement velocity.
//
// Parameters:
//   - velocity: The speed of the mouse movement. This is a value that determines how fast the mouse moves from one point to another.
//		Omit this field or set it to 0 for instant movement.
func VelocityOpt(velocity int) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.Velocity = velocity
	}
}

// DoneSignalOpt is the option to specify a done signal channel for mouse movement.
//
// Parameters:
//   - done: A channel that signals when the mouse movement is done. This is useful for synchronizing mouse movements with other operations.
func DoneSignalOpt(done chan struct{}) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.Done = done
	}
}
