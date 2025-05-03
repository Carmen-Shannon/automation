package windows

import "syscall"

var (
	User32              = syscall.NewLazyDLL("user32.dll")
	EnumDisplayDevices  = User32.NewProc("EnumDisplayDevicesW")
	EnumDisplaySettings = User32.NewProc("EnumDisplaySettingsW")
	GetSystemMetrics    = User32.NewProc("GetSystemMetrics")
	SetCursorPos        = User32.NewProc("SetCursorPos")
	GetCursorPos        = User32.NewProc("GetCursorPos")
)

const (
	SM_XVIRTUALSCREEN  = 76 // The x-coordinate of the top-left corner of the virtual screen
	SM_YVIRTUALSCREEN  = 77 // The y-coordinate of the top-left corner of the virtual screen
	SM_CXVIRTUALSCREEN = 78 // The width of the virtual screen
	SM_CYVIRTUALSCREEN = 79 // The height of the virtual screen
)
