package ui

import (
	"io/ioutil"
	"math"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/drw"
	"github.com/jakubDoka/mlok/ggl/pck"
	"github.com/jakubDoka/mlok/ggl/txt"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/sterr"

	"github.com/jakubDoka/goml/goss"
)

// Error thrown when processors scene is nil
var ErrNoScene = sterr.New("processor is missing scene to process (use p.SetScene)")

// Processor handles scene composition, it resizes all elements and draws them
type Processor struct {
	scene  *Scene
	frame  mat.AABB
	canvas drw.Geom

	margins, relativeMargins []*float64
	filled                   []*Element
	verticalFormatter        VerticalFormatter
	horizontalFormatter      HorizontalFormatter
}

// SetScene sets current scene processor uses
//
// you have to set scene before performing any other operations
func (p *Processor) SetScene(s *Scene) {
	p.scene = s
	s.Resize.Notify()
}

// Fetch passes triangles to given target
//
// panics if scene is not set
func (p *Processor) Fetch(t ggl.Target) {
	p.assertScene()

	p.scene.Batch.Fetch(t)
}

// Render make Renderer render ui (yes)
//
// panics if scene is not set
func (p *Processor) Render(r ggl.Renderer) {
	p.assertScene()

	p.scene.Batch.Draw(r)
}

// SetFrame sets the frame of Processor witch also updates all elements inside
// you can call this every frame if your window is resizable as update will not get
// triggered if value == oldFrame
//
// panics if scene is not set
func (p *Processor) SetFrame(value mat.AABB) {
	p.assertScene()

	if p.frame != value {
		p.frame = value
		p.scene.Resize.Notify()
	}
}

// Update calls update on all elements and performs resizing and redrawing if needed
// call this every frame, resizing and should not happen if user is not
// interacting with ui. Of corse you can trigger either one by:
//
// 	scene.Resize.Notify()
//	scene.Redraw.Notify()
//
// note that resizing also triggers consequent Redrawing
//
// panics if scene is not set
func (p *Processor) Update(w *ggl.Window, delta float64) {
	p.assertScene()

	p.scene.Root.update(p, w, delta)
	if p.scene.Resize.Should() {
		p.Resize()
	}
	if p.scene.Redraw.Should() {
		p.Redraw()
	}
}

// Redraw redraws the scene
//
// should not be called manually, use scene.Redraw.Notify() instead
//
// panics if scene is not set
func (p *Processor) Redraw() {
	p.assertScene()

	p.scene.TextSelected = false
	p.scene.Batch.Clear()
	p.scene.Root.redraw(&p.scene.Batch, &p.canvas)
	p.scene.Redraw.Done()
}

// Resize prefroms scene resizing
//
// should not be called manually, use scene.Resize.Notify() instead
//
// panics if scene is not set
func (p *Processor) Resize() {
	p.assertScene()

	p.resize(&p.scene.Root, p.frame.W(), &p.horizontalFormatter, X)
	p.resize(&p.scene.Root, p.frame.H(), &p.verticalFormatter, Y)

	p.scene.Root.move(p.frame.Min, false)

	p.scene.Resize.Done()
	p.scene.Redraw.Notify()
}

// assertScene panics if processor scene is nil
func (p *Processor) assertScene() {
	if p.scene == nil {
		panic(errNoScene)
	}
}

// resize is complex mathod that performs resizing of scene on given dimension, order of dimensions should
// X and then Y because of how text resizing works
func (p *Processor) resize(e *Element, takable float64, formatter Formatter, dim Dimension) (taken float64) {
	formatter.Set(e)
	takable -= formatter.Sum(e.Margin)
	ptr, size := formatter.Ptr(), formatter.Size()

	if size != Fill {
		switch e.Resizing[dim] {
		case Ignore:
			takable = size
			*ptr = takable
		case Shrink:
			takable = math.Min(size, takable)
		}
	}

	takable -= formatter.Sum(e.Padding)

	// splitter performs space splitting between len targets and
	// prevents negative sizes
	splitter := func(total float64, len int) float64 {
		return math.Max(total/float64(len), 0)
	}

	// filler calculates size that element should be offered with
	// if it is standalone in current dimension
	filler := func(e *Element) float64 {
		formatter.Set(e)
		sz := formatter.Size()
		if sz == Fill || e.Expands(dim) {
			sz = takable
		} else {
			sz += formatter.Sum(e.Margin)
		}
		return sz
	}

	s := e.children.Slice()
	if formatter.Condition(e.Horizontal()) { // resolving public sizes
		e.processed = e.processed[:0]
		// sizes that are not fill, also relative sizes
		for i := 0; i < len(s); i++ {
			ch := s[i].Value
			if ch.hidden {
				continue
			}

			formatter.Set(ch)
			sz := formatter.Size()
			if ch.Relative {
				p.resize(ch, filler(ch), formatter, dim)
			} else if sz != Fill {
				taken += p.resize(ch, sz+formatter.Sum(ch.Margin), formatter, dim)
			} else {
				e.processed = append(e.processed, ch)
			}
		}

		cont := float64(len(e.processed))
		calc := func() float64 { return math.Max((takable-taken)/(cont), 0) }

		split := calc()
		// evaluating fill sizes
		for _, ch := range e.processed {
			diff := p.resize(ch, split, formatter, dim)
			taken += diff
			cont -= 1
			if diff != split {
				split = calc()
			}
		}

		p.margins = p.margins[:0]
		// resolving margins
		for i := 0; i < len(s); i++ {
			ch := s[i].Value
			formatter.Set(ch)
			if ch.Relative {
				p.relativeMargins = p.relativeMargins[:0]
				for _, v := range formatter.MarginPtr() {
					if *v == Fill {
						p.relativeMargins = append(p.relativeMargins, v)
					}
				}

				split := splitter(takable-*formatter.Ptr(), len(p.relativeMargins))
				for _, v := range p.relativeMargins {
					*v = split
				}
			} else {
				for _, v := range formatter.MarginPtr() {
					if *v == Fill {
						p.margins = append(p.margins, v)
					}
				}
			}
		}

		// evaluating fill margins
		split = splitter(takable-taken, len(p.margins))
		for _, v := range p.margins {
			*v = split
		}
	} else { // resolving provate sizes
		e.processed = e.processed[:0]
		// resolve sizes
		for i := 0; i < len(s); i++ {
			ch := s[i].Value
			if ch.hidden {
				continue
			}
			taken = math.Max(p.resize(ch, filler(ch), formatter, dim), taken)
			e.processed = append(e.processed, ch)
		}

		if takable < taken && e.Expands(dim) {
			for _, ch := range e.processed {
				if ch.hidden {
					continue
				}
				formatter.Set(ch)
				if formatter.Size() == Fill && taken > *formatter.Ptr() {
					taken = math.Max(p.resize(ch, taken, formatter, dim), taken)
				}
			}
			takable = taken
		}

		// resolve margins
		for _, ch := range e.processed {
			formatter.Set(ch)
			p.margins = p.margins[:0]
			for _, v := range formatter.MarginPtr() {
				if *v == Fill {
					p.margins = append(p.margins, v)
				}
			}

			split := splitter(takable-*formatter.Ptr(), len(p.margins))
			for _, v := range p.margins {
				*v = split
			}
		}
	}

	formatter.Set(e)
	formatter.ChildSize(taken)
	taken = formatter.Provide(takable, taken)

	if size == Fill {
		if e.Resizing[dim] == Ignore {
			*ptr = takable
		} else {
			*ptr = math.Max(takable, taken)
		}
	} else {
		// shrink or expand
		switch e.Resizing[dim] {
		case Expand:
			*ptr = math.Max(taken, size)
		case Shrink:
			*ptr = math.Min(taken, size)
		case Exact:
			*ptr = taken
		}

	}

	*ptr += formatter.Sum(e.Padding)
	return *ptr + formatter.Sum(e.Margin)
}

// HorizontalFormatter is used when resolving Horizontal sizes and margins
// methods are not documented, for doc look for Formatter interface
type HorizontalFormatter struct {
	FormatterBase
}

func (r *HorizontalFormatter) Size() float64 { return r.Props.Size.X }

func (r *HorizontalFormatter) Ptr() *float64 { return &r.size.X }

func (r *HorizontalFormatter) Sum(a mat.AABB) float64 { return fill.hSum(a) }

func (r *HorizontalFormatter) Provide(takable, taken float64) float64 {
	return r.Module.Width(takable, taken)
}

func (r *HorizontalFormatter) MarginPtr() [2]*float64 {
	r.margin.Min.X, r.margin.Max.X = r.Margin.Min.X, r.Margin.Max.X
	return [2]*float64{&r.margin.Min.X, &r.margin.Max.X}
}

func (r *HorizontalFormatter) Condition(b bool) bool { return b }

func (r *HorizontalFormatter) ChildSize(value float64) { r.Element.ChildSize.X = value }

// VerticalFormatter is used when resolving Vertical sizes and margins
// methods are not documented, for doc look for Formatter interface
type VerticalFormatter struct {
	FormatterBase
}

func (r *VerticalFormatter) Size() float64 { return r.Props.Size.Y }

func (r *VerticalFormatter) Ptr() *float64 { return &r.size.Y }

func (r *VerticalFormatter) Sum(a mat.AABB) float64 { return fill.vSum(a) }

func (r *VerticalFormatter) Provide(takable, taken float64) float64 {
	return r.Module.Height(takable, taken)
}

func (r *VerticalFormatter) MarginPtr() [2]*float64 {
	r.margin.Min.Y, r.margin.Max.Y = r.Margin.Min.Y, r.Margin.Max.Y
	return [2]*float64{&r.margin.Min.Y, &r.margin.Max.Y}
}

func (r *VerticalFormatter) Condition(b bool) bool { return !b }

func (r *VerticalFormatter) ChildSize(value float64) { r.Element.ChildSize.Y = value }

type FormatterBase struct {
	*Element
}

func (r *FormatterBase) Set(e *Element) {
	r.Element = e
}

// Formatter handles horizontal or vertical resizing
type Formatter interface {
	// Size should returns element size stored in elements configuration
	Size() float64
	// Set stets currently processed element
	Set(*Element)
	// Ptr returns pointer to currently resized element dimension
	Ptr() *float64
	// Sum sums up the margin/padding simension
	Sum(mat.AABB) float64
	// Provide calls method on elements module with given arguments
	Provide(float64, float64) float64
	// MarginPtr returns pointers to margin dimensions
	MarginPtr() [2]*float64
	// Condition modifies boolean as needed for resizing
	Condition(bool) bool
	// ChildSize sets final child size fo element
	ChildSize(float64)
}

// Scene is a ui scene, it stores all elements and notifiers
// for Processor to process
type Scene struct {
	Redraw, Resize Notifier
	Root           Element
	Assets         *Assets
	Batch          ggl.Batch

	TextSelected bool

	*Parser

	ids    map[string]*Element
	groups map[string][]*Element
}

// NScene returns ready-to-use scene, do not use Scene{}
//
// use only after you created the window
//
// method panics if this is called before window creation
func NScene() *Scene {
	s := &Scene{
		ids:    map[string]*Element{},
		groups: map[string][]*Element{},
		Assets: &Assets{
			Styles: goss.Styles{},
			Markdowns: map[string]*txt.Markdown{
				"default": txt.NMarkdown(),
			},
			Cursors: map[string]CursorDrawer{},
		},
		Parser: NParser(),
	}

	s.SetSheet(pck.Sheet{
		Pic: txt.Atlas7x13.Pic,
	})

	s.Root.init(s)

	return s
}

// in comparison to NScene this function can me used before window creation
func NEmptyScene() *Scene {
	s := &Scene{
		ids:    map[string]*Element{},
		groups: map[string][]*Element{},
		Assets: &Assets{},
	}

	s.Root.init(s)

	return s
}

func (s *Scene) LoadGoss(paths ...string) error {
	for _, p := range paths {
		bts, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		err = s.AddGoss(bts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scene) AddGoss(source []byte) error {
	if s.Parser == nil {
		panic("parser is nil, direct adding is not avaliable")
	}

	if s.Assets == nil {
		panic("assets are missing, direct adding is not avaliable")
	}

	stl, err := s.Parser.GS.Parse(source)
	if err != nil {
		return err
	}

	s.Assets.Styles.Add(stl)
	return nil
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
func (s *Scene) SetSheet(sheet pck.Sheet) {
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
// panics if assets are nil
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
//
// panics if assets are nil
func (s *Scene) InitStyle(e *Element) {
	if s.Assets == nil {
		panic("assets are missing, style initialization is not avaliable")
	}
	e.Style = e.Module.DefaultStyle()
	if e.Style == nil {
		e.Style = goss.Style{}
	}
	for _, st := range e.Styles {
		style, ok := s.Assets.Styles[st]
		if ok {
			style.Overwrite(e.Style)
		}
	}
	if e.Raw.Style != nil {
		e.Raw.Style.Overwrite(e.Style)
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
