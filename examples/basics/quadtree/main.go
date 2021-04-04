package main

import (
	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/drw"
	"github.com/jakubDoka/mlok/logic/frame"
	"github.com/jakubDoka/mlok/logic/spatial"
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/rgba"
	"github.com/jakubDoka/mlok/mat/rnd"
)

//go:generate genny -pkg=main -in=$GOPATH\pkg\mod\github.com\jakub!doka\mlok@v0.4.0\logic\memory\storage.go -out=gen-storage.go gen "Element=Entity"

// i also ganges some values to fit quadtree better, if you think that is cheating, try using same values
// for naive approach, results will be even worse
func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}
	// bit more space as even quadtree is ineffective with grouped up entities
	win.SetCamera(mat.IM.Scaled(mat.ZV, .01))

	rect := win.Rect() // viewport bounds

	// Three does not need any initialization, here we are just defining where it will be,
	// how match entities we allow per node and how deep it can be
	tree := spatial.QuadTree{
		Bounds:   rect,
		NodeCap:  20,
		DepthCap: 20,
	}

	const group = 0
	rnd := rnd.Time() // random source seeded on time
	storage := EntityStorage{}
	for i := 0; i < 10000; i++ {
		e, id := storage.Allocate()
		e.pos = rnd.AABB(rect)                           // random point within bounds
		e.vel = mat.Rad(rnd.Angle(), rnd.Range(50, 300)) // velocity in random direction with random length
		e.size = rnd.Range(10, 100)                      // random size
		if i%100 == 0 {
			e.size *= 10 * rnd.Float64()
			e.vel.Scaled(100 * rnd.Float64())
		}
		// Now we are inserting into tree. Notice the address we added to entity. With this very little
		// modification, quadtree can easily locate the entity next time we update it. Next are bounds that
		// will determinate where shape will be inserted. Last two values are identifiers. Tree can manage
		// ids sorted into groups so you can withdraw only ids that you are interested in, and this is not
		// just convenience, Tree is structuring ids in a way that no loop over all ids are needed to get the
		// all ids from one group.
		tree.Insert(&e.address, mat.Square(e.pos, e.size), id, group)
	}

	drawer := drw.Geom{}
	drawer.Restart()
	batch := ggl.Batch{}
	ticker := frame.Delta{}
	ticker.Tick() // get rid of zero value time information witch would make first delta time too long

	// some reusable slices that tree will use to store queries, these are optional, though query will produce
	// lot of allocations othervise. On the other hand these slices cannot be stored in tree ans you may perform
	// querying from multiple threads. Ths makes it safe.
	var buffer, frontier, temp []int
	for !win.ShouldClose() {
		delta := ticker.Tick()
		ticker.Log(1)

		occupied := storage.Occupied()
		for _, id := range occupied {
			e := storage.Item(id)
			area := mat.Square(e.pos, e.size)

			// finding intersection
			intersection := false
			// There we are performing query witch almost newer yields all entities. First we specify filtering
			// (witch group and whether we need it or everithing else), if you need all groups, just specify group
			// that does not exist and pass false. Nest we are passing area with witch ids should intersect, but as
			// they are just ids, intersection check have to be performed by us.
			buffer, frontier, temp = tree.Query(group, true, area, buffer, frontier, temp)
			for _, oid := range buffer {
				if id == oid {
					continue
				}
				o := storage.Item(oid)
				if mat.Square(o.pos, o.size).Intersects(area) {
					intersection = true
					break
				}
			}

			// drawing
			if intersection {
				drawer.Color(rgba.Red)
			} else {
				drawer.Color(rgba.Green)
			}
			drawer.AABB(area)

			// moving
			e.pos.AddE(e.vel.Scaled(delta))
			// Update moves ids between quadrants as needed to keep tree up to date, if entity does not move
			// there is no need to update it though. Though update is mostly noop as entities does not change
			// quadrants every frame, unless they are moving with great speed, with also means it is hard to
			// detect collisions.
			tree.Update(&e.address, area, id, group)

			// keeping on screen
			if e.pos.X > rect.Max.X {
				e.pos.X = rect.Min.X
			}
			if e.pos.X < rect.Min.X {
				e.pos.X = rect.Max.X
			}
			if e.pos.Y > rect.Max.Y {
				e.pos.Y = rect.Min.Y
			}
			if e.pos.Y < rect.Min.Y {
				e.pos.Y = rect.Max.Y
			}
		}

		// final draw
		drawer.Fetch(&batch)
		drawer.Clear()
		batch.Draw(win)
		batch.Clear()

		win.Update()
		win.Clear(rgba.Black)
	}
}

type Entity struct {
	vel, pos mat.Vec
	size     float64
	address  int
}

/*func main() {
	win, err := ggl.NWindow(nil)
	if err != nil {
		panic(err)
	}

	win.SetCamera(mat.IM.Scaled(mat.ZV, .1))

	r := win.Rect()

	tree := spatial.QuadTree{
		Bounds:   r,
		NodeCap:  5,
		DepthCap: 10,
	}

	geom := drw.Geom{}
	geom.Restart()
	batch := ggl.Batch{}

	rnd := rnd.Time()

	obj := make([]Object, 10000)
	for i := range obj {
		obj[i] = Object{
			vel:  rnd.Circ(mat.C(0, 0, 100)),
			size: rnd.Range(5, 40),
			pos:  rnd.AABB(r),
		}
		o := &obj[i]
		tree.Insert(&o.address, mat.Square(o.pos, o.size), i, i)
	}

	ticker := frame.Delta{}.Init()
	var delta float64

	var coll, frontier, temp []int
	for !win.ShouldClose() {
		delta = math.Min(ticker.Tick(), 1.0/30.0)
		ticker.Log(1)
		for i := range tree.Nodes {
			n := &tree.Nodes[i]
			if tree.Nodes[n.Parent].Branch {
				geom.Color(mat.Alpha(.3)).Thickness(1).AABB(n.AABB)
			}
		}

		for i := range obj {
			o := &obj[i]

			o.pos.AddE(o.vel.Scaled(delta))

			if o.pos.X > r.Max.X {
				o.pos.X = r.Min.X
			}
			if o.pos.X < r.Min.X {
				o.pos.X = r.Max.X
			}
			if o.pos.Y > r.Max.Y {
				o.pos.Y = r.Min.Y
			}
			if o.pos.Y < r.Min.Y {
				o.pos.Y = r.Max.Y
			}

			a := mat.Square(o.pos, o.size)
			o.intersecting = false
			coll, frontier, temp = tree.Query(-1, false, a, coll[:0], frontier[:0], temp[:0])
			for _, id := range coll {
				if id == i {
					continue
				}
				oo := &obj[id]
				if a.Intersects(mat.Square(oo.pos, oo.size)) {
					o.intersecting = true
					break
				}
			}
			tree.Update(&o.address, a, i, i)
			if o.intersecting {
				geom.Color(rgba.Red)
			} else {
				geom.Color(rgba.Green)
			}
			geom.AABB(a)
		}

		geom.Fetch(&batch)
		geom.Clear()
		batch.Draw(win)
		batch.Clear()

		win.Update()
		win.Clear(rgba.Black)
	}
}

*/
