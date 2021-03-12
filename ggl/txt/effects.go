package txt

import (
	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/mat"
)

// Effect is interface for eny effect applied on text
// TriangleData is slice that should be modified, delta is to change state of effect
// in case this is permanent effect
type Effect interface {
	Apply(try ggl.Vertexes, delta float64)
	Kind() int8
	Start() int
	Close(endIdx int)
	Copy(startIdx int) Effect
}

// Effect types
const (
	Instant int8 = iota
	Changing
	TextType
)

// ColorEffect is one shot effect that changes color of a text
type ColorEffect struct {
	EffectBase
	Color mat.RGBA
}

// NColorEffect constructs new color effect as long as s is valid hex color
func NColorEffect(hex string, start int) (*ColorEffect, error) {
	c, err := mat.HexToRGBA(hex)
	if err != nil {
		return nil, err
	}

	return &ColorEffect{
		EffectBase{start * ggl.SpriteVertexSize, 0},
		c,
	}, nil
}

// Kind implements Effect
func (e *ColorEffect) Kind() int8 {
	return Instant
}

// Apply implements Effect
func (e *ColorEffect) Apply(try ggl.Vertexes, _ float64) {
	try = try[e.start:e.End]
	for i := range try {
		try[i].Color = e.Color
	}
}

// Close implements Effect
func (e *ColorEffect) Close(endIdx int) {
	e.End = endIdx * ggl.SpriteVertexSize
}

// Copy implements Effect
func (e *ColorEffect) Copy(start int) Effect {
	ev := *e
	ev.start = start * ggl.SpriteVertexSize
	return &ev
}

// FontEffect stores what font should be used for given slice of text
type FontEffect struct {
	EffectBase
	Font string
}

// NFontEffect is here for consistency
func NFontEffect(font string, start, end int) *FontEffect {
	return &FontEffect{
		EffectBase{start, end},
		font,
	}
}

// Apply implements Effect
func (f *FontEffect) Apply(_ ggl.Vertexes, _ float64) {
	panic("unimplemented")
}

// Kind implements Effect
func (f *FontEffect) Kind() int8 {
	return TextType
}

// Copy implements Effect
func (f *FontEffect) Copy(start int) Effect {
	panic("unimplemented")
}

// EffectBase hold properties that every effect has
type EffectBase struct {
	start, End int
}

// Close implements Effect
func (e *EffectBase) Close(endIdx int) {
	e.End = endIdx
}

// Redundant ...
func (e *EffectBase) Redundant() bool {
	return e.End == e.start
}

// Start implements Effect
func (e *EffectBase) Start() int {
	return e.start
}
