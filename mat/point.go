package mat

// Clamp clamps point between two points
func (v Point) Clamp(min, max Point) Point {
	return v.Max(min).Min(max)
}

// Max returns max of two points componentvise
func (v Point) Max(b Point) Point {
	return Point{Maxi(v.X, b.X), Maxi(v.Y, b.Y)}
}

// Min returns min of two points componentvise
func (v Point) Min(b Point) Point {
	return Point{Mini(v.X, b.X), Mini(v.Y, b.Y)}
}

// Maxi ...
func Maxi(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Mini ...
func Mini(a, b int) int {
	if a < b {
		return a
	}
	return b
}
