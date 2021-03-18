package mat

import (
	"testing"
)

func BenchmarkIntersect(b *testing.B) {
	a, c := C(10, 10, 10), C(0, 0, 10)
	for i := 0; i < b.N; i++ {
		a.Intersect(c)
	}
}

func TestCircleIntersectY(t *testing.T) {
	testCases := []struct {
		desc string
		a, b Circ
		c, d Vec
	}{
		{
			"no intersection",
			C(10, 10, 10),
			C(30, 10, 7),
			ZV, ZV,
		},
		{
			"touch",
			C(0, 10, 10),
			C(20, 10, 10),
			V(10, 10), V(10, 10),
		},
		{
			"touch",
			C(10, 0, 10),
			C(10, 20, 10),
			V(10, 10), V(10, 10),
		},
		{
			"intersection",
			C(10, 0, 10),
			C(10, 20, 10.1),
			V(8.998763295968432, 9.94975), V(11.001236704031568, 9.94975),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			a, b := tC.a.Intersect(tC.b)
			if !a.Approx(tC.c, 8) || !b.Approx(tC.d, 8) {
				t.Error(a, b)
			}
		})
	}
}

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
		desc           string
		val            Vec
		r1, r2, r3, r4 float64
	}{
		{
			desc: "on center",
			val:  V(5, 5),
			r1:   -5, r2: 15,
			r3: -5, r4: 15,
		},
		{
			desc: "touching",
			val:  V(15, 15),
			r1:   5, r2: 5,
			r3: 5, r4: 5,
		},
		{
			desc: "not",
			val:  V(16, 16),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			a, b := c.ProjectX(tC.val.X)
			if !Approx(a, tC.r1, 8) || !Approx(b, tC.r2, 8) {
				t.Error(a, b, "x")
			}

			a, b = c.ProjectY(tC.val.Y)
			if !Approx(a, tC.r3, 8) || !Approx(b, tC.r4, 8) {
				t.Error(a, b, "y")
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
		c.ProjectX(1)
	}
}
