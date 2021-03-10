package ui

import (
	"gobatch/ggl/pck"
	"gobatch/ggl/txt"
	"gobatch/mat"
	"gobatch/mat/rgba"
	"math"

	"github.com/jakubDoka/goml/goss"
)

// Assets contains all data that is important for ui elements
type Assets struct {
	// spritesheet that scene will use
	*pck.Sheet
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
const Fill float64 = 100000.767898765556788777667787666
const unknown float64 = math.MaxFloat64

// RawStyle is wrapper of goss.Style and adds extra functionality
type RawStyle struct {
	goss.Style
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

// Bool returns boolean value under the key of def, if retrieval fails
func (r RawStyle) Bool(key string, def bool) (v bool) {
	v = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	switch v := val[0].(type) {
	case bool:
		return v
	case int:
		return v == 1
	}

	return
}

// Float returns float under the key of default value if obtaining failed
func (r RawStyle) Float(key string, def float64) (v float64) {
	v = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	switch v := val[0].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}

	return v
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

// Ident returns string from Style or def if not present
func (r RawStyle) Ident(key, def string) string {
	val, ok := r.Style.Ident(key)
	if !ok {
		return def
	}
	return val
}

// Vec returns vector under tha key, if parsing fails or vec is not present def is returned
func (r RawStyle) Vec(key string, def mat.Vec) (u mat.Vec) {
	u = def

	val, ok := r.Style[key]
	if !ok {
		return
	}

	components := [2]float64{}
	for i := 0; i < len(components) && i < len(val); i++ {
		switch v := val[i].(type) {
		case float64:
			components[i] = v
		case int:
			components[i] = float64(v)
		case string:
			if v != "fill" {
				return
			}
			components[i] = Fill
		default:
			return
		}
	}

	switch len(val) {
	case 1:
		return mat.V(components[0], components[0])
	case 2:
		return mat.V(components[0], components[1])
	}

	return
}

// Region returns a texture region, it can be hardcoded aabb or region name.
func (r RawStyle) Region(key string, regions map[string]mat.AABB, def mat.AABB) mat.AABB {
	m, ok := regions[r.Ident(key, "")]
	if !ok {
		return r.AABB(key, def)
	}
	return m
}

// AABB parser margin under the key, if parsing fails or margin is not present, default is returned
func (r RawStyle) AABB(key string, def mat.AABB) (m mat.AABB) {
	m = def

	val, ok := r.Style[key]
	if !ok {
		return
	}

	sides := [4]float64{}
	for i := 0; i < len(sides) && i < len(val); i++ {
		switch v := val[i].(type) {
		case float64:
			sides[i] = v
		case int:
			sides[i] = float64(v)
		case string:
			if v != "fill" {
				return
			}
			sides[i] = Fill
		default:
			return
		}

	}

	switch len(val) {
	case 1:
		return mat.A(sides[0], sides[0], sides[0], sides[0])
	case 2:
		return mat.A(sides[0], sides[1], sides[0], sides[1])
	case 4:
		return mat.A(sides[0], sides[1], sides[2], sides[3])
	}

	return
}

// RGBA returns a color under the key, if color parsing fails or color is not present, def is returned
// this also accepts names mapped in gobatch/mat/rgba.Colors
func (r RawStyle) RGBA(key string, def mat.RGBA) (c mat.RGBA) {
	c = def

	val, ok := r.Style[key]
	if !ok {
		return
	}

	if v, ok := val[0].(string); ok {
		if v, ok := rgba.Colors[v]; ok {
			return v
		}
		return
	}

	channels := [4]float64{}
	for i := 0; i < len(channels) && i < len(val); i++ {

		switch vl := val[i].(type) {
		case float64:
			channels[i] = vl
		case int:
			channels[i] = float64(vl)
		default:
			return
		}
	}

	switch len(val) {
	case 1:
		return mat.Alpha(channels[0])
	case 3:
		return mat.RGB(channels[0], channels[1], channels[2])
	case 4:
		return mat.RGBA{
			R: channels[0],
			G: channels[1],
			B: channels[2],
			A: channels[3],
		}
	}

	return
}
