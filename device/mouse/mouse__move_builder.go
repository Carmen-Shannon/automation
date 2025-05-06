package mouse

import "automation/device/display"

type mouseMoveOption struct {
	Velocity int
	Jitter   int
	Done     chan struct{}
	Display  *display.Display
}

type MouseMoveOption func(*mouseMoveOption)

func JitterOpt(jitter int) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.Jitter = jitter
	}
}

func DisplayOpt(display *display.Display) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.Display = display
	}
}

func VelocityOpt(velocity int) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.Velocity = velocity
	}
}

func DoneSignalOpt(done chan struct{}) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.Done = done
	}
}
