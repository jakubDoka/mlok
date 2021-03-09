package main

import (
	"fmt"
	"gobatch/ggl"
	"gobatch/ggl/pck"
	"gobatch/ggl/ui"
	"gobatch/logic/events"
	"gobatch/logic/frame"
	"gobatch/mat"
	_ "image/png"
)

func main() {
	// creates window with default config
	window, err := ggl.NWindow(&ggl.WindowConfig{
		Resizable: true,
		Width:     1000,
		Height:    600,
	})
	if err != nil {
		panic(err)
	}

	s := ui.NScene()

	sheet := pck.Sheet{}

	sheet.AddMarkdown(s.Assets.Markdowns["default"])
	err = sheet.AddImages("square.png")

	if err != nil {
		panic(err)
	}

	sheet.Pack()

	s.SetSheet(&sheet)

	s.Parser = ui.NParser()

	err = s.Root.AddGoml([]byte(`
	<div style="
		background: red;
		size: fill;
		composition: horizontal;
	">
		<area style="size: 100; margin: 0 20;background: black;"/>
		<scroll style="
			resize_mode: ignore;
			size: 800 400;
			margin: fill;
			background: green;
			friction: 5;
			bars: true;
			scroll_sensitivity: 5;
			bar_color: wheat;
		">
			<button id="button" style="
				all_regions: square;
				text_margin: fill;
				text_scale: 10;
				text_size: 0; 
				patch_scale: 0.4 0.4;
				size: fill;
			" idle_text="idle" hover_text="hover" pressed_text="pressed"/>
			
			<div style="
				text_scale: 1; 
				text_margin: fill;
				text_size: 100 fill;
				background: almond; 
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

	button := s.ID("button")
	button.Events.Add(&events.Listener{
		Name: ui.Click,
		Runner: func(i interface{}) bool {
			fmt.Println("button clicked")
			return false
		},
	})

	p := ui.Processor{}

	p.SetScene(s)

	d := frame.Delta{}.Init()
	l := frame.Limitter{}

	l.SetFPS(60)

	// stay open and update so os will not report issue
	for !window.ShouldClose() {
		d.Log(1)
		l.Regulate()

		p.Update(window, d.Tick())

		p.SetFrame(window.Frame())
		// making background white
		window.Clear(mat.Green)

		// drawing ui to window
		p.Render(window)
		window.Update()
	}
}

// 70203421
// noK.uqo.3.ixo
