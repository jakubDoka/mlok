package ai

import (
	"math"

	"github.com/jakubDoka/mlok/mat"
)

// Predict calculates optimal trajectory for shooter to hit moving target
func Predict(shooter, target, targetVelocity mat.Vec, bulletSpeed float64) (pos mat.Vec, ok bool) {
	/*
		Equation is derivated via analytic math so do not search match obvious logic in there.

		If distance target<>shooter == target<>bullet we assume that we hit the target, all we need is
		right coefficient by witch we multiply targets velocity to get the final position.

		At the end wi just decide what is correct direction as shooting back in time is not an
		option for usual game.
	*/

	d := target.Sub(shooter)

	a := targetVelocity.X*targetVelocity.X + targetVelocity.Y*targetVelocity.Y - bulletSpeed*bulletSpeed
	b := 2 * (d.X*targetVelocity.X + d.Y*targetVelocity.Y)
	c := d.X*d.X + d.Y*d.Y

	// polynomial
	cof := b*b - 4*a*c

	if cof < 0 {
		return
	}

	cof = math.Sqrt(cof)
	a *= 2
	t1, t2 := (-b+cof)/a, (-b-cof)/a

	// deciding witch cof is correct
	if t1 <= 0 || t1 > t2 {
		t1 = t2
	}

	return target.Add(targetVelocity.Scaled(t1)), true
}
