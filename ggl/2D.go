package ggl

import (
	"gobatch/mt"
)

// Sprite2DTarget is something sprite can draw on to
type Sprite2DTarget interface {
	Accept(data VS2D, indices Indices)
}

// SpriteIndices is slice of constant indices used by Sprite
var SpriteIndices = Indices{0, 1, 2, 0, 3, 2}

// Sprite2D is most efficient way of drawing textures to Batch2D (if you find faster way i welcome your pr)
// sprite does not allocate any memory all data is on stack, its designed to be easily copied by value.
type Sprite2D struct {
	tex  [4]mt.V2
	data [4]Vertex2D
}

// NSprite2D creates new sprite out of frame. Frame should be the region where the texture
// you want to draw is
func NSprite2D(frame mt.AABB) Sprite2D {
	s := Sprite2D{}
	vert := frame.Vertices()

	c := frame.Center()

	for i, v := range vert {
		s.data[i].Tex = v
		s.data[i].Color = mt.Alpha(1)
		s.tex[i] = v.Sub(c)
	}

	return s
}

// Draw draws sprite to Batch projected by given matrix, and colored by mask
func (s *Sprite2D) Draw(b Sprite2DTarget, mat mt.Mat2, mask mt.RGBA) {
	s.Update(mat, mask)
	b.Accept(s.data[:], SpriteIndices)
}

// LazyDraw draws sprite as it is, to change draw result call Update
func (s *Sprite2D) LazyDraw(b Sprite2DTarget) {
	b.Accept(s.data[:], SpriteIndices)
}

// Update only updates sprite data but does not draw it
func (s *Sprite2D) Update(mat mt.Mat2, mask mt.RGBA) {
	for i := range s.data {
		s.data[i].Pos = mat.Project(s.tex[i])
		s.data[i].Color = mask
	}
}
