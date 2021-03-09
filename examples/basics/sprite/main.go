package main

import (
	"gobatch/ggl"
	"gobatch/logic/frame"
	"gobatch/mat"
	_ "image/jpeg"
	"math"
)

func main() {
	// example shows how you can render the sprite on s screen
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	// loading texture from disk and setting up opengl object in one step
	// returned value is pointer with some data simple data like texture size
	texture, err := ggl.LoadTexture("../../assets/plane.jpg")
	if err != nil {
		panic(err)
	}

	// Batch only holds a components together, does not need any initialization,
	// not even texture is needed if you intend to draw only triangles
	batch := ggl.Batch{Texture: texture}

	// Sprite has to be set up but don't get misled by constructor, creating sprite
	// is very cheap as sprite live on stack (simplyfed, of course go have no stack or heap)
	sprite := ggl.NSprite(texture.Frame())

	// little helper, it measures a frame time
	ticker := frame.Delta{}

	var time, rotation float64

	for !win.ShouldClose() {
		delta := ticker.Tick()
		time += delta * .5               // tracking the time program runs for / 2 because then it looks better
		rotation += delta * math.Pi * .2 // one rotation per 10 seconds

		// some math magic
		scale := mat.V(math.Sin(time), math.Cos(time)).Scaled(.3)
		position := mat.Rad(-rotation, 220)
		color := mat.RGB(
			math.Abs(math.Sin(time)),
			math.Abs(math.Sin(time+math.Pi*.3)),
			math.Abs(math.Sin(time+math.Pi*.6)),
		)

		// we update the sprite state so its projection changes
		sprite.Update(mat.M(position, scale, rotation), color)

		// We then draw sprite in its current stats to batch batch is composed of ggl.Data
		// which accepts the vertex data sprite produces, all it does is appending to slice
		sprite.Fetch(&batch)

		// now we are getting little too fancy
		win.Clear(color.Inverted())

		// batch now draws to window which does the draw call
		batch.Draw(win)
		// don't forget to clear batch or you will run out of memory
		batch.Clear()
		// also important or you will end up with frozen window
		win.Update()
	}
}
