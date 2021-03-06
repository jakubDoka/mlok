package txt

import (
	"gobatch/ggl"
	"gobatch/mat"
	"math"
	"unicode"

	"github.com/jakubDoka/gogen/str"
)

// Drawer draws text for TextBox
type Drawer struct {
	*Atlas
	Region mat.Vec
	glyph  ggl.Sprite
	tab    float64
}

// NDrawer is drawer constructor
func NDrawer(atlas *Atlas) *Drawer {
	t := &Drawer{
		Atlas: atlas,
	}

	t.glyph.Update(mat.IM, mat.Alpha(1))
	t.glyph.SetIntensity(1)

	t.tab = t.Glyph(' ').Advance * 4

	return t
}

// Write writes string to paragraph
func (d *Drawer) Write(p *Paragraph, text string) {
	s := str.NString(text)
	start := len(p.raw)
	p.raw = append(p.Text, s...)
	d.Draw(p, start, len(p.raw))
}

// Draw draws a slice of p.raw to p.data, text continused where last draw stopped
func (d *Drawer) Draw(p *Paragraph, start, end int) {
	var (
		prevRune rune = -1
		// last stores data about last seen space
		last struct {
			present             bool
			idx, vertex, indice int
		}
		rect, frame, bounds mat.AABB
		control             bool
	)

	for i := start; i < end; i++ {
		r := p.raw[i]

		if r == ' ' {
			last.idx = i
			last.vertex = p.data.Vertexes.Len()
			last.indice = len(p.data.Indices)
			last.present = true
		} else {
			control = d.ControlRune(r, p)
		}

		if control {
			// we don't want our effects to get offset so we are appending empty spaces anyway. As long as
			// its a glyph it should hold a place
			d.glyph.Clear()
		} else {
			rect, frame, bounds, p.dot = d.DrawRune(prevRune, r, p.dot)
			// text is overflowing bounds so erase last word and write it on new line
			// but only if there is a space to break it on
			if p.Width != 0 && last.present && p.dot.X > p.Width {
				p.dots = p.dots[:last.idx]
				d.ControlRune('\n', p)
				d.glyph.Clear()

				// truncating data to previous state
				p.data.Vertexes = p.data.Vertexes[:last.vertex]
				p.data.Indices = p.data.Indices[:last.indice]

				// space is now replaced with newline, reusing it would create endless loop
				last.present = false
				prevRune = '\n'

				i = last.idx
				p.raw[i] = '\n'
			} else {
				prevRune = r
				p.dots = append(p.dots, p.dot)
				d.glyph.Set(frame, rect)
				p.bounds = p.bounds.Union(bounds)
			}
		}

		d.glyph.Fetch(&p.data)
	}
}

// ControlRune changes dot accordingly if inputted rune is control rune, also returns whether
// change happened, it also appends a new dot to slice
func (d *Drawer) ControlRune(r rune, p *Paragraph) bool {
	var nDot mat.Vec

	switch r {
	case '\n':
		// little hack to make cursor snap to end of line
		nDot = mat.V(math.MaxFloat64, p.dot.Y)
		p.dot.X = 0
		p.dot.Y -= p.LineHeight
	case '\r':
		p.dot.X = 0
		nDot = p.dot
	case '\t':
		p.dot.X += d.tab
		nDot = p.dot
	default:
		return false
	}

	p.dots = append(p.dots, nDot)
	return true
}

// Advance calculates glyph advance for this text
func (d *Drawer) Advance(prev, r rune) (l float64) {
	if !d.Contains(r) {
		r = unicode.ReplacementChar
	}
	if !d.Contains(unicode.ReplacementChar) {
		return
	}
	if !d.Contains(prev) {
		prev = unicode.ReplacementChar
	}

	if prev >= 0 {
		l += d.Kern(prev, r)
	}

	return l + d.Glyph(r).Advance
}
