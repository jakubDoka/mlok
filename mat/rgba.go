package mat

import (
	"errors"
	"image/color"
)

// Lerp linearly interpolates between two floats
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// LerpColor does linear interpolation between two colors
func LerpColor(a, b RGBA, t float64) RGBA {
	return RGBA{
		R: Lerp(a.R, b.R, t),
		G: Lerp(a.G, b.G, t),
		B: Lerp(a.B, b.B, t),
		A: Lerp(a.A, b.A, t),
	}
}

// RGBA represents an alpha-premultiplied RGBA color with components within range [0, 1].
//
// The difference between color.RGBA is that the value range is [0, 1] and the values are floats.
type RGBA struct {
	R, G, B, A float64
}

// Color constants
var (
	Transparent = RGBA{}
	Black       = RGB(0, 0, 0)
	White       = RGB(1, 1, 1)
	Red         = RGB(1, 0, 0)
	Green       = RGB(0, 1, 0)
	Blue        = RGB(0, 0, 1)
)

// RGB returns a fully opaque RGBA color with the given RGB values.
//
// A common way to construct a transparent color is to create one with RGB constructor, then
// multiply it by a color obtained from the Alpha constructor.
func RGB(r, g, b float64) RGBA {
	return RGBA{r, g, b, 1}
}

// Alpha returns a white RGBA color with the given alpha component.
func Alpha(a float64) RGBA {
	return RGBA{a, a, a, a}
}

// Add adds color d to color r component-wise and returns the result (the components are not
// clamped).
func (r RGBA) Add(d RGBA) RGBA {
	return RGBA{
		R: r.R + d.R,
		G: r.G + d.G,
		B: r.B + d.B,
		A: r.A + d.A,
	}
}

// Sub subtracts color d from color r component-wise and returns the result (the components
// are not clamped).
func (r RGBA) Sub(d RGBA) RGBA {
	return RGBA{
		R: r.R - d.R,
		G: r.G - d.G,
		B: r.B - d.B,
		A: r.A - d.A,
	}
}

// Mul multiplies color r by color d component-wise (the components are not clamped).
func (r RGBA) Mul(d RGBA) RGBA {
	return RGBA{
		R: r.R * d.R,
		G: r.G * d.G,
		B: r.B * d.B,
		A: r.A * d.A,
	}
}

// Div divides r by d component-wise (the components are not clamped).
func (r RGBA) Div(d RGBA) RGBA {
	return RGBA{
		A: r.A / d.A,
		B: r.B / d.B,
		R: r.R / d.R,
		G: r.G / d.G,
	}
}

// Scaled multiplies each component of color r by scale and returns the result (the components
// are not clamped).
func (r RGBA) Scaled(scale float64) RGBA {
	return RGBA{
		R: r.R * scale,
		G: r.G * scale,
		B: r.B * scale,
		A: r.A * scale,
	}
}

// RGBA returns alpha-premultiplied red, green, blue and alpha components of the RGBA color.
func (r RGBA) RGBA() (rc, g, b, a uint32) {
	rc = uint32(0xffff * r.R)
	g = uint32(0xffff * r.G)
	b = uint32(0xffff * r.B)
	a = uint32(0xffff * r.A)
	return
}

// ErrInvalidHex is returned by HexToRGBA if hex string contains non ex characters
var ErrInvalidHex = errors.New("byte is not a hex code")

// ErrTooShort is returned by HexToRGBA if hex string is too short to parse a color
var ErrTooShort = errors.New("hex string is too short (min is 6)")

// HexToRGBA converts hex string to RGBA
func HexToRGBA(s string) (r RGBA, err error) {
	if len(s) < 6 {
		return r, ErrTooShort
	}

	hexToByte := func(b byte) (r byte) {
		switch {
		case b >= '0' && b <= '9':
			r = b - '0'
		case b >= 'a' && b <= 'f':
			r = b - 'a' + 10
		case b >= 'A' && b <= 'F':
			r = b - 'A' + 10
		default:
			err = ErrInvalidHex
		}

		return
	}

	r.R = float64(hexToByte(s[0])<<4+hexToByte(s[1])) / 0xFF
	r.G = float64(hexToByte(s[2])<<4+hexToByte(s[3])) / 0xFF
	r.B = float64(hexToByte(s[4])<<4+hexToByte(s[5])) / 0xFF
	if len(s) == 8 {
		r.A = float64(hexToByte(s[6])<<4+hexToByte(s[7])) / 0xFF
	} else {
		r.A = 1
	}

	return r, err
}

// ToRGBA converts a color to RGBA format. Using this function is preferred to using RGBAModel, for
// performance (using RGBAModel introduces additional unnecessary allocations).
func ToRGBA(c color.Color) RGBA {
	if r, ok := c.(RGBA); ok {
		return r
	}
	r, g, b, a := c.RGBA()
	return RGBA{
		float64(r) / 0xffff,
		float64(g) / 0xffff,
		float64(b) / 0xffff,
		float64(a) / 0xffff,
	}
}

// RGBAModel converts colors to RGBA format.
var RGBAModel = color.ModelFunc(rgbaModel)

func rgbaModel(r color.Color) color.Color {
	return ToRGBA(r)
}
