package memory

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/jakubDoka/goml/core"
	"github.com/jakubDoka/mlok/logic/memory/gen"
)

func TestSortVec(t *testing.T) {
	testCases := []struct {
		desc string
		inp  gen.Int32Vec
	}{
		{
			desc: "reverse",
			inp:  gen.Int32Vec{0, 1, 2, 3, 4, 5},
		},
		{
			desc: "nothing",
			inp:  gen.Int32Vec{5, 4, 3, 2, 1, 0},
		},
		{
			desc: "random",
			inp:  gen.Int32Vec{3, 2, 4, 0, 1, 5},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.inp.Sort(nil, func(a, b int32) bool { return a > b })
			core.TestEqual(t, tC.inp, gen.Int32Vec{5, 4, 3, 2, 1, 0})
		})
	}
}

func BenchmarkSortVec(b *testing.B) {
	v := make(gen.Int32Vec, 1000)
	for i := range v {
		v[i] = rand.Int31()
	}
	c := make(gen.Int32Vec, 1000)

	buff := make([]int, 1000)
	for i := 0; i < b.N; i++ {
		buff = c.Sort(buff[:0], func(a, b int32) bool { return a > b })
		copy(c, v)
	}
}

type SortBy []int32

func (a SortBy) Len() int           { return len(a) }
func (a SortBy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortBy) Less(i, j int) bool { return a[i] < a[j] }

func BenchmarkSort(b *testing.B) {
	v := make(gen.Int32Vec, 1000)
	for i := range v {
		v[i] = rand.Int31()
	}
	c := make(gen.Int32Vec, 1000)
	for i := 0; i < b.N; i++ {
		sort.Sort(SortBy(c))
		copy(c, v)
	}
}

func TestByInsertVec(t *testing.T) {
	v := make(gen.Int32Vec, 1000)
	for i := range v {
		v[i] = int32(rand.Intn(1000))
	}

	b := make(gen.Int32Vec, 0, 1000)
	for _, v := range v {
		b.BiInsert(v, gen.Int32BiComp)
	}
	v.Sort(nil, func(a, b int32) bool { return a < b })

	core.TestEqual(t, b, v)
}

func BenchmarkInsetVec(b *testing.B) {
	v := make(gen.Int32Vec, 1000)
	for i := 0; i < b.N; i++ {
		v.Insert(0, 1)
		v = v[:len(v)-1]
	}
}
