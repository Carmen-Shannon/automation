package display

type displayCaptureOption struct {
	Displays []Display
}

type DisplayCaptureOption func(*displayCaptureOption)

func DisplaysOpt(displays []Display) DisplayCaptureOption {
	return func(opt *displayCaptureOption) {
		opt.Displays = displays
	}
}