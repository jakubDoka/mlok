package drw

import (
	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/mat"
)

type LineProcessor struct {
	Points  []mat.Vec
	Indices ggl.Indices
}

func (l *LineProcessor) Process(g *Geom, ld LineDrawer, points ...mat.Vec) {
	if ld == nil {
		ld = Default
	}

	ld.Init(g.thickness)

	l.Indices.Clear()
	l.Points = l.Points[:0]

	ln := len(points)
	rl := ln - 2
	if g.loop {
		rl++
	} else {
		var e End
		e.Init(points[0], points[1], g.thickness)
		ld.Start(&e, l)
	}

	var e Edge
	for i := 0; i < rl; i++ {
		e.Init(points, g.thickness, i, ln)
		ld.Edge(&e, l)
	}

	if g.loop {
		e.Init(points, g.thickness, rl, ln)
		ld.Close(&e, l)
	} else {
		var e End
		e.Init(points[ln-1], points[ln-2], g.thickness)
		ld.End(&e, l)
	}
}

func (l *LineProcessor) AppendPoints(center mat.Vec, points ...mat.Vec) {
	l.Points = append(l.Points, points...)
	slice := l.Points[len(l.Points)-len(points):]
	for i := range slice {
		slice[i].AddE(center)
	}
}

func (l *LineProcessor) AppendIndices(indices ...uint32) {
	l.Indices = append(l.Indices, indices...)
	l.Indices[len(l.Indices)-len(indices):].Shift(uint32(len(l.Points)))
}

type RoundLine struct {
	LineDrawerBase
	Circle
}

func (r *RoundLine) Init(thickness float64) {
	r.Filled(thickness, 0, 0, AutoResolution(thickness, 0, 0, 1))
}

func (r *RoundLine) Start(e *End, lp *LineProcessor) {
	lp.Indices = append(lp.Indices, r.Indices...)
	lp.AppendPoints(e.A, r.Base...)
	r.LineDrawerBase.Start(e, lp)
}

func (r *RoundLine) End(e *End, lp *LineProcessor) {
	r.LineDrawerBase.End(e, lp)
	lp.AppendIndices(r.Indices...)
	lp.AppendPoints(e.A, r.Base...)
}

func (r *RoundLine) Edge(e *Edge, lp *LineProcessor) {
	r.PreEdge(e, lp)
	lp.AppendIndices(r.Indices...)
	lp.AppendPoints(e.B, r.Base...)
	r.PostEnge(e, lp)
}

func (r *RoundLine) Close(e *Edge, lp *LineProcessor) {
	r.LineDrawerBase.Close(e, lp)
	lp.AppendIndices(r.Indices...)
	lp.AppendPoints(e.B, r.Base...)
}

type SharpLine struct{ LineDrawerBase }

func (s *SharpLine) Start(e *End, lp *LineProcessor) {
	lp.Indices = append(lp.Indices, 0, 1, 2)
	s.AddEnd(e, lp)
	s.LineDrawerBase.Start(e, lp)
}

func (s *SharpLine) End(e *End, lp *LineProcessor) {
	s.LineDrawerBase.End(e, lp)
	l := uint32(len(lp.Points))
	lp.Indices = append(lp.Indices, l-2, l-1, l)
	s.AddEnd(e, lp)
}

func (s *SharpLine) AddEnd(e *End, lp *LineProcessor) {
	lp.Points = append(lp.Points, e.A.Add(e.A.To(e.B1).Normal()))
}

func (s *SharpLine) Edge(e *Edge, lp *LineProcessor) {
	s.PreEdge(e, lp)
	l := uint32(len(lp.Points))
	if e.Convex {
		lp.Indices = append(lp.Indices, l-2, l-1, l)
	} else {
		lp.Indices = append(lp.Indices, l-2, l-1, l+1)
	}
	s.PostEnge(e, lp)
}

func (s *SharpLine) Close(e *Edge, lp *LineProcessor) {
	s.LineDrawerBase.Close(e, lp)
	l := uint32(len(lp.Points))
	if e.Convex {
		lp.Indices = append(lp.Indices, l-1, l-2, l-4)
	} else {
		lp.Indices = append(lp.Indices, l-1, l-2, l-3)
	}
}

type LineDrawerBase struct{}

func (*LineDrawerBase) Init(thickness float64) {}

func (*LineDrawerBase) Start(e *End, lp *LineProcessor) {
	lp.AppendIndices(LineIndicePatern...)
	lp.Points = append(lp.Points, e.B1, e.B2)
}

func (*LineDrawerBase) End(e *End, lp *LineProcessor) {
	lp.Points = append(lp.Points, e.B2, e.B1)
}

func (l *LineDrawerBase) Edge(e *Edge, lp *LineProcessor) {
	l.PreEdge(e, lp)
	l.PostEnge(e, lp)
}

func (*LineDrawerBase) PostEnge(e *Edge, lp *LineProcessor) {
	lp.AppendIndices(LineIndicePatern...)
	lp.Points = append(lp.Points, e.B3, e.B4)
}

func (*LineDrawerBase) PreEdge(e *Edge, lp *LineProcessor) {
	lp.Points = append(lp.Points, e.B1, e.B2)
}

func (l *LineDrawerBase) Close(e *Edge, lp *LineProcessor) {
	l.Edge(e, lp)
	ln := len(lp.Indices)
	lp.Indices[ln-1] = 0
	lp.Indices[ln-2] = 1
	lp.Indices[ln-4] = 1
}

var LineIndicePatern = ggl.Indices{0, 1, 3, 0, 3, 2}

type LineDrawer interface {
	Init(thickness float64)
	Start(*End, *LineProcessor)
	End(*End, *LineProcessor)
	Edge(*Edge, *LineProcessor)
	Close(*Edge, *LineProcessor)
}

// Line Modes
var (
	Default = &LineDrawerBase{}
	Sharp   = &SharpLine{}
	// Round holds some inner state so if you are drawing on multiple threads, you need to create own instances
	Round = &RoundLine{}
)

type End struct {
	A, B, BA, B1, B2 mat.Vec
	Thickness        float64
}

func (e *End) Init(a, b mat.Vec, thickness float64) {
	e.A, e.B, e.BA = a, b, b.To(a)
	nba := e.BA.Normal().Normalized().Scaled(thickness)
	e.B1, e.B2 = a.Sub(nba), a.Add(nba)
}

type Edge struct {
	End
	C, BC, B3, B4 mat.Vec
	Convex        bool
}

func (e *Edge) Init(points []mat.Vec, thickness float64, i, ln int) {
	e.A, e.B, e.C = points[i], points[(i+1)%ln], points[(i+2)%ln]
	e.BA, e.BC = e.B.To(e.A), e.B.To(e.C)
	e.Convex = e.BA.Cross(e.BC) >= 0
	nba, nbc := e.BA.Normal().Normalized().Scaled(thickness), e.BC.Normal().Normalized().Scaled(thickness)
	e.B1, e.B2, e.B3, e.B4 = e.B.Sub(nba), e.B.Add(nba), e.B.Add(nbc), e.B.Sub(nbc)
}
