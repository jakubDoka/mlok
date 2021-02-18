package mat

import (
	"fmt"
	"math"
)

// Vec is a  vector type with X and Y coordinates.
//
// Create vectors with the V constructor:
//
//   u := mat.V(1, 2)
//   v := mat.V(8, -3)
//
// Use various methods to manipulate them:
//
//   w := u.Add(v)
//   fmt.Println(w)        // Vec(9, -1)
//   fmt.Println(u.Sub(v)) // Vec(-7, 5)
//   u = mat.V(2, 3)
//   v = mat.V(8, 1)
//   if u.X < 0 {
//	     fmt.Println("this won't happen")
//   }
//   x := u.Unit().Dot(v.Unit())
type Vec struct {
	X, Y float64
}

// Vector related constants
var (
	Origin = Vec{}
	Scale  = Vec{1, 1}
)

// V returns a new  vector with the given coordinates.
func V(x, y float64) Vec {
	return Vec{x, y}
}

// Rad returns vector from representation of radial cordenates
func Rad(angle, length float64) Vec {
	s, c := math.Sincos(angle)
	return Vec{c * length, s * length}
}

// String returns the string representation of the vector u.
func (u Vec) String() string {
	return fmt.Sprintf("V(%.3f %.3f)", u.X, u.Y)
}

// XY returns the components of the vector in two return values.
func (u Vec) XY() (x, y float64) {
	return u.X, u.Y
}

// Add returns the sum of vectors u and v.
func (u Vec) Add(v Vec) Vec {
	return Vec{
		u.X + v.X,
		u.Y + v.Y,
	}
}

// Sub returns the difference betweeen vectors u and v.
func (u Vec) Sub(v Vec) Vec {
	return Vec{
		u.X - v.X,
		u.Y - v.Y,
	}
}

// Div returns the vector u divided by the vector v component-wise.
func (u Vec) Div(v Vec) Vec {
	return Vec{u.X / v.X, u.Y / v.Y}
}

// Mul returns the vector u multiplied by the vector v component-wise.
func (u Vec) Mul(v Vec) Vec {
	return Vec{u.X * v.X, u.Y * v.Y}
}

// To returns the vector from u to v. Equivalent to v.Sub(u).
func (u Vec) To(v Vec) Vec {
	return Vec{
		v.X - u.X,
		v.Y - u.Y,
	}
}

// AddE modifies caller
func (u *Vec) AddE(v Vec) {
	u.X += v.X
	u.Y += v.Y
}

// SubE modifies caller
func (u *Vec) SubE(v Vec) {
	u.X -= v.X
	u.Y -= v.Y
}

// MulE modifies caller
func (u *Vec) MulE(v Vec) {
	u.X *= v.X
	u.Y *= v.Y
}

// DivE modifies caller
func (u *Vec) DivE(v Vec) {
	u.X /= v.X
	u.Y /= v.Y
}

// Floor converts x and y to their integer equivalents.
func (u Vec) Floor() Vec {
	return Vec{
		math.Floor(u.X),
		math.Floor(u.Y),
	}
}

// Scaled returns the vector u multiplied by c.
func (u Vec) Scaled(c float64) Vec {
	return Vec{u.X * c, u.Y * c}
}

// Divided returns the vector u divided by c.
func (u Vec) Divided(c float64) Vec {
	return Vec{u.X / c, u.Y / c}
}

// Inv returns u with both components inverted
func (u Vec) Inv() Vec {
	u.X = -u.X
	u.Y = -u.Y
	return u
}

// Len returns the length of the vector u.
func (u Vec) Len() float64 {
	return math.Hypot(u.X, u.Y)
}

// Len2 returns length*length of vector witch is faster then Len
// you can for example comparing Len2 of two vectors yields same results
// as comparing Len
func (u Vec) Len2() float64 {
	return u.X*u.X + u.Y*u.Y
}

// Angle returns the angle between the vector u and the x-axis. The result is in range [-Pi, Pi].
func (u Vec) Angle() float64 {
	return math.Atan2(u.Y, u.X)
}

// Normalized returns a vector of length 1 facing the direction of u (has the same angle).
func (u Vec) Normalized() Vec {
	if u.X == 0 && u.Y == 0 {
		return Vec{1, 0}
	}
	return u.Scaled(1 / u.Len())
}

// Rotated returns the vector u rotated by the given angle in radians.
func (u Vec) Rotated(angle float64) Vec {
	sin, cos := math.Sincos(angle)
	return Vec{
		u.X*cos - u.Y*sin,
		u.X*sin + u.Y*cos,
	}
}

// Normal returns a vector normal to u. Equivalent to u.Rotated(math.Pi / 2), but faster.
func (u Vec) Normal() Vec {
	return Vec{-u.Y, u.X}
}

// Norm returns normal vector of line, size can be adjusted
func (u Vec) Norm(len float64) Vec {
	return u.Normal().Normalized().Scaled(len)
}

// Dot returns the dot product of vectors u and v.
func (u Vec) Dot(v Vec) float64 {
	return u.X*v.X + u.Y*v.Y
}

// Cross return the cross product of vectors u and v.
func (u Vec) Cross(v Vec) float64 {
	return u.X*v.Y - v.X*u.Y
}

// AngleTo returns angle between u and v.
func (u Vec) AngleTo(v Vec) float64 {
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
//   v := u.Map(math.Floor)   // v is Vec(10, -2), both components of u floored
func (u Vec) Map(f func(float64) float64) Vec {
	return Vec{
		f(u.X),
		f(u.Y),
	}
}

// Approx same as normal approx but for vector
func (u Vec) Approx(b Vec, precision int) bool {
	return Approx(u.X, b.X, precision) && Approx(u.Y, b.Y, precision)
}

// Lerp returns a linear interpolation between vectors a and b.
//
// This function basically returns a point along the line between a and b and t chooses which one.
// If t is 0, then a will be returned, if t is 1, b will be returned. Anything between 0 and 1 will
// return the appropriate point between a and b and so on.
func (u Vec) Lerp(b Vec, t float64) Vec {
	return u.Scaled(1 - t).Add(b.Scaled(t))
}
