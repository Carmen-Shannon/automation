//go:build linux
// +build linux

package display

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	linux "github.com/Carmen-Shannon/automation/tools/_linux"
)

func NewVirtualScreen() VirtualScreen {
	var vs virtualScreen
	displays, err := vs.DetectDisplays()
	if err != nil {
		return &vs
	} else if len(displays) == 0 {
		return &vs
	}

	left, bottom := displays[0].X, displays[0].Y
	right, top := displays[0].X+int32(displays[0].Width), displays[0].Y+int32(displays[0].Height)

	for _, d := range displays {
		if d.X < left {
			left = d.X
		}
		if d.Y < bottom {
			bottom = d.Y
		}
		if d.X+int32(d.Width) > right {
			right = d.X + int32(d.Width)
		}
		if d.Y+int32(d.Height) > top {
			top = d.Y + int32(d.Height)
		}
	}

	vs = virtualScreen{
		Left:     left,
		Right:    right,
		Top:      top,
		Bottom:   bottom,
		Displays: displays,
	}
	return &vs

}

func (vs *virtualScreen) CaptureBmp(options ...DisplayCaptureOption) ([]BMP, error) {
	return nil, nil
}

func (vs *virtualScreen) DetectDisplays() ([]Display, error) {
	// Execute the `xrandr` command to get display information
	output, err := linux.ExecuteXrandr()
	if err != nil {
		return nil, err
	}

	// Parse the output of the xrandr command
	return extractDisplaysFromXrandrOutput(string(output)), nil
}

func extractDisplaysFromXrandrOutput(output string) []Display {
	lines := strings.Split(output, "\n")
	var displays []Display
	var currentDisplay *Display

	for _, line := range lines {
		if isDisplayDetails(line) {
			var displayEntry Display
			if isPrimaryDisplay(line) {
				displayEntry.Primary = true
			}
			// checking for the connected displays example: eDP-1 connected primary 1920x1080+0+0
			// Regular expression to match the resolution format
			re := regexp.MustCompile(`\d+x\d+\+\d+\+\d+`)
			match := re.FindString(line)
			if match != "" {
				match = strings.Split(match, " ")[0]
				res := strings.Split(match, "x")
				// at this point res looks like ["1920","1080+0+-69"]
				width, _ := strconv.Atoi(res[0])
				yRes := strings.Split(res[1], "+")
				// at this point yRes looks like ["1080","0","-69"]
				height, _ := strconv.Atoi(yRes[0])
				x, _ := strconv.ParseInt(yRes[1], 10, 32)
				y, _ := strconv.ParseInt(yRes[2], 10, 32)

				displayEntry.Width = width
				displayEntry.Height = height
				displayEntry.X = int32(x)
				displayEntry.Y = int32(y)
				if x == 0 && y == 0 {
					displayEntry.Primary = true
				}
				currentDisplay = &displayEntry
			}
		} else if currentDisplay != nil && strings.Contains(line, "*+") {
			re := regexp.MustCompile(`\d+\.\d+\*\+`)
			match := re.FindString(line)
			if match != "" {
				refreshRateStr := strings.TrimSuffix(match, "*+")
				refreshRate, _ := strconv.ParseFloat(refreshRateStr, 32)
				currentDisplay.RefreshRate = float32(refreshRate)
				displays = append(displays, *currentDisplay)
				currentDisplay = nil
			}
		}
	}

	return displays
}

func extractRawPixelData(xwdOutput []byte, width, height int) ([]byte, error) {
	// The XWD file format includes a header before the pixel data.
	// The header size is typically 100 bytes, but this may vary depending on the X server.
	const xwdHeaderSize = 100

	if len(xwdOutput) <= xwdHeaderSize {
		return nil, fmt.Errorf("invalid xwd output: too small")
	}

	// Extract the raw pixel data (BGRA format)
	rawPixelData := xwdOutput[xwdHeaderSize:]

	// Verify the size of the raw pixel data
	expectedSize := width * height * 4 // 4 bytes per pixel (32-bit color)
	if len(rawPixelData) < expectedSize {
		return nil, fmt.Errorf("invalid xwd output: insufficient pixel data")
	}

	// Convert the pixel data to BGRA format (if necessary)
	// Note: Depending on the X server, the pixel format may already be BGRA.
	// If conversion is needed, implement it here.

	return rawPixelData, nil
}

func isDisplayDetails(xrandrOutput string) bool {
	return strings.Contains(xrandrOutput, " connected ")
}

func isPrimaryDisplay(xrandrOutput string) bool {
	return strings.Contains(xrandrOutput, " primary ")
}
