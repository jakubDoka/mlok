package main

import (
	"fmt"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/ui"
	"github.com/jakubDoka/mlok/logic/frame"
	"github.com/jakubDoka/mlok/mat/rgba"
)

func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	// using this scene constructor does the font initialization for us
	scene := ui.NScene()
	// one way to add some elements to the scene, we are actually adding goml script here
	err = scene.Root.LoadGoml("scene.goml")
	if err != nil {
		panic(err)
	}

	// new stuff
	name := scene.ID("name").Module.(*ui.Area)
	ip := scene.ID("ip").Module.(*ui.Area)
	connect := scene.ID("connect")

	// listening to event
	connect.Listen(ui.Click, func(i interface{}) {
		fmt.Println(string(name.Content))
		fmt.Println(string(ip.Content))
	})

	// processor does not need any initialization we just have to set scene
	processor := ui.Processor{}
	processor.SetScene(scene)
	ticker := frame.Delta{}
	for !win.ShouldClose() {
		win.Clear(rgba.AirForceBlueRaf)

		processor.SetFrame(win.Frame())      // first we set ui frame to whole window
		processor.Update(win, ticker.Tick()) // then we perform update on all elements
		processor.Render(win)                // last step is to render all elements

		win.Update()
	}
}
