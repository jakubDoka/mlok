package pck

import (
	"fmt"
	"image"
	"image/draw"
	"reflect"
	"testing"

	_ "image/png"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/txt"
)

func TestPathParsing(t *testing.T) {
	testCases := []struct {
		desc             string
		path, root, name string
		w, h             int
		ok               bool
	}{
		{
			desc: "normal",
			path: "hello/root/meme_30_30.png",
			root: "root",
			name: "meme",
			w:    30,
			h:    30,
			ok:   true,
		},
		{
			desc: "invalid",
			path: "hello/root/meme_E0_30.png",
			root: "root",
			name: "meme_E0_30",
			w:    0,
			h:    0,
			ok:   false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			name, w, h, ok := DetectSpritesheet(tC.path, tC.root)
			if name != tC.name || w != tC.w || h != tC.h || ok != tC.ok {
				t.Error(name, w, h, ok)
			}
		})
	}
}

func TestNSheet(t *testing.T) {
	ttf, err := txt.LoadTTF("C:/Users/jakub/Documents/programming/golang/src/github.com/jakubDoka/mlok/ggl/pck/test_data/else.ttf", 100)
	if err != nil {
		panic(err)
	}
	atlas := txt.NewAtlas("else", ttf, 0, txt.ASCII)

	fmt.Println(atlas.Pic.Bounds())

	sheet := Sheet{Root: "test_data"}

	sheet.Data = append(sheet.Data, PicData{
		Name: "else",
		Img:  atlas.Pic,
	})
	sheet.AddImages("C:/Users/jakub/Documents/programming/golang/src/github.com/jakubDoka/mlok/ggl/pck/test_data/nest.png")
	sheet.Pack()

	image1, err := ggl.LoadImage("C:/Users/jakub/Documents/programming/golang/src/github.com/jakubDoka/mlok/ggl/pck/test_data/nest.png")
	if err != nil {
		panic(err)
	}
	r := sheet.Regions["nest"].ToImage()

	image2 := image.NewNRGBA(r)
	draw.Draw(image2, image2.Bounds(), sheet.Pic, r.Min, 0)
	ggl.FlipNRGBA(image1)
	LogImage(image2)

	if !reflect.DeepEqual(image1.Pix, image2.Pix) {
		t.Errorf("\n%v\n%v", image1.Pix, image2.Pix)
	}
}

func LogImage(img *image.NRGBA) {
	for i := 0; i < len(img.Pix); i += 4 {
		if img.Pix[i+3] != 0 {
			fmt.Print("#")
		} else {
			fmt.Print(" ")
		}
		if i%(img.Stride) == 0 {
			fmt.Println()
		}
	}
}
