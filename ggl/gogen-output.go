package ggl

// VS implements essential utility methods for any
// struct satisfying VertexData interface. Its a gogen TEMPLATE:
//
// 	/*gen(
//		VS<Vertex, 9, VS>
//	)*/
//
// this block generates VertexSlice with Vertex and 8 is the size of Vertex,
// VS is name of generated struct divided by float64 byte size for more info
// search github.com/jakubDoka/gogen.
type VS []Vertex

// Rewrite revrites elements from index to o values
func (v VS) Rewrite(o VS, idx int) {
	copy(v[idx:], o)
}

// Clear clears slice
func (v *VS) Clear() {
	*v = (*v)[:0]
}

// Len implements VertexData interface
func (v VS) Len() int {
	return len(v)
}

// VertexSize implements VertexData interface
func (v VS) VertexSize() int {
	return 9
}

// Data is Vertex and indice collector, mainly utility that handles vertex offsets
// it also stores one aditionall slice as space for preporsessing
type Data struct {
	Vertexes VS
	Indices  Indices
}

// Copy copies Data to another resulting into two deeply equal objects
func (d *Data) Copy(dst *Data) {
	dst.Clear()
	dst.Indices = append(dst.Indices, d.Indices...)
	dst.Vertexes = append(dst.Vertexes, d.Vertexes...)
}

// Clear clears t batch but leaves allocated Data
func (d *Data) Clear() {
	d.Vertexes.Clear()
	d.Indices.Clear()
}

// Accept accepts vertex Data, this is only correct way of feeding batch with Vertexes
// along side indices, if you don't use indices append directly to Data
func (d *Data) Accept(Data VS, indices Indices) {
	l1 := len(d.Indices)
	l2 := uint32(d.Vertexes.Len())

	d.Vertexes = append(d.Vertexes, Data...)
	d.Indices = append(d.Indices, indices...)

	l3 := len(d.Indices)
	for i := l1; i < l3; i++ {
		d.Indices[i] += l2
	}
}

// Batch is main drawer, it performs direct draw to canvas and is used as target for Sprite.
// Batch acts like canvas i some ways but performance difference of drawing Batch to canvas and
// drawing canvas to canvas is significant. If you need image to ber redrawn ewer frame draw Batch
// to canvas and use canvas for drawing.
type Batch struct {
	Data

	buffer  *Buffer
	program *Program
	texture *Texture
}

// NBatch allows constructing Batch with custom Buffer and Program for applying
// per Batch shader and related buffer structure. Passing nil absolutely fine,
// as canvas or vindow will use theier own, if you don't even need texture use struct
// literal (Batch{}) to construct Batch
func NBatch(texture *Texture, buffer *Buffer, program *Program) *Batch {
	return &Batch{
		texture: texture,
		buffer:  buffer,
		program: program,
	}
}

// Draw draws all Data to target
func (b *Batch) Draw(target RenderTarget) {
	target.Accept(b.Vertexes, b.Indices, b.texture, b.program, b.buffer)
}

// Program returns Batch program, it can be nil
func (b *Batch) Program() *Program {
	return b.program
}
