package pck

import (
	"github.com/jakubDoka/gogen/templates"
)

// Data is a standard Vector type with utility methods
type Data []PicData

// Clone creates new Data copies content of v to it and returns it
func (v Data) Clone() Data {
	nv := make(Data, len(v))
	copy(nv, v)

	return nv
}

// Clear is equivalent to Truncate(0)
func (v *Data) Clear() {
	v.Truncate(0)
}

// Rewrite revrites elements from index to o values
func (v Data) Rewrite(o Data, idx int) {
	copy(v[idx:], o)
}

// Len implements VertexData interface
func (v Data) Len() int {
	return len(v)
}

// UClear does not care about memory leaks, it just sets length to 0
func (v *Data) UClear() {
	*v = (*v)[:0]
}

// Truncate in comparison to truncating by bracket operator also sets all
// forgoten elements to default value, witch is useful if this is slice of pointers
// Data will have length you specify
func (v *Data) Truncate(l int) {
	var nil PicData
	dv := *v
	for i := l; i < len(dv); i++ {
		dv[i] = nil
	}

	*v = dv[:l]
}

// Extend extends vec size by amount so then len(Data) = len(Data) + amount
func (v *Data) Extend(amount int) {
	vv := *v
	l := len(vv) + amount
	if cap(vv) >= l { // no need to allocate
		*v = vv[:l]
		return
	}

	nv := make(Data, l)
	copy(nv, vv)
	*v = nv
}

// Remove removes element and returns it
func (v *Data) Remove(idx int) (val PicData) {
	var nil PicData

	dv := *v

	val = dv[idx]
	*v = append(dv[:idx], dv[1+idx:]...)

	dv[len(dv)-1] = nil

	return val
}

// RemoveSlice removes sequence of slice
func (v *Data) RemoveSlice(start, end int) {
	dv := *v

	*v = append(dv[:start], dv[end:]...)

	v.Truncate(len(dv) - (end - start))
}

// PopFront removes first element and returns it
func (v *Data) PopFront() PicData {
	return v.Remove(0)
}

// Pop removes last element
func (v *Data) Pop() PicData {
	return v.Remove(len(*v) - 1)
}

// Insert inserts value to given index
func (v *Data) Insert(idx int, val PicData) {
	dv := *v
	*v = append(append(append(make(Data, 0, len(dv)+1), dv[:idx]...), val), dv[idx:]...)
}

// InsertSlice inserts slice to given index
func (v *Data) InsertSlice(idx int, val []PicData) {
	dv := *v
	*v = append(append(append(make(Data, 0, len(dv)+1), dv[:idx]...), val...), dv[idx:]...)
}

// Reverse reverses content of slice
func (v Data) Reverse() {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v.Swap(i, j)
	}
}

// Last returns last element of slice
func (v Data) Last() PicData {
	return v[len(v)-1]
}

// Sort is quicksort for Data, because this is part of a template comp function is necessary
func (v Data) Sort(comp func(a, b PicData) bool) {
	if len(v) < 2 {
		return
	}
	// Template is part  of its self, how amazing
	ps := make(templates.IntVec, 2, len(v))
	ps[0], ps[1] = -1, len(v)

	var (
		p PicData

		l, e, s, j int
	)

	for {
		l = len(ps)

		e = ps[l-1] - 1
		if e <= 0 {
			return
		}

		s = ps[l-2] + 1
		p = v[e]

		if s < e {
			for j = s; j < e; j++ {
				if comp(v[j], p) {
					v.Swap(s, j)
					s++
				}
			}

			v.Swap(s, e)
			ps.Insert(l-1, s)
		} else {
			ps = ps[:l-1]
		}
	}
}

// Swap swaps two elements
func (v Data) Swap(a, b int) {
	v[a], v[b] = v[b], v[a]
}

// ForEach is a standard foreach method. Its shortcut for modifying all elements
func (v Data) ForEach(con func(i int, e PicData) PicData) {
	for i, e := range v {
		v[i] = con(i, e)
	}
}

// Filter leaves only elements for with filter returns true
func (v *Data) Filter(filter func(e PicData) bool) {
	dv := *v

	var i int
	for _, e := range dv {
		if filter(e) {
			dv[i] = e
			i++
		}
	}

	v.Truncate(i)
}

// Find returns first element for which find returns true along with index,
// if there is none, index equals -1
func (v Data) Find(find func(e PicData) bool) (idx int, res PicData) {
	for i, e := range v {
		if find(e) {
			return i, e
		}
	}

	idx = -1
	return
}

//BiSearch performs a binary search on Ves assuming it is sorted. cmp consumer should
// return 0 if a == b equal, 1 if a > b and 2 if b > a, even if value wos not found it returns
// it returns closest index and false. If Data is empty -1 and false is returned
func (v Data) BiSearch(value PicData, cmp func(a, b PicData) uint8) (int, bool) {
	start, end := 0, len(v)
	if start == end {
		return -1, false
	}
	for {
		mid := start + (end-start)/2
		switch cmp(v[mid], value) {
		case 0:
			return mid, true
		case 1:
			end = mid + 0
		case 2:
			start = mid + 1
		}

		if start == end {
			return start, start < len(v) && cmp(v[start], value) == 0
		}
	}
}

// BiInsert inserts inserts value in a way that keeps vec sorted, binary search is used to determinate
// where to insert
func (v *Data) BiInsert(value PicData, cmp func(a, b PicData) uint8) {
	i, _ := v.BiSearch(value, cmp)
	v.Insert(i, value)
}

// Move moves value from old to new shifting elements in between
func (v Data) Move(old, new int) {
	if old == new {
		return // well
	}

	shifting := -1
	o, n := old, new
	if old > new {
		shifting = 1
		old, new = new+1, old+1
	}

	cell := v[o]
	copy(v[old:new], v[old-shifting:new-shifting])
	v[n] = cell
}

