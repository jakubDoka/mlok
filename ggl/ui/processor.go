package ui

import (
	"gobatch/ggl"
	"gobatch/ggl/dw"
	"gobatch/ggl/pck"
	"gobatch/ggl/txt"
	"gobatch/mat"
	"math"

	"github.com/jakubDoka/goml/goss"
)

// Processor is wrapper that handles a Element, it stores some data global to all elements
// that can be reused for reduction of allocations
type Processor struct {
	scene  *Scene
	frame  mat.AABB
	canvas dw.Geom

	margins           []*float64
	filled, processed []*Element
}

// NProcessor create processor with blanc scene, so it can be used right away
func NProcessor() *Processor {
	return &Processor{scene: NScene()}
}

// SetScene ...
func (p *Processor) SetScene(s *Scene) {
	p.scene = s
}

// Fetch implements ggl.Fetcher interface
func (p *Processor) Fetch(t ggl.Target) {
	p.scene.Batch.Fetch(t)
}

// Render make Renderer render ui (yes)
func (p *Processor) Render(r ggl.Renderer) {
	p.scene.Batch.Draw(r)
}

// SetFrame sets the frame of Processor witch also updates all elements inside
func (p *Processor) SetFrame(value mat.AABB) {
	if p.frame != value {
		p.frame = value
		p.scene.Resize.Notify()
	}
}

// Update calls update on all elements, and performs resizing and redrawing if needed
func (p *Processor) Update(w *ggl.Window, delta float64) {
	p.scene.Root.update(p, w, delta)
	if p.scene.Resize.Should() {
		p.Resize()
	}
	if p.scene.Redraw.Should() {
		p.Redraw()
	}
}

// Redraw redraws ewerithing if needed
func (p *Processor) Redraw() {
	p.scene.TextSelected = false
	p.scene.Batch.Clear()
	p.scene.Root.redraw(&p.scene.Batch, &p.canvas)
	p.scene.Redraw.Done()
}

// Resize has to be called upon Frame change, its not recommended
// to call this manually, instead call Deformed() to notify processor
// about change
func (p *Processor) Resize() {

	p.scene.Root.size = p.frame.Size()
	p.scene.Root.calcMinSize()
	p.scene.Root.resize(p)
	p.scene.Root.move(p.frame.Min, false)

	p.scene.Resize.Done()
	p.scene.Redraw.Notify()
}

func (p *Processor) calcMargin(d *Element, remain mat.Vec) {
	var (
		fm formatter
		// don want to allocate if we don't have to
		hfm hFormatter
		vfm vFormatter
	)

	if d.Horizontal() {
		fm = &hfm
	} else {
		fm = &vfm
	}

	p.margins = p.margins[:0]

	s := d.children.Slice()
	for i := 0; i < len(s); i++ {
		ch := s[i].Value
		if ch.hidden {
			continue
		}

		fm.set(ch)
		ch.margin = ch.Margin

		a, b := fm.marginPtrX()
		if *a == Fill {
			p.margins = append(p.margins, a)
		}
		if *b == Fill {
			p.margins = append(p.margins, b)
		}

		sz := remain.Y - fm.y()
		a, b = fm.marginPtrY()
		if *a == Fill && *b == Fill {
			sz *= .5
			*a = sz
			*b = sz
		} else if *b == Fill {
			sz -= *a
			*b = sz
		} else if *a == Fill {
			sz -= *b
			*a = sz
		}
	}

	if len(p.margins) == 0 {
		return
	}

	perChild := remain.X / float64(len(p.margins))
	for _, v := range p.margins {
		*v = perChild
	}
}

func (p *Processor) calcSize(d *Element) (remain mat.Vec) {
	var (
		fm formatter
		// don want to allocate if we don't have to
		hfm hFormatter
		vfm vFormatter
	)

	if d.Horizontal() {
		fm = &hfm
	} else {
		fm = &vfm
	}

	size := fm.space(d.size)
	p.filled = p.filled[:0]
	p.processed = p.processed[:0]

	s := d.children.Slice()
	for i := 0; i < len(s); i++ {
		ch := s[i].Value
		if ch.hidden {
			continue
		}

		fm.set(ch)
		p.processed = append(p.processed, s[i].Value)

		c := fm.constantX(size.Y - fm.marginY())
		size.X -= fm.marginX()
		if c == Fill {
			p.filled = append(p.filled, ch)
		} else {
			size.X -= c
			fm.setX(c)
		}
	}

	ln := float64(len(p.filled))
	if ln != 0 {
		perChild := size.X / ln

		for _, ch := range p.filled {
			ln--
			fm.set(ch)
			val := fm.offer(perChild)
			fm.setX(val)
			size.X = math.Max(size.X-val, 0)
			if val != perChild {
				perChild = size.X / ln
			}
		}
		p.filled = p.filled[:0]
	}

	for _, ch := range p.processed {
		fm.set(ch)
		c := fm.constantY()
		m := fm.marginY()
		if c == Fill {

			val := fm.private(size.Y - m)
			size.Y = math.Max(val+m, size.Y)
			p.filled = append(p.filled, ch)
		} else {
			size.Y = math.Max(c+m, size.Y)
			fm.setY(c)
		}
	}

	for _, ch := range p.filled {
		fm.set(ch)
		val := size.Y - fm.marginY()
		fm.final(val)
		fm.setY(val)

	}

	return size
}

type vFormatter struct {
	e *Element
}

func (v *vFormatter) set(e *Element)              { v.e = e }
func (v *vFormatter) space(value mat.Vec) mat.Vec { return value.Swapped() }

func (v *vFormatter) constantX(y float64) float64 {
	return math.Max(v.e.Module.Height(y), v.e.ChildSize.Y)
}

func (v *vFormatter) constantY() float64 {
	return math.Max(v.e.Module.Width(-1), v.e.ChildSize.X)
}

func (v *vFormatter) marginX() (spc float64)        { return marginY(v.e) }
func (v *vFormatter) marginY() (spc float64)        { return marginX(v.e) }
func (v *vFormatter) marginPtrX() (l, r *float64)   { return marginPtrY(v.e) }
func (v *vFormatter) marginPtrY() (b, t *float64)   { return marginPtrX(v.e) }
func (v *vFormatter) setX(value float64)            { v.e.size.Y = value }
func (v *vFormatter) setY(value float64)            { v.e.size.X = value }
func (v *vFormatter) offer(value float64) float64   { return v.e.Module.OfferHeight(value) }
func (v *vFormatter) private(value float64) float64 { return v.e.Module.PrivateWidth(value) }
func (v *vFormatter) final(value float64)           { v.e.Module.FinalWidth(value) }
func (v *vFormatter) y() float64                    { return v.e.size.X }

type hFormatter struct {
	e *Element
}

func (h *hFormatter) set(e *Element)              { h.e = e }
func (h *hFormatter) space(value mat.Vec) mat.Vec { return value }

func (h *hFormatter) constantX(y float64) float64 {
	return math.Max(h.e.Module.Width(y), h.e.ChildSize.X)
}

func (h *hFormatter) constantY() float64 {
	return math.Max(h.e.Module.Height(-1), h.e.ChildSize.Y)
}

func (h *hFormatter) marginX() (spc float64)        { return marginX(h.e) }
func (h *hFormatter) marginY() (spc float64)        { return marginY(h.e) }
func (h *hFormatter) marginPtrX() (l, r *float64)   { return marginPtrX(h.e) }
func (h *hFormatter) marginPtrY() (b, t *float64)   { return marginPtrY(h.e) }
func (h *hFormatter) setX(value float64)            { h.e.size.X = value }
func (h *hFormatter) setY(value float64)            { h.e.size.Y = value }
func (h *hFormatter) offer(value float64) float64   { return h.e.Module.OfferWidth(value) }
func (h *hFormatter) private(value float64) float64 { return h.e.Module.PrivateHeight(value) }
func (h *hFormatter) final(value float64)           { h.e.Module.FinalHeight(value) }
func (h *hFormatter) y() float64                    { return h.e.size.Y }

type formatter interface {
	set(e *Element)
	space(mat.Vec) mat.Vec
	constantX(y float64) float64
	constantY() float64
	marginX() float64
	marginY() float64
	marginPtrX() (l, r *float64)
	marginPtrY() (b, t *float64)
	setX(float64)
	setY(float64)
	offer(float64) float64
	private(float64) float64
	final(float64)
	y() float64
}

func marginX(e *Element) (spc float64) {
	l, _, r, _ := e.Margin.Deco()
	if l != Fill {
		spc += l
	}
	if r != Fill {
		spc += r
	}
	return
}

func marginY(e *Element) (spc float64) {
	_, b, _, t := e.Margin.Deco()
	if b != Fill {
		spc += b
	}
	if t != Fill {
		spc += t
	}
	return
}

func marginPtrX(e *Element) (a, b *float64) {
	return &e.margin.Min.X, &e.margin.Max.X
}

func marginPtrY(e *Element) (a, b *float64) {
	return &e.margin.Min.Y, &e.margin.Max.Y
}

// Scene is base of an ui scene, it stores all elements and notifiers
// for Processor to process
type Scene struct {
	Redraw, Resize Notifier
	Root           Element
	Assets         Assets
	Batch          ggl.Batch

	TextSelected bool

	*Parser

	ids    map[string]*Element
	groups map[string][]*Element
}

// NScene returns ready-to-use scene, do not use Scene{}
func NScene() *Scene {
	s := &Scene{
		ids:    map[string]*Element{},
		groups: map[string][]*Element{},
		Assets: Assets{
			Styles: goss.Styles{},
			Markdowns: map[string]*txt.Markdown{
				"default": txt.NMarkdown(),
			},
			Cursors: map[string]CursorDrawer{},
		},
		Parser: NParser(),
	}

	s.SetSheet(&pck.Sheet{
		Pic: txt.Atlas7x13.Pic,
	})

	s.Root.init(s)

	return s
}

func (s *Scene) addElement(e *Element) {
	if e.id != "" {
		s.ids[e.id] = e
	}

	if e.group != "" {
		s.groups[e.group] = append(s.groups[e.group], e)
	}
}

// ID returns element or null if no element is under id
func (s *Scene) ID(id string) *Element {
	return s.ids[id]
}

// Group returns all elements in given group
func (s *Scene) Group(group string) []*Element {
	return s.groups[group]
}

// SetSheet sets the sprite sheet Scene will use
func (s *Scene) SetSheet(sheet *pck.Sheet) {
	s.Assets.Sheet = sheet
	s.Batch.Texture = ggl.NTexture(sheet.Pic, false)
}

// Log Invokes Error event on scene root and passes ErrorEventData as event argument
func (s *Scene) Log(e *Element, err error) {
	if err == nil {
		return
	}
	s.Root.Events.Invoke(Error, ErrorEventData{e, err})
}

// ReloadStyle reloads style of an element and all its children f.e.
//
//	s.ReloadStyle(&s.Root) // reload ewerithing
//
func (s *Scene) ReloadStyle(e *Element) {
	s.InitStyle(e)
	e.Module.Init(e)
	ch := e.children.Slice()
	for i := 0; i < len(ch); i++ {
		s.ReloadStyle(ch[i].Value)
	}
	s.Resize.Notify()
}

// InitStyle initializes the stile on given element, this can be called
// multiple times if style changes
func (s *Scene) InitStyle(e *Element) {
	e.Style = e.Module.DefaultStyle()
	if e.Style == nil {
		e.Style = goss.Style{}
	}
	if e.Raw.Style != nil {
		e.Raw.Style.Overwrite(e.Style)
	}
	for _, st := range e.Styles {
		style, ok := s.Assets.Styles[st]
		if ok {
			style.Overwrite(e.Style)
		}
	}
	e.Init()
	if e.Parent != nil {
		e.Inherit(e.Parent.Style)
	}
}

// Notifier ...
type Notifier bool

// Notify initiates notification
func (n *Notifier) Notify() {
	*n = false
}

// Done turns of the notification
func (n *Notifier) Done() {
	*n = true
}

// Should returns whether there is notification
func (n Notifier) Should() bool {
	return !bool(n)
}

// ErrorEventData ...
type ErrorEventData struct {
	Element *Element
	Err     error
}
