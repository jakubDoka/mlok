package particles

import (
	"gobatch/mat"
	"gobatch/mat/lerp"
	"gobatch/mat/rgba"
	"reflect"
	"testing"

	"github.com/jakubDoka/goml/goss"
)

func TestGradient(t *testing.T) {
	testCases := []struct {
		desc  string
		style goss.Style
		res   lerp.Gradient
	}{
		{
			desc: "constant",
			style: goss.Style{
				"": {"white"},
			},
			res: lerp.ConstColor(rgba.White),
		},
		{
			desc: "constant",
			style: goss.Style{
				"": {1},
			},
			res: lerp.ConstColor(rgba.White),
		},
		{
			desc: "constant",
			style: goss.Style{
				"": {1, 1, 1},
			},
			res: lerp.ConstColor(rgba.White),
		},
		{
			desc: "constant",
			style: goss.Style{
				"": {1, 1, 1, 1},
			},
			res: lerp.ConstColor(rgba.White),
		},
		{
			desc: "linear",
			style: goss.Style{
				"": {"black", "d", "white"},
			},
			res: lerp.LC(rgba.Black, rgba.White),
		},
		{
			desc: "chained",
			style: goss.Style{
				"": {"black", "d", 1, "d", 1, 0, 1},
			},
			res: lerp.ChainedColor{
				lerp.CP(0, rgba.Black),
				lerp.CP(0, rgba.White),
				lerp.CP(0, mat.RGB(1, 0, 1)),
			},
		},
		{
			desc: "chained",
			style: goss.Style{
				"": {"black", .2, "d", 1, "p", .5, "d", 1, 0, 1, "p", 1},
			},
			res: lerp.ChainedColor{
				lerp.CP(.2, rgba.Black),
				lerp.CP(.5, rgba.White),
				lerp.CP(1, mat.RGB(1, 0, 1)),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			stl := WrapStyle(tC.style)
			grad := stl.Gradient("", nil, nil)
			if !reflect.DeepEqual(grad, tC.res) {
				t.Error(grad)
			}
		})
	}
}
