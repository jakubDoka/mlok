package mat


// Point is a  vector type with X and Y coordinates.
type Point struct {
	X, Y int
}

// ZP is zero value Point
var ZP Point

// P returns a new  vector with the given coordinates.
func P(x, y int) Point {
	return Point{x, y}
}

// XY returns the components of the vector in two return values.
func (v Point) XY() (x, y int) {
	return v.X, v.Y
}

// Add returns the sum of vectors v and v.
func (v Point) Add(u Point) Point {
	return Point{
		v.X + u.X,
		v.Y + u.Y,
	}
}

// Sub subtracts u from v and returns recult.
func (v Point) Sub(u Point) Point {
	return Point{
		v.X - u.X,
		v.Y - u.Y,
	}
}

// Div returns the vector v divided by the vector u component-wise.
func (v Point) Div(u Point) Point {
	return Point{v.X / u.X, v.Y / u.Y}
}

// Mul returns the vector v multiplied by the vector u component-wise.
func (v Point) Mul(u Point) Point {
	return Point{v.X * u.X, v.Y * u.Y}
}

// To returns the vector from v to u. Equivalent to u.Sub(v).
func (v Point) To(u Point) Point {
	return Point{
		u.X - v.X,
		u.Y - v.Y,
	}
}

// AddE same as v = v.Add(u)
func (v *Point) AddE(u Point) {
	v.X += u.X
	v.Y += u.Y
}

// SubE same as v = v.Sub(u)
func (v *Point) SubE(u Point) {
	v.X -= u.X
	v.Y -= u.Y
}

// MulE same as v = v.Mul(u)
func (v *Point) MulE(u Point) {
	v.X *= u.X
	v.Y *= u.Y
}

// DivE same as v = v.Div(u)
func (v *Point) DivE(u Point) {
	v.X /= u.X
	v.Y /= u.Y
}

// Scaled returns the vector v multiplied by c.
func (v Point) Scaled(c int) Point {
	return Point{v.X * c, v.Y * c}
}

// Divided returns the vector v divided by c.
func (v Point) Divided(c int) Point {
	return Point{v.X / c, v.Y / c}
}

// Inv returns v with both components inverted
func (v Point) Inv() Point {
	v.X = -v.X
	v.Y = -v.Y
	return v
}

// Flatten flatens the Point into Array, values are
// ordered as they would on stack
func (v Point) Flatten() [2]int {
	return [...]int{v.X, v.Y}
}

// Mutator similar to Flatten returns array with vector components
// though these are pointers to componenets instead
func (v *Point) Mutator() [2]*int {
	return [...]*int{&v.X, &v.Y}
}

