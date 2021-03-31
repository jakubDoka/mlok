package main

import (
	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/pck"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/rgba"

	_ "image/png"
)

func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	sheet := pck.Sheet{}

	sheet.AddImages("fish.png", "face_32_32.png", "something.png")

	sheet.Pack()

	batch := sheet.Batch()
	sprite, _ := sheet.Sprite("All")
	face1, _ := sheet.Sprite("face1")
	something, _ := sheet.Sprite("something")

	win.Clear(mat.White)

	sprite.Draw(&batch, mat.IM.Scaled(mat.ZV, 2), rgba.White)
	face1.Draw(&batch, mat.IM.Scaled(mat.ZV, 5).Move(mat.V(300, 200)), rgba.White)
	something.Draw(&batch, mat.IM.Scaled(mat.ZV, 3).Move(mat.V(-300, 200)), rgba.White)
	batch.Draw(win)

	for !win.ShouldClose() {
		win.Update()
	}
}
