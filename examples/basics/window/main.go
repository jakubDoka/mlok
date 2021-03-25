package main

import (
	"math"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/logic/frame"
	"github.com/jakubDoka/mlok/mat"
)

func main() {
	win, err := ggl.NWindow(&ggl.WindowConfig{
		Width:                  1000,
		Height:                 600,
		Resizable:              true,
		TransparentFramebuffer: true,
	})
	if err != nil {
		panic(err)
	}

	ticker := frame.Delta{}
	var time float64
	for !win.ShouldClose() {
		time += ticker.Tick()

		win.Clear(mat.RGBA{
			R: math.Abs(math.Sin(time / 2)),
			G: math.Abs(math.Sin(time / 3)),
			B: math.Abs(math.Sin(time / 5)),
			A: math.Abs(math.Sin(time / 2)),
		})

		win.Update()
	}
}
