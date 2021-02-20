package ui

import (
	"gobatch/ggl"
	"gobatch/ggl/dw"
	"gobatch/mat"
)

// Processor is wrapper that handles a Div, it stores some data global to all divs
// that can be reused for reduction of allocations
type Processor struct {
	Horizontal, redraw, resize bool

	Root   Div
	frame  mat.AABB
	canvas dw.Geom
	batch  *ggl.Batch
	Assets *Assets

	ids           map[string]*Div
	pfTmp, pfTmp2 []*float64
	divTemp       []*Div
}

// NProcessor creates ready-to-ues processor
func NProcessor(assets *Assets) *Processor {
	t := ggl.NTexture(assets.Pic, false)
	b := ggl.NBatch(t, nil, nil)
	return NProcessorFromBatch(b, assets)
}

// NProcessorFromBatch creates new processor with custom batch
func NProcessorFromBatch(batch *ggl.Batch, assets *Assets) *Processor {
	p := &Processor{
		batch:  batch,
		Assets: assets,
		redraw: true,
	}

	p.Root.children = NChildren()
	return p
}

// DivByID returns div or null if no div is under id
func (p *Processor) DivByID(id string) *Div {
	return p.ids[id]
}

// Fetch implements ggl.Fetcher interface
func (p *Processor) Fetch(t ggl.Target) {
	p.batch.Fetch(t)
}

// Render make Renderer render ui (yes)
func (p *Processor) Render(r ggl.Renderer) {
	p.batch.Draw(r)
}

// Redraw redraws ewerithing if needed
func (p *Processor) Redraw() {
	p.batch.Clear()
	p.Root.redraw(p.batch, &p.canvas)
	p.redraw = false
}

// SetFrame sets the frame of Processor witch also updates all elements inside
func (p *Processor) SetFrame(value mat.AABB) {
	p.frame = value
	p.Resize()
}

// Resize has to be called upon Frame change, its not recommended
// to call this manually, instead call Deformed() to notify processor
// about change
func (p *Processor) Resize() {
	p.Dirty()
	p.Root.size = p.frame.Size()
	p.Root.resize()
	p.Root.move(p.frame.Min, p.Horizontal)
	p.resize = false
}

// Dirty makes processor redraw all elements upon update call
// but after all divs are updated, best time to call this is during
// update or after Render
func (p *Processor) Dirty() {
	p.redraw = true
}

// Deformed makes processor resize all elements upon update call
// but after all divs are updated, best time to call this is during
// update or after Render
func (p *Processor) Deformed() {
	p.resize = true
}

// calcSize calculates all sizes equal to Fill amongst children of d
func (p *Processor) calcSize(d *Div) {
	offset, space, individualSpace := p.setup(d.Horizontal, d.size)

	for _, d := range d.children.Slice() {
		d := d.Value
		flt := d.Margin.Flatten()
		is := individualSpace
		for i, v := range flt {
			if i%2 == offset {
				if v != Fill {
					space -= v
				}
			} else {
				if v != Fill {
					is -= v
				}
			}
		}

		sft := d.Size.Flatten()
		smt := d.size.Mutator()

		if sft[1-offset] == Fill {
			*smt[1-offset] = is
		} else {
			*smt[1-offset] = sft[1-offset]
		}

		param := sft[offset]
		if param == Fill {
			p.pfTmp = append(p.pfTmp, smt[offset])
		} else {
			*smt[offset] = param
			space -= param
		}
	}
	feed(space, p.pfTmp)
}

// calculates margin in case it is Equal to Fill for all
// children of div
func (p *Processor) calcMargin(d *Div) {
	/*
		goal is to calculate how match free space is in div and divide it
		between margin equal to Fill that supports it, function
		collects the pointers to all fill fields and feeds tham with supposed
		values
	*/
	offset, space, individualSpace := p.setup(d.Horizontal, d.size)

	// little heck but i saves us repetitive logic
	for _, d := range d.children.Slice() {
		d := d.Value
		p.pfTmp2 = p.pfTmp2[:0]
		mut := d.margin.Mutator()
		flt := d.Margin.Flatten()
		is := individualSpace
		for i, v := range flt {
			// deciding how to treat margin value, notice how var offset relates to var horizontal
			if i%2 == offset {
				if v == Fill {
					p.pfTmp = append(p.pfTmp, mut[i])
				} else {
					*mut[i] = v
					space -= v
				}
			} else {
				if v == Fill {
					p.pfTmp2 = append(p.pfTmp2, mut[i])
				} else {
					*mut[i] = v
					is -= v
				}
			}

		}

		// subtracting the size of div, its like this because this way it works for
		// vertical and horizontal case
		sft := d.Frame.Size().Flatten()
		is -= sft[1-offset]
		space -= sft[offset]

		feed(is, p.pfTmp2)
	}

	feed(space, p.pfTmp)
}

func (p *Processor) setup(horizontal bool, size mat.Vec) (offset int, space, individualSpace float64) {
	p.pfTmp = p.pfTmp[:0]
	offset = 1
	individualSpace, space = size.XY()
	if horizontal {
		offset = 0
		space, individualSpace = size.XY()
	}

	return
}

// feed performs final space division between elements
func feed(space float64, targets []*float64) {
	if space <= 0 {
		// make sure they are zero, this gets rid of old values
		for _, v := range targets {
			*v = 0
		}
	} else {
		// split equally
		perOne := space / float64(len(targets))
		for _, v := range targets {
			*v = perOne
		}
	}
}
