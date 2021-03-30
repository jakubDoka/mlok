package spatial

import (
	"fmt"
	"strings"

	"github.com/jakubDoka/mlok/mat"
)

// Tree is a quadtree that prefers updating over reinseting. If you
// have to detect collizions between big amount of objects with variety
// of shape sizes this is a good pick, use hasher instead if you are using
// lot of same sized objects
//
// Tree has to be initialized with NTree as Tree.Globals is pointer
//
// Tree is not fully tested feature
type Tree struct {
	*Globals
	TNode

	Bounds          mat.AABB
	Ch              []Tree
	Branch, NotRoot bool
}

// Globals stored data shared by all tree nodes
type Globals struct {
	tmp, tmp2 []TreeEntity
	Cap       int
}

// NTree returns ready to use tree
func NTree(cap int, bounds mat.AABB) *Tree {
	return &Tree{
		Globals: &Globals{
			Cap: cap,
		},
		Bounds: bounds,
	}
}

// Query returns all entities that can intersect with area, if include is false only entities with
// different group are returned, othervise entities with sam group are returned, using QueryAll if
// you don't need all groups and filtering them your self is slower as Quadterr si structured to
// optimize the process
func (t *Tree) Query(group int, include bool, coll []TreeEntity, area mat.AABB) []TreeEntity {
	t.Collect(group, include, coll)
	for i := range t.Ch {
		if t.Bounds.Intersects(area) {
			coll = t.Ch[i].Query(group, include, coll, area)
		}
	}
	return coll
}

// QueryAll returns all Entities that area can intersect with
func (t *Tree) QueryAll(coll *[]TreeEntity, area mat.AABB) {
	t.CollectAll(coll)
	for i := range t.Ch {
		if t.Bounds.Intersects(area) {
			t.Ch[i].QueryAll(coll, area)
		}
	}
}

// Insert inserts entity to tree, returns false if entity cannot
// be inserted
func (t *Tree) Insert(te TreeEntity) bool {
	if t.Branch {
		for i := range t.Ch {
			if t.Ch[i].Insert(te) {
				return true
			}
		}
	} else {
		if t.Count+1 > t.Cap {
			t.tmp = t.tmp[:0]
			t.Bail()
			t.Split()
			for _, e := range t.tmp {
				if !t.Insert(e) { // does not fit to any quadrant
					t.TNode.Insert(e, e.Group())
				}
			}
			if t.Insert(te) {
				return true
			}
		}
	}

	if te.Bounds().Fits(t.Bounds) || !t.NotRoot {
		t.TNode.Insert(te, te.Group())
		return true
	}
	return false
}

// Remove removes the entity, if there is no such entity, false is returned
func (t *Tree) Remove(te TreeEntity) bool {
	if t.Branch {
		for i := range t.Ch {
			if te.Bounds().Fits(t.Ch[i].Bounds) {
				return t.Ch[i].Remove(te)
			}
		}
	}
	return t.TNode.Remove(te, te.Group())
}

// Update does two things, it moves
func (t *Tree) Update() (count int) {
	if t.Branch {
		for i := range t.Ch {
			count += t.Ch[i].Update()
		}
		if count < t.Cap {
			t.tmp = t.tmp[:0]
			t.Branch = false
			for i := range t.Ch {
				t.Ch[i].Bail()
			}
			for _, e := range t.tmp {
				t.TNode.Insert(e, e.Group())
			}
			count += t.Cap // this will keep update time stable
		}
	}

	var j int
	for _, e := range t.tmp2 {
		if e.Dead() {
			continue
		}
		if !e.Bounds().Fits(t.Bounds) && t.NotRoot {
			t.tmp2[j] = e
			j++
			continue
		}
		t.TNode.Insert(e, e.Group())
	}
	for k := j; k < len(t.tmp2); k++ {
		t.tmp2[k] = nil
	}
	t.tmp2 = t.tmp2[:j]

	t.Count = 0
	for i := range t.Sets {
		ids := t.Sets[i].IDs
		var j int
		for _, e := range ids {
			if e.Dead() {
				continue
			}
			if e.Bounds().Fits(t.Bounds) {
				ids[j] = e
				j++
				t.Count++
				continue
			}
			t.tmp2 = append(t.tmp2, e)
		}
		for k := j; k < len(ids); k++ {
			ids[k] = nil
		}
		t.Sets[i].IDs = ids[:j]
	}

	return count + t.Count
}

// TotalCount returns total count of children in Tree
func (t *Tree) TotalCount() (total int) {
	if t.Branch {
		for i := range t.Ch {
			total += t.Ch[i].TotalCount()
		}
	}

	return total + t.Count
}

// Split allocates Quadtree children
func (t *Tree) Split() {
	t.Branch = true
	if len(t.Ch) == 0 {
		t.Ch = make([]Tree, 4)

		cet := t.Bounds.Center()
		t.Ch[0].Bounds = mat.AABB{
			Min: t.Bounds.Min,
			Max: cet,
		}
		t.Ch[1].Bounds = mat.AABB{
			Min: mat.V(cet.X, t.Bounds.Min.Y),
			Max: mat.V(t.Bounds.Max.X, cet.Y),
		}
		t.Ch[2].Bounds = mat.AABB{
			Min: cet,
			Max: t.Bounds.Max,
		}
		t.Ch[3].Bounds = mat.AABB{
			Min: mat.V(t.Bounds.Min.X, cet.Y),
			Max: mat.V(cet.X, t.Bounds.Max.Y),
		}

		for i := range t.Ch {
			t.Ch[i].Globals = t.Globals
			t.Ch[i].NotRoot = true
		}
	}
}

// Bail retrives all objects from TNode into coll
func (t *Tree) Bail() {
	t.Count = 0
	for i := range t.Sets {
		ids := t.Sets[i].IDs
		t.tmp = append(t.tmp, ids...)
		for j := range ids {
			ids[j] = nil
		}
		t.Sets[i].IDs = ids[:0]
	}
}

// FormatDebug makes a readable formatting of tree structure
func (t *Tree) FormatDebug(depth int) string {
	depth++
	if !t.Branch {
		return fmt.Sprintf("%s%d", strings.Repeat("  ", depth), t.Count)
	}
	return fmt.Sprintf(
		"%s%d\n%s\n%s\n%s\n%s\n",
		strings.Repeat("  ", depth),
		t.Count,
		t.Ch[0].FormatDebug(depth),
		t.Ch[1].FormatDebug(depth),
		t.Ch[2].FormatDebug(depth),
		t.Ch[3].FormatDebug(depth),
	)
}

// TreeEntity represents insertabel data for Tree
type TreeEntity interface {
	// Bounds returns bounding rectangle
	Bounds() mat.AABB
	// Group returns a entity group
	Group() int
	// Dead returns whether entity should be removed
	Dead() bool
}
