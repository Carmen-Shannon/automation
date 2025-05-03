//go:build windows
// +build windows

package windows

import "syscall"

var (
	User32              = syscall.NewLazyDLL("user32.dll")
	EnumDisplayDevices  = User32.NewProc("EnumDisplayDevicesW")
	EnumDisplaySettings = User32.NewProc("EnumDisplaySettingsW")
	GetSystemMetrics    = User32.NewProc("GetSystemMetrics")
	SetCursorPos        = User32.NewProc("SetCursorPos")
	GetCursorPos        = User32.NewProc("GetCursorPos")
	MouseEvent          = User32.NewProc("mouse_event")
)

const (
	SM_XVIRTUALSCREEN  = 76 // The x-coordinate of the top-left corner of the virtual screen
	SM_YVIRTUALSCREEN  = 77 // The y-coordinate of the top-left corner of the virtual screen
	SM_CXVIRTUALSCREEN = 78 // The width of the virtual screen
	SM_CYVIRTUALSCREEN = 79 // The height of the virtual screen

	MOUSEEVENTF_LEFTDOWN   = 0x0002 // The left button is down flag
	MOUSEEVENTF_LEFTUP     = 0x0004 // The left button is up flag
	MOUSEEVENTF_RIGHTDOWN  = 0x0008 // The right button is down flag
	MOUSEEVENTF_RIGHTUP    = 0x0010 // The right button is up flag
	MOUSEEVENTF_MIDDLEDOWN = 0x0020 // The middle button is down flag
	MOUSEEVENTF_MIDDLEUP   = 0x0040 // The middle button is up flag
)
