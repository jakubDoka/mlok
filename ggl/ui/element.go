package ui

import (
	"io/ioutil"
	"math"
	"strings"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/drw"
	"github.com/jakubDoka/mlok/logic/event"
	"github.com/jakubDoka/mlok/mat"

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

// Element is bas and only building piece of ui. It has a lot of components. Most significant is Module
// witch defines appearance and behavior. Then there are Events witch serves a role of comunication with
// game by connecting listeners. Then there are attributes and style, witch is provided by goml and goss.
// Also if it is not obvious, elements create recursive tree that is encapsulated in scene, and scene is
// handled by processor, this allows switching between scenes every element queries for folloving style
// propertyes:
//
// style:
// 	margin: 			aabb						// space that element should have around it self, can be fill*
//  size: 				vec							// size element should spam, can be fill*
//	composition: 		vertical/horizontal			// composition of children (on top of each other or next to each other)
//  resize_mode/_x/_y: 	expand/shrink/exact/ignore	// how element will react to size of its children
// attributes:
//	- name - every element has to have unique identifier among other children, think of it as a name in file system
//  you can access any child Element.Child method, if there is child named "car" and it has child "wheel" you can access
//  it from cars parent by "car.wheel", if you don't specify name, string index is assigned
//  - id - id has to be unique for all elements, if there are two or more elements with same id, last one loaded
// 	keep the id, you can then quickly access the node by using Scene.ID method which is only reason why ids exist
//  - group - similar to id but multiple elements can be in one group, if you use Scene.Group slice of elements gets
//  returned
//  - style - in style you can specify styling of element with goss syntax. For example "margin: fill;size: 100;" will
// 	is equivalent to setting properties of element manually like:
//		e.Margin = mat.A(Fill, Fill, Fill, FIll)
//		e.Size = mat.V(100, 100)
//  - styles - there you can specify scene loaded styles that you want to include in element, this way you can write jus
//  style name on multiple places and not hardcode attribute for each element with same styling. List or space separated
//  string is accepted as multiple styles can be used ("style1 style2" or  ["stile1" "style2"])
//  - hidden - hidden make element hidden from the start, you can write just hidden with no value, and it will be
//  considered true
//
// *fill = reminding space inside parent will be taken, if there is more children with fill prop, space is split equally
//
// style behavior works very match css, if you specify list of stiles they will be merged together, each overriding previous
// in list. The hierarchy is default_style < style_attribute < styles_attribute < inheritance. Yes there is also inheritance.
// if you provide something like "margin: inherit;" margin will be copied from parents style if he has a margin. As order of
// children matters, ordered map is used for storing the children. You can index by Element.ChildAt and negative python
// indexing is supported.
type Element struct {
	Props
	InputState

	Proc   drw.Preprocessor
	Parent *Element
	Scene  *Scene
	Module Module
	Raw    goml.Element
	Events event.String

	Frame     mat.AABB
	ChildSize mat.Vec

	Styles              []string
	relative, processed []*Element

	margin mat.AABB
	size   mat.Vec

	children        Children
	hidden          bool
	id, group, name string
	index           int
}

// NElement constructor initialides internal maps
func NElement() *Element {
	return &Element{
		children: NChildren(),
		Events:   event.String{},
		Raw:      goml.NDiv(),
	}
}

// AddChild adds child to Element, child is initialized if root of e is scene root
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
// index can differ if you insert with name that some element already have as
// old element will get removed
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

// PopChild removes child by index
func (e *Element) PopChild(index int) *Element {
	e.projectIndex(&index)
	o := e.children.RemoveIndex(index)
	o.Value.onRemove()
	return o.Value
}

// Child takes dot separated path, by which you can get any
// child ony any level so if you have element with child "a" and that
// child has child "b" you ll get that child by "b" by "a.b"
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
// children is enough
//
// passing negative value makes cursor unlimited
func (e *Element) FindChild(name string, cap int, cursor *[]*Element) bool {
	if cap == 0 {
		return false
	}

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

// LoadGoml abstracts loading of file
func (e *Element) LoadGoml(paths ...string) error {
	for _, p := range paths {
		bts, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		err = e.AddGoml(bts)
		if err != nil {
			return err
		}
	}

	return nil
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

func (e *Element) Listen(name string, runner event.StringRunner) {
	e.Events.Add(&event.Listener{
		Name:   name,
		Runner: runner,
	})
}

// ChildCount returns child count on first layer, useful
// with ChildAt to loop over all children
func (e *Element) ChildCount() int {
	return len(e.children.Slice())
}

// ChildAt gets child by index, usefull with ChildCount
// to loop over all children
func (e *Element) ChildAt(index int) *Element {
	e.projectIndex(&index)
	return e.children.Slice()[index].Value
}

// ForChild loops over children
func (e *Element) ForChild(con func(ch *Element)) {
	e.forChild(FCfg{}, con)
}

func (e *Element) forChild(cfg FCfg, con func(ch *Element)) {
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
// it does not make sense to set index if element has no parent
// and nothin will happen
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

// ID is id getter
func (e *Element) ID() string {
	return e.id
}

// SetID changes element id, element is no longer accessable from old id
func (e *Element) SetID(id string) {
	if e.Scene != nil {
		delete(e.Scene.ids, e.id)
		e.Scene.ids[id] = e
	}
	e.id = id
}

// Group is group getter
func (e *Element) Group() string {
	return e.group
}

// SetGroup moves element from one group to another
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

// Hidden is hidden getter
func (e *Element) Hidden() bool {
	return e.hidden
}

// SetHidden is hidden setter
func (e *Element) SetHidden(value bool) {
	e.hidden = value
	e.onHiddenChange()
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

// MarginRealSize returns real size Margin will spam, Fill is considered 0
func (e *Element) MarginRealSize() mat.Vec {
	return sumAABB(e.Margin)
}

// PaddingSize returns total size element padding spams
func (e *Element) PaddingSize() mat.Vec {
	return e.Padding.Max.Add(e.Padding.Min)
}

func sumAABB(aabb mat.AABB) (sz mat.Vec) {
	l, b, r, t := aabb.Deco()
	if l != Fill {
		sz.X += l
	}
	if r != Fill {
		sz.X += r
	}
	if b != Fill {
		sz.Y += b
	}
	if t != Fill {
		sz.Y += t
	}

	return
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

// projectIndex is used when manipulating with child indexes
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
	// handle mouse exit, enter
	contains := e.Frame.Contains(w.MousePos())
	if e.Hovering && !contains {
		e.Hovering = false
		e.Events.Invoke(MouseExited, nil)
	} else if !e.Hovering && contains {
		e.Hovering = true
		e.Events.Invoke(MouseEntered, nil)
	}

	e.Module.Update(w, delta)

	e.forChild(IgnoreHidden, func(ch *Element) {
		ch.update(p, w, delta)
	})
}

// Redraw draws element and all its children to target, if preprocessor is not nil
// triangles are also preprocessed
func (e *Element) redraw(t ggl.Target, canvas *drw.Geom) {
	tar := t
	if e.Proc != nil {
		tar = e.Proc
	}

	canvas.Restart()
	e.Module.Draw(tar, canvas)
	e.forChild(IgnoreHidden, func(ch *Element) {
		ch.redraw(tar, canvas)
	})
	canvas.Restart()
	e.Module.DrawOnTop(tar, canvas)

	if e.Proc != nil {
		e.Proc.Fetch(t)
		e.Proc.Clear()
	}
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
	e.Module.PostInit()
}

// Move is next step after resize, size of all elements is calculated,
// now we can move them all to correct place
func (e *Element) move(offset mat.Vec, horizontal bool) mat.Vec {
	off := offset.Add(e.margin.Min).Add(e.Offest)
	e.Frame = e.size.ToAABB().Moved(off)
	off.AddE(e.Padding.Min)
	oOff := off
	e.Module.OnFrameChange()
	e.forChild(FCfg{
		Filter:  IgnoreHidden.Filter,
		Reverse: !e.Horizontal(),
	}, func(ch *Element) {
		if ch.Relative {
			ch.move(oOff, false)
		} else {
			off = ch.move(off, e.Horizontal())
		}

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

// FCfg is configuration for Element.forChild method
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

// even constants
const (
	MouseEntered = "mouse_entered"
	MouseExited  = "mouse_exited"
	Click        = "click"
	Deselect     = "deselect"
	Select       = "select"
	TextChanged  = "text_changed"
	Error        = "error"
	Enter        = "enter"
)

// InputState ...
type InputState struct {
	Hovering bool
}

// Module is most significant part of element, it implements all behavior, if you want modules to comunicate with each
// other put them into one module and make that instantiate elements for them, this way you can perform comunication in
// main module with no reflection
type Module interface {
	// DefaultStyle should returns default style of element that will be used as base, returning zero value is fine
	DefaultStyle() goss.Style
	// Init is called when module is inserted into element that is part of a scene
	Init(*Element)
	// PostInit is called after all children of element called Init
	PostInit()
	// Draw should draw the div, draw your triangles onto given target, you can use Geom as canvas
	// though you have to draw it to target too, Geom is cleared and restarted before draw call
	Draw(ggl.Target, *drw.Geom)
	// DrawOnTop does the same thing as draw, but on top of children
	DrawOnTop(ggl.Target, *drw.Geom)
	// Update is stage where your event handling and visual updates should happen
	Update(*ggl.Window, float64)
	// OnFrameChange is called by processor when frame of element changes
	OnFrameChange()

	MinSize() mat.Vec

	Width(takable, taken float64) float64
	Height(takable, taken float64) float64
}

// ModuleBase is a base of every module, you should embed this struct in your module
// and "override" default methods, though don't forget to call original Init that
// initializes the styles, if you don't give your element a module, this will be paced as placeholder
type ModuleBase struct {
	*Element
	Background mat.RGBA
}

// New implements ModuleFactory interface
func (m *ModuleBase) New() Module {
	return &ModuleBase{}
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

func (m *ModuleBase) PostInit() {}

// Draw implements Module interface
func (m *ModuleBase) Draw(t ggl.Target, g *drw.Geom) {
	g.Color(m.Background).AABB(m.Frame)
	g.Fetch(t)
}

// DrawOnTop implements Module interface
func (m *ModuleBase) DrawOnTop(t ggl.Target, g *drw.Geom) {
}

// Update implements Module interface
func (*ModuleBase) Update(*ggl.Window, float64) {}

// OnFrameChange implements Module interface
func (*ModuleBase) OnFrameChange() {}

// MinSize implements Module interface
func (m *ModuleBase) MinSize() (sz mat.Vec) {
	if m.Size.X != Fill {
		sz.X += m.Size.X
	}
	if m.Size.Y != Fill {
		sz.Y += m.Size.Y
	}
	return sz.Add(m.MarginRealSize())
}

// Width implements Module interface
func (m *ModuleBase) Width(takable, taken float64) float64 {
	return taken
}

// Height implements Module interface
func (m *ModuleBase) Height(takable, taken float64) float64 {
	return taken
}

// HSum calculates size of elements in horizontal composition
type HSum struct {
	*mat.Vec
}

// Add performs calculation
func (h HSum) Add(size mat.Vec) {
	h.X += size.X
	h.Y = math.Max(h.Y, size.Y)
}

// VSum is analogous to HSum just for horizontal composition
type VSum struct {
	*mat.Vec
}

// Add performs calculation
func (h VSum) Add(size mat.Vec) {
	h.Y += size.Y
	h.X = math.Max(h.X, size.X)
}

type summer interface {
	Add(mat.Vec)
}
