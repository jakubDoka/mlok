package drw

import (
	"math"

	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/mat"
	"github.com/jakubDoka/gobatch/mat/angle"
)

// AutoResolutionSpacing is size of fraction of the circle that has Auto resolution
var AutoResolutionSpacing float64 = 1

// Auto is constant that if you pass as resolution to Circle, auto resolution will be used
const Auto = -1

// Circle is a drawing tool that can draw circles with different
// resolution efficiently as everithing is precalculated upon creation
type Circle struct {
	ggl.Data
	Base
	resolution                    int
	radius, thickness, start, end float64
	outline                       bool
}

// NCircle creates ready-to-use Circle, choose resolution based of how big circles you want to draw
// bigger the circle more resolution matters, radius is not that important ac you can scale circle
// how ever you like
func NCircle(radius, thickness float64, resolution int) Circle {
	return NArc(radius, thickness, 0, 0, resolution)
}

func NArc(radius, thickness, start, end float64, resolution int) Circle {
	c := Circle{}
	if thickness == 0 {
		c.Filled(radius, start, end, resolution)
	} else {
		c.Outline(radius, thickness, start, end, resolution)
	}
	return c
}

func (c *Circle) Outline(radius, thickness, start, end float64, resolution int) {
	if !c.outline || radius != c.radius || thickness != c.thickness ||
		resolution != c.resolution || start != c.start || end != c.end {
		c.Base.Resize(resolution * 2)
		ang, step := c.setup(start, end, resolution)
		for i := 0; i < len(c.Base); i, ang = i+2, ang+step {
			c.Base[i] = mat.Rad(ang, radius+thickness)
			c.Base[i+1] = mat.Rad(ang, radius-thickness)
		}
		if start != end {
			c.Base = append(c.Base, mat.Rad(ang, radius+thickness), mat.Rad(ang, radius-thickness))
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

		if start == end {
			c.Indices[l-4] = 0
			c.Indices[l-2] = 1
			c.Indices[l-1] = 0
		}

		c.resolution = resolution
	}

	c.outline = true
}

func (c *Circle) Filled(radius, start, end float64, resolution int) {
	if c.outline || radius != c.radius || resolution != c.resolution {
		c.Base.Resize(resolution + 1)
		ang, step := c.setup(start, end, resolution)
		for i := 1; i < len(c.Base); i, ang = i+1, ang+step {
			c.Base[i] = mat.Rad(ang, radius)
		}
		if start != end {
			c.Base = append(c.Base, mat.Rad(ang, radius))
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

		if start == end {
			c.Indices[l-1] = 1
		} else {
			c.Indices = c.Indices[:l-3]
		}

		c.resolution = resolution
	}

	c.outline = false
}

func AutoResolution(radius, start, end, spacing float64) int {
	if start == end {
		end = start + angle.Pi2
	}
	cof := math.Abs(start-end) / angle.Pi2

	return mat.Maxi(int(radius*angle.Pi2*cof/spacing), 3)
}

func (c *Circle) setup(start, end float64, resolution int) (s, step float64) {
	if start == end {
		end = start + angle.Pi2
	} else if start > end {
		start, end = end, start
	}

	c.start = start
	c.end = end

	return start, (end - start) / float64(resolution)
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
