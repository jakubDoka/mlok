package ui

import (
	"gobatch/ggl"
	"gobatch/ggl/dw"
	"gobatch/mat"
	"math"
	"strings"

	"github.com/jakubDoka/goml"
	"github.com/jakubDoka/goml/goss"
	"github.com/jakubDoka/sterr"
)

/*imp(
	github.com/jakubDoka/gogen/templates
)*/

/*gen(
	templates.OrderedMap<string, *Element, Children>
)*/

var (
	errIndexOutOfBounds = sterr.New("index out of bounds (min: %e max: %e index: %e)")
	errMissingParent    = sterr.New("invalid operation, element does not have parent")
	errCannotIndex      = sterr.New("cannot index into empty div")
	errNoScene          = sterr.New("element is not part of a scene so scene parser is not avaliable")
	errNoParser         = sterr.New("parser is not avaliable in current scene, set Scene.Parser to enable this method")
)

// Element is a most basic ui element from witch all ui componenets consists of
type Element struct {
	Props

	Proc   dw.Preprocessor
	Parent *Element
	Scene  *Scene
	Module Module
	Raw    goml.Element

	children Children

	Frame, prev, margin mat.AABB
	size                mat.Vec

	hidden, initted bool
	id, group, name string
	Styles          []string
	index           int
}

// NElement initialides internal maps
func NElement() *Element {
	return &Element{
		children: NChildren(),
	}
}

// AddChild adds child to Element, child is initialized if root of e is Processor root
func (e *Element) AddChild(name string, o *Element) {
	o.onAdd(name, e.ChildCount(), e)
	e.children.Put(name, o)
}

// Path returns elements path, path rebuilding is slow and mainly mant for debbuging purposes
func (e *Element) Path() string {
	if e.Parent == nil {
		return "root"
	}
	return e.Parent.Path() + "." + e.name
}

// InsertChild allows specifying destination index were child should be inserted
//
// negative indexing is allowed, though if you want to insert child at the end
// use AddChild
//
// index can differ if you insert with name that some element already have as
// old element will get removed
//
// method panics if index is out of bounds
func (e *Element) InsertChild(name string, index int, o *Element) {
	if index == e.ChildCount() {
		e.AddChild(name, o)
		return
	}

	e.projectIndex(&index)
	o.onAdd(name, index, e)
	e.children.Insert(name, index, o)
	e.updateIndexes(index, e.ChildCount()-1)
}

// RemoveChild removes child by name and returns it of nil if there is no
// child with that name on this level
func (e *Element) RemoveChild(name string) *Element {
	div, _, _ := e.children.Remove(name)
	div.onRemove()
	return div
}

// PopChild removes child by index, negative indexing is supported
func (e *Element) PopChild(index int) *Element {
	e.projectIndex(&index)
	o := e.children.RemoveIndex(index)
	o.Value.onRemove()
	return o.Value
}

// Child takes dot separated path, by which you can get any
// child ony any level so if you have element with child "a" and that
// child has child "b" you ll get that child by "a.a"
func (e *Element) Child(path string) (*Element, bool) {
	comps := strings.Split(path, ".")
	for _, c := range comps {
		nd, _, ok := e.children.Value(c)
		if !ok {
			return nil, false
		}
		e = nd
	}

	return e, true
}

// FindChild performs recursive search for child, cap specifies how match
// children is enough, even passing 0 can result into one child in cursor
//
// passing negative value makes cursor unlimited
func (e *Element) FindChild(name string, cap int, cursor *[]*Element) bool {
	if val, _, ok := e.children.Value(name); ok {
		*cursor = append(*cursor, val)
		if cap >= 0 && len(*cursor) >= cap {
			return true
		}
	}

	s := e.children.Slice()
	for i := 0; i < len(s); i++ {
		if s[i].Value.FindChild(name, cap, cursor) {
			return true
		}
	}

	return false
}

// AddGoml parses inputted goml source and adds parsed elements to e
//
// panics if e.Scene == nil or e.Scene.Parser == nil
func (e *Element) AddGoml(source []byte) error {
	if e.Scene == nil {
		panic(errNoScene)
	}
	if e.Scene.Parser == nil {
		panic(errNoParser)
	}

	elems, err := e.Scene.Parse(source)
	if err != nil {
		return err
	}

	for _, ch := range elems {
		e.AddChild(ch.name, ch)
	}

	return nil
}

// ChildCount returns child count on first layer, useful
// with ChildAt to loop over all children
func (e *Element) ChildCount() int {
	return len(e.children.Slice())
}

// ChildAt gets child by index, usefull with ChildCount
// to loop over all children, supports negative indexing
//
// panics if index is out of bounds
func (e *Element) ChildAt(index int) *Element {
	e.projectIndex(&index)
	return e.children.Slice()[index].Value
}

// ForEachChild performs action on each child for which filter return true
// its unreasonable to use filter on just one place, just pass nil and decide in con
func (e *Element) ForEachChild(cfg FCfg, con func(ch *Element)) {
	if cfg.Filter == nil {
		cfg.Filter = func(ch *Element) bool { return true }
	}

	s := e.children.Slice()
	if cfg.Reverse {
		for i := len(s) - 1; i >= 0; i-- {
			if cfg.Filter(s[i].Value) {
				con(s[i].Value)
			}
		}
	} else {
		for i := 0; i < len(s); i++ {
			if cfg.Filter(s[i].Value) {
				con(s[i].Value)
			}
		}
	}

}

// Index getter
func (e *Element) Index() int {
	return e.index
}

// SetIndex sets index of element amongst other children
//
// negative python like indexing is allowed
//
// it does not make sense to set index if element has no parent
func (e *Element) SetIndex(value int) {
	if e.Scene != nil {
		e.Scene.Resize.Notify()
	}

	if e.Parent != nil {
		e.Parent.projectIndex(&value)
		e.Parent.children.ReIndex(e.index, value)
		e.Parent.updateIndexes(e.index, value) // this will set the index
	}
}

// ReIndex moves child from old to new index
//
// negative python like indexing is allowed
func (e *Element) ReIndex(old, new int) {
	div := e.children.Slice()[old]
	div.Value.SetIndex(new)
}

// Name getter
func (e *Element) Name() string {
	return e.name
}

// SetName will replace element with name equal to value if there is such
// div, of corse if value == e.Name() nothing happens
//
// this method will panic if e.Parent is nil
func (e *Element) SetName(value string) {
	if e.Parent != nil {
		e.Parent.children.Rename(e.name, value)
	}
	e.name = value
}

// Rename changes name of element with original name, returns false
// if no element with that name wos found, the element index is preserved
func (e *Element) Rename(old, new string) bool {
	div, _, ok := e.children.Value(old)
	if !ok {
		return false
	}
	div.SetName(new)
	return true
}

// ID ...
func (e *Element) ID() string {
	return e.id
}

// SetID ...
func (e *Element) SetID(id string) {
	if e.Scene != nil {
		delete(e.Scene.ids, e.id)
		e.Scene.ids[id] = e
	}
	e.id = id
}

// Group ...
func (e *Element) Group() string {
	return e.group
}

// SetGroup ...
func (e *Element) SetGroup(group string) {
	if e.Scene != nil {
		all, ok := e.Scene.groups[e.group]
		if ok {
			for i, v := range all {
				if v == e {
					all = append(all[:i], all[i+1:]...)
					break
				}
			}
			e.Scene.groups[e.group] = all
		}
		e.Scene.groups[group] = append(e.Scene.groups[group], e)
	}
	e.group = group
}

// Hide hides the div, when element is hidden its size is ignored
func (e *Element) Hide() {
	e.hidden = true
	e.onHiddenChange()
}

// Show does reverse of Hide
func (e *Element) Show() {
	e.hidden = false
	e.onHiddenChange()
}

func (e *Element) onHiddenChange() {
	if e.Scene != nil {
		e.Scene.Resize.Notify()
	}
}

func (e *Element) updateIndexes(old, new int) {
	if old > new {
		old, new = new, old
	}
	new++

	s := e.children.Slice()
	for i := old; i < new; i++ {
		s[i].Value.index = i
	}
}

func (e *Element) onAdd(name string, index int, parent *Element) {
	if e.Parent != nil { // just in case
		e.Parent.RemoveChild(e.name)
	}

	e.Parent = parent
	e.name = name
	e.index = index
	if parent.Scene != nil {
		e.init(parent.Scene)
	}
}

func (e *Element) onRemove() {
	e.Parent = nil
	e.Scene = nil
}

// projectIndex is used when manipulating with child indexes, it allows negative indexing
func (e *Element) projectIndex(i *int) {
	v := *i
	l := len(e.children.Slice())
	if l == 0 {
		panic(errCannotIndex)
	}

	if v < 0 {
		v += l
		if v < 0 {
			panic(errIndexOutOfBounds.Args(-l, l-1, v))
		}
		*i = v
	} else if v >= l {
		panic(errIndexOutOfBounds.Args(-l, l-1, v))
	}
}

// update propagates update call on modules
func (e *Element) update(p *Processor, w *ggl.Window, delta float64) {
	e.Module.Update(w, delta)

	e.ForEachChild(IgnoreHidden, func(ch *Element) {
		ch.update(p, w, delta)
	})
}

// Redraw draws element and all its children to target, if preprocessor is not nil
// triangles are also preprocessed
func (e *Element) redraw(t ggl.Target, canvas *dw.Geom) {
	canvas.Clear()
	canvas.Restart()

	var tar ggl.Target = e.Proc
	if tar == nil {
		tar = t
	}

	e.Module.Draw(tar, canvas)
	e.ForEachChild(IgnoreHidden, func(ch *Element) {
		ch.redraw(tar, canvas)
	})

	if e.Proc != nil {
		e.Proc.Fetch(t)
	}
}

// Resize resizes all children to fit each other, though this
// does not move them
func (e *Element) resize(p *Processor) {
	p.calcSize(e)

	e.ForEachChild(IgnoreHidden, func(ch *Element) {
		ch.resize(p)
	})

	p.calcMargin(e)

	e.evalSize(e.calcChildSize())

	e.Frame = e.size.ToAABB() // main jazz, resize the frame
}

// Init initializes element and its children
func (e *Element) init(s *Scene) {
	if e.Module == nil {
		e.Module = &ModuleBase{}
	}
	if e.children.IsNil() {
		e.children = NChildren()
	}

	e.Scene = s
	s.InitStyle(e)
	s.addElement(e)

	e.Module.Init(e)
	ch := e.children.Slice()
	for i := 0; i < len(ch); i++ {
		ch[i].Value.init(s)
	}
}

// EvalSize evaluates final size for element based of final size of children
func (e *Element) evalSize(chSize mat.Vec) {
	switch e.ResizeMode {
	case Shrink:
		e.size = chSize.Min(e.size)
	case Expand:
		e.size = chSize.Max(e.size)
	case Exact:
		e.size = chSize
	case Ignore:
		e.size = e.Props.Size
	}
}

// CalcChildSize calculates children size according to element orientation
func (e *Element) calcChildSize() (chSize mat.Vec) {
	if e.Horizontal() {
		sum := HSum{&chSize}
		e.ForEachChild(IgnoreHidden, func(ch *Element) {
			sum.Add(ch.spaceNeeded())
		})
	} else {
		sum := VSum{&chSize}
		e.ForEachChild(IgnoreHidden, func(ch *Element) {
			sum.Add(ch.spaceNeeded())
		})
	}

	return
}

// SizeNeeded returns how match space the element spams
func (e *Element) spaceNeeded() mat.Vec {
	return e.Frame.Size().Add(e.margin.Min).Add(e.margin.Max)
}

// Move is next step after resize, size of all elements is calculated,
// now we can move them all to correct place
func (e *Element) move(offset mat.Vec, horizontal bool) mat.Vec {
	off := offset.Add(e.margin.Min)
	e.Frame = e.Frame.Moved(off)
	e.Module.OnFrameChange()

	e.ForEachChild(IgnoreHiddenReverse, func(ch *Element) {
		off = ch.move(off, e.Horizontal())
	})

	if horizontal {
		l, _, r, _ := e.margin.Deco()
		offset.X += l + r + e.Frame.W()
	} else {
		_, b, _, t := e.margin.Deco()
		offset.Y += b + t + e.Frame.H()
	}

	return offset
}

// FCfg is configuration for Element.ForEachChild method
type FCfg struct {
	Filter  func(ch *Element) bool
	Reverse bool
}

var (
	hf = func(e *Element) bool {
		return !e.hidden
	}
	// IgnoreHidden filters out all hidden children
	IgnoreHidden = FCfg{
		Filter: hf,
	}
	// IgnoreHiddenReverse does the same as ignore hidden but also loops in reverse order
	IgnoreHiddenReverse = FCfg{
		Filter:  hf,
		Reverse: true,
	}
)

// Module is what makes Element alive, it defines its behavior. There is quite a but of
// functions you might not even need, so use ModuleBase as core of your module to implement default
// methods
type Module interface {
	// DefaultStyle should returns default style of element that will be used as base, returning zero value is fine
	DefaultStyle() goss.Style
	// Init is called when module is inserted into element that is already initted, assets and div
	// should cower all needs of initialization
	Init(*Element)
	// Draw should draw the div, draw your triangles onto given target, you can use Geom as canvas
	// though you have to draw it to target too, Geom is cleared and restarted before draw call
	Draw(ggl.Target, *dw.Geom)
	// Update is stage where your event handling and visual updates should happen, it gives you access to
	// Processor so you can trigger global updates and mainly call p.Dirty().
	Update(*ggl.Window, float64)
	// OnFrameChange is called by processor when frame of element changes
	OnFrameChange()
	// Size should return a size that element will take, BaseModule will just return Style size for example
	// but size can depend on state of element
	Size() mat.Vec
	// PrivateWidth should calculate the horizontal size incase of vertical parent composition that
	// element needs in case it changes dynamically this will be called by processor only if parent
	// can expand based of children
	PrivateWidth(supposed float64) (desired float64)
	// PrivateHeight is analogous to PrivateWidth but it should calculate vertical size in case of
	// horizontal composition this will be called by processor only if parent can expand based of
	// children
	PrivateHeight(supposed float64) (desired float64)
	// PublicWidth gives an option to modify horizontal size of element in case it has a fill
	// property, this will be called by processor only if parent can expand based of children
	PublicWidth(supposed float64)
	// PublicHeight does same as PublicWidth thought for Height
	PublicHeight(supposed float64)
}

// ModuleBase is a base of every module, you should embed this struct in your module
// and "override" default methods, though don't forget to call original Init that
// initializes the styles, if you don't give your element a module, this will be paced as placeholder
type ModuleBase struct {
	*Element
	Background mat.RGBA
}

// DefaultStyle implements Module interface
func (m *ModuleBase) DefaultStyle() goss.Style {
	return goss.Style{}
}

// Init implements Module interface
func (m *ModuleBase) Init(div *Element) {
	m.Element = div
	m.Background = m.RGBA("background", mat.Transparent)
}

// Draw implements Module interface
func (m *ModuleBase) Draw(t ggl.Target, g *dw.Geom) {
	g.Color(m.Background).AABB(m.Frame)
	g.Fetch(t)
}

// Update implements Module interface
func (*ModuleBase) Update(*ggl.Window, float64) {}

// OnFrameChange implements Module interface
func (*ModuleBase) OnFrameChange() {}

// Size implements Module interface
func (m *ModuleBase) Size() mat.Vec {
	return m.Props.Size
}

// PrivateWidth implements Module interface
func (*ModuleBase) PrivateWidth(supposed float64) (desired float64) { return }

// PrivateHeight implements Module interface
func (*ModuleBase) PrivateHeight(supposed float64) (desired float64) { return }

// PublicWidth implements Module interface
func (*ModuleBase) PublicWidth(supposed float64) {}

// PublicHeight implements Module interface
func (*ModuleBase) PublicHeight(supposed float64) {}

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

// EventHandler handles event registration for elements
type EventHandler map[string][]*EventListener

// Add adds listener to handler, keep the listener accessable if you want to
// remove it later
func (e EventHandler) Add(listener *EventListener) {
	evs := e[listener.Name]
	listener.idx = len(evs)
	e[listener.Name] = append(evs, listener)
}

// Invoke invokes the event listeners, removed listeners are skipped and deleted
func (e EventHandler) Invoke(name string, ed *EventData) {
	evs := e[name]
	for i := len(evs) - 1; i >= 0; i-- {
		if evs[i].Runner(ed) {
			break
		}
	}
}

// EventListener holds function tha gets called when event is triggered
// if events returns true, all consequent events will get blocked, execution
// goes from newest to oldest event listener
type EventListener struct {
	Name   string
	Runner func(*EventData) bool
	idx    int
	evs    EventHandler
}

// Remove removes the listener from event handler
func (e *EventListener) Remove() {
	if e.evs == nil {
		return
	}

	evs := e.evs[e.Name]
	for i := e.idx; i < len(evs); i++ {
		evs[i].idx--
	}

	evs = append(evs[:e.idx], evs[e.idx+1:]...)
	e.evs[e.Name] = evs
}

// EventData ...
type EventData struct{}
