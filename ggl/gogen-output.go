package ggl


// Vertexes implements essential utility methods for any
// struct satisfying VertexData interface. Its a gogen TEMPLATE:
//
// 	/*gen(
//		Vertexes<Vertex, 9, Vertexes>
//	)*/
//
// this block generates VertexSlice with Vertex and 8 is the size of Vertex,
// Vertexes is name of generated struct divided by float64 byte size for more info
// search github.com/jakubDoka/gogen.
type Vertexes []Vertex

// Rewrite revrites elements from index to o values
func (v Vertexes) Rewrite(o Vertexes, idx int) {
	copy(v[idx:], o)
}

// Clear clears slice
func (v *Vertexes) Clear() {
	*v = (*v)[:0]
}

// Len implements VertexData interface
func (v Vertexes) Len() int {
	return len(v)
}

// VertexSize implements VertexData interface
func (v Vertexes) VertexSize() int {
	return 9
}


// Data is Vertex and indice collector, mainly utility that handles vertex offsets
// it also stores one aditionall slice as space for preporsessing
type Data struct {
	Vertexes Vertexes
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
func (d *Data) Accept(Data Vertexes, indices Indices) {
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

	Buffer  *Buffer
	Program *Program
	Texture *Texture
}

// Draw draws all Data to target
func (b *Batch) Draw(target Renderer) {
	target.Render(b.Vertexes, b.Indices, b.Texture, b.Program, b.Buffer)
}

