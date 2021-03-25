package lerpc

import (
	"math/rand"

	"github.com/jakubDoka/mlok/mat"
)

// ConstantTween is wrapper around mat.RGBA that implements Gradient,
// but it does nothing with the value
type ConstantTween mat.RGBA

// Const is ConstantTween constructor
func Const(value mat.RGBA) ConstantTween {
	return ConstantTween(value)
}

// Value implements Gradient interface
func (r ConstantTween) Value(t float64) mat.RGBA {
	return mat.RGBA(r)
}

// Gen implements Generator interface
func (r ConstantTween) Gen() mat.RGBA {
	return mat.RGBA(r)
}

// LinearTween supports linear interpolation between two colors encapsulated
// in struct for interface reasons
type LinearTween struct {
	A, B mat.RGBA
}

// Linear is LinearTween constructor
func Linear(a, b mat.RGBA) LinearTween {
	return LinearTween{a, b}
}

// Value implements Tween interface
func (c LinearTween) Value(t float64) mat.RGBA {
	return mat.LerpColor(c.A, c.B, t)
}

// TweenPoint ...
type TweenPoint struct {
	Position float64
	Color    mat.RGBA
}

// Point is TweenPoint constructor
func Point(position float64, color mat.RGBA) TweenPoint {
	return TweenPoint{position, color}
}

// ChainedTween allows linear interpolation between multiple colors
// Simple resulting color will be interpolation between that two points
// that contain t, interpolation is of corse relative and smooth
type ChainedTween []TweenPoint

// Chained is ChainedTween constructor
func Chained(tps ...TweenPoint) ChainedTween {
	return ChainedTween(tps)
}

// Value implements Tween interface
//
// if t < then position of first color, first color will be returned
// if t > then last color position, last color will be returned
func (c ChainedTween) Value(t float64) mat.RGBA {
	for i := 0; i < len(c); i++ {
		if t < c[i].Position {
			if i == 0 {
				return c[i].Color
			}

			return mat.LerpColor(
				c[i-1].Color,
				c[i].Color,
				(t-c[i-1].Position)/(c[i].Position-c[i-1].Position),
			)
		}
	}
	return c[len(c)-1].Color
}

// RandomGenerator generates random color in range between
// two colors
type RandomGenerator struct {
	Min, Max mat.RGBA
}

// Random is RandomGenerator constructor
func Random(min, max mat.RGBA) RandomGenerator {
	return RandomGenerator{min, max}
}

// Gen implements generator interface
func (r RandomGenerator) Gen() mat.RGBA {
	return mat.RGBA{
		R: r.Min.R + (r.Max.R-r.Min.R)*rand.Float64(),
		G: r.Min.G + (r.Max.G-r.Min.G)*rand.Float64(),
		B: r.Min.B + (r.Max.B-r.Min.B)*rand.Float64(),
		A: r.Min.A + (r.Max.A-r.Min.A)*rand.Float64(),
	}
}

// Tween is something that returns color based of t, alike tween
// is used for interpolation, just with colors
type Tween interface {
	Value(t float64) mat.RGBA
}

// Generator should generate a value
type Generator interface {
	Gen() mat.RGBA
}
