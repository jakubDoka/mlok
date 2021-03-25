package main

import (
	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/drw"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/rgba"
)

// this example demonstrates automatic resolution for circles
func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	geom := drw.Geom{}
	geom.Restart()

	radius := 100.
	for !win.ShouldClose() {
		batch := ggl.Batch{}

		radius += win.MouseScroll().Y

		geom.Spacing(40).Circle(mat.C(250, 0, radius))
		geom.Spacing(1).Circle(mat.C(-250, 0, radius))

		geom.Fetch(&batch)
		geom.Clear()

		win.Clear(rgba.Black)

		batch.Draw(win)
		batch.Clear()

		win.Update()
	}
}
