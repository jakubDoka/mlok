package mat

import (
	"math"
	"testing"
)

func TestRayString(t *testing.T) {
	res := R(0, 0, 10, 10).String()
	if res != "R(O(0.000 0.000) V(10.000 10.000))" {
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
			ray:      R(5, 5, 0, -10),
			ray2:     R(0, 0, 10, 0),
			pos:      V(5, 0),
			contains: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res, contains := tC.ray.IntersectionPoint(tC.ray2)
			if tC.contains != contains {
				t.Fail()
				return
			}

			if !res.Approx(tC.pos, 8) {
				t.Error(res, tC.pos)
			}
		})
	}
}

func BenchmarkRayIntersection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		R(0, 0, 10, 0).IntersectionPoint(R(5, 5, 0, -10))
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

func TestRayIntersectY(t *testing.T) {
	y, ok := R(0, 0, 10, 0).IntersectY(R(5, 5, 0, -10))
	if !ok {
		t.Error(0)
		return
	}
	if y != 0 {
		t.Error(y, 1)
	}

	y, ok = R(0, 0, 10, 0).IntersectY(R(5, 5, 1, 0))
	if ok {
		t.Error(y, 2)
	}
}

func TestRayProjectY(t *testing.T) {
	res, _ := R(0, 0, 10, 10).ProjectY(5)
	if res != 5 {
		t.Error(res, 0)
	}

	res, ok := R(0, 0, 10, 0).ProjectY(5)
	if ok {
		t.Error(res, 0)
	}

}

func TestRaySimmetricPoint(t *testing.T) {
	testCases := []struct {
		desc      string
		pos       Vec
		ray, ray2 Ray
		contains  bool
		distance  float64
	}{
		{
			desc:     "middle cross",
			pos:      V(5, 5).Add(V(0, -math.Hypot(1, 1))),
			ray2:     R(0, 10, 10, -10),
			ray:      R(0, 0, 10, 10),
			contains: true,
			distance: 1,
		},
		{
			desc:     "colinear",
			ray2:     R(0, 0, 10, 10),
			ray:      R(0, 0, 10, 10),
			contains: false,
		},
		{
			desc:     "horizontal middle",
			ray2:     R(5, 5, 0, -10),
			ray:      R(0, 0, 10, 0),
			pos:      V(4, -1),
			contains: true,
			distance: 1,
		},
		{
			desc:     "horizontal middle reverse",
			ray:      R(5, 5, 0, -10),
			ray2:     R(0, 0, 10, 0),
			pos:      V(4, -1),
			contains: true,
			distance: 1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res, contains := tC.ray.SymmetricPoint(tC.ray2, tC.distance)
			if tC.contains != contains {
				t.Fail()
				return
			}

			if !res.Approx(tC.pos, 8) {
				t.Error(res, tC.pos)
			}
		})
	}
}
