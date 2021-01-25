package txt

import (
	"gogen/templates"
)

// Effs is a standard Vector type with utility methods
type Effs []Effect

// Clone creates new Effs copies content of v to it and returns it
func (v Effs) Clone() Effs {
	nv := make(Effs, len(v))
	copy(nv, v)

	return nv
}

// Clear is equivalent to Truncate(0)
func (v *Effs) Clear() {
	v.Truncate(0)
}

// Truncate in comparison to truncating by bracket operator also sets all
// forgoten elements to default value, witch is useful if this is slice of pointers
// Effs will have length you specify
func (v *Effs) Truncate(l int) {
	var nil Effect
	dv := *v
	for i := l; i < len(dv); i++ {
		dv[i] = nil
	}

	*v = dv[:l]
}

// Remove removes element and returns it
func (v *Effs) Remove(idx int) (val Effect) {
	var nil Effect

	dv := *v

	val = dv[idx]
	*v = append(dv[:idx], dv[1+idx:]...)

	dv[len(dv)-1] = nil

	return val
}

// RemoveSlice removes sequence of slice
func (v *Effs) RemoveSlice(start, end int) {
	dv := *v

	*v = append(dv[:start], dv[end:]...)

	v.Truncate(len(dv) - (end - start))
}

// PopFront removes first element and returns it
func (v *Effs) PopFront() Effect {
	return v.Remove(0)
}

// Pop removes last element
func (v *Effs) Pop() Effect {
	return v.Remove(len(*v) - 1)
}

// Insert inserts value to given index
func (v *Effs) Insert(idx int, val Effect) {
	dv := *v
	*v = append(append(append(make(Effs, 0, len(dv)+1), dv[:idx]...), val), dv[idx:]...)
}

// InsertSlice inserts slice to given index
func (v *Effs) InsertSlice(idx int, val []Effect) {
	dv := *v
	*v = append(append(append(make(Effs, 0, len(dv)+1), dv[:idx]...), val...), dv[idx:]...)
}

// Reverse reverses content of slice
func (v Effs) Reverse() {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v.Swap(i, j)
	}
}

// Last returns last element of slice
func (v Effs) Last() Effect {
	return v[len(v)-1]
}

// Sort is quicksort for Effs, because this is part of a template comp function is necessary
func (v Effs) Sort(comp func(a, b Effect) bool) {
	if len(v) < 2 {
		return
	}
	// Template is part  of its self, how amazing
	ps := make(templates.IntVec, 2, len(v))
	ps[0], ps[1] = -1, len(v)

	var (
		p Effect

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
func (v Effs) Swap(a, b int) {
	v[a], v[b] = v[b], v[a]
}

// ForEach is a standard foreach method. Its shortcut for modifying all elements
func (v Effs) ForEach(con func(i int, e Effect) Effect) {
	for i, e := range v {
		v[i] = con(i, e)
	}
}

// Filter leaves only elements for with filter returns true
func (v *Effs) Filter(filter func(e Effect) bool) {
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
func (v Effs) Find(find func(e Effect) bool) (idx int, res Effect) {
	for i, e := range v {
		if find(e) {
			return i, e
		}
	}

	idx = -1
	return
}


// FEffs is a standard Vector type with utility methods
type FEffs []*FontEffect

// Clone creates new FEffs copies content of v to it and returns it
func (v FEffs) Clone() FEffs {
	nv := make(FEffs, len(v))
	copy(nv, v)

	return nv
}

// Clear is equivalent to Truncate(0)
func (v *FEffs) Clear() {
	v.Truncate(0)
}

// Truncate in comparison to truncating by bracket operator also sets all
// forgoten elements to default value, witch is useful if this is slice of pointers
// FEffs will have length you specify
func (v *FEffs) Truncate(l int) {
	var nil *FontEffect
	dv := *v
	for i := l; i < len(dv); i++ {
		dv[i] = nil
	}

	*v = dv[:l]
}

// Remove removes element and returns it
func (v *FEffs) Remove(idx int) (val *FontEffect) {
	var nil *FontEffect

	dv := *v

	val = dv[idx]
	*v = append(dv[:idx], dv[1+idx:]...)

	dv[len(dv)-1] = nil

	return val
}

// RemoveSlice removes sequence of slice
func (v *FEffs) RemoveSlice(start, end int) {
	dv := *v

	*v = append(dv[:start], dv[end:]...)

	v.Truncate(len(dv) - (end - start))
}

// PopFront removes first element and returns it
func (v *FEffs) PopFront() *FontEffect {
	return v.Remove(0)
}

// Pop removes last element
func (v *FEffs) Pop() *FontEffect {
	return v.Remove(len(*v) - 1)
}

// Insert inserts value to given index
func (v *FEffs) Insert(idx int, val *FontEffect) {
	dv := *v
	*v = append(append(append(make(FEffs, 0, len(dv)+1), dv[:idx]...), val), dv[idx:]...)
}

// InsertSlice inserts slice to given index
func (v *FEffs) InsertSlice(idx int, val []*FontEffect) {
	dv := *v
	*v = append(append(append(make(FEffs, 0, len(dv)+1), dv[:idx]...), val...), dv[idx:]...)
}

// Reverse reverses content of slice
func (v FEffs) Reverse() {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v.Swap(i, j)
	}
}

// Last returns last element of slice
func (v FEffs) Last() *FontEffect {
	return v[len(v)-1]
}

// Sort is quicksort for FEffs, because this is part of a template comp function is necessary
func (v FEffs) Sort(comp func(a, b *FontEffect) bool) {
	if len(v) < 2 {
		return
	}
	// Template is part  of its self, how amazing
	ps := make(templates.IntVec, 2, len(v))
	ps[0], ps[1] = -1, len(v)

	var (
		p *FontEffect

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
func (v FEffs) Swap(a, b int) {
	v[a], v[b] = v[b], v[a]
}

// ForEach is a standard foreach method. Its shortcut for modifying all elements
func (v FEffs) ForEach(con func(i int, e *FontEffect) *FontEffect) {
	for i, e := range v {
		v[i] = con(i, e)
	}
}

// Filter leaves only elements for with filter returns true
func (v *FEffs) Filter(filter func(e *FontEffect) bool) {
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
func (v FEffs) Find(find func(e *FontEffect) bool) (idx int, res *FontEffect) {
	for i, e := range v {
		if find(e) {
			return i, e
		}
	}

	idx = -1
	return
}

