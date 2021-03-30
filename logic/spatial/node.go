package spatial

import (
	"github.com/jakubDoka/sterr"
)

var (
	errRemoval = sterr.New("failed to remove %v from %v with group %d")
)

/*gen(
	Node<TreeEntity, TNode>
)*/

// T is template parameter
type T = int

//def(
//rules Node<T>
//dep Set<T, Set>

// Node keeps over all count of ids in sets for quick lookup
// Its a main building piece of hasher, as we do not expect big amounts of entities in a single node
// it does not use maps to store ids witch is not that elegant but fatser
type Node struct {
	Count int
	Sets  []Set
}

// Insert ...
func (n *Node) Insert(id T, group int) {
	n.Count++
	for i := range n.Sets {
		s := &n.Sets[i]
		if s.Group == group {
			s.IDs = append(s.IDs, id)
			return
		}
	}

	l := len(n.Sets)
	if cap(n.Sets) != l {
		n.Sets = n.Sets[:l+1]
		s := &n.Sets[l]
		s.Group = group
		s.IDs[0] = id
		return
	}
	n.Sets = append(n.Sets, Set{group, []T{id}})

}

// Remove panics if id does not exist within the node, you always have to make sure
// you are removing correctly as leaving dead ids in a hasher is leaking of memory
//
// method panics if object you tried to remove is not present to remove
func (n *Node) Remove(id T, group int) bool {
	n.Count--
	ll := len(n.Sets)
	var nil T // because this is a template
	for i := range n.Sets {
		s := &n.Sets[i]
		if s.Group == group {
			l := len(s.IDs)

			if l == 1 {
				if s.IDs[0] != id {
					return false
				}
				s.IDs[0] = nil
				n.Sets[i], n.Sets[ll-1] = n.Sets[ll-1], *s
				n.Sets = n.Sets[:ll-1]
				return true
			}

			for j := 0; j < l; j++ {
				if id == s.IDs[j] {
					l--
					s.IDs[j] = s.IDs[l]
					s.IDs[l] = nil
					s.IDs = s.IDs[:l]
					return true
				}
			}
		}
	}

	return false
}

// Collect retrieve ids from a node to coll, if include is true only ids of given group
// will get collected, otherwise ewerithing but specified group is returned
func (n *Node) Collect(group int, include bool, coll []T) []T {
	if include {
		for _, s := range n.Sets {
			if s.Group == group {
				coll = append(coll, s.IDs...)
				return coll
			}
		}
	} else {
		for _, s := range n.Sets {
			if s.Group != group {
				coll = append(coll, s.IDs...)
			}
		}
	}
	return coll
}

// CollectAll colects all objects withoud differentiating a group
func (n *Node) CollectAll(coll *[]T) {
	for _, s := range n.Sets {
		*coll = append(*coll, s.IDs...)
	}
}

//)

//def(
//rules Set<T>

// Set is an id set that also has a group important for node
type Set struct {
	Group int
	IDs   []T
}

//)
