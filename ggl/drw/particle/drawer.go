package particle

import (
	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/ggl/drw"
	"github.com/jakubDoka/gobatch/mat"
)

// Sprite extends ggl.Sprite to implement Drawer
type Sprite struct {
	ggl.Sprite
}

// Square creates a square particle drawer
func Square(size float64) *Sprite {
	spr := ggl.NSprite(mat.A(0, 0, size, size))
	spr.SetIntensity(0)
	return &Sprite{spr}
}

// Copy implements Drawer interface
func (s *Sprite) Copy() Drawer {
	cp := *s
	return &cp
}

// Metrics implements Drawer interface
func (s *Sprite) Metrics() (indices, vertexes int) {
	return ggl.SpriteIndicesSize, ggl.SpriteVertexSize
}

// Circle extends dw circle to implement the Drawer interface
type Circle struct {
	drw.Circle
}

// Copy implements Drawer interface
func (c *Circle) Copy() Drawer {
	cp := Circle{}
	cp.Vertexes = append(cp.Vertexes, c.Vertexes...)
	cp.Indices = append(cp.Indices, c.Indices...)
	cp.Base = append(cp.Base, c.Base...)

	return &cp
}

// Metrics implements Drawer interface
func (c *Circle) Metrics() (indices, vertexes int) {
	return len(c.Indices), len(c.Vertexes)
}

// Drawer is a particle drawer, core is a ggl.Drawer but it also has to be able to copy it self
// and provide how math indices and vertexes one draw takes
type Drawer interface {
	ggl.Drawer
	Copy() Drawer
	Metrics() (indices, vertexes int)
}
