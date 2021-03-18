package mat

import (
	"fmt"
	"math"
)

// Circ is an circle
type Circ struct {
	C Vec
	R float64
}

// C is circle constructor
func C(x, y, r float64) Circ {
	return Circ{Vec{x, y}, r}
}

func (c Circ) String() string {
	return fmt.Sprintf("C(%v %v %v)", ff(c.C.X), ff(c.C.Y), ff(c.R))
}

// Union returns smallest possible circle that encloses both circles
func (c Circ) Union(o Circ) Circ {
	d := c.C.To(o.C)
	l := d.Len()

	c.C.AddE(Rad(d.Angle(), l*.5+o.R-c.R))
	c.R += l + o.R

	return c
}

func (c Circ) Intersect(s Circ) (k, p Vec) {
	/*
		a*a = ca*ca + h*h
		b*b = cb*cb + h*h
		c = ca + cb

		a*a - ca*ca = b*b - cb*cb
		a*a - ca*ca = b*b - (c - ca)^2
		a*a - ca*ca = b*b - c*c + 2*c*ca - ca*ca
		2*c*ca = - b*b + c*c + a*a
		ca = (c*c - b*b + a*a) / 2*c
		h = a*a - ca*ca
	*/
	d := c.C.To(s.C)
	l := d.Len2()
	sl := math.Sqrt(l)
	if sl > c.R+s.R {
		return
	}
	a := (l - s.R*s.R + c.R*c.R) / (2 * sl)
	h := math.Sqrt(c.R*c.R - a*a)

	n := d.Divided(sl)
	f := n.Scaled(a)
	g := n.Scaled(h).Normal()
	c.C.AddE(f)

	return c.C.Add(g), c.C.Sub(g)
}

// ProjectX projects x to y using Polynomial
func (c Circ) ProjectX(x float64) (a, b float64) {
	/*
		(X - c.C.X)^2 + (Y - c.C.Y)^2 = c.R^2
		Y^2 - 2*Y*s.C.Y + s.C.Y^2 + (X - c.C.X)^2 - c.R^2 = 0
	*/
	x -= c.C.X
	return Polynomial(1, -2*c.C.Y, c.C.Y*c.C.Y+x*x-c.R*c.R)
}

// ProjectY projects y to x using Polynomial
func (c Circ) ProjectY(y float64) (a, b float64) {
	/*
		(X - c.C.X)^2 + (Y - c.C.Y)^2 = c.R^2
		X^2 - 2*X*c.C.X + c.C.X^2 + (Y - c.C.Y)^2 - c.R^2 = 0
	*/
	y -= c.C.Y
	return Polynomial(1, -2*c.C.X, c.C.X*c.C.X+y*y-c.R*c.R)
}

// SimplePrjX performs simple projection x into y on circle
//
// if x cannot be projected NaN is returned
func (c Circ) SimplePrjX(x float64) float64 {
	/*
		formula is derivated from circle formula

			(x - c.C.X)**2 + (y - c.C.Y)**2 == c.R**2

		though solution is simplified that loses one possible solution
		(using sqr), this might be usefull as it is 20x faster the ProjectX
	*/
	x -= c.C.X
	return math.Sqrt(c.R*c.R-x*x) + c.C.Y
}

// SimplePrjY is analogous to SimplePrjX and projects y to x
//
// if y cannot be projected NaN is returned
func (c Circ) SimplePrjY(y float64) float64 {
	y -= c.C.Y
	return math.Sqrt(c.R*c.R-y*y) + c.C.X
}

// Intersects detects intersection between two circles
func (c Circ) Intersects(o Circ) bool {
	/*
		we are avoiding math.Sqrt as multiplication is faster
	*/
	ln := c.R + o.R
	return c.C.To(o.C).Len2() <= ln*ln
}

// Contains returns whether circle contains pos
func (c Circ) Contains(pos Vec) bool {
	return pos.To(c.C).Len2() <= c.R*c.R
}

// Approx compers Circles approximately to surthen presision
func (c Circ) Approx(o Circ, precision int) bool {
	return Approx(c.R, o.R, precision) && c.C.Approx(o.C, precision)
}
