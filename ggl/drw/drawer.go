package drw

import (
	"math"

	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/mat"
)

// Geom is is small abstraction around Data that allows
// drawing of geometric shapes, geomdrawer can be used as canvas for creating
// complex shapes, Drawers can draw to each other, though if you don't want
// triangles to get modified by target drawer use drawer.Data as target
//
//  // use as is or create with NGeom that sets some default values
//  d := drw.Geom{}
//
//  // red rectangle
//  d.Color(mat.Red).Fill(true).AABB(mat.A(0, 0, 100, 100))
//
//  // draws green outline to red rectangel with total thickness 10
//  d.Color(mat.Green).Width(5).Fill(false).AABB(mat.A(0, 0, 100, 100))
//
//  // draws line closed in triangle with custom edge style
//  d.Color(mat.RGB(0, 1, 1)).Loop(true).Width(10).Edge(CutEdge{})
//  d.Line(mat.V(0, -100), mat.V(-100, -200), mat.V(100, -200))
//
//  // draw everithing to target with no transformation
//  d.Fetch(t)
//
//  // draw to target with ewerithing transformed by matrix and masked, in this case
//	// everithing is shifted by 100 to right, twice as big and rotated by 90 degrees
//	d.Draw(t, mat.M(mat.V(100, 0), mat.V(2, 2), math.Phi*.5), mat.RGB(.5, .5, .5))
type Geom struct {
	ggl.Data
	tmp ggl.Data
	geomCfg
	convexes []bool
	circle   Circle
}

// NGeomDrawer sets some nice default values
func NGeomDrawer() Geom {
	return Geom{geomCfg: nGeomCfg()}
}

// Restart restarts the configuration to default one
func (g *Geom) Restart() {
	g.Clear()
	g.geomCfg = nGeomCfg()
}

// ExampleGD show example use of geom drawer
func ExampleGD(t ggl.Target) {
	// use as is or create with NGeom that sets some default values
	d := Geom{}

	// red rectangle
	d.Color(mat.Red).Fill(true).AABB(mat.A(0, 0, 100, 100))

	// green outline for red rectangel with total thickness 10
	d.Color(mat.Green).Width(5).Fill(false).AABB(mat.A(0, 0, 100, 100))

	// line closed in triangle with custom edge style
	d.Color(mat.RGB(0, 1, 1)).Loop(true).Width(10).Edge(CutEdge{})
	d.Line(mat.V(0, -100), mat.V(-100, -200), mat.V(100, -200))

	// draw everithing to target with no transformation
	d.Fetch(t)

	// draw to target with ewerithing transformed by matrix and masked, in this case
	// everithing is shifted by 100 to right, twice as big and rotated by 90 degrees
	d.Draw(t, mat.M(mat.V(100, 0), mat.V(2, 2), math.Pi*.25), mat.RGB(.5, .5, .5))
}

// Accept implements ggl.Target interface
func (g *Geom) Accept(vertexes ggl.Vertexes, indices ggl.Indices) {
	start := g.Vertexes.Len()
	g.Data.Accept(vertexes, indices)
	g.Apply(start, g.Vertexes.Len())
}

// Draw implements ggl.Drawer interface
func (g *Geom) Draw(t ggl.Target, mat mat.Mat, rgba mat.RGBA) {
	g.tmp.Vertexes = g.tmp.Vertexes[:0]
	g.tmp.Vertexes = append(g.tmp.Vertexes, g.Vertexes...)

	for i := range g.tmp.Vertexes {
		g.tmp.Vertexes[i].Color = g.tmp.Vertexes[i].Color.Mul(rgba)
		g.tmp.Vertexes[i].Pos = mat.Project(g.tmp.Vertexes[i].Pos)
	}

	t.Accept(g.tmp.Vertexes, g.Indices)
}

// Fetch implements ggl.Fetcher interface
func (g *Geom) Fetch(t ggl.Target) {
	t.Accept(g.Vertexes, g.Indices)
}

// Color sets drawind color
func (g *Geom) Color(value mat.RGBA) *Geom {
	g.col = value
	return g
}

// Intensity sets drawing intensity
func (g *Geom) Intensity(value float64) *Geom {
	g.intens = value
	return g
}

// Edge sets line drawing edge style
func (g *Geom) Edge(edge Edge) *Geom {
	g.edge = edge
	return g
}

// Loop sets whether lines should be closed into loops
func (g *Geom) Loop(loop bool) *Geom {
	g.loop = loop
	return g
}

// Fill decides whether shapes should be filled or just outlines
func (g *Geom) Fill(fill bool) *Geom {
	g.fill = fill
	return g
}

// Width sets line width of drawer
func (g *Geom) Width(width float64) *Geom {
	g.width = width
	return g
}

// Resolution sets resolution of circle
func (g *Geom) Resolution(resolution int) *Geom {
	g.resolution = resolution
	return g
}

// AABB draws AABB appliable(Fill, Edge, Width)
func (g *Geom) AABB(value mat.AABB) {
	g.Rect(value.Vertices())
}

// Rect draws rectangle appliable(Fill, Edge, Width)
func (g *Geom) Rect(corners [4]mat.Vec) {
	if g.fill {
		g.Accept(nil, ggl.SpriteIndices)
		vs := g.Reserve(4)

		for i := range vs {
			vs[i].Pos = corners[i]
		}
	} else {
		loop := g.loop
		g.loop = true
		g.Line(corners[:]...)
		g.loop = loop
	}
}

// Reserve reserves vertexes, sets theier intensity and color and returns slice that points to them
func (g *Geom) Reserve(amount int) ggl.Vertexes {
	ol := len(g.Vertexes)
	l := ol + amount
	if cap(g.Vertexes) >= l {
		g.Vertexes = g.Vertexes[:l]
	} else {
		nv := make(ggl.Vertexes, l)
		copy(nv, g.Vertexes)
		g.Vertexes = nv
	}

	g.Apply(ol, l)

	return g.Vertexes[ol:l]
}

// Apply applies curent color and intensity to slice of vertices
func (g *Geom) Apply(start, end int) {
	for i := start; i < end; i++ {
		g.Vertexes[i].Color = g.col
		g.Vertexes[i].Intensity = g.intens
	}
}

// Line creates line based of configuration
func (g *Geom) Line(points ...mat.Vec) {
	g.convexes = g.convexes[:0]

	e := g.edge
	if e == nil {
		e = EdgeBase{}
	}

	var (
		vl   = len(g.Vertexes)
		l    = len(points)
		size = e.Size(l, g.loop)
		vs   = g.Reserve(size)
		edge EdgeData
	)

	if !g.loop {
		l--
	}

	for i := 0; i < l; i++ {
		edge.Init(i, points, g.width)
		g.convexes = append(g.convexes, edge.Convex)
		e.Process(&edge, vs, i, size)
	}

	e.Indices(l, vl, g.loop, &g.Indices, g.convexes)
}

func (g *Geom) Circle(c mat.Circ) {

	var v ggl.Vertexes

	if g.fill {
		g.circle.Filled(1, g.resolution)
		g.Accept(nil, g.circle.Indices)
		g.circle.Vertexes, v = g.Reserve(len(g.circle.Vertexes)), g.circle.Vertexes
		g.circle.Update(mat.M(c.C, mat.V(c.R, c.R), 0), g.col)
	} else {
		g.circle.Outline(c.R, g.width, g.resolution)
		g.Accept(nil, g.circle.Indices)
		g.circle.Vertexes, v = g.Reserve(len(g.circle.Vertexes)), g.circle.Vertexes
		g.circle.Update(mat.M(c.C, mat.V(1, 1), 0), g.col)
	}

	g.circle.Vertexes = v
}

// Edge ...
type Edge interface {
	Size(size int, loop bool) int
	Process(e *EdgeData, vertexes ggl.Vertexes, idx, size int)
	Indices(size, shift int, loop bool, buff *ggl.Indices, convexes []bool)
}

// EdgeData stores information about line edge that is processed
type EdgeData struct {
	A, B, C       mat.Vec
	Prev, Segment [4]mat.Vec
	Convex        bool
}

// Init initializes line edge, its to avoid allocation and data moving, edge is
// created once and initted multiple times for each interation
func (l *EdgeData) Init(i int, points []mat.Vec, width float64) {
	ln := len(points)
	l.A, l.B, l.C = points[i], points[(i+1)%ln], points[(i+2)%ln]
	l.Convex = l.B.To(l.A).Cross(l.B.To(l.C)) > 0
	l.Prev = l.Segment
	vec := l.A.To(l.B).Norm(width)
	l.Segment = [4]mat.Vec{l.A.Add(vec), l.A.Sub(vec), l.B.Sub(vec), l.B.Add(vec)}
}

// CutEdge ...
type CutEdge struct {
	EdgeBase
}

// Indices implements Edge interface
func (c CutEdge) Indices(size, shift int, loop bool, buff *ggl.Indices, convexes []bool) {
	var (
		es = uint32(c.EdgeSize())
		b  = *buff
		l  = ggl.SpriteIndicesSize + 3
		j  = uint32(shift)
		k  = len(*buff)
	)

	for i := 0; i < size; i++ {
		b = append(b, ggl.SpriteIndices...)
		if convexes[i] {
			b = append(b, 2, 3, 4)
		} else {
			b = append(b, 2, 3, 5)
		}

		b[k:].Shift(j)
		j += es
		k += l
	}

	if loop {
		b[k-1] -= j - uint32(shift)
	} else {
		b = b[:k-3]
	}

	*buff = b
}

// EdgeBase ...
type EdgeBase struct{}

// Size implements Edge interface
func (c EdgeBase) Size(i int, loop bool) int {
	var a = c.EdgeSize()
	if loop {
		a = 0
	}
	return i*c.EdgeSize() - a
}

// Process implements Edge interface
func (c EdgeBase) Process(e *EdgeData, vertexes ggl.Vertexes, idx, _ int) {
	idx *= c.EdgeSize()
	for i, v := range e.Segment {
		vertexes[idx+i].Pos = v
	}
}

// Indices implements Edge interface
func (c EdgeBase) Indices(size, shift int, loop bool, buff *ggl.Indices, convexes []bool) {
	es := uint32(c.EdgeSize())
	b := *buff

	var j = uint32(shift)
	var k = len(*buff)
	for i := 0; i < size; i++ {
		b = append(b, ggl.SpriteIndices...)
		b[k:].Shift(j)
		j += es
		k += ggl.SpriteIndicesSize
	}

	*buff = b
}

// EdgeSize returns size of one edge of cut Edge
func (c EdgeBase) EdgeSize() int {
	return 4
}

type geomCfg struct {
	col                   mat.RGBA
	edge                  Edge
	loop, fill            bool
	resolution            int
	width, intens, radius float64
}

func nGeomCfg() geomCfg {
	return geomCfg{
		col:   mat.Alpha(1),
		width: 10,
		fill:  true,
	}
}
