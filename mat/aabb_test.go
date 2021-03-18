package mat

import "testing"

func TestAABBIntersectCircle(t *testing.T) {
	b := A(0, 0, 10, 10)
	testCases := []struct {
		desc string
		a    AABB
		c    Circ
		i    bool
	}{
		{"inside", b, C(5, 5, 100), true},

		{"right", b, C(15, 5, 10), true},
		{"not right", b, C(15, 5, 3), false},

		{"left", b, C(-5, 5, 10), true},
		{"not left", b, C(-5, 5, 3), false},

		{"top", b, C(5, 15, 10), true},
		{"not top", b, C(5, 15, 3), false},

		{"bottom", b, C(5, -5, 10), true},
		{"not bottom", b, C(5, -5, 3), false},

		{"bottom left", b, C(-5, -5, 10), true},
		{"not bottom left", b, C(-5, -5, 3), false},

		{"top left", b, C(-5, 15, 10), true},
		{"not top left", b, C(-5, 15, 3), false},

		{"bottom right", b, C(15, -5, 10), true},
		{"not bottom right", b, C(15, -5, 3), false},

		{"top right", b, C(15, 15, 10), true},
		{"not top right", b, C(15, 15, 3), false},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			re := tC.a.IntersectCircle(tC.c)
			if re != tC.i {
				t.Error("fail")
			}
		})
	}
}

func BenchmarkRectIntersectCircle(b *testing.B) {
	a, c := A(0, 0, 10, 10), C(-5, -5, 10)
	for i := 0; i < b.N; i++ {
		a.IntersectCircle(c)
	}
}
