package ui


// Children stores its items in underlying slice and map just keeps indexes
type Children struct {
	m map[string]int
	s []OMC2
}

// NChildren initializes inner map
func NChildren() Children {
	return Children{
		m: map[string]int{},
	}
}

// Value ...
func (o *Children) Value(key string) (val *Div, idx int, ok bool) {
	idx, k := o.m[key]
	if !k {
		return
	}
	return o.s[idx].Value, idx, true
}

// Put ...
func (o *Children) Put(key string, value *Div) {
	if i, ok := o.m[key]; ok {
		o.s[i].Value = value
	} else {
		o.m[key] = len(o.s)
		o.s = append(o.s, OMC2{key, value})
	}
}

// Remove can be very slow if map is huge
func (o *Children) Remove(key string) (val *Div, idx int, ok bool) {
	val, idx, ok = o.Value(key)
	if ok {
		o.RemoveIndex(idx)
	}
	return
}

// RemoveIndex removes by index
func (o *Children) RemoveIndex(idx int) (cell OMC2) {
	cell = o.s[idx]
	delete(o.m, o.s[idx].Key)
	o.shift(idx+1, len(o.s), -1)
	o.s = append(o.s[:idx], o.s[idx+1:]...)
	return
}

// Insert insets element under index and key
func (o *Children) Insert(key string, idx int, value *Div) {
	o.Remove(key)
	o.m[key] = idx
	o.shift(idx, len(o.s), 1)
	o.s = append(append(append(make([]OMC2, 0, len(o.s)+1), o.s[:idx]...), OMC2{key, value}), o.s[idx:]...)
}

// Slice returns underlying slice
func (o *Children) Slice() []OMC2 {
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


// OMC2 is component of ordered map that stores key the index is under
type OMC2 struct {
	Key   string
	Value *Div
}

