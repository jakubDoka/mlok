package memory

import (
	"github.com/cheekybits/genny/generic"
	"github.com/jakubDoka/mlok/logic/memory/gen"
)

type Element generic.Type

// ElementCapsule is something like an optional type, it holds boolean about whether
// it contains value though it does not hold pointer
type ElementCapsule struct {
	occupied bool
	value    Element
}

// ElementStorage generates IDs witch makes no need to use hashing,
// only drawback is that you cannot choose the id, it will be assigned
// like a pointer, but without putting presure no gc, brilliant ElementStorage for
// components. Its highly unlikely you will run out of ids as they are reused
type ElementStorage struct {
	vec      []ElementCapsule
	freeIDs  gen.Int32Vec
	occupied []int32
	count    int
	outdated bool
}

// Blanc allocates blanc space adds
func (s *ElementStorage) Blanc() {
	s.freeIDs = append(s.freeIDs, int32(len(s.vec)))
	s.vec = append(s.vec, ElementCapsule{})
}

// Allocate id allocates if it is free, else it returns nil
func (s *ElementStorage) AllocateID(id int32) *Element {
	if int(id) >= len(s.vec) || s.vec[id].occupied {
		return nil
	}

	idx, _ := s.freeIDs.BiSearch(id, gen.Int32BiComp)
	s.freeIDs.Remove(idx)

	return &s.vec[id].value
}

// Allocate allocates an value and returns id and pointer to it. Note that
// allocate does not always allocate at all and just reuses freed space,
// returned pointer also does not point to zero value and you have to overwrite all
// properties to get expected behavior
func (s *ElementStorage) Allocate() (*Element, int32) {
	s.count++
	s.outdated = true

	l := len(s.freeIDs)
	if l != 0 {
		id := s.freeIDs[l-1]
		s.freeIDs = s.freeIDs[:l-1]
		s.vec[id].occupied = true
		return &s.vec[id].value, id
	}

	s.vec = append(s.vec, ElementCapsule{})
	id := int32(len(s.vec)) - 1
	s.vec[id].occupied = true
	return &s.vec[id].value, id
}

// Remove removes a value and frees memory for something else
//
// panic if there is nothing to free
func (s *ElementStorage) Remove(id int32) {
	if !s.vec[id].occupied {
		panic("removeing already removed value")
	}

	s.count--
	s.outdated = true

	s.freeIDs.BiInsert(id, gen.Int32BiComp)
	s.vec[id].occupied = false
}

// Item returns pointer to value under the "id", accessing random id can result in
// random value that can be considered unoccupied
//
// method panics if id is not occupied
func (s *ElementStorage) Item(id int32) *Element {
	if !s.vec[id].occupied {
		panic("accessing non occupied id")
	}

	return &s.vec[id].value
}

// Used returns whether id is used
func (s *ElementStorage) Used(id int32) bool {
	return s.vec[id].occupied
}

// Len returns size of ElementStorage
func (s *ElementStorage) Len() int {
	return len(s.vec)
}

// Count returns amount of values stored
func (s *ElementStorage) Count() int {
	return s.count
}

// update updates state of occupied slice, every time you remove or add
// element, ElementStorage gets outdated, this makes it up to date
func (s *ElementStorage) update() {
	s.outdated = false
	s.occupied = s.occupied[:0]
	l := int32(len(s.vec))
	for i := int32(0); i < l; i++ {
		if s.vec[i].occupied {
			s.occupied = append(s.occupied, i)
		}
	}
}

// Occupied return all occupied ids in ElementStorage, this method panics if ElementStorage is outdated
// See Update method.
func (s *ElementStorage) Occupied() []int32 {
	if s.outdated {
		s.update()
	}

	return s.occupied
}

// Clear clears ElementStorage, but keeps allocated space
func (s *ElementStorage) Clear() {
	s.vec = s.vec[:0]
	s.occupied = s.occupied[:0]
	s.freeIDs = s.freeIDs[:0]
	s.count = 0
}

// SlowClear clears the the ElementStorage slowly with is tradeoff for having faster allocating speed
func (s *ElementStorage) SlowClear() {
	for i := range s.vec {
		if s.vec[i].occupied {
			s.freeIDs = append(s.freeIDs, int32(i))
			s.vec[i].occupied = false
		}
	}

	s.occupied = s.occupied[:0]
	s.count = 0
}
