package mouse

import "automation/device/display"

type mouse struct {
	x int32
	y int32
}

var (
	// the virtual screen to use for mouse movement, cached on the first call to Move so it isn't initialized on every call
	vs *display.VirtualScreen
	// the primary display to use for mouse movement, cached on the first call to Move so it isn't initialized on every call
	pd *display.Display
)

type Mouse interface {
	// Move moves the mouse to the specified coordinates on the given displays.
	// If no displays are provided, it defaults to the primary display - this is OS dependent.
	//
	// If the coordinates are outside of the display area bounds on any given display, then the function will return an error.
	//
	// Parameters:
	//   - x: The x-coordinate to move the mouse to.
	//   - y: The y-coordinate to move the mouse to.
	//   - display: Optional list of displays to consider for the move operation.
	//     If no displays are provided, the primary display is used.
	//
	// Returns:
	//   - error: An error if the move operation fails, otherwise nil.
	Move(x, y int, displays ...display.Display) error

	// GetCurrentPosition retrieves the current position of the mouse cursor.
	// The position is returned as a tuple of (x, y) coordinates.
	// If the position cannot be determined, (0, 0) is returned
	//
	// Returns:
	//   - x: The current x-coordinate of the mouse cursor.
	//   - y: The current y-coordinate of the mouse cursor.
	GetCurrentPosition() (int, int)
}

var _ Mouse = (*mouse)(nil) // compile-time check to ensure that mouse implements Mouse

func (m *mouse) GetCurrentPosition() (int, int) {
	return int(m.x), int(m.y)
}
