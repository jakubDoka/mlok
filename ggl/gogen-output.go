package ggl


// VS2D implements essential utility methods for any
// struct satisfying VertexData interface. Its a gogen TEMPLATE:
//
// 	/*gen(
//		VS2D<Vertex2D, 9, VS2D>
//	)*/
//
// this block generates VertexSlice with Vertex2D and 8 is the size of Vertex2D,
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
	return 9
}


// Data2D is Vertex and indice collector, mainly utility that handles vertex offsets
type Data2D struct {
	Vertexes VS2D
	Indices  Indices
}

// Clear clears t batch but leaves allocated Data2D
func (d *Data2D) Clear() {
	d.Vertexes.Clear()
	d.Indices.Clear()
}

// Accept accepts vertex Data2D, this is only correct way of feeding batch with Vertexes
// along side indices, if you don't use indices append directly to Data
func (d *Data2D) Accept(Data2D VS2D, indices Indices) {
	l1 := len(d.Indices)
	l2 := uint32(d.Vertexes.Len())

	d.Vertexes = append(d.Vertexes, Data2D...)
	d.Indices = append(d.Indices, indices...)

	l3 := len(d.Indices)
	for i := l1; i < l3; i++ {
		d.Indices[i] += l2
	}
}


// Batch2D is main drawer, it performs direct draw to canvas and is used as target for Sprite.
// Batch2D acts like canvas i some ways but performance difference of drawing Batch2D to canvas and
// drawing canvas to canvas is significant. If you need image to ber redrawn ewer frame draw Batch2D
// to canvas and use canvas for drawing.
type Batch2D struct {
	Data2D

	buffer  *Buffer
	program *Program
	texture *Texture
}

// NBatch2D allows constructing Batch2D with custom Buffer and Program for applying
// per Batch2D shader and related buffer structure. Passing nil absolutely fine,
// as canvas or vindow will use theier own, if you don't even need texture use struct
// literal (Batch{}) to construct Batch2D
func NBatch2D(texture *Texture, buffer *Buffer, program *Program) *Batch2D {
	return &Batch2D{
		texture: texture,
		buffer:  buffer,
		program: program,
	}
}

// Draw draws all Data2D to target
func (b *Batch2D) Draw(target Target) {
	target.Accept(b.Vertexes, b.Indices, b.texture, b.program, b.buffer)
}

// Program returns Batch2D program, it can be nil
func (b *Batch2D) Program() *Program {
	return b.program
}

