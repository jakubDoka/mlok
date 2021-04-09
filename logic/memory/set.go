package memory

import (
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -pkg=gen -in=$GOFILE -out=gen/$GOFILE gen "Elem=NUMBERS"

type Elem generic.Number

// ElemSet uses binary search to store comparable elements
// it acts as normal Set though, with performance advantage
// over map
type ElemSet []Elem

// Has returns whether set has the value
func (s *ElemSet) Has(val Elem) bool {
	_, ok := s.find(val)
	return ok
}

// Remove removes the value or returns false if value is not present
func (s *ElemSet) Remove(val Elem) bool {
	idx, ok := s.find(val)
	if !ok {
		return false
	}
	v := *s
	*s = append(v[:idx], v[idx+1:]...)
	return true
}

// Insert inserts the element or returns false if element already is present
func (s *ElemSet) Insert(val Elem) bool {
	idx, ok := s.find(val)
	if ok {
		return false
	}

	dv := *s
	e := len(dv)
	dv = append(dv, val)
	for i := e - 1; i >= idx; i-- {
		dv[i+1] = dv[i]
	}
	dv[idx] = val
	*s = dv
	return true
}

func (s *ElemSet) find(e Elem) (int, bool) {
	v := *s
	start, end := 0, len(v)
	if start == end {
		return 0, false
	}

	if end < 17 { // benchmarked value
		for i := range v {
			if v[i] == e {
				return i, true
			}
		}

		return end, false
	}

	for {
		mid := start + (end-start)/2
		val := v[mid]
		if val == e {
			return mid, true
		} else if val > e {
			end = mid + 0
		} else {
			start = mid + 1
		}

		if start == end {
			return start, start < len(v) && v[start] == e
		}
	}
}
