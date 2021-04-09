package rnd

import (
	"math"
	"math/rand"
	"time"

	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/angle"
)

// Rnd is extencion of rand.Rand, it adds some utility methods relevant for mat package
type Rnd struct {
	*rand.Rand
}

// Time returns Rnd seeded by current time
func Time() Rnd {
	return New(time.Now().Unix())
}

// New creates new Rnd with given seed
func New(seed int64) Rnd {
	return Rnd{Rand: rand.New(rand.NewSource(seed))}
}

// Circ returns random point in given circle
func (r *Rnd) Circ(c mat.Circ) mat.Vec {
	return mat.Rad(r.Angle(), c.R).Add(c.C)
}

// Ray returns random point on ray
func (r *Rnd) Ray(ray mat.Ray) mat.Vec {
	return ray.O.Add(ray.V.Scaled(r.Float64()))
}

// AABB returns random point contained in aabb
func (r *Rnd) AABB(a mat.AABB) mat.Vec {
	return mat.V(
		r.Range(a.Min.X, a.Max.X),
		r.Range(a.Min.Y, a.Max.Y),
	)
}

// Range returns point in given range
func (r *Rnd) Range(min, max float64) float64 {
	return min + (max-min)*r.Float64()
}

// angle returns random angle in range <-Pi, Pi>
func (r *Rnd) Angle() float64 {
	return angle.Pi2*r.Float64() - math.Pi
}
