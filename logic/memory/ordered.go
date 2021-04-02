package memory

import "github.com/cheekybits/genny/generic"

type Key generic.Type
type Value generic.Type

// KeyValueCapsule is component of ordered map that stores key and a value
type KeyValueCapsule struct {
	K Key
	V Value
}

// KeyValueOrdered stores its items in underlying slice and map just keeps indexes
type KeyValueOrdered struct {
	m map[Key]int
	s []KeyValueCapsule
}

// NOrderedMap initializes inner map
func NKeyValueOrdered() KeyValueOrdered {
	return KeyValueOrdered{
		m: map[Key]int{},
	}
}

// IsNil reports whether KeyValueOrdered instance is uninitialized
func (o *KeyValueOrdered) IsNil() bool {
	return o.m == nil
}

// Value returns value under key
func (o *KeyValueOrdered) Value(key Key) (val *Value, idx int, ok bool) {
	idx, k := o.m[key]
	if !k {
		return
	}
	return &o.s[idx].V, idx, true
}

// Put puts a value under key
func (o *KeyValueOrdered) Put(key Key, value Value) {
	if i, ok := o.m[key]; ok {
		o.s[i].V = value
	} else {
		o.m[key] = len(o.s)
		o.s = append(o.s, KeyValueCapsule{key, value})
	}
}

// Remove removes the key value pair
func (o *KeyValueOrdered) Remove(key Key) (v Value, i int, b bool) {
	val, idx, ok := o.Value(key)

	if ok {
		o.RemoveIndex(idx)
	} else {
		return
	}

	return *val, idx, ok
}

// RemoveIndex removes by index
func (o *KeyValueOrdered) RemoveIndex(idx int) (cell KeyValueCapsule) {
	cell = o.s[idx]
	delete(o.m, o.s[idx].K)
	o.shift(idx+1, len(o.s), -1)
	o.s = append(o.s[:idx], o.s[idx+1:]...)
	return
}

// Insert insets element under index and key
func (o *KeyValueOrdered) Insert(key Key, idx int, value Value) {
	o.Remove(key)
	o.m[key] = idx
	o.shift(idx, len(o.s), 1)
	o.s = append(append(append(make([]KeyValueCapsule, 0, len(o.s)+1), o.s[:idx]...), KeyValueCapsule{key, value}), o.s[idx:]...)
}

// Slice returns underlying slice
func (o *KeyValueOrdered) Slice() []KeyValueCapsule {
	return o.s
}

// Index returns index of a key's value
func (o *KeyValueOrdered) Index(name Key) (int, bool) {
	val, ok := o.m[name]
	return val, ok
}

// Clear removes all elements
func (o *KeyValueOrdered) Clear() {
	for k := range o.m {
		delete(o.m, k)
	}
	o.s = o.s[:0]
}

// ReIndex changes index of an element
func (o *KeyValueOrdered) ReIndex(old, new int) {
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
	o.m[cell.K] = n
	o.s[n] = cell
}

// Rename renames element and keeps index
func (o *KeyValueOrdered) Rename(old, new Key) bool {
	val, ok := o.m[old]
	if ok {
		o.Remove(new)
		delete(o.m, old)
		o.m[new] = val
		o.s[val].K = new
		return true
	}
	return false
}

func (o *KeyValueOrdered) shift(start, end, dif int) {
	for i := start; i < end; i++ {
		o.m[o.s[i].K] += dif
	}
}
