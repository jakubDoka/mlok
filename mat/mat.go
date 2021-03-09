// Package mat focuses on mathematics, in case of more intense operations related to ray cast
// i left some documentation inside functions in case you want to understand how math here works
//
// missing features: Circle, Polygon
package mat

import (
	"fmt"
	"math"
	"strconv"
)

// Tran is standard transform with some utility, it plays well with ggl.Sprite
type Tran struct {
	Pos, Scl Vec
	Rot      float64
}

// Mat returns matrix reperesenting transform
func (t *Tran) Mat() Mat {
	return M(t.Pos, t.Scl, t.Rot)
}

// Mat is  matrix used for  transformation
// its a standard 3x3 matrix with setup:
//
// |X.X Y.X C.X|
// |X.Y Y.Y C.Y|
// | 0   0   1 |
//
// Structure can tell that X is X-axis direction vector Y si Y-axis direction vector and C
// is where they intersect IM represents system as it is, if you increase lengths of axis vectors
// Projected image will increase size, is you rotate both axises, image will be tilted, if you move it
// Image will get offset. You can combine these how every you like, just mind that order matters.
type Mat struct {
	X, Y, C Vec
}

// IM is Matrix used as base for all transformations
var IM = Mat{Vec{1, 0}, Vec{0, 1}, Vec{0, 0}}

// ZM is zero value Mat
var ZM Mat

// M is equivalent to:
//
// 	IM.ScaledXY(Vec{}, scl).Rotated(Vec{}, rot).Moved(pos)
//
// but 3x faster and shorter and covers mostly all you need
func M(pos, scl Vec, rot float64) Mat {
	sin, cos := math.Sincos(rot)
	return Mat{Vec{cos * scl.X, sin * scl.X}, Vec{-sin * scl.Y, cos * scl.Y}, pos}
}

// Raw returns raw representation of matrix
func (m Mat) Raw() [9]float32 {
	return [...]float32{
		float32(m.X.X), float32(m.X.Y), 0,
		float32(m.Y.X), float32(m.Y.Y), 0,
		float32(m.C.X), float32(m.C.Y), 1,
	}
}

// String returns a string representation of the Mat.
//
//   m := mat.IM
//   fmt.Println(m) // Mat(1 0 0 | 0 1 0)
func (m Mat) String() string {
	return fmt.Sprintf(
		"Mat(%v %v %v | %v %v %v)",
		ff(m.X.X), ff(m.Y.X), ff(m.C.X),
		ff(m.X.Y), ff(m.Y.Y), ff(m.C.Y),
	)
}

// Mv moves everything by the delta vector.
func (m Mat) Mv(delta Vec) Mat {
	m.C.AddE(delta)
	return m
}

// ScaledXY scales everything around a given point by the scale factor in each axis respectively.
func (m Mat) ScaledXY(around Vec, scale Vec) Mat {
	m.C.SubE(around)

	m.C.MulE(scale)
	m.X.MulE(scale)
	m.Y.MulE(scale)

	m.C.AddE(around)
	return m
}

// Scaled scales everything around a given point by the scale factor.
func (m Mat) Scaled(around Vec, scale float64) Mat {
	return m.ScaledXY(around, Vec{scale, scale})
}

// Rotated rotates everything around a given point by the given angle in radians.
func (m Mat) Rotated(around Vec, angle float64) Mat {
	sin, cos := math.Sincos(angle)
	m.C.SubE(around)
	m = m.Chained(Mat{Vec{cos, sin}, Vec{-sin, cos}, Vec{}})
	m.C.AddE(around)
	return m
}

// Chained adds another Mat to this one. All tranformations by the next Mat will be applied
// after the transformations of this Mat.
func (m Mat) Chained(next Mat) Mat {
	return Mat{
		Vec{
			next.X.X*m.X.X + next.Y.X*m.X.Y,
			next.X.Y*m.X.X + next.Y.Y*m.X.Y,
		},
		Vec{
			next.X.X*m.Y.X + next.Y.X*m.Y.Y,
			next.X.Y*m.Y.X + next.Y.Y*m.Y.Y,
		},
		Vec{
			next.X.X*m.C.X + next.Y.X*m.C.Y + next.C.X,
			next.X.Y*m.C.X + next.Y.Y*m.C.Y + next.C.Y,
		},
	}
}

// Project applies all transformations added to the Mat to a vector u and returns the result.
//
// Time complexity is O(1).
func (m Mat) Project(u Vec) Vec {
	return Vec{m.X.X*u.X + m.Y.X*u.Y + m.C.X, m.X.Y*u.X + m.Y.Y*u.Y + m.C.Y}
}

// Unproject does the inverse operation to Project.
//
// Time complexity is O(1).
func (m Mat) Unproject(u Vec) Vec {
	det := m.X.X*m.Y.Y - m.Y.X*m.X.Y
	return Vec{
		(m.Y.Y*(u.X-m.C.X) - m.Y.X*(u.Y-m.C.Y)) / det,
		(-m.X.Y*(u.X-m.C.X) + m.X.X*(u.Y-m.C.Y)) / det,
	}
}

// Approx same as normal approx but for matrix
func (m Mat) Approx(b Mat, precision int) bool {
	return m.X.Approx(b.X, precision) && m.Y.Approx(b.Y, precision) && m.C.Approx(b.C, precision)
}

// Approx returns whether two floats are same with certain precision
func Approx(a, b float64, precision int) bool {
	return math.Abs(Round(a, precision)-Round(b, precision)) < math.Pow(10, -float64(precision-1))
}

// Round rounds float with to given decimal points
func Round(v float64, precision int) float64 {
	scl := math.Pow(10, float64(precision))
	return math.Trunc(v*scl) / scl
}

// Clamp ...
func Clamp(val, min, max float64) float64 {
	if val > max {
		return max
	}
	if val < min {
		return min
	}
	return val
}

func ff(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
