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
	Base
	resolution        int
	radius, thickness float64
	outline           bool
}

// NCircle creates ready-to-use Circle, choose resolution based of how big circles you want to draw
// bigger the circle more resolution matters, radius is not that important ac you can scale circle
// how ever you like
func NCircle(radius, thickness float64, resolution int) Circle {
	// setting up vectors
	c := Circle{}
	if thickness == 0 {
		c.Filled(radius, resolution)
	} else {
		c.Outline(radius, thickness, resolution)
	}
	return c
}

func (c *Circle) Outline(radius, thickness float64, resolution int) {
	if !c.outline || radius != c.radius || thickness != c.thickness || resolution != c.resolution {
		c.Base.Resize(resolution * 2)
		frac := math.Pi * 2 / float64(resolution)
		for i, ang := 0, 0.0; i < len(c.Base); i, ang = i+2, ang+frac {
			c.Base[i] = mat.Rad(ang, radius+thickness)
			c.Base[i+1] = mat.Rad(ang, radius-thickness)
		}
		c.radius = radius
		c.thickness = thickness
		c.Vertexes.Resize(len(c.Base))
	}

	if !c.outline || resolution != c.resolution {
		// setting up indices
		l := resolution * 6
		c.Indices.Resize(l)
		ln := uint32(l)
		for i, j := uint32(0), uint32(0); i < ln; i, j = i+6, j+2 {
			c.Indices[i+0] = j
			c.Indices[i+1] = j + 1
			c.Indices[i+2] = j + 2
			c.Indices[i+3] = j + 1
			c.Indices[i+4] = j + 3
			c.Indices[i+5] = j + 2
		}
		c.Indices[l-4] = 0
		c.Indices[l-2] = 1
		c.Indices[l-1] = 0

		c.resolution = resolution
	}

	c.outline = true
}

func (c *Circle) Filled(radius float64, resolution int) {
	if c.outline || radius != c.radius || resolution != c.resolution {
		c.Base.Resize(resolution + 1)
		frac := math.Pi * 2 / float64(resolution)
		for i, ang := 1, 0.0; i < len(c.Base); i, ang = i+1, ang+frac {
			c.Base[i] = mat.Rad(ang, radius)
		}
		c.radius = radius
		c.Vertexes.Resize(len(c.Base))
	}

	if c.outline || resolution != c.resolution {
		// setting up indices
		l := resolution * 3
		c.Indices.Resize(l)
		res := uint32(resolution)
		for i := uint32(0); i < res; i++ {
			j := i * 3
			c.Indices[j+1] = i + 1
			c.Indices[j+2] = i + 2
		}
		c.Indices[l-1] = 1
		c.resolution = resolution
	}

	c.outline = false
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

type Base []mat.Vec

/*imp(
	github.com/jakubDoka/gogen/templates
)*/

/*gen(
	templates.Resize<Base, Resize>
)*/
