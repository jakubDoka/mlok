package mat

import (
	"fmt"
	"math"
)

// Vec is a  vector type with X and Y coordinates.
//
// Create vectors with the V constructor:
//
//   v := mat.V(1, 2)
//   v := mat.V(8, -3)
//
// Use various methods to manipulate them:
//
//   w := v.Add(v)
//   fmt.Println(w)        // Vec(9, -1)
//   fmt.Println(v.Sub(v)) // Vec(-7, 5)
//   v = mat.V(2, 3)
//   v = mat.V(8, 1)
//   if v.X < 0 {
//	     fmt.Println("this won't happen")
//   }
//   x := v.Unit().Dot(v.Unit())
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

// ToAABB turns Vec into AABB where Min is V(0 0) and Max is v
func (v Vec) ToAABB() AABB {
	return AABB{Max: v}
}

// Rad returns vector from representation of radial cordenates
func Rad(angle, length float64) Vec {
	s, c := math.Sincos(angle)
	return Vec{c * length, s * length}
}

// String returns the string representation of the vector v.
func (v Vec) String() string {
	return fmt.Sprintf("V(%.3f %.3f)", v.X, v.Y)
}

// XY returns the components of the vector in two return values.
func (v Vec) XY() (x, y float64) {
	return v.X, v.Y
}

// Add returns the sum of vectors v and v.
func (v Vec) Add(u Vec) Vec {
	return Vec{
		v.X + u.X,
		v.Y + u.Y,
	}
}

// Sub subtracts u from v and returns recult.
func (v Vec) Sub(u Vec) Vec {
	return Vec{
		v.X - u.X,
		v.Y - u.Y,
	}
}

// Div returns the vector v divided by the vector u component-wise.
func (v Vec) Div(u Vec) Vec {
	return Vec{v.X / u.X, v.Y / u.Y}
}

// Mul returns the vector v multiplied by the vector u component-wise.
func (v Vec) Mul(u Vec) Vec {
	return Vec{v.X * u.X, v.Y * u.Y}
}

// To returns the vector from v to u. Equivalent to u.Sub(v).
func (v Vec) To(u Vec) Vec {
	return Vec{
		u.X - v.X,
		u.Y - v.Y,
	}
}

// AddE same as v = v.Add(u)
func (v *Vec) AddE(u Vec) {
	v.X += u.X
	v.Y += u.Y
}

// SubE same as v = v.Sub(u)
func (v *Vec) SubE(u Vec) {
	v.X -= u.X
	v.Y -= u.Y
}

// MulE same as v = v.Mul(u)
func (v *Vec) MulE(u Vec) {
	v.X *= u.X
	v.Y *= u.Y
}

// DivE same as v = v.Div(u)
func (v *Vec) DivE(u Vec) {
	v.X /= u.X
	v.Y /= u.Y
}

// Floor converts x and y to their integer equivalents.
func (v Vec) Floor() Vec {
	return Vec{
		math.Floor(v.X),
		math.Floor(v.Y),
	}
}

// Scaled returns the vector v multiplied by c.
func (v Vec) Scaled(c float64) Vec {
	return Vec{v.X * c, v.Y * c}
}

// Divided returns the vector v divided by c.
func (v Vec) Divided(c float64) Vec {
	return Vec{v.X / c, v.Y / c}
}

// Inv returns v with both components inverted
func (v Vec) Inv() Vec {
	v.X = -v.X
	v.Y = -v.Y
	return v
}

// Len returns the length of the vector v.
func (v Vec) Len() float64 {
	return math.Hypot(v.X, v.Y)
}

// Len2 returns length*length of vector witch is faster then Len
// you can for example comparing Len2 of two vectors yields same results
// as comparing Len
func (v Vec) Len2() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Angle returns the angle between the vector v and the x-axis. The result is in range [-Pi, Pi].
func (v Vec) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

// Normalized returns a vector of length 1 facing the direction of v (has the same angle).
func (v Vec) Normalized() Vec {
	if v.X == 0 && v.Y == 0 {
		return Vec{1, 0}
	}
	return v.Scaled(1 / v.Len())
}

// Rotated returns the vector v rotated by the given angle in radians.
func (v Vec) Rotated(angle float64) Vec {
	sin, cos := math.Sincos(angle)
	return Vec{
		v.X*cos - v.Y*sin,
		v.X*sin + v.Y*cos,
	}
}

// Normal returns a vector normal to v. Equivalent to v.Rotated(math.Pi / 2), but faster.
func (v Vec) Normal() Vec {
	return Vec{-v.Y, v.X}
}

// Norm returns normal vector of line, size can be adjusted
func (v Vec) Norm(len float64) Vec {
	return v.Normal().Normalized().Scaled(len)
}

// Dot returns the dot product of vectors v and u.
func (v Vec) Dot(u Vec) float64 {
	return v.X*u.X + v.Y*u.Y
}

// Cross return the cross product of vectors v and u.
func (v Vec) Cross(u Vec) float64 {
	return v.X*u.Y - u.X*v.Y
}

// Max uses math.Max on both components and returns resulting vector
func (v Vec) Max(u Vec) Vec {
	return Vec{
		math.Max(v.X, u.X),
		math.Max(v.Y, u.Y),
	}
}

// Min uses math.Min on both components and returns resulting vector
func (v Vec) Min(u Vec) Vec {
	return Vec{
		math.Min(v.X, u.X),
		math.Min(v.Y, u.Y),
	}
}

// AngleTo returns angle between v and v.
func (v Vec) AngleTo(u Vec) float64 {
	a := math.Abs(v.Angle() - u.Angle())
	if a > math.Pi {
		return 2*math.Pi - a
	}
	return a
}

// Map applies the function f to both x and y components of the vector v and returns the modified
// vector.
//
//   v := mat.V(10.5, -1.5)
//   v := v.Map(math.Floor)   // v is Vec(10, -2), both components of v floored
func (v Vec) Map(f func(float64) float64) Vec {
	return Vec{
		f(v.X),
		f(v.Y),
	}
}

// Approx same as normal approx but for vector
func (v Vec) Approx(b Vec, precision int) bool {
	return Approx(v.X, b.X, precision) && Approx(v.Y, b.Y, precision)
}

// Lerp returns a linear interpolation between vectors a and b.
//
// This function basically returns a point along the line between a and b and t chooses which one.
// If t is 0, then a will be returned, if t is 1, b will be returned. Anything between 0 and 1 will
// return the appropriate point between a and b and so on.
func (v Vec) Lerp(b Vec, t float64) Vec {
	return v.Scaled(1 - t).Add(b.Scaled(t))
}

// Flatten flatens the Vec into Array, values are
// ordered as they would on stack
func (v Vec) Flatten() [2]float64 {
	return [...]float64{v.X, v.Y}
}

// Mutator similar to Flatten returns array with vector components
// though these are pointers to componenets instead
func (v *Vec) Mutator() [2]*float64 {
	return [...]*float64{&v.X, &v.Y}
}
