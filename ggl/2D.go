package ggl

import (
	"github.com/jakubDoka/mlok/mat"
)

// Fetcher is something that only passes triangle data, and does no preprocessing
// use Fetch if you don't need to modify triangles, its faster for all structs that
// implement this interface
type Fetcher interface {
	Fetch(b Target)
}

// Target is something that accepts triangle data, but data is just copied,
// is is not used for anything, though shifting of indices is Targets roles
type Target interface {
	Accept(vertexes Vertexes, indices Indices)
}

// Drawer is triangle drawer, it should always preporcess triangles with given
// matrix and color and then give them to target
type Drawer interface {
	Draw(b Target, mat mat.Mat, mask mat.RGBA)
}

// Updater si something that can update its state based of color and transformation
type Updater interface {
	Update(mat mat.Mat, mask mat.RGBA)
}

// Vertex is essentia vertex struct for rendering
type Vertex struct {
	Pos, Tex  mat.Vec
	Color     mat.RGBA
	Intensity float64
}

// Fetch implements Fetcher interface
func (d *Data) Fetch(t Target) {
	t.Accept(d.Vertexes, d.Indices)
}

// Sprite related constants
const (
	SpriteVertexSize  = 4
	SpriteIndicesSize = 6
	NinePatchSide     = 3
)

// SpriteIndices is slice of constant indices used by Sprite
var SpriteIndices = Indices{0, 1, 2, 0, 3, 2}

// Patch consists of grid layout of 3x3 sprites that together form
// a continuous rectangle. Even though it has same properties as sprite, every
// operation is 9x as expensive, now little showcase of how NPS works:
//
//	+-+-+-+    	 	+-+------+-+
//	| | | |			| |      | |
//  +-+-+-+			+-+------+-+
//	| | | |			| |      | |
//  +-+-+-+			+-+------+-+
//	| | | |			| |      | |
//  +-+-+-+	(7, 7)	+-+------+-+ (12, 7)
//
// Patch is not window thread dependant
type Patch struct {
	s       [NinePatchSide][NinePatchSide]Sprite
	Padding mat.AABB
}

// NPatch creates ready-for-use Patch
func NPatch(frame, padding mat.AABB) Patch {
	v, h := PadMap(frame, padding)
	n := Patch{Padding: padding}

	for y := 0; y < NinePatchSide; y++ {
		for x := 0; x < NinePatchSide; x++ {
			vert := mat.A(v[x], h[y], v[x+1], h[y+1]).Vertices()
			for z := 0; z < 4; z++ {
				n.s[y][x].data[z] = Vertex{
					Tex:       vert[z],
					Intensity: 1,
					Color:     mat.Alpha(1),
				}
			}
		}
	}

	n.Resize(frame.W(), frame.H())

	return n
}

// SetColor sets modulation with witch Patch gets drawn
func (n *Patch) SetColor(value mat.RGBA) {
	for y := 0; y < NinePatchSide; y++ {
		for x := 0; x < NinePatchSide; x++ {
			n.s[y][x].SetColor(value)
		}
	}
}

// SetIntensity sets intensity of all inner sprites, so it has same effect
// as sprite intensity
func (n *Patch) SetIntensity(value float64) {
	for y := 0; y < NinePatchSide; y++ {
		for x := 0; x < NinePatchSide; x++ {
			n.s[y][x].SetIntensity(value)
		}
	}
}

// Resize resizes NPS to given width and height, corners will stay as same scale while
// other 5 parts scale up accordingly to create continuos Rectangle. This is mainly usefull
// for ui panels and flexible frames.
func (n *Patch) Resize(width, height float64) {
	/*
		function first takes mapping of new frame
		then i creates AABBs in loop for each sprite
		mapping is centered so that sprite can be drawn
		centered
	*/
	width *= .5
	height *= .5

	n.SetDist(mat.A(-width, -height, width, height))
}

// SetDist makes sprite drawn into dst area just by fetching it
func (n *Patch) SetDist(dst mat.AABB) {
	/*
		function first takes mapping of new frame
		then i creates AABBs in loop for each sprite
		mapping is centered so that sprite can be drawn
		centered
	*/
	v, h := PadMap(dst, n.Padding)
	for y := 0; y < NinePatchSide; y++ {
		for x := 0; x < NinePatchSide; x++ {
			n.s[y][x].tex = mat.A(v[x], h[y], v[x+1], h[y+1]).Vertices()
		}
	}
}

// Update transforms NPS vertices with matrix and sets mask
func (n *Patch) Update(mat mat.Mat, mask mat.RGBA) {
	for y := 0; y < NinePatchSide; y++ {
		for x := 0; x < NinePatchSide; x++ {
			for z, v := range n.s[y][x].tex {
				n.s[y][x].data[z].Pos = mat.Project(v)
				n.s[y][x].data[z].Color = mask
			}
		}
	}
}

// Fetch implements Fetcher interface
func (n *Patch) Fetch(t Target) {
	for y := 0; y < NinePatchSide; y++ {
		for x := 0; x < NinePatchSide; x++ {
			t.Accept(n.s[y][x].data[:], SpriteIndices)
		}
	}
}

// Draw implements Drawer interface
func (n *Patch) Draw(t Target, mat mat.Mat, mask mat.RGBA) {
	n.Update(mat, mask)
	n.Fetch(t)
}

// Size returns sprite size when its drawn with mat.IM transforation
func (n *Patch) Size() mat.Vec {
	return n.s[0][0].tex[0].To(n.s[2][2].tex[2])
}

// PadMap creates padding break points with help of witch we can determinate
// vertices of Patch, think of v as four vertical lines and h as 4 horizontal
// lines, if you draw them on paper you will get 9 rectangles
func PadMap(frame, padding mat.AABB) (v, h [4]float64) {
	v = [4]float64{frame.Min.X, frame.Min.X + padding.Min.X, frame.Max.X - padding.Max.X, frame.Max.X}
	h = [4]float64{frame.Min.Y, frame.Min.Y + padding.Min.Y, frame.Max.Y - padding.Max.Y, frame.Max.Y}
	return
}

// Sprite is most efficient way of drawing textures to Batch (if you find faster way i welcome your pr)
// sprite does not allocate any memory all data is on stack, its designed to be easily copied by value.
//
// Sprite is not window thread dependant
type Sprite struct {
	tex  [4]mat.Vec
	data [4]Vertex
}

// NSprite creates new sprite out of frame. Frame should be the region where the texture,
// you want to draw, is located on spritesheet
//
//	ggl.NSprite(yourTexture.Frame()) // draws whole texture
//
func NSprite(frame mat.AABB) Sprite {
	s := Sprite{}
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

// SetColor sets the color of sprite
func (s *Sprite) SetColor(value mat.RGBA) {
	for i := range s.data {
		s.data[i].Color = value
	}
}

// SetIntensity sets the intensity of sprite image, if you set it to 0 the rectangle in a color of sprite mask will
// be drawn, if you set it to 1 (which is default) it will draw texture as it is or black area if batch does not
// have texture.
func (s *Sprite) SetIntensity(value float64) {
	for i := range s.data {
		s.data[i].Intensity = value
	}
}

// Clear makes sprite invisible when its drawn, this is to impelemt paceholder data
func (s *Sprite) Clear() {
	for i := range s.data {
		s.data[i].Pos = mat.Vec{}
	}
}

// Set sets sprites source texture region and destination rectangle, this is mainly used when drawing text
func (s *Sprite) Set(dst, src mat.AABB) {
	tex, pos := src.Vertices(), dst.Vertices()
	for i := range s.data {
		s.data[i].Tex = tex[i]
		s.data[i].Pos = pos[i]
	}
}

// SetDist sets destination area where sprite will be drawn
func (s *Sprite) SetDist(dst mat.AABB) {
	pos := dst.Vertices()
	for i := range s.data {
		s.data[i].Pos = pos[i]
	}
}

// Draw draws sprite to Batch projected by given matrix, and colored by mask
func (s *Sprite) Draw(b Target, mat mat.Mat, mask mat.RGBA) {
	s.Update(mat, mask)
	b.Accept(s.data[:], SpriteIndices)
}

// Fetch draws sprite as it is, to change draw result call Update
func (s *Sprite) Fetch(b Target) {
	b.Accept(s.data[:], SpriteIndices)
}

// Update only updates sprite data but does not draw it
func (s *Sprite) Update(mat mat.Mat, mask mat.RGBA) {
	for i := range s.data {
		s.data[i].Pos = mat.Project(s.tex[i])
		s.data[i].Color = mask
	}
}

// Size returns sprite size when its drawn with mat.IM transforation
func (s *Sprite) Size() mat.Vec {
	return s.data[0].Tex.To(s.data[2].Tex)
}
