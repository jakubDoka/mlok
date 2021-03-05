package ui

import (
	"gobatch/ggl"
	"gobatch/ggl/dw"
	"gobatch/ggl/txt"
	"gobatch/mat"

	"github.com/jakubDoka/gogen/str"
	"github.com/jakubDoka/goml"
	"github.com/jakubDoka/goml/goss"
)

// TextFactory instantiates text modules
type TextFactory struct{}

// New implements module factory interface
func (t *TextFactory) New(elem goml.Element) (Module, bool) {
	tx := &Text{}
	if val, ok := elem.Attributes["text"]; ok {
		tx.text = val[0]
	}

	return tx, true
}

// Text handles text rendering
type Text struct {
	ModuleBase
	txt.Paragraph
	*txt.Markdown
	dw.SpriteViewport

	text   string
	Offset mat.Vec
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

	//t.Proc = &t.SpriteViewport
	t.Text = str.NString(t.text)
}

// Draw implements Module interface
func (t *Text) Draw(tr ggl.Target, g *dw.Geom) {
	g.Fetch(tr)
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
	t.Width = supposed / t.Scl.X
	t.Markdown.Parse(&t.Paragraph)
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
	t.Pos = t.Offset.Add(mat.V(t.Frame.Min.X, t.Frame.Max.Y))
	t.SpriteViewport.Area = t.Frame
	t.Paragraph.Update(0)
}

// Size implements Module interface
func (t *Text) Size() mat.Vec {
	return t.Paragraph.Bounds().Size().Mul(t.Scl).Max(t.Props.Size)
}

// SetText sets text and displays the change
func (t *Text) SetText(text string) {
	t.Paragraph.Text = str.NString(t.text)
	t.Markdown.Parse(&t.Paragraph)
}
