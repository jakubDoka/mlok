package main

import (
	"flag"
	"fmt"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/drw"
	"github.com/jakubDoka/mlok/ggl/key"
	"github.com/jakubDoka/mlok/logic/frame"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/angle"
	"github.com/jakubDoka/mlok/mat/rgba"
)

var (
	wallCount   = flag.Int("walls", 3, "amount of walls used")
	circleCount = flag.Int("circles", 3, "amount of circles used")
	rayCount    = flag.Int("rays", 100, "amount of rays used")
)

// example demonstrates raycasting capability of mat package
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

	// generating random circles
	circles := make([]mat.Circ, *circleCount)
	for i := range circles {
		circles[i] = mat.Circ{
			C: mat.V(mat.Range(f.Min.X, f.Max.X), mat.Range(f.Min.Y, f.Max.Y)),
			R: mat.Range(0, 50),
		}
	}

	// generating rays, this can take several seconds if there is more then 1000 rays.
	// Not that math is slow but allocating memory, *rayCount * 4 * 8 bytes in one swoop
	rays := make([]mat.Ray, *rayCount)
	cof := angle.Pi2 / float64(*rayCount)
	angle := .0
	for i := range rays {
		rays[i] = mat.Ray{
			O: mat.Rad(angle, 10),
			V: mat.Rad(angle, 4000),
		}
		// too little numbers can create visual artifacts
		rays[i].V.X = mat.Round(rays[i].V.X, 4)
		rays[i].V.Y = mat.Round(rays[i].V.Y, 4)
		rays[i].O.X = mat.Round(rays[i].O.X, 4)
		rays[i].O.Y = mat.Round(rays[i].O.Y, 4)
		angle += cof
	}

	fmt.Println(rays)

	batch := ggl.Batch{}

	// Geom draws geometric shapes witch is ideal for us
	geom := drw.Geom{}
	geom.Restart()
	geom.Thickness(1)
	geom.Fill(false)

	// maybe poor design decision but circle intersections are buffered
	// as there can be 0 to 2 intersections
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

		// dragging rays
		if win.Pressed(key.MouseLeft) {
			d := win.MousePrevPos().To(win.MousePos())
			for i := range rays {
				rays[i].O.AddE(d)
			}
		}

		// drawing it all
		win.Clear(rgba.Black)
		geom.Fetch(&batch)
		geom.Clear()
		batch.Draw(win)
		batch.Clear()
		win.Update()
	}
}
