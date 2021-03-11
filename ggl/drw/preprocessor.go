package drw

import (
	"gobatch/ggl"
	"gobatch/mat"
)

// SpriteViewport is able to crop sprite triangles that are not rotated
// thus i makes effect like sprites are only visible from Area, if you
// try to clamp rotated sprites they will get deformed, trying to clamp
// something other then sprites will also result in random deformation
// of corse vertices will be visible only inside the View but for
// unreasonable performance cost.
type SpriteViewport struct {
	ggl.Data
	Area mat.AABB
}

// Accept implements ggl.Target interface
func (p *SpriteViewport) Accept(vertexes ggl.Vertexes, indices ggl.Indices) {
	var (
		count int
		prev  = len(p.Vertexes)
	)
	for i := 0; i < len(vertexes); i += 4 {
		vs := [4]ggl.Vertex{}
		copy(vs[:], vertexes[i:i+4])

		area := mat.AABB{Min: vs[0].Pos, Max: vs[2].Pos}
		if !area.Intersects(p.Area) {
			continue
		}

		count += ggl.SpriteIndicesSize

		in := area.Intersect(p.Area)
		if in == area {
			p.Vertexes = append(p.Vertexes, vs[:]...)
			continue
		}

		for i, v := range in.Vertices() {
			vs[i].Pos = v
		}

		var (
			cet = area.Center()
			sz  = area.Size().Scaled(.5)
			min = in.Min.Sub(cet).Div(sz)
			max = in.Max.Sub(cet).Div(sz)

			tex  = mat.AABB{Min: vs[0].Tex, Max: vs[2].Tex}
			tCet = tex.Center()
			tSz  = tex.Size().Scaled(.5)
		)

		tex.Min.X = tCet.X + min.X*tSz.X
		tex.Max.X = tCet.X + max.X*tSz.X
		tex.Min.Y = tCet.Y + min.Y*tSz.Y
		tex.Max.Y = tCet.Y + max.Y*tSz.Y

		for i, v := range tex.Vertices() {
			vs[i].Tex = v
		}

		p.Vertexes = append(p.Vertexes, vs[:]...)
	}

	ip := len(p.Indices)
	p.Indices = append(p.Indices, indices[:count]...)
	p.Indices[ip:].Shift(uint32(prev))
}

// ClampedViewport clamps all vertices into Area
type ClampedViewport struct {
	ggl.Data
	Area mat.AABB
}

// Accept implements ggl.Target interface
func (p *ClampedViewport) Accept(vertexes ggl.Vertexes, indices ggl.Indices) {
	ln := p.Vertexes.Len()
	p.Data.Accept(vertexes, indices)
	for i := ln; i < p.Vertexes.Len(); i++ {
		p.Vertexes[i].Pos = p.Area.Clamp(p.Vertexes[i].Pos)
	}
}

// Preprocessor is something that changes triangles in some way
// and can draw them, for example multiplying of seting color mask
type Preprocessor interface {
	ggl.Fetcher
	ggl.Target
	Clear()
}
