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

// SpriteFactory instantiates Sprite module
type SpriteFactory struct{}

// New implements ModuleFactory interface
func (t *SpriteFactory) New() Module {
	return &Sprite{}
}

// Sprite is sprite for ui elements
type Sprite struct {
	ModuleBase
	Sprite ggl.Sprite

	Mask mat.RGBA
}

// Init implements module interface
func (s *Sprite) Init(e *Element) {
	s.Init(e)
	s.Mask = e.RGBA("sprite_mask", mat.White)
	reg, ok := e.Scene.Assets.Regions[e.Ident("region", "")]
	if !ok {
		reg = e.AABB("region", mat.ZA)
	}
	s.Sprite.Set(mat.ZA, reg)
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
type Scroll struct {
	ModuleBase
	dw.SpriteViewport

	BarWidth, Friction, ScrollSensitivity  float64
	BarColor, RailColor, IntersectionColor mat.RGBA
	Bars                                   [2]Bar

	offset, vel, ratio, corner mat.Vec
}

// Init implements module interface
func (s *Scroll) Init(e *Element) {
	s.ModuleBase.Init(e)
	s.Proc = &s.SpriteViewport
	s.BarWidth = s.Float("bar_width", 20)
	s.Friction = s.Float("friction", -1) // instant
	s.ScrollSensitivity = s.Float("scroll_sensitivity", 30)
	s.BarColor = s.RGBA("bar_color", rgba.White)
	s.RailColor = s.RGBA("bar_rail_color", mat.RGBA{})
	s.IntersectionColor = s.RGBA("bar_intersection_color", mat.RGBA{})
	a, b := s.prt()
	a.Use = s.Bool("bar_x", false)
	b.Use = s.Bool("bar_y", false)
	a.position = 1 // to prevent snap
	b.position = 1
}

// DrawOnTop implements module interface
func (s *Scroll) DrawOnTop(t ggl.Target, c *dw.Geom) {
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
		c.Color(s.BarColor).AABB(rect)
	}

	c.Fetch(t)
}

// Update implements module interface
func (s *Scroll) Update(w *ggl.Window, delta float64) {
	a, b := s.prt()
	if !a.Use && !b.Use && s.vel.Len2() > .01 {
		s.offset.AddE(s.vel)
		if s.Friction < 0 {
			s.vel = mat.ZV
		} else {
			// we don't want to get crazy if frames are low
			s.vel.SubE(s.vel.Scaled(math.Min(s.Friction*delta, 1)))
		}
		s.Scene.Resize.Notify()
	}

	if !s.Hovering {
		return
	}

	scroll := w.MouseScroll()
	if scroll.Y != 0 {
		s.vel.Y -= scroll.Y * s.ScrollSensitivity
	}

	if w.JustReleased(key.MouseButtonLeft) {
		for i := range s.Bars {
			s.Bars[i].selected = false
		}
		return
	}

	if !w.Pressed(key.MouseButtonLeft) {
		return
	}

	var (
		mouse  = w.MousePrevPos()
		move   = mouse.To(w.MousePos())
		as, ae = s.barBounds(0)
		bs, be = s.barBounds(1)
	)

	/*if move.Len2() == 0 {
		s.vel = mat.ZV
		return
	}*/

	if a.Use && a.use && a.selected || (s.corner.Y > mouse.Y && mouse.X >= as && mouse.X <= ae) {
		if w.JustPressed(key.MouseButtonLeft) {
			a.selected = true
		}
		if a.selected {
			a.Move(-move.X)
			s.Scene.Resize.Notify()
		}
	} else {
		s.vel.X = move.X
	}

	if b.Use && b.use && b.selected || (s.corner.X < mouse.X && mouse.Y >= bs && mouse.Y <= be) {
		if w.JustPressed(key.MouseButtonLeft) {
			b.selected = true
		}
		if b.selected {
			b.Move(move.Y)
			s.Scene.Resize.Notify()
		}
	} else {
		s.vel.Y = move.Y
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
		mut  = s.ratio.Mutator()
		size = s.Frame.Size().Flatten()
		ch   = s.ChildSize.Flatten()
		both = true
	)

	for i, v := range mut {
		b := &s.Bars[i]
		b.space = size[i]
		if b.Use {
			fn := *v + s.BarWidth
			b.use = fn > 0
			if b.use {
				b.space -= s.BarWidth
				*v += s.BarWidth
			} else {
				both = false
			}
			b.length = math.Max(b.space*b.space/ch[i], 5)
			b.reminder = b.space - b.length
		} else {
			both = false
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
	if both {
		s.offset.Y += s.BarWidth
	}
	s.updateOffset()
	s.SpriteViewport.Area = s.Frame
}

// move applies velocity to offset
func (s *Scroll) update() {

	a, b := s.prt()
	dif := s.ratio.Inv()
	if dif.X < 0 {
		if a.Use {
			s.offset.X = dif.X * (1 - a.position)
		} else {
			s.offset.X = mat.Clamp(s.offset.X, dif.X, 0)
		}
	} else {
		s.offset.X = 0
	}

	if dif.Y < 0 {
		if b.Use {
			s.offset.Y = dif.Y * b.position
		} else {
			s.offset.Y = mat.Clamp(s.offset.Y, dif.Y, 0)
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
	for i := 0; i < len(ch); i++ {
		ch[i].Value.Offest = s.offset
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

// Shift shifts the bar assuming diff is projected
func (b *Bar) Shift(diff float64) {
	b.position = mat.Clamp(b.position+diff, 0, 1)
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

	t.Text = str.NString(t.Raw.Attributes.Ident("text", ""))
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
	width := supposed / t.Scl.X

	if width != t.Width {
		t.Width = width
		t.Markdown.Parse(&t.Paragraph)
	}

	t.size.X = t.Bounds().W() * t.Scl.X
	return t.size.X
}

// PublicHeight implements Module interface
func (t *Text) PublicHeight(supposed float64) {
	t.size.Y = t.Paragraph.Bounds().H() * t.Scl.Y
}

// PublicWidth implements Module interface
func (t *Text) PublicWidth(supposed float64) {
	t.PrivateWidth(supposed)
	t.PublicHeight(0)
}

// OnFrameChange implements Module interface
func (t *Text) OnFrameChange() {
	t.Pos = mat.V(t.Frame.Min.X, t.Frame.Max.Y)
	t.Paragraph.Update(0)
}

// Size implements Module interface
func (t *Text) Size() mat.Vec {
	return t.Paragraph.Bounds().Size().Mul(t.Scl).Max(t.Props.Size)
}

// SetText sets text and displays the change
func (t *Text) SetText(text string) {
	t.Paragraph.Text = str.NString(text)
	t.Scene.Resize.Notify()
}
