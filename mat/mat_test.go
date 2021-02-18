package mat

import (
	"math"
	"testing"
)

func TestNMat(t *testing.T) {
	res := M(Vec{10, 10}, Vec{10, 10}, math.Pi/2)
	res2 := IM.Scaled(Vec{}, 10).Rotated(Vec{}, math.Pi/2).Mv(Vec{10, 10})
	sup := Mat{Vec{0, 10}, Vec{-10, 0}, Vec{10, 10}}
	if !res.Approx(sup, 6) || !res.Approx(res2, 6) {
		t.Error(res, res2)
	}
}

func BenchmarkMatSetupSlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IM.Scaled(Vec{}, 10).Rotated(Vec{}, math.Pi/2).Mv(Vec{10, 10})
	}
}

func BenchmarkMatSetupFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		M(Vec{10, 10}, Vec{10, 10}, math.Pi/2)
	}
}

func TestMatString(t *testing.T) {
	r := IM.Mv(V(0.28374, 0.3972)).String()
	if r != "Mat(1.000 0.000 0.284 | 0.000 1.000 0.397)" {
		t.Error(r)
	}
}

func TestMatRaw(t *testing.T) {
	r := IM.Raw()
	if r != [...]float32{1, 0, 0, 0, 1, 0, 0, 0, 1} {
		t.Error(r)
	}
}

func TestMatProjection(t *testing.T) {
	m := M(Vec{1, 1}, Vec{1, 2}, math.Pi)
	r := m.Project(Vec{1, 1})
	if !r.Approx(Vec{0, -1}, 6) {
		t.Error(r)
	}
	r = m.Unproject(r)
	if !r.Approx(Vec{1, 1}, 6) {
		t.Error(r)
	}
}

func TestTran(t *testing.T) {
	tran := Tran{Vec{10, 10}, Vec{10, 10}, math.Pi / 2}
	sup := Mat{Vec{0, 10}, Vec{-10, 0}, Vec{10, 10}}
	if !tran.Mat().Approx(sup, 6) {
		t.Fail()
	}
}
