package main

import (
	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/ggl/ui"
	"github.com/jakubDoka/gobatch/logic/events"
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
		text_scale: 2;
		size: 100 0;
		margin: 100 0;
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
			margin: 10 fill fill fill;
			size: 0;
			text_margin: 0;
			padding: 10;
			hover_mask: red;
		">
			button
		</>
	</>
	<div id="poppup" styles="label" hidden style="
		relative: true; 
		background: .5 .5 0;
		size: 0;
		margin: fill;
	">
		Do you want to exit?
		<div style="
			composition: horizontal;
			margin: fill 10;
		">
			<button id="yes" styles="popup_button">yes</>
			<button id="no" styles="popup_button">no</>
		</>
	</>
	`))

	if err != nil {
		panic(err)
	}

	opened := true

	opener := scene.ID("opener")
	poppup := scene.ID("poppup")
	yes := scene.ID("yes")
	no := scene.ID("no")

	opener.Events.Add(&events.Listener{
		Name: ui.Click,
		Runner: func(i interface{}) {
			poppup.Show()
		},
	})

	yes.Events.Add(&events.Listener{
		Name: ui.Click,
		Runner: func(i interface{}) {
			opened = false
		},
	})

	no.Events.Add(&events.Listener{
		Name: ui.Click,
		Runner: func(i interface{}) {
			poppup.Hide()
		},
	})

	// Processor just performs actions on scene, you can easily switch between scenes
	processor := ui.Processor{}
	processor.SetScene(scene)

	// delta time (frame time) is needed for processor update
	ticker := frame.Delta{}

	for !win.ShouldClose() && opened {
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
