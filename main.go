package main

import (
	"gobatch/ggl"
	"gobatch/ggl/ui"
	"gobatch/mat"
	"strconv"
)

func main() {
	// creates window with default config
	window, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	// defining assets, assets can contain spritesheet, markdown for text drawing and styles
	a := ui.Assets{
		Styles: map[string]*ui.Style{
			"": {
				Subs: []ui.SubStyle{
					{
						Background: mat.Blue,
					},
				},
			},
			"stl": {
				Size:   mat.V(ui.Fill, ui.Fill),
				Margin: mat.A(10, 10, 10, 10),
				Subs: []ui.SubStyle{
					{
						Background: mat.Green,
					},
				},
			},
		},
	}

	// ui processor is Root Div wrapper for easier use
	p := ui.NProcessorFromBatch(ggl.NBatch(nil, nil, nil), &a)
	// adding 10 divs with given style, there can be multiple stiles separated by spaces
	for i := 0; i < 10; i++ {
		p.Root.children.Put(strconv.Itoa(i), &ui.Div{Styles: "stl"})
	}

	// i have to address these
	p.Root.Init(p)
	p.SetFrame(window.Frame().Moved(window.Frame().Center().Inv()))

	// making background white
	window.Clear(mat.White)

	// drawing ui to window
	p.Render(window)

	// stay open and update so os will not report issue
	for !window.ShouldClose() {
		window.Update()
	}
}
