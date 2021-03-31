package txt

import (
	"fmt"
	"image"
	"math"
	"testing"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/rgba"
)

func TestAtlas(t *testing.T) {
	ttf, err := LoadTTF("C:/Users/jakub/Documents/programming/golang/src/github.com/jakubDoka/mlok/ggl/txt/test_data/else.ttf", 100)
	if err != nil {
		panic(err)
	}

	atlas := NAtlas("", ttf, 0, ASCII)
	if atlas.Pic.Bounds().Min != image.ZP {
		t.Error(atlas.Pic.Bounds(), atlas.descent, atlas.ascent, atlas.lineHeight)
	}
}

func TestDrawer(t *testing.T) {
	win, err := ggl.NWindow(nil)
	if err != nil {
		return
	}

	batch := ggl.Batch{
		Texture: ggl.NTexture(Atlas7x13.Pic, false),
	}

	drawer := NDrawer(Atlas7x13)
	text := Text{}

	drawer.Draw(&text, "Hello world!")
	fmt.Println(text.Bounds)
	text.DrawCentered(&batch, mat.M(mat.ZV, mat.V(10, 10), math.Pi*.5), rgba.AirForceBlueRaf)

	win.SetCamera(mat.IM.Scaled(mat.ZV, 1))

	batch.Draw(win)

	for !win.ShouldClose() {
		win.Update()
	}

	t.Fail()
}
