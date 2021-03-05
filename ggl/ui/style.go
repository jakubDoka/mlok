package ui

import (
	"gobatch/ggl/pck"
	"gobatch/ggl/txt"
	"gobatch/mat"

	"github.com/jakubDoka/goml/goss"
)

// Assets contains all data that is important for ui elements
type Assets struct {
	// spritesheet that scene will use
	*pck.Sheet
	// Markdown is for text rendering
	Markdowns map[string]*txt.Markdown
	// styles should be supplied from .goss files
	goss.Styles
}

// Props determinate look of Element and its properties
type Props struct {
	RawStyle
	// Margin defines spacing between elements, it supports Fill mode
	Margin mat.AABB
	// Size defines default size of element
	Size mat.Vec
	// Composition defines orientation of children in div, if horizontal
	// or vertical, if Composition.None() then it is initialized to be Vertical
	Composition
	// resize mode sets how element should react to size of children, see
	// constants documentation, if ResizeMode.None() then it is initialized with Expand
	ResizeMode
}

// Horizontal reports whether style composition is horizontal
func (s *Props) Horizontal() bool {
	return s.Composition == Horizontal
}

// Init initializes the style
func (s *Props) Init() {
	s.Margin = s.AABB("margin", mat.ZA)
	s.Size = s.Vec("size", mat.ZV)
	c, _ := s.Int("composition")
	s.Composition = Composition(c)
	c, _ = s.Int("resize_mode")
	s.ResizeMode = ResizeMode(c)
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

// Fill - if Margin is equal Fill it will take a remaining space in element
// if there are multiple elements with fill margin, space is split between them
const Fill float64 = 100000.767898765556788777667787666

// RawStyle is wrapper of goss.Style and adds extra functionality
type RawStyle struct {
	goss.Style
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
func (r RawStyle) RGBA(key string, def mat.RGBA) (c mat.RGBA) {
	c = def

	val, ok := r.Style[key]
	if !ok {
		return
	}

	channels := [4]float64{}
	for i := 0; i < len(channels) && i < len(val); i++ {
		channels[i], ok = val[i].(float64)
		if !ok {
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
