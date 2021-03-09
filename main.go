package main

import (
	"gobatch/ggl"
	"gobatch/ggl/pck"
	"gobatch/ggl/ui"
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
	<div style="background: 1;size: fill;composition: horizontal;margin: 10;">
		<div style="background: 0.5;size: fill;margin: 10;"> 
			<div style="background: 0.5;size: fill;"/>
			hello
			<div style="background: 0.5;size: fill;"/>
		</>
		<div style="background: 0;size: fill;"/>
	</>
	<div style="background: 0.5;size: fill;"/>
	`))

	if err != nil {
		panic(err)
	}

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
