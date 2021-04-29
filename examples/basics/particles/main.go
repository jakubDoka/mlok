package main

import (
	"math"
	"runtime"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/drw/particle"
	"github.com/jakubDoka/mlok/logic/frame"
	"github.com/jakubDoka/mlok/logic/gate"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/lerp"
	"github.com/jakubDoka/mlok/mat/lerpc"
	"github.com/jakubDoka/mlok/mat/rgba"
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
		ScaleX:            lerp.Const(2),
		ScaleY:            lerp.Const(2),
		ScaleMultiplier:   lerp.Bezier(0, 4, 0, 0), // particle will grow and then shrink
		TwerkAcceleration: lerp.Const(0),
		Acceleration:      lerp.Const(0),
		Twerk:             lerp.Random(-20, 20),
		Livetime:          lerp.Const(2), // constant one second livetime
		Rotation:          lerp.Const(0),
		Velocity:          lerp.Random(400, 800),
		Spread:            lerp.Random(0, math.Pi*.7), // random number in range between 0 - math.Pi*.7

		Color: lerpc.Chained(lerpc.Point(0, mat.Alpha(0)), lerpc.Point(.2, mat.Alpha(.2)), lerpc.Point(1, mat.Alpha(0))), // fading out
		//random color in folloving range, each channel is independently randomized
		Mask: lerpc.Random(rgba.OrangePeel, rgba.OrangeRed),

		EmissionShape: particle.Point{},
		Gravity:       mat.V(0, -250),
		Friction:      1,
	}

	drawer := particle.Sprite{Sprite: ggl.NSprite(mat.A(0, 0, 10, 10))}
	drawer.SetIntensity(0)
	drawer.SetPivot(mat.V(5, 5))
	// ew need something to draw the partiles, circle is good enough
	// and to make it more exciting it will be 2D sphere
	tp.SetDrawer(&drawer)

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
						Amount: 10,
						Pos:    mat.V(-100, -200),
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
