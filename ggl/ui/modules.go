package ui

import (
	"math"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/drw"
	"github.com/jakubDoka/mlok/ggl/key"
	"github.com/jakubDoka/mlok/ggl/txt"
	"github.com/jakubDoka/mlok/logic/timer"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/rgba"

	"github.com/atotto/clipboard"
	"github.com/jakubDoka/gogen/str"
	"github.com/jakubDoka/goml/goss"
)

// Area is a text input element, you can get its content by turning Area.Content into string.
//
// style:
//	cursor_drawer:			custom_type // something that draws the cursor, it is taken by name from assets
//	cursor_thickness:		float		// how thick the cursor is
//	cursor_mask:			rgba		// mask of cursor
//	cursor_blink_frequency:	float		// how often cursor blinks
//	auto_frequency:			float		// when you hold some button that controls the input, action starts
//										// repeating and this sets how often it repeats
//	hold_responce_speed:	float		// how long you have to hold on to button until it starts repeating
type Area struct {
	Text
	HoldMap

	drw CursorDrawer

	selected, dirty, noEffects, shown bool

	LineIdx, Line int

	Blinker         timer.Timer
	CursorThickness float64
	CursorMask      mat.RGBA
}

// New implements ModuleFactory interface
func (a *Area) New() Module {
	return &Area{}
}

// Init implements Module interface
func (a *Area) Init(e *Element) {
	a.Composed = true // important for next call
	a.Text.Init(e)

	a.drw = a.CursorDrawer("cursor_drawer", a.Scene.Assets.Cursors, defaultCursor{})
	a.CursorThickness = e.Float("cursor_thickness", 2)
	a.CursorMask = e.RGBA("cursor_mask", mat.White)
	a.Blinker = timer.Period(e.Float("cursor_blinking_frequency", .6))

	a.AutoFrequency = e.Float("auto_frequency", .03)
	a.HoldResponceSpeed = e.Float("hold_responce_speed", .5)

	a.binds = map[key.Key]float64{}
}

// Update implements Module interface
func (a *Area) Update(w *ggl.Window, delta float64) {
	// Text.Update sets up lot of things
	a.Text.Update(w, delta)
	if a.Blinker.Period < 0 || w.Pressed(key.MouseLeft) && a.Start != a.End {
		a.shown = true
	} else {
		if a.Blinker.TickDoneReset(delta) {
			a.shown = !a.shown
			a.Scene.Redraw.Notify()
		}
	}

	if !a.selected && !a.Hovering {
		return
	}

	// we don't want effects to be applied when user is editing text
	if w.JustPressed(key.MouseLeft) {
		if a.selected && !a.Hovering {
			a.NoEffects = a.noEffects
			a.selected = false
			a.Dirty()
			a.Events.Invoke(Deselect, nil)
		} else if !a.selected && a.Hovering {
			a.noEffects = a.NoEffects
			a.NoEffects = true
			a.selected = true
			a.Dirty()
			a.Events.Invoke(Select, nil)
			// returning as text has to get redrawn and reindexed so we can restore the cursor
			return
		}
	}

	if !a.selected {
		return
	}

	// it is shortest way to handle arrow navigation
	if (a.Line > 0 && a.Hold(key.Up, w, delta, func() {
		a.Line--
		a.Start = a.ProjectLine(a.LineIdx, a.Line)
	})) || (a.Line < a.Lines()-1 && a.Hold(key.Down, w, delta, func() {
		a.Line++
		a.Start = a.ProjectLine(a.LineIdx, a.Line)
	})) || (a.Start > 0 && a.Hold(key.Left, w, delta, func() {
		a.Start--
	})) || (a.Start < len(a.Content) && a.Hold(key.Right, w, delta, func() {
		a.Start++
	})) || a.dirty {
		a.LineIdx, a.Line = a.UnprojectLine(a.Start)
		a.End = a.Start
		a.dirty = false
		a.Blinker.Reset()
		a.shown = true
		a.Scene.Redraw.Notify()
	}

	typed := w.Typed()

	// cut paste thing
	var cut bool
	if w.Pressed(key.LeftControl) {
		cut = w.JustPressed(key.X)
		if a.Start != a.End && cut {
			a.Scene.Log(a.Element, a.Clip(a.Start, a.End))
		} else if w.JustPressed(key.V) {
			var err error
			typed, err = clipboard.ReadAll()
			a.Scene.Log(a.Element, err)
		}
	}

	if !a.Hold(key.Enter, w, delta, func() {
		typed = "\n"
		a.Events.Invoke(Enter, nil)
	}) && !a.Hold(key.Backspace, w, delta, func() {
		if a.Start != 0 && a.Start == a.End {
			a.Start--
		}
	}) && !a.Hold(key.Tab, w, delta, func() {
		typed = "\t"
	}) && typed == "" && !cut {
		return
	}

	if cut { // we don't want to accidentally write something
		typed = ""
	}

	a.Blinker.Reset()
	a.shown = true

	nv := str.NString(typed)
	if a.Start > a.End {
		a.Start, a.End = a.End, a.Start
	}
	a.Content.RemoveSlice(a.Start, a.End)
	a.Content.InsertSlice(a.Start, nv)
	a.Events.Invoke(TextChanged, typed)
	a.Start += len(nv)
	a.End = a.Start
	a.Dirty()
	a.dirty = true
}

// DrawOnTop implements Module interface
func (a *Area) DrawOnTop(tg ggl.Target, canvas *drw.Geom) {
	if a.selected && a.shown {
		a.drw.Draw(
			tg,
			canvas,
			a.Dot(mat.Maxi(a.Start, a.End)).Sub(mat.V(0, a.Descent*a.Scl.Y)),
			mat.V(a.CursorThickness, a.LineHeight*a.Scl.Y),
			a.CursorMask,
		)
		canvas.Clear()
	}
	a.Text.DrawOnTop(tg, canvas)
}

// Dirty is similar to Text.Dirty but it preserves the Start and End.
func (a *Area) Dirty() {
	start, end := a.Start, a.End
	a.Text.Dirty()
	a.Start, a.End = start, end
}

// CursorDrawer is something that draws the cursor inside the Area when area is selected.
type CursorDrawer interface {
	Draw(t ggl.Target, canvas *drw.Geom, base, size mat.Vec, mask mat.RGBA)
}

type defaultCursor struct{}

func (d defaultCursor) Draw(t ggl.Target, canvas *drw.Geom, base, size mat.Vec, mask mat.RGBA) {
	size.X *= .5
	canvas.Color(mask).AABB(mat.A(base.X-size.X, base.Y, base.X+size.X, base.Y+size.Y))
	canvas.Fetch(t)
}

// HoldMap is extracted behavior of text edit, see HoldMat.Hold
type HoldMap struct {
	HoldResponceSpeed, AutoFrequency float64

	binds map[key.Key]float64
}

// Hold supports Hold and repeat effect, if you hold button for long enough, action starts repeating.
func (h HoldMap) Hold(b key.Key, win *ggl.Window, delta float64, do func()) bool {
	if win.JustPressed(b) {
		do()
		return true
	} else if win.JustReleased(b) {
		h.binds[b] = 0
	} else if win.Pressed(b) {
		tm := h.binds[b] + delta
		if tm > h.HoldResponceSpeed+h.AutoFrequency {
			do()
			h.binds[b] = h.HoldResponceSpeed
			return true
		}
		h.binds[b] = tm
	}
	return false
}

// Button is a button, has lot of space for customization, it can be in three states:
// idle, hover and pressed, each mode can have different message, texture, padding and mask
// use of texture is optional and button uses Patch for more flexibility, whe button is
// initialized, it creates its own Text element witch visible text
//
// style:
// 	all_text: 								string		// sets text on all state
//  idle/hover/pressed/disabled+_text: 		string		// sets text for each state
//	all_masks: 								rgba		// sets mask on all states
//  idle/hover/pressed/disabled+_mask: 		rgba		// sets mask for each state
//  all_regions:                			aabb|name   // sets region on all states
// 	idle/hover/pressed/disabled+_region:	aabb|name	// sets region for each state
//  all_padding:                			aabb	   	// sets padding on all states
// 	idle/hover/pressed/disabled+_padding:	aabb		// sets padding for each state
type Button struct {
	Patch

	Text   Text
	States [len(buttonStates)]ButtonState

	Current            ButtonStateEnum
	selected, Disabled bool
}

// New implements ModuleFactory interface
func (b *Button) New() Module {
	return &Button{}
}

// Init implements Module interface
func (b *Button) Init(e *Element) {
	b.Patch.Init(e)

	text := b.Raw.Attributes.Ident("all_text", "")
	parsed := str.NString(text)
	mask := b.RGBA("all_masks", mat.White)
	region := e.Region("all_regions", e.Scene.Assets.Regions, mat.ZA)
	padding := b.AABB("all_padding", mat.ZA)
	for i := range b.States {
		bs := &b.States[i]
		bs.Mask = mask
		bs.Region = region
		bs.Padding = padding
		bs.Text = parsed
	}

	for i, s := range buttonStates {
		bs := &b.States[i]
		bs.Text = str.NString(b.Raw.Attributes.Ident(s+"_text", text))
		bs.Mask = b.RGBA(s+"_mask", bs.Mask)
		bs.Region = e.Region(s+"_region", e.Scene.Assets.Regions, bs.Region)
		bs.Padding = b.AABB(s+"_padding", bs.Padding)
	}

}

func (b *Button) PostInit() {
	if b.ChildCount() == 0 {
		textElem := NElement()
		textElem.Module = &b.Text
		b.AddChild("buttonText", textElem)
	} else {
		for i := 0; i < b.ChildCount(); i++ {
			ch := b.ChildAt(i)
			val, ok := ch.Module.(*Text)
			if ok {
				// copy the module to this text and change address
				// also apply text to states
				b.Text = *val
				for i := range b.States {
					if len(b.States[i].Text) == 0 {
						b.States[i].Text = b.Text.Content
					}
				}
				ch.Module = &b.Text
				ch.SetName("buttonText")
				break
			}
		}
	}

	b.Current = None
	b.ApplyState(Idle)
}

// Update implements Module interface
func (b *Button) Update(w *ggl.Window, delta float64) {
	if b.Disabled {
		b.ApplyState(Disabled)
		return
	}

	if !b.Hovering {
		b.ApplyState(Idle)
		b.selected = false
		return
	}

	if w.JustPressed(key.MouseLeft) {
		b.selected = true
	}

	if w.JustReleased(key.MouseLeft) {
		if b.selected {
			b.Events.Invoke(Click, nil)
		}
		b.selected = false
	}

	if b.selected {
		b.ApplyState(Pressed)
	} else {
		b.ApplyState(Hover)
	}
}

// ApplyState applies the button state by index
func (b *Button) ApplyState(state ButtonStateEnum) {
	if b.Current == state {
		return
	}
	b.Current = state
	bs := &b.States[state]
	b.Patch.Padding = bs.Padding
	b.Patch.SetRegion(bs.Region)
	b.Patch.Mask = bs.Mask
	b.Text.Content = bs.Text
	b.Text.Dirty()
}

// SetText sets text on all states to given value
func (b *Button) SetText(text string) {
	str := str.NString(text)
	for i := range b.States {
		b.States[i].Text = str
	}
}

// ButtonState ...
type ButtonState struct {
	Mask            mat.RGBA
	Region, Padding mat.AABB
	Text            str.String
}

type ButtonStateEnum uint8

const (
	Idle ButtonStateEnum = iota
	Hover
	Pressed
	Disabled
	None
)

var buttonStates = [...]string{"idle", "hover", "pressed", "disabled"}

// Patch is similar tor Sprite but uses ggl.Patch instead thus has more styling options
//
// style:
//	region:			aabb|name	// rendert texture region
//  padding:		aabb	 	// patch padding
//	patch_mask:		rgba		// path mask, if there is no region it acts as background color
//  patch_scale:	vec			// because Patch is more flexible, specifying scale can create more variations
type Patch struct {
	ModuleBase

	Patch           ggl.Patch
	Padding, Region mat.AABB
	Mask            mat.RGBA
	Scale           mat.Vec
}

// New implements ModuleFactory interface
func (p *Patch) New() Module {
	return &Patch{}
}

// Init implements Module interface
func (p *Patch) Init(e *Element) {
	p.ModuleBase.Init(e)
	p.Mask = e.RGBA("patch_mask", mat.White)
	p.Region = e.Region("region", e.Scene.Assets.Regions, mat.ZA)

	w, h := p.Region.W()*.5, p.Region.H()*.5
	p.Padding = e.AABB("patch_padding", mat.A(w, h, w, h))

	p.Scale = e.Vec("patch_scale", mat.V(1, 1))

	p.SetRegion(p.Region)
}

// Draw implements Module interface
func (p *Patch) Draw(t ggl.Target, canvas *drw.Geom) {
	p.Patch.Fetch(t)
}

// OnFrameChange implements Module interface
func (p *Patch) OnFrameChange() {
	size := p.Frame.Size().Div(p.Scale)
	p.Patch.Resize(size.X, size.Y)
	p.Patch.Update(mat.M(p.Frame.Center(), p.Scale, 0), p.Mask)
}

// SetRegion sets patch region, and deals with all technical issues
func (p *Patch) SetRegion(value mat.AABB) {
	if value == mat.ZA {
		p.Patch.SetIntensity(0)
	} else {
		p.Patch = ggl.NPatch(value, p.Padding)
	}
}

// SetPadding sets patch padding and deals with all technical issues
func (p *Patch) SetPadding(value mat.AABB) {
	p.Padding = value
	p.SetRegion(p.Region)
}

// Sprite is sprite for ui elements
// style:
//	mask:	rgba		// texture modulation
//	region:	aabb|name	// from where texture should be sampled from
type Sprite struct {
	ModuleBase
	Sprite ggl.Sprite

	Mask mat.RGBA
}

// New implements ModuleFactory interface
func (s *Sprite) New() Module {
	return &Sprite{}
}

// Init implements Module interface
func (s *Sprite) Init(e *Element) {
	s.ModuleBase.Init(e)
	s.Mask = e.RGBA("sprite_mask", mat.White)
	reg := e.Region("region", s.Scene.Assets.Regions, mat.ZA)

	if reg == mat.ZA {
		s.Sprite.SetIntensity(0) // sprite will act as normal div
	} else {
		s.Sprite = ggl.NSprite(reg)
	}
	s.Sprite.SetColor(s.Mask)
}

// OnFrameChange implements Module interface
func (s *Sprite) OnFrameChange() {
	s.Sprite.SetDist(s.Frame)
}

// Draw implements Module interface
func (s *Sprite) Draw(t ggl.Target, canvas *drw.Geom) {
	s.Sprite.Fetch(t)
}

// Scroll can make element visible trough scrollable viewport
//
// style:
//	bar_width:			float	// sets bar thickness
//	friction:			float	// higher the value, less sliding, negative == no sliding
//	scroll_sensitivity:	float	// how match will scroll slide
//	bar_color:			rgba	// bar handle color
//	rail_color:			rgba	// bar rail color
//	intersection_color:	rgba	// if both bars are active, rectangle appears in a corner
//	bar_x/bar_y/bars:	bool	// makes bars visible and active
type Scroll struct {
	ModuleBase
	drw.SpriteViewport

	BarWidth, Friction, ScrollSensitivity  float64
	BarColor, RailColor, IntersectionColor mat.RGBA
	X, Y                                   Bar

	offset, vel, ratio, corner mat.Vec
	dirty, useVel              bool
}

// New implements ModuleFactory interface
func (s *Scroll) New() Module {
	return &Scroll{}
}

// Init implements module interface
func (s *Scroll) Init(e *Element) {
	s.ModuleBase.Init(e)
	s.Proc = &s.SpriteViewport
	s.BarWidth = s.Float("bar_width", 20)
	s.Friction = s.Float("friction", -1) // instant
	s.ScrollSensitivity = s.Float("scroll_sensitivity", 30)
	s.BarColor = s.RGBA("bar_color", rgba.White)
	s.RailColor = s.RGBA("rail_color", mat.RGBA{})
	s.IntersectionColor = s.RGBA("intersection_color", mat.RGBA{})
	s.X.Use = s.Bool("bar_x", false)
	s.Y.Use = s.Bool("bar_y", false)
	if !s.Y.Use && !s.X.Use {
		c := s.Bool("bars", false)
		s.X.Use = c
		s.Y.Use = c
	}
	s.X.position = 1 // to prevent snap
	s.Y.position = 1
}

// DrawOnTop implements module interface
func (s *Scroll) DrawOnTop(t ggl.Target, c *drw.Geom) {
	if s.X.use {
		rect := mat.AABB{Min: s.Frame.Min, Max: s.corner}
		c.Color(s.RailColor).AABB(rect)
		rect.Min.X, rect.Max.X = s.barBounds(&s.X, 0)
		c.Color(s.BarColor).AABB(rect)
	}
	if s.Y.use {
		rect := mat.AABB{Min: s.corner, Max: s.Frame.Max}
		c.Color(s.RailColor).AABB(rect)
		rect.Min.Y, rect.Max.Y = s.barBounds(&s.Y, 1)
		c.Color(s.BarColor).AABB(rect)
	}

	if s.X.use && s.Y.use {
		rect := mat.A(s.corner.X, s.Frame.Min.Y, s.Frame.Max.X, s.corner.Y)
		c.Color(s.IntersectionColor).AABB(rect)
	}

	c.Fetch(t)
}

// Update implements module interface
func (s *Scroll) Update(w *ggl.Window, delta float64) {
	if s.dirty {
		s.dirty = false
		s.update()
		s.updateOffset()
		s.move(s.Frame.Min.Sub(s.margin.Min).Sub(s.Offest), false)
		s.Scene.Redraw.Notify()
	}

	if s.useVel {
		s.useVel = s.vel.Len2() > .01
		if !s.X.selected && !s.Y.selected {
			s.offset.AddE(s.vel)
		}
		if s.Friction < 0 {
			s.vel = mat.ZV
		} else {
			// we don't want to get crazy if frames are low
			s.vel.SubE(s.vel.Scaled(math.Min(s.Friction*delta, 1)))
		}
		s.dirty = true

	}

	if !s.Hovering {
		return
	}

	scroll := w.MouseScroll()
	if scroll.Y != 0 {
		s.vel.Y -= scroll.Y * s.ScrollSensitivity
		s.useVel = true
	}

	if w.JustReleased(key.MouseLeft) {
		s.X.selected = false
		s.Y.selected = false
		return
	}

	if !w.Pressed(key.MouseLeft) {
		return
	}

	var (
		mouse  = w.MousePrevPos()
		move   = mouse.To(w.MousePos())
		as, ae = s.barBounds(&s.X, 0)
		bs, be = s.barBounds(&s.Y, 1)
	)

	if s.X.Use && s.X.use && s.X.selected || (s.corner.Y > mouse.Y && mouse.X >= as && mouse.X <= ae) {
		if w.JustPressed(key.MouseLeft) {
			s.X.selected = true
		}
		if s.X.selected {
			s.X.Move(-move.X)
			s.dirty = true
		}
	} else if !s.Scene.TextSelected {
		s.vel.X = move.X
		s.useVel = true
	}

	if s.Y.Use && s.Y.use && s.Y.selected || (s.corner.X < mouse.X && mouse.Y >= bs && mouse.Y <= be) {
		if w.JustPressed(key.MouseLeft) {
			s.Y.selected = true
		}
		if s.Y.selected {
			s.Y.Move(move.Y)
			s.dirty = true
		}
	} else if !s.Scene.TextSelected {
		s.vel.Y = move.Y
		s.useVel = true
	}

}

func (s *Scroll) barBounds(b *Bar, side int) (float64, float64) {
	prj := b.reminder * b.position
	if side == 0 {
		prj = -prj - b.length
	}
	prj += s.corner.Flatten()[side]

	return prj, prj + b.length
}

// OnFrameChange implements module interface
func (s *Scroll) OnFrameChange() {
	size := s.Frame.Size()
	s.ratio = s.ChildSize.Sub(size)
	s.corner = mat.V(s.Frame.Max.X, s.Frame.Min.Y)

	ratio := s.ratio
	if s.X.Use {
		ratio.Y += s.BarWidth
		s.X.space = size.X
	}
	if s.Y.Use {
		ratio.X += s.BarWidth
		s.Y.space = size.Y
	}

	s.X.use = s.X.Use && ratio.X > 0
	s.Y.use = s.Y.Use && ratio.Y > 0

	if s.X.use && s.Y.use {
		s.X.space -= s.BarWidth
		s.ratio.Y += s.BarWidth
		s.Y.space -= s.BarWidth
		s.ratio.X += s.BarWidth
	}

	if s.X.use {
		s.X.CalcRatio(s.ChildSize.X)
		s.corner.Y += s.BarWidth
	}
	if s.Y.use {
		s.Y.CalcRatio(s.ChildSize.Y)
		s.corner.X -= s.BarWidth
	}

	s.SpriteViewport.Area = s.Frame
}

// move applies velocity to offset
func (s *Scroll) update() {
	dif := s.ratio.Inv()
	if dif.X < 0 {
		if s.X.Use && (!s.useVel || s.X.selected) {
			s.offset.X = dif.X * (1 - s.X.position) // needs to be inverted or it ll look unnatural
		} else {
			s.offset.X = mat.Clamp(s.offset.X, dif.X, 0)
			s.X.position = 1 - s.offset.X/dif.X // make sure to move bar too
		}
	} else {
		s.offset.X = 0
	}

	if dif.Y < 0 {
		if s.Y.Use && (!s.useVel || s.Y.selected) {
			s.offset.Y = dif.Y * s.Y.position
		} else {
			s.offset.Y = mat.Clamp(s.offset.Y, dif.Y, 0)
			s.Y.position = s.offset.Y / dif.Y
		}
	} else {
		s.offset.Y = 0
	}
}

// move moves all elements by scroll offset
func (s *Scroll) updateOffset() {
	ch := s.children.Slice()
	off := s.offset
	if s.X.use && s.Y.use { // have to shift whole thing because of strange dimensions
		off.Y += s.BarWidth
	}
	for i := 0; i < len(ch); i++ {
		ch[i].Value.Offest = off
	}
}

// Bar holds information about scroll bar
type Bar struct {
	Use      bool
	position float64

	use, selected           bool
	length, space, reminder float64
}

// CalcRatio calculates length and reminder of bar
func (b *Bar) CalcRatio(size float64) {
	b.length = math.Max(b.space*b.space/size, 5)
	b.reminder = b.space - b.length
}

// Move moves the bar
func (b *Bar) Move(vel float64) {
	b.position = mat.Clamp(b.position+vel/b.reminder, 0, 1)
}

// Text handles text rendering
//
// style:
//	text_scale:				vec						// size of font
//	text_color:				rgba					// color by with text is pre-multiplied
//	text_size:				vec						// works as element size
//	text_margin:			aabb					// works as element margin
//  text_padding:			aabb					// works as element padding
//	text_background:		rgba					// works as moduleBase background
//	text_selection_color:	rgba					// color if text selection
//	text_align:				float|left|middle|right	// text align
//	text_no_effects:		bool					// makes text effects like color and differrent fonts disabled
//	text_markdown:			name					// sets a markdown that text will use to render
type Text struct {
	ModuleBase
	txt.Paragraph
	*txt.Markdown
	dirty, Composed, selected bool
	SelectionColor            mat.RGBA
	Start, End, LineIdx, Line int
}

// New implements module factory interface
func (t *Text) New() Module {
	return &Text{}
}

// DefaultStyle implements Module interface
func (t *Text) DefaultStyle() goss.Style {
	return goss.Style{
		"text_scale":           {"inherit"},
		"text_color":           {"inherit"},
		"text_size":            {"inherit"},
		"text_margin":          {"inherit"},
		"text_padding":         {"inherit"},
		"text_background":      {"inherit"},
		"text_selection_color": {"inherit"},
		"text_align":           {"inherit"},
		"text_no_effects":      {"inherit"},
		"text_markdown":        {"inherit"},
	}
}

// Init implements Module interface
func (t *Text) Init(e *Element) {
	t.ModuleBase.Init(e)

	ident := t.Ident("text_markdown", "default")
	if ident == "inherit" {
		ident = "default"
	}
	mkd, ok := t.Scene.Assets.Markdowns[ident]
	if !ok {
		panic(t.Path() + ": markdown with name '" + ident + "' is not present in assets")
	}

	t.Markdown = mkd
	t.Align = t.Props.Align("text_align", txt.Left)
	t.Scl = t.Vec("text_scale", mat.V(1, 1))
	t.Mask = t.RGBA("text_color", mat.White)
	t.SelectionColor = t.RGBA("text_selection_color", mat.Alpha(.5))
	if !t.Composed {
		t.Props.Size = t.Vec("text_size", mat.ZV)
		t.Props.Margin = t.AABB("text_margin", mat.A(4, 4, 4, 4))
		t.Background = t.RGBA("text_background", t.Background)
		t.Props.Padding = t.AABB("text_padding", mat.ZA)
	}
	t.NoEffects = t.Bool("text_no_effects", false)
	t.Content = str.NString(t.Raw.Attributes.Ident("text", string(t.Content)))

	t.Dirty()
}

// Draw implements Module interface
func (t *Text) Draw(tr ggl.Target, g *drw.Geom) {
	t.ModuleBase.Draw(tr, g)
	t.Paragraph.Draw(tr)
}

// Update implements Module interface
func (t *Text) Update(w *ggl.Window, delta float64) {
	t.Paragraph.Update(delta)
	if t.Changes() {
		t.Scene.Redraw.Notify()
	}

	start, end := t.Start, t.End
	if start > end {
		start, end = end, start
	}

	if w.Pressed(key.LeftControl) {
		if start != end && w.JustPressed(key.C) {
			t.Scene.Log(t.Element, t.Clip(start, end))
		}
	}

	if !t.Hovering && start == end {
		return
	}

	// selection start
	if w.JustPressed(key.MouseLeft) {
		if !t.Hovering {
			t.selected = false
			t.End = t.Start
		} else {
			t.Start, t.LineIdx, t.Line = t.CursorFor(w.MousePos())
			t.selected = true
		}
		t.Scene.Redraw.Notify()
	}

	// selection dragging
	if w.Pressed(key.MouseLeft) && t.selected {
		oldI, oldL, oldE := t.LineIdx, t.Line, t.End
		t.End, t.LineIdx, t.Line = t.CursorFor(w.MousePos())
		if oldE != t.End {
			t.Scene.Redraw.Notify()
		}
		if t.End > t.Start {
			t.LineIdx, t.Line = oldI, oldL
		}
	}
}

// DrawOnTop implements Module interface
func (t *Text) DrawOnTop(tg ggl.Target, canvas *drw.Geom) {
	start, end := t.Start, t.End
	if start != end {
		t.Scene.TextSelected = true
		if start > end {
			start, end = end, start
		}

		for i := start; i < end; i++ {
			min := t.Dot(i)
			min.Y -= t.Descent * t.Scl.Y

			max := t.Dot(i + 1)
			max.Y += t.Ascent * t.Scl.Y

			canvas.Color(t.SelectionColor).AABB(mat.AABB{Min: min, Max: max})
		}
	}

	canvas.Fetch(tg)
}

// OnFrameChange implements Module interface
func (t *Text) OnFrameChange() {
	t.Pos = mat.V(t.Frame.Min.X+t.Padding.Min.X, t.Frame.Max.Y-t.Padding.Max.Y)
	t.Paragraph.Update(0)
}

// Width implements Module interface
func (t *Text) Width(takable, taken float64) float64 {
	t.UpdateParagraph(takable)
	if t.Align != txt.Left {
		return t.Paragraph.Width * t.Scl.X
	}
	return t.Bounds().W() * t.Scl.X
}

// Height implements Module interface
func (t *Text) Height(takable, taken float64) float64 {
	return t.Bounds().H() * t.Scl.Y
}

// UpdateParagraph updates the paragraph to fit given width, though can end up bigger
func (t *Text) UpdateParagraph(width float64) {
	if t.dirty || width != t.Paragraph.Width*t.Scl.X {
		t.Paragraph.Width = width / t.Scl.X
		t.Markdown.Parse(&t.Paragraph)
		t.dirty = false
	}
}

// Clip copies the range into clipboard
func (t *Text) Clip(start, end int) error {
	return clipboard.WriteAll(string(t.Compiled[start:end]))
}

// SetText sets text and displays the change
func (t *Text) SetText(text string) {
	t.Content = str.NString(text)
	t.Dirty()
}

// Dirty forces text to update
func (t *Text) Dirty() {
	t.dirty = true
	t.Start, t.End = 0, 0
	t.Scene.Resize.Notify()
}
