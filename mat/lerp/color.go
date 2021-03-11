package lerp

import "gobatch/mat"

// ConstColor is wrapper around mat.RGBA that implements Gradient
type ConstColor mat.RGBA

// Color implements Gradient interface
func (r ConstColor) Color(t float64) mat.RGBA {
	return mat.RGBA(r)
}

// LinearColor supports linear interpolation between two colors encapsulated
// in Struct so is implements Gradient interface
type LinearColor struct {
	A, B mat.RGBA
}

// LC is LinearColor constructor
func LC(a, b mat.RGBA) LinearColor {
	return LinearColor{a, b}
}

// Color implements gradient interface
func (c LinearColor) Color(t float64) mat.RGBA {
	return mat.LerpColor(c.A, c.A, t)
}

// ColorPoint ...
type ColorPoint struct {
	Position float64
	Color    mat.RGBA
}

// CP is colorPoint constructor
func CP(position float64, color mat.RGBA) ColorPoint {
	return ColorPoint{position, color}
}

// ChainedColor is used for color interpolation, mainly when you need to interpolate
// trough multiple colors
type ChainedColor []ColorPoint

// Color implements Gradient interface
//
// if t < then position of first color, first color will be returned
// if t > then last color position, last color will be returned
func (c ChainedColor) Color(t float64) mat.RGBA {
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

// Gradient is something that returns color based of t, alike tween
// is used for interpolation, just with colors
type Gradient interface {
	Color(t float64) mat.RGBA
}
