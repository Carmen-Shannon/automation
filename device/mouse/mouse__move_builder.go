package mouse

import "automation/device/display"

type mouseMoveOption struct {
	ToX       int32
	ToY       int32
	Display *display.Display
	Velocity *uint8
}

type MouseMoveOption func(*mouseMoveOption)

func ToXOpt(toX int32) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.ToX = toX
	}
}

func ToYOpt(toY int32) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.ToY = toY
	}
}

func DisplayOpt(display *display.Display) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.Display = display
	}
}

func VelocityOpt(velocity *uint8) MouseMoveOption {
	return func(opt *mouseMoveOption) {
		opt.Velocity = velocity
	}
}
