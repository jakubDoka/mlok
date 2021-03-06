package ggl

import (
	"gobatch/ggl/key"
	"gobatch/mat"
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Setup is something that sets up the window drawing state, as library will support  and 3D different
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
	Mask mat.RGBA
	Canvas

	cursorInsideWindow bool

	prevInp, currInp, tempInp struct {
		mouse   mat.Vec
		buttons [key.Last + 1]bool
		repeat  [key.Last + 1]bool
		scroll  mat.Vec
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
		Mask:   mat.Alpha(1),
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
		win.Canvas.Resize(mat.A(0, 0, float64(width), float64(height)))
	})

	config.Setup.Modify(&win)

	return &win, nil
}

// Update updates the window so you can see it changing and also updates input
func (w *Window) Update() {
	width, height := w.GetSize()
	w.Canvas.RenderToScreen(mat.IM, w.Mask, int32(width), int32(height))
	w.SwapBuffers()
	glfw.PollEvents()
	w.doUpdateInput()
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
func (w *Window) Pressed(button key.Key) bool {
	return w.currInp.buttons[button]
}

// JustPressed returns whether the Button has just been pressed down.
func (w *Window) JustPressed(button key.Key) bool {
	return w.currInp.buttons[button] && !w.prevInp.buttons[button]
}

// JustReleased returns whether the Button has just been released up.
func (w *Window) JustReleased(button key.Key) bool {
	return !w.currInp.buttons[button] && w.prevInp.buttons[button]
}

// Repeated returns whether a repeat event has been triggered on button.
//
// Repeat event occurs repeatedly when a button is held down for some time.
func (w *Window) Repeated(button key.Key) bool {
	return w.currInp.repeat[button]
}

// MousePos returns the current mouse position in the Window's Bounds.
func (w *Window) MousePos() mat.Vec {
	return w.currInp.mouse
}

// MousePrevPos returns the previous mouse position in the Window's Bounds.
func (w *Window) MousePrevPos() mat.Vec {
	return w.prevInp.mouse
}

// MouseIsInside returns true if the mouse position is within the Window's Bounds.
func (w *Window) MouseIsInside() bool {
	return w.cursorInsideWindow
}

// MouseScroll returns the mouse scroll amount (in both axes) since the last call to Window.Update.
func (w *Window) MouseScroll() mat.Vec {
	return w.currInp.scroll
}

// Typed returns the text typed on the keyboard since the last call to Window.Update.
func (w *Window) Typed() string {
	return w.currInp.typed
}

func (w *Window) initInput() {
	w.Window.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			w.tempInp.buttons[key.Key(button)] = true
		case glfw.Release:
			w.tempInp.buttons[key.Key(button)] = false
		}
	})

	w.Window.SetKeyCallback(func(_ *glfw.Window, k glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if k == glfw.KeyUnknown {
			return
		}
		switch action {
		case glfw.Press:
			w.tempInp.buttons[key.Key(k)] = true
		case glfw.Release:
			w.tempInp.buttons[key.Key(k)] = false
		case glfw.Repeat:
			w.tempInp.repeat[key.Key(k)] = true
		}
	})

	w.Window.SetCursorEnterCallback(func(_ *glfw.Window, entered bool) {
		w.cursorInsideWindow = entered
	})

	w.Window.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		w.tempInp.mouse = mat.V(x, y).Sub(w.viewpot.Scaled(.5)).Mul(mat.V(1, -1))
	})

	w.Window.SetScrollCallback(func(_ *glfw.Window, xoff, yoff float64) {
		w.tempInp.scroll.X += xoff
		w.tempInp.scroll.Y += yoff
	})

	w.Window.SetCharCallback(func(_ *glfw.Window, r rune) {
		w.tempInp.typed += string(r)
	})
}

// internal input bookkeeping
func (w *Window) doUpdateInput() {
	w.prevInp = w.currInp
	w.currInp = w.tempInp

	w.tempInp.repeat = [key.Last + 1]bool{}
	w.tempInp.scroll = mat.Vec{}
	w.tempInp.typed = ""
}

// WindowHint allows to specify vindow hints (revolution right here)
type WindowHint struct {
	Hint  glfw.Hint
	Value int
}
