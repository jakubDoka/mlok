package lerp

// Bezier is cubic Bézier curve used for interpolation. For more info
// see https://en.wikipedia.org/wiki/B%C3%A9zier_curve,
// In case you are looking for visualization see https://www.desmos.com/calculator/d1ofwre0fr
// We are using only one dimensional curve as i do not consider to dimensions usefull for lib
// if you need to dimensions just use two Bezier-s
type Bezier struct {
	Start, StartHandle, EndHandle, End float64
}

// ZB is Zero Bezier Curve
var ZB Bezier

// B is Bezier constructor
func B(start, startHandle, endHandle, end float64) Bezier {
	return Bezier{start, startHandle + start, end - endHandle, end}
}

// Float returns Float along the curve interpolated by t
func (b Bezier) Float(t float64) float64 {
	inv := 1.0 - t
	return b.Start*inv*inv*inv + b.StartHandle*inv*inv*t*3.0 + b.EndHandle*inv*t*t*3.0 + b.End*t*t*t
}

// Linear supports linear interpolation
type Linear struct {
	Start, End float64
}

// ZL is Linear zero value
var ZL Linear

// L is lerp
func L(start, end float64) Linear {
	return Linear{start, end}
}

// Point implements Interpolator interface
func (l Linear) Point(t float64) float64 {
	return l.Start + (l.End-l.Start)*t
}

// Constant does nothing, its a placeholder with no overhead
type Constant struct {
	Value float64
}

// ZC is Constant zero value
var ZC Constant

// Float implements Tween interface
func (p Constant) Float(t float64) float64 {
	return p.Value
}

// Tween is something that projects t into some other value
// that can make smooth movement
type Tween interface {
	Float(t float64) float64
}
