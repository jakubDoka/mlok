package main

import (
	"flag"
	"gobatch/ggl"
	"gobatch/ggl/drw/particles"
	"gobatch/logic/frame"
	"gobatch/logic/gate"
	"gobatch/mat"
	"gobatch/mat/lerp"
	"gobatch/mat/rgba"
	"log"
	"math"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write heap profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	tp := particles.Type{
		Scale:             lerp.Const(1),
		ScaleMultiplier:   lerp.L(1, 0),
		TwerkAcceleration: lerp.Const(0),
		Acceleration:      lerp.Const(0),
		Twerk:             lerp.Const(0),
		Livetime:          lerp.Const(1),
		Rotation:          lerp.Const(0),
		Velocity:          lerp.Const(400),
		Spread:            lerp.R(0, math.Pi),
		EmissionShape:     particles.Point{},
		Color:             lerp.ConstColor(mat.White),
	}

	sprite := ggl.NSprite(mat.A(0, 0, 10, 10))
	sprite.SetIntensity(0)
	tp.SetDrawer(&particles.Sprite{Sprite: sprite})

	batch := ggl.Batch{}

	ticker := frame.Delta{}
	limitter := frame.Limitter{}

	//limitter.SetFPS(60)

	var delta float64

	system := particles.System{}
	system.SetThreads(4)

	gt := gate.Gate{}
	for i := 0; i < 4; i++ {
		thr := system.Thread(i)
		gt.Add(&gate.FazerBase{
			Fazes: []gate.FazeRunner{
				func(tIdx, count int) {
					thr.Update(delta)
					thr.Request(particles.Request{
						Amount: 100,
						Mask:   mat.White,
						Type:   &tp,
					})
				},
			},
		})
	}

	for !win.ShouldClose() {
		delta = ticker.Tick()
		ticker.Log(2)
		limitter.Regulate()

		gt.Run()
		gt.Wait()

		system.Spawn()

		batch.Data = system.Data

		win.Clear(rgba.Black)
		// batch now draws to window which does the draw call
		batch.Draw(win)
		// don't forget to clear batch or you will run out of memory
		batch.Clear()
		// also important or you will end up with frozen window
		win.Update()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}
