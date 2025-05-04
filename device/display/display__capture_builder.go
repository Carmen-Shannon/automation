package display

import "slices"

type displayCaptureOption struct {
	Displays []Display
	BitCount int      // acceptable values: 1, 4, 8, 16, 24, 32
	Bounds   [4]int32 // left, right, top, bottom bounds for the capture area
}

type DisplayCaptureOption func(*displayCaptureOption)

func DisplaysOpt(displays []Display) DisplayCaptureOption {
	return func(opt *displayCaptureOption) {
		opt.Displays = displays
	}
}

var validBitCounts = []int{1, 4, 8, 16, 24, 32}

func BitCountOpt(bitCount int) DisplayCaptureOption {
	if !slices.Contains(validBitCounts, bitCount) {
		return func(opt *displayCaptureOption) {}
	}
	return func(opt *displayCaptureOption) {
		if !slices.Contains(validBitCounts, bitCount) {
			opt.BitCount = 24 // Default to 24 bits per pixel if not valid input
		} else {
			opt.BitCount = bitCount
		}
	}
}

func BoundsOpt(bounds [4]int32) DisplayCaptureOption {
	return func(opt *displayCaptureOption) {
		opt.Bounds = bounds
	}
}
