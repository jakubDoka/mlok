package main

import (
	"gobatch/ggl"
	"gobatch/mt"
	_ "image/png"
	"log"
	"segmentation/src/core/tm"

	"golang.org/x/image/colornames"
)

const windowWidth = 800
const windowHeight = 600

func main() {
	window, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	texture, err := ggl.LoadTexture("square.png")
	if err != nil {
		log.Fatalln(err)
	}

	batch := ggl.NBatch2D(texture, nil, nil)
	sprite := ggl.NSprite2D(texture.Frame())
	tm := tm.NTime()
	for !window.ShouldClose() {
		window.Clear2D(mt.ToRGBA(colornames.Aliceblue))

		for i := 0; i < 50000; i++ {
			sprite.Draw(batch, mt.NMat2(mt.NV2(1000, 1000), mt.NV2(.5, .5), 2), mt.Alpha(1))
		}

		batch.Draw(window)
		batch.Clear()

		window.Update()
		tm.Update()
	}
}
