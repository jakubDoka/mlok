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

	// this can be called on any element but root is most obvious one
	err = scene.Root.AddGoml([]byte(`
	<div style="
		composition: horizontal;
		size: fill;
	">
		<div id="chat" style="
			background: almond;
			size: 500 fill;
			resize_mode: ignore;
		"> 
			<text text="The Chat" style="
				text_color: 1 0.5 0.5;
				text_size: 0;
				text_scale: 4;
				text_margin: fill 0;
			"/>
			<scroll style="
				size: fill;
				resize_mode: ignore;
				bar_y: true;
				padding: 0 0 20 0;
				background: 0.3;
				margin: 10;
				text_scale: 2;
			">
				<text id="chat-text" name="target"/>
			</>
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
				<button id="send" style="
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

	// now we can used ids to access elements
	chat := scene.ID("chat")
	// as ui.Module embeds ui.Element you don't have to keep two variables in
	// case you need access to module and element
	chatText := scene.ID("chat-text").Module.(*ui.Text)
	alter := scene.ID("alter").Module.(*ui.Button)
	input := scene.ID("input").Module.(*ui.Area)
	send := scene.ID("send")

	// registering event listener for hiding the chat. Why listeners and not just runners?
	// Well if you store your listener you can remove it by calling its method,
	// you can also change Anything except Name and it will have an effect
	alter.Events.Add(&events.Listener{
		Name: ui.Click,
		Runner: func(i interface{}) {
			hidden := chat.Hidden()
			if hidden {
				alter.SetText("hide")
			} else {
				alter.SetText("show")
			}
			chat.SetHidden(!hidden)

		},
	})

	// listenner for message send event, its a minimal logic
	send.Events.Add(&events.Listener{
		Name: ui.Click,
		Runner: func(i interface{}) {
			// text grows from bottom to top so one line above first message is unnoticeable
			chatText.Content = append(chatText.Content, '\n')
			// yes we have markdown and ']]', if used inside of markdown closure, is replaced with ']'
			chatText.Content = append(chatText.Content, []rune("#ff00ff[[you]]:] ")...)
			chatText.Content = append(chatText.Content, input.Content...)
			input.SetText("") // clear the input
		},
	})

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
