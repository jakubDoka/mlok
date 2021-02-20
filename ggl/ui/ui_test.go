package ui

import (
	"gobatch/mat"
	"strconv"
	"testing"
)

func TestCalcMarginSize(t *testing.T) {
	p := Processor{}
	s := []Style{
		{},
		{
			Margin: mat.A(Fill, Fill, Fill, Fill),
			Size:   mat.V(10, 10),
		},
		{
			Size:       mat.V(100, 100),
			Horizontal: true,
		},
		{
			Margin: mat.A(0, Fill, 0, Fill),
			Size:   mat.V(10, 10),
		},
		{
			Size: mat.V(100, 100),
		},
		{
			Margin: mat.A(Fill, 0, Fill, 0),
			Size:   mat.V(10, 10),
		},
		{
			Margin: mat.A(10, 10, 10, 10),
			Size:   mat.V(10, 10),
		},
		{
			ResizeMode: Exact,
		},
		{
			Margin: mat.A(10, 10, 10, 10),
			Size:   mat.V(Fill, Fill),
		},
		{
			Size: mat.V(90, 90),
		},
	}

	testCases := []struct {
		desc          string
		parent        Style
		input, result []Div
	}{
		{
			desc:   "all fill",
			parent: s[2],
			input: []Div{
				{Style: s[1]},
				{Style: s[1]},
			},
			result: []Div{
				{margin: mat.A(20, 45, 20, 45), Frame: mat.A(20, 45, 30, 55)},
				{margin: mat.A(20, 45, 20, 45), Frame: mat.A(70, 45, 80, 55)},
			},
		},
		{
			desc:   "combined",
			parent: s[2],
			input: []Div{
				{Style: s[1]},
				{Style: s[3]},
				{Style: s[1]},
			},
			result: []Div{
				{margin: mat.A(17.5, 45, 17.5, 45), Frame: mat.A(17.5, 45, 27.5, 55)},
				{margin: mat.A(0, 45, 0, 45), Frame: mat.A(45, 45, 55, 55)},
				{margin: mat.A(17.5, 45, 17.5, 45), Frame: mat.A(72.5, 45, 82.5, 55)},
			},
		},
		{
			desc:   "vertical combined",
			parent: s[4],
			input: []Div{
				{Style: s[1]},
				{Style: s[5]},
				{Style: s[1]},
			},
			result: []Div{
				{margin: mat.A(45, 17.5, 45, 17.5), Frame: mat.A(45, 17.5, 55, 27.5)},
				{margin: mat.A(45, 0, 45, 0), Frame: mat.A(45, 45, 55, 55)},
				{margin: mat.A(45, 17.5, 45, 17.5), Frame: mat.A(45, 72.5, 55, 82.5)},
			},
		},
		{
			desc:   "shrink",
			parent: s[7],
			input: []Div{
				{Style: s[6]},
				{Style: s[6]},
				{Style: s[6]},
			},
			result: []Div{
				{margin: s[6].Margin, Frame: mat.A(10, 10, 20, 20)},
				{margin: s[6].Margin, Frame: mat.A(10, 40, 20, 50)},
				{margin: s[6].Margin, Frame: mat.A(10, 70, 20, 80)},
			},
		},
		{
			desc:   "width",
			parent: s[9],
			input: []Div{
				{Style: s[8]},
				{Style: s[8]},
				{Style: s[8]},
			},
			result: []Div{
				{margin: s[8].Margin, Frame: mat.A(10, 10, 80, 20)},
				{margin: s[8].Margin, Frame: mat.A(10, 40, 80, 50)},
				{margin: s[8].Margin, Frame: mat.A(10, 70, 80, 80)},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p.Root.Style = tC.parent
			for i, ch := range tC.input {
				ch := ch
				p.Root.children.Put(strconv.Itoa(i), &ch)
			}
			p.SetFrame(p.Root.Size.ToAABB())
			for i, v := range p.Root.children.Slice() {
				v := v.Value
				if v.margin != tC.result[i].margin {
					t.Error(i, "margin", v.margin, "!=", tC.result[i].margin)
				}
				if v.Frame != tC.result[i].Frame {
					t.Error(i, "frame", v.Frame, "!=", tC.result[i].Frame)
				}
			}
		})
	}
}
