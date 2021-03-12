package lerp

import "math/rand"

// BezierTween uses cubic Bézier curve equation for interpolation. For more info
// see https://en.wikipedia.org/wiki/B%C3%A9zier_curve,
// In case you are looking for visualization see https://www.desmos.com/calculator/d1ofwre0fr
// We are using only one dimensional curve as i do not consider to dimensions usefull for lib
// if you need to dimensions just use two Bezier-s
type BezierTween struct {
	Start, StartHandle, EndHandle, End float64
}

// ZB is Zero Bezier Curve
var ZB BezierTween

// Bezier is Bezier constructor
func Bezier(start, startHandle, endHandle, end float64) BezierTween {
	return BezierTween{start, startHandle, endHandle, end}
}

// Value returns Value along the curve interpolated by t
func (b BezierTween) Value(t float64) float64 {
	inv := 1.0 - t
	return b.Start*inv*inv*inv + b.StartHandle*inv*inv*t*3.0 + b.EndHandle*inv*t*t*3.0 + b.End*t*t*t
}

// LinearTween supports linear interpolation
type LinearTween struct {
	Start, End float64
}

// ZL is LinearTween zero value
var ZL LinearTween

// Linear is lerp
func Linear(start, end float64) Tween {
	return LinearTween{start, end}
}

// Value implements Tween interface
func (l LinearTween) Value(t float64) float64 {
	return l.Start + (l.End-l.Start)*t
}

// ConstantTween does nothing, its a placeholder with no overhead
type ConstantTween float64

// ZC is ConstantTween zero value
var ZC ConstantTween

// Const is COnstantTween constructor
func Const(value float64) ConstantTween {
	return ConstantTween(value)
}

// Value implements Tween interface
func (p ConstantTween) Value(t float64) float64 {
	return float64(p)
}

// Gen implements Generator interface
func (p ConstantTween) Gen() float64 {
	return float64(p)
}

// RandomGenerator generates random value between Min and Max
type RandomGenerator struct {
	Min, Max float64
}

// ZR is RandomGenerator zerovalue
var ZR RandomGenerator

// Random is RandomGenerator constructor
func Random(value, offset float64) RandomGenerator {
	return RandomGenerator{value, offset}
}

// Gen implements Generator interface
func (r RandomGenerator) Gen() float64 {
	return r.Min + (r.Max-r.Min)*rand.Float64()
}

// Tween is something that projects t into some other value
// that can make smooth movement
type Tween interface {
	Value(t float64) float64
}

// Generator should generate a value
type Generator interface {
	Gen() float64
}
