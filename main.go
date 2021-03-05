package main

import (
	"gobatch/ggl"
	"gobatch/ggl/pck"
	"gobatch/ggl/txt"
	"gobatch/ggl/ui"
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
	<div style="background: 1f 0f 0f;size: fill;">
		<div style="background: 1f; text_color: 0f 0f 0f;">
			hello
		</>
		<div style="text_scale: 10f; background: 1f 0f 1f; text_margin: 0f;">
			there
		</>
	</>
	`))

	if err != nil {
		panic(err)
	}

	p := ui.Processor{}

	p.SetScene(s)

	p.SetFrame(window.Frame())

	p.Update(window, 0)

	// making background white
	window.Clear(mat.Green)

	// drawing ui to window
	p.Render(window)

	// stay open and update so os will not report issue
	for !window.ShouldClose() {
		window.Update()
	}
}

// 70203421
// noK.uqo.3.ixo
