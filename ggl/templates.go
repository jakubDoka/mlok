package ggl

const elementsize = 0

/*imp(
	github.com/jakubDoka/gogen/templates
)*/

/*gen(
	vertexSlice<Vertex, 9, Vertexes>
	data<Vertexes, Data>
	batch<Data, NBatch, Batch>
	templates.Resize<Indices, Resize>
	templates.Resize<Vertexes, Resize>
)*/

//def(
//rules vertexSlice<interface{}, elementsize>

// vertexSlice implements essential utility methods for any
// struct satisfying VertexData interface. Its a gogen TEMPLATE:
//
// 	/*gen(
//		vertexSlice<Vertex, 9, Vertexes>
//	)*/
//
// this block generates VertexSlice with Vertex and 8 is the size of Vertex,
// Vertexes is name of generated struct divided by float64 byte size for more info
// search github.com/jakubDoka/gogen.
type vertexSlice []interface{}

// Rewrite revrites elements from index to o values
func (v vertexSlice) Rewrite(o vertexSlice, idx int) {
	copy(v[idx:], o)
}

// Clear clears slice
func (v *vertexSlice) Clear() {
	*v = (*v)[:0]
}

// Len implements VertexData interface
func (v vertexSlice) Len() int {
	return len(v)
}

// VertexSize implements VertexData interface
func (v vertexSlice) VertexSize() int {
	return elementsize
}

//)

//def(
//rules data<vertexSlice>

// data is Vertex and indice collector, mainly utility that handles vertex offsets
// it also stores one aditionall slice as space for preporsessing
type data struct {
	Vertexes vertexSlice
	Indices  Indices
}

// Copy copies data to another resulting into two deeply equal objects
func (d *data) Copy(dst *data) {
	dst.Clear()
	dst.Indices = append(dst.Indices, d.Indices...)
	dst.Vertexes = append(dst.Vertexes, d.Vertexes...)
}

// Clear clears t batch but leaves allocated data
func (d *data) Clear() {
	d.Vertexes.Clear()
	d.Indices.Clear()
}

// Accept accepts vertex data, this is only correct way of feeding batch with Vertexes
// along side indices, if you don't use indices append directly to Data
func (d *data) Accept(data vertexSlice, indices Indices) {
	l1 := len(d.Indices)
	l2 := uint32(d.Vertexes.Len())

	d.Vertexes = append(d.Vertexes, data...)
	d.Indices = append(d.Indices, indices...)

	if l2 != 0 {
		l3 := len(d.Indices)
		for i := l1; i < l3; i++ {
			d.Indices[i] += l2
		}
	}
}

//)

//def(
//rules batch<data, nbatch>

// batch is main drawer, it performs direct draw to canvas and is used as target for Sprite.
// batch acts like canvas i some ways but performance difference of drawing batch to canvas and
// drawing canvas to canvas is significant. If you need image to ber redrawn ewer frame draw batch
// to canvas and use canvas for drawing.
type batch struct {
	data

	Buffer  *Buffer
	Program *Program
	Texture *Texture
}

// Draw draws all data to target
func (b *batch) Draw(target Renderer) {
	target.Render(b.Vertexes, b.Indices, b.Texture, b.Program, b.Buffer)
}

//)
