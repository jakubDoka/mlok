package spatial

import (
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "Item=int"

type Item generic.Type

type ItemNode struct {
	Groups []ItemGroup
	Items  []Item
}

func (n *ItemNode) Collect(group int, including bool, coll []Item) []Item {
	if including {
		for _, g := range n.Groups {
			if g.ID == group {
				coll = append(coll, n.Items[g.Idx:g.Idx+g.Len]...)
				break
			}
		}
	} else {
		for _, g := range n.Groups {
			if g.ID != group {
				coll = append(coll, n.Items[g.Idx:g.Idx+g.Len]...)
			}
		}
	}

	return coll
}

func (n *ItemNode) Insert(item Item, group int) {
	var g *ItemGroup
	var idx int

	for i := range n.Groups {
		if n.Groups[i].ID == group {
			g = &n.Groups[i]
			idx = i
		}
	}

	if g == nil {
		n.Groups = append(n.Groups, ItemGroup{
			ID:  group,
			Idx: len(n.Items),
		})
		g = &n.Groups[len(n.Groups)-1]
	} else {
		for i := idx + 1; i < len(n.Groups); i++ {
			n.Groups[i].Idx++
		}
	}

	g.Len++

	n.Items = append(n.Items, item)
	copy(n.Items[g.Idx+1:], n.Items[g.Idx:])
	n.Items[g.Idx] = item
}

func (n *ItemNode) Remove(item Item, group int) bool {
	var g *ItemGroup
	var idx int

	for i := range n.Groups {
		if n.Groups[i].ID == group {
			g = &n.Groups[i]
			idx = i
			break
		}
	}

	if g == nil {
		return false
	}

	e := g.Idx + g.Len

	for i := g.Idx; ; i++ {
		if i >= e {
			return false
		}
		if n.Items[i] == item {
			n.Items = append(n.Items[:i], n.Items[i+1:]...)
			break
		}
	}

	for i := idx + 1; i < len(n.Groups); i++ {
		n.Groups[i].Idx--
	}

	g.Len--
	if g.Len == 0 {
		n.Groups = append(n.Groups[:idx], n.Groups[idx+1:]...)
	}

	return true
}

type ItemGroup struct {
	ID, Idx, Len int
}

// ItemNode keeps over all count of ids in sets for quick lookup
// Its a main building piece of hasher, as we do not expect big amounts of entities in a single node
// it does not use maps to store ids witch is not that elegant but fatser
/*type ItemNode struct {
	Count int
	Sets  []ItemSet
}

// Insert ...
func (n *ItemNode) Insert(id Item, group int) {
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
	n.Sets = append(n.Sets, ItemSet{group, []Item{id}})

}

// Remove panics if id does not exist within the node, you always have to make sure
// you are removing correctly as leaving dead ids in a hasher is leaking of memory
//
// method panics if object you tried to remove is not present to remove
func (n *ItemNode) Remove(id Item, group int) bool {
	n.Count--
	ll := len(n.Sets)
	var nil Item // because this is a template
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
func (n *ItemNode) Collect(group int, include bool, coll []Item) []Item {
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
func (n *ItemNode) CollectAll(coll *[]Item) {
	for _, s := range n.Sets {
		*coll = append(*coll, s.IDs...)
	}
}

// ItemSet is an id set that also has a group important for node
type ItemSet struct {
	Group int
	IDs   []Item
}*/
