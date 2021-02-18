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

func TestRGBAMath(t *testing.T) {
	if Alpha(1).Add(Alpha(1)) != Alpha(2) {
		t.Error("add")
	}
	if Alpha(1).Sub(Alpha(1)) != Alpha(0) {
		t.Error("sub")
	}
	if Alpha(2).Mul(Alpha(2)) != Alpha(4) {
		t.Error("mul")
	}
	if Alpha(2).Div(Alpha(2)) != Alpha(1) {
		t.Error("div")
	}
	if Alpha(2).Scaled(3) != Alpha(6) {
		t.Error("scl")
	}
	if LerpColor(Alpha(0), Alpha(10), .5) != Alpha(5) {
		t.Error(LerpColor(Alpha(0), Alpha(10), .5))
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
