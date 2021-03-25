package angle

import (
	"math"
	"math/rand"
)

const (
	Pi2  = math.Pi * 2
	IPi2 = 1 / Pi2
)

// Random returns random angle in range <-math.Pi, math.Pi>
func Random() float64 {
	return Pi2*rand.Float64() - math.Pi
}

// Turn turns the angle from start to dest by velocity, always the shorter way
//
// assuming angles are ranging from <-math.Pi, math.Pi>
func Turn(start, dest, vel float64) float64 {
	d := dest - start
	if math.Abs(d) < vel {
		return dest
	}

	if d > 0 {
		if d > math.Pi {
			start -= vel
		} else {
			start += vel
		}
	} else {
		if d < -math.Pi {
			start += vel
		} else {
			start -= vel
		}
	}

	return start
}

// To returns shortest step form a to b
//
// assuming angles are ranging from <-math.Pi, math.Pi>
func To(a, b float64) float64 {
	d := b - a
	if d > math.Pi {
		d -= Pi2
	} else if d < -math.Pi {
		d += Pi2
	}

	return d
}

// Norm normalizes angle between <-math.Pi, math.Pi>
func Norm(a float64) float64 {
	if a > math.Pi {
		a -= math.Floor(a*IPi2)*Pi2 + Pi2
	} else if a < -math.Pi {
		a -= math.Ceil(a*IPi2)*Pi2 - Pi2
	}
	return a
}
