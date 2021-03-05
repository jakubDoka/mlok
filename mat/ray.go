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
	return fmt.Sprintf("R(O(%.3f %.3f) V(%.3f %.3f))", r.O.X, r.O.Y, r.V.X, r.V.Y)
}

// IntersectionPoint calculates the intersection point between two rays
// unless rays are colinear, intersection will be returned, but if false is
// returned, intersection does not include tha both ray segments
func (r Ray) IntersectionPoint(s Ray) (Vec, bool) {
	/*
		first we calculate x of a point and the we project it by non horizontal
		ray, last step is to check if point belongs to both rays
	*/
	x, ok := r.IntersectX(s)
	if !ok {
		return Vec{}, false
	}

	y, ok := r.ProjectX(x)
	if !ok { // can happen if r.V.X == 0
		y, _ = s.ProjectX(x) // other way around works a lines are not colinear at this point
	}

	res := Vec{x, y}

	return res, r.InReach(res) && s.InReach(res)
}

// Contains returns whether ray contains the pos
func (r Ray) Contains(pos Vec) bool {
	return r.Formula(pos) == 0 && r.InReach(pos)
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

// InReach is special method usefull in Contains method, inspect the source code if you want to know
// when it returns true
func (r Ray) InReach(pos Vec) bool {
	/*
		its little tough to explane but whe we check if point belongs to a ray
		and we already know it belongs to a line that ray is on, then if r.V.len()
		is bigger both distances from ray endpoints to given point, the point belongs
		to ray
	*/
	l := r.V.Len2()
	o := r.O.To(pos)
	return l >= o.Len2() && l >= o.Sub(r.V).Len2()
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

// IntersectX returns x coordinate of intersection between two rays
//
// it false if there is no intersection
func (r Ray) IntersectX(s Ray) (float64, bool) {
	/*
		equation is derived from system of equations with
		two unknowns where equations are r.Formula and s.Formula
		from which we can derive x of intersection point
	*/
	c := s.V.Y*r.V.X - r.V.Y*s.V.X
	if c == 0 {
		return 0, false
	}
	a := (r.V.Y*r.O.X - r.V.X*r.O.Y) * s.V.X
	b := s.V.Y*s.O.X*r.V.X - s.V.X*s.O.Y*r.V.Y

	return (b - a) / c, true
}

// IntersectY returns y coordinate of intersection between two rays
//
// returns false if there is no intersection
func (r Ray) IntersectY(s Ray) (float64, bool) {
	/*
		analogous to r.IntersectX
	*/
	c := s.V.X*r.V.Y - r.V.X*s.V.Y
	if c == 0 {
		return 0, false
	}
	a := (r.V.X*r.O.Y - r.V.Y*r.O.X) * s.V.Y
	b := s.V.X*s.O.Y*r.V.X - s.V.Y*s.O.X*r.V.Y

	return (b - a) / c, true
}

// ProjectX returns y coordinate for x coordinate, resulting point is on line
//
// returns false if r is vertical
func (r Ray) ProjectX(x float64) (float64, bool) {
	/*
		derived by evaluating y form r.Formula
	*/
	if r.V.X == 0 {
		return 0, false
	}
	return (r.V.Y*x - r.V.Y*r.O.X + r.V.X*r.O.Y) / r.V.X, true
}

// ProjectY returns x coordinate for y coordinate, resulting point is on line
//
// return false if r is horizontal
func (r Ray) ProjectY(y float64) (float64, bool) {
	/*
		derived by evaluating x form r.Formula
	*/
	if r.V.Y == 0 {
		return 0, false
	}

	return (r.V.X*y - r.V.Y*r.O.X + r.V.X*r.O.Y) / r.V.Y, true
}

// SymmetricPoint finds point of symmetry of two lines that is at a distance form them
func (r Ray) SymmetricPoint(s Ray, distance float64) (Vec, bool) {
	/*
		using the equation for calculation of distance r <-> p

			d = r.Formula(p)/r.V.Len()

		we use equation for s and r to express x of the point that has
		same distance from both lines

		step two is to project x on one of the rays with r.ProjectX
		to get y coordinate of point
	*/
	c := r.V.Y*-s.V.X - s.V.Y*-r.V.X
	if c == 0 { // colinear
		return Vec{}, false
	}

	cof1, cof2 := r.V.X*r.O.Y-r.V.Y*r.O.X, s.V.X*s.O.Y-s.V.Y*s.O.X
	len1, len2 := r.V.Len(), s.V.Len()

	a := distance * (len1*-s.V.X - len2*-r.V.X)
	b := cof2*-r.V.X - cof1*-s.V.X

	x := (a + b) / c
	y := (distance*len1 - r.V.Y*x - cof1) / -r.V.X
	if math.IsNaN(y) {
		y = (distance*len2 - s.V.Y*x - cof2) / -s.V.X
	}

	return Vec{x, y}, true
}
