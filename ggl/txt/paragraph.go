package txt

import (
	"gobatch/ggl"
	"gobatch/mat"

	"github.com/jakubDoka/gogen/str"
)

/*imp(
	gogen/templates
)*/

/*gen(
	templates.Vec<Effect, Effs>
	templates.Vec<*FontEffect, FEffs>
)*/

// Paragraph stores triangles that can form text, there is a lot of passive configurations
// and to apply changes you have to pass Paragraph to Markdown.Parse method thought you can also
// use Drawer directly to draw text though then you have to initialize lineheight
type Paragraph struct {
	// saves the paragraph transformation, you have to call Update to apply changes
	mat.Tran

	// this data can be freely modified red, serialized, after Update its overwritten
	// anyway
	Data ggl.Data

	// determines how text should wrap, Drawer will tri to display text so
	// it does not overflows Width, though it only breaks on spaces
	Width float64
	// this field is only used if it isn't negative and it alters LineHight (spooks)
	LineHeight float64
	// if this is true no effects are displayed
	NoEffects bool
	// set this to true if you want it custom
	CustomLineheight bool

	// mask is what all triangles solors will be multiplied by
	Mask mat.RGBA

	// Font determinate base
	Font string
	// Text is only important to markdown, when you are drawing directly with
	// drawer, Text is not used
	Text str.String

	data ggl.Data

	progress float64
	dots     []mat.Vec
	dot      mat.Vec
	bounds   mat.AABB

	initted bool

	raw str.String

	changing, instant Effs
	chunks            FEffs
}

// Init initializes paragraph if key fields were omitted, its only for user experience
func (p *Paragraph) Init() {
	p.initted = true
	if p.Tran.Scl == mat.Origin {
		p.Tran.Scl = mat.Scale
	}

	if p.Mask == mat.Transparent {
		p.Mask = mat.White
	}
}

// Clear is only usefull when drawing to paragraph directly with drawer
// it clears triangles
func (p *Paragraph) Clear() {
	p.data.Clear()
	p.dots = p.dots[:0]
	p.dot = mat.Origin
}

// AddEff appends chunk of test to paragraph
func (p *Paragraph) AddEff(e Effect) {
	switch e.Kind() {
	case Instant:
		p.instant = append(p.instant, e)
	case Changing:
		p.changing = append(p.changing, e)
	case TextType:
		p.chunks = append(p.chunks, e.(*FontEffect))
	default:
		panic("invalid event kind")
	}
}

// Sort is part of markdown building procedure, it sorts all effects so they can be applied properly
func (p *Paragraph) Sort() {
	s := func(a, b Effect) bool {
		return a.Start() < b.Start()
	}

	p.changing.Sort(s)
	p.instant.Sort(s)

	p.chunks.Sort(func(a, b *FontEffect) bool {
		return a.start < b.start
	})
}

// Update has to be called after changes or they will not be visible, it returns false if
// there is nothing to draw
func (p *Paragraph) Update(delta float64, start, end float64) bool {
	if !p.initted {
		p.Init()
	}

	p.Data.Clear()
	p.Data.Indices = p.data.Indices

	clamp := start != end
	visible := !clamp
	mat := p.Mat()
	for _, t := range p.data.Vertexes {
		// just clamping them isn't really memory efficient but logic that would determinate whether
		// glyph is not visible and remove it would be slower and more bug prone
		if clamp {
			if t.Pos.Y > start {
				t.Pos.Y = start
			} else if t.Pos.Y < end {
				t.Pos.Y = end
			} else {
				visible = true
			}
		}

		t.Pos = mat.Project(t.Pos)
		t.Color = t.Color.Mul(p.Mask)

		p.Data.Vertexes = append(p.Data.Vertexes, t)
	}

	p.progress += delta

	for _, e := range p.changing {
		e.Apply(p.Data.Vertexes, p.progress)
	}

	return visible || len(p.Data.Vertexes) == 0
}

// Draw draws its triangles to given target
func (p *Paragraph) Draw(t ggl.Target) {
	t.Accept(p.Data.Vertexes, p.Data.Indices)
}

// CursorFor returns snapped position of cursor and its index, index is between 0 and
// len(displayedGlyphs)+1 if you are using Effects and try to use cursor index to insert into
// Text your attempt will fail as with effects present len(Text) > len(displayedGlyphs)
// Note that you first have to call Markdown.Parse on paragraph that creates mapping for
// finding cursor, othervise you will end up with zero values or invalid values.
func (p *Paragraph) CursorFor(mouse mat.Vec) (pos mat.Vec, idx int) {
	mouse = p.Mat().Unproject(mouse)

	for i, pos := range p.dots {
		if pos.X >= mouse.X && pos.Y <= mouse.Y {
			return pos, i
		}
	}

	return
}

// Bounds returns bounding rectangle of untransformed text, bounds are valid only after
// Markdown.Parse call
func (p *Paragraph) Bounds() mat.AABB {
	return p.bounds
}

// SetCenter moves paragraph so it is centered at given position
func (p *Paragraph) SetCenter(pos mat.Vec) {
	v := p.bounds.ToVec().Mul(mat.V(.5, -.5))
	p.Pos = pos.Sub(v)
}
