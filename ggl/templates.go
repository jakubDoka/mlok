package ggl

import "gobatch/mt"

const elementsize = 0

// Vertex2D is essentia vertex struct for 2D rendering
type Vertex2D struct {
	Pos, Tex mt.V2
	Color    mt.RGBA
}

/*gen(
	VertexSlice<Vertex2D, 8, VS2D>
	Data<VS2D, Data2D>
	Batch<Data2D, Batch2D>
)*/

//def(
//rules VertexSlice<interface{}, elementsize>

// VertexSlice implements essential utility methods for any
// struct satisfying VertexData interface. Its a gogen TEMPLATE:
//
// 	/*gen(
//		VertexSlice<Vertex2D, 8, VS2D>
//	)*/
//
// this block generates VertexSlice with Vertex2D and 8 is the size of Vertex2D,
// VS2D is name of generated struct divided by float64 byte size for more info
// search github.com/jakubDoka/gogen.
type VertexSlice []interface{}

// Rewrite revrites elements from index to o values
func (v VertexSlice) Rewrite(o VertexSlice, idx int) {
	copy(v[idx:], o)
}

// Clear clears slice
func (v *VertexSlice) Clear() {
	*v = (*v)[:0]
}

// Len implements VertexData interface
func (v VertexSlice) Len() int {
	return len(v)
}

// VertexSize implements VertexData interface
func (v VertexSlice) VertexSize() int {
	return elementsize
}

//)

//def(
//rules Data<VertexSlice>

// Data is Vertex and indice collector, mainly utility that handles vertex offsets
type Data struct {
	Vertexes VertexSlice
	indices  Indices
}

// Clear clears t batch but leaves allocated data
func (d *Data) Clear() {
	d.Vertexes.Clear()
	d.indices.Clear()
}

// Accept accepts vertex data, this is only correct way of feeding batch with Vertexes
// along side indices, if you don't use indices append directly to Data
func (d *Data) Accept(data VertexSlice, indices Indices) {
	l1 := len(d.indices)
	l2 := uint32(d.Vertexes.Len())

	d.Vertexes = append(d.Vertexes, data...)
	d.indices = append(d.indices, indices...)

	l3 := len(d.indices)
	for i := l1; i < l3; i++ {
		d.indices[i] += l2
	}
}

//)

//def(
//rules Batch<Data>

// Batch is main drawer, it performs direct draw to canvas and is used as target for Sprite
// Batch acts like canvas i some ways but performance difference of drawing batch to canvas and
// drawing canvas to canvas is significant. if you need image to remain use canvas, as its name hints
// it in deed works like its called.
type Batch struct {
	Data

	buffer  *Buffer
	program *Program
	texture *Texture
}

// NBatch allows constructing batch with custom Buffer and Program for applying
// per batch shader and related buffer structure. Passing nil absolutely fine,
// as canvas or vindow will use theier own, if you don't even need texture use struct
// literal (Batch{}) to construct batch
func NBatch(texture *Texture, buffer *Buffer, program *Program) *Batch {
	return &Batch{
		texture: texture,
		buffer:  buffer,
		program: program,
	}
}

// Draw draws all data to target
func (b *Batch) Draw(target Target) {
	target.Accept(b.Vertexes, b.indices, b.texture, b.program, b.buffer)
}

// Program returns batch program, it can be nil
func (b *Batch) Program() *Program {
	return b.program
}

//)
