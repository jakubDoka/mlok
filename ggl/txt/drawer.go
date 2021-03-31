package txt

import (
	"unicode"
	"unicode/utf8"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/mat"
)

// Drawer draws text for ContentBox
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

// Draw draws a text to Target, string is decoded by utf8
func (d *Drawer) Draw(t Target, text string) {
	var (
		dot, bounds           = t.Metrics()
		raw                   = []byte(text)
		prev                  = rune(-1)
		rBounds, region, rect mat.AABB
	)

	control := func(r rune) bool {
		switch r {
		case '\n':
			dot.X = 0
			dot.Y += d.LineHeight()
		case '\r':
			dot.X = 0
		case '\t':
			dot.X += d.tab
		default:
			return false
		}
		return true
	}

	for r, size := utf8.DecodeRune(raw); r != utf8.RuneError; r, size = utf8.DecodeRune(raw) {
		raw = raw[size:]

		if control(r) {
			continue
		}

		rect, region, rBounds, dot = d.DrawRune(prev, r, dot)
		bounds = bounds.Union(rBounds)

		d.glyph.Set(rect, region.Moved(d.Region))
		d.glyph.Fetch(t)
	}

	t.SetMetrics(dot, bounds)
}

// Draw draws a slice of p.Compiled to p.data, text continused where last draw stopped
func (d *Drawer) drawParagraph(p *Paragraph, start, end int) {
	var (
		prev rune = -1
		// last stores data about last seen space
		last struct {
			present             bool
			idx, vertex, indice int
			dot                 mat.Vec
			bounds              mat.AABB
		}
		rect, frame, bounds mat.AABB
		control             bool
	)

	for i := start; i < end; i++ {
		r := p.Compiled[i]

		if r == ' ' {
			last.idx = i
			last.vertex = p.data.Vertexes.Len()
			last.indice = len(p.data.Indices)
			last.present = true
			last.dot = p.dot
			last.bounds = p.bounds
			control = false
		} else {
			control = d.paragraphControlRune(r, p)
		}

		if control {
			// we don't want our effects to get offset so we are appending empty spaces anyway. As long as
			// its a glyph it should hold a place
			d.glyph.Clear()
		} else {
			rect, frame, bounds, p.dot = d.DrawRune(prev, r, p.dot)
			// text is overflowing bounds so erase last word and write it on new line
			// but only if there is a space to break it on
			if p.Width != 0 && last.present && p.dot.X > p.Width {
				p.dots = p.dots[:last.idx+1]
				p.dot = last.dot
				p.bounds = last.bounds

				d.paragraphControlRune('\n', p)
				d.glyph.Clear()

				// truncating data to previous state
				p.data.Vertexes = p.data.Vertexes[:last.vertex]
				p.data.Indices = p.data.Indices[:last.indice]

				// space is now replaced with newline, reusing it would create endless loop
				last.present = false
				r = '\n'

				i = last.idx
				p.Compiled[i] = r
			} else {
				p.dots = append(p.dots, p.dot)
				d.glyph.Set(rect, frame.Moved(d.Region))
				p.bounds = p.bounds.Union(bounds)
			}
		}

		prev = r
		d.glyph.Fetch(&p.data)
	}

	p.bounds = p.bounds.Union(mat.Square(p.dot, 0))
}

// ControlRune changes dot accordingly if inputted rune is control rune, also returns whether
// change happened, it also appends a new dot to slice
func (d *Drawer) paragraphControlRune(r rune, p *Paragraph) bool {
	switch r {
	case '\n':
		p.lines[len(p.lines)-1].end = len(p.dots)
		p.dot.X = 0
		p.dot.Y -= p.LineHeight
		p.lines = append(p.lines, line{p.dot.Y, len(p.dots), -1})
	case '\r':
		p.dot.X = 0
	case '\t':
		p.dot.X += d.tab
	default:
		return false
	}

	p.dots = append(p.dots, p.dot)
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

// Text is a builtin Drawer target, it can act as sprite
type Text struct {
	ggl.Data
	projected ggl.Data
	Dot       mat.Vec
	Bounds    mat.AABB
}

// Draws the text centered on matrix center
func (t *Text) DrawCentered(tg ggl.Target, mat mat.Mat, color mat.RGBA) {
	tr := mat.Move(mat.C.Inv())
	t.Draw(tg, mat.Move(tr.Project(t.Bounds.Center().Inv())), color)
}

// Draw implements ggl.Drawer interface
func (t *Text) Draw(tg ggl.Target, mat mat.Mat, color mat.RGBA) {
	t.Update(mat, color)
	t.Fetch(tg)
}

// Update updates the text color and transforation
func (t *Text) Update(mat mat.Mat, color mat.RGBA) {
	t.projected.Clear()
	t.projected.Indices = t.Indices
	t.projected.Vertexes = append(t.projected.Vertexes, t.Vertexes...)
	for i := range t.projected.Vertexes {
		v := &t.projected.Vertexes[i]
		v.Color = color
		v.Pos = mat.Project(v.Pos)
	}
}

// Fetch implements ggl.Fetcher interface
func (t *Text) Fetch(tg ggl.Target) {
	t.projected.Fetch(tg)
}

// Metrics implements Target interface
func (t *Text) Metrics() (mat.Vec, mat.AABB) {
	return t.Dot, t.Bounds
}

// SetMetrics implements Target interface
func (t *Text) SetMetrics(dot mat.Vec, bounds mat.AABB) {
	t.Dot, t.Bounds = dot, bounds
}

// Target is something that Drawer can draw to
type Target interface {
	ggl.Target
	// Metrics returns current dot and bounds of Target
	Metrics() (mat.Vec, mat.AABB)
	// SetMetrics sets the metrics of Target to new ones
	SetMetrics(mat.Vec, mat.AABB)
}
