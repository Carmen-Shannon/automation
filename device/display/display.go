package display

import "errors"

type Display struct {
	X           int32
	Y           int32
	Width       int
	Height      int
	RefreshRate float32
	Primary     bool
}

type VirtualScreen struct {
	Left     int32
	Right    int32
	Top      int32
	Bottom   int32
	Displays []Display
}

func GetPrimaryDisplay() (Display, error) {
	displays, err := DetectDisplays()
	if err != nil || len(displays) == 0 {
		return Display{}, err
	}
	for _, display := range displays {
		if display.Primary {
			return display, nil
		}
	}
	return Display{}, errors.New("no primary display found")
}
