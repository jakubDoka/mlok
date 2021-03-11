package particles

import (
	"gobatch/ggl"
	"gobatch/mat"
)

/*imp(
	github.com/jakubDoka/gogen/templates
)*/

/*gen(
	templates.Resize<particles, Resize>
)*/

// System is a particle renderer and updater.
// Its designed to be updated from multiple threads,
// and also to spawn particles off thread. System has its own goroutine
// that is reused for spawning so Drop is necessary to prevent a memory leak
type System struct {
	threadCount, vertex, indice int
	intensity                   float64
	dropped, spawning           bool

	particles particles

	threads []Thread

	ggl.Data

	spawner struct {
		in  chan bool
		out chan struct{}
	}
}

// RunSpawner spawns all spawn requests on separate thread
func (s *System) RunSpawner() {
	if s.spawning {
		panic("already spawning, call Wait first")
	}
	s.spawning = true
	s.spawner.in <- true
}

// Wait waits for spawner to finish
func (s *System) Wait() {
	if !s.spawning {
		panic("spawner is asleep, call RunSpawner to wake him up")
	}
	<-s.spawner.out
	s.spawning = false
}

// Spawn spawns all particles on current thread
func (s *System) Spawn() {
	s.clear()
	s.spawn()
	s.allocate()
}

// Drop drops the particle system (spawner disposing)
func (s *System) Drop() {
	if s.spawner.in != nil {
		s.spawner.in <- false
		s.spawner.in = nil
	}
}

// Thread returns handle to one of used threads, you have to call update, but you can call it
// form different thread
func (s *System) Thread(threadIndex int) *Thread {
	return &s.threads[threadIndex]
}

// SetThreads sets the thread count the System is targetting
func (s *System) SetThreads(count int) {
	s.threadCount = count
	s.threads = make([]Thread, count)
	for i := range s.threads {
		s.threads[i].idx = i
		s.threads[i].System = s
	}

	s.setupSpawner()
}

func (s *System) setupSpawner() {
	if s.spawner.in != nil {
		s.spawner.in <- false // kill current spawner
	} else {
		s.spawner.in = make(chan bool)
		s.spawner.out = make(chan struct{})
	}

	go func() {
		for <-s.spawner.in {
			s.clear()
			s.spawn()
			s.allocate()
			s.spawner.out <- struct{}{}
		}
	}()
}

func (s *System) spawn() {
	for i := range s.threads {
		t := &s.threads[i]
		for _, sr := range t.requests {
			str := len(s.particles)
			s.particles.Resize(str + sr.Amount)
			for i := 0; i < sr.Amount; i++ {
				r := sr.Rotation.Float(0)
				vel := mat.Rad(sr.Spread.Float(0)+sr.Dir, sr.Velocity.Float(0))
				if sr.RotationRelativeToVelocity {
					r += vel.Angle()
				}

				p := &s.particles[i+str]

				p.Type = sr.Type

				p.vel = vel
				p.orig = sr.Pos
				p.pos = sr.Pos.Add(sr.Gen(sr.Dir))

				p.mask = sr.Mask

				p.scl = sr.Scale.Float(0)
				p.livetime = sr.Livetime.Float(0)
				p.twerk = sr.Twerk.Float(0)
				p.rot = r
				p.progress = 0

				p.vertex = s.vertex
				p.indice = s.indice

				s.vertex += p.vertexes
				s.indice += p.indices
			}
		}
		t.requests = t.requests[:0]
	}
}

func (s *System) clear() {
	s.vertex = 0
	s.indice = 0

	var i int
	for j := range s.particles {
		p := &s.particles[j]
		if p.progress >= 1 {
			continue
		}

		p.vertex = s.vertex
		p.indice = s.indice

		s.vertex += p.vertexes
		s.indice += p.indices

		s.particles[i] = *p
		i++
	}

	s.particles = s.particles[:i]

	return
}

func (s *System) allocate() {
	s.Vertexes.Resize(s.vertex)
	s.Indices.Resize(s.indice)
}

//  ...
type particle struct {
	*Type

	vertex, indice int

	mask mat.RGBA

	pos, vel, orig mat.Vec

	scl, rot, twerk, progress, livetime float64
}

func (p *particle) update(delta float64) {
	p.vel.AddE(p.vel.Scaled(p.Acceleration.Float(p.progress) * delta))
	p.vel.AddE(p.pos.To(p.orig).Normalized().Scaled(p.OriginGravity * delta))
	p.vel.SubE(p.vel.Scaled(p.Friction * delta))

	p.pos.AddE(p.vel.Scaled(delta))

	p.twerk += p.TwerkAcceleration.Float(p.progress) * delta
	p.twerk -= p.twerk * p.Friction * delta

	p.rot += p.twerk * delta

	p.progress += delta / p.livetime
}

// Thread ...
type Thread struct {
	idx int
	*System
	requests []Request
}

// Request requests particle spawn, particles should be spawned within frame of this call
func (t *Thread) Request(r Request) {
	if t.spawning {
		panic("cannot request particles when System is spawning particles, this is sync issue on your side")
	}
	if t.threadCount > r.threadCount {
		r.setThreads(t.threadCount)
	}
	t.requests = append(t.requests, r)
}

// Update updates particle state of current thread
func (t *Thread) Update(delta float64) {
	if t.spawning {
		panic("cannot update when System is spawning particles, this is sync issue on your side")
	}
	for i := t.idx; i < len(t.particles); i += t.threadCount {
		p := &t.particles[i]

		p.update(delta)

		// we just overwrite data
		data := ggl.Data{
			Indices:  t.Indices[:p.indice],
			Vertexes: t.Vertexes[:p.vertex],
		}

		scl := p.scl * p.ScaleMultiplier.Float(p.progress)

		p.dws[t.idx].Draw(
			&data,
			mat.M(p.pos, mat.V(scl, scl), p.rot),
			p.Color.Color(p.progress).Mul(p.mask),
		)
	}
}

// Request holds data for particle spawning
type Request struct {
	Pos    mat.Vec
	Dir    float64
	Amount int
	Mask   mat.RGBA
	*Type
}

type particles []particle
