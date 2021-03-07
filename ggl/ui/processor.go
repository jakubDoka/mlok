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

	pfTmp, pfTmp2 []*float64
	divTemp       []*Element
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
	p.scene.Batch.Clear()
	p.scene.Root.redraw(&p.scene.Batch, &p.canvas)
	p.scene.Redraw.Done()
}

// Resize has to be called upon Frame change, its not recommended
// to call this manually, instead call Deformed() to notify processor
// about change
func (p *Processor) Resize() {

	p.scene.Root.size = p.frame.Size()
	p.scene.Root.resize(p)
	p.scene.Root.move(p.frame.Min, false)

	p.scene.Resize.Done()
	p.scene.Redraw.Notify()
}

// calcSize calculates all sizes equal to Fill amongst children of d
func (p *Processor) calcSize(d *Element) {
	offset, space, privateSize := p.setup(d.Horizontal(), d.size)

	prev := privateSize
	if d.Horizontal() {
		d.forChild(IgnoreHidden, func(ch *Element) {
			m := sumMargin(1, ch)
			privateSize = math.Max(ch.Module.PrivateHeight(privateSize-m)+m, privateSize)
		})
	} else {
		d.forChild(IgnoreHidden, func(ch *Element) {
			m := sumMargin(0, ch)
			privateSize = math.Max(ch.Module.PrivateWidth(privateSize-m)+m, privateSize)
		})
	}

	// we set newly obtained private space but only if it makes sene for resize mode
	if d.Resizing[offset] >= Shrink {
		privateSize = prev
	}

	p.divTemp = p.divTemp[:0]
	d.forChild(IgnoreHidden, func(ch *Element) {
		flt := ch.Margin.Flatten()
		is := privateSize
		for i, v := range flt {
			if i%2 == offset {
				if v != Fill {
					space -= v
				}
			} else {
				if v != Fill {
					is -= v
				}
			}
		}

		sft := ch.Size.Flatten()
		smt := ch.size.Mutator()

		if sft[1-offset] == Fill {
			*smt[1-offset] = is
		} else {
			*smt[1-offset] = sft[1-offset]
		}

		param := sft[offset]
		if param == Fill {
			p.divTemp = append(p.divTemp, ch)
			p.pfTmp = append(p.pfTmp, smt[offset])
		} else {
			*smt[offset] = param
			space -= param
		}
	})

	feed(space, p.pfTmp)

	if d.Resizing[offset] < Shrink {
		if d.Horizontal() {
			for _, ch := range p.divTemp {
				ch.size.X = ch.Module.PublicWidth(ch.size.X)
			}
		} else {
			for _, ch := range p.divTemp {
				ch.size.Y = ch.Module.PublicHeight(ch.size.Y)
			}
		}
	}
}

func sumMargin(offset int, ch *Element) (r float64) {
	arr := ch.Margin.Flatten()
	for i := offset; i < len(arr); i += 2 {
		if arr[i] != Fill {
			r += arr[i]
		}
	}

	return
}

// calculates margin in case it is Equal to Fill for all
// children of div
func (p *Processor) calcMargin(d *Element) {
	/*
		goal is to calculate how match free space is in element and divide it
		between margin equal to Fill that supports it, function
		collects the pointers to all fill fields and feeds tham with supposed
		values
	*/
	offset, space, privateSize := p.setup(d.Horizontal(), d.size)

	d.forChild(IgnoreHidden, func(ch *Element) {
		p.pfTmp2 = p.pfTmp2[:0]
		mut := ch.margin.Mutator()
		flt := ch.Margin.Flatten()
		is := privateSize
		for i, v := range flt {
			// deciding how to treat margin value, notice how var offset relates to var horizontal
			if i%2 == offset {
				if v == Fill {
					p.pfTmp = append(p.pfTmp, mut[i])
				} else {
					*mut[i] = v
					space -= v
				}
			} else {
				if v == Fill {
					p.pfTmp2 = append(p.pfTmp2, mut[i])
				} else {
					*mut[i] = v
					is -= v
				}
			}

		}

		// subtracting the size of div, its like this because this way it works for
		// vertical and horizontal case
		sft := ch.size.Flatten()
		is -= sft[1-offset]
		space -= sft[offset]

		feed(is, p.pfTmp2)
	})

	feed(space, p.pfTmp)
}

func (p *Processor) setup(horizontal bool, size mat.Vec) (offset int, space, privateSize float64) {
	p.pfTmp = p.pfTmp[:0]
	offset = 1
	privateSize, space = size.XY()
	if horizontal {
		offset = 0
		space, privateSize = size.XY()
	}

	return
}

// feed performs final space division between elements
func feed(space float64, targets []*float64) {
	if space <= 0 {
		// make sure they are zero, this gets rid of old values
		for _, v := range targets {
			*v = 0
		}
	} else {
		// split equally
		perOne := space / float64(len(targets))
		for _, v := range targets {
			*v = perOne
		}
	}
}

// Scene is base of an ui scene, it stores all elements and notifiers
// for Processor to process
type Scene struct {
	Redraw, Resize Notifier
	Root           Element
	Assets         Assets
	Batch          ggl.Batch

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
		},
	}

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
