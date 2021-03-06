// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package main

import "github.com/jakubDoka/mlok/logic/memory/gen"

// EntityCapsule is something like an optional type, it holds boolean about whether
// it contains value though it does not hold pointer
type EntityCapsule struct {
	occupied bool
	value    Entity
}

// EntityStorage generates IDs witch makes no need to use hashing,
// only drawback is that you cannot choose the id, it will be assigned
// like a pointer, but without putting presure no gc, brilliant EntityStorage for
// components. Its highly unlikely you will run out of ids as they are reused
type EntityStorage struct {
	vec      []EntityCapsule
	freeIDs  gen.IntVec
	occupied []int
	count    int
	outdated bool
}

// Blanc allocates blanc space adds
func (s *EntityStorage) Blanc() {
	s.freeIDs = append(s.freeIDs, len(s.vec))
	s.vec = append(s.vec, EntityCapsule{})
}

// Allocate id allocates if it is free, else it returns nil
func (s *EntityStorage) AllocateID(id int) *Entity {
	if int(id) >= len(s.vec) || s.vec[id].occupied {
		return nil
	}

	idx, _ := s.freeIDs.BiSearch(id, gen.IntBiComp)
	s.freeIDs.Remove(idx)

	return &s.vec[id].value
}

// Allocate allocates an value and returns id and pointer to it. Note that
// allocate does not always allocate at all and just reuses freed space,
// returned pointer also does not point to zero value and you have to overwrite all
// properties to get expected behavior
func (s *EntityStorage) Allocate() (*Entity, int) {
	s.count++
	s.outdated = true

	l := len(s.freeIDs)
	if l != 0 {
		id := s.freeIDs[l-1]
		s.freeIDs = s.freeIDs[:l-1]
		s.vec[id].occupied = true
		return &s.vec[id].value, id
	}

	id := len(s.vec)
	s.vec = append(s.vec, EntityCapsule{})

	s.vec[id].occupied = true
	return &s.vec[id].value, id
}

// Remove removes a value and frees memory for something else
//
// panic if there is nothing to free
func (s *EntityStorage) Remove(id int) {
	if !s.vec[id].occupied {
		panic("removeing already removed value")
	}

	s.count--
	s.outdated = true

	s.freeIDs.BiInsert(id, gen.IntBiComp)
	s.vec[id].occupied = false
}

// Item returns pointer to value under the "id", accessing random id can result in
// random value that can be considered unoccupied
//
// method panics if id is not occupied
func (s *EntityStorage) Item(id int) *Entity {
	if !s.vec[id].occupied {
		panic("accessing non occupied id")
	}

	return &s.vec[id].value
}

// Used returns whether id is used
func (s *EntityStorage) Used(id int) bool {
	return s.vec[id].occupied
}

// Len returns size of EntityStorage
func (s *EntityStorage) Len() int {
	return len(s.vec)
}

// Count returns amount of values stored
func (s *EntityStorage) Count() int {
	return s.count
}

// update updates state of occupied slice, every time you remove or add
// element, EntityStorage gets outdated, this makes it up to date
func (s *EntityStorage) update() {
	s.outdated = false
	s.occupied = s.occupied[:0]
	l := len(s.vec)
	for i := 0; i < l; i++ {
		if s.vec[i].occupied {
			s.occupied = append(s.occupied, i)
		}
	}
}

// Occupied return all occupied ids in EntityStorage, this method panics if EntityStorage is outdated
// See Update method.
func (s *EntityStorage) Occupied() []int {
	if s.outdated {
		s.update()
	}

	return s.occupied
}

// Clear clears EntityStorage, but keeps allocated space
func (s *EntityStorage) Clear() {
	s.vec = s.vec[:0]
	s.occupied = s.occupied[:0]
	s.freeIDs = s.freeIDs[:0]
	s.count = 0
}

// SlowClear clears the the EntityStorage slowly with is tradeoff for having faster allocating speed
func (s *EntityStorage) SlowClear() {
	for i := range s.vec {
		if s.vec[i].occupied {
			s.freeIDs = append(s.freeIDs, i)
			s.vec[i].occupied = false
		}
	}

	s.occupied = s.occupied[:0]
	s.count = 0
}
