// Package memory offers gogen templates for memory managemant, structs
// can help to compose structure that takes garbage collection into account
// it gets rid of pointers, or lowers the allocations
package memory

// T is template parameter
type T struct{}

//def(
//rules Capsule<T>

// Capsule is something like an optional type, it holds boolean about whether
// it contains value though it does not hold pointer
type Capsule struct {
	occupied bool
	value    T
}

//)

// Storage ...
//def(
//rules Storage<T, int32>
//dep Capsule<T, Capsule>
// Storage generates IDs witch makes no need to use hashing,
// only drawback is that you cannot choose the id, it will be assigned
// like a pointer, but without putting presure no gc, brilliant storage for
// components. Its highly unlikely you will run out of ids as they are reused
type Storage struct {
	vec      []Capsule
	freeIDs  []int32
	occupied []int32
	count    int
	outdated bool
}

// Allocate allocates an value and returns id and pointer to it. Note that
// allocate does not always allocate at all and just reuses freed space,
// returned pointer also does not point to zero value and you have to overwrite all
// properties to get expected behavior
func (s *Storage) Allocate() (*T, int32) {
	s.count++
	s.outdated = true

	l := len(s.freeIDs)
	if l != 0 {
		id := s.freeIDs[l-1]
		s.freeIDs = s.freeIDs[:l-1]
		s.vec[id].occupied = true
		return &s.vec[id].value, id
	}

	s.vec = append(s.vec, Capsule{})
	id := int32(len(s.vec)) - 1
	s.vec[id].occupied = true
	return &s.vec[id].value, id
}

// Remove removes a value and frees memory for something else
//
// panic if there is nothing to free
func (s *Storage) Remove(id int32) {
	if !s.vec[id].occupied {
		panic("removeing already removed value")
	}

	s.count--
	s.outdated = true

	s.freeIDs = append(s.freeIDs, id)
	s.vec[id].occupied = false
}

// Item returns pointer to value under the "id", accessing random id can result in
// random value that can be considered unoccupied
//
// method panics if id is not occupied
func (s *Storage) Item(id int32) *T {
	if !s.vec[id].occupied {
		panic("accessing non occupied id")
	}

	return &s.vec[id].value
}

// Used returns whether id is used
func (s *Storage) Used(id int32) bool {
	return s.vec[id].occupied
}

// Len returns size of storage
func (s *Storage) Len() int {
	return len(s.vec)
}

// Count returns amount of values stored
func (s *Storage) Count() int {
	return s.count
}

// update updates state of occupied slice, every time you remove or add
// element, storage gets outdated, this makes it up to date
func (s *Storage) update() {
	s.outdated = false
	s.occupied = s.occupied[:0]
	l := int32(len(s.vec))
	for i := int32(0); i < l; i++ {
		if s.vec[i].occupied {
			s.occupied = append(s.occupied, i)
		}
	}
}

// Occupied return all occupied ids in storage, this method panics if Storage is outdated
// See Update method.
func (s *Storage) Occupied() []int32 {
	if s.outdated {
		s.update()
	}

	return s.occupied
}

// Clear clears storage, but keeps allocated space
func (s *Storage) Clear() {
	s.vec = s.vec[:0]
	s.occupied = s.occupied[:0]
	s.freeIDs = s.freeIDs[:0]
	s.count = 0
}

// SlowClear clears the the Storage slowly with is tradeoff for having faster allocating speed
func (s *Storage) SlowClear() {
	for i := range s.vec {
		if s.vec[i].occupied {
			s.freeIDs = append(s.freeIDs, int32(i))
			s.vec[i].occupied = false
		}
	}

	s.occupied = s.occupied[:0]
	s.count = 0
}

//)

// QuickPool ...
//def(
//rules QuickPool<T>
// QuickPool is template that can store interface referenced
// objects on one place reducing allocation, if you are creating
// lot of date stored behind interface that has very short livetime
// QuickPool can be nice optimization
type QuickPool struct {
	vec    []T
	cursor int
}

// Item returns pointer to pooling struct that is free, if no free struct is present
// new one is allocated
func (q *QuickPool) Item() *T {
	if q.cursor == len(q.vec) {
		q.vec = append(q.vec, T{})
	}
	q.cursor++
	return &q.vec[q.cursor-1]
}

// Restart makes QuickPool reuse old objects, if you cannot call this
// as pooling structs you taken are constantly referenced, its useles to
// use QuickPool in first place
func (q *QuickPool) Restart() {
	q.cursor = 0
}

//)
