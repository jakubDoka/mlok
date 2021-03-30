package main

import (
	"fmt"
	"math"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/logic/frame"
	"github.com/jakubDoka/mlok/logic/spatial"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/angle"
	"github.com/jakubDoka/mlok/mat/rgba"

	_ "image/png"
)

const (
	RepelCof           = 7.5
	AlignCof           = 0.045
	CohesionCof        = 0.03
	MaxSpeed, MinSpeed = 25.0, 150.0
	Sight              = 10
)

var Scale = mat.V(.5, .5)

func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	t, err := ggl.LoadTexture("fish.png")
	if err != nil {
		panic(err)
	}

	b := ggl.Batch{
		Texture: t,
	}

	// calculating how many tiles MinHash needs to spam whole screen
	size := win.Frame().Size().Divided(Sight).Point()

	e := BoidEngine{
		Sp: ggl.NSprite(t.Frame()),
		Mh: spatial.NMinHash(size.X+1, size.Y+1, mat.V(Sight, Sight)),
	}

	d := frame.Delta{}.Init()

	// We have to shift the camera as hasher has origin in V(0 0) and grows positively.
	// Thats just how hashing is designed.
	win.SetCamera(mat.IM.Move(mat.V(-500, -300))) // everithing about camera is inverted

	e.Spawn(10000, win.Rect()) // rect returns viewport rect in world coordinates

	useless := 0

	for !win.ShouldClose() {
		useless += e.Update(d.Tick(), win.Rect())
		d.CustomLog(1, func() {
			fmt.Println("we made", useless, "iterations in last second")
			useless = 0
		})

		win.Clear(rgba.DarkBlue)
		e.Draw(&b)
		b.Draw(win)
		b.Clear()

		win.Update()
	}
}

// Draw all boids
func (b *BoidEngine) Draw(t ggl.Target) {
	for _, boid := range b.Bs {
		b.Sp.Draw(t, mat.M(boid.Position, Scale, boid.Velocity.Angle()), rgba.White)
	}
}

type Boid struct {
	Position, Velocity mat.Vec
	Adders             mat.Point
}

type BoidEngine struct {
	Bs   []Boid
	Sp   ggl.Sprite
	Buff []int
	Mh   spatial.MinHash
}

func (b *BoidEngine) Spawn(amount int, bounds mat.AABB) {
	// This just showcases how to remove objects
	for i, boid := range b.Bs {
		// Removing of object assumes addres is correct and object with id i and group 0 do exist
		if !b.Mh.Remove(boid.Adders, i, 0) {
			panic("removal of object with incorrect addres or removal of nonexistant object")
		}
	}

	b.Bs = make([]Boid, amount)
	for i := range b.Bs {
		boid := &b.Bs[i]
		boid.Position = mat.V(
			mat.Range(bounds.Min.X, bounds.Max.X),
			mat.Range(bounds.Min.Y, bounds.Max.Y),
		)
		boid.Velocity = mat.Rad(angle.Random(), 100)
		// We are inserting new object to hasher. All that really gets stored is id though.
		// Addres will be modified, last argument is group, in case you need to detect collisions
		// between multiple groups, hasher has optimized solution for this.
		b.Mh.Insert(&boid.Adders, boid.Position, i, 0)
	}
}

// Update velocity and position, it also retruns amount of useles iterations
func (b *BoidEngine) Update(delta float64, bounds mat.AABB) int {
	var useless int
	for i := range b.Bs {
		boid := &b.Bs[i]

		count := 1.0

		// tree rules
		var repel mat.Vec
		cohesion := boid.Position
		alignmant := boid.Velocity

		// Now we are querying for ids that are nearby the boid (all tiles that intersect rectangle)
		b.Buff = b.Mh.Query(mat.Square(boid.Position, Sight), b.Buff[:0], 0, true)
		for _, id := range b.Buff {
			if i == id {
				continue
			}
			other := &b.Bs[id]
			dif := other.Position.To(boid.Position)

			len2 := dif.Len2()

			if math.Sqrt(len2) > Sight {
				useless++
				continue
			}

			count++
			repel.AddE(dif.Divided(len2))
			alignmant.AddE(other.Velocity)
			cohesion.AddE(other.Position)
		}

		cohesion = boid.Position.To(cohesion.Divided(count)).Scaled(CohesionCof)
		alignmant = alignmant.Scaled(AlignCof / count)
		if !bounds.Contains(boid.Position) { // so that boids will come back
			repel.AddE(boid.Position.To(bounds.Center()).Normalized().Scaled(.1))
		}
		repel = repel.Scaled(RepelCof)

		boid.Velocity.AddE(cohesion.Add(alignmant).Add(repel))
		boid.Velocity = mat.Rad(boid.Velocity.Angle(), mat.Clamp(boid.Velocity.Len(), 0, MaxSpeed))

		boid.Position.AddE(boid.Velocity.Scaled(delta))

		// Lastly we are updating the boid addres, this is no-op if addres does not change,
		// group and id has to be provided
		b.Mh.Update(&boid.Adders, boid.Position, i, 0)
	}
	return useless
}

/*import (
	"fmt"
	"math"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/logic/frame"
	"github.com/jakubDoka/mlok/logic/spatial"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/angle"
	"github.com/jakubDoka/mlok/mat/rgba"

	_ "image/png"
)

const (
	RepelCof           = 7.5
	AlignCof           = 0.045
	CohesionCof        = 0.03
	MaxSpeed, MinSpeed = 25.0, 1000.0
	Sight              = 10
)

var Scale = mat.V(.5, .5)

func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	t, err := ggl.LoadTexture("fish.png")
	if err != nil {
		panic(err)
	}

	b := ggl.Batch{
		Texture: t,
	}

	// calculating how many tiles MinHash needs to spam whole screen
	size := win.Frame().Size().Divided(Sight).Point()

	e := BoidEngine{
		Sp: ggl.NSprite(t.Frame()),
		Mh: spatial.NMinHash(size.X+1, size.Y+1, mat.V(Sight, Sight)),
	}

	d := frame.Delta{}.Init()

	// We have to shift the camera as hasher has origin in V(0 0) and grows positively.
	// Thats just how hashing is designed.
	win.SetCamera(mat.IM.Move(mat.V(-500, -300))) // everithing about camera is inverted

	e.Spawn(10000, win.Rect()) // rect returns viewport rect in world coordinates
	fmt.Println(win.Rect())

	useless := 0

	for !win.ShouldClose() {

		useless += e.Update(d.Tick(), win.Rect())
		d.CustomLog(1, func() {
			fmt.Println("we made", useless, "iterations in last second")
			useless = 0
		})

		win.Clear(rgba.DarkBlue)
		e.Draw(&b)
		b.Draw(win)
		b.Clear()

		win.Update()
	}
}

type BoidEngine struct {
	Bs   []Boid
	Sp   ggl.Sprite
	Buff []int
	Mh   spatial.MinHash
}

// Spawn boids on random positions
func (b *BoidEngine) Spawn(amount int, bounds mat.AABB) {
	for i, boid := range b.Bs {
		// Removing of object assumes addres is correct and object with id i and group 0 do exist
		if !b.Mh.Remove(boid.Adders, i, 0) {
			panic("removal of object with incorrect addres or removal of nonexistant object")
		}
	}

	b.Bs = make([]Boid, amount)
	for i := range b.Bs {
		boid := &b.Bs[i]
		boid.Position = mat.V(
			mat.Range(bounds.Min.X, bounds.Max.X),
			mat.Range(bounds.Min.Y, bounds.Max.Y),
		)
		boid.Velocity = mat.Rad(angle.Random(), 100)
		// We are inserting new object to hasher. All that really gets stored is id though.
		// Addres will be modified, last argument is group, in case you need to detect collisions
		// between multiple groups, hasher has optimized solution for this.
		b.Mh.Insert(&boid.Adders, boid.Position, i, 0)
	}
}

// Draw all boids
func (b *BoidEngine) Draw(t ggl.Target) {
	for _, boid := range b.Bs {
		b.Sp.Draw(t, mat.M(boid.Position, Scale, boid.Velocity.Angle()), rgba.White)
	}
}

// Update velocity and position, it also retruns amount of useles iterations
func (b *BoidEngine) Update(delta float64, bounds mat.AABB) int {
	var useless int
	for i := range b.Bs {
		boid := &b.Bs[i]

		// tree rules
		var repel mat.Vec
		count := 1.0
		cohesion := boid.Position
		alignmant := boid.Velocity

		// Now we are querying for ids that are nearby the boid
		b.Buff = b.Mh.Query(mat.Square(boid.Position, Sight), b.Buff[:0], 0, true)
		for _, id := range b.Buff {
			if i == id {
				continue
			}
			other := &b.Bs[id]
			dif := other.Position.To(boid.Position)

			len2 := dif.Len2()

			if math.Sqrt(len2) > Sight {
				useless++
				continue
			}

			count++
			repel.AddE(dif.Divided(len2))
			alignmant.AddE(other.Velocity)
			cohesion.AddE(other.Position)
		}

		cohesion = boid.Position.To(cohesion.Divided(count)).Scaled(CohesionCof)
		alignmant = alignmant.Scaled(AlignCof / count)
		if !bounds.Contains(boid.Position) { // so that boids will come back
			repel.AddE(boid.Position.To(bounds.Center()).Normalized().Scaled(.1))
		}
		repel = repel.Scaled(RepelCof)

		boid.Velocity.AddE(cohesion.Add(alignmant).Add(repel))
		boid.Velocity = mat.Rad(boid.Velocity.Angle(), mat.Clamp(boid.Velocity.Len(), 0, MaxSpeed))

		boid.Position.AddE(boid.Velocity.Scaled(delta))

		// Lastly we are updating the boid addres, this is no-op if addres does not change
		b.Mh.Update(&boid.Adders, boid.Position, i, 0)
	}
	return useless
}

type Boid struct {
	Position, Velocity mat.Vec
	Adders             mat.Point
}
*/
