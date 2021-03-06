package ui

import (
	"math"

	"github.com/jakubDoka/mlok/ggl/pck"
	"github.com/jakubDoka/mlok/ggl/txt"
	"github.com/jakubDoka/mlok/load"
	"github.com/jakubDoka/mlok/mat"

	"github.com/jakubDoka/goml/goss"
)

// Assets contains all data that is important for ui elements
type Assets struct {
	// spritesheet that scene will use
	pck.Sheet
	// Markdown is for text rendering
	Markdowns map[string]*txt.Markdown
	// Cursors are avaliable cursor drawers
	Cursors map[string]CursorDrawer
	// styles should be supplied from .goss files
	goss.Styles
}

// Props determinate look of Element and its properties
type Props struct {
	RawStyle
	// Margin defines spacing between elements, it supports Fill mode
	Margin mat.AABB
	// Padding defines how match space inside the elements should be free
	// Fill is not supported
	Padding mat.AABB
	// Size defines default size of element
	Size mat.Vec
	// Composition defines orientation of children in div, if horizontal
	// or vertical, if Composition.None() then it is initialized to be Vertical
	Composition
	// resize mode sets how element should react to size of children, see
	// constants documentation, if ResizeMode.None() then it is initialized with Expand
	Resizing [2]ResizeMode
	// Relative if true makes element ignore size of neighbors this it will be sized
	// and positioned merely by its offset
	Relative bool
	// Offset is used when moving elements inside scroll but if Relative is true
	// offset will be applied as offset from position where element would end up
	// if it were only element in the parent
	Offest mat.Vec
}

// Horizontal reports whether style composition is horizontal
func (s *Props) Horizontal() bool {
	return s.Composition == Horizontal
}

func (s *Props) Expands(dim Dimension) bool {
	return s.Resizing[dim] < Shrink
}

func (s *Props) Shrinks(dim Dimension) bool {
	return s.Resizing[dim] == Shrink || s.Resizing[dim] == Exact
}

func (s *Props) Ingors(dim Dimension) bool {
	return s.Resizing[dim] == Ignore
}

// Init initializes the style
func (s *Props) Init() {
	s.Margin = s.AABB("margin", s.Margin)
	s.Padding = s.AABB("padding", s.Padding)
	s.Size = s.Vec("size", s.Size)
	s.Composition = s.RawStyle.Composition("composition")

	s.Resizing[0] = s.ResizeMode("resizing_x")
	s.Resizing[1] = s.ResizeMode("resizing_y")
	if s.Resizing == [2]ResizeMode{} {
		r := s.ResizeMode("resizing")
		s.Resizing[0] = r
		s.Resizing[1] = r
	}

	s.Relative = s.Bool("relative", false)
	s.Offest = s.Vec("offset", mat.ZV)
}

type Dimension uint8

const (
	X Dimension = iota
	Y
)

// Composition ...
type Composition uint8

// None returns whether this is zero value
func (c Composition) None() bool {
	return c == 0
}

const (
	// Vertical makes children ordered vertically
	Vertical Composition = iota
	// Horizontal makes children ordered horizontally
	Horizontal
)

// Compositions maps each composition to its string representation
var Compositions = map[string]Composition{
	"vertical":   Vertical,
	"horizontal": Horizontal,
}

// ResizeMode ...
type ResizeMode uint8

// None returns whether this is zero value
func (c ResizeMode) None() bool {
	return c == 0
}

// Resize modes
const (
	// element expands to size of children but does not shrink under default size
	Expand ResizeMode = iota
	// element will shrink or expand based of children
	Exact
	// element will shrink to size of children, but children can overflow
	Shrink
	// element will ignore children size and leave overflow or underflow
	Ignore
)

// ResizeModes maps each resize mode to its string reperesentation
var ResizeModes = map[string]ResizeMode{
	"expand": Expand,
	"exact":  Exact,
	"shrink": Shrink,
	"ignore": Ignore,
}

// Fill - if Margin is equal Fill it will take a remaining space in element
// if there are multiple elements with fill margin, space is split between them
const Fill float64 = load.Fill
const unknown float64 = math.MaxFloat64

var fill filler

type filler struct{}

func (f filler) hSum(a mat.AABB) float64 {
	return f.value(a.Min.X) + f.value(a.Max.X)
}

func (f filler) vSum(a mat.AABB) float64 {
	return f.value(a.Min.Y) + f.value(a.Max.Y)
}

func (filler) add(dest *float64, scr float64) {
	if scr != Fill {
		*dest += scr
	}
}

func (filler) set(dest *float64, scr float64) {
	if scr != Fill {
		*dest = scr
	}
}

func (filler) value(scr float64) float64 {
	if scr != Fill {
		return scr
	}
	return 0
}

// RawStyle extends load.RawStyle by some functionality
type RawStyle struct {
	load.RawStyle
}

func (r RawStyle) Align(key string, def txt.Align) (v txt.Align) {
	v = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	switch vl := val[0].(type) {
	case int:
		return txt.Align(vl)
	case float64:
		return txt.Align(vl)
	case string:
		return txt.Aligns[vl]
	}

	return
}

// CursorDrawer retrieves a cursor drawer from style
func (r RawStyle) CursorDrawer(key string, drawers map[string]CursorDrawer, def CursorDrawer) (v CursorDrawer) {
	v = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	switch vl := val[0].(type) {
	case string:
		val, ok := drawers[vl]
		if ok {
			return val
		}
	}

	return
}

// Composition parses style composition, if parsing fails, Vertical is returned
func (r RawStyle) Composition(key string) (c Composition) {
	val, ok := r.Style[key]
	if !ok {
		return
	}

	switch v := val[0].(type) {
	case int:
		return Composition(v)
	case string:
		return Compositions[v]
	}
	return
}

// ResizeMode parser resize mode, if pasring fails Expand is returned
func (r RawStyle) ResizeMode(key string) (e ResizeMode) {
	val, ok := r.Style[key]
	if !ok {
		return
	}

	switch v := val[0].(type) {
	case int:
		return ResizeMode(v)
	case string:
		return ResizeModes[v]
	}
	return
}
