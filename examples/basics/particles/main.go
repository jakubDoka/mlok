package main

import (
	"math"
	"runtime"

	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/ggl/drw"
	"github.com/jakubDoka/gobatch/ggl/drw/particle"
	"github.com/jakubDoka/gobatch/logic/frame"
	"github.com/jakubDoka/gobatch/logic/gate"
	"github.com/jakubDoka/gobatch/mat"
	"github.com/jakubDoka/gobatch/mat/lerp"
	"github.com/jakubDoka/gobatch/mat/lerpc"
	"github.com/jakubDoka/gobatch/mat/rgba"
)

// example demonstrates particle drawing and multithreading capability of ggl/drw/particle
// and logic/gate
func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	// almost all parameters are interfaces, which forces you to specify all value
	// but gives great deal of customization
	tp := particle.Type{
		ScaleX:            lerp.Const(1),
		ScaleY:            lerp.Const(2),
		ScaleMultiplier:   lerp.Bezier(0, 4, 0, 0), // particle will grow and then shrink
		TwerkAcceleration: lerp.Const(0),
		Acceleration:      lerp.Const(0),
		Twerk:             lerp.Random(-20, 20),
		Livetime:          lerp.Const(1), // constant one second livetime
		Rotation:          lerp.Const(0),
		Velocity:          lerp.Random(800, 1600),
		Spread:            lerp.Random(0, math.Pi*.7), // random number in range between 0 - math.Pi*.7

		Color: lerpc.Linear(mat.Alpha(0), mat.Alpha(1)), // fading out
		//random color in folloving range, each channel is independently randomized
		Mask: lerpc.Random(mat.Black, mat.White),

		EmissionShape: particle.Point{},
		Gravity:       mat.V(0, -500),
		Friction:      2,
	}

	// ew need something to draw the partiles, circle is good enough
	// and to make it more exciting it will be 2D sphere
	tp.SetDrawer(&particle.Circle{Circle: drw.NCircle(10, 3, 20)})

	batch := ggl.Batch{}

	ticker := frame.Delta{}
	var delta float64

	// we can run particle system at as many threads as we have cores
	threadCount := runtime.NumCPU()
	system := particle.System{}
	system.SetThreads(threadCount)

	// setting up gate to run the particle system on multiple threads
	// gate will run one faze on each thread simultaneously
	gt := gate.Gate{}
	for i := 0; i < threadCount; i++ {
		thr := system.Thread(i)
		gt.Add(&gate.FazerBase{
			Fazes: []gate.FazeRunner{
				func(tIdx, count int) {
					thr.Update(delta)
					thr.Request(particle.Request{
						Amount: 1,
						Pos:    mat.V(-100, -300),
						Mask:   mat.White,
						Type:   &tp,
						Dir:    math.Pi / 8,
					})
				},
			},
		})
	}

	for !win.ShouldClose() {
		delta = ticker.Tick()
		ticker.Log(2) // logging a frame rate, why not

		// this will run all threads at the same time
		// under the hud every thread always proceses particles like:
		//	for i := threadIndex; i < particleCount; i += threadCount {
		gt.Run()
		// you can of corse do something in between but all threads are
		// already used anyway
		gt.Wait()

		// RunSpawner awakens separate thread for spawning new particles
		// particle spawning and drawing does not share any state so we can
		// safely do this, you can optionally use system.Spawn() if you don't
		// need it on separate thread
		system.RunSpawner()

		win.Clear(rgba.Black)

		system.Fetch(&batch) // no need to clear system, you actually should never do that
		batch.Draw(win)
		batch.Clear()
		win.Update()

		// waiting for spawner to finish
		system.Wait()
	}

	// system has its own worker thread that has to be terminated
	system.Drop()
}
