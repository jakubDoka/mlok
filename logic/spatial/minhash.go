package spatial

import (
	"github.com/jakubDoka/mlok/mat"
)

// MinHash is for efficient spatial hasher, tuned to perfection
// is able to handle thousands of entities for cost of using ids and limited
// space, if you will be dealing with concentrated entity use this over
type MinHash struct {
	nodeSize mat.Vec
	w, h     int
	Nodes    []IntNode
	buff     []int
}

// NMinHash is MinHash constructor
func NMinHash(w, h int, tileSize mat.Vec) MinHash {
	return MinHash{
		nodeSize: mat.V(1, 1).Div(tileSize),
		w:        w,
		h:        h,
		Nodes:    make([]IntNode, h*w),
	}
}

// TileSize ...
func (h *MinHash) TileSize() mat.Vec {
	return mat.V(1, 1).Div(h.nodeSize)
}

// Insert adds shape to MinHash
func (h *MinHash) Insert(adr *mat.Point, pos mat.Vec, id, group int) {
	*adr = h.Adr(pos)
	h.Nodes[prj(adr.X, adr.Y, h.w)].Insert(id, group)
}

// Remove removes shape from MinHash. If operation fails, false is returned
func (h *MinHash) Remove(adr mat.Point, id, group int) bool {
	return h.Nodes[prj(adr.X, adr.Y, h.w)].Remove(id, group)
}

// Update updates state of object if it changed quadrant, if operation fails, false is returned
func (h *MinHash) Update(old *mat.Point, pos mat.Vec, id, group int) bool {
	p := h.Adr(pos)
	if *old == p {
		return true
	}

	if h.Nodes[prj(old.X, old.Y, h.w)].Remove(id, group) {
		h.Nodes[prj(p.X, p.Y, h.w)].Insert(id, group)
		*old = p
		return true
	}

	return false
}

// Query returns colliding shapes with given rect
func (h *MinHash) Query(rect mat.AABB, coll []int, group int, including bool) []int {
	max := h.Adr(rect.Max).Add(mat.P(2, 2)).Min(mat.P(h.w, h.h))
	min := h.Adr(rect.Min).Sub(mat.P(1, 1)).Max(mat.P(0, 0))

	for y := min.Y; y < max.Y; y++ {
		for x := min.X; x < max.X; x++ {
			n := &h.Nodes[prj(x, y, h.w)]
			if len(n.Ints) != 0 {
				coll = n.Collect(group, including, coll)
			}
		}
	}

	return coll
}

func prj(x, y, stride int) int {
	return y*stride + x
}

// Adr returns node, position belongs to
func (h *MinHash) Adr(pos mat.Vec) mat.Point {
	// we want this inlined
	x, y := int(pos.X*h.nodeSize.X), int(pos.Y*h.nodeSize.Y)
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= h.w {
		x = h.w - 1
	}
	if y >= h.h {
		y = h.h - 1
	}
	return mat.P(x, y)
}
