package main

import (
	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/ggl/ui"
	"github.com/jakubDoka/gobatch/logic/frame"
	"github.com/jakubDoka/gobatch/mat/rgba"
)

// example showcases ui capability of ggl/ui
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
	err = scene.Root.LoadGoml("scene.goml")

	if err != nil {
		panic(err)
	}

	// now we can used ids to access elements
	chat := scene.ID("chat")
	// as ui.Module embeds ui.Element you don't have to keep two variables in
	// case you need access to module and the element
	chatText := scene.ID("chat-text").Module.(*ui.Text)
	alter := scene.ID("alter").Module.(*ui.Button)
	input := scene.ID("input").Module.(*ui.Area)
	send := scene.ID("send")

	// registering event listener for hiding the chat. Why listeners and not just runners?
	// Well if you store your listener you can remove it by calling its method,
	// you can also change Anything except Name and it will have an effect
	alter.Listen(ui.Click, func(i interface{}) {
		hidden := chat.Hidden()
		if hidden {
			alter.SetText("hide")
		} else {
			alter.SetText("show")
		}
		chat.SetHidden(!hidden)
	})

	// listenner for message send event, its a minimal logic
	send.Listen(ui.Click, func(i interface{}) {
		// text grows from bottom to top so one line above first message is unnoticeable
		chatText.Content = append(chatText.Content, '\n')
		// yes we have markdown and ']]', if used inside of markdown closure, is replaced with ']'
		chatText.Content = append(chatText.Content, []rune("#ff00ff[[you]]:] ")...)
		chatText.Content = append(chatText.Content, input.Content...)
		chatText.Dirty()  // text will not update every frame, you have to notify about change
		input.SetText("") // clear the input
	})

	// Processor just performs actions on scene, you can easily switch between scenes
	// and transition should be instant
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
