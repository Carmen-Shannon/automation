//go:build windows
// +build windows

package display

import (
	"automation/tools/windows"
	"fmt"
	"unsafe"
)

// rect represents a rectangle with coordinates for the display.
type rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

// monitorInfo contains information about a display monitor.
type monitorInfo struct {
	Size    uint32
	Monitor rect
	Work    rect
	Flags   uint32
}

// devMode represents the device mode for a display
type devMode struct {
	DeviceName    [32]uint16 // dmDeviceName: Friendly name of the display device
	SpecVersion   uint16     // dmSpecVersion: Version number of the DEVMODE structure
	DriverVersion uint16     // dmDriverVersion: Driver version number
	Size          uint16     // dmSize: Size of the DEVMODE structure
	DriverExtra   uint16     // dmDriverExtra: Size of private driver data
	Fields        uint32     // dmFields: Flags indicating which fields are initialized

	// Union: dmPosition or printer-specific fields
	PositionX          int32  // dmPosition.x: X-coordinate of the display position
	PositionY          int32  // dmPosition.y: Y-coordinate of the display position
	DisplayOrientation uint32 // dmDisplayOrientation: Orientation of the display
	DisplayFixedOutput uint32 // dmDisplayFixedOutput: Fixed output settings for the display

	Color       int16      // dmColor: Color or monochrome setting
	Duplex      int16      // dmDuplex: Duplex printing setting
	YResolution int16      // dmYResolution: Y-resolution for printers
	TTOption    int16      // dmTTOption: TrueType font option
	Collate     int16      // dmCollate: Collation setting for printers
	FormName    [32]uint16 // dmFormName: Form name for printers
	LogPixels   uint16     // dmLogPixels: Logical pixels per inch

	BitsPerPel uint32 // dmBitsPerPel: Color resolution in bits per pixel
	PelsWidth  uint32 // dmPelsWidth: Width of the display in pixels
	PelsHeight uint32 // dmPelsHeight: Height of the display in pixels

	// Union: dmDisplayFlags or dmNup
	DisplayFlags     uint32 // dmDisplayFlags: Display mode flags
	DisplayFrequency uint32 // dmDisplayFrequency: Refresh rate in hertz

	ICMMethod     uint32 // dmICMMethod: ICM handling method
	ICMIntent     uint32 // dmICMIntent: ICM intent
	MediaType     uint32 // dmMediaType: Media type for printers
	DitherType    uint32 // dmDitherType: Dithering method
	Reserved1     uint32 // dmReserved1: Reserved; must be zero
	Reserved2     uint32 // dmReserved2: Reserved; must be zero
	PanningWidth  uint32 // dmPanningWidth: Panning width; must be zero
	PanningHeight uint32 // dmPanningHeight: Panning height; must be zero
}

// displayDevice represents the DISPLAY_DEVICE structure.
type displayDevice struct {
	Size         uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

type bitmapInfoHeader struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type bitmapInfo struct {
	BmiHeader bitmapInfoHeader
	BmiColors [1]uint32
}

func Init() VirtualScreen {
	// Retrieve the virtual screen's top-left corner
	left, _, _ := windows.GetSystemMetrics.Call(uintptr(windows.SM_XVIRTUALSCREEN))
	bottom, _, _ := windows.GetSystemMetrics.Call(uintptr(windows.SM_YVIRTUALSCREEN))

	// Retrieve the virtual screen's dimensions
	right, _, _ := windows.GetSystemMetrics.Call(uintptr(windows.SM_CXVIRTUALSCREEN))
	top, _, _ := windows.GetSystemMetrics.Call(uintptr(windows.SM_CYVIRTUALSCREEN))

	// Construct the VirtualScreen struct
	vs := virtualScreen{
		Left:   int32(left),
		Right:  int32(right),
		Top:    int32(top),
		Bottom: int32(bottom),
	}
	displays, err := vs.DetectDisplays()
	if err != nil {
		return &vs
	}
	vs.Displays = displays

	return &vs
}

func (vs *virtualScreen) CaptureBmp(options ...DisplayCaptureOption) ([][]byte, error) {
	displayCaptureOptions := &displayCaptureOption{}
	for _, opt := range options {
		opt(displayCaptureOptions)
	}

	// Parse the DisplayCaptureOption varargs
	var displays []Display
	if len(options) == 0 || len(displayCaptureOptions.Displays) == 0 {
		// Default to capturing the primary display
		pd, err := vs.GetPrimaryDisplay()
		if err != nil {
			return nil, err
		}
		displays = append(displays, pd)
	} else {
		// Use the specified displays
		displays = displayCaptureOptions.Displays
	}

	// Prepare a slice to hold the bitmap data for each display
	var bitmaps [][]byte

	// Iterate over the displays and capture each one
	for _, display := range displays {
		// Get the device context of the screen
		hdcScreen, _, err := windows.GetDC.Call(0)
		if hdcScreen == 0 {
			return nil, fmt.Errorf("failed to get screen device context: %w", err)
		}
		defer windows.ReleaseDC.Call(0, hdcScreen)

		// Create a compatible device context
		hdcMem, _, err := windows.CreateCompatibleDC.Call(hdcScreen)
		if hdcMem == 0 {
			return nil, fmt.Errorf("failed to create compatible device context: %w", err)
		}
		defer windows.DeleteDC.Call(hdcMem)

		// Calculate the width and height of the display
		width := display.Width
		height := display.Height

		// Create a compatible bitmap
		hBitmap, _, err := windows.CreateCompatibleBitmap.Call(hdcScreen, uintptr(width), uintptr(height))
		if hBitmap == 0 {
			return nil, fmt.Errorf("failed to create compatible bitmap: %w", err)
		}
		defer windows.DeleteObject.Call(hBitmap)

		// Select the bitmap into the memory device context
		oldBitmap, _, err := windows.SelectObject.Call(hdcMem, hBitmap)
		if oldBitmap == 0 {
			return nil, fmt.Errorf("failed to select bitmap into device context: %w", err)
		}
		defer windows.SelectObject.Call(hdcMem, oldBitmap)

		// Copy the screen contents into the memory device context
		ret, _, err := windows.BitBlt.Call(hdcMem, 0, 0, uintptr(width), uintptr(height), hdcScreen, uintptr(display.X), uintptr(display.Y), uintptr(windows.SRCCOPY))
		if ret == 0 {
			return nil, fmt.Errorf("failed to copy screen contents: %w", err)
		}

		// Retrieve the bitmap data
		var bmpInfo bitmapInfo
		bmpInfo.BmiHeader.BiSize = uint32(unsafe.Sizeof(bmpInfo.BmiHeader))
		bmpInfo.BmiHeader.BiWidth = int32(width)
		bmpInfo.BmiHeader.BiHeight = -int32(height) // Pnegative height for top-down DIB
		bmpInfo.BmiHeader.BiPlanes = 1
		bmpInfo.BmiHeader.BiBitCount = 32
		bmpInfo.BmiHeader.BiCompression = windows.BI_RGB

		// Calculate the size of the bitmap data
		rowSize := (int(width)*4 + 3) & ^3 // Row size must be a multiple of 4 bytes
		bitmapSize := rowSize * int(height)

		// Allocate memory for the bitmap data
		bitmapData := make([]byte, bitmapSize)

		// Get the bitmap data
		ret, _, err = windows.GetDIBits.Call(hdcMem, hBitmap, 0, uintptr(height), uintptr(unsafe.Pointer(&bitmapData[0])), uintptr(unsafe.Pointer(&bmpInfo)), uintptr(windows.DIB_RGB_COLORS))
		if ret == 0 {
			return nil, fmt.Errorf("failed to retrieve bitmap data: %w", err)
		}

		// Append the complete BMP data to the result slice
		bitmaps = append(bitmaps, bitmapData)
	}

	return bitmaps, nil
}

func (vs *virtualScreen) DetectDisplays() ([]Display, error) {
	var displays []Display
	var device displayDevice
	device.Size = uint32(unsafe.Sizeof(device))

	for i := 0; ; i++ {
		ret, _, _ := windows.EnumDisplayDevices.Call(0, uintptr(i), uintptr(unsafe.Pointer(&device)), uintptr(0x00000001))
		if ret == 0 {
			break
		}

		// Skip devices that are not attached to the desktop
		if device.StateFlags&0x00000001 == 0 { // DISPLAY_DEVICE_ATTACHED_TO_DESKTOP
			continue
		}

		var dm devMode
		dm.Size = uint16(unsafe.Sizeof(dm))
		ret, _, _ = windows.EnumDisplaySettings.Call(uintptr(unsafe.Pointer(&device.DeviceName)), uintptr(0xFFFFFFFF), uintptr(unsafe.Pointer(&dm)))
		if ret == 0 {
			continue
		}
		var primary bool
		if dm.PositionX == 0 && dm.PositionY == 0 {
			primary = true
		}

		displays = append(displays, Display{
			X:           dm.PositionX,
			Y:           dm.PositionY,
			Width:       int(dm.PelsWidth),
			Height:      int(dm.PelsHeight),
			RefreshRate: float32(dm.DisplayFrequency),
			Primary:     primary,
		})

	}
	vs.Displays = displays
	return displays, nil
}
