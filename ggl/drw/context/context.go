package context

import (
	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/rgba"
)

// C helps you draw complex objects built from multiple sprites
type C []Part

// Init Initialized context with given Defaults. This can be called multiple times on same context
// and it will be restarted. C is mant to be reused to reduce allocations.
func (c *C) Init(parts ...PartDefs) {
	v := *c

	v = v[:0]

	for len(parts) > len(v) {
		v = append(v, Part{
			Def:   parts[len(v)],
			Spr:   ggl.NSprite(parts[len(v)].Region),
			Mask:  rgba.White,
			Scale: mat.V(1, 1),
		})
	}

	*c = v
}

// Draw draws context to target with applied transform and mask
func (c C) Draw(t ggl.Target, matrix mat.Mat, mask mat.RGBA) {
	c.Update(matrix, mask)
	c.Fetch(t)
}

// Update updates the sprite transform of context
func (c C) Update(matrix mat.Mat, mask mat.RGBA) {
	for i := range c {
		p := &c[i]
		p.Spr.Update(
			mat.M(
				p.Offset.Add(p.Def.Offset),
				p.Def.Scale.Mul(p.Scale),
				p.Rotation+p.Def.Rotation,
			).Chained(matrix),
			mask.Mul(p.Mask).Mul(p.Def.Mask),
		)
	}
}

func (c C) Fetch(t ggl.Target) {
	for i := range c {
		c[i].Spr.Fetch(t)
	}
}

// Part is a building piece of context, if contains lot of configuration that is combined with
// Default configuration to make a final transformation and color
type Part struct {
	Def           PartDefs
	Spr           ggl.Sprite
	Offset, Scale mat.Vec
	Mask          mat.RGBA
	Rotation      float64
}

// TotalOffset returns total offset of part, taking PartDefs into account
func (p *Part) TotalOffset() mat.Vec {
	return p.Offset.Add(p.Def.Offset)
}

// TotalMask returns total mask of part, taking PartDefs into account
func (p *Part) TotalMask() mat.RGBA {
	return p.Mask.Mul(p.Def.Mask)
}

// TotalRotation returns total rotation of part, taking PartDefs into account
func (p *Part) TotalRotation() float64 {
	return p.Rotation + p.Def.Rotation
}

// PartDefs stores the default values for Part
type PartDefs struct {
	Offset, Pivot, Scale mat.Vec
	Rotation             float64
	Mask                 mat.RGBA
	Region               mat.AABB
}

var DefaultPartDefs = PartDefs{
	Scale: mat.V(1, 1),
	Mask:  rgba.White,
}
