package ui

import (
	"gobatch/ggl"
	"gobatch/ggl/dw"
	"gobatch/ggl/key"
	"gobatch/ggl/txt"
	"gobatch/mat"
	"gobatch/mat/rgba"
	"math"

	"github.com/atotto/clipboard"
	"github.com/jakubDoka/gogen/str"
	"github.com/jakubDoka/goml/goss"
)

// Area is a text input element, you can get its content by turning Area.Text.Text into string
type Area struct {
	Text
	HoldMap

	dw CursorDrawer

	selected, dirty, noEffects bool

	LineIdx, Line int

	CursorWidth                float64
	CursorMask, SelectionColor mat.RGBA
}

// New implements ModuleFactory interface
func (a *Area) New() Module {
	return &Area{}
}

// Init implements Module interface
func (a *Area) Init(e *Element) {
	a.Composed = true // important for next call
	a.Text.Init(e)

	a.dw = a.CursorDrawer("cursor_drawer", a.Scene.Assets.Cursors, defaultCursor{})
	a.CursorWidth = e.Float("cursor_width", 2)
	a.CursorMask = e.RGBA("cursor_mask", mat.White)
	a.SelectionColor = e.RGBA("selection_color", mat.Alpha(.5))
	a.AutoFrequency = e.Float("auto_frequency", .03)
	a.HoldResponceSpeed = e.Float("hold_responce_speed", .5)

	a.binds = map[key.Key]float64{}
}

// Update implements Module interface
func (a *Area) Update(w *ggl.Window, delta float64) {
	// Text.Update sets up lot of things
	a.Text.Update(w, delta)

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
	})) || (a.Start < len(a.Text.Text) && a.Hold(key.Right, w, delta, func() {
		a.Start++
	})) || a.dirty {
		a.LineIdx, a.Line = a.UnprojectLine(a.Start)
		a.End = a.Start
		a.dirty = false
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
		typed = "\a"
	}) && typed == "" && !cut {
		return
	}

	if cut { // we don't want to accidentally write something
		typed = ""
	}

	a.Text.Text.RemoveSlice(a.Start, a.End)
	a.Text.Text.InsertSlice(a.Start, str.NString(typed))
	a.Events.Invoke(TextChanged, typed)
	a.Start += len(typed)
	a.End = a.Start
	a.Dirty()
	a.dirty = true
}

// DrawOnTop implements Module interface
func (a *Area) DrawOnTop(tg ggl.Target, canvas *dw.Geom) {
	if !a.selected {
		return
	}
	a.dw.Draw(tg, canvas, a.Dot(mat.Maxi(a.Start, a.End)), mat.V(a.CursorWidth, a.Ascent*a.Scl.Y), a.CursorMask)
	canvas.Clear()
	a.Text.DrawOnTop(tg, canvas)
}

// Dirty is similar to Text.Dirty but it preserves the Start and End
func (a *Area) Dirty() {
	start, end := a.Start, a.End
	a.Text.Dirty()
	a.Start, a.End = start, end
}

// CursorDrawer is something that draws the cursor inside the Area when area is selected
type CursorDrawer interface {
	Draw(t ggl.Target, canvas *dw.Geom, base, size mat.Vec, mask mat.RGBA)
}

type defaultCursor struct{}

func (d defaultCursor) Draw(t ggl.Target, canvas *dw.Geom, base, size mat.Vec, mask mat.RGBA) {
	size.X *= .5
	canvas.Color(mask).AABB(mat.A(base.X-size.X, base.Y, base.X+size.X, base.Y+size.Y))
	canvas.Fetch(t)
}

// HoldMap is extracted behavior of text edit, see HoldFunction
type HoldMap struct {
	HoldResponceSpeed, AutoFrequency float64

	binds map[key.Key]float64
}

// Hold supports Hold and repeat effect, if you hold button for long enough, action starts repeating
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

const buttonStateLen = 3

const (
	idle = iota
	hover
	pressed
)

var buttonStates = [buttonStateLen]string{"idle", "hover", "pressed"}

// Button is a button, has lot of space for customization, it can be in three states:
// idle, hover and pressed, each mode can have different message, texture, padding and mask
// use of texture is optional and button uses Patch for more flexibility, whe button is
// initialized, it creates its own Text element witch visible text
//
// style:
// 	all_text: 						string		// sets text on all state
//  idle/hover/pressed+_text: 		string		// sets text for each state
//	all_masks: 						color/name	// sets mask on all states
//  idle/hover/pressed+_mask: 		color/name	// sets mask for each state
//  all_regions:                	aabb/name   // sets region on all states
// 	idle/hover/pressed+_region:		aabb/name	// sets region for each state
//  all_padding:                	aabb	   	// sets padding on all states
// 	idle/hover/pressed+_padding:	aabb		// sets padding for each state
type Button struct {
	Patch

	Text   Text
	States [buttonStateLen]ButtonState

	current  int
	selected bool
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
	for i := range b.States {
		b.States[i].Text = parsed
	}
	mask := b.RGBA("all_masks", mat.White)

	for i := range b.States {
		b.States[i].Mask = mask
	}
	region := e.Region("all_regions", e.Scene.Assets.Regions, mat.ZA)
	for i := range b.States {
		b.States[i].Region = region
	}
	padding := b.AABB("all_padding", mat.ZA)
	for i := range b.States {
		b.States[i].Padding = padding
	}
	for i, s := range buttonStates {
		bs := &b.States[i]
		bs.Text = str.NString(b.Raw.Attributes.Ident(s+"_text", text))
		bs.Mask = b.RGBA(s+"_mask", bs.Mask)
		bs.Region = e.Region(s+"_region", e.Scene.Assets.Regions, bs.Region)
		bs.Padding = b.AABB(s+"_padding", bs.Padding)
	}
	textElem := NElement()
	textElem.Module = &b.Text
	b.AddChild("buttonText", textElem)
	b.current = -1
	b.ApplyState(idle)
}

// Update implements Module interface
func (b *Button) Update(w *ggl.Window, delta float64) {
	if !b.Hovering {
		b.ApplyState(idle)
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
		b.ApplyState(pressed)
	} else {
		b.ApplyState(hover)
	}
}

// ApplyState applies the button state by index
func (b *Button) ApplyState(state int) {
	if b.current == state {
		return
	}
	b.current = state
	bs := &b.States[state]
	b.Patch.Padding = bs.Padding
	b.Patch.SetRegion(bs.Region)
	b.Patch.Mask = bs.Mask
	b.Text.Text = bs.Text
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

// Patch is similar tor Sprite but uses ggl.Patch instead thus has more styling options
//
// style:
//	region: 		aabb/name 	// rendert texture region
//  padding: 		aabb     	// patch padding
//	patch_mask: 	color/name 	// path mask, if there is no region it acts as background color
//  patch_scale: 	vec 		// because Patch is more flexible, specifying scale can create more variations
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
	p.Padding = e.AABB("padding", mat.A(w, h, w, h))

	p.Scale = e.Vec("patch_scale", mat.V(1, 1))

	p.SetRegion(p.Region)
}

// Draw implements Module interface
func (p *Patch) Draw(t ggl.Target, canvas *dw.Geom) {
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
// 	mask: 	color 		// texture modulation
//  region: aabb/name 	// from where texture should be sampled from
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
func (s *Sprite) Draw(t ggl.Target, canvas *dw.Geom) {
	s.Sprite.Fetch(t)
}

// Scroll can make element visible trough scrollable viewport
//
// style:
// 	bar_width: 			float/int 		// sets bar thickness
//	friction: 			float/int  		// higher the value, less sliding, negative == no sliding
// 	scroll_sensitivity: float/int 		// how match will scroll slide
//	bar_color:          color			// bar handle color
//  rail_color:         color           // bar rail color
//	intersection_color: color			// if both bars are active, rectangle appears in a corner
//  bar_x/bar_y/bars:   bool			// makes bars visible and active
//  outside:            bool            // if there is only one bar, it will be displayed outside
type Scroll struct {
	ModuleBase
	dw.SpriteViewport

	BarWidth, Friction, ScrollSensitivity  float64
	BarColor, RailColor, IntersectionColor mat.RGBA
	X, Y                                   Bar
	Outside                                bool

	offset, vel, ratio, corner mat.Vec
	dirty, useVel, useles      bool
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
	s.Outside = s.Bool("outside", false)
	s.X.position = 1 // to prevent snap
	s.Y.position = 1
}

// DrawOnTop implements module interface
func (s *Scroll) DrawOnTop(t ggl.Target, c *dw.Geom) {
	if s.useles {
		return
	}

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
	if s.useles {
		return
	}

	if s.dirty {
		s.dirty = false
		s.move(s.Frame.Min.Sub(s.margin.Min).Sub(s.Offest), false)
		s.Scene.Redraw.Notify()
	}

	if s.useVel {
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
		s.useVel = s.vel.Len2() > .01
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

	s.update()
	s.updateOffset()
	s.SpriteViewport.Area = s.Frame
}

// Size implements Module interface
func (s *Scroll) Size(supposed mat.Vec) mat.Vec {
	if !s.Outside {
		return supposed
	}
	x, y := supposed.X < s.ChildSize.X, supposed.Y < s.ChildSize.Y
	if x && !y {
		supposed.Y += s.BarWidth
	} else if y && !x {
		supposed.X += s.BarWidth
	}
	return supposed
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

// move moves all elements by delta
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

// Bar ...
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
type Text struct {
	ModuleBase
	txt.Paragraph
	*txt.Markdown
	dirty, Composed           bool
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
		"text_background":      {"inherit"},
		"text_selection_color": {"inherit"},
	}
}

// Init implements Module interface
func (t *Text) Init(e *Element) {
	t.ModuleBase.Init(e)

	ident := t.Ident("markdown", "default")
	mkd, ok := t.Scene.Assets.Markdowns[ident]
	if !ok {
		panic(t.Path() + ": markdown with name '" + ident + "' is not present in assets")
	}

	t.Markdown = mkd
	t.Scl = t.Vec("text_scale", mat.V(1, 1))
	t.Mask = t.RGBA("text_color", mat.White)
	t.SelectionColor = t.RGBA("text_selection_color", mat.Alpha(.5))
	if !t.Composed {
		t.Props.Size = t.Vec("text_size", mat.V(Fill, Fill))
		t.Props.Margin = t.AABB("text_margin", mat.A(4, 4, 4, 4))
		t.Background = t.RGBA("text_background", t.Background)
	}
	t.NoEffects = t.Bool("no_effects", false)
	t.Text = str.NString(t.Raw.Attributes.Ident("text", string(t.Text)))
	t.Dirty()
}

// Draw implements Module interface
func (t *Text) Draw(tr ggl.Target, g *dw.Geom) {
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
		t.Start, t.LineIdx, t.Line = t.CursorFor(w.MousePos())
		t.Scene.Redraw.Notify()
	}

	// selection dragging
	if w.Pressed(key.MouseLeft) {
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
func (t *Text) DrawOnTop(tg ggl.Target, canvas *dw.Geom) {
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

// PrivateWidth implements Module interface
func (t *Text) PrivateWidth(supposed float64) (desired float64) {
	if t.Props.Size.X != Fill {
		return supposed
	}

	width := supposed / t.Scl.X

	if width != t.Width || t.dirty {
		t.Width = width
		t.Markdown.Parse(&t.Paragraph)
		t.dirty = false
	}

	t.size.X = math.Max(t.Bounds().W()*t.Scl.X, supposed)
	return t.size.X
}

// PublicHeight implements Module interface
func (t *Text) PublicHeight(supposed float64) float64 {
	return math.Max(t.Paragraph.Bounds().H()*t.Scl.Y, supposed)
}

// PublicWidth implements Module interface
func (t *Text) PublicWidth(supposed float64) float64 {
	width := t.PrivateWidth(supposed)
	return width
}

// OnFrameChange implements Module interface
func (t *Text) OnFrameChange() {
	t.Pos = mat.V(t.Frame.Min.X, t.Frame.Max.Y)
	t.Paragraph.Update(0)
}

// Size implements Module interface
func (t *Text) Size(supposed mat.Vec) mat.Vec {
	if t.Props.Size.X != Fill {
		width := supposed.X / t.Scl.X
		if width != t.Width || t.dirty {
			t.Width = width
			t.Markdown.Parse(&t.Paragraph)
			t.dirty = false
		}
	}

	return supposed.Max(t.Bounds().Size().Mul(t.Scl))
}

// Clip copies the range into clipboard
func (t *Text) Clip(start, end int) error {
	return clipboard.WriteAll(string(t.Compiled[start:end]))
}

// SetText sets text and displays the change
func (t *Text) SetText(text string) {
	t.Paragraph.Text = str.NString(text)
	t.Dirty()
}

// Dirty forces text to update
func (t *Text) Dirty() {
	t.dirty = true
	t.Start, t.End = 0, 0
	t.Scene.Resize.Notify()
}
