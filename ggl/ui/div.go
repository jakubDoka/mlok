package ui

import (
	"fmt"
	"gobatch/ggl"
	"gobatch/ggl/dw"
	"gobatch/mat"
	"gobatch/refl"
	"math"
	"strings"
)

/*imp(
	github.com/jakubDoka/gogen/templates
)*/

/*gen(
	templates.OrderedMap<string, *Div, Children>
)*/

// Div is a most basic ui element from witch all ui componenets consists of
type Div struct {
	Module
	Style

	Proc   dw.Preprocessor
	Parent *Div

	p        *Processor
	children Children

	Frame, prev, margin mat.AABB
	size                mat.Vec

	Styles, name string
	index        int
}

// AddChild adds child to Div, child is initialized if root of d is Processor root
func (d *Div) AddChild(name string, o *Div) {
	d.children.Put(o.name, o)
	if d.p != nil {
		o.init(d.p)
	}
}

// InsertChild allows specifying index
//
// method panics if index negative or > len(children)
func (d *Div) InsertChild(name string, index int, o *Div) {
	if index < 0 || index > len(d.children.Slice()) {
		panic(fmt.Sprintf("invalid index, accepts only range 0..%d", len(d.children.Slice())))
	}
	d.children.Insert(o.name, index, o)
	if d.p != nil {
		o.init(d.p)
	}
}

// Child takes dot separated path, by which you can get any
// child ony any level so if you have div with child "a" and that
// child has child "b" you ll get that child by "a.a"
func (d *Div) Child(path string) (*Div, bool) {
	comps := strings.Split(path, ".")
	for _, c := range comps {
		nd, _, ok := d.children.Value(c)
		if !ok {
			return nil, false
		}
		d = nd
	}

	return d, true
}

// Index getter
func (d *Div) Index() int {
	return d.index
}

// SetIndex sets index of div amongst other children
//
// this of course will fail if div has no parent
func (d *Div) SetIndex(value int) {
	d.p.Dirty()
	d.Parent.children.ReIndex(d.index, value)
	d.index = value
}

// ReIndex moves child from old to new index
func (d *Div) ReIndex(old, new int) {
	div := d.children.Slice()[old]
	div.Value.SetIndex(new)
}

// Name getter
func (d *Div) Name() string {
	return d.name
}

// SetName will replace div with name equal to value if there is such
// div, of corse if value == d.Name() nothing happens
//
// this of course will fail if div has no parent
func (d *Div) SetName(value string) {
	d.Parent.children.Rename(d.name, value)
	d.name = value
}

// Rename changes name of div with original name, returns false
// if no div with that name wos found, the div index is preserved
func (d *Div) Rename(old, new string) bool {
	div, _, ok := d.children.Value(old)
	if !ok {
		return false
	}
	div.SetName(new)
	return true
}

// FindChild performs recursive search for child, cap specifies how match
// children is enough, even passing 0 can result into one child in cursor
func (d *Div) FindChild(name string, cap int, cursor *[]*Div) bool {
	if val, _, ok := d.children.Value(name); ok {
		*cursor = append(*cursor, val)
		if len(*cursor) >= cap {
			return true
		}
	}

	for _, ch := range d.children.Slice() {
		if ch.Value.FindChild(name, cap, cursor) {
			return true
		}
	}

	return false
}

// Redraw draws div and all its children to target, if preprocessor is not nil
// triangles are also preprocessed
func (d *Div) redraw(t ggl.Target, canvas *dw.Geom) {
	canvas.Clear()
	canvas.Restart()

	var tar ggl.Target = d.Proc
	if tar == nil {
		tar = t
	}

	d.Draw(tar, canvas)
	for _, ch := range d.children.Slice() {
		ch.Value.redraw(tar, canvas)
	}

	if d.Proc != nil {
		d.Proc.Fetch(t)
	}
}

// Resize resizes all children to fit each other, though this
// does not move them
func (d *Div) resize() {
	d.p.calcSize(d)

	for _, ch := range d.children.Slice() {
		ch.Value.resize()
	}

	d.p.calcMargin(d)

	chSize := d.calcChildSize()
	chSize = d.evalSize(chSize)

	d.Frame = chSize.ToAABB() // main jazz, resize the frame
}

// Init initializes div and its children
func (d *Div) init(p *Processor) {
	if d.Module == nil {
		d.Module = &ModuleBase{}
	}

	d.p = p

	d.Init(d, p.Assets)
	for _, ch := range d.children.Slice() {
		ch.Value.init(p)
	}
}

// EvalSize evaluates final size for div based of final size of children
func (d *Div) evalSize(chSize mat.Vec) mat.Vec {
	switch d.ResizeMode {
	case Shrink:
		chSize = chSize.Min(d.size)
	case Expand:
		chSize = chSize.Max(d.size)
	case Exact:
		// pass
	case Ignore:
		// calculated data turn out to be useles though update part is still important
		chSize = d.size
	}

	return chSize
}

// CalcChildSize calculates children size according to div orientation
func (d *Div) calcChildSize() (chSize mat.Vec) {
	if d.Horizontal {
		sum := HSum{&chSize}
		for _, ch := range d.children.Slice() {
			sum.Add(ch.Value.spaceNeeded())
		}
	} else {
		sum := VSum{&chSize}
		for _, ch := range d.children.Slice() {
			sum.Add(ch.Value.spaceNeeded())
		}
	}

	return
}

// SpaceNeeded returns how match space the element spams
func (d *Div) spaceNeeded() mat.Vec {
	return d.Frame.Size().Add(d.margin.Min).Add(d.margin.Max)
}

// Move is next step after resize, size of all elements is calculated,
// now we can move them all to correct place
func (d *Div) move(offset mat.Vec, horizontal bool) mat.Vec {
	off := offset.Add(d.margin.Min)
	d.Frame = d.Frame.Moved(off)
	for _, ch := range d.children.Slice() {
		off = ch.Value.move(off, d.Horizontal)
	}

	if horizontal {
		l, _, r, _ := d.margin.Deco()
		offset.X += l + r + d.Frame.W()
	} else {
		_, b, _, t := d.margin.Deco()
		offset.Y += b + t + d.Frame.H()
	}

	return offset
}

// Module is what makes Div alive, it defines its behavior, can even create more dependant
// Divs. Best practice is to use ModuleBase or any module that extends it to create new module
// Mainly because you can overwrite behavior or extend it. If you want modules to depend on each
// other embed them into ine main module and pass pointer to gem into dependant elements. This way
// you can easily connect the behavior
type Module interface {
	// Init is called when module is inserted into div that is already initted, assets and div
	// should cower all needs of initialization
	Init(*Div, *Assets)
	// Draw should draw the div, draw your triangles onto given target, you can use Geom as canvas
	// though you have to draw it to target too, Geom is cleared and restarted before draw call
	Draw(ggl.Target, *dw.Geom)
	// Update is stage where your event handling and visual updates should happen, it gives you access to
	// Processor so you can trigger global updates and mainly call p.Dirty().
	Update(*Processor, *ggl.Window, float64)
}

// ModuleBase is a base of every module, you should embed this struct in your module
// and "override" default methods, though don't forget to call original Init that
// initializes the styles, if you don't give your div a module, this will be paced as placeholder
type ModuleBase struct {
	*Div
}

// Init implements Module interface
func (m *ModuleBase) Init(div *Div, asstets *Assets) {
	m.Div = div
	for _, s := range strings.Split(m.Styles, " ") {
		style := asstets.Styles[s]
		refl.Overwrite(&m.Style, style, false)
	}
}

// Draw implements Module interface
func (m *ModuleBase) Draw(t ggl.Target, g *dw.Geom) {
	if len(m.Subs) == 0 {
		return
	}

	g.Color(m.Subs[m.Current].Background).AABB(m.Frame)
	g.Fetch(t)
}

// Update implements Module interface
func (*ModuleBase) Update(*Processor, *ggl.Window, float64) {}

// HSum calculates size of elements in horizontal composition
type HSum struct {
	*mat.Vec
}

// Add performs calculation
func (h *HSum) Add(size mat.Vec) {
	h.X += size.X
	h.Y = math.Max(h.Y, size.Y)
}

// VSum is analogous to HSum just for horizontal composition
type VSum struct {
	*mat.Vec
}

// Add performs calculation
func (h *VSum) Add(size mat.Vec) {
	h.Y += size.Y
	h.X = math.Max(h.X, size.X)
}

// Move moves child to idex among other children
func (c *Children) Move(src, dst int) {
	h := c.RemoveIndex(src)
	c.Insert(h.Key, dst, h.Value)
}
