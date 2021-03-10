package netw

import (
	"encoding/binary"
	"math"
)

// Size stores byte sizes of types
const (
	Byte    = 1
	Uint16  = 2
	Uint32  = 4
	Uint64  = 8
	Uint    = 8
	Int8    = 1
	Int16   = 2
	Int32   = 4
	Int64   = 8
	Int     = 8
	Float32 = 4
	Float64 = 8
)

// Buffer si for marshaling Data into bites
type Buffer struct {
	Data []byte
	ops  [8]byte

	cursor int
	Failed bool
}

// PutBool puts bool to Buffer
func (b *Buffer) PutBool(input bool) {
	b.Data = append(b.Data, btu(input))
}

// Bool reads bool from Buffer
func (b *Buffer) Bool() bool {
	return utb(b.Byte())
}

// Byte reads byte from Buffer
func (b *Buffer) Byte() byte {
	if b.Advance(Byte) {
		return 0
	}
	return b.Data[b.cursor-Byte]
}

//def(
//rules PutInt32<int32, uint32, Int32, Uint32, PutUint32>

// PutInt32 writes int32 to Buffer
func (b *Buffer) PutInt32(input int32) {
	binary.LittleEndian.PutUint32(b.ops[:Int32], uint32(input))
	b.Data = append(b.Data, b.ops[:Int32]...)
}

// Int32 reads int32 from Buffer
func (b *Buffer) Int32() int32 {
	if b.Advance(Int32) {
		return 0
	}

	return int32(binary.LittleEndian.Uint32(b.Data[b.cursor-Int32 : b.cursor]))
}

//)

/*gen(
	PutInt32<int16, uint16, Int16, Uint16, PutUint16, PutInt16>
	PutInt32<int64, uint64, Int64, Uint64, PutUint64, PutInt64>
	PutInt32<int, uint64, Int, Uint64, PutUint64, PutInt>
	PutInt32<uint16, uint16, Uint16, Uint16, PutUint16, PutUint16>
	PutInt32<uint32, uint32, Uint32, Uint32, PutUint32, PutUint32>
	PutInt32<uint64, uint64, Uint64, Uint64, PutUint64, PutUint64>
	PutInt32<uint, uint64, Uint, Uint64, PutUint64, PutUint>
)*/

// PutFloat64 writes float64 to Buffer
func (b *Buffer) PutFloat64(input float64) error {
	binary.LittleEndian.PutUint64(b.ops[:], math.Float64bits(input))
	b.Data = append(b.Data, b.ops[:]...)
	return nil
}

// Float64 reads float64 from Buffer
func (b *Buffer) Float64() float64 {
	if b.Advance(Float64) {
		return 0
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(b.Data[b.cursor-Float64 : b.cursor]))
}

// PutFloat32 writes float32 to Buffer
func (b *Buffer) PutFloat32(input float32) error {
	binary.LittleEndian.PutUint32(b.ops[:Float32], math.Float32bits(input))
	b.Data = append(b.Data, b.ops[:Float32]...)
	return nil
}

// Float32 reads float32 from Buffer
func (b *Buffer) Float32() float32 {
	if b.Advance(Float32) {
		return 0
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(b.Data[b.cursor-Float32 : b.cursor]))
}

// PutString writes string to byteArray
func (b *Buffer) PutString(input string) {
	b.PutUint32(uint32(len(input)))

	b.Data = append(b.Data, []byte(input)...)
}

// String reads string from Buffer
func (b *Buffer) String() string {
	l := int(b.Uint32())

	if l == 0 || b.Advance(l) {
		return ""
	}

	return string(b.Data[b.cursor-l : b.cursor])
}

// Clear puts buffer to its default value, except Data, that is just truncated to 0 len
func (b *Buffer) Clear() {
	b.Data = b.Data[:0]
	b.cursor = 0
	b.Failed = false
}

// Finished returns whether Buffers cursor is at the end
func (b *Buffer) Finished() bool {
	return b.cursor == len(b.Data)
}

// Advance shifts a Buffer cursor by value and reports if operation failed (true if failed)
func (b *Buffer) Advance(add int) bool {
	if b.Failed {
		return true
	}
	b.cursor += add
	b.Failed = b.cursor > len(b.Data)
	return b.Failed
}

func btu(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func utb(u uint8) bool {
	return u == 1
}

// Builder collects data and encodes it with its length.
// Data in this form can be then send over network and split by Splitter
type Builder struct {
	Buffer
}

// Write writes to builder, this never fails and always writes all data
func (b *Builder) Write(data []byte) (n int, err error) {
	b.PutUint32(uint32(len(data)))
	b.Data = append(b.Data, data...)
	return len(data), nil
}

// Flush returns resulting data and clears the Buffer
func (b *Builder) Flush() []byte {
	dat := b.Data
	b.Clear()
	return dat
}

// Splitter handles splitting of incoming packets into separate Buffers, assuming that
// Buffers were put together by Builder
type Splitter struct {
	Buffer
	leftover   Buffer
	supposed   int
	pool, buff []Buffer
}

func (s *Splitter) Read(data []byte) (n int, err error) {
	if s.supposed != 0 {
		supposed := s.supposed - len(s.leftover.Data)
		if len(data) < supposed {
			s.leftover.Data = append(s.leftover.Data, data...)
			return
		}

		s.leftover.Data = append(s.leftover.Data, data[:supposed]...)
		s.buff = append(s.buff, s.leftover)
		data = data[supposed:]
		s.supposed = 0
	}

	s.Data = append(s.Data, data...)

	for {
		length := int(s.Uint32())
		if s.Failed {
			s.Data = append(s.Data[:0], s.Data[s.cursor-4:]...)
			s.reset()
			return
		}

		if s.Advance(length) {
			s.leftover = s.reuse(s.Data[s.cursor-length:])
			s.supposed = length
			s.Clear()
			return
		}

		s.buff = append(s.buff, s.reuse(s.Data[s.cursor-length:s.cursor]))

		if s.Finished() {
			s.Clear()
			return
		}
	}
}

// Recycle gives Buffer back to splitter so it can reuse it in next split
func (s *Splitter) Recycle(bs ...Buffer) {
	s.pool = append(s.pool, bs...)
}

// Flush returns collected buffers and clears the Splitter
//
// IMPORTANT: Splitter will reuse the slice next time you call Read on him
// so if you want to preserv buffers for later, copy them to slice you
// own. Buffers are cheap to copy.
func (s *Splitter) Flush() []Buffer {
	data := s.buff
	s.buff = s.buff[:0]
	return data
}

func (s *Splitter) reset() {
	s.Failed = false
	s.cursor = 0
}

func (s *Splitter) reuse(data []byte) (buff Buffer) {
	l := len(s.pool) - 1
	if l == -1 {
		buff.Data = append(buff.Data, data...)
		return
	}

	buff, s.pool[l] = s.pool[l], buff
	s.pool = s.pool[:l]
	buff.Data = append(buff.Data, data...)

	return buff
}
