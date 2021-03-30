package spatial

import (
	"reflect"
	"testing"

	"github.com/jakubDoka/mlok/mat"
)

func TestMinHahs(t *testing.T) {
	coll := []int{}
	h := NMinHash(4, 4, mat.V(10, 10))
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
			h.Query(tC.area, coll, 0, true)
			if !reflect.DeepEqual(tC.res, coll) {
				t.Error(coll, tC.res)
				return
			}
		})
	}
}

func BenchmarkHasherReinsert(b *testing.B) {
	h := NMinHash(4, 4, mat.V(10, 10))
	adr := mat.Point{}
	for i := 0; i < b.N; i++ {
		h.Insert(&adr, mat.ZV, 0, 0)
		h.Remove(adr, 0, 0)
	}
}

func BenchmarkHasherUpdate(b *testing.B) {
	h := NMinHash(4, 4, mat.V(10, 10))
	adr := mat.Point{}
	h.Insert(&adr, mat.ZV, 0, 0)
	for i := 0; i < b.N; i++ {
		h.Update(&adr, mat.ZV, 0, 0)
	}
}

func BenchmarkHasherQuery(b *testing.B) {
	h := NMinHash(4, 4, mat.V(10, 10))
	adr := mat.Point{}
	h.Insert(&adr, mat.ZV, 0, 0)
	var buff []int
	for i := 0; i < b.N; i++ {
		buff = h.Query(mat.ZA, buff[:0], 0, true)
	}
}
