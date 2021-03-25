package netw

import (
	"encoding/binary"
)

// PutInt16 writes int16 to Buffer
func (b *Buffer) PutInt16(input int16) {
	binary.LittleEndian.PutUint16(b.ops[:Int16], uint16(input))
	b.Data = append(b.Data, b.ops[:Int16]...)
}

// Int16 reads int16 from Buffer
func (b *Buffer) Int16() int16 {
	if b.Advance(Int16) {
		return 0
	}

	return int16(binary.LittleEndian.Uint16(b.Data[b.cursor-Int16 : b.cursor]))
}


// PutInt64 writes int64 to Buffer
func (b *Buffer) PutInt64(input int64) {
	binary.LittleEndian.PutUint64(b.ops[:Int64], uint64(input))
	b.Data = append(b.Data, b.ops[:Int64]...)
}

// Int64 reads int64 from Buffer
func (b *Buffer) Int64() int64 {
	if b.Advance(Int64) {
		return 0
	}

	return int64(binary.LittleEndian.Uint64(b.Data[b.cursor-Int64 : b.cursor]))
}


// PutInt writes int to Buffer
func (b *Buffer) PutInt(input int) {
	binary.LittleEndian.PutUint64(b.ops[:Int], uint64(input))
	b.Data = append(b.Data, b.ops[:Int]...)
}

// Int reads int from Buffer
func (b *Buffer) Int() int {
	if b.Advance(Int) {
		return 0
	}

	return int(binary.LittleEndian.Uint64(b.Data[b.cursor-Int : b.cursor]))
}


// PutUint16 writes uint16 to Buffer
func (b *Buffer) PutUint16(input uint16) {
	binary.LittleEndian.PutUint16(b.ops[:Uint16], uint16(input))
	b.Data = append(b.Data, b.ops[:Uint16]...)
}

// Uint16 reads uint16 from Buffer
func (b *Buffer) Uint16() uint16 {
	if b.Advance(Uint16) {
		return 0
	}

	return uint16(binary.LittleEndian.Uint16(b.Data[b.cursor-Uint16 : b.cursor]))
}


// PutUint32 writes uint32 to Buffer
func (b *Buffer) PutUint32(input uint32) {
	binary.LittleEndian.PutUint32(b.ops[:Uint32], uint32(input))
	b.Data = append(b.Data, b.ops[:Uint32]...)
}

// Uint32 reads uint32 from Buffer
func (b *Buffer) Uint32() uint32 {
	if b.Advance(Uint32) {
		return 0
	}

	return uint32(binary.LittleEndian.Uint32(b.Data[b.cursor-Uint32 : b.cursor]))
}


// PutUint64 writes uint64 to Buffer
func (b *Buffer) PutUint64(input uint64) {
	binary.LittleEndian.PutUint64(b.ops[:Uint64], uint64(input))
	b.Data = append(b.Data, b.ops[:Uint64]...)
}

// Uint64 reads uint64 from Buffer
func (b *Buffer) Uint64() uint64 {
	if b.Advance(Uint64) {
		return 0
	}

	return uint64(binary.LittleEndian.Uint64(b.Data[b.cursor-Uint64 : b.cursor]))
}


// PutUint writes uint to Buffer
func (b *Buffer) PutUint(input uint) {
	binary.LittleEndian.PutUint64(b.ops[:Uint], uint64(input))
	b.Data = append(b.Data, b.ops[:Uint]...)
}

// Uint reads uint from Buffer
func (b *Buffer) Uint() uint {
	if b.Advance(Uint) {
		return 0
	}

	return uint(binary.LittleEndian.Uint64(b.Data[b.cursor-Uint : b.cursor]))
}


// Resize resizes the buff
func (v *buff) Resize(size int) {
	if cap(*v) >= size {
		*v = (*v)[:size]
	} else {
		ns := make(buff, size)
		copy(ns, *v)
		*v = ns
	}
}

