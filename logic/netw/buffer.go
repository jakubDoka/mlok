package netw

import (
	"encoding/binary"
	"math"
	"net"
	"reflect"

	"github.com/jakubDoka/mlok/mat"
)

/*imp(
	github.com/jakubDoka/gogen/templates
)*/

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

var ByteType = reflect.TypeOf(byte(0))

// Buffer si for marshaling Data into bites
type Buffer struct {
	Data []byte
	ops  [8]byte

	cursor int
	Failed bool
}

func (b *Buffer) PutVec(input mat.Vec) {
	b.PutFloat64(input.X)
	b.PutFloat64(input.Y)
}

func (b *Buffer) Vec() mat.Vec {
	return mat.V(b.Float64(), b.Float64())
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
	templates.Resize<buff, Resize>
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

// Writer writes data to connection in a way that Reader can read it on other side
type Writer struct {
	buff []byte
	size [4]byte
}

// Write writes to builder, this never fails and always writes all data
func (w *Writer) Write(conn net.Conn, data []byte) (n int, err error) {
	binary.LittleEndian.PutUint32(w.size[:], uint32(len(data)))

	w.buff = w.buff[:0]
	w.buff = append(w.buff, w.size[:]...)
	w.buff = append(w.buff, data...)

	return conn.Write(w.buff)
}

// Reader handles splitting of incoming packets into separate Buffers, assuming that
// Buffers were put together by Builder
type Reader struct {
	buff buff
	pool []Buffer
	size [4]byte
}

// Read reads one buffer from connection
func (r *Reader) Read(conn net.Conn) (b Buffer, err error) {
	_, err = conn.Read(r.size[:])
	if err != nil {
		return
	}

	ln := int(binary.LittleEndian.Uint32(r.size[:]))
	r.buff.Resize(ln)
	for i := 0; i < ln; {
		n, err := conn.Read(r.buff[i:])
		if err != nil {
			return b, err
		}
		i += n
	}

	b = r.reuse(r.buff)

	return
}

// Recycle gives Buffer back to Reader so it can reuse it in next read
func (r *Reader) Recycle(bs ...Buffer) {
	r.pool = append(r.pool, bs...)
}

// reuse takes an buffer from pool or creates new one if there is none
func (r *Reader) reuse(data []byte) (buff Buffer) {
	l := len(r.pool) - 1
	if l != -1 {
		buff, r.pool[l] = r.pool[l], buff
		r.pool = r.pool[:l]
	}

	buff.Data = append(buff.Data[:0], data...)

	return buff
}

type buff []byte
