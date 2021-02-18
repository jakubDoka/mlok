package mat

import (
	"math"
	"testing"
)

func TestVec(t *testing.T) {
	res := Rad(0, 1)
	if !res.Approx(V(1, 0), 6) {
		t.Error(res, 0)
	}

	res = Rad(math.Pi/2, 1)
	if !res.Approx(V(0, 1), 6) {
		t.Error(res, 1)
	}

	if V(0, 0) != (Vec{0, 0}) {
		t.Error(V(0, 0), 2)
	}

	x, y := V(0, 0).XY()
	if x != 0 || y != 0 {
		t.Error(x, y, 3)
	}

	res = V(0.79080, 0.879879).Floor()
	if res != V(0, 0) {
		t.Error(4)
	}
}

func TestVecString(t *testing.T) {
	res := V(0, 0).String()
	if res != "V(0.000 0.000)" {
		t.Error(res)
	}
}

func TestVecSimpleOperands(t *testing.T) {
	base := V(2, 4)

	testCases := []struct {
		desc     string
		operand  func(Vec) Vec
		res, arg Vec
	}{
		{
			desc:    "add",
			operand: base.Add,
			res:     V(4, 6),
			arg:     V(2, 2),
		},
		{
			desc:    "sub",
			operand: base.Sub,
			res:     V(0, 2),
			arg:     V(2, 2),
		},
		{
			desc:    "mul",
			operand: base.Mul,
			res:     V(4, 8),
			arg:     V(2, 2),
		},
		{
			desc:    "div",
			operand: base.Div,
			res:     V(1, 2),
			arg:     V(2, 2),
		},
		{
			desc:    "to",
			operand: base.To,
			res:     V(2, 0),
			arg:     V(4, 4),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := tC.operand(tC.arg)
			if res != tC.res {
				t.Error(res, "!=", tC.res)
			}
		})
	}
}

func TestVecPointerOperands(t *testing.T) {
	base := V(2, 2)

	testCases := []struct {
		desc     string
		operand  func(Vec)
		res, arg Vec
	}{
		{
			desc:    "add",
			operand: base.AddE,
			res:     V(4, 4),
			arg:     V(2, 2),
		},
		{
			desc:    "sub",
			operand: base.SubE,
			res:     V(2, 2),
			arg:     V(2, 2),
		},
		{
			desc:    "add",
			operand: base.MulE,
			res:     V(4, 4),
			arg:     V(2, 2),
		},
		{
			desc:    "add",
			operand: base.DivE,
			res:     V(2, 2),
			arg:     V(2, 2),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.operand(tC.arg)
			if base != tC.res {
				t.Error(base, "!=", tC.res)
			}
		})
	}
}

func TestVecReducers(t *testing.T) {
	base := V(2, 2)

	testCases := []struct {
		desc    string
		operand func() float64
		res     float64
	}{
		{
			desc:    "angle",
			operand: base.Angle,
			res:     math.Pi / 4,
		},
		{
			desc:    "len",
			operand: base.Len,
			res:     math.Sqrt(8),
		},
		{
			desc:    "len2",
			operand: base.Len2,
			res:     8,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := tC.operand()
			if res != tC.res {
				t.Error(res, "!=", tC.res)
			}
		})
	}
}

func TestAngleTo(t *testing.T) {
	tests := []struct {
		a, b   Vec
		result float64
		name   string
	}{
		{
			Vec{1, 1},
			Vec{1, 1},
			0,
			"identical",
		},
		{
			Vec{-1, -1},
			Vec{1, 1},
			math.Pi,
			"opposite",
		},
		{
			Vec{1, 0},
			Vec{0, 1},
			math.Pi / 2,
			"left",
		},
		{
			Vec{-1, 0},
			Vec{0, -1},
			math.Pi / 2,
			"left inv",
		},
	}

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			res := te.a.AngleTo(te.b)
			if !Approx(res, te.result, 6) {
				t.Errorf("%v != %v", res, te.result)
			}
		})
	}
}
