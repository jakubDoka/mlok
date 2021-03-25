package mat3

import (
	"math"

	"github.com/jakubDoka/mlok/mat"
)

type Vec struct {
	X, Y, Z float64
}

func V(x, y, z float64) Vec {
	return Vec{x, y, z}
}

// Rotated rotates vector around pivot by angle with a right hand rule
// so tongue is pivot and curved fingers point to the direction of rotation
func (v Vec) Rotated(angle float64, pivot Vec) Vec {
	colinear := pivot.Scaled(v.Dot(pivot) / v.Dot(v))
	orthogonal := v.Sub(colinear)
	length := orthogonal.Len()

	// x and y are two vectors orthogonal and normalized, we can now use tham as local
	// coordinate system
	y := v.Cross(pivot).Normalized()
	x := orthogonal.Divided(length)
	sin, cos := math.Sincos(angle)

	// now we use coordinate system to project the angle, then scale it to original
	// length and finally add the previously subtracted component
	return x.Scaled(cos).Add(y.Scaled(sin)).Scaled(length).Add(colinear)
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

func (v Vec) Sub(o Vec) Vec {
	return Vec{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

func (v Vec) Mul(o Vec) Vec {
	return Vec{v.X * o.X, v.Y * o.Y, v.Z * o.Z}
}

func (v *Vec) AddE(o Vec) {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
}

func (v *Vec) MulE(o Vec) {
	v.X *= o.X
	v.Y *= o.Y
	v.Z *= o.Z
}

func (v Vec) Scaled(scalar float64) Vec {
	return Vec{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vec) Divided(scalar float64) Vec {
	return Vec{v.X / scalar, v.Y / scalar, v.Z / scalar}
}

func (v Vec) Cross(o Vec) Vec {
	return Vec{
		v.Y*o.Z - v.Z*o.Y,
		v.X*o.Z - v.Z*o.X,
		v.Y*o.X - v.X*o.Y,
	}
}

func (v Vec) Dot(o Vec) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vec) Normalized() Vec {
	len := v.Len()
	return Vec{
		v.X / len,
		v.Y / len,
		v.Z / len,
	}
}

func (v Vec) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec) Inv() Vec {
	return Vec{-v.X, -v.Y, -v.Z}
}

func (v Vec) Approx(o Vec, precision int) bool {
	return mat.Approx(v.X, o.X, precision) && mat.Approx(v.Y, o.Y, precision) && mat.Approx(v.Z, o.Z, precision)
}
