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
func ToAABB(r image.Rectangle) AABB {
	return AABB{
		Vec{float64(r.Min.X), float64(r.Min.Y)},
		Vec{float64(r.Max.X), float64(r.Max.Y)},
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
func FromRect(r image.Rectangle) AABB {
	return AABB{
		Vec{float64(r.Min.X), float64(r.Min.Y)},
		Vec{float64(r.Max.X), float64(r.Max.Y)},
	}
}

// ToImage converts AABB to image.AABB
func (r AABB) ToImage() image.Rectangle {
	return image.Rect(
		int(r.Min.X),
		int(r.Min.Y),
		int(r.Max.X),
		int(r.Max.Y),
	)
}

// ToVec converts AABB to vec where x is AABB width anc y is rect Height
func (r AABB) ToVec() Vec {
	return r.Min.To(r.Max)
}

// String returns the string representation of the AABB.
func (r AABB) String() string {
	return fmt.Sprintf("A(%.3f %.3f %.3f %.3f)", r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

// Norm returns the AABB in normal form, such that Max is component-wise greater or equal than Min.
func (r AABB) Norm() AABB {
	return AABB{
		Min: Vec{
			math.Min(r.Min.X, r.Max.X),
			math.Min(r.Min.Y, r.Max.Y),
		},
		Max: Vec{
			math.Max(r.Min.X, r.Max.X),
			math.Max(r.Min.Y, r.Max.Y),
		},
	}
}

// W returns the width of the AABB.
func (r AABB) W() float64 {
	return r.Max.X - r.Min.X
}

// H returns the height of the AABB.
func (r AABB) H() float64 {
	return r.Max.Y - r.Min.Y
}

// Size returns the vector of width and height of the AABB.
func (r AABB) Size() Vec {
	return Vec{r.W(), r.H()}
}

// Area returns the area of r. If r is not normalized, area may be negative.
func (r AABB) Area() float64 {
	return r.W() * r.H()
}

// Center returns the position of the center of the AABB.
func (r AABB) Center() Vec {
	return r.Min.Lerp(r.Max, 0.5)
}

// Moved returns the AABB moved (both Min and Max) by the given vector delta.
func (r AABB) Moved(delta Vec) AABB {
	return AABB{
		Min: r.Min.Add(delta),
		Max: r.Max.Add(delta),
	}
}

// Resized returns the AABB resized to the given size while keeping the position of the given
// anchor.
//
//   r.Resized(r.Min, size)      // resizes while keeping the position of the lower-left corner
//   r.Resized(r.Max, size)      // same with the top-right corner
//   r.Resized(r.Center(), size) // resizes around the center
//
// This function does not make sense for resizing a rectangle of zero area and will panic. Use
// ResizedMin in the case of zero area.
func (r AABB) Resized(anchor, size Vec) AABB {
	if r.W()*r.H() == 0 {
		panic(fmt.Errorf("(%T).Resize: zero area", r))
	}
	fraction := Vec{size.X / r.W(), size.Y / r.H()}
	return AABB{
		Min: anchor.Add(r.Min.Sub(anchor).Mul(fraction)),
		Max: anchor.Add(r.Max.Sub(anchor).Mul(fraction)),
	}
}

// ResizedMin returns the AABB resized to the given size while keeping the position of the AABB's
// Min.
//
// Sizes of zero area are safe here.
func (r AABB) ResizedMin(size Vec) AABB {
	return AABB{
		Min: r.Min,
		Max: r.Min.Add(size),
	}
}

// Contains checks whether a vector u is contained within this AABB (including it's borders).
func (r AABB) Contains(u Vec) bool {
	return r.Min.X <= u.X && u.X <= r.Max.X && r.Min.Y <= u.Y && u.Y <= r.Max.Y
}

// Union returns the minimal AABB which covers both r and s. AABBs r and s must be normalized.
func (r AABB) Union(s AABB) AABB {
	return A(
		math.Min(r.Min.X, s.Min.X),
		math.Min(r.Min.Y, s.Min.Y),
		math.Max(r.Max.X, s.Max.X),
		math.Max(r.Max.Y, s.Max.Y),
	)
}

// Intersect returns the maximal AABB which is covered by both r and s. AABBs r and s must be normalized.
//
// If r and s don't overlap, this function returns a zero-rectangle.
func (r AABB) Intersect(s AABB) AABB {
	t := A(
		math.Max(r.Min.X, s.Min.X),
		math.Max(r.Min.Y, s.Min.Y),
		math.Min(r.Max.X, s.Max.X),
		math.Min(r.Max.Y, s.Max.Y),
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
func (r AABB) Intersects(s AABB) bool {
	return !(s.Max.X < r.Min.X ||
		s.Min.X > r.Max.X ||
		s.Max.Y < r.Min.Y ||
		s.Min.Y > r.Max.Y)
}

// Vertices returns a slice of the four corners which make up the rectangle.
func (r AABB) Vertices() [4]Vec {
	return [4]Vec{
		r.Min,
		{r.Min.X, r.Max.Y},
		r.Max,
		{r.Max.X, r.Min.Y},
	}
}

// LocalVertices creates array of vertices relative to center of rect
func (r AABB) LocalVertices() [4]Vec {
	v := r.Vertices()
	c := r.Center()

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
