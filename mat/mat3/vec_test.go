package mat3

import (
	"math"
	"testing"
)

func BenchmarkCross(b *testing.B) {
	v := Vec{1, 0, 6}
	u := Vec{4, 8, 9}
	for i := 0; i < b.N; i++ {
		v.Cross(u)
	}
}

func TestVecRotated(t *testing.T) {
	testCases := []struct {
		desc            string
		inp, out, pivot Vec
		ang             float64
	}{
		{
			desc:  "around x",
			inp:   V(0, 1, 0),
			out:   V(0, 0, 1),
			pivot: V(1, 0, 0),
			ang:   math.Pi / 2,
		},
		{
			desc:  "around x tilted",
			inp:   V(0, 1, 1),
			out:   V(0, -1, 1),
			pivot: V(1, 0, 0),
			ang:   math.Pi / 2,
		},
		{
			desc:  "around y",
			inp:   V(1, 0, 0),
			out:   V(0, 0, -1),
			pivot: V(0, 1, 0),
			ang:   math.Pi / 2,
		},
		{
			desc:  "around z",
			inp:   V(1, 0, 0),
			out:   V(0, 1, 0),
			pivot: V(0, 0, 1),
			ang:   math.Pi / 2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := tC.inp.Rotated(tC.ang, tC.pivot)
			if !tC.out.Approx(res, 8) {
				t.Error(res)
			}
		})
	}
}
