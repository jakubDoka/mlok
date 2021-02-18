package main

import (
	"gobatch/ggl"
	"gobatch/mat"
)

func main() {
	window, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	t, err := ggl.LoadTexture("square.png")
	if err != nil {
		panic(err)
	}

	b := ggl.NBatch(t, nil, nil)

	w, h := float64(t.W/4), float64(t.H/4)

	n := ggl.NNinePatchSprite(t.Frame(), mat.A(w, h, w, h))

	n.Resize(1000, 600)

	n.Draw(b, mat.IM.Scaled(mat.Origin, .5), mat.Alpha(1))

	b.Draw(window)

	for !window.ShouldClose() {
		window.Update()
	}
}
