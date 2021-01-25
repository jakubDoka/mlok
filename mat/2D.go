package mat

import (
	"fmt"
	"image"
	"math"
)

// V2 is a 2D vector type with X and Y coordinates.
//
// Create vectors with the V constructor:
//
//   u := mat.NV2(1, 2)
//   v := mat.NV2(8, -3)
//
// Use various methods to manipulate them:
//
//   w := u.Add(v)
//   fmt.Println(w)        // V2(9, -1)
//   fmt.Println(u.Sub(v)) // V2(-7, 5)
//   u = mat.NV2(2, 3)
//   v = mat.NV2(8, 1)
//   if u.X < 0 {
//	     fmt.Println("this won't happen")
//   }
//   x := u.Unit().Dot(v.Unit())
type V2 struct {
	X, Y float64
}

// Vector related constants
var (
	Origin2 = V2{}
	Scale2  = V2{1, 1}
)

// NV2 returns a new 2D vector with the given coordinates.
func NV2(x, y float64) V2 {
	return V2{x, y}
}

// Rad2 returns vector from representation of radial cordenates
func Rad2(angle, length float64) V2 {
	return V2{math.Cos(angle) * length, math.Sin(angle) * length}
}

// String returns the string representation of the vector u.
//
//   u := mat.V(4.5, -1.3)
//   u.String()     // returns "V2(4.5, -1.3)"
//   fmt.Println(u) // V2(4.5, -1.3)
func (u V2) String() string {
	return fmt.Sprintf("V2(%v, %v)", u.X, u.Y)
}

// XY returns the components of the vector in two return values.
func (u V2) XY() (x, y float64) {
	return u.X, u.Y
}

// Add returns the sum of vectors u and v.
func (u V2) Add(v V2) V2 {
	return V2{
		u.X + v.X,
		u.Y + v.Y,
	}
}

// Sub returns the difference betweeen vectors u and v.
func (u V2) Sub(v V2) V2 {
	return V2{
		u.X - v.X,
		u.Y - v.Y,
	}
}

// Div returns the vector u divided by the vector v component-wise.
func (u V2) Div(v V2) V2 {
	return V2{u.X / v.X, u.Y / v.Y}
}

// Mul returns the vector u multiplied by the vector v component-wise.
func (u V2) Mul(v V2) V2 {
	return V2{u.X * v.X, u.Y * v.Y}
}

// AddE modifies caller
func (u *V2) AddE(v V2) {
	u.X += v.X
	u.Y += v.Y
}

// SubE modifies caller
func (u *V2) SubE(v V2) {
	u.X -= v.X
	u.Y -= v.Y
}

// MulE modifies caller
func (u *V2) MulE(v V2) {
	u.X *= v.X
	u.Y *= v.Y
}

// DivE modifies caller
func (u *V2) DivE(v V2) {
	u.X /= v.X
	u.Y /= v.Y
}

// Floor converts x and y to their integer equivalents.
func (u V2) Floor() V2 {
	return V2{
		math.Floor(u.X),
		math.Floor(u.Y),
	}
}

// To returns the vector from u to v. Equivalent to v.Sub(u).
func (u V2) To(v V2) V2 {
	return V2{
		v.X - u.X,
		v.Y - u.Y,
	}
}

// Scaled returns the vector u multiplied by c.
func (u V2) Scaled(c float64) V2 {
	return V2{u.X * c, u.Y * c}
}

// Divided returns the vector u divided by c.
func (u V2) Divided(c float64) V2 {
	return V2{u.X / c, u.Y / c}
}

// Inv returns u with both components inverted
func (u V2) Inv() V2 {
	u.X = -u.X
	u.Y = -u.Y
	return u
}

// Len returns the length of the vector u.
func (u V2) Len() float64 {
	return math.Hypot(u.X, u.Y)
}

// Angle returns the angle between the vector u and the x-axis. The result is in range [-Pi, Pi].
func (u V2) Angle() float64 {
	return math.Atan2(u.Y, u.X)
}

// Normalized returns a vector of length 1 facing the direction of u (has the same angle).
func (u V2) Normalized() V2 {
	if u.X == 0 && u.Y == 0 {
		return V2{1, 0}
	}
	return u.Scaled(1 / u.Len())
}

// Rotated returns the vector u rotated by the given angle in radians.
func (u V2) Rotated(angle float64) V2 {
	sin, cos := math.Sincos(angle)
	return V2{
		u.X*cos - u.Y*sin,
		u.X*sin + u.Y*cos,
	}
}

// Normal returns a vector normal to u. Equivalent to u.Rotated(math.Pi / 2), but faster.
func (u V2) Normal() V2 {
	return V2{-u.Y, u.X}
}

// Dot returns the dot product of vectors u and v.
func (u V2) Dot(v V2) float64 {
	return u.X*v.X + u.Y*v.Y
}

// Cross return the cross product of vectors u and v.
func (u V2) Cross(v V2) float64 {
	return u.X*v.Y - v.X*u.Y
}

// Project returns a projection (or component) of vector u in the direction of vector v.
//
// Behaviour is undefined if v is a zero vector.
func (u V2) Project(v V2) V2 {
	len := u.Dot(v) / v.Len()
	return v.Normalized().Scaled(len)
}

// AngleTo returns angle between u and v.
func (u V2) AngleTo(v V2) float64 {
	a := math.Abs(u.Angle() - v.Angle())
	if a > math.Pi {
		return 2*math.Pi - a
	}
	return a
}

// Map applies the function f to both x and y components of the vector u and returns the modified
// vector.
//
//   u := mat.V(10.5, -1.5)
//   v := u.Map(math.Floor)   // v is V2(10, -2), both components of u floored
func (u V2) Map(f func(float64) float64) V2 {
	return V2{
		f(u.X),
		f(u.Y),
	}
}

// Approx same as normal approx but for vector
func (u V2) Approx(b V2, precision int) bool {
	return Approx(u.X, b.X, precision) && Approx(u.Y, b.Y, precision)
}

// Lerp returns a linear interpolation between vectors a and b.
//
// This function basically returns a point along the line between a and b and t chooses which one.
// If t is 0, then a will be returned, if t is 1, b will be returned. Anything between 0 and 1 will
// return the appropriate point between a and b and so on.
func (u V2) Lerp(b V2, t float64) V2 {
	return u.Scaled(1 - t).Add(b.Scaled(t))
}

// Mat2 is 2D matrix used for 2D transformation
// its a standard 3x3 matrix with setup:
//
// |X.X Y.X C.X|
// |X.Y Y.Y C.Y|
// | 0   0   1 |
//
// Structure can tell that X is X-axis direction vector Y si Y-axis direction vector and C
// is where they intersect IM2 represents system as it is, if you increase lengths of axis vectors
// Projected image will increase size, is you rotate both axises, image will be tilted, if you move it
// Image will get offset. You can combine these how every you like, just mind that order matters.
type Mat2 struct {
	X, Y, C V2
}

// IM2 is Matrix used as base for all transformations
var IM2 = Mat2{V2{1, 0}, V2{0, 1}, V2{0, 0}}

// NMat2 is equivalent to:
//
// 	IM2.ScaledXY(V2{}, scl).Rotated(V2{}, rot).Moved(pos)
//
// but 3x faster and shorter and covers mostly all you need
func NMat2(pos, scl V2, rot float64) Mat2 {
	sin, cos := math.Sincos(rot)
	return Mat2{V2{cos * scl.X, sin * scl.X}, V2{-sin * scl.Y, cos * scl.Y}, pos}
}

// Raw returns raw representation of matrix
func (m Mat2) Raw() [9]float32 {
	return [...]float32{
		float32(m.X.X), float32(m.X.Y), 0,
		float32(m.Y.X), float32(m.Y.Y), 0,
		float32(m.C.X), float32(m.C.Y), 1,
	}
}

// String returns a string representation of the Mat2.
//
//   m := mat.IM2
//   fmt.Println(m) // Mat2(1 0 0 | 0 1 0)
func (m Mat2) String() string {
	return fmt.Sprintf(
		"Mat2(%v %v %v | %v %v %v)",
		m.X.X, m.Y.X, m.C.X,
		m.X.Y, m.Y.Y, m.C.Y,
	)
}

// Mv moves everything by the delta vector.
func (m Mat2) Mv(delta V2) Mat2 {
	m.C.AddE(delta)
	return m
}

// ScaledXY scales everything around a given point by the scale factor in each axis respectively.
func (m Mat2) ScaledXY(around V2, scale V2) Mat2 {
	m.C.SubE(around)

	m.C.MulE(scale)
	m.X.MulE(scale)
	m.Y.MulE(scale)

	m.C.AddE(around)
	return m
}

// Scaled scales everything around a given point by the scale factor.
func (m Mat2) Scaled(around V2, scale float64) Mat2 {
	return m.ScaledXY(around, V2{scale, scale})
}

// Rotated rotates everything around a given point by the given angle in radians.
func (m Mat2) Rotated(around V2, angle float64) Mat2 {
	sin, cos := math.Sincos(angle)
	m.C.SubE(around)
	m = m.Chained(Mat2{V2{cos, sin}, V2{-sin, cos}, V2{}})
	m.C.AddE(around)
	return m
}

// Chained adds another Mat2 to this one. All tranformations by the next Mat2 will be applied
// after the transformations of this Mat2.
func (m Mat2) Chained(next Mat2) Mat2 {
	return Mat2{
		V2{
			next.X.X*m.X.X + next.Y.X*m.X.Y,
			next.X.Y*m.X.X + next.Y.Y*m.X.Y,
		},
		V2{
			next.X.X*m.Y.X + next.Y.X*m.Y.Y,
			next.X.Y*m.Y.X + next.Y.Y*m.Y.Y,
		},
		V2{
			next.X.X*m.C.X + next.Y.X*m.C.Y + next.C.X,
			next.X.Y*m.C.X + next.Y.Y*m.C.Y + next.C.Y,
		},
	}
}

// Project applies all transformations added to the Mat2 to a vector u and returns the result.
//
// Time complexity is O(1).
func (m Mat2) Project(u V2) V2 {
	return V2{m.X.X*u.X + m.Y.X*u.Y + m.C.X, m.X.Y*u.X + m.Y.Y*u.Y + m.C.Y}
}

// Unproject does the inverse operation to Project.
//
// Time complexity is O(1).
func (m Mat2) Unproject(u V2) V2 {
	det := m.X.X*m.Y.Y - m.Y.X*m.X.Y
	return V2{
		(m.Y.Y*(u.X-m.C.X) - m.Y.X*(u.Y-m.C.Y)) / det,
		(-m.X.Y*(u.X-m.C.X) + m.X.X*(u.Y-m.C.Y)) / det,
	}
}

// Approx same as normal approx but for matrix
func (m Mat2) Approx(b Mat2, precision int) bool {
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

// AABB is a 2D rectangle aligned with the axes of the coordinate system. It is defined by two
// points, Min and Max.
//
// The invariant should hold, that Max's components are greater or equal than Min's components
// respectively.
type AABB struct {
	Min, Max V2
}

// NAABB returns a new AABB with given the Min and Max coordinates.
//
// Note that the returned rectangle is not automatically normalized.
func NAABB(minX, minY, maxX, maxY float64) AABB {
	return AABB{
		Min: V2{minX, minY},
		Max: V2{maxX, maxY},
	}
}

// ToAABB converts image.Rectangle to AABB
func ToAABB(r image.Rectangle) AABB {
	return AABB{
		V2{float64(r.Min.X), float64(r.Min.Y)},
		V2{float64(r.Max.X), float64(r.Max.Y)},
	}
}

// Cube returns AABB with center in c and width and height both equal to size * 2
func Cube(c V2, size float64) AABB {
	return AABB{V2{c.X - size, c.Y - size}, V2{c.X + size, c.Y + size}}
}

// CR returns rect witch center is equal to c, width equal to w, likewise height equal to h
func CR(c V2, w, h float64) AABB {
	w, h = w/2, h/2
	return AABB{V2{X: c.X - w, Y: c.Y - h}, V2{X: c.X + w, Y: c.Y + h}}
}

// FromRect converts image.Rectangle to AABB
func FromRect(r image.Rectangle) AABB {
	return AABB{
		V2{float64(r.Min.X), float64(r.Min.Y)},
		V2{float64(r.Max.X), float64(r.Max.Y)},
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
func (r AABB) ToVec() V2 {
	return r.Min.To(r.Max)
}

// String returns the string representation of the AABB.
//
//   r := mat.NAABB(100, 50, 200, 300)
//   r.String()     // returns "AABB(100, 50, 200, 300)"
//   fmt.Println(r) // AABB(100, 50, 200, 300)
func (r AABB) String() string {
	return fmt.Sprintf("AABB(%v, %v, %v, %v)", r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

// Norm returns the AABB in normal form, such that Max is component-wise greater or equal than Min.
func (r AABB) Norm() AABB {
	return AABB{
		Min: V2{
			math.Min(r.Min.X, r.Max.X),
			math.Min(r.Min.Y, r.Max.Y),
		},
		Max: V2{
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
func (r AABB) Size() V2 {
	return V2{r.W(), r.H()}
}

// Area returns the area of r. If r is not normalized, area may be negative.
func (r AABB) Area() float64 {
	return r.W() * r.H()
}

// Center returns the position of the center of the AABB.
func (r AABB) Center() V2 {
	return r.Min.Lerp(r.Max, 0.5)
}

// Moved returns the AABB moved (both Min and Max) by the given vector delta.
func (r AABB) Moved(delta V2) AABB {
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
func (r AABB) Resized(anchor, size V2) AABB {
	if r.W()*r.H() == 0 {
		panic(fmt.Errorf("(%T).Resize: zero area", r))
	}
	fraction := V2{size.X / r.W(), size.Y / r.H()}
	return AABB{
		Min: anchor.Add(r.Min.Sub(anchor).Mul(fraction)),
		Max: anchor.Add(r.Max.Sub(anchor).Mul(fraction)),
	}
}

// ResizedMin returns the AABB resized to the given size while keeping the position of the AABB's
// Min.
//
// Sizes of zero area are safe here.
func (r AABB) ResizedMin(size V2) AABB {
	return AABB{
		Min: r.Min,
		Max: r.Min.Add(size),
	}
}

// Contains checks whether a vector u is contained within this AABB (including it's borders).
func (r AABB) Contains(u V2) bool {
	return r.Min.X <= u.X && u.X <= r.Max.X && r.Min.Y <= u.Y && u.Y <= r.Max.Y
}

// Union returns the minimal AABB which covers both r and s. AABBs r and s must be normalized.
func (r AABB) Union(s AABB) AABB {
	return NAABB(
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
	t := NAABB(
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
func (r AABB) Vertices() [4]V2 {
	return [4]V2{
		r.Min,
		{r.Min.X, r.Max.Y},
		r.Max,
		{r.Max.X, r.Min.Y},
	}
}

// LocalVertices creates array of vertices relative to center of rect
func (r AABB) LocalVertices() [4]V2 {
	v := r.Vertices()
	c := r.Center()

	for i, e := range v {
		v[i] = e.Sub(c)
	}

	return v
}

// VecBounds2 gets the smallest rectangle in witch all provided points fit in
func VecBounds2(vectors ...V2) (base AABB) {
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

// Ray2 is a standard 2D raycast, it supports calculating intersections with all collizion shapes
type Ray2 struct {
	A, B V2
}

// Tran2 is standard 2D transform with some utility, it plays well with ggl.Sprite2D
type Tran2 struct {
	Pos, Scl V2
	Rot      float64
}

// Mat returns matrix reperesenting transform
func (t *Tran2) Mat() Mat2 {
	return NMat2(t.Pos, t.Scl, t.Rot)
}
