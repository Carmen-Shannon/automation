//go:build windows
// +build windows

package key_codes

const (
	// Letters
	KeyCodeA KeyCode = 0x41
	KeyCodeB KeyCode = 0x42
	KeyCodeC KeyCode = 0x43
	KeyCodeD KeyCode = 0x44
	KeyCodeE KeyCode = 0x45
	KeyCodeF KeyCode = 0x46
	KeyCodeG KeyCode = 0x47
	KeyCodeH KeyCode = 0x48
	KeyCodeI KeyCode = 0x49
	KeyCodeJ KeyCode = 0x4A
	KeyCodeK KeyCode = 0x4B
	KeyCodeL KeyCode = 0x4C
	KeyCodeM KeyCode = 0x4D
	KeyCodeN KeyCode = 0x4E
	KeyCodeO KeyCode = 0x4F
	KeyCodeP KeyCode = 0x50
	KeyCodeQ KeyCode = 0x51
	KeyCodeR KeyCode = 0x52
	KeyCodeS KeyCode = 0x53
	KeyCodeT KeyCode = 0x54
	KeyCodeU KeyCode = 0x55
	KeyCodeV KeyCode = 0x56
	KeyCodeW KeyCode = 0x57
	KeyCodeX KeyCode = 0x58
	KeyCodeY KeyCode = 0x59
	KeyCodeZ KeyCode = 0x5A

	// Numbers
	KeyCode0 KeyCode = 0x30
	KeyCode1 KeyCode = 0x31
	KeyCode2 KeyCode = 0x32
	KeyCode3 KeyCode = 0x33
	KeyCode4 KeyCode = 0x34
	KeyCode5 KeyCode = 0x35
	KeyCode6 KeyCode = 0x36
	KeyCode7 KeyCode = 0x37
	KeyCode8 KeyCode = 0x38
	KeyCode9 KeyCode = 0x39

	// Function Keys
	KeyCodeF1  KeyCode = 0x70
	KeyCodeF2  KeyCode = 0x71
	KeyCodeF3  KeyCode = 0x72
	KeyCodeF4  KeyCode = 0x73
	KeyCodeF5  KeyCode = 0x74
	KeyCodeF6  KeyCode = 0x75
	KeyCodeF7  KeyCode = 0x76
	KeyCodeF8  KeyCode = 0x77
	KeyCodeF9  KeyCode = 0x78
	KeyCodeF10 KeyCode = 0x79
	KeyCodeF11 KeyCode = 0x7A
	KeyCodeF12 KeyCode = 0x7B

	// Control Keys
	KeyCodeShift      KeyCode = 0x10
	KeyCodeCtrl       KeyCode = 0x11
	KeyCodeAlt        KeyCode = 0x12
	KeyCodeCaps       KeyCode = 0x14
	KeyCodeTab        KeyCode = 0x09
	KeyCodeEnter      KeyCode = 0x0D
	KeyCodeEscape     KeyCode = 0x1B
	KeyCodeSpace      KeyCode = 0x20
	KeyCodeBack       KeyCode = 0x08
	KeyCodeDelete     KeyCode = 0x2E
	KeyCodeInsert     KeyCode = 0x2D
	KeyCodeHome       KeyCode = 0x24
	KeyCodeEnd        KeyCode = 0x23
	KeyCodePageUp     KeyCode = 0x21
	KeyCodePageDown   KeyCode = 0x22
	KeyCodeLeftShift  KeyCode = 0xA0
	KeyCodeRightShift KeyCode = 0xA1
	KeyCodeLeftCtrl   KeyCode = 0xA2
	KeyCodeRightCtrl  KeyCode = 0xA3
	KeyCodeLeftAlt    KeyCode = 0xA4
	KeyCodeRightAlt   KeyCode = 0xA5

	// Arrow Keys
	KeyCodeLeft  KeyCode = 0x25
	KeyCodeUp    KeyCode = 0x26
	KeyCodeRight KeyCode = 0x27
	KeyCodeDown  KeyCode = 0x28

	// Numpad Keys
	KeyCodeNumpad0  KeyCode = 0x60
	KeyCodeNumpad1  KeyCode = 0x61
	KeyCodeNumpad2  KeyCode = 0x62
	KeyCodeNumpad3  KeyCode = 0x63
	KeyCodeNumpad4  KeyCode = 0x64
	KeyCodeNumpad5  KeyCode = 0x65
	KeyCodeNumpad6  KeyCode = 0x66
	KeyCodeNumpad7  KeyCode = 0x67
	KeyCodeNumpad8  KeyCode = 0x68
	KeyCodeNumpad9  KeyCode = 0x69
	KeyCodeMultiply KeyCode = 0x6A
	KeyCodeAdd      KeyCode = 0x6B
	KeyCodeSubtract KeyCode = 0x6D
	KeyCodeDecimal  KeyCode = 0x6E
	KeyCodeDivide   KeyCode = 0x6F

	// Special Keys
	KeyCodePrintScreen  KeyCode = 0x2C
	KeyCodeScrollLock   KeyCode = 0x91
	KeyCodePause        KeyCode = 0x13
	KeyCodeNumLock      KeyCode = 0x90
	KeyCodeSemicolon    KeyCode = 0xBA // ;:
	KeyCodeEqual        KeyCode = 0xBB // =
	KeyCodeComma        KeyCode = 0xBC // ,
	KeyCodeMinus        KeyCode = 0xBD // -
	KeyCodePeriod       KeyCode = 0xBE // .
	KeyCodeFwdSlash     KeyCode = 0xBF // /?
	KeyCodeTilde        KeyCode = 0xC0 // `~
	KeyCodeLeftBracket  KeyCode = 0xDB // [{
	KeyCodeBackslash    KeyCode = 0xDC // \|
	KeyCodeRightBracket KeyCode = 0xDD // ]}
	KeyCodeQuote        KeyCode = 0xDE // '"
)
