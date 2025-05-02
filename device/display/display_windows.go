//go:build windows
// +build windows

package display

import (
	"syscall"
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
	Size         uint32      `json:"size"`
	DeviceName   [32]uint16  `json:"device_name"`
	DeviceString [128]uint16 `json:"device_string"`
	StateFlags   uint32      `json:"state_flags"`
	DeviceID     [128]uint16 `json:"device_id"`
	DeviceKey    [128]uint16 `json:"device_key"`
}

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	enumDisplayDevices  = user32.NewProc("EnumDisplayDevicesW")
	enumDisplaySettings = user32.NewProc("EnumDisplaySettingsW")
)

func DetectDisplays() ([]Display, error) {
	var displays []Display
	var device displayDevice
	device.Size = uint32(unsafe.Sizeof(device))

	for i := 0; ; i++ {
		ret, _, _ := enumDisplayDevices.Call(0, uintptr(i), uintptr(unsafe.Pointer(&device)), uintptr(0x00000001))
		if ret == 0 {
			break
		}

		// Skip devices that are not attached to the desktop
		if device.StateFlags&0x00000001 == 0 { // DISPLAY_DEVICE_ATTACHED_TO_DESKTOP
			continue
		}

		var dm devMode
		dm.Size = uint16(unsafe.Sizeof(dm))
		ret, _, _ = enumDisplaySettings.Call(uintptr(unsafe.Pointer(&device.DeviceName)), uintptr(0xFFFFFFFF), uintptr(unsafe.Pointer(&dm)))
		if ret == 0 {
			continue
		}

		displays = append(displays, Display{
			Width:       int(dm.PelsWidth),
			Height:      int(dm.PelsHeight),
			RefreshRate: float32(dm.DisplayFrequency),
		})

	}
	return displays, nil
}
