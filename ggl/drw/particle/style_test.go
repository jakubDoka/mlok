package particle

import (
	"reflect"
	"testing"

	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/lerp"
	"github.com/jakubDoka/mlok/mat/lerpc"
	"github.com/jakubDoka/mlok/mat/rgba"

	"github.com/jakubDoka/goml/goss"
)

func TestGradient(t *testing.T) {
	testCases := []struct {
		desc  string
		style goss.Style
		res   lerpc.Tween
	}{
		{
			desc: "constant",
			style: goss.Style{
				"": {"white"},
			},
			res: lerpc.Const(rgba.White),
		},
		{
			desc: "constant",
			style: goss.Style{
				"": {1},
			},
			res: lerpc.Const(rgba.White),
		},
		{
			desc: "constant",
			style: goss.Style{
				"": {1, 1, 1},
			},
			res: lerpc.Const(rgba.White),
		},
		{
			desc: "constant",
			style: goss.Style{
				"": {1, 1, 1, 1},
			},
			res: lerpc.Const(rgba.White),
		},
		{
			desc: "linear",
			style: goss.Style{
				"": {"black", "d", "white"},
			},
			res: lerpc.Linear(rgba.Black, rgba.White),
		},
		{
			desc: "chained",
			style: goss.Style{
				"": {"black", "d", 1, "d", 1, 0, 1},
			},
			res: lerpc.Chained(
				lerpc.Point(0, rgba.Black),
				lerpc.Point(0, rgba.White),
				lerpc.Point(0, mat.RGB(1, 0, 1)),
			),
		},
		{
			desc: "chained",
			style: goss.Style{
				"": {"black", .2, "d", 1, "p", .5, "d", 1, 0, 1, "p", 1},
			},
			res: lerpc.Chained(
				lerpc.Point(.2, rgba.Black),
				lerpc.Point(.5, rgba.White),
				lerpc.Point(1, mat.RGB(1, 0, 1)),
			),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			stl := WrapStyle(tC.style)
			grad := stl.ColorTween("", nil, nil)
			if !reflect.DeepEqual(grad, tC.res) {
				t.Error(grad)
			}
		})
	}
}

func TestParser(t *testing.T) {
	a := Assets{
		FloatTweens: map[string]lerp.Tween{
			"a": lerp.Const(10),
		},
	}

	p := Parser{
		Assets: a,
	}

	source := `
	style:
		emission_shape: rect 10 5;
		drawer: rect 10;
		color: green d white p 0.5 d purple p 1;
		mask: wheat;
		scale_multiplier: linear 0 1;
		acceleration: bezier 0 1 1 0;
		twerk_acceleration: 1.0
		velocity: random 100 200;
		livetime: 1;
		scale_x: 1;
		scale_y: 1;
		gravity: 0;
		friction: 1;
		rotation_relative_to_velocity: true;
	;`

	res := Type{
		EmissionShape: Rectangle{10, 5},
		base:          Square(10),

		Color: lerpc.Chained(
			lerpc.Point(0, rgba.Green),
			lerpc.Point(.5, rgba.White),
			lerpc.Point(1, rgba.Purple),
		),
		Mask: lerpc.Const(rgba.Wheat),

		ScaleMultiplier:   lerp.Linear(0, 1),
		Acceleration:      lerp.Bezier(0, 1, 1, 0),
		TwerkAcceleration: lerp.Const(0),

		Velocity: lerp.Random(100, 200),
		Livetime: lerp.Const(1),
		Rotation: lerp.Const(0),
		Spread:   lerp.Const(0),
		ScaleX:   lerp.Const(1),
		ScaleY:   lerp.Const(1),

		Friction:                   1,
		RotationRelativeToVelocity: true,
	}

	err := p.AddGoss([]byte(source))
	if err != nil {
		panic(err)
	}

	r := p.Construct("style")

	if !reflect.DeepEqual(r, &res) {
		t.Errorf("\n%#v\n%#v", r, &res)
	}
}
