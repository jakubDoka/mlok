package spatial

import (
	"testing"

	"github.com/jakubDoka/mlok/mat"
)

func TestTree(t *testing.T) {
	tree := NTree(1, mat.A(0, 0, 1000, 1000))
	d := dummy{
		pos:  mat.V(200, 100),
		size: 10,
	}
	for i := 0; i < 10; i++ {
		tree.Insert(&d)
	}
	t.Error(tree.TotalCount(), "\n", tree.FormatDebug(-1))
	for i := 0; i < 10; i++ {
		tree.Remove(&d)
	}
	t.Error(tree.TotalCount(), "\n", tree.FormatDebug(-1))
	tree.Update()
	t.Error(tree.TotalCount(), "\n", tree.FormatDebug(-1))
	for i := 0; i < 10; i++ {
		tree.Insert(&d)
	}
	d.dead = true
	tree.Update()
	tree.Update()
	tree.Update()
	tree.Update()
	tree.Update()
	tree.Update()
	tree.Update()
	tree.Update()
	tree.Update()
	tree.Update()
	t.Error(tree.TotalCount(), "\n", tree.FormatDebug(-1))

}

type dummy struct {
	pos   mat.Vec
	size  float64
	group int
	dead  bool
}

func (d *dummy) Group() int {
	return d.group
}

func (d *dummy) Dead() bool {
	return d.dead
}

func (d *dummy) Pos() mat.Vec {
	return d.pos
}

func (d *dummy) Bounds() mat.AABB {
	return mat.Square(d.pos, d.size)
}
