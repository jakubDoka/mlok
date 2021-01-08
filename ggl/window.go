package ggl

import (
	"gobatch/mt"
	"log"
	"time"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Setup is something that sets up the window drawing state, as library will support 2D and 3D different
// Setups can be used. its preferale to use lib setup with composition and override methods
type Setup interface {
	// VertexShader returns vertex shader source code of Setup
	FragmentShader() string
	// FragmentShader returns fragment shader source of Setup
	VertexShader() string
	// Buffer returns buffer that window will use
	Buffer() Buffer
	// Modify leaves space for some additional modification
	Modify(win *Window)
}

// Window enables use of almost all other structs avaliable in this package, Window initializes opengl context.
//
type Window struct {
	*glfw.Window
	Mask mt.RGBA
	Canvas

	cursorInsideWindow bool

	prevInp, currInp, tempInp struct {
		mouse   mt.V2
		buttons [KeyLast + 1]bool
		repeat  [KeyLast + 1]bool
		scroll  mt.V2
		typed   string
	}
}

// NWindow creates new window from WindowConfig, if config is nil, default one will  be used
func NWindow(config *WindowConfig) (*Window, error) {

	if config == nil {
		config = &WindowConfig{
			Setup:  Setup2D{},
			Width:  1000,
			Height: 600,
			Title:  "Hello there",
		}
	} else if config.Setup == nil {
		config.Setup = Setup2D{}
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	bti := func(b bool) int {
		if b {
			return glfw.True
		}
		return glfw.False
	}

	hints := append([]WindowHint{
		{glfw.ContextVersionMajor, 3},
		{glfw.ContextVersionMinor, 3},
		{glfw.OpenGLProfile, glfw.OpenGLCoreProfile},
		{glfw.OpenGLForwardCompatible, glfw.True},

		{glfw.Resizable, bti(config.Resizable)},
		{glfw.Decorated, bti(!config.Undecorated)},
		{glfw.Floating, bti(config.AlwaysOnTop)},
		{glfw.AutoIconify, bti(!config.NoIconify)},
		{glfw.TransparentFramebuffer, bti(config.TransparentFramebuffer)},
		{glfw.Maximized, bti(config.Maximized)},
		{glfw.Visible, bti(!config.Invisible)},
	}, config.Hints...)

	for _, h := range hints {
		glfw.WindowHint(h.Hint, h.Value)
	}

	window, err := glfw.CreateWindow(config.Width, config.Height, config.Title, config.Monitor, config.Share)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	win := Window{
		Window: window,
		Mask:   mt.Alpha(1),
	}

	win.initInput()

	program, err := NProgramFromSource(config.Setup.VertexShader(), config.Setup.FragmentShader())
	if err != nil {
		return nil, err
	}

	win.Canvas = *NCanvas(
		*RawTexture(int32(config.Width), int32(config.Height), nil, DefaultTextureConfig...),
		*program,
		config.Setup.Buffer(),
	)

	win.SetSizeCallback(func(w *glfw.Window, width, height int) {
		win.Canvas.Resize(mt.NAABB(0, 0, float64(width), float64(height)))
	})

	config.Setup.Modify(&win)

	return &win, nil
}

// Update updates the window so you can see it changing
func (w *Window) Update() {
	width, height := w.GetSize()
	w.Canvas.Render2D(mt.IM2, w.Mask, int32(width), int32(height))
	w.SwapBuffers()
	glfw.PollEvents()
}

// WindowConfig stores configuration for window, But that des not mean you cannot set these propertis
// later on, Its just nicer to have ewerithing on one place
type WindowConfig struct {
	Setup         Setup
	Width, Height int
	Title         string

	Resizable              bool
	Undecorated            bool
	TransparentFramebuffer bool
	AlwaysOnTop            bool
	NoIconify              bool
	Maximized              bool
	Invisible              bool

	Hints   []WindowHint
	Monitor *glfw.Monitor
	Share   *glfw.Window
}

// Pressed returns whether the Button is currently pressed down.
func (w *Window) Pressed(button Button) bool {
	return w.currInp.buttons[button]
}

// JustPressed returns whether the Button has just been pressed down.
func (w *Window) JustPressed(button Button) bool {
	return w.currInp.buttons[button] && !w.prevInp.buttons[button]
}

// JustReleased returns whether the Button has just been released up.
func (w *Window) JustReleased(button Button) bool {
	return !w.currInp.buttons[button] && w.prevInp.buttons[button]
}

// Repeated returns whether a repeat event has been triggered on button.
//
// Repeat event occurs repeatedly when a button is held down for some time.
func (w *Window) Repeated(button Button) bool {
	return w.currInp.repeat[button]
}

// MousePos returns the current mouse position in the Window's Bounds.
func (w *Window) MousePos() mt.V2 {
	return w.currInp.mouse
}

// MousePrevPos returns the previous mouse position in the Window's Bounds.
func (w *Window) MousePrevPos() mt.V2 {
	return w.prevInp.mouse
}

// MouseIsInside returns true if the mouse position is within the Window's Bounds.
func (w *Window) MouseIsInside() bool {
	return w.cursorInsideWindow
}

// MouseScroll returns the mouse scroll amount (in both axes) since the last call to Window.Update.
func (w *Window) MouseScroll() mt.V2 {
	return w.currInp.scroll
}

// Typed returns the text typed on the keyboard since the last call to Window.Update.
func (w *Window) Typed() string {
	return w.currInp.typed
}

// Button is a keyboard or mouse button. Why distinguish?
type Button int

// List of all mouse buttons.
const (
	MouseButton1      = Button(glfw.MouseButton1)
	MouseButton2      = Button(glfw.MouseButton2)
	MouseButton3      = Button(glfw.MouseButton3)
	MouseButton4      = Button(glfw.MouseButton4)
	MouseButton5      = Button(glfw.MouseButton5)
	MouseButton6      = Button(glfw.MouseButton6)
	MouseButton7      = Button(glfw.MouseButton7)
	MouseButton8      = Button(glfw.MouseButton8)
	MouseButtonLast   = Button(glfw.MouseButtonLast)
	MouseButtonLeft   = Button(glfw.MouseButtonLeft)
	MouseButtonRight  = Button(glfw.MouseButtonRight)
	MouseButtonMiddle = Button(glfw.MouseButtonMiddle)
)

// List of all keyboard buttons.
const (
	KeyUnknown      = Button(glfw.KeyUnknown)
	KeySpace        = Button(glfw.KeySpace)
	KeyApostrophe   = Button(glfw.KeyApostrophe)
	KeyComma        = Button(glfw.KeyComma)
	KeyMinus        = Button(glfw.KeyMinus)
	KeyPeriod       = Button(glfw.KeyPeriod)
	KeySlash        = Button(glfw.KeySlash)
	Key0            = Button(glfw.Key0)
	Key1            = Button(glfw.Key1)
	Key2            = Button(glfw.Key2)
	Key3            = Button(glfw.Key3)
	Key4            = Button(glfw.Key4)
	Key5            = Button(glfw.Key5)
	Key6            = Button(glfw.Key6)
	Key7            = Button(glfw.Key7)
	Key8            = Button(glfw.Key8)
	Key9            = Button(glfw.Key9)
	KeySemicolon    = Button(glfw.KeySemicolon)
	KeyEqual        = Button(glfw.KeyEqual)
	KeyA            = Button(glfw.KeyA)
	KeyB            = Button(glfw.KeyB)
	KeyC            = Button(glfw.KeyC)
	KeyD            = Button(glfw.KeyD)
	KeyE            = Button(glfw.KeyE)
	KeyF            = Button(glfw.KeyF)
	KeyG            = Button(glfw.KeyG)
	KeyH            = Button(glfw.KeyH)
	KeyI            = Button(glfw.KeyI)
	KeyJ            = Button(glfw.KeyJ)
	KeyK            = Button(glfw.KeyK)
	KeyL            = Button(glfw.KeyL)
	KeyM            = Button(glfw.KeyM)
	KeyN            = Button(glfw.KeyN)
	KeyO            = Button(glfw.KeyO)
	KeyP            = Button(glfw.KeyP)
	KeyQ            = Button(glfw.KeyQ)
	KeyR            = Button(glfw.KeyR)
	KeyS            = Button(glfw.KeyS)
	KeyT            = Button(glfw.KeyT)
	KeyU            = Button(glfw.KeyU)
	KeyV            = Button(glfw.KeyV)
	KeyW            = Button(glfw.KeyW)
	KeyX            = Button(glfw.KeyX)
	KeyY            = Button(glfw.KeyY)
	KeyZ            = Button(glfw.KeyZ)
	KeyLeftBracket  = Button(glfw.KeyLeftBracket)
	KeyBackslash    = Button(glfw.KeyBackslash)
	KeyRightBracket = Button(glfw.KeyRightBracket)
	KeyGraveAccent  = Button(glfw.KeyGraveAccent)
	KeyWorld1       = Button(glfw.KeyWorld1)
	KeyWorld2       = Button(glfw.KeyWorld2)
	KeyEscape       = Button(glfw.KeyEscape)
	KeyEnter        = Button(glfw.KeyEnter)
	KeyTab          = Button(glfw.KeyTab)
	KeyBackspace    = Button(glfw.KeyBackspace)
	KeyInsert       = Button(glfw.KeyInsert)
	KeyDelete       = Button(glfw.KeyDelete)
	KeyRight        = Button(glfw.KeyRight)
	KeyLeft         = Button(glfw.KeyLeft)
	KeyDown         = Button(glfw.KeyDown)
	KeyUp           = Button(glfw.KeyUp)
	KeyPageUp       = Button(glfw.KeyPageUp)
	KeyPageDown     = Button(glfw.KeyPageDown)
	KeyHome         = Button(glfw.KeyHome)
	KeyEnd          = Button(glfw.KeyEnd)
	KeyCapsLock     = Button(glfw.KeyCapsLock)
	KeyScrollLock   = Button(glfw.KeyScrollLock)
	KeyNumLock      = Button(glfw.KeyNumLock)
	KeyPrintScreen  = Button(glfw.KeyPrintScreen)
	KeyPause        = Button(glfw.KeyPause)
	KeyF1           = Button(glfw.KeyF1)
	KeyF2           = Button(glfw.KeyF2)
	KeyF3           = Button(glfw.KeyF3)
	KeyF4           = Button(glfw.KeyF4)
	KeyF5           = Button(glfw.KeyF5)
	KeyF6           = Button(glfw.KeyF6)
	KeyF7           = Button(glfw.KeyF7)
	KeyF8           = Button(glfw.KeyF8)
	KeyF9           = Button(glfw.KeyF9)
	KeyF10          = Button(glfw.KeyF10)
	KeyF11          = Button(glfw.KeyF11)
	KeyF12          = Button(glfw.KeyF12)
	KeyF13          = Button(glfw.KeyF13)
	KeyF14          = Button(glfw.KeyF14)
	KeyF15          = Button(glfw.KeyF15)
	KeyF16          = Button(glfw.KeyF16)
	KeyF17          = Button(glfw.KeyF17)
	KeyF18          = Button(glfw.KeyF18)
	KeyF19          = Button(glfw.KeyF19)
	KeyF20          = Button(glfw.KeyF20)
	KeyF21          = Button(glfw.KeyF21)
	KeyF22          = Button(glfw.KeyF22)
	KeyF23          = Button(glfw.KeyF23)
	KeyF24          = Button(glfw.KeyF24)
	KeyF25          = Button(glfw.KeyF25)
	KeyKP0          = Button(glfw.KeyKP0)
	KeyKP1          = Button(glfw.KeyKP1)
	KeyKP2          = Button(glfw.KeyKP2)
	KeyKP3          = Button(glfw.KeyKP3)
	KeyKP4          = Button(glfw.KeyKP4)
	KeyKP5          = Button(glfw.KeyKP5)
	KeyKP6          = Button(glfw.KeyKP6)
	KeyKP7          = Button(glfw.KeyKP7)
	KeyKP8          = Button(glfw.KeyKP8)
	KeyKP9          = Button(glfw.KeyKP9)
	KeyKPDecimal    = Button(glfw.KeyKPDecimal)
	KeyKPDivide     = Button(glfw.KeyKPDivide)
	KeyKPMultiply   = Button(glfw.KeyKPMultiply)
	KeyKPSubtract   = Button(glfw.KeyKPSubtract)
	KeyKPAdd        = Button(glfw.KeyKPAdd)
	KeyKPEnter      = Button(glfw.KeyKPEnter)
	KeyKPEqual      = Button(glfw.KeyKPEqual)
	KeyLeftShift    = Button(glfw.KeyLeftShift)
	KeyLeftControl  = Button(glfw.KeyLeftControl)
	KeyLeftAlt      = Button(glfw.KeyLeftAlt)
	KeyLeftSuper    = Button(glfw.KeyLeftSuper)
	KeyRightShift   = Button(glfw.KeyRightShift)
	KeyRightControl = Button(glfw.KeyRightControl)
	KeyRightAlt     = Button(glfw.KeyRightAlt)
	KeyRightSuper   = Button(glfw.KeyRightSuper)
	KeyMenu         = Button(glfw.KeyMenu)
	KeyLast         = Button(glfw.KeyLast)
)

// String returns a human-readable string describing the Button.
func (b Button) String() string {
	name, ok := buttonNames[b]
	if !ok {
		return "Invalid"
	}
	return name
}

var buttonNames = map[Button]string{
	MouseButton4:      "MouseButton4",
	MouseButton5:      "MouseButton5",
	MouseButton6:      "MouseButton6",
	MouseButton7:      "MouseButton7",
	MouseButton8:      "MouseButton8",
	MouseButtonLeft:   "MouseButtonLeft",
	MouseButtonRight:  "MouseButtonRight",
	MouseButtonMiddle: "MouseButtonMiddle",
	KeyUnknown:        "Unknown",
	KeySpace:          "Space",
	KeyApostrophe:     "Apostrophe",
	KeyComma:          "Comma",
	KeyMinus:          "Minus",
	KeyPeriod:         "Period",
	KeySlash:          "Slash",
	Key0:              "0",
	Key1:              "1",
	Key2:              "2",
	Key3:              "3",
	Key4:              "4",
	Key5:              "5",
	Key6:              "6",
	Key7:              "7",
	Key8:              "8",
	Key9:              "9",
	KeySemicolon:      "Semicolon",
	KeyEqual:          "Equal",
	KeyA:              "A",
	KeyB:              "B",
	KeyC:              "C",
	KeyD:              "D",
	KeyE:              "E",
	KeyF:              "F",
	KeyG:              "G",
	KeyH:              "H",
	KeyI:              "I",
	KeyJ:              "J",
	KeyK:              "K",
	KeyL:              "L",
	KeyM:              "M",
	KeyN:              "N",
	KeyO:              "O",
	KeyP:              "P",
	KeyQ:              "Q",
	KeyR:              "R",
	KeyS:              "S",
	KeyT:              "T",
	KeyU:              "U",
	KeyV:              "V",
	KeyW:              "W",
	KeyX:              "X",
	KeyY:              "Y",
	KeyZ:              "Z",
	KeyLeftBracket:    "LeftBracket",
	KeyBackslash:      "Backslash",
	KeyRightBracket:   "RightBracket",
	KeyGraveAccent:    "GraveAccent",
	KeyWorld1:         "World1",
	KeyWorld2:         "World2",
	KeyEscape:         "Escape",
	KeyEnter:          "Enter",
	KeyTab:            "Tab",
	KeyBackspace:      "Backspace",
	KeyInsert:         "Insert",
	KeyDelete:         "Delete",
	KeyRight:          "Right",
	KeyLeft:           "Left",
	KeyDown:           "Down",
	KeyUp:             "Up",
	KeyPageUp:         "PageUp",
	KeyPageDown:       "PageDown",
	KeyHome:           "Home",
	KeyEnd:            "End",
	KeyCapsLock:       "CapsLock",
	KeyScrollLock:     "ScrollLock",
	KeyNumLock:        "NumLock",
	KeyPrintScreen:    "PrintScreen",
	KeyPause:          "Pause",
	KeyF1:             "F1",
	KeyF2:             "F2",
	KeyF3:             "F3",
	KeyF4:             "F4",
	KeyF5:             "F5",
	KeyF6:             "F6",
	KeyF7:             "F7",
	KeyF8:             "F8",
	KeyF9:             "F9",
	KeyF10:            "F10",
	KeyF11:            "F11",
	KeyF12:            "F12",
	KeyF13:            "F13",
	KeyF14:            "F14",
	KeyF15:            "F15",
	KeyF16:            "F16",
	KeyF17:            "F17",
	KeyF18:            "F18",
	KeyF19:            "F19",
	KeyF20:            "F20",
	KeyF21:            "F21",
	KeyF22:            "F22",
	KeyF23:            "F23",
	KeyF24:            "F24",
	KeyF25:            "F25",
	KeyKP0:            "KP0",
	KeyKP1:            "KP1",
	KeyKP2:            "KP2",
	KeyKP3:            "KP3",
	KeyKP4:            "KP4",
	KeyKP5:            "KP5",
	KeyKP6:            "KP6",
	KeyKP7:            "KP7",
	KeyKP8:            "KP8",
	KeyKP9:            "KP9",
	KeyKPDecimal:      "KPDecimal",
	KeyKPDivide:       "KPDivide",
	KeyKPMultiply:     "KPMultiply",
	KeyKPSubtract:     "KPSubtract",
	KeyKPAdd:          "KPAdd",
	KeyKPEnter:        "KPEnter",
	KeyKPEqual:        "KPEqual",
	KeyLeftShift:      "LeftShift",
	KeyLeftControl:    "LeftControl",
	KeyLeftAlt:        "LeftAlt",
	KeyLeftSuper:      "LeftSuper",
	KeyRightShift:     "RightShift",
	KeyRightControl:   "RightControl",
	KeyRightAlt:       "RightAlt",
	KeyRightSuper:     "RightSuper",
	KeyMenu:           "Menu",
}

func (w *Window) initInput() {
	w.Window.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			w.tempInp.buttons[Button(button)] = true
		case glfw.Release:
			w.tempInp.buttons[Button(button)] = false
		}
	})

	w.Window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyUnknown {
			return
		}
		switch action {
		case glfw.Press:
			w.tempInp.buttons[Button(key)] = true
		case glfw.Release:
			w.tempInp.buttons[Button(key)] = false
		case glfw.Repeat:
			w.tempInp.repeat[Button(key)] = true
		}
	})

	w.Window.SetCursorEnterCallback(func(_ *glfw.Window, entered bool) {
		w.cursorInsideWindow = entered
	})

	w.Window.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		w.tempInp.mouse = mt.NV2(x, float64(w.H)-y)
	})

	w.Window.SetScrollCallback(func(_ *glfw.Window, xoff, yoff float64) {
		w.tempInp.scroll.X += xoff
		w.tempInp.scroll.Y += yoff
	})

	w.Window.SetCharCallback(func(_ *glfw.Window, r rune) {
		w.tempInp.typed += string(r)
	})
}

// UpdateInput polls window events. Call this function to poll window events
// without swapping buffers. Note that the Update method invokes UpdateInput.
func (w *Window) UpdateInput() {
	mainthread.Call(func() {
		glfw.PollEvents()
	})
	w.doUpdateInput()
}

// UpdateInputWait blocks until an event is received or a timeout. If timeout is 0
// then it will wait indefinitely
func (w *Window) UpdateInputWait(timeout time.Duration) {
	mainthread.Call(func() {
		if timeout <= 0 {
			glfw.WaitEvents()
		} else {
			glfw.WaitEventsTimeout(timeout.Seconds())
		}
	})
	w.doUpdateInput()
}

// internal input bookkeeping
func (w *Window) doUpdateInput() {
	w.prevInp = w.currInp
	w.currInp = w.tempInp

	w.tempInp.repeat = [KeyLast + 1]bool{}
	w.tempInp.scroll = mt.V2{}
	w.tempInp.typed = ""

}

// WindowHint allows to specify vindow hints (revolution right here)
type WindowHint struct {
	Hint  glfw.Hint
	Value int
}

// Setup2D is basic 2D game rendering setup, if you would like to change some part of it just embed it
// in your setup struct and override methods
type Setup2D struct{}

// VertexShader implements Setup interface
func (s Setup2D) VertexShader() string {
	return `
	#version 330
	
	layout (location = 0) in vec2 vert;
	layout (location = 1) in vec2 tex;
	layout (location = 2) in vec4 mask;
	
	uniform mat3 camera2D;
	uniform vec2 viewportSize;
	uniform vec2 textureSize;
	
	out vec2 fragTex;
	out vec4 fragMask;
	void main() {
		fragMask = mask;
		fragTex = tex/textureSize;
		gl_Position = vec4(camera2D * vec3(vert/viewportSize, 0), 1);
	}
	`
}

// FragmentShader implements Setup interface
func (s Setup2D) FragmentShader() string {
	return `
	#version 330

	uniform sampler2D tex;
	uniform int useTexture;

	in vec2 fragTex;
	in vec4 fragMask;

	out vec4 outputColor;

	void main() {
		if (useTexture == 1) {
			outputColor = texture(tex, fragTex) * fragMask;
		} else {
			outputColor = fragMask;
		}
	}
	`
}

// Program implements Setup interface
func (s Setup2D) Program() Program {
	prog, err := NProgramFromSource(s.VertexShader(), s.FragmentShader())
	if err != nil {
		panic(err)
	}

	return *prog
}

// Buffer implements Setup interface
func (s Setup2D) Buffer() Buffer {
	return NBuffer(2, 2, 4)
}

// Modify implements Setup interface
func (s Setup2D) Modify(win *Window) {
	win.SetCamera2D(&mt.IM2)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
