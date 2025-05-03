//go:build linux
// +build linux

package key_codes

const (
	// Letters
	KeyCodeA KeyCode = 30
	KeyCodeB KeyCode = 48
	KeyCodeC KeyCode = 46
	KeyCodeD KeyCode = 32
	KeyCodeE KeyCode = 18
	KeyCodeF KeyCode = 33
	KeyCodeG KeyCode = 34
	KeyCodeH KeyCode = 35
	KeyCodeI KeyCode = 23
	KeyCodeJ KeyCode = 36
	KeyCodeK KeyCode = 37
	KeyCodeL KeyCode = 38
	KeyCodeM KeyCode = 50
	KeyCodeN KeyCode = 49
	KeyCodeO KeyCode = 24
	KeyCodeP KeyCode = 25
	KeyCodeQ KeyCode = 16
	KeyCodeR KeyCode = 19
	KeyCodeS KeyCode = 31
	KeyCodeT KeyCode = 20
	KeyCodeU KeyCode = 22
	KeyCodeV KeyCode = 47
	KeyCodeW KeyCode = 17
	KeyCodeX KeyCode = 45
	KeyCodeY KeyCode = 21
	KeyCodeZ KeyCode = 44

	// Numbers
	KeyCode0 KeyCode = 11
	KeyCode1 KeyCode = 2
	KeyCode2 KeyCode = 3
	KeyCode3 KeyCode = 4
	KeyCode4 KeyCode = 5
	KeyCode5 KeyCode = 6
	KeyCode6 KeyCode = 7
	KeyCode7 KeyCode = 8
	KeyCode8 KeyCode = 9
	KeyCode9 KeyCode = 10

	// Function Keys
	KeyCodeF1  KeyCode = 59
	KeyCodeF2  KeyCode = 60
	KeyCodeF3  KeyCode = 61
	KeyCodeF4  KeyCode = 62
	KeyCodeF5  KeyCode = 63
	KeyCodeF6  KeyCode = 64
	KeyCodeF7  KeyCode = 65
	KeyCodeF8  KeyCode = 66
	KeyCodeF9  KeyCode = 67
	KeyCodeF10 KeyCode = 68
	KeyCodeF11 KeyCode = 87
	KeyCodeF12 KeyCode = 88

	// Control Keys
	KeyCodeShift      KeyCode = 42
	KeyCodeCtrl       KeyCode = 29
	KeyCodeAlt        KeyCode = 56
	KeyCodeCaps       KeyCode = 58
	KeyCodeTab        KeyCode = 15
	KeyCodeEnter      KeyCode = 28
	KeyCodeEscape     KeyCode = 1
	KeyCodeSpace      KeyCode = 57
	KeyCodeBack       KeyCode = 14
	KeyCodeDelete     KeyCode = 111
	KeyCodeInsert     KeyCode = 110
	KeyCodeHome       KeyCode = 102
	KeyCodeEnd        KeyCode = 107
	KeyCodePageUp     KeyCode = 104
	KeyCodePageDown   KeyCode = 109
	KeyCodeLeftShift  KeyCode = 42
	KeyCodeRightShift KeyCode = 54
	KeyCodeLeftCtrl   KeyCode = 29
	KeyCodeRightCtrl  KeyCode = 97
	KeyCodeLeftAlt    KeyCode = 56
	KeyCodeRightAlt   KeyCode = 100

	// Arrow Keys
	KeyCodeLeft  KeyCode = 105
	KeyCodeUp    KeyCode = 103
	KeyCodeRight KeyCode = 106
	KeyCodeDown  KeyCode = 108

	// Numpad Keys
	KeyCodeNumpad0  KeyCode = 82
	KeyCodeNumpad1  KeyCode = 79
	KeyCodeNumpad2  KeyCode = 80
	KeyCodeNumpad3  KeyCode = 81
	KeyCodeNumpad4  KeyCode = 75
	KeyCodeNumpad5  KeyCode = 76
	KeyCodeNumpad6  KeyCode = 77
	KeyCodeNumpad7  KeyCode = 71
	KeyCodeNumpad8  KeyCode = 72
	KeyCodeNumpad9  KeyCode = 73
	KeyCodeMultiply KeyCode = 55
	KeyCodeAdd      KeyCode = 78
	KeyCodeSubtract KeyCode = 74
	KeyCodeDecimal  KeyCode = 83
	KeyCodeDivide   KeyCode = 98

	// Special Keys
	KeyCodePrintScreen  KeyCode = 99
	KeyCodeScrollLock   KeyCode = 70
	KeyCodePause        KeyCode = 119
	KeyCodeNumLock      KeyCode = 69
	KeyCodeSemicolon    KeyCode = 39 // ;:
	KeyCodeEqual        KeyCode = 13 // =
	KeyCodeComma        KeyCode = 51 // ,
	KeyCodeMinus        KeyCode = 12 // -
	KeyCodePeriod       KeyCode = 52 // .
	KeyCodeFwdSlash     KeyCode = 53 // /?
	KeyCodeTilde        KeyCode = 41 // `~
	KeyCodeLeftBracket  KeyCode = 26 // [{
	KeyCodeBackslash    KeyCode = 43 // \|
	KeyCodeRightBracket KeyCode = 27 // ]}
	KeyCodeQuote        KeyCode = 40 // '"
)
