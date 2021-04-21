package main

import (
	_ "image/png"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/pck"
	"github.com/jakubDoka/mlok/ggl/ui"
	"github.com/jakubDoka/mlok/logic/frame"
	"github.com/jakubDoka/mlok/mat"
)

func main() {
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

	if err != nil {
		panic(err)
	}

	sheet.Pack()

	s.SetSheet(sheet)

	s.Parser = ui.NParser()

	err = s.Root.AddGoml([]byte(`
	<div style="background: .5; size: fill;composition: horizontal;margin: 10;">
		<div style="background: .5;size: fill;margin: 10;text_size: 0;">
			<div style="background: .5 .5 .5 .5;size: fill;"/>
			hello
			<div style="background: .5 0 .5 .5;size: fill;"/>
		</>
		<div style="background: .5 .5 0 .5;size: fill;"/>
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

	gl.Enable(gl.SCISSOR_TEST)
	gl.Scissor(0, 0, 1000, 1000)

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

// 14016
// specificky simbol 5757
