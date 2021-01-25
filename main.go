package main

import (
	"gobatch/ggl"
	"gobatch/ggl/txt"
	"gobatch/mat"
	"gogen/str"
)

func main() {
	window, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	m := txt.NMarkdown()

	t := ggl.NTexture(&txt.Atlas7x13.Pic)

	b := ggl.NBatch2D(t, nil, nil)

	p := txt.Paragraph{
		Width: 100,
		Text:  str.NString("Hello there, as you can see i can draw !green[green] !g[text] and also !33333300[transparent] text"),
	}

	m.Parse(&p)

	p.Update(0, 0, 0)

	p.Draw(b)

	window.SetCamera2D(mat.IM2.Scaled(mat.V2{}, 2))

	b.Draw(window)

	for !window.ShouldClose() {
		b.Draw(window)
		window.Update()
	}
}
