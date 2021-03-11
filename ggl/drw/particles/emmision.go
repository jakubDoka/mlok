package particles

import (
	"gobatch/mat"
	"math/rand"
)

// Point is most basic emission shape, all particles will get spawned o a same position
type Point struct{}

// Gen implements EmissionShape interface
func (p Point) Gen(t float64) mat.Vec {
	return mat.ZV
}

// Circular spawns particles in specified arc, possibly circle
type Circular struct {
	Radius, Spread float64
}

// Gen implements EmissionShape interface
func (c Circular) Gen(dir float64) mat.Vec {
	return mat.Rad(c.Spread*2*rand.Float64()-c.Spread+dir, c.Radius*rand.Float64())
}

// Rectangle spans particles in a rectangle area
type Rectangle struct {
	Width, Height float64
}

// Gen implements EmissionShape interface
func (r Rectangle) Gen(dir float64) mat.Vec {
	return mat.V(r.Width*rand.Float64()-r.Width*.5, r.Height*rand.Float64()-r.Height*.5).Rotated(dir)
}

// EmissionShape is a vector generator, usefull for particle system (determining particle spawn position)
type EmissionShape interface {
	Gen(direction float64) mat.Vec
}
