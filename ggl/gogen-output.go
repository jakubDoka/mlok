package ggl


// VS2D implements essential utility methods for any
// struct satisfying VertexData interface. Its a gogen TEMPLATE:
//
// 	/*gen(
//		VS2D<Vertex2D, 8, VS2D>
//	)*/
//
// this block generates VS2D with Vertex2D and 8 is the size of Vertex2D,
// VS2D is name of generated struct divided by float64 byte size for more info
// search github.com/jakubDoka/gogen.
type VS2D []Vertex2D

// Rewrite revrites elements from index to o values
func (v VS2D) Rewrite(o VS2D, idx int) {
	copy(v[idx:], o)
}

// Clear clears slice
func (v *VS2D) Clear() {
	*v = (*v)[:0]
}

// Len implements VertexData interface
func (v VS2D) Len() int {
	return len(v)
}

// VertexSize implements VertexData interface
func (v VS2D) VertexSize() int {
	return 8
}


// Data2D is Vertex and indice collector, mainly utility that handles vertex offsets
type Data2D struct {
	Vertexes VS2D
	indices  Indices
}

// Clear clears t batch but leaves allocated data
func (d *Data2D) Clear() {
	d.Vertexes.Clear()
	d.indices.Clear()
}

// Accept accepts vertex data, this is only correct way of feeding batch with Vertexes
// along side indices, if you don't use indices append directly to Data2D
func (d *Data2D) Accept(data VS2D, indices Indices) {
	l1 := len(d.indices)
	l2 := uint32(d.Vertexes.Len())

	d.Vertexes = append(d.Vertexes, data...)
	d.indices = append(d.indices, indices...)

	l3 := len(d.indices)
	for i := l1; i < l3; i++ {
		d.indices[i] += l2
	}
}


// Batch2D is main drawer, it performs direct draw to canvas and is used as target for Sprite
// Batch2D acts like canvas i some ways but performance difference of drawing batch to canvas and
// drawing canvas to canvas is significant. if you need image to remain use canvas, as its name hints
// it in deed works like its called.
type Batch2D struct {
	Data2D

	buffer  *Buffer
	program *Program
	texture *Texture
}

// NBatch2D allows constructing batch with custom Buffer and Program for applying
// per batch shader and related buffer structure. Passing nil absolutely fine,
// as canvas or vindow will use theier own, if you don't even need texture use struct
// literal (Batch2D{}) to construct batch
func NBatch2D(texture *Texture, buffer *Buffer, program *Program) *Batch2D {
	return &Batch2D{
		texture: texture,
		buffer:  buffer,
		program: program,
	}
}

// Draw draws all data to target
func (b *Batch2D) Draw(target Target) {
	target.Accept(b.Vertexes, b.indices, b.texture, b.program, b.buffer)
}

