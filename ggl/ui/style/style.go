package style

import "gobatch/mat"

type Style struct {
	Margin, Padding AABB
	Size, Position  Vec
	Relation, Side  Anchor

	Background, Mask        RGBA
	Font, Markdown, Texture string
}

type AABB []float64

func (a AABB) Dec() (res mat.AABB) {
	switch len(a) {
	case 1:
		res = mat.A(a[0], a[0], a[0], a[0])
	case 2:
		res = mat.A(a[0], a[1], a[0], a[1])
	case 3:
		res = mat.A(a[0], a[1], a[0], a[2])
	case 4:
		res = mat.A(a[0], a[1], a[2], a[3])
	}

	return
}

type RGBA [4]float64

func (c RGBA) Dec() mat.RGBA {
	return mat.RGBA{
		R: c[0],
		G: c[1],
		B: c[2],
		A: c[3],
	}
}

type Vec [2]float64

func (v Vec) Dec() mat.Vec {
	return mat.V(v[0], v[1])
}

// Anchor ...
type Anchor int8

// Anchor enumeration
const (
	Left Anchor = iota
	Right
	Top
	Bottom

	BottomLeft
	TopLeft
	TopRight
	BottomRight

	Center
	Fill
	Children
)

var convertor = map[string]Anchor{
	"left":        Left,
	"right":       Right,
	"top":         Top,
	"bottom":      Bottom,
	"bottomleft":  BottomLeft,
	"topleft":     TopLeft,
	"topright":    TopRight,
	"bottomright": BottomRight,
	"center":      Center,
	"fill":        Fill,
	"children":    Children,
}
