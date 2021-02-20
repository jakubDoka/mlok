package mat

import (
	"fmt"
	"image"
	"math"
)

// AABB is a  rectangle aligned with the axes of the coordinate system. It is defined by two
// points, Min and Max.
//
// The invariant should hold, that Max's components are greater or equal than Min's components
// respectively.
type AABB struct {
	Min, Max Vec
}

// A returns a new AABB with given the Min and Max coordinates.
//
// Note that the returned rectangle is not automatically normalized.
func A(minX, minY, maxX, maxY float64) AABB {
	return AABB{
		Min: Vec{minX, minY},
		Max: Vec{maxX, maxY},
	}
}

// ToAABB converts image.Rectangle to AABB
func ToAABB(a image.Rectangle) AABB {
	return AABB{
		Vec{float64(a.Min.X), float64(a.Min.Y)},
		Vec{float64(a.Max.X), float64(a.Max.Y)},
	}
}

// Cube returns AABB with center in c and width and height both equal to size * 2
func Cube(c Vec, size float64) AABB {
	return AABB{Vec{c.X - size, c.Y - size}, Vec{c.X + size, c.Y + size}}
}

// CR returns rect witch center is equal to c, width equal to w, likewise height equal to h
func CR(c Vec, w, h float64) AABB {
	w, h = w/2, h/2
	return AABB{Vec{X: c.X - w, Y: c.Y - h}, Vec{X: c.X + w, Y: c.Y + h}}
}

// FromRect converts image.Rectangle to AABB
func FromRect(a image.Rectangle) AABB {
	return AABB{
		Vec{float64(a.Min.X), float64(a.Min.Y)},
		Vec{float64(a.Max.X), float64(a.Max.Y)},
	}
}

// ToImage converts AABB to image.AABB
func (a AABB) ToImage() image.Rectangle {
	return image.Rect(
		int(a.Min.X),
		int(a.Min.Y),
		int(a.Max.X),
		int(a.Max.Y),
	)
}

// ToVec converts AABB to vec where x is AABB width anc y is rect Height
func (a AABB) ToVec() Vec {
	return a.Min.To(a.Max)
}

// String returns the string representation of the AABB.
func (a AABB) String() string {
	return fmt.Sprintf("A(%.3f %.3f %.3f %.3f)", a.Min.X, a.Min.Y, a.Max.X, a.Max.Y)
}

// Norm returns the AABB in normal form, such that Max is component-wise greater or equal than Min.
func (a AABB) Norm() AABB {
	return AABB{
		Min: Vec{
			math.Min(a.Min.X, a.Max.X),
			math.Min(a.Min.Y, a.Max.Y),
		},
		Max: Vec{
			math.Max(a.Min.X, a.Max.X),
			math.Max(a.Min.Y, a.Max.Y),
		},
	}
}

// W returns the width of the AABB.
func (a AABB) W() float64 {
	return a.Max.X - a.Min.X
}

// H returns the height of the AABB.
func (a AABB) H() float64 {
	return a.Max.Y - a.Min.Y
}

// Size returns the vector of width and height of the AABB.
func (a AABB) Size() Vec {
	return Vec{a.W(), a.H()}
}

// Area returns the area of a. If a is not normalized, area may be negative.
func (a AABB) Area() float64 {
	return a.W() * a.H()
}

// Center returns the position of the center of the AABB.
func (a AABB) Center() Vec {
	return a.Min.Lerp(a.Max, 0.5)
}

// Moved returns the AABB moved (both Min and Max) by the given vector delta.
func (a AABB) Moved(delta Vec) AABB {
	return AABB{
		Min: a.Min.Add(delta),
		Max: a.Max.Add(delta),
	}
}

// Resized returns the AABB resized to the given size while keeping the position of the given
// anchor.
//
//   a.Resized(a.Min, size)      // resizes while keeping the position of the lower-left corner
//   a.Resized(a.Max, size)      // same with the top-right corner
//   a.Resized(a.Center(), size) // resizes around the center
func (a AABB) Resized(anchor, size Vec) AABB {
	fraction := Vec{size.X / a.W(), size.Y / a.H()}
	return AABB{
		Min: anchor.Add(a.Min.Sub(anchor).Mul(fraction)),
		Max: anchor.Add(a.Max.Sub(anchor).Mul(fraction)),
	}
}

// ResizedMin returns the AABB resized to the given size while keeping the position of the AABB's
// Min.
//
// Sizes of zero area are safe here.
func (a AABB) ResizedMin(size Vec) AABB {
	return AABB{
		Min: a.Min,
		Max: a.Min.Add(size),
	}
}

// Contains checks whether a vector u is contained within this AABB (including it's borders).
func (a AABB) Contains(u Vec) bool {
	return a.Min.X <= u.X && u.X <= a.Max.X && a.Min.Y <= u.Y && u.Y <= a.Max.Y
}

// Union returns the minimal AABB which covers both a and s. AABBs a and s must be normalized.
func (a AABB) Union(s AABB) AABB {
	return A(
		math.Min(a.Min.X, s.Min.X),
		math.Min(a.Min.Y, s.Min.Y),
		math.Max(a.Max.X, s.Max.X),
		math.Max(a.Max.Y, s.Max.Y),
	)
}

// Intersect returns the maximal AABB which is covered by both a and s. AABBs a and s must be normalized.
//
// If a and s don't overlap, this function returns a zero-rectangle.
func (a AABB) Intersect(s AABB) AABB {
	t := A(
		math.Max(a.Min.X, s.Min.X),
		math.Max(a.Min.Y, s.Min.Y),
		math.Min(a.Max.X, s.Max.X),
		math.Min(a.Max.Y, s.Max.Y),
	)

	if t.Min.X >= t.Max.X || t.Min.Y >= t.Max.Y {
		return AABB{}
	}

	return t
}

// Intersects returns whether or not the given AABB intersects at any point with this AABB.
//
// This function is overall about 5x faster than Intersect, so it is better
// to use if you have no need for the returned AABB from Intersect.
func (a AABB) Intersects(s AABB) bool {
	return !(s.Max.X < a.Min.X ||
		s.Min.X > a.Max.X ||
		s.Max.Y < a.Min.Y ||
		s.Min.Y > a.Max.Y)
}

// Vertices returns a slice of the four corners which make up the rectangle.
func (a AABB) Vertices() [4]Vec {
	return [4]Vec{
		a.Min,
		{a.Min.X, a.Max.Y},
		a.Max,
		{a.Max.X, a.Min.Y},
	}
}

// LocalVertices creates array of vertices relative to center of rect
func (a AABB) LocalVertices() [4]Vec {
	v := a.Vertices()
	c := a.Center()

	for i, e := range v {
		v[i] = e.Sub(c)
	}

	return v
}

// VecBounds gets the smallest rectangle in witch all provided points fit in
func VecBounds(vectors ...Vec) (base AABB) {
	base.Min.X = math.MaxFloat64
	base.Min.Y = math.MaxFloat64
	base.Max.X = -math.MaxFloat64
	base.Max.Y = -math.MaxFloat64

	for _, v := range vectors {
		if base.Min.X > v.X {
			base.Min.X = v.X
		}
		if base.Min.Y > v.Y {
			base.Min.Y = v.Y
		}
		if base.Max.X < v.X {
			base.Max.X = v.X
		}
		if base.Max.Y < v.Y {
			base.Max.Y = v.Y
		}
	}

	return base
}

// Clamp clamps Vec inside AABB area
func (a AABB) Clamp(v Vec) Vec {
	return Vec{
		math.Max(math.Min(v.X, a.Max.X), a.Min.X),
		math.Max(math.Min(v.Y, a.Max.Y), a.Min.Y),
	}
}

// Flatten returns AABB flattened into Array, values are
// in same order as they would be stored on stack
func (a AABB) Flatten() [4]float64 {
	return [...]float64{a.Min.X, a.Min.Y, a.Max.X, a.Max.Y}
}

// Mutator is similar to Iterator but this gives option to mutate
// state of AABB trough Array Entries
func (a *AABB) Mutator() [4]*float64 {
	return [...]*float64{&a.Min.X, &a.Min.Y, &a.Max.X, &a.Max.Y}
}

// Deco returns edge values
func (a *AABB) Deco() (left, bottom, right, top float64) {
	return a.Min.X, a.Min.Y, a.Max.X, a.Max.Y
}
