package main

import (
	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/ui"
	"github.com/jakubDoka/mlok/logic/frame"
	"github.com/jakubDoka/mlok/mat/rgba"
)

// example showcases how you can create a popup with ggl/ui and also use
// predefined styles
func main() {
	win, err := ggl.NWindow(&ggl.WindowConfig{
		Width:     1000,
		Height:    600,
		Resizable: true,
	})

	if err != nil {
		panic(err)
	}

	scene := ui.NScene() // comes with default font and parser

	scene.LoadGoss("style.goss")
	if err != nil {
		panic(err)
	}
	err = scene.Root.LoadGoml("scene.goml")
	if err != nil {
		panic(err)
	}

	opened := true

	opener := scene.ID("opener")
	poppup := scene.ID("poppup")
	yes := scene.ID("yes")
	no := scene.ID("no")

	opener.Listen(ui.Click, func(i interface{}) {
		poppup.Show()
	})
	yes.Listen(ui.Click, func(i interface{}) {
		opened = false
	})
	no.Listen(ui.Click, func(i interface{}) {
		poppup.Hide()
	})

	processor := ui.Processor{}
	processor.SetScene(scene)

	ticker := frame.Delta{}

	for !win.ShouldClose() && opened {
		win.Clear(rgba.BabyBlue)

		processor.SetFrame(win.Frame())
		processor.Update(win, ticker.Tick())
		processor.Render(win)

		win.Update()
	}
}
