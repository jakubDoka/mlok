package spatial

import (
	"testing"

	"github.com/jakubDoka/goml/core"
)

func TestNode(t *testing.T) {
	n := IntNode{}
	n.Insert(0, 0)
	n.Remove(0, 0)

	core.TestEqual(t, len(n.Groups), 0)
	core.TestEqual(t, len(n.Ints), 0)

	n.Insert(0, 0)
	n.Insert(3, 0)
	n.Insert(1, 1)
	n.Insert(2, 1)

	coll := n.Collect(1, true, nil)
	core.TestEqual(t, coll, []int{2, 1})

	coll = n.Collect(1, false, nil)
	core.TestEqual(t, coll, []int{3, 0})

	core.TestEqual(t, len(n.Groups), 2)
	core.TestEqual(t, len(n.Ints), 4)

	n.Remove(0, 0)
	n.Remove(3, 0)

	core.TestEqual(t, len(n.Groups), 1)
	core.TestEqual(t, len(n.Ints), 2)

}

func Benchmark(b *testing.B) {
	n := IntNode{}
	for i := 0; i < b.N; i++ {
		n.Insert(0, 0)
		n.Remove(0, 0)
	}
}
