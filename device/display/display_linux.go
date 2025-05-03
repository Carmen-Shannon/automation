//go:build linux
// +build linux

package display

import (
	"regexp"
	"strconv"
	"strings"

	"automation/tools/linux"
)

// DetectDisplays detects all connected displays and ensures the primary display is at index 0.
func (vs *virtualScreen) DetectDisplays() ([]Display, error) {
	// Execute the `xrandr` command to get display information
	output, err := linux.ExecuteXrandr()
	if err != nil {
		return nil, err
	}

	// Parse the output of the xrandr command
	return extractDisplaysFromXrandrOutput(string(output)), nil
}

func Init() VirtualScreen {
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

func isDisplayDetails(xrandrOutput string) bool {
	return strings.Contains(xrandrOutput, " connected ")
}

func isPrimaryDisplay(xrandrOutput string) bool {
	return strings.Contains(xrandrOutput, " primary ")
}
