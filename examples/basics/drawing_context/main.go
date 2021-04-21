package main

import (
	"math"

	_ "image/png"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/drw/context"
	"github.com/jakubDoka/mlok/logic/frame"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/rgba"
)

func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	texture, err := ggl.LoadTexture("square.png")
	if err != nil {
		panic(err)
	}

	// batch uses texture to render
	batch := ggl.Batch{
		Texture: texture,
	}

	var con context.C

	con.Init([]context.PartDefs{
		{
			Region: texture.Frame(),
			Mask:   rgba.Green,
			Scale:  mat.V(1, 1),
		},
		{
			Region:   texture.Frame(),
			Offset:   mat.V(100, 100),
			Rotation: 10,
			Mask:     rgba.Red,
			Scale:    mat.V(1, .5),
		},
		{
			Region:   texture.Frame(),
			Offset:   mat.V(-100, -100),
			Rotation: 30,
			Mask:     rgba.Wheat,
			Scale:    mat.V(.5, 1),
		},
	}...)

	ticker := frame.Delta{}

	var time float64
	for !win.ShouldClose() {
		time += ticker.Tick()
		ticker.Log(1)

		// i love math
		con.Draw(
			&batch,
			mat.M( // transformation encoded in matrix
				mat.V(100, 100),
				mat.V(math.Sin(time), math.Cos(time)),
				time/2,
			),
			mat.Alpha(math.Abs(math.Cos(time))),
		)

		win.Clear(rgba.AirForceBlueRaf)

		batch.Draw(win)
		batch.Clear() // ype clear goes here

		win.Update()
	}
}
