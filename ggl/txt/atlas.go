package txt

import (
	"gobatch/ggl"
	"gobatch/mat"
	"image"
	"image/draw"
	"math"
	"sort"
	"unicode"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Atlas7x13 is an Atlas using basicfont.Face7x13 with the ASCII rune set
var Atlas7x13 *Atlas

// ASCII is a set of all ASCII runes. These runes are codepoints from 32 to 127 inclusive.
var ASCII []rune

func init() {
	ASCII = make([]rune, unicode.MaxASCII-32)
	for i := range ASCII {
		ASCII[i] = rune(32 + i)
	}
	Atlas7x13 = NewAtlas("default", basicfont.Face7x13, 0, ASCII)
}

// Glyph describes one glyph in an Atlas.
type Glyph struct {
	Dot     mat.Vec
	Frame   mat.AABB
	Advance float64
}

// Atlas is a set of pre-drawn glyphs of a fixed set of runes. This allows for efficient text drawing.
type Atlas struct {
	face       font.Face
	Pic        *image.NRGBA
	mapping    map[rune]Glyph
	ascent     float64
	descent    float64
	lineHeight float64
	spacing    float64
	Name       string
}

// NewAtlas creates a new Atlas containing glyphs of the union of the given sets of runes (plus
// unicode.ReplacementChar) from the given font face. Spacing is space in pixels that is added around
// each glyph. This is very usefull when applying shaders to your text, mind that bigger sprites burdens
// gpu more.
//
// Creating an Atlas is rather expensive, do not create a new Atlas each frame.
//
// Do not destroy or close the font.Face after creating the Atlas. Atlas still uses it.
func NewAtlas(name string, face font.Face, spacing int, runeSets ...[]rune) *Atlas {
	seen := make(map[rune]bool)
	runes := []rune{unicode.ReplacementChar}
	for _, set := range runeSets {
		for _, r := range set {
			if !seen[r] {
				runes = append(runes, r)
				seen[r] = true
			}
		}
	}

	fixedMapping, fixedBounds := makeSquareMapping(face, runes, fixed.I(2), spacing)

	atlasImg := image.NewNRGBA(image.Rect(
		fixedBounds.Min.X.Floor(),
		fixedBounds.Min.Y.Floor(),
		fixedBounds.Max.X.Ceil(),
		fixedBounds.Max.Y.Ceil(),
	))

	for i := range fixedMapping {
		fg := &fixedMapping[i]
		dr, mask, maskp, _, _ := face.Glyph(fg.dot, fg.r)
		draw.Draw(atlasImg, dr, mask, maskp, draw.Src)
	}

	ggl.FlipNRGBA(atlasImg)

	bounds := mat.A(
		i2f(fixedBounds.Min.X),
		i2f(fixedBounds.Min.Y),
		i2f(fixedBounds.Max.X),
		i2f(fixedBounds.Max.Y),
	)

	mapping := make(map[rune]Glyph)
	for i := range fixedMapping {
		fg := &fixedMapping[i]
		mapping[fg.r] = Glyph{
			Dot: mat.V(
				i2f(fg.dot.X),
				bounds.Max.Y-i2f(fg.dot.Y),
			),
			Frame: mat.A(
				i2f(fg.frame.Min.X),
				bounds.Max.Y-i2f(fg.frame.Max.Y),
				i2f(fg.frame.Max.X),
				bounds.Max.Y-i2f(fg.frame.Min.Y),
			),
			Advance: i2f(fg.advance),
		}
	}

	return &Atlas{
		face:       face,
		Pic:        atlasImg,
		mapping:    mapping,
		Name:       name,
		spacing:    float64(spacing),
		ascent:     i2f(face.Metrics().Ascent),
		descent:    i2f(face.Metrics().Descent),
		lineHeight: i2f(face.Metrics().Height),
	}
}

// Contains reports wheter r in contained within the Atlas.
func (a *Atlas) Contains(r rune) bool {
	_, ok := a.mapping[r]
	return ok
}

// Glyph returns the description of r within the Atlas.
func (a *Atlas) Glyph(r rune) Glyph {
	return a.mapping[r]
}

// Kern returns the kerning distance between runes r0 and r1. Positive distance means that the
// glyphs should be further apart.
func (a *Atlas) Kern(r0, r1 rune) float64 {
	return i2f(a.face.Kern(r0, r1))
}

// Ascent returns the distance from the top of the line to the baseline.
func (a *Atlas) Ascent() float64 {
	return a.ascent
}

// Descent returns the distance from the baseline to the bottom of the line.
func (a *Atlas) Descent() float64 {
	return a.descent
}

// LineHeight returns the recommended vertical distance between two lines of text.
func (a *Atlas) LineHeight() float64 {
	return a.lineHeight
}

// DrawRune returns parameters necessary for drawing a rune glyph.
//
// Rect is a rectangle where the glyph should be positioned. Frame is the glyph frame inside the
// Atlas's Picture. NewDot is the new position of the dot.
func (a *Atlas) DrawRune(prevR, r rune, dot mat.Vec) (rect, frame, bounds mat.AABB, newDot mat.Vec) {
	if !a.Contains(r) {
		r = unicode.ReplacementChar
	}
	if !a.Contains(unicode.ReplacementChar) {
		newDot = dot
		return
	}
	if !a.Contains(prevR) {
		prevR = unicode.ReplacementChar
	}

	var kern float64
	if prevR >= 0 {
		kern = a.Kern(prevR, r)
	}

	dot.X += kern

	glyph := a.Glyph(r)

	rect = glyph.Frame.Moved(dot.Sub(glyph.Dot))

	bounds = mat.A(dot.X-kern, dot.Y-a.Descent(), dot.X+glyph.Advance, dot.Y+a.Ascent())

	dot.X += glyph.Advance

	return rect, glyph.Frame, bounds, dot
}

type fixedGlyph struct {
	r       rune
	dot     fixed.Point26_6
	frame   fixed.Rectangle26_6
	advance fixed.Int26_6
}

// makeSquareMapping finds an optimal glyph arrangement of the given runes, so that their common
// bounding box is as square as possible.
func makeSquareMapping(face font.Face, runes []rune, padding fixed.Int26_6, spacing int) ([]fixedGlyph, fixed.Rectangle26_6) {
	var bounds fixed.Rectangle26_6

	buff := make([]fixedGlyph, len(runes))
	bounds = makeMapping(face, runes, padding, math.MaxInt32, spacing, &buff) // find longest possible composition

	sort.Search(int(bounds.Max.X-bounds.Min.X), func(i int) bool {
		width := fixed.Int26_6(i)
		bounds = makeMapping(face, runes, padding, width, spacing, &buff)
		return bounds.Max.X-bounds.Min.X >= bounds.Max.Y-bounds.Min.Y
	})

	return buff, bounds
}

// makeMapping arranges glyphs of the given runes into rows in such a way, that no glyph is located
// fully to the right of the specified width. Specifically, it places glyphs in a row one by one and
// once it reaches the specified width, it starts a new row.
func makeMapping(face font.Face, runes []rune, padding, width fixed.Int26_6, spacing int, buffer *[]fixedGlyph) fixed.Rectangle26_6 {
	buff := (*buffer)[:0]
	bounds := fixed.Rectangle26_6{}
	additional := fixed.I(2 * spacing)

	dot := fixed.P(0, 0)
	dot.Y = face.Metrics().Ascent + fixed.I(spacing)

	for _, r := range runes {
		b, advance, ok := face.GlyphBounds(r)
		if !ok {
			continue
		}

		// this is important for drawing, artifacts arise otherwise
		frame := fixed.Rectangle26_6{
			Min: fixed.P(b.Min.X.Floor()-spacing, b.Min.Y.Floor()-spacing),
			Max: fixed.P(b.Max.X.Ceil()+spacing, b.Max.Y.Ceil()+spacing),
		}

		dot.X -= frame.Min.X
		frame = frame.Add(dot)

		buff = append(buff, fixedGlyph{
			r:       r,
			dot:     dot.Add(fixed.P(0, spacing)),
			frame:   frame,
			advance: advance,
		})
		bounds = bounds.Union(frame)

		dot.X = frame.Max.X

		// padding + align to integer
		dot.X += padding
		dot.X = fixed.I(dot.X.Ceil())

		// width exceeded, new row
		if frame.Max.X >= width {
			dot.X = 0
			dot.Y += face.Metrics().Ascent + face.Metrics().Descent + additional

			// padding + align to integer
			dot.Y += padding
			dot.Y = fixed.I(dot.Y.Ceil())
		}
	}

	if dot.X != 0 {
		bounds.Max.Y += fixed.I(spacing)
	}

	*buffer = buff
	return bounds
}

func i2f(i fixed.Int26_6) float64 {
	return float64(i) / (1 << 6)
}
