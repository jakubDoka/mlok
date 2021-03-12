package main

import (
	"gobatch/ggl"
	"gobatch/ggl/drw"
	"gobatch/ggl/drw/particle"
	"gobatch/logic/frame"
	"gobatch/logic/gate"
	"gobatch/mat"
	"gobatch/mat/lerp"
	"gobatch/mat/lerpc"
	"gobatch/mat/rgba"
	"math"
	"runtime"
)

func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	// almost all parameters are interfaces, which forces you to specify all value
	// but gives great deal of customization
	tp := particle.Type{
		Scale:             lerp.Const(1),
		ScaleMultiplier:   lerp.Bezier(0, 4, 0, 0), // particle will grow and then shrink
		TwerkAcceleration: lerp.Const(0),
		Acceleration:      lerp.Const(0),
		Twerk:             lerp.Const(0),
		Livetime:          lerp.Const(1), // constant one second livetime
		Rotation:          lerp.Const(0),
		Velocity:          lerp.Random(800, 1600),
		Spread:            lerp.Random(0, math.Pi*.7), // random number in range between 0 - math.Pi*.7

		Color: lerpc.Linear(mat.Alpha(1), mat.Alpha(0)), // fading out
		//random color in folloving range, each channel is independently randomized
		Mask: lerpc.Random(mat.Black, mat.White),

		EmissionShape: particle.Point{},
		Gravity:       mat.V(0, -500),
		Friction:      3,
	}

	// ew need somethong to draw the partiles, circle is good enough
	tp.SetDrawer(&particle.Circle{Circle: drw.NCircle(10, 20)})

	batch := ggl.Batch{}

	ticker := frame.Delta{}
	var delta float64

	// particle system as many threads as we have cores
	threadCount := runtime.NumCPU()
	system := particle.System{}
	system.SetThreads(threadCount)

	// setting up gate to run the particle system on multiple threads
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
		// under the hud every thread always proceses particles like follows
		//	for i := threadIndex; i < particleCount; i += threadCount {
		gt.Run()
		gt.Wait()

		// spawn requests has to be runned on single thread though you can choose to
		// run spawning on different thread
		system.Spawn()

		// drawing system onto batch
		system.Fetch(&batch)

		win.Clear(rgba.Black)
		// batch now draws to window which does the draw call
		batch.Draw(win)
		// don't forget to clear batch or you will run out of memory
		batch.Clear()
		// also important or you will end up with frozen window
		win.Update()
	}
}
