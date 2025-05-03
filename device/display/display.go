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

type virtualScreen struct {
	Left     int32
	Right    int32
	Top      int32
	Bottom   int32
	Displays []Display
}

type VirtualScreen interface {
	// DetectDisplays detects all displays connected to the system and returns a slice of display structs.
	// It also modifies the virtual screen Displays field to include the detected displays.
	// If no displays are found, it returns an error.
	//
	// Returns:
	//   - []Display: A slice of Display structs representing the detected displays.
	//   - error: An error if the detection fails or no displays are found.
	DetectDisplays() ([]Display, error)

	// GetPrimaryDisplay retrieves the primary display from the virtual screen.
	// If no primary display is found, it returns an error.
	//
	// Returns:
	//   - Display: The primary display struct.
	//   - error: An error if no primary display is found.
	GetPrimaryDisplay() (Display, error)

	// Displays returns a slice of all displays connected to the system.
	// Returns:
	//   - []Display: A slice of Display structs representing all connected displays.
	GetDisplays() []Display

	// Left returns the left bound of the virtual screen.
	// Returns:
	//   - int32: The left bound of the virtual screen.
	GetLeft() int32

	// Right returns the right bound of the virtual screen.
	// Returns:
	//   - int32: The right bound of the virtual screen.
	GetRight() int32

	// Top returns the top bound of the virtual screen.
	// Returns:
	//   - int32: The top bound of the virtual screen.
	GetTop() int32

	// Bottom returns the bottom bound of the virtual screen.
	// Returns:
	//   - int32: The bottom bound of the virtual screen.
	GetBottom() int32
}

var _ VirtualScreen = (*virtualScreen)(nil) // compile-time check to ensure that virtualScreen implements VirtualScreen

func (vs *virtualScreen) GetPrimaryDisplay() (Display, error) {
	displays := vs.Displays

	if displays == nil {
		displays, err := vs.DetectDisplays()
		if err != nil || len(displays) == 0 {
			return Display{}, err
		}
	}
	for _, display := range displays {
		if display.Primary {
			return display, nil
		}
	}
	return Display{}, errors.New("no primary display found")
}

func (vs *virtualScreen) GetDisplays() []Display {
	return vs.Displays
}

func (vs *virtualScreen) GetLeft() int32 {
	return vs.Left
}

func (vs *virtualScreen) GetRight() int32 {
	return vs.Right
}

func (vs *virtualScreen) GetTop() int32 {
	return vs.Top
}

func (vs *virtualScreen) GetBottom() int32 {
	return vs.Bottom
}
