package ui

import (
	"math"

	"github.com/jakubDoka/gobatch/ggl/pck"
	"github.com/jakubDoka/gobatch/ggl/txt"
	"github.com/jakubDoka/gobatch/load"
	"github.com/jakubDoka/gobatch/mat"

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

// Init initializes the style
func (s *Props) Init() {
	s.Margin = s.AABB("margin", s.Margin)
	s.Padding = s.AABB("padding", s.Padding)
	s.Size = s.Vec("size", s.Size)
	s.Composition = s.RawStyle.Composition("composition")

	s.Resizing[0] = s.ResizeMode("resize_mode_x")
	s.Resizing[1] = s.ResizeMode("resize_mode_y")
	if s.Resizing == [2]ResizeMode{} {
		r := s.ResizeMode("resize_mode")
		s.Resizing[0] = r
		s.Resizing[1] = r
	}

	s.Relative = s.Bool("relative", false)
	s.Offest = s.Vec("offset", mat.ZV)
}

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

var compositions = map[string]Composition{
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

var resizeModes = map[string]ResizeMode{
	"Expand": Expand,
	"exact":  Exact,
	"shring": Shrink,
	"ignore": Ignore,
}

// Fill - if Margin is equal Fill it will take a remaining space in element
// if there are multiple elements with fill margin, space is split between them
const Fill float64 = load.Fill
const unknown float64 = math.MaxFloat64

// RawStyle extends load.RawStyle by some functionality
type RawStyle struct {
	load.RawStyle
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
		return compositions[v]
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
		return resizeModes[v]
	}
	return
}
