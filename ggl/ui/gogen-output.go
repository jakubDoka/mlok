package ui


// Children stores its items in underlying slice and map just keeps indexes
type Children struct {
	m map[string]int
	s []OMC
}

// NChildren initializes inner map
func NChildren() Children {
	return Children{
		m: map[string]int{},
	}
}

// IsNil reports whether Children instance is uninitialized
func (o *Children) IsNil() bool {
	return o.m == nil
}

// Value ...
func (o *Children) Value(key string) (val *Element, idx int, ok bool) {
	idx, k := o.m[key]
	if !k {
		return
	}
	return o.s[idx].Value, idx, true
}

// Put ...
func (o *Children) Put(key string, value *Element) {
	if i, ok := o.m[key]; ok {
		o.s[i].Value = value
	} else {
		o.m[key] = len(o.s)
		o.s = append(o.s, OMC{key, value})
	}
}

// Remove can be very slow if map is huge
func (o *Children) Remove(key string) (val *Element, idx int, ok bool) {
	val, idx, ok = o.Value(key)
	if ok {
		o.RemoveIndex(idx)
	}
	return
}

// RemoveIndex removes by index
func (o *Children) RemoveIndex(idx int) (cell OMC) {
	cell = o.s[idx]
	delete(o.m, o.s[idx].Key)
	o.shift(idx+1, len(o.s), -1)
	o.s = append(o.s[:idx], o.s[idx+1:]...)
	return
}

// Insert insets element under index and key
func (o *Children) Insert(key string, idx int, value *Element) {
	o.Remove(key)
	o.m[key] = idx
	o.shift(idx, len(o.s), 1)
	o.s = append(append(append(make([]OMC, 0, len(o.s)+1), o.s[:idx]...), OMC{key, value}), o.s[idx:]...)
}

// Slice returns underlying slice
func (o *Children) Slice() []OMC {
	return o.s
}

// Index returns index of a key's value
func (o *Children) Index(name string) (int, bool) {
	val, ok := o.m[name]
	return val, ok
}

// Clear removes all elements
func (o *Children) Clear() {
	for k := range o.m {
		delete(o.m, k)
	}
	o.s = o.s[:0]
}

// ReIndex changes index of an element
func (o *Children) ReIndex(old, new int) {
	if old == new {
		return // well
	}

	shifting := -1
	ol, n := old, new
	if old > new {
		shifting = 1
		old, new = new+1, old+1
	}

	cell := o.s[ol]
	o.shift(old-shifting, new-shifting, shifting)
	copy(o.s[old:new], o.s[old-shifting:new-shifting])
	o.m[cell.Key] = n
	o.s[n] = cell
}

// Rename renames element and keeps index
func (o *Children) Rename(old, new string) bool {
	val, ok := o.m[old]
	if ok {
		o.Remove(new)
		delete(o.m, old)
		o.m[new] = val
		o.s[val].Key = new
		return true
	}
	return false
}

func (o *Children) shift(start, end, dif int) {
	for i := start; i < end; i++ {
		o.m[o.s[i].Key] += dif
	}
}


// OMC is component of ordered map that stores key the index is under
type OMC struct {
	Key   string
	Value *Element
}

