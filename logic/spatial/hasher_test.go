package spatial

import (
	"gobatch/mat"
	"reflect"
	"testing"
)

func TestMinHahs(t *testing.T) {
	coll := []int{}
	h := NMinHahs(4, 4, mat.V(10, 10))
	h.Insert(&mat.Point{}, mat.V(1, 1), 0, 0)
	h.Insert(&mat.Point{}, mat.V(40, 40), 1, 0)
	h.Insert(&mat.Point{}, mat.V(20, 20), 2, 0)
	testCases := []struct {
		desc string
		res  []int
		area mat.AABB
	}{
		{"corner", []int{0, 2}, mat.A(0, 0, 15, 15)},
		{"all", []int{0, 2, 1}, mat.A(-20, -20, 100, 100)},
		{"nothing", []int{0}, mat.A(-100, -100, -100, -100)},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			coll = coll[:0]
			h.Query(tC.area, &coll, 0, true)
			if !reflect.DeepEqual(tC.res, coll) {
				t.Error(coll, tC.res)
				return
			}
		})
	}
}

func BenchmarkHasher(b *testing.B) {
	h := NMinHahs(4, 4, mat.V(10, 10))
	adr := mat.Point{}
	for i := 0; i < b.N; i++ {
		h.Insert(&adr, mat.ZV, 0, 0)
		h.Remove(&adr, 0, 0)
	}
}
