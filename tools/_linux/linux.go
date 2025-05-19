//go:build linux
// +build linux

package linux

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <X11/keysym.h>
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

// XKeysymToString converts an X KeySym value to its string representation.
func XKeysymToString(keysym uint32) string {
	// Call the XKeysymToString function from the X11 library
	cStr := C.XKeysymToString(C.KeySym(keysym))
	if cStr == nil {
		return ""
	}
	// Convert the C string to a Go string
	return C.GoString(cStr)
}

func ExecuteXrandr() ([]byte, error) {
	return exec.Command("xrandr", "--query").Output()
}

func ExecuteXdotoolMouseMove(x, y int32) error {
	err := exec.Command("xdotool", "mousemove", fmt.Sprintf("%d", x), fmt.Sprintf("%d", y)).Run()
	if err != nil {
		return err
	}
	return nil
}

func ExecuteXdotoolGetMousePosition() (int32, int32, error) {
	cmd := exec.Command("xdotool", "getmouselocation")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get mouse position: %w", err)
	}

	var x, y int32
	_, err = fmt.Sscanf(string(output), "x:%d y:%d", &x, &y)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse mouse position: %w", err)
	}

	return x, y, nil
}

func ExecuteXdotoolClick(button int, duration int) error {
	// Simulate the button press
	if duration == 0 {
		err := exec.Command("xdotool", "click", fmt.Sprintf("%d", button)).Run()
		if err != nil {
			return fmt.Errorf("failed to click mouse button %d: %w", button, err)
		}
		return nil
	}
	
	err := exec.Command("xdotool", "mousedown", fmt.Sprintf("%d", button)).Run()
	if err != nil {
		return fmt.Errorf("failed to press mouse button %d: %w", button, err)
	}

	time.Sleep(time.Duration(duration) * time.Millisecond)

	// Simulate the button release
	err = exec.Command("xdotool", "mouseup", fmt.Sprintf("%d", button)).Run()
	if err != nil {
		return fmt.Errorf("failed to release mouse button %d: %w", button, err)
	}

	return nil
}

func ExecuteXdotoolKeyDown(keySym string) error {
	return exec.Command("xdotool", "keydown", keySym).Run()
}

func ExecuteXdotoolKeyUp(keySym string) error {
	return exec.Command("xdotool", "keyup", keySym).Run()
}

func ExecuteXwd(x, y, width, height int) ([]byte, error) {
	// Construct the `xwd` command
	cmd := exec.Command("xwd", "-root", "-silent", "-geometry", fmt.Sprintf("%dx%d+%d+%d", width, height, x, y))

	// Capture the output of the command
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute xwd: %w", err)
	}

	return out.Bytes(), nil
}
