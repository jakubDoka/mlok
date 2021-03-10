package mat

import "math"

// Circ is an circle
type Circ struct {
	C Vec
	R float64
}

// C is circle constructor
func C(x, y, r float64) Circ {
	return Circ{Vec{x, y}, r}
}

// Union returns smallest possible circle that encloses both circles
func (c Circ) Union(o Circ) Circ {
	d := c.C.To(o.C)
	l := d.Len()

	c.C.AddE(Rad(d.Angle(), l*.5+o.R-c.R))
	c.R += l + o.R

	return c
}

// IntersectsAABB returns whether circle intersects AABB
/*func (c Circ) IntersectsAABB(o AABB) bool {
	l, b, r, t := c.C.X > o.Min.X, c.C.X < o.Max.X, c.C.Y > o.Min.Y, c.C.X < o.Min.Y

	if l && b && r && t {
		return true
	}

	if t && b && l && c.C.X-c.R < o.Max.X {
		return true
	}

	if t && b && r && c.C.X+c.R > o.Min.X {
		return true
	}

	if l && r && t && c.C.Y-c.R < o.Max.Y {
		return true
	}

	if l && r && b && c.C.Y+c.R > o.Min.Y {
		return true
	}

	v := o.Vertices()

	if b && r && c.Contains(v[0]) {
		return true
	}

	if b && l && c.Contains(v[1]) {
		return true
	}

	if l && t && c.Contains(v[2]) {
		return true
	}

	if r && t && c.Contains(v[3]) {
		return true
	}

	return false
}*/

// PrjX projects x to y using Polynomial
func (c Circ) PrjX(x float64) [2]float64 {
	x -= c.C.X
	return Polynomial(1, 2*c.C.Y, c.C.Y*c.C.Y+x*x-c.R*c.R)
}

// PrjY projects y to x using Polynomial
func (c Circ) PrjY(y float64) [2]float64 {
	y -= c.C.Y
	return Polynomial(1, 2*c.C.X, c.C.X*c.C.X+y*y-c.R*c.R)
}

// SimplePrjX performs simple projection x into y on circle
//
// if x cannot be projected NaN is returned
func (c Circ) SimplePrjX(x float64) float64 {
	/*
		formula is derivated from circle formula

			(x - c.C-X)**2 + (y - c.C.Y)**2 == c.R**2

		though solution is simplified that loses one possible solution
		(using sqr), this might be usefull as it is 20x faster the PrjX
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
