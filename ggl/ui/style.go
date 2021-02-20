package ui

import (
	"gobatch/ggl/pck"
	"gobatch/ggl/txt"
	"gobatch/mat"
)

// Assets contains all data that is important for ui elements
type Assets struct {
	*pck.Sheet
	// Markdown is for text rendering
	Markdowns map[string]*txt.Markdown
	// styles all on one place can be easily loaded from json
	Styles map[string]*Style
}

// Style determinate look of Div and its properties
type Style struct {
	// Margin defines spacing between elements, it supports Fill mode
	Margin mat.AABB
	// Size defines default size of element
	Size mat.Vec
	// Horizontal defines orientation of children in div, if horizontal
	// is true, children will be rendered next to each other, if false
	// on top of each other
	Horizontal bool
	// resize mode sets how div should react to size of children, see
	// constants documentation
	ResizeMode
	// Current saves substile
	Current int
	// Substiles are for greater flexibility, they store data that can very
	// how your module uses data in substile is on your implementation
	Subs []SubStyle
	// As Style cannot cover all kinds of data you need for your custom modules
	// you can store it in Data
	Data interface{}
}

// SubStyle stores data that can have multiple variants for one div,
// so div can easily switch between them, data that is currently in SubStile
// is specific to build-in div modules
type SubStyle struct {
	Background, Mask mat.RGBA

	Font, Markdown, Texture string

	Data interface{}
}

// ResizeMode ...
type ResizeMode uint8

// Resize modes
const (
	// div expands to size of children but does not shrink under default size
	Expand ResizeMode = iota
	// div will shrink to size of children, but children can overflow
	Shrink
	// div will shrink or expand based of children
	Exact
	// div will ignore children size and leave overflow or underflow
	Ignore
)

// Fill - if Margin is equal Fill it will take a remaining space in element
// if there are multiple elements with fill margin, space is split between them
const Fill float64 = 100000.767898765556788777667787666
