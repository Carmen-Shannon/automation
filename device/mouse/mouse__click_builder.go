package mouse

type mouseClickOption struct {
	Left     bool
	Right    bool
	Middle   bool
	Duration int
}

type MouseClickOption func(*mouseClickOption)

func LeftClickOpt() MouseClickOption {
	return func(opt *mouseClickOption) {
		opt.Left = true
	}
}

func RightClickOpt() MouseClickOption {
	return func(opt *mouseClickOption) {
		opt.Right = true
	}
}

func MiddleClickOpt() MouseClickOption {
	return func(opt *mouseClickOption) {
		opt.Middle = true
	}
}

func DurationOpt(duration int) MouseClickOption {
	return func(opt *mouseClickOption) {
		opt.Duration = duration
	}
}
