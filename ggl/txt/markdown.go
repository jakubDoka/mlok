package txt

import (
	"io/ioutil"
	"math"
	"os"

	"github.com/jakubDoka/mlok/mat"

	"github.com/golang/freetype/truetype"
	"github.com/jakubDoka/gogen/str"
	"golang.org/x/image/font"
)

// key markdown sintax is stored in variable  so it can be customized
var (
	MarkdownIdent = '!'
	ColorIdent    = '#'
	BlockStart    = '['
	BlockEnd      = ']'
	NullIdent     = string([]rune{1})
	DefaultFont   = "default"
	noConstructor = "markdown is missing default font, use use constructor that adds it"
)

// Markdown handles text markdown parsing, markdown sintax is as follows:
//
// 	#FF0000[hello] 	// hello will appear ugly red
// 	#FF000099[hello] // hello will appear ugly red and slightly transparent
//
// Now writing effects like this is maybe flexible but not always convenient to user
// thats why you can define your own effects by adding effects to Markdown.Effects. Name
// of the effect is a key in map thus:
//
// 	m.Effects["red"] = ColorEffect{Color: mat.RGB(1, 0, 0)}
//
// will do the trick. user then can use this effect like:
//
//	!red[hello] // witch does the same
//
// now you can make any kind of triangle mutator you like, it even allows dinamic effects
// so something like fading text or exploding text is possible. You can also use multiple
// fonts in one paragraph by adding more atlases to your markdown. Pattern is same as for
// adding custom effects:
//
//	m.Effects["italic"] = NDrawer(italicFontAtlas)
//
// User then can use font like:
//
// 	!italic[hello] // hello will be italic, syntax will not appear
//
// Last feature of markdown are shortcuts. After you have added all effects you wanted you can call
// GenerateShortcuts method. This will map all effect names to its starting rune thus user can be very
// lazy:
//
//	!i[hello] // hello will be italic with low effort
//
// Markdown is only compatible with paragraph, mind that parsing markdown is slow and grows linearly with text
// length, O(n), if course if you want effects to even display you have to set DisplayEffects to true in paragraph
type Markdown struct {
	Shortcuts map[rune]string
	Fonts     map[string]*Drawer
	Effects   map[string]Effect

	buff, stack FEffs
	stack2      Effs
}

// NMarkdown initializes inner maps and adds default drawer
func NMarkdown() *Markdown {
	m := &Markdown{
		Shortcuts: map[rune]string{},
		Effects: map[string]Effect{
			"red":   &ColorEffect{Color: mat.Red},
			"green": &ColorEffect{Color: mat.Green},
			"blue":  &ColorEffect{Color: mat.Blue},
		},
		Fonts: map[string]*Drawer{DefaultFont: NDrawer(Atlas7x13)},
	}

	m.GenerateShortcuts()

	return m
}

// GenerateShortcuts creates shortcuts for all effects, if names overlap random one is bind
func (m *Markdown) GenerateShortcuts() {
	for k := range m.Fonts {
		if k == "" {
			continue
		}
		m.Shortcuts[rune(k[0])] = k
	}

	for k := range m.Effects {
		if k == "" {
			continue
		}
		m.Shortcuts[rune(k[0])] = k
	}
}

// Parse turns markdown stored in p.Content into final text with effects, for markdown syntax see
// struct documentation
func (m *Markdown) Parse(p *Paragraph) {
	if _, ok := m.Fonts[p.Font]; !ok {
		p.Font = DefaultFont
		if _, ok := m.Fonts[p.Font]; !ok {
			panic(noConstructor)
		}
	}

	p.Compiled = append(p.Compiled[:0], p.Content...)

	p.changing.Clear()
	p.instant.Clear()
	p.chunks.Clear()

	if !p.NoEffects {
		m.CollectEffects(p)
		p.Sort()
	}

	m.ResolveChunks(p)
	m.MakeTriangles(p)
	return
}

// CollectEffects removes all valid effect syntax and stores parsed effects in paragraph
func (m *Markdown) CollectEffects(p *Paragraph) {

	var (
		mv, i int
		ident string
		ok    bool
	)

	m.stack2 = m.stack2[:0]

	push := func() {
		ef := m.stack2.Pop()
		ef.Close(i)
		p.AddEff(ef)
	}
o:
	for ; i < len(p.Compiled); i += mv {
		b := p.Compiled[i]
		mv = 1

		switch b {
		case BlockEnd: // fond text that should be skipped
			if len(m.stack2) != 0 {
				p.Compiled.Remove(i)
				if i < len(p.Compiled) && p.Compiled[i] == ']' {
					continue
				}
				mv = 0
				push()
				continue
			}
		case ColorIdent, MarkdownIdent:
			// ingoreing
		default:
			continue
		}

		if i+2 >= len(p.Compiled) || !str.IsIdent(byte(p.Compiled[i+1])) { //shortcut can't fit there or space right after explanation mark
			continue
		}

		if p.Compiled[i+2] == BlockStart { // in case of shortcut - shortcut is always just one rune
			ident, ok = m.Shortcuts[p.Compiled[i+1]]
			if !ok { // invalid shortcut so ignore it
				continue
			}
			p.Compiled.RemoveSlice(i, i+3)
			mv = 0
		} else { // find out full identifier
			k := i + 1
			for {
				if k >= len(p.Compiled) {
					continue o //out of bounds and we haven't even found non ident byte, ignoring
				}

				if !str.IsIdent(byte(p.Compiled[k])) {
					if p.Compiled[k] != BlockStart {
						continue o //ident should end with BlockStart, ignoring
					}
					break
				}
				k++
			}

			ident = string(p.Compiled[i+1 : k]) // i+1 because we are not including ident

			if b == ColorIdent { // this can also be color ident so handle it
				ce, err := NColorEffect(ident, i)
				if err != nil {
					continue
				}
				m.stack2 = append(m.stack2, ce)
				ident = NullIdent // we don't want to handle ident twice
			}

			p.Compiled.RemoveSlice(i, k+1)
			mv = 0
		}

		if ident == NullIdent {
			continue
		}

		if _, ok := m.Fonts[ident]; ok {
			m.stack2 = append(m.stack2, NFontEffect(ident, i, 0))
		} else if val, ok := m.Effects[ident]; ok {
			m.stack2 = append(m.stack2, val.Copy(i))
		}

	}

	for len(m.stack2) != 0 { // close all reminding effects
		push()
	}
}

// MakeTriangles creates triangles, these are not drawn to screen as only
// instant effects are applied on them and you have to call Update on paragraph
// for anything to show up
func (m *Markdown) MakeTriangles(p *Paragraph) {
	p.data.Clear()

	p.dot = mat.V(0, -p.Ascent)
	p.bounds = mat.A(0, -p.LineHeight, 0, 0)

	p.dots = append(p.dots[:0], p.dot)
	p.lines = append(p.lines[:0], line{p.dot.Y, 0, -1})

	for _, c := range p.chunks {
		m.Fonts[c.Font].drawParagraph(p, c.start, c.End)
	}

	end := &p.lines[len(p.lines)-1]
	if end.end == -1 {
		end.end = len(p.dots)
	}

	// do aligning
	if p.Width != 0 && p.Align != Left {
		if p.lines[0].end == 1 { // no content case
			p.dots[0].X += p.Width * float64(p.Align)
		} else {
			for _, l := range p.lines {
				end := l.end*4 - 4
				shift := (p.Width - p.data.Vertexes[end-1].Pos.X) * float64(p.Align)
				for i := l.start * 4; i < end; i++ {
					p.data.Vertexes[i].Pos.X += shift
				}
				for i := l.start; i < l.end; i++ {
					p.dots[i].X += shift
				}
			}
		}
	}

	for _, e := range p.instant { //instant effects are applied to base data
		e.Apply(p.data.Vertexes, 0)
	}
}

// ResolveChunks gets rid of nested FontEffects as nesting of then does not make sense
// it turns ranges like 0-10 3-7 to 0-3 3-7 7-10 so no ranges overlap.
func (m *Markdown) ResolveChunks(p *Paragraph) {
	fef := NFontEffect(p.Font, 0, len(p.Compiled))
	p.chunks = p.chunks[:0]
	if p.NoEffects {
		p.chunks = append(p.chunks, fef)
		return
	}

	m.buff = m.buff[:0]
	m.stack = m.stack[:0]
	m.stack = append(m.stack, fef)

	if !p.CustomLineheight {
		f := m.Fonts[p.Font]
		p.LineHeight = f.LineHeight()
		p.Ascent = f.Ascent()
		p.Descent = f.Descent()
	}

	for _, c := range p.chunks {
		if !p.CustomLineheight {
			f := m.Fonts[c.Font]
			p.LineHeight = math.Max(p.LineHeight, f.LineHeight())
			p.Ascent = math.Max(p.Ascent, f.Ascent())
			p.Descent = math.Max(p.Descent, f.Descent())
		}

		for len(m.stack) != 0 {
			l := m.stack.Last()
			if l.End > c.start {
				m.buff = append(m.buff, NFontEffect(l.Font, l.start, c.start))
				l.start = c.End
				break
			} else {
				m.buff = append(m.buff, m.stack.Pop())
			}
		}
		m.stack = append(m.stack, c)
	}

	m.stack.Reverse()
	m.buff = append(m.buff, m.stack...)
	m.buff.Filter(func(e *FontEffect) bool {
		return !e.Redundant()
	})

	p.chunks = append(p.chunks, m.buff...)
}

// LoadTTF loads TTF file into font.Face
func LoadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
