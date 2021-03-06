package main

import (
	"gobatch/ggl"
	"gobatch/ggl/pck"
	"gobatch/ggl/txt"
	"gobatch/ggl/ui"
	"gobatch/logic/frame"
	"gobatch/mat"
)

func main() {
	// creates window with default config
	window, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	s := ui.NScene()

	s.SetSheet(&pck.Sheet{
		Pic: txt.Atlas7x13.Pic,
	})

	s.Parser = ui.NParser()

	err = s.Root.AddGoml([]byte(`
	<div style="
		background: red;
		size: fill;
	">
		<scroll style="
			resize_mode: ignore;
			size: 400 fill;
			margin: 100;
			background: black;
			friction: 5;
			bar_x: true;
			bar_y: true;
			scroll_sensitivity: 5;
			bar_rail_color: wheat;
		">
			<div id="h" style="
				text_scale: 5f; 
				background: almond; 
				text_margin: 10 0;
				text_color: gray;
				size: fill;
			">
				There hes to be a way to write code that makes it more pleasant and bug free,
				though i probably will never find out judging from how stupid i am already. 
				Newertheless this cost me lot of pain but here it is, some decent ui system.
				All elements can be defined in #FF00FF[goml] and styled with #FF00FF[goss]. 
			</>
		</>
	</>
	`))

	if err != nil {
		panic(err)
	}

	p := ui.Processor{}

	p.SetScene(s)

	p.SetFrame(window.Frame())

	d := frame.Delta{}.Init()

	// stay open and update so os will not report issue
	for !window.ShouldClose() {
		d.Log(1)

		p.Update(window, d.Tick())

		// making background white
		window.Clear(mat.Green)

		// drawing ui to window
		p.Render(window)
		window.Update()
	}
}

// 70203421
// noK.uqo.3.ixo
