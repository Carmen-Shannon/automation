//go:build windows
// +build windows

package display

import (
    "syscall"
    "unsafe"
)

var (
    user32              = syscall.NewLazyDLL("user32.dll")
    getSystemMetrics    = user32.NewProc("GetSystemMetrics")
    enumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
    getMonitorInfo      = user32.NewProc("GetMonitorInfoW")
    enumDisplaySettings = user32.NewProc("EnumDisplaySettingsW")
)

type rect struct {
    Left   int32
    Top    int32
    Right  int32
    Bottom int32
}

type monitorInfo struct {
    Size    uint32
    Monitor rect
    Work    rect
    Flags   uint32
}

type devMode struct {
    DeviceName       [32]uint16
    SpecVersion      uint16
    DriverVersion    uint16
    Size             uint16
    DriverExtra      uint16
    Fields           uint32
    PositionX        int32
    PositionY        int32
    DisplayOrientation uint32
    DisplayFixedOutput uint32
    Width            uint32
    Height           uint32
    PelsWidth        uint32
    PelsHeight       uint32
    BitsPerPel       uint32
    DisplayFrequency uint32 // Refresh rate
}

const (
    ENUM_CURRENT_SETTINGS = 0xFFFFFFFF
)

// DetectDisplays detects all connected displays and ensures the primary display is at index 0.
func DetectDisplays() ([]Display, error) {
    var displays []Display

    // Callback function for EnumDisplayMonitors
    callback := syscall.NewCallback(func(hMonitor, hdcMonitor, lprcMonitor, dwData uintptr) uintptr {
        // Get monitor info
        var mi monitorInfo
        mi.Size = uint32(unsafe.Sizeof(mi))
        if _, _, err := getMonitorInfo.Call(hMonitor, uintptr(unsafe.Pointer(&mi))); err != nil && err.Error() != "The operation completed successfully." {
            return 1 // Continue enumeration
        }

        // Get refresh rate using EnumDisplaySettings
        var dm devMode
        dm.Size = uint16(unsafe.Sizeof(dm))
        refreshRate := float32(0)
        if _, _, err := enumDisplaySettings.Call(uintptr(unsafe.Pointer(&dm.DeviceName)), ENUM_CURRENT_SETTINGS, uintptr(unsafe.Pointer(&dm))); err == nil {
            refreshRate = float32(dm.DisplayFrequency)
        }

        // Add monitor to the list
        displays = append(displays, Display{
            Width:       int(mi.Monitor.Right - mi.Monitor.Left),
            Height:      int(mi.Monitor.Bottom - mi.Monitor.Top),
            RefreshRate: refreshRate,
        })
        return 1 // Continue enumeration
    })

    // Enumerate all monitors
    if _, _, err := enumDisplayMonitors.Call(0, 0, callback, 0); err != nil && err.Error() != "The operation completed successfully." {
        return nil, err
    }

    // Ensure the primary display is at index 0
    primaryWidth, _, _ := getSystemMetrics.Call(0) // SM_CXSCREEN
    primaryHeight, _, _ := getSystemMetrics.Call(1) // SM_CYSCREEN
    for i, display := range displays {
        if display.Width == int(primaryWidth) && display.Height == int(primaryHeight) {
            // Move primary display to the first position
            displays[0], displays[i] = displays[i], displays[0]
            break
        }
    }

    return displays, nil
}