package key

import "github.com/go-gl/glfw/v3.3/glfw"

// Key can be any key, keyboard or mouse
type Key int

// all key constants
const (
	Unknown      = Key(glfw.KeyUnknown)        //
	Mouse1       = Key(glfw.MouseButton1)      //
	Mouse2       = Key(glfw.MouseButton2)      //
	Mouse3       = Key(glfw.MouseButton3)      //
	Mouse4       = Key(glfw.MouseButton4)      //
	Mouse5       = Key(glfw.MouseButton5)      //
	Mouse6       = Key(glfw.MouseButton6)      //
	Mouse7       = Key(glfw.MouseButton7)      //
	Mouse8       = Key(glfw.MouseButton8)      //
	MouseLast    = Key(glfw.MouseButtonLast)   //
	MouseLeft    = Key(glfw.MouseButtonLeft)   //
	MouseRight   = Key(glfw.MouseButtonRight)  //
	MouseMiddle  = Key(glfw.MouseButtonMiddle) //
	Space        = Key(glfw.KeySpace)          //
	Apostrophe   = Key(glfw.KeyApostrophe)     // '
	Comma        = Key(glfw.KeyComma)          // ,
	Minus        = Key(glfw.KeyMinus)          // -
	Period       = Key(glfw.KeyPeriod)         // .
	Slash        = Key(glfw.KeySlash)          // /
	_0           = Key(glfw.Key0)              //
	_1           = Key(glfw.Key1)              //
	_2           = Key(glfw.Key2)              //
	_3           = Key(glfw.Key3)              //
	_4           = Key(glfw.Key4)              //
	_5           = Key(glfw.Key5)              //
	_6           = Key(glfw.Key6)              //
	_7           = Key(glfw.Key7)              //
	_8           = Key(glfw.Key8)              //
	_9           = Key(glfw.Key9)              //
	Semicolon    = Key(glfw.KeySemicolon)      // ;
	Equal        = Key(glfw.KeyEqual)          // =
	A            = Key(glfw.KeyA)              //
	B            = Key(glfw.KeyB)              //
	C            = Key(glfw.KeyC)              //
	D            = Key(glfw.KeyD)              //
	E            = Key(glfw.KeyE)              //
	F            = Key(glfw.KeyF)              //
	G            = Key(glfw.KeyG)              //
	H            = Key(glfw.KeyH)              //
	I            = Key(glfw.KeyI)              //
	J            = Key(glfw.KeyJ)              //
	K            = Key(glfw.KeyK)              //
	L            = Key(glfw.KeyL)              //
	M            = Key(glfw.KeyM)              //
	N            = Key(glfw.KeyN)              //
	O            = Key(glfw.KeyO)              //
	P            = Key(glfw.KeyP)              //
	Q            = Key(glfw.KeyQ)              //
	R            = Key(glfw.KeyR)              //
	S            = Key(glfw.KeyS)              //
	T            = Key(glfw.KeyT)              //
	U            = Key(glfw.KeyU)              //
	V            = Key(glfw.KeyV)              //
	W            = Key(glfw.KeyW)              //
	X            = Key(glfw.KeyX)              //
	Y            = Key(glfw.KeyY)              //
	Z            = Key(glfw.KeyZ)              //
	LeftBracket  = Key(glfw.KeyLeftBracket)    // [
	Backslash    = Key(glfw.KeyBackslash)      // \
	RightBracket = Key(glfw.KeyRightBracket)   // ]
	GraveAccent  = Key(glfw.KeyGraveAccent)    // `
	World1       = Key(glfw.KeyWorld1)         // non-US #1
	World2       = Key(glfw.KeyWorld2)         // non-US #2
	Escape       = Key(glfw.KeyEscape)         //
	Enter        = Key(glfw.KeyEnter)          //
	Tab          = Key(glfw.KeyTab)            //
	Backspace    = Key(glfw.KeyBackspace)      //
	Insert       = Key(glfw.KeyInsert)         //
	Delete       = Key(glfw.KeyDelete)         //
	Right        = Key(glfw.KeyRight)          //
	Left         = Key(glfw.KeyLeft)           //
	Down         = Key(glfw.KeyDown)           //
	Up           = Key(glfw.KeyUp)             //
	PageUp       = Key(glfw.KeyPageUp)         //
	PageDown     = Key(glfw.KeyPageDown)       //
	Home         = Key(glfw.KeyHome)           //
	End          = Key(glfw.KeyEnd)            //
	CapsLock     = Key(glfw.KeyCapsLock)       //
	ScrollLock   = Key(glfw.KeyScrollLock)     //
	NumLock      = Key(glfw.KeyNumLock)        //
	PrintScreen  = Key(glfw.KeyPrintScreen)    //
	Pause        = Key(glfw.KeyPause)          //
	F1           = Key(glfw.KeyF1)             //
	F2           = Key(glfw.KeyF2)             //
	F3           = Key(glfw.KeyF3)             //
	F4           = Key(glfw.KeyF4)             //
	F5           = Key(glfw.KeyF5)             //
	F6           = Key(glfw.KeyF6)             //
	F7           = Key(glfw.KeyF7)             //
	F8           = Key(glfw.KeyF8)             //
	F9           = Key(glfw.KeyF9)             //
	F10          = Key(glfw.KeyF10)            //
	F11          = Key(glfw.KeyF11)            //
	F12          = Key(glfw.KeyF12)            //
	F13          = Key(glfw.KeyF13)            //
	F14          = Key(glfw.KeyF14)            //
	F15          = Key(glfw.KeyF15)            //
	F16          = Key(glfw.KeyF16)            //
	F17          = Key(glfw.KeyF17)            //
	F18          = Key(glfw.KeyF18)            //
	F19          = Key(glfw.KeyF19)            //
	F20          = Key(glfw.KeyF20)            //
	F21          = Key(glfw.KeyF21)            //
	F22          = Key(glfw.KeyF22)            //
	F23          = Key(glfw.KeyF23)            //
	F24          = Key(glfw.KeyF24)            //
	F25          = Key(glfw.KeyF25)            //
	LeftShift    = Key(glfw.KeyLeftShift)      //
	LeftControl  = Key(glfw.KeyLeftControl)    //
	LeftAlt      = Key(glfw.KeyLeftAlt)        //
	LeftSuper    = Key(glfw.KeyLeftSuper)      //
	RightShift   = Key(glfw.KeyRightShift)     //
	RightControl = Key(glfw.KeyRightControl)   //
	RightAlt     = Key(glfw.KeyRightAlt)       //
	RightSuper   = Key(glfw.KeyRightSuper)     //
	Menu         = Key(glfw.KeyMenu)           //
	Last         = Key(glfw.KeyMenu)           //
)

func (k Key) String() string {
	val, ok := Names[k]
	if !ok {
		return "Invalid"
	}

	return val
}

// Names is helper for Key.String() method
var Names = map[Key]string{
	Unknown:      "Unknown",
	Mouse4:       "Mouse4",
	Mouse5:       "Mouse5",
	Mouse6:       "Mouse6",
	Mouse7:       "Mouse7",
	MouseLast:    "MouseLast",
	MouseLeft:    "MouseLeft",
	MouseRight:   "MouseRight",
	MouseMiddle:  "MouseMiddle",
	Space:        "Space",
	Apostrophe:   "Apostrophe",
	Comma:        "Comma",
	Minus:        "Minus",
	Period:       "Period",
	Slash:        "Slash",
	_0:           "0",
	_1:           "1",
	_2:           "2",
	_3:           "3",
	_4:           "4",
	_5:           "5",
	_6:           "6",
	_7:           "7",
	_8:           "8",
	_9:           "9",
	Semicolon:    "Semicolon",
	Equal:        "Equal",
	A:            "A",
	B:            "B",
	C:            "C",
	D:            "D",
	E:            "E",
	F:            "F",
	G:            "G",
	H:            "H",
	I:            "I",
	J:            "J",
	K:            "K",
	L:            "L",
	M:            "M",
	N:            "N",
	O:            "O",
	P:            "P",
	Q:            "Q",
	R:            "R",
	S:            "S",
	T:            "T",
	U:            "U",
	V:            "V",
	W:            "W",
	X:            "X",
	Y:            "Y",
	Z:            "Z",
	LeftBracket:  "LeftBracket",
	Backslash:    "Backslash",
	RightBracket: "RightBracket",
	GraveAccent:  "GraveAccent",
	World1:       "World1",
	World2:       "World2",
	Escape:       "Escape",
	Enter:        "Enter",
	Tab:          "Tab",
	Backspace:    "Backspace",
	Insert:       "Insert",
	Delete:       "Delete",
	Right:        "Right",
	Left:         "Left",
	Down:         "Down",
	Up:           "Up",
	PageUp:       "PageUp",
	PageDown:     "PageDown",
	Home:         "Home",
	End:          "End",
	CapsLock:     "CapsLock",
	ScrollLock:   "ScrollLock",
	NumLock:      "NumLock",
	PrintScreen:  "PrintScreen",
	Pause:        "Pause",
	F1:           "F1",
	F2:           "F2",
	F3:           "F3",
	F4:           "F4",
	F5:           "F5",
	F6:           "F6",
	F7:           "F7",
	F8:           "F8",
	F9:           "F9",
	F10:          "F10",
	F11:          "F11",
	F12:          "F12",
	F13:          "F13",
	F14:          "F14",
	F15:          "F15",
	F16:          "F16",
	F17:          "F17",
	F18:          "F18",
	F19:          "F19",
	F20:          "F20",
	F21:          "F21",
	F22:          "F22",
	F23:          "F23",
	F24:          "F24",
	F25:          "F25",
	LeftShift:    "LeftShift",
	LeftControl:  "LeftControl",
	LeftAlt:      "LeftAlt",
	LeftSuper:    "LeftSuper",
	RightShift:   "RightShift",
	RightControl: "RightControl",
	RightAlt:     "RightAlt",
	RightSuper:   "RightSuper",
	Menu:         "Menu",
}
