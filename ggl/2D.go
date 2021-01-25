package ggl

import (
	"gobatch/mat"
)

// Target2D is something sprite can draw on to
type Target2D interface {
	Accept(data VS2D, indices Indices)
}

// Drawer2D is triangle drawer oriented to 2D targets
type Drawer2D interface {
	Draw(b Target2D, mat mat.Mat2, mask mat.RGBA)
}

// Vertex2D is essentia vertex struct for 2D rendering
type Vertex2D struct {
	Pos, Tex  mat.V2
	Color     mat.RGBA
	Intensity float64
}

// Sprite related constants
const (
	SpriteVertexSize  = 4
	SpriteIndicesSize = 6
)

// SpriteIndices is slice of constant indices used by Sprite
var SpriteIndices = Indices{0, 1, 2, 0, 3, 2}

// Sprite2D is most efficient way of drawing textures to Batch2D (if you find faster way i welcome your pr)
// sprite does not allocate any memory all data is on stack, its designed to be easily copied by value.
type Sprite2D struct {
	tex  [4]mat.V2
	data [4]Vertex2D
}

// NSprite2D creates new sprite out of frame. Frame should be the region where the texture
// you want to draw is
func NSprite2D(frame mat.AABB) Sprite2D {
	s := Sprite2D{}
	vert := frame.Vertices()

	c := frame.Center()

	for i, v := range vert {
		s.data[i].Tex = v
		s.data[i].Color = mat.Alpha(1)
		s.data[i].Intensity = 1
		s.tex[i] = v.Sub(c)
	}

	return s
}

// SetIntensity sets the intensity of sprite image, if you set it to 0 the rectangle in a color of sprite mask will
// be drawn, if you set it to 1 (which is default) it will draw texture as it is or black area if batch does not
// have texture.
func (s *Sprite2D) SetIntensity(value float64) {
	for i := range s.data {
		s.data[i].Intensity = value
	}
}

// Clear makes sprite invisible when its drawn, this is to impelemt paceholder data
func (s *Sprite2D) Clear() {
	for i := range s.data {
		s.data[i].Pos = mat.V2{}
	}
}

// Set sets sprites source texture region and destination rectangle, this is mainly used when drawing text
func (s *Sprite2D) Set(dst, src mat.AABB) {
	tex, pos := dst.Vertices(), src.Vertices()
	for i := range s.data {
		s.data[i].Tex = tex[i]
		s.data[i].Pos = pos[i]
	}
}

// Draw draws sprite to Batch projected by given matrix, and colored by mask
func (s *Sprite2D) Draw(b Target2D, mat mat.Mat2, mask mat.RGBA) {
	s.Update(mat, mask)
	b.Accept(s.data[:], SpriteIndices)
}

// LazyDraw draws sprite as it is, to change draw result call Update
func (s *Sprite2D) LazyDraw(b Target2D) {
	b.Accept(s.data[:], SpriteIndices)
}

// Update only updates sprite data but does not draw it
func (s *Sprite2D) Update(mat mat.Mat2, mask mat.RGBA) {
	for i := range s.data {
		s.data[i].Pos = mat.Project(s.tex[i])
		s.data[i].Color = mask
	}
}
