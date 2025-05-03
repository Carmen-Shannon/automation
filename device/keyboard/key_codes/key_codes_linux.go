//go:build linux
// +build linux

package key_codes

type KeyCode uint32

const (
	// Letters
	KeyCodeA KeyCode = 0x0061 // XK_a
	KeyCodeB KeyCode = 0x0062 // XK_b
	KeyCodeC KeyCode = 0x0063 // XK_c
	KeyCodeD KeyCode = 0x0064 // XK_d
	KeyCodeE KeyCode = 0x0065 // XK_e
	KeyCodeF KeyCode = 0x0066 // XK_f
	KeyCodeG KeyCode = 0x0067 // XK_g
	KeyCodeH KeyCode = 0x0068 // XK_h
	KeyCodeI KeyCode = 0x0069 // XK_i
	KeyCodeJ KeyCode = 0x006a // XK_j
	KeyCodeK KeyCode = 0x006b // XK_k
	KeyCodeL KeyCode = 0x006c // XK_l
	KeyCodeM KeyCode = 0x006d // XK_m
	KeyCodeN KeyCode = 0x006e // XK_n
	KeyCodeO KeyCode = 0x006f // XK_o
	KeyCodeP KeyCode = 0x0070 // XK_p
	KeyCodeQ KeyCode = 0x0071 // XK_q
	KeyCodeR KeyCode = 0x0072 // XK_r
	KeyCodeS KeyCode = 0x0073 // XK_s
	KeyCodeT KeyCode = 0x0074 // XK_t
	KeyCodeU KeyCode = 0x0075 // XK_u
	KeyCodeV KeyCode = 0x0076 // XK_v
	KeyCodeW KeyCode = 0x0077 // XK_w
	KeyCodeX KeyCode = 0x0078 // XK_x
	KeyCodeY KeyCode = 0x0079 // XK_y
	KeyCodeZ KeyCode = 0x007a // XK_z

	// Numbers
	KeyCode0 KeyCode = 0x0030 // XK_0
	KeyCode1 KeyCode = 0x0031 // XK_1
	KeyCode2 KeyCode = 0x0032 // XK_2
	KeyCode3 KeyCode = 0x0033 // XK_3
	KeyCode4 KeyCode = 0x0034 // XK_4
	KeyCode5 KeyCode = 0x0035 // XK_5
	KeyCode6 KeyCode = 0x0036 // XK_6
	KeyCode7 KeyCode = 0x0037 // XK_7
	KeyCode8 KeyCode = 0x0038 // XK_8
	KeyCode9 KeyCode = 0x0039 // XK_9

	// Function Keys
	KeyCodeF1  KeyCode = 0xffbe // XK_F1
	KeyCodeF2  KeyCode = 0xffbf // XK_F2
	KeyCodeF3  KeyCode = 0xffc0 // XK_F3
	KeyCodeF4  KeyCode = 0xffc1 // XK_F4
	KeyCodeF5  KeyCode = 0xffc2 // XK_F5
	KeyCodeF6  KeyCode = 0xffc3 // XK_F6
	KeyCodeF7  KeyCode = 0xffc4 // XK_F7
	KeyCodeF8  KeyCode = 0xffc5 // XK_F8
	KeyCodeF9  KeyCode = 0xffc6 // XK_F9
	KeyCodeF10 KeyCode = 0xffc7 // XK_F10
	KeyCodeF11 KeyCode = 0xffc8 // XK_F11
	KeyCodeF12 KeyCode = 0xffc9 // XK_F12

	// Control Keys
	KeyCodeShift      KeyCode = 0xffe1 // XK_Shift_L
	KeyCodeCtrl       KeyCode = 0xffe3 // XK_Control_L
	KeyCodeAlt        KeyCode = 0xffe9 // XK_Alt_L
	KeyCodeCaps       KeyCode = 0xffe5 // XK_Caps_Lock
	KeyCodeTab        KeyCode = 0xff09 // XK_Tab
	KeyCodeEnter      KeyCode = 0xff0d // XK_Return
	KeyCodeEscape     KeyCode = 0xff1b // XK_Escape
	KeyCodeSpace      KeyCode = 0x0020 // XK_space
	KeyCodeBack       KeyCode = 0xff08 // XK_BackSpace
	KeyCodeDelete     KeyCode = 0xffff // XK_Delete
	KeyCodeInsert     KeyCode = 0xff63 // XK_Insert
	KeyCodeHome       KeyCode = 0xff50 // XK_Home
	KeyCodeEnd        KeyCode = 0xff57 // XK_End
	KeyCodePageUp     KeyCode = 0xff55 // XK_Page_Up
	KeyCodePageDown   KeyCode = 0xff56 // XK_Page_Down
	KeyCodeLeftShift  KeyCode = 0xffe1 // XK_Shift_L
	KeyCodeRightShift KeyCode = 0xffe2 // XK_Shift_R
	KeyCodeLeftCtrl   KeyCode = 0xffe3 // XK_Control_L
	KeyCodeRightCtrl  KeyCode = 0xffe4 // XK_Control_R
	KeyCodeLeftAlt    KeyCode = 0xffe9 // XK_Alt_L
	KeyCodeRightAlt   KeyCode = 0xffea // XK_Alt_R

	// Arrow Keys
	KeyCodeLeft  KeyCode = 0xff51 // XK_Left
	KeyCodeUp    KeyCode = 0xff52 // XK_Up
	KeyCodeRight KeyCode = 0xff53 // XK_Right
	KeyCodeDown  KeyCode = 0xff54 // XK_Down

	// Numpad Keys
	KeyCodeNumpad0  KeyCode = 0xffb0 // XK_KP_0
	KeyCodeNumpad1  KeyCode = 0xffb1 // XK_KP_1
	KeyCodeNumpad2  KeyCode = 0xffb2 // XK_KP_2
	KeyCodeNumpad3  KeyCode = 0xffb3 // XK_KP_3
	KeyCodeNumpad4  KeyCode = 0xffb4 // XK_KP_4
	KeyCodeNumpad5  KeyCode = 0xffb5 // XK_KP_5
	KeyCodeNumpad6  KeyCode = 0xffb6 // XK_KP_6
	KeyCodeNumpad7  KeyCode = 0xffb7 // XK_KP_7
	KeyCodeNumpad8  KeyCode = 0xffb8 // XK_KP_8
	KeyCodeNumpad9  KeyCode = 0xffb9 // XK_KP_9
	KeyCodeMultiply KeyCode = 0xffaa // XK_KP_Multiply
	KeyCodeAdd      KeyCode = 0xffab // XK_KP_Add
	KeyCodeSubtract KeyCode = 0xffad // XK_KP_Subtract
	KeyCodeDecimal  KeyCode = 0xffae // XK_KP_Decimal
	KeyCodeDivide   KeyCode = 0xffaf // XK_KP_Divide

	// Special Keys
	KeyCodePrintScreen  KeyCode = 0xff61 // XK_Print
	KeyCodeScrollLock   KeyCode = 0xff14 // XK_Scroll_Lock
	KeyCodePause        KeyCode = 0xff13 // XK_Pause
	KeyCodeNumLock      KeyCode = 0xff7f // XK_Num_Lock
	KeyCodeSemicolon    KeyCode = 0x003b // XK_semicolon
	KeyCodeEqual        KeyCode = 0x003d // XK_equal
	KeyCodeComma        KeyCode = 0x002c // XK_comma
	KeyCodeMinus        KeyCode = 0x002d // XK_minus
	KeyCodePeriod       KeyCode = 0x002e // XK_period
	KeyCodeFwdSlash     KeyCode = 0x002f // XK_slash
	KeyCodeTilde        KeyCode = 0x0060 // XK_grave
	KeyCodeLeftBracket  KeyCode = 0x005b // XK_bracketleft
	KeyCodeBackslash    KeyCode = 0x005c // XK_backslash
	KeyCodeRightBracket KeyCode = 0x005d // XK_bracketright
	KeyCodeQuote        KeyCode = 0x0027 // XK_apostrophe
)
