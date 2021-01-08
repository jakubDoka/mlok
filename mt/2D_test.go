package mt

import (
	"math"
	"testing"
)

func TestAngleTo(t *testing.T) {
	tests := []struct {
		a, b   V2
		result float64
		name   string
	}{
		{
			V2{1, 1},
			V2{1, 1},
			0,
			"identical",
		},
		{
			V2{-1, -1},
			V2{1, 1},
			math.Pi,
			"opposite",
		},
		{
			V2{1, 0},
			V2{0, 1},
			math.Pi / 2,
			"left",
		},
		{
			V2{-1, 0},
			V2{0, -1},
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

func TestNMat2(t *testing.T) {
	res := NMat2(V2{10, 10}, V2{10, 10}, math.Pi/2)
	res2 := IM2.Scaled(V2{}, 10).Rotated(V2{}, math.Pi/2).Mv(V2{10, 10})
	sup := Mat2{V2{0, 10}, V2{-10, 0}, V2{10, 10}}
	if !res.Approx(sup, 6) || !res.Approx(res2, 6) {
		t.Error(res, res2)
	}
}

func BenchmarkMat2SetupSlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IM2.Scaled(V2{}, 10).Rotated(V2{}, math.Pi/2).Mv(V2{10, 10})
	}
}

func BenchmarkMat2SetupFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NMat2(V2{10, 10}, V2{10, 10}, math.Pi/2)
	}
}

func TestMat2String(t *testing.T) {
	r := IM2.String()
	if r != "Mat2(1 0 0 | 0 1 0)" {
		t.Error(r)
	}
}

func TestMat2Raw(t *testing.T) {
	r := IM2.Raw()
	if r != [...]float32{1, 0, 0, 0, 1, 0, 0, 0, 1} {
		t.Error(r)
	}
}

func TestMat2Projection(t *testing.T) {
	m := NMat2(V2{1, 1}, V2{1, 2}, math.Pi)
	r := m.Project(V2{1, 1})
	if !r.Approx(V2{0, -1}, 6) {
		t.Error(r)
	}
	r = m.Unproject(r)
	if !r.Approx(V2{1, 1}, 6) {
		t.Error(r)
	}
}
