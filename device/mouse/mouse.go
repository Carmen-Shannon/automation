package mouse

import (
	"automation/device/display"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type mouse struct {
	mu   sync.Mutex
	done chan struct{}
	x    int32
	y    int32
}

var (
	// the virtual screen to use for mouse movement, cached on the first call to Move so it isn't initialized on every call
	vs display.VirtualScreen
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
	//   - options: Optional parameters for the mouse movement, such as display and velocity.
	//
	// Returns:
	//   - error: An error if the move operation fails, otherwise nil.
	Move(x, y int32, options ...MouseMoveOption) error

	// Click performs a mouse click at the current mouse position.
	// The default click is a left click with no duration, an instant click down and up.
	// To modify this behavior, you can pass in a list of MouseClickOptions to customize the click action.
	//
	// Parameters:
	//   - options: Optional parameters for the mouse click, such as button type and click count.
	//
	// Returns:
	//   - error: An error if the click operation fails, otherwise nil.
	Click(options ...MouseClickOption) error

	// GetCurrentPosition retrieves the current position of the mouse cursor.
	// The position is returned as a tuple of (x, y) coordinates.
	// If the position cannot be determined, (0, 0) is returned.
	// The Init function should be called prior to calling this function, otherwise it will always return (0, 0)
	//
	// Returns:
	//   - x: The current x-coordinate of the mouse cursor.
	//   - y: The current y-coordinate of the mouse cursor.
	GetCurrentPosition() (int, int)
}

var _ Mouse = (*mouse)(nil) // compile-time check to ensure that mouse implements Mouse

func Init() *mouse {
	var m mouse
	m.mu = sync.Mutex{}
	m.done = nil

	x, y, err := doGetMousePosition()
	if err != nil {
		return &m
	}

	m.x = x
	m.y = y
	return &m
}

func (m *mouse) Click(options ...MouseClickOption) error {
	clickOptions := &mouseClickOption{}
	for _, opt := range options {
		opt(clickOptions)
	}
	// default to left click if no options are provided
	if !clickOptions.Left && !clickOptions.Right && !clickOptions.Middle {
		clickOptions.Left = true
	}

	// Perform the click(s) based on the options
	if clickOptions.Left {
		err := m.doMouseClick(1, clickOptions.Duration)
		if err != nil {
			return fmt.Errorf("failed to perform left click: %w", err)
		}
	}

	if clickOptions.Right {
		err := m.doMouseClick(3, clickOptions.Duration)
		if err != nil {
			return fmt.Errorf("failed to perform right click: %w", err)
		}
	}

	if clickOptions.Middle {
		err := m.doMouseClick(2, clickOptions.Duration)
		if err != nil {
			return fmt.Errorf("failed to perform middle click: %w", err)
		}
	}

	return nil
}

func (m *mouse) GetCurrentPosition() (int, int) {
	return int(m.x), int(m.y)
}

// Move moves the mouse to the specified coordinates on the given displays.
func (m *mouse) Move(x, y int32, options ...MouseMoveOption) error {
	moveOptions := &mouseMoveOption{}
	for _, opt := range options {
		opt(moveOptions)
	}
	if moveOptions.Done != nil {
		m.done = moveOptions.Done
		defer func() {
			close(moveOptions.Done)
		}()
	}

	if vs == nil {
		vs = display.NewVirtualScreen()
	}
	if moveOptions.Display == nil {
		if pd == nil {
			d, err := vs.GetPrimaryDisplay()
			if err != nil {
				return err
			}
			pd = &d
		}
		moveOptions.Display = pd
	}

	absoluteX := moveOptions.Display.X + x
	absoluteY := moveOptions.Display.Y + y

	// Validate the coordinates against the virtual screen bounds
	if (absoluteX < vs.GetLeft() || absoluteX > vs.GetRight()) ||
		(absoluteY > vs.GetTop() || absoluteY < vs.GetBottom()) {
		return errors.New("coordinates are outside the virtual screen bounds for display")
	}

	// If velocity is not set or is zero, perform the movement in one step
	if moveOptions.Velocity <= 0 {
		err := m.doMouseMove(absoluteX, absoluteY)
		if err != nil {
			return err
		}
		m.x = absoluteX
		m.y = absoluteY
		return nil
	} else {
		err := m.moveWithVelocity(absoluteX, absoluteY, moveOptions.Velocity, moveOptions.Jitter, moveOptions.Display)
		if err != nil {
			return err
		}
		m.x = absoluteX
		m.y = absoluteY
		return nil
	}
}

func (m *mouse) moveWithVelocity(x, y int32, velocity, jitter int, disp *display.Display) error {
	startX, startY := m.x, m.y
	deltaX := float64(x - startX)
	deltaY := float64(y - startY)
	distance := math.Sqrt(deltaX*deltaX + deltaY*deltaY)
	refreshRate := 120.0 // Defaulting to 120 FPS for mouse movement
	if disp != nil {
		refreshRate = math.Max(120.0, float64(disp.RefreshRate))
	} else if pd != nil {
		refreshRate = math.Max(120.0, float64(pd.RefreshRate))
	}

	steps := int(math.Ceil(distance / float64(velocity) * refreshRate)) // Number of steps based on refresh rate
	stepDuration := time.Second / time.Duration(refreshRate)            // Base time per step

	// Create a ticker for consistent timing
	ticker := time.NewTicker(stepDuration)
	defer ticker.Stop() // Ensure the ticker is stopped when the function exits

	// Define control points for the parabolic curve
	controlX := float64(startX) + deltaX/2 + float64(rand.Intn(2*jitter+1)-jitter)
	controlY := float64(startY) + deltaY/2 + float64(rand.Intn(2*jitter+1)-jitter)

	m.mu.Lock()
	defer m.mu.Unlock()

	currentVelocity := float64(velocity) // Start with the base velocity

	for i := 1; i <= steps; i++ {
		select {
		case <-ticker.C:
			// Adjust velocity based on jitter
			if jitter > 0 {
				velocityFluctuation := float64(rand.Intn(2*jitter+1)-jitter) * 0.1    // Fluctuation scaled by jitter
				currentVelocity = math.Max(10, float64(velocity)+velocityFluctuation) // Ensure velocity doesn't drop too low
			}

			// Recalculate step duration based on the new velocity
			stepDuration = time.Second / time.Duration(refreshRate*currentVelocity/float64(velocity))
			ticker.Reset(stepDuration)

			// Calculate the t parameter (progress along the curve)
			t := float64(i) / float64(steps)

			// Apply the easing function to t
			easedT := 3*t*t - 2*t*t*t

			// Calculate the parabolic curve point using the quadratic bezier formula
			currentX := (1-easedT)*(1-easedT)*float64(startX) + 2*(1-easedT)*easedT*controlX + easedT*easedT*float64(x)
			currentY := (1-easedT)*(1-easedT)*float64(startY) + 2*(1-easedT)*easedT*controlY + easedT*easedT*float64(y)

			// Move the mouse to the calculated position
			err := m.doMouseMove(int32(currentX), int32(currentY))
			if err != nil {
				return fmt.Errorf("failed to move mouse: %w", err)
			}
		}
	}

	// Ensure the final position is set
	err := m.doMouseMove(x, y)
	if err != nil {
		return fmt.Errorf("failed to move mouse to final position: %w", err)
	}

	m.x = x
	m.y = y
	return nil
}
