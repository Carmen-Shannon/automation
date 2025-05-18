# Automation
This module serves as an automation tool and supports various system actions related to interacting with the display, mouse and keyboard.

## Installation
Install the latest version (v1.0.0):
```bash
go get github.com/Carmen-Shannon/automation@v1.0.0
```

Depending on your operating system there are different requirements for this module:

### Windows
- user32.dll
- gdi32.dll

These are included in all windows builds by default, so no additional setup is required.

### Linux
- xdotool
- xrandr
- xwd
- x11
- ImageMagick

Run the following apt to get all dependencies for linux:

```bash
sudo apt install libx11-dev xdotool x11-xserver-utils x11-utils x11-apps ImageMagick
```

## Documentation
The following section will highlight the available interfaces and their usage, for full documentation please refer to the inline function docs.

### Packages
#### Device
- `Display`
    - `VirtualScreen`
        - Handles the virtual screen space
        - Includes references to all connected displays
        - Can capture displays or specified window boundaries in BMP format
    - `BMP`
        - The struct for a BMP image, contains a ToBinary function that converts the struct to a valid BMP image file in bytes
- `Keyboard`
    - `KeyPress`
        - The only function currently in the keyboard package, it allows simulation of a key press.
        - Includes support for linux and windows english utf-8 keys
        - Supports a combination of keys, such as modifiers like shift
- `Mouse`
    - `Mouse`
        - Look at that I named one package consistently...
        - The mouse interface handles all mouse actions, such as clicking and moving
        - Has options that allow for parabolic/smoothed movement and jitter

#### Tools
- `Matcher`
    - `Matcher`
        - This interface handles template-matching a sub-image to it's relative x and y positions of the scanned image
        - Takes advantage of concurrency to scan multiple parts of the image at a time
        - Has threshold values and timeout options that can be set to control the fuzzy matching
- `Worker`
    - `DynamicWorkerPool`
        - This interface allows for concurrent tasks to be scheduled and completed within a controlled environment
        - It has signal flags available to listen for when work has completed
        - It has granular control over each worker
        - It has a separate process to handle scheduling the work in the queue to the available workers


## Example Usage
See `main.go.example` for some example usage of the automation tooling.

