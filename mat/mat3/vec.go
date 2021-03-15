package mat3

import (
	"math"

	"github.com/jakubDoka/gobatch/mat"
)

type Vec struct {
	X, Y, Z float64
}

func V(x, y, z float64) Vec {
	return Vec{x, y, z}
}

func (v Vec) Sub(o Vec) Vec {
	return Vec{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

func (v Vec) Cross2(o Vec) Vec {
	y := -(v.Z*o.X + o.Z*v.X) / (v.Y*o.X + o.Y*v.X)
	x := (v.Y*y + v.Z) / v.X
	return Vec{x, y, 1}
}

func (v Vec) Cross(o Vec) Vec {
	return Vec{
		v.Y*o.Z - v.Z*o.Y,
		v.X*o.Z - v.Z*o.X,
		v.Y*o.X - v.X*o.Y,
	}
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
