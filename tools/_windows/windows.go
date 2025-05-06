//go:build windows
// +build windows

package windows

import (
	"fmt"
	"syscall"
)

var (
	// User32 DLL calls
	User32              = syscall.NewLazyDLL("user32.dll")
	EnumDisplayDevices  = User32.NewProc("EnumDisplayDevicesW")
	EnumDisplaySettings = User32.NewProc("EnumDisplaySettingsW")
	GetSystemMetrics    = User32.NewProc("GetSystemMetrics")
	SetCursorPos        = User32.NewProc("SetCursorPos")
	GetCursorPos        = User32.NewProc("GetCursorPos")
	MouseEvent          = User32.NewProc("mouse_event")
	KeybdEvent          = User32.NewProc("keybd_event")
	getDC               = User32.NewProc("GetDC")
	ReleaseDC           = User32.NewProc("ReleaseDC")
	MonitorFromRect     = User32.NewProc("MonitorFromRect")
	MonitorFromWindow   = User32.NewProc("MonitorFromWindow")
	EnumWindows         = User32.NewProc("EnumWindows")

	// GDI32 DLL calls
	Gdi32                  = syscall.NewLazyDLL("gdi32.dll")
	createCompatibleDC     = Gdi32.NewProc("CreateCompatibleDC")
	DeleteDC               = Gdi32.NewProc("DeleteDC")
	createCompatibleBitmap = Gdi32.NewProc("CreateCompatibleBitmap")
	selectObject           = Gdi32.NewProc("SelectObject")
	DeleteObject           = Gdi32.NewProc("DeleteObject")
	bitBlt                 = Gdi32.NewProc("BitBlt")
	GetDIBits              = Gdi32.NewProc("GetDIBits")
	GetDeviceCaps          = Gdi32.NewProc("GetDeviceCaps")
)

const (
	// System metrics constants
	SM_XVIRTUALSCREEN  = 76 // The x-coordinate of the top-left corner of the virtual screen
	SM_YVIRTUALSCREEN  = 77 // The y-coordinate of the top-left corner of the virtual screen
	SM_CXVIRTUALSCREEN = 78 // The width of the virtual screen
	SM_CYVIRTUALSCREEN = 79 // The height of the virtual screen

	// Mouse event flags
	MOUSEEVENTF_LEFTDOWN   = 0x0002 // The left button is down flag
	MOUSEEVENTF_LEFTUP     = 0x0004 // The left button is up flag
	MOUSEEVENTF_RIGHTDOWN  = 0x0008 // The right button is down flag
	MOUSEEVENTF_RIGHTUP    = 0x0010 // The right button is up flag
	MOUSEEVENTF_MIDDLEDOWN = 0x0020 // The middle button is down flag
	MOUSEEVENTF_MIDDLEUP   = 0x0040 // The middle button is up flag

	// these are for the SendInput function as flags, they are unused because SendInput sucks and doesn't work????
	INPUT_KEYBOARD        = 1      // Keyboard input type
	KEYEVENTF_EXTENDEDKEY = 0x0001 // Extended key flag for keyboard input
	KEYEVENTF_KEYUP       = 0x0002 // Key up flag for keyboard input
	KEYEVENTF_UNICODE     = 0x0004 // Unicode flag for keyboard input
	KEYEVENTF_SCANCODE    = 0x0008 // Scan code flag for keyboard input

	// GDI constants
	SRCCOPY                  = 0x00CC0020
	BI_RGB                   = 0
	DIB_RGB_COLORS           = 0
	LOGPIXELSX               = 88         // Logical pixels/inch in the X direction
	LOGPIXELSY               = 90         // Logical pixels/inch in the Y direction
	MONITOR_DEFAULTTONEAREST = 0x00000002 // Default monitor option for MonitorFromRect function
)

type BitmapInfoHeader struct {
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

type BitmapInfo struct {
	BmiHeader BitmapInfoHeader
	BmiColors [1]uint32
}

type BitmapHeader struct {
	Type      uint16
	Size      uint32
	Reserved1 uint16
	Reserved2 uint16
	OffBits   uint32
}

func GetScreenDC() (uintptr, error) {
	hdc, _, err := getDC.Call(0)
	if hdc == 0 {
		return 0, fmt.Errorf("failed to get screen device context: %w", err)
	}
	return hdc, nil
}

func CreateMemoryDC(hdc uintptr) (uintptr, error) {
	hdcMem, _, err := createCompatibleDC.Call(hdc)
	if hdcMem == 0 {
		return 0, fmt.Errorf("failed to create compatible device context: %w", err)
	}
	return hdcMem, nil
}

func CreateBitmap(hdc uintptr, width, height int) (uintptr, error) {
	hBitmap, _, err := createCompatibleBitmap.Call(hdc, uintptr(width), uintptr(height))
	if hBitmap == 0 {
		return 0, fmt.Errorf("failed to create compatible bitmap: %w", err)
	}
	return hBitmap, nil
}

func SelectBitmap(hdc uintptr, hBitmap uintptr) (uintptr, error) {
	oldBitmap, _, err := selectObject.Call(hdc, hBitmap)
	if oldBitmap == 0 {
		return 0, fmt.Errorf("failed to select bitmap into device context: %w", err)
	}
	return oldBitmap, nil
}

func CopyScreenToMemory(hdcDest, hdcSrc uintptr, xDest, yDest, width, height, xSrc, ySrc int) error {
	ret, _, err := bitBlt.Call(
		hdcDest, uintptr(xDest), uintptr(yDest), uintptr(width), uintptr(height),
		hdcSrc, uintptr(xSrc), uintptr(ySrc),
		uintptr(SRCCOPY),
	)
	if ret == 0 {
		return fmt.Errorf("failed to copy screen contents: %w", err)
	}
	return nil
}
