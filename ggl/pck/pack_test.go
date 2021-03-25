package pck

import (
	"image"
	"image/draw"
	"reflect"
	"testing"

	"github.com/jakubDoka/mlok/ggl"

	"github.com/jakubDoka/gogen/dirs"
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
	names, err := dirs.ListFilePaths("C:/Users/jakub/Documents/programming/golang/src/mlok/t1", ".png")
	if err != nil {
		panic(err)
	}

	sheet, err := NSheet("t1", names...)
	if err != nil {
		panic(err)
	}

	image1, err := ggl.LoadImage("C:/Users/jakub/Documents/programming/golang/src/mlok/t1/beckup.png")
	if err != nil {
		panic(err)
	}

	r := sheet.Regions["beckup"].ToImage()

	image2 := image.NewNRGBA(r)
	draw.Draw(image2, image2.Bounds(), sheet.Pic, r.Min, 0)
	ggl.FlipNRGBA(image2)

	if !reflect.DeepEqual(image1.Pix, image2.Pix) {
		t.Errorf("\n%v\n%v", image1.Pix, image2.Pix)
	}
}
