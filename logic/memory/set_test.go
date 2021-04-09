package memory

import (
	"testing"

	"github.com/jakubDoka/goml/core"
)

func TestSet(t *testing.T) {
	testCases := []struct {
		desc     string
		inp, out ElemSet
		insert   bool
	}{
		{
			desc: "simple",
			inp:  ElemSet{5},
			out:  ElemSet{},
		},
		{
			desc:   "simple insert",
			inp:    ElemSet{},
			out:    ElemSet{5},
			insert: true,
		},
		{
			desc: "no remove",
			inp:  ElemSet{0, 1, 4, 7, 8},
			out:  ElemSet{0, 1, 4, 7, 8},
		},
		{
			desc:   "no insert",
			inp:    ElemSet{0, 1, 5, 7, 8},
			out:    ElemSet{0, 1, 5, 7, 8},
			insert: true,
		},
		{
			desc:   "insert",
			inp:    ElemSet{0, 1, 7, 8},
			out:    ElemSet{0, 1, 5, 7, 8},
			insert: true,
		},
		{
			desc: "remove",
			inp:  ElemSet{0, 1, 5, 7, 8},
			out:  ElemSet{0, 1, 7, 8},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.insert {
				tC.inp.Insert(5)
			} else {
				tC.inp.Remove(5)
			}

			core.TestEqual(t, tC.inp, tC.out)
		})
	}
}
