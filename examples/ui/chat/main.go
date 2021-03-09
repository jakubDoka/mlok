package main

import (
	"gobatch/ggl"
	"gobatch/ggl/ui"
	"gobatch/logic/events"
	"gobatch/logic/frame"
	"gobatch/mat/rgba"
)

func main() {
	win, err := ggl.NWindow(&ggl.WindowConfig{
		Width:     1000,
		Height:    600,
		Resizable: true,
	})

	if err != nil {
		panic(err)
	}

	scene := ui.NScene()

	err = scene.Root.AddGoml([]byte(`
	<div style="
		composition: horizontal;
		size: fill;
	">
		<scroll id="chat" style="
			background: almond;
			size: 500 fill;
			bars: true;
			outside: true;
			text_scale: 4; 
			resize_mode: ignore;
		"> 
			<text text="The Chat" style="
				scale: 4;
				text_color: 1 0.5 0.5;
				text_size: 0 100;
				text_margin: fill 10;
			"/>
			<text id="chat-text"/>
			<div style="
				composition: horizontal;
				margin: 10;
				size: fill 0;
				background: 1;
			">
				<area id="input" style="
					background: 0.3 0.3 0.3;
					text_scale: 2;
					size: fill 0;
				"/>
				<button style="
					all_masks: 1 0.5 0.5;
					hover_mask: green;
					size: 0 fill;
				" all_text="send"/>
			</>
			
		</>
		<button id="alter" style="
			all_masks: gray;
			hover_mask: green;
		" all_text="hide"/>
	</>
	`))

	if err != nil {
		panic(err)
	}

	chat := scene.ID("chat")
	//chatText := scene.ID("chat-text").Module.(*ui.Text)
	alter := scene.ID("alter").Module.(*ui.Button)

	alter.Events.Add(&events.Listener{
		Name: ui.Click,
		Runner: func(i interface{}) bool {
			hidden := chat.Hidden()
			if hidden {
				alter.SetText("hide")
			} else {
				alter.SetText("show")
			}
			chat.SetHidden(!hidden)
			return false
		},
	})

	processor := ui.Processor{}
	processor.SetScene(scene)

	ticker := frame.Delta{}

	for !win.ShouldClose() {
		win.Clear(rgba.BabyBlue)

		processor.SetFrame(win.Frame())
		processor.Update(win, ticker.Tick())
		processor.Render(win)

		win.Update()
	}
}
