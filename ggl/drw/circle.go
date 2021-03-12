package drw

import (
	"math"

	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/mat"
)

// Circle is a drawing tool that can draw circles with different
// resolution efficiently as everithing is precalculated upon creation
type Circle struct {
	ggl.Data
	Base []mat.Vec
}

// NCircle creates ready-to-use Circle, choose resolution based of how big circles you want to draw
// bigger the circle more resolution matters, radius is not that important ac you can scale circle
// how ever you like
func NCircle(radius float64, resolution int) Circle {
	// setting up vectors
	c := Circle{}
	c.Base = make([]mat.Vec, resolution+1)
	frac := math.Pi * 2 / float64(resolution)
	for i, ang := 1, 0.0; i < len(c.Base); i, ang = i+1, ang+frac {
		c.Base[i] = mat.Rad(ang, radius)
	}

	// setting up indices
	l := resolution * 3
	c.Indices = make(ggl.Indices, l)
	res := uint32(resolution)
	for i := uint32(0); i < res; i++ {
		j := i * 3
		c.Indices[j+1] = i + 1
		c.Indices[j+2] = i + 2
	}
	c.Indices[l-1] = 1

	// setup vertexes
	c.Vertexes = make(ggl.Vertexes, len(c.Base))

	return c
}

// Draw implements ggl.Drawer interface
func (c *Circle) Draw(t ggl.Target, tran mat.Mat, rgba mat.RGBA) {
	c.Update(tran, rgba)
	c.Fetch(t)
}

// Update updates circle state to given transformation and color
func (c *Circle) Update(tran mat.Mat, rgba mat.RGBA) {
	for i, v := range c.Base {
		c.Vertexes[i].Pos = tran.Project(v)
		c.Vertexes[i].Color = rgba
	}
}
