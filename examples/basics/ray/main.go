package main

import (
	"flag"

	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/ggl/drw"
	"github.com/jakubDoka/gobatch/ggl/key"
	"github.com/jakubDoka/gobatch/logic/frame"
	"github.com/jakubDoka/gobatch/mat"
	"github.com/jakubDoka/gobatch/mat/angle"
	"github.com/jakubDoka/gobatch/mat/rgba"
)

var (
	wallCount   = flag.Int("walls", 3, "amount of walls used")
	circleCount = flag.Int("circles", 3, "amount of circles used")
	rayCount    = flag.Int("rays", 100, "amount of rays used")
)

func main() {
	flag.Parse()

	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	// generating random walls
	f := win.Frame()
	walls := make([]mat.Ray, *wallCount)
	for i := range walls {
		walls[i] = mat.Ray{
			O: mat.V(mat.Range(f.Min.X, f.Max.X), mat.Range(f.Min.Y, f.Max.Y)),
			V: mat.Rad(angle.Random(), mat.Range(100, 300)),
		}
	}

	// generating circles
	circles := make([]mat.Circ, *circleCount)
	for i := range circles {
		circles[i] = mat.Circ{
			C: mat.V(mat.Range(f.Min.X, f.Max.X), mat.Range(f.Min.Y, f.Max.Y)),
			R: mat.Range(0, 50),
		}
	}

	// generating rays, this can take several seconds if there is more then 1000 rays
	// not that math is slow but allocating memory
	rays := make([]mat.Ray, *rayCount)
	cof := angle.Pi2 / float64(*rayCount)
	angle := .0
	for i := range rays {
		rays[i] = mat.Ray{
			O: mat.Rad(angle, 10),
			V: mat.Rad(angle, 4000),
		}
		angle += cof
	}

	batch := ggl.Batch{}

	// Geom draws geometric shapes witch is ideal for us
	geom := drw.Geom{}
	geom.Restart()
	geom.Width(1)
	geom.Resolution(100)
	geom.Fill(false)

	// maybe poor design decision but circle intersections are buffered
	buff := make([]mat.Vec, 0, 2)
	delta := frame.Delta{}

	for !win.ShouldClose() {
		// log framerate
		delta.Tick()
		delta.Log(2)

		// calculating intersections and drawing rays
		geom.Color(rgba.Gray)
		for _, r := range rays {
			p := r.O.Add(r.V)
			l := r.V.Len2()

			for _, w := range walls {
				pt, ok := w.Intersect(r)
				ln := pt.To(r.O).Len2()
				if ok && ln < l {
					l = ln
					p = pt
				}
			}

			for _, c := range circles {
				buff := r.IntersectCircle(c, buff)
				for _, v := range buff {
					ln := v.To(r.O).Len2()
					if ln < l {
						l = ln
						p = v
					}
				}
			}

			geom.Line(r.O, p)
		}

		// drawing walls and circles
		geom.Color(rgba.White)
		for _, w := range walls {
			geom.Line(w.O, w.O.Add(w.V))
		}
		for _, c := range circles {
			geom.Circle(c)
		}

		// drag for rays
		if win.Pressed(key.MouseLeft) {
			d := win.MousePrevPos().To(win.MousePos())
			for i := range rays {
				rays[i].O.AddE(d)
			}
		}

		win.Clear(rgba.Black)

		// geom can be drawn as sprite, though don't forget to clear
		geom.Fetch(&batch)
		geom.Clear()

		batch.Draw(win)
		batch.Clear()

		// also important or you will end up with frozen window
		win.Update()
	}
}
