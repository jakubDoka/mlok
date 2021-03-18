package mat

import (
	"math"
	"testing"
)

func TestRayString(t *testing.T) {
	res := R(0, 0, 10, 10).String()
	if res != "R(0 0 10 10)" {
		t.Error(res)
	}
}

func TestRayContains(t *testing.T) {
	testCases := []struct {
		desc     string
		pos      Vec
		ray      Ray
		contains bool
	}{
		{
			desc:     "middle",
			pos:      V(5, 5),
			ray:      R(0, 0, 10, 10),
			contains: true,
		},
		{
			desc:     "in front",
			pos:      V(11, 11),
			ray:      R(0, 0, 10, 10),
			contains: false,
		},
		{
			desc:     "horizontal middle",
			pos:      V(5, 0),
			ray:      R(0, 0, 10, 0),
			contains: true,
		},
		{
			desc:     "vertical middle",
			pos:      V(0, 5),
			ray:      R(0, 0, 0, 10),
			contains: true,
		},
		{
			desc:     "horizontal in front",
			pos:      V(11, 0),
			ray:      R(0, 0, 10, 0),
			contains: false,
		},
		{
			desc:     "horizontal next to",
			pos:      V(5, 1),
			ray:      R(0, 0, 10, 0),
			contains: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := tC.ray.Contains(tC.pos)
			if tC.contains != res {
				t.Fail()
			}
		})
	}
}

func TestRayIntersection(t *testing.T) {
	testCases := []struct {
		desc      string
		pos       Vec
		ray, ray2 Ray
		contains  bool
	}{
		{
			desc:     "middle cross",
			pos:      V(5, 5),
			ray2:     R(0, 10, 10, -10),
			ray:      R(0, 0, 10, 10),
			contains: true,
		},
		{
			desc:     "in front",
			ray2:     R(0, 10, 0, 1),
			ray:      R(0, 0, 10, 10),
			contains: false,
		},
		{
			desc:     "colinear",
			ray2:     R(0, 10, 0, 1),
			ray:      R(0, 0, 0, 10),
			contains: false,
		},
		{
			desc:     "horizontal middle",
			ray2:     R(5, 5, 0, -10),
			ray:      R(0, 0, 10, 0),
			pos:      V(5, 0),
			contains: true,
		},
		{
			desc:     "horizontal middle reverse",
			ray:      R(5, -5, 0, 10),
			ray2:     R(0, 0, 10, 0),
			pos:      V(5, 0),
			contains: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res, contains := tC.ray.Intersect(tC.ray2)
			if tC.contains != contains {
				t.Error(res, contains)
				return
			}

			if !res.Approx(tC.pos, 8) {
				t.Error(res, tC.pos)
			}
		})
	}
}

func TestIntersectCircle(t *testing.T) {
	testCases := []struct {
		desc string
		r    Ray
		c    Circ
		a, b Vec
	}{
		{
			"touch",
			R(10, 10, 10, 0),
			C(0, 0, 10),
			V(0, 10), V(0, 10),
		},
		{
			"touch",
			R(10, 10, -10, 0),
			C(0, 0, 10),
			V(0, 10), V(0, 10),
		},
		{
			"touch",
			R(10, 10, 0, 10),
			C(0, 0, 10),
			V(10, 0), V(10, 0),
		},
		{
			"touch",
			R(10, 10, 0, -10),
			C(0, 0, 10),
			V(10, 0), V(10, 0),
		},
		{
			"intersect",
			R(0, 0, 0, -10),
			C(0, 0, 10),
			V(0, -10), V(0, 10),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			a, b := tC.r.LineIntersectCircle(tC.c)
			if !a.Approx(tC.a, 8) || !b.Approx(tC.b, 8) {
				t.Error(a, b)
			}
		})
	}
}

func BenchmarkRayIntersection(b *testing.B) {
	a, c := R(0, 0, 10, 0), R(5, 5, 0, -10)
	for i := 0; i < b.N; i++ {
		a.Intersect(c)
	}
}

func TestRayCollinear(t *testing.T) {
	if !R(0, 0, 10, 4).Colinear(R(0, 0, 5, 2)) {
		t.Error(0)
	}
	if R(0, 0, 10, 4).Colinear(R(0, 0, 5, 3)) {
		t.Error(1)
	}
}

func TestIntersectX(t *testing.T) {
	testCases := []struct {
		desc string
		a, b Ray
		x    Vec
	}{
		{
			"colinear",
			R(0, 0, 10, 0),
			R(0, 1, 1, 0),
			V(math.NaN(), math.NaN()),
		},

		{
			"intersecting",
			R(0, 0, 10, 0),
			R(0, 10, 1, -1),
			V(10, 0),
		},
		{
			"intersecting",
			R(0, 0, 10, 0),
			R(0, 10, 1, -2),
			V(5, 0),
		},
		{
			"intersecting",
			R(0, 0, 10, 0),
			R(0, 10, -4, -2),
			V(-20, 0),
		},
		{
			"intersecting",
			R(0, 0, 1, 1),
			R(0, 10, 1, -1),
			V(5, 5),
		},
		{
			"intersecting",
			R(1, 0, 1, 1),
			R(0, 10, 1, -1),
			V(5.5, 4.5),
		},
		{
			"intersecting",
			R(10, 0, 1, 1),
			R(0, 10, 1, -1),
			V(10, 0),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			x := tC.a.LineIntersect(tC.b)
			if !x.Approx(tC.x, 8) {
				t.Error(x)
			}

		})
	}
}

func TestRayProjectY(t *testing.T) {
	res := R(0, 0, 10, 10).ProjectY(5)
	if res != 5 {
		t.Error(res, 0)
	}

	res = R(0, 0, 10, 0).ProjectY(5)
	if !math.IsInf(res, 1) && !math.IsInf(res, -1) {
		t.Error(res, 0)
	}

}
