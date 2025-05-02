//go:build linux
// +build linux

package display

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// DetectDisplays detects all connected displays and ensures the primary display is at index 0.
func DetectDisplays() ([]Display, error) {
	// Execute the `xrandr` command to get display information
	output, err := exec.Command("xrandr", "--query").Output()
	if err != nil {
		return nil, err
	}

	// Parse the output of the xrandr command
	displays := parseXrandrOutput(string(output))

	// Ensure the primary display is at index 0
	for i, display := range displays {
		if display.Width > 0 && display.Height > 0 && display.RefreshRate > 0 {
			// Move primary display to the first position
			displays[0], displays[i] = displays[i], displays[0]
			break
		}
	}

	return displays, nil
}

// parseXrandrOutput parses the output of the `xrandr --query` command and returns a slice of Display structs.
func parseXrandrOutput(output string) []Display {
    lines := strings.Split(output, "\n")
    var displays []Display

    for _, line := range lines {
		if strings.Contains(line," connected") {
			var displayEntry Display
			// checking for the connected displays example: eDP-1 connected primary 1920x1080+0+0
			// Regular expression to match the resolution format
			re := regexp.MustCompile(`\d+x\d+(?=\+)`)
			match := re.FindString(line)
			if match != "" {
				res := strings.Split(match, "x")
				width, _ := strconv.Atoi(res[0])
				heightSplit := strings.Split(res[1], "+")[0]
				height, _ := strconv.Atoi(heightSplit)
				displayEntry.Width = width
				displayEntry.Height = height
				displays = append(displays, displayEntry)
			}
		}
	}

    return displays
}
