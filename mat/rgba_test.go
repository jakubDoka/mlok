package mat

import (
	"image/color"
	"math"
	"testing"
)

func TestRGBA(t *testing.T) {
	if RGB(1, 1, 1) != Alpha(1) {
		t.Fail()
	}

	r, g, b, a := Alpha(1).RGBA()
	if r != math.MaxUint16 || g != math.MaxUint16 || b != math.MaxUint16 || a != math.MaxUint16 {
		t.Error(r, g, b, a)
	}

	white := color.RGBA{255, 255, 255, 255}
	res := rgbaModel(white)
	if res != White {
		t.Error(res, 0)
	}

	res = rgbaModel(Alpha(0))
	if res != Alpha(0) {
		t.Error(res, 1)
	}
}

func ch(a float64) RGBA {
	return RGBA{a, a, a, a}
}

func TestRGBAMath(t *testing.T) {
	if ch(1).Add(ch(1)) != ch(2) {
		t.Error("add")
	}
	if ch(1).Sub(ch(1)) != ch(0) {
		t.Error("sub")
	}
	if ch(2).Mul(ch(2)) != ch(4) {
		t.Error("mul")
	}
	if ch(2).Div(ch(2)) != ch(1) {
		t.Error("div")
	}
	if ch(2).Scaled(3) != ch(6) {
		t.Error("scl")
	}
	if LerpColor(ch(0), ch(10), .5) != ch(5) {
		t.Error(LerpColor(ch(0), ch(10), .5))
	}
}

func TestHexToRGBA(t *testing.T) {
	testCases := []struct {
		desc string
		hex  string
		fail bool
		res  RGBA
	}{
		{
			desc: "no alpha",
			hex:  "ffffff",
			res:  Alpha(1),
		},
		{
			desc: "alpha",
			hex:  "ffffff00",
			res:  RGBA{1, 1, 1, 0},
		},
		{
			desc: "invalid char",
			hex:  "fffffk",
			fail: true,
		},
		{
			desc: "too short",
			hex:  "fffff",
			fail: true,
		},
		{
			desc: "all types",
			hex:  "00ffFF",
			res:  RGB(0, 1, 1),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res, err := HexToRGBA(tC.hex)
			if tC.fail != (err != nil) {
				t.Fail()
			}

			if tC.fail {
				return
			}

			if res != tC.res {
				t.Error(res)
			}
		})
	}
}
