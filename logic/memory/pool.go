package memory

import "github.com/cheekybits/genny/generic"

type Quick generic.Type

// QuickPool is template that can store interface referenced
// objects on one place reducing allocation, if you are creating
// lot of date stored behind interface that has very short livetime
// QuickPool can be nice optimization
type QuickPool struct {
	vec    []Quick
	cursor int
}

// Item returns pointer to pooling struct that is free, if no free struct is present
// new one is allocated
func (q *QuickPool) Item(new Quick) *Quick {
	if q.cursor == len(q.vec) {
		var nil Quick
		q.vec = append(q.vec, nil)
	}
	q.cursor++
	return &q.vec[q.cursor-1]
}

// Over overwrites by raw value and returns pointer to its new position
func (q *QuickPool) Over(val Quick) *Quick {
	if q.cursor == len(q.vec) {
		q.vec = append(q.vec, val)
	} else {
		q.vec[q.cursor] = val
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
