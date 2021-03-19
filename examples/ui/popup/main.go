package main

import (
	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/ggl/ui"
	"github.com/jakubDoka/gobatch/logic/frame"
	"github.com/jakubDoka/gobatch/mat/rgba"
)

func main() {
	win, err := ggl.NWindow(&ggl.WindowConfig{
		Width:     1000,
		Height:    600,
		Resizable: true, // to showcase that resizing window is responsive
	})

	if err != nil {
		panic(err)
	}

	scene := ui.NScene() // comes with default font and parser

	scene.AddGoss([]byte(`
	label{
		margin: 10; 
		size: fill; 
		background: .5; 
		text_margin: fill;
		text_scale: 4;
		text_color: black;
		text_size: 0;
	}
	popup_button{
		all_masks: gray;
		hover_mask: green;
		size: 100 100;
		text_scale: 2;
		text_margin: fill;
	}
	`))

	if err != nil {
		panic(err)
	}

	// this can be called on any element but root is most obvious one
	err = scene.Root.AddGoml([]byte(`
	<div styles="label">this is here just</>
	<div styles="label">to show how popup</>
	<div styles="label" style="
		composition: horizontal; 
		text_margin: fill fill 0 fill;
	">
		works, just click this 
		<button id="opener" styles="label" style="
			padding: 10;
			margin: 0 fill fill fill; 
			size: 0;
			hover_mask: red;
		">
			button
		</>
	</>
	<#><div styles="label" style="
		margin: fill;
		size: 300;
		text_margin: fill 0;
		text_size: 0;
	">
		This is so called popup. Do you want to exit?
		<div style="
			composition: horizontal;
		">
			<button styles="popup_button">yes</>
			<button styles="popup_button">no</>
		</>
	</><#>
	`))

	if err != nil {
		panic(err)
	}

	// Processor just performs actions on scene, you can easily switch between scenes
	processor := ui.Processor{}
	processor.SetScene(scene)

	// delta time (frame time) is needed for processor update
	ticker := frame.Delta{}

	for !win.ShouldClose() {
		// we have lot of stupid colors
		win.Clear(rgba.BabyBlue)

		// window is resizable so we have to update frame, though if you pass same frame twice
		// processor will not resize elements
		processor.SetFrame(win.Frame())
		processor.Update(win, ticker.Tick())
		processor.Render(win)

		win.Update()
	}
}
