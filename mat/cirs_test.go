package mat

import (
	"testing"
)

func TestCircUnion(t *testing.T) {
	testCases := []struct {
		desc    string
		a, b, c Circ
	}{
		{
			desc: "symmetirc",
			a:    C(1, 0, 1),
			b:    C(-1, 0, 1),
			c:    C(0, 0, 4),
		},
		{
			desc: "touching",
			a:    C(2, 0, 2),
			b:    C(-1, 0, 1),
			c:    C(1.5, 0, 6),
		},
		{
			desc: "apart",
			a:    C(6, 0, 2),
			b:    C(-1, 0, 1),
			c:    C(3.5, 0, 10),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := tC.a.Union(tC.b)
			if !res.Approx(tC.c, 7) {
				t.Error(res)
			}
		})
	}
}

func TestCircPrj(t *testing.T) {
	c := C(5, 5, 10)
	testCases := []struct {
		desc   string
		val    Vec
		rx, ry [2]float64
	}{
		{
			desc: "on center",
			val:  V(5, 5),
			rx:   [2]float64{15, -5},
			ry:   [2]float64{15, -5},
		},
		{
			desc: "touching",
			val:  V(15, 15),
			rx:   [2]float64{5, 5},
			ry:   [2]float64{5, 5},
		},
		{
			desc: "not",
			val:  V(16, 16),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := c.PrjX(tC.val.X)
			if res != tC.rx {
				t.Error(res, tC.rx, "x")
			}

			res = c.PrjY(tC.val.Y)
			if res != tC.ry {
				t.Error(res, tC.ry, "y")
			}
		})
	}
}

func TestCircIntersect(t *testing.T) {
	testCases := []struct {
		desc       string
		a, b       Circ
		intersects bool
	}{
		{
			desc:       "intersects",
			a:          C(0, 0, 3),
			b:          C(0, 10, 8),
			intersects: true,
		},
		{
			desc:       "touch",
			a:          C(0, 0, 3),
			b:          C(0, 10, 7),
			intersects: true,
		},
		{
			desc: "no",
			a:    C(0, 0, 3),
			b:    C(0, 10, 3),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := tC.a.Intersects(tC.b)
			if res != tC.intersects {
				t.Fail()
			}
		})
	}
}

func BenchmarkCircSimplePrj(b *testing.B) {
	c := C(0, 0, 10)
	for i := 0; i < b.N; i++ {
		c.SimplePrjX(1)
	}
}

func BenchmarkCircPrj(b *testing.B) {
	c := C(0, 0, 10)
	for i := 0; i < b.N; i++ {
		c.PrjX(1)
	}
}
