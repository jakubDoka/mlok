package spatial

import (
	"fmt"
	"strings"

	"github.com/jakubDoka/mlok/mat"
)

type QuadTree struct {
	Bounds            mat.AABB
	Nodes             []QuadNode
	NodeCap, DepthCap int
}

func (t *QuadTree) Query(group int, include bool, area mat.AABB, helpers ...[]int) (coll, frontier, temp []int) {
	if len(t.Nodes) == 0 {
		return
	}

	switch len(helpers) {
	case 3:
		coll = helpers[0][:0]
		fallthrough
	case 2:
		frontier = helpers[1][:0]
		fallthrough
	case 1:
		temp = helpers[2][:0]
	}

	frontier = append(frontier[:0], 0)
	for len(frontier) != 0 {
		for _, i := range frontier {
			n := &t.Nodes[i]
			if !n.Intersects(area) {
				continue
			}
			coll = n.Collect(group, include, coll)
			if !n.Branch {
				continue
			}
			ptr := n.Children
			temp = append(temp, ptr, ptr+1, ptr+2, ptr+3)
		}
		frontier, temp = temp, frontier
		temp = temp[:0]
	}

	return
}

func (t *QuadTree) Insert(address *int, bounds mat.AABB, id, group int) bool {
	if len(t.Nodes) == 0 {
		t.Nodes = append(t.Nodes, QuadNode{AABB: t.Bounds})
	}

	n := &t.Nodes[0]
	for n.Branch {
		var next *QuadNode
		for i, e := slice(n.Children); i < e; i++ {
			n := &t.Nodes[i]
			if !n.Closed && bounds.Fits(n.AABB) {
				next = n
				break
			}
		}

		if next == nil {
			break
		}
		n = next
	}

	n.Insert(id, group)
	*address = n.Self

	if t.ShouldSplit(n) {
		t.Split(n.Self)
	}

	return false
}

func (t *QuadTree) Update(address *int, bounds mat.AABB, id, group int) {
	n := &t.Nodes[*address]
	if !bounds.Fits(n.AABB) {
		n.Remove(id, group)
		if len(n.Ints) == 0 && t.Count(n.Self) == 0 {
			n.Branch = false
		}
		t.Insert(address, bounds, id, group)
	} else if n.Closed {
		n.Remove(id, group)
		if len(n.Ints) == 0 && t.Count(n.Self) == 0 {
			n.Branch = false
		}
		for n.Closed && n.Self != 0 {
			n = &t.Nodes[n.Parent]
		}
		n.Insert(id, group)
		*address = n.Self
	} else {
		if n.Branch {
			for i, e := slice(n.Children); i < e; i++ {
				if bounds.Fits(t.Nodes[i].AABB) {
					n.Remove(id, group)
					t.Nodes[i].Insert(id, group)
					*address = i
					if t.ShouldSplit(n) {
						t.Split(i)
					}
				}
			}
		}
	}

}

func (t *QuadTree) ShouldSplit(n *QuadNode) bool {
	return len(n.Ints) >= t.NodeCap && n.Level < t.DepthCap || t.DepthCap == 0
}

func (t *QuadTree) Remove(address, id, group int) bool {
	ok := t.Nodes[address].Remove(id, group)
	if ok && t.Count(address) < t.NodeCap {
		t.Close(address)
	}
	return ok
}

func (t *QuadTree) Count(node int) (count int) {
	n := &t.Nodes[node]
	if n.Branch {
		ptr := n.Children
		count += t.Count(ptr) + t.Count(ptr+1) + t.Count(ptr+2) + t.Count(ptr+3)
	}
	return count + len(n.Ints)
}

func (t *QuadTree) Close(node int) {
	n := &t.Nodes[node]

	if n.Branch {
		ptr := n.Children
		t.Close(ptr)
		t.Close(ptr + 1)
		t.Close(ptr + 2)
		t.Close(ptr + 3)
	}
}

func (t *QuadTree) Split(node int) {
	n := &t.Nodes[node]
	n.Branch = true
	if n.Children != 0 {
		return
	}

	n.Children = len(t.Nodes)
	level := n.Level + 1
	cet := n.Center()

	t.Nodes = append(t.Nodes,
		QuadNode{
			AABB: mat.AABB{
				Min: n.Min,
				Max: cet,
			},
			Level:  level,
			Self:   n.Children,
			Parent: n.Self,
		},
		QuadNode{
			AABB: mat.AABB{
				Min: mat.V(cet.X, n.Min.Y),
				Max: mat.V(n.Max.X, cet.Y),
			},
			Level:  level,
			Self:   n.Children + 1,
			Parent: n.Self,
		},
		QuadNode{
			AABB: mat.AABB{
				Min: cet,
				Max: n.Max,
			},
			Level:  level,
			Self:   n.Children + 2,
			Parent: n.Self,
		},
		QuadNode{
			AABB: mat.AABB{
				Min: mat.V(n.Min.X, cet.Y),
				Max: mat.V(cet.X, n.Max.Y),
			},
			Level:  level,
			Self:   n.Children + 3,
			Parent: n.Self,
		},
	)
}

// FormatDebug makes a readable formatting of tree structure
func (t *QuadTree) Debug(depth, node int) string {
	depth++
	n := &t.Nodes[node]
	if !n.Branch {
		return fmt.Sprintf("%s%d", strings.Repeat("  ", depth), len(n.Ints))
	}
	ptr := n.Children
	return fmt.Sprintf(
		"%s%d\n%s\n%s\n%s\n%s\n",
		strings.Repeat("  ", depth),
		len(n.Ints),
		t.Debug(depth, ptr),
		t.Debug(depth, ptr+1),
		t.Debug(depth, ptr+2),
		t.Debug(depth, ptr+3),
	)
}

func slice(start int) (s, e int) {
	return start, start + 4
}

type QuadNode struct {
	IntNode
	mat.AABB

	Children, Level, Parent, Self int
	Branch, Closed                bool
}
