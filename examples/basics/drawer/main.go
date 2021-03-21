package main

import (
	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/ggl/drw"
	"github.com/jakubDoka/gobatch/mat"
	"github.com/jakubDoka/gobatch/mat/rgba"
)

// this example demonstrates drawing capability of drw.Geom
func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	geom := drw.Geom{}
	geom.Restart() // set values to defaults

	// circles and arcs
	geom.Color(rgba.Blue).Arc(0, 5).Circle(mat.C(0, 0, 40))
	geom.Color(rgba.Red).Resolution(8).Arc(3, 8).Fill(false).Circle(mat.C(0, 0, 70))
	geom.Color(rgba.Yellow).Resolution(100).Arc(0, 0).Thickness(5).Circle(mat.C(0, 0, 100))

	// rectangles
	geom.Color(rgba.Purple).Fill(true).AABB(mat.A(100, 100, 500, 300))
	geom.Color(rgba.Amber).Fill(false).AABB(mat.A(200, 50, 300, 150))
	geom.Color(rgba.HanBlue).Thickness(10).LineType(drw.Sharp).AABB(mat.A(250, -100, 400, 200))

	// lines
	geom.Color(rgba.White).LineType(drw.Round).Thickness(20).Line(
		mat.V(-400, -100),
		mat.V(-100, 200),
		mat.V(100, 100),
		mat.V(100, 250),
	)

	batch := ggl.Batch{}

	geom.Fetch(&batch)
	win.Clear(rgba.Black)

	batch.Draw(win)

	for !win.ShouldClose() {
		win.Update()
	}
}
