package ui

import (
	"gobatch/ggl"
	"gobatch/ggl/dw"
	"gobatch/ggl/key"
	"gobatch/ggl/txt"
	"gobatch/mat"
	"gobatch/mat/rgba"
	"math"

	"github.com/jakubDoka/gogen/str"
	"github.com/jakubDoka/goml/goss"
)

// ButtonFactory instantiates Button module
type ButtonFactory struct{}

// New implements ModuleFactory interface
func (t *ButtonFactory) New() Module {
	return &Button{}
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
	b.Patch.Patch.SetColor(bs.Mask)
	b.Text.Text = bs.Text
	b.Scene.Resize.Notify()
}

// ButtonState ...
type ButtonState struct {
	Mask            mat.RGBA
	Region, Padding mat.AABB
	Text            str.String
}

// PatchFactory instantiates Patch module
type PatchFactory struct{}

// New implements ModuleFactory interface
func (t *PatchFactory) New() Module {
	return &Patch{}
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

// SpriteFactory instantiates Sprite module
type SpriteFactory struct{}

// New implements ModuleFactory interface
func (t *SpriteFactory) New() Module {
	return &Sprite{}
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

// ScrollFactory instantiates scroll modules
type ScrollFactory struct{}

// New implements ModuleFactory interface
func (s *ScrollFactory) New() Module {
	return &Scroll{}
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
	Bars                                   [2]Bar
	Outside                                bool

	offset, vel, ratio, corner mat.Vec
	dirty, useVel, useles      bool
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
	a, b := s.prt()
	a.Use = s.Bool("bar_x", false)
	b.Use = s.Bool("bar_y", false)
	if !b.Use && !a.Use {
		c := s.Bool("bars", false)
		a.Use = c
		b.Use = c
	}
	s.Outside = s.Bool("outside", false)
	a.position = 1 // to prevent snap
	b.position = 1
}

// DrawOnTop implements module interface
func (s *Scroll) DrawOnTop(t ggl.Target, c *dw.Geom) {
	if s.useles {
		return
	}

	a, b := s.prt()
	if a.use {
		rect := mat.AABB{Min: s.Frame.Min, Max: s.corner}
		c.Color(s.RailColor).AABB(rect)
		rect.Min.X, rect.Max.X = s.barBounds(0)
		c.Color(s.BarColor).AABB(rect)
	}
	if b.use {
		rect := mat.AABB{Min: s.corner, Max: s.Frame.Max}
		c.Color(s.RailColor).AABB(rect)
		rect.Min.Y, rect.Max.Y = s.barBounds(1)
		c.Color(s.BarColor).AABB(rect)
	}

	if a.use && b.use {
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

	a, b := s.prt()
	if s.useVel {
		if !a.selected && !b.selected {
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
		s.Bars[0].selected = false
		s.Bars[1].selected = false
		return
	}

	if !w.Pressed(key.MouseLeft) {
		return
	}

	var (
		mouse  = w.MousePrevPos()
		move   = mouse.To(w.MousePos())
		as, ae = s.barBounds(0)
		bs, be = s.barBounds(1)
	)

	if a.Use && a.use && a.selected || (s.corner.Y > mouse.Y && mouse.X >= as && mouse.X <= ae) {
		if w.JustPressed(key.MouseLeft) {
			a.selected = true
		}
		if a.selected {
			a.Move(-move.X)
			s.dirty = true
		}
	} else {
		s.vel.X = move.X
		s.useVel = true
	}

	if b.Use && b.use && b.selected || (s.corner.X < mouse.X && mouse.Y >= bs && mouse.Y <= be) {
		if w.JustPressed(key.MouseLeft) {
			b.selected = true
		}
		if b.selected {
			b.Move(move.Y)
			s.dirty = true
		}
	} else {
		s.vel.Y = move.Y
		s.useVel = true
	}

}

func (s *Scroll) barBounds(side int) (float64, float64) {
	b := &s.Bars[side]
	prj := b.reminder * b.position
	if side == 0 {
		prj = -prj - b.length
	}
	prj += s.corner.Flatten()[side]

	return prj, prj + b.length
}

// OnFrameChange implements module interface
func (s *Scroll) OnFrameChange() {
	s.ratio = s.ChildSize.Sub(s.Frame.Size())
	s.corner = mat.V(s.Frame.Max.X, s.Frame.Min.Y)

	var (
		mut   = s.ratio.Mutator()
		final = s.ratio.Add(mat.V(s.BarWidth, s.BarWidth)).Flatten()
		size  = s.Frame.Size().Flatten()
		ch    = s.ChildSize.Flatten()
	)

	for i, v := range mut {
		oi := (i + 1) % 2
		b := &s.Bars[i]
		b.space = size[i]
		if b.Use {
			b.use = final[i] > 0
			if b.use {
				if final[oi] > 0 && s.Bars[oi].Use {
					b.space -= s.BarWidth
					*v += s.BarWidth
				}
			}
			b.length = math.Max(b.space*b.space/ch[i], 5)
			b.reminder = b.space - b.length
		}
	}

	a, b := s.prt()
	if a.use {
		s.corner.Y += s.BarWidth
	}
	if b.use {
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
	a, b := s.prt()
	dif := s.ratio.Inv()
	if dif.X < 0 {
		if a.Use && (!s.useVel || a.selected) {
			s.offset.X = dif.X * (1 - a.position) // needs to be inverted or it ll look unnatural
		} else {

			s.offset.X = mat.Clamp(s.offset.X, dif.X, 0)
			a.position = 1 - s.offset.X/dif.X // make sure to move bar too
		}
	} else {
		s.offset.X = 0
	}

	if dif.Y < 0 {
		if b.Use && (!s.useVel || b.selected) {
			s.offset.Y = dif.Y * b.position
		} else {
			s.offset.Y = mat.Clamp(s.offset.Y, dif.Y, 0)
			b.position = s.offset.Y / dif.Y
		}
	} else {
		s.offset.Y = 0
	}
}

func (s *Scroll) prt() (a, b *Bar) {
	return &s.Bars[0], &s.Bars[1]
}

// move moves all elements by delta
func (s *Scroll) updateOffset() {
	ch := s.children.Slice()
	off := s.offset
	if s.Bars[0].use && s.Bars[1].use { // have to shift whole thing because of strange dimensions
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

// Move moves the bar
func (b *Bar) Move(vel float64) {
	b.position = mat.Clamp(b.position+vel/b.reminder, 0, 1)
}

// TextFactory instantiates text modules
type TextFactory struct{}

// New implements module factory interface
func (t *TextFactory) New() Module {
	return &Text{}
}

// Text handles text rendering
type Text struct {
	ModuleBase
	txt.Paragraph
	*txt.Markdown
}

// DefaultStyle implements Module interface
func (t *Text) DefaultStyle() goss.Style {
	return goss.Style{
		"text_scale":  {"inherit"},
		"text_color":  {"inherit"},
		"text_size":   {"inherit"},
		"text_margin": {"inherit"},
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
	t.Props.Size = t.Vec("text_size", mat.V(Fill, Fill))
	t.Props.Margin = t.AABB("text_margin", mat.A(4, 4, 4, 4))
	t.Text = str.NString(t.Raw.Attributes.Ident("text", string(t.Text)))
}

// Draw implements Module interface
func (t *Text) Draw(tr ggl.Target, g *dw.Geom) {
	t.Paragraph.Draw(tr)
}

// Update implements Module interface
func (t *Text) Update(w *ggl.Window, delta float64) {
	t.Paragraph.Update(delta)
	if t.Changes() {
		t.Scene.Redraw.Notify()
	}
}

// PrivateWidth implements Module interface
func (t *Text) PrivateWidth(supposed float64) (desired float64) {
	if t.Props.Size.X != Fill {
		return
	}

	width := supposed / t.Scl.X

	if width != t.Width {
		t.Width = width
		t.Markdown.Parse(&t.Paragraph)
	}

	t.size.X = t.Bounds().W() * t.Scl.X
	return t.size.X
}

// PublicHeight implements Module interface
func (t *Text) PublicHeight(supposed float64) float64 {
	return t.Paragraph.Bounds().H() * t.Scl.Y
}

// PublicWidth implements Module interface
func (t *Text) PublicWidth(supposed float64) float64 {
	width := t.PrivateWidth(supposed)
	t.PublicHeight(0)
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
		t.Width = supposed.X / t.Scl.X
		t.Markdown.Parse(&t.Paragraph)
	}

	return supposed.Max(t.Bounds().Size().Mul(t.Scl))
}

// SetText sets text and displays the change
func (t *Text) SetText(text string) {
	t.Paragraph.Text = str.NString(text)
	t.Scene.Resize.Notify()
}
