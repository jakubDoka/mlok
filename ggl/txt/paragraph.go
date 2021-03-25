package txt

import (
	"fmt"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/mat"

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
	// it does not overflows Width, though it only breaks on spaces, if width
	// is 0 it will never wrap
	Width float64
	// Fields are only relevent if Custom Lineheight is true, othervise they will
	// get overwritten
	LineHeight, Ascent, Descent float64
	// if this is true no effects are displayed
	NoEffects bool
	// set this to true if you want it custom
	CustomLineheight bool

	// mask is what all triangles solors will be multiplied by
	Mask mat.RGBA

	// Font determinate base
	Font string
	// Content is only important to markdown, when you are drawing directly with
	// drawer, Content is not used
	Content str.String

	// Align defines how the text is aligned
	Align

	data, selection ggl.Data

	progress float64
	lines    []line
	dots     []mat.Vec
	dot      mat.Vec
	bounds   mat.AABB

	Compiled str.String

	changing, instant Effs
	chunks            FEffs
}

// Clear is only usefull when drawing to paragraph directly with drawer
// it clears triangles
func (p *Paragraph) Clear() {
	p.data.Clear()
	p.dots = p.dots[:0]
	p.dot = mat.ZV
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
func (p *Paragraph) Update(delta float64) {
	p.Data.Clear()
	p.Data.Indices = p.data.Indices
	mat := p.Mat()
	for _, t := range p.data.Vertexes {
		t.Pos = mat.Project(t.Pos)
		t.Color = t.Color.Mul(p.Mask)

		p.Data.Vertexes = append(p.Data.Vertexes, t)
	}

	p.progress += delta

	for _, e := range p.changing {
		e.Apply(p.Data.Vertexes, p.progress)
	}
}

// Changes returns whether p.Update will change triangles
func (p *Paragraph) Changes() bool {
	return len(p.changing) != 0
}

// Draw draws its triangles to given target
func (p *Paragraph) Draw(t ggl.Target) {
	t.Accept(p.Data.Vertexes, p.Data.Indices)
}

// CursorFor returns index of a dot on witch the cursor is nad insertcion index
// that is where you should insert the text to make it look like cursor is writing.
// Then there is line local position and index of a line
func (p *Paragraph) CursorFor(mouse mat.Vec) (global, local, line int) {
	mouse = p.Mat().Unproject(mouse)

	y := 0
	for ; y < len(p.lines); y++ {
		if p.lines[y].level < mouse.Y {
			break
		}
	}

	if y >= len(p.lines) { // its under
		y--
	}

	r := p.lines[y]
	x := r.start
	for ; x < r.end; x++ {
		if p.dots[x].X > mouse.X {
			break
		}
	}

	if x > 0 { // its before the charater
		x--
	}

	return x, x - r.start, y
}

// Dot returns dot at given index
func (p *Paragraph) Dot(i int) mat.Vec {
	return p.Mat().Project(p.dots[i])
}

// ProjectLine projects line and local index intro global index
// complexity is O(1)
func (p *Paragraph) ProjectLine(i, line int) int {
	l := p.lines[line]
	return mat.Mini(l.start+i, l.end-1)
}

// UnprojectLine does reverse of Project line complexity is O(n)
// where n is len(p.lines)
func (p *Paragraph) UnprojectLine(idx int) (local, line int) {
	for i, l := range p.lines {
		if idx < l.end {
			return idx - l.start, i
		}
	}

	return -1, -1
}

// Lines returns count of lines in paragraph
func (p *Paragraph) Lines() int {
	return len(p.lines)
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

// line stores start end and level of a line
type line struct {
	level      float64
	start, end int
}

// Allign determinate text align
type Align float64

func (a Align) String() string {
	for k, v := range Aligns {
		if v == a {
			return k
		}
	}

	return fmt.Sprint(float64(a))
}

// Align constants
const (
	Left   Align = 0
	Middle Align = .5
	Right  Align = 1
)

var Aligns = map[string]Align{
	"left":   Left,
	"middle": Middle,
	"right":  Right,
}
