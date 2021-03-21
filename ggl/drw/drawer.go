package drw

import (
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
//  d.Color(mat.Green).Thickness(5).Fill(false).AABB(mat.A(0, 0, 100, 100))
//
//  // draws line closed in triangle with custom edge style
//  d.Color(mat.RGB(0, 1, 1)).Loop(true).Thickness(10).Edge(CutEdge{})
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
	lineProc LineProcessor
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

/*// Edge sets line drawing edge style
func (g *Geom) Edge(edge Edge) *Geom {
	g.edge = edge
	return g
}*/

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

// Thickness sets line thickness of drawer
func (g *Geom) Thickness(thickness float64) *Geom {
	g.thickness = thickness
	return g
}

// Resolution sets resolution of circle, it can be set to Auto
// but if spacing is 0 nothing will be drawn
func (g *Geom) Resolution(resolution int) *Geom {
	g.resolution = resolution
	return g
}

// Spacing sets circle spacing, if you pass 0 old resolution will be used
// if you pass any positive number, resolution will be set to auto
func (g *Geom) Spacing(value float64) *Geom {
	if g.spacing == 0 {
		g.resolution = g.oldResolution
	} else {
		g.oldResolution = g.resolution
		g.resolution = Auto
	}

	g.spacing = value
	return g
}

func (g *Geom) Arc(start, end float64) *Geom {
	g.start = start
	g.end = end
	return g
}

func (g *Geom) LineType(ld LineDrawer) *Geom {
	g.lineDrawer = ld
	return g
}

func (g *Geom) Line(points ...mat.Vec) {
	g.lineProc.Process(g, g.lineDrawer, points...)
	g.Data.Accept(nil, g.lineProc.Indices)
	v := g.Reserve(len(g.lineProc.Points))
	for i, p := range g.lineProc.Points {
		v[i].Pos = p
	}
}

func (g *Geom) Circle(c mat.Circ) {
	resolution := g.resolution
	if resolution == Auto {
		resolution = AutoResolution(c.R, g.start, g.end, g.spacing)
	}
	var v ggl.Vertexes
	if g.fill {
		g.circle.Filled(1, g.start, g.end, resolution)
		g.Accept(nil, g.circle.Indices)
		g.circle.Vertexes, v = g.Reserve(len(g.circle.Vertexes)), g.circle.Vertexes
		g.circle.Update(mat.M(c.C, mat.V(c.R, c.R), 0), g.col)
	} else {
		g.circle.Outline(c.R, g.thickness, g.start, g.end, resolution)
		g.Accept(nil, g.circle.Indices)
		g.circle.Vertexes, v = g.Reserve(len(g.circle.Vertexes)), g.circle.Vertexes
		g.circle.Update(mat.M(c.C, mat.V(1, 1), 0), g.col)
	}
	g.circle.Vertexes = v
}

// AABB draws AABB appliable(Fill, Edge, Thickness)
func (g *Geom) AABB(value mat.AABB) {
	g.Rect(value.Vertices())
}

// Rect draws rectangle appliable(Fill, Edge, Thickness)
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

type geomCfg struct {
	col                       mat.RGBA
	loop, fill                bool
	resolution, oldResolution int

	thickness, intens, start, end, spacing float64

	lineDrawer LineDrawer
}

func nGeomCfg() geomCfg {
	return geomCfg{
		col:        mat.Alpha(1),
		thickness:  10,
		fill:       true,
		resolution: Auto,
		spacing:    1,
	}
}
