package mat

import (
	"fmt"
	"math"
)

// Ray is a standard  raycast, it supports calculating intersections with all collizion shapes
type Ray struct {
	// O is an origin and V is a directional vector
	O, V Vec
}

// ZR is zero value ray
var ZR Ray

// R creates new raycast, o stands for origin and v for directional vector
func R(ox, oy, vx, vy float64) Ray {
	return Ray{Vec{ox, oy}, Vec{vx, vy}}
}

func (r Ray) String() string {
	return fmt.Sprintf("R(%v %v %v %v)", ff(r.O.X), ff(r.O.Y), ff(r.V.X), ff(r.V.Y))
}

// Contains returns whether ray contains the pos
func (r Ray) Contains(pos Vec) bool {
	return r.Formula(pos) == 0 && r.InAABB(pos)
}

// IntersectCircle returns points of intersection between Ray and circle and whether
// they are valid
func (r Ray) IntersectCircle(c Circ, buff []Vec) []Vec {
	if !r.LineIntersectsCircle(c) {
		return buff
	}

	a, b := r.LineIntersectCircle(c)

	if r.InAABB(a) {
		buff = append(buff, a)
	}
	if r.InAABB(b) {
		buff = append(buff, b)
	}

	return buff
}

// LineIntersectCircle calculates intersection points between circle and line
//
// this does not make sense for line and circle that does not intersect
func (r Ray) LineIntersectCircle(c Circ) (g, h Vec) {
	/*
		(X - c.C.X)^2 + (Y - c.C.Y)^2 = c.R*c.R
		r.V.Y*X - r.V.X*Y - r.V.Y*r.O.X + r.V.X*r.O.Y = 0

		X*X - 2*c.C.X*X + c.C.X*c.C.X + Y*Y - 2*c.C.Y*Y + c.C.Y*c.C.Y = c.R*c.R
		(r.V.X*Y + r.V.Y*r.O.X - r.V.X*r.O.Y) / r.V.Y = X

		a := c.C.X*c.C.X + c.C.Y*c.C.Y - c.R*c.R
		b := r.V.Y*r.O.X - r.V.X*r.O.Y

		X*X - 2*c.C.X*X + Y*Y - 2*c.C.Y*Y + a = 0
		(r.V.X*Y + b) / r.V.Y = X

		(r.V.X*Y + b)^2 / r.V.Y*r.V.Y - 2*c.C.X * (r.V.X*Y + b) / r.V.Y + Y*Y - 2*c.C.Y*Y + a = 0
		// * r.V.Y*r.V.Y
		(r.V.X*Y + b)^2 - 2*c.C.X*r.V.Y * (r.V.X*Y + b) + Y*Y*r.V.Y*r.V.Y - 2*c.C.Y*Y*r.V.Y*r.V.Y + a*r.V.Y*r.V.Y = 0
		// brackets
		r.V.X*r.V.X*Y*Y + 2*b*r.V.X*Y + b*b - 2*c.C.X*r.V.Y*r.V.X*Y - 2*c.C.X*r.V.Y*b + Y*Y*r.V.Y*r.V.Y - 2*c.C.Y*Y*r.V.Y*r.V.Y + a*r.V.Y*r.V.Y = 0
		// simplify
		Y*Y * (r.V.X*r.V.X + r.V.Y*r.V.Y) + Y * (2*b*r.V.X - 2*c.C.X*r.V.Y*r.V.X - 2*c.C.Y*r.V.Y*r.V.Y) + b*b - 2*c.C.X*r.V.Y*b  + a*r.V.Y*r.V.Y = 0

		d := r.V.X*r.V.X + r.V.Y*r.V.Y
		e := 2*b*r.V.X - 2*c.C.X*r.V.Y*r.V.X - 2*c.C.Y*r.V.Y*r.V.Y
		f := b*b - 2*c.C.X*r.V.Y*b  + a*r.V.Y*r.V.Y

		e := 2 * (b*r.V.X - r.V.Y * (c.C.X*r.V.X + c.C.Y*r.V.Y))
	*/

	a := c.C.X*c.C.X + c.C.Y*c.C.Y - c.R*c.R
	b := r.V.Y*r.O.X - r.V.X*r.O.Y
	d := r.V.X*r.V.X + r.V.Y*r.V.Y
	e := 2 * (b*r.V.X - r.V.Y*(c.C.X*r.V.X+c.C.Y*r.V.Y))
	f := b*b - r.V.Y*(2*c.C.X*b-a*r.V.Y)

	g.Y, h.Y = Polynomial(d, e, f)
	if r.V.Y == 0 {
		g.X, h.X = c.ProjectY(g.Y)
	} else {
		g.X = r.ProjectY(g.Y)
		h.X = r.ProjectY(h.Y)
	}

	return
}

// IntersectsCircle returns whether ray intersects circle
func (r Ray) IntersectsCircle(c Circ) bool {
	return r.LineIntersectsCircle(c) && r.InAABB(c.C)
}

// LineIntersectsCircle returns whether line and circle intersects
func (r Ray) LineIntersectsCircle(c Circ) bool {
	return math.Abs(r.Formula(c.C))/r.V.Len() <= c.R
}

// Formula returns 0 if point belongs to line that ray is on and negative or positive
// number depending on witch half plane point is in comparison to line
func (r Ray) Formula(pos Vec) float64 {
	/*
		formula is derivated from general line equation

			a*x + b*y + c = 0

		where a is equal to r.V.Y and b is equal to - r.V.X
		as V(a, b) is normal to r.V

		to calculate c we have to put r.O into equation

			c = - r.V.Y*r.O.X + r.V.X*r.O.Y
	*/
	return r.V.Y*pos.X - r.V.X*pos.Y - r.V.Y*r.O.X + r.V.X*r.O.Y
}

// Intersect calculates the intersection point between two rays
// unless rays are colinear, intersection will be returned, but if false is
// returned, intersection does not include tha both ray segments
func (r Ray) Intersect(s Ray) (v Vec, ok bool) {
	if r.Colinear(s) {
		return
	}

	v = r.LineIntersect(s)

	return v, r.InAABB(v) && s.InAABB(v)
}

// LineIntersect returns the point of intersection of two lines expressed by rays
//
// function does not make sense for colinear lines
func (r Ray) LineIntersect(s Ray) (point Vec) {
	/*
		equation is derived from system of equations with
		two unknowns where equations are r.Formula and s.Formula
		from which we can derive x of intersection point

		starting with:
			r.V.Y*X - r.V.X*Y - r.V.Y*r.O.X + r.V.X*r.O.Y = 0
		and:
			s.V.Y*X - s.V.X*Y - s.V.Y*s.O.X + s.V.X*s.O.Y = 0

		get y from first one:
			r.V.Y*X - r.V.Y*r.O.X + r.V.X*r.O.Y = r.V.X*Y
			(r.V.Y*X - r.V.Y*r.O.X + r.V.X*r.O.Y)/r.V.X = Y

		then we substitute and get x:
			s.V.Y*X - s.V.X * (r.V.Y*X - r.V.Y*r.O.X + r.V.X*r.O.Y) / r.V.X - s.V.Y*s.O.X + s.V.X*s.O.Y = 0 // * r.V.X
			s.V.Y*X*r.V.X - s.V.X*r.V.Y*X + s.V.X*r.V.Y*r.O.X - s.V.X*r.V.X*r.O.Y - s.V.Y*s.O.X*r.V.X + s.V.X*s.O.Y*r.V.X = 0 // - s.V.Y*X*r.V.X + s.V.X*r.V.Y*X
			s.V.X*r.V.Y*r.O.X - s.V.X*r.V.X*r.O.Y - s.V.Y*s.O.X*r.V.X + s.V.X*s.O.Y*r.V.X = s.V.X*r.V.Y*X - s.V.Y*X*r.V.X // simplify
			s.V.X * (r.V.Y*r.O.X + r.V.X * (s.O.Y - r.O.Y)) - s.V.Y*s.O.X*r.V.X = X * (s.V.X*r.V.Y - s.V.Y*r.V.X) // / (s.V.X*r.V.Y - s.V.Y*r.V.X)
			(s.V.X * (r.V.Y*r.O.X + r.V.X * (s.O.Y - r.O.Y)) - s.V.Y*s.O.X*r.V.X) / (s.V.X*r.V.Y - s.V.Y*r.V.X) = X
	*/

	point.X = (s.V.X*(r.V.Y*r.O.X+r.V.X*(s.O.Y-r.O.Y)) - s.V.Y*s.O.X*r.V.X) / (s.V.X*r.V.Y - s.V.Y*r.V.X)

	if r.V.X == 0 {
		point.Y = s.ProjectX(point.X)
	} else {
		point.Y = r.ProjectX(point.X)
	}

	return
}

// ProjectX returns y coordinate for x coordinate, resulting point is on line
//
// this does not make sense for vertical ray
func (r Ray) ProjectX(x float64) float64 {
	/*
		derived by evaluating y form r.Formula
	*/
	return (r.V.Y*x - r.V.Y*r.O.X + r.V.X*r.O.Y) / r.V.X
}

// ProjectY returns x coordinate for y coordinate, resulting point is on line
//
// this does not make sense for horizontal ray
func (r Ray) ProjectY(y float64) float64 {
	/*
		derived by evaluating x form r.Formula
	*/
	return (r.V.X*y + r.V.Y*r.O.X - r.V.X*r.O.Y) / r.V.Y
}

// Colinear returns whether two rays are colinear
func (r Ray) Colinear(s Ray) bool {
	/*
		if two vectors are colinear then cross product has to return 0
		though we skip subtraction by comparing to sides of it as if
		a == b then a - b == 0
	*/
	return r.V.X*s.V.Y == r.V.Y*s.V.X // division is slower
}

// InAABB returns whether pos is inside AABB expressed by ray
func (r Ray) InAABB(pos Vec) bool {
	p := r.O.To(pos)
	if r.V.Y == 0 {
		p.Y = 0
	} else {
		p.Y /= r.V.Y
	}
	if r.V.X == 0 {
		p.X = 0
	} else {
		p.X /= r.V.X
	}
	return p.X <= 1 && p.Y <= 1 && p.X >= 0 && p.Y >= 0
}
