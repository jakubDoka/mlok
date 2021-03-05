package mat

// Clamp clamps point between two points
func (v Point) Clamp(min, max Point) Point {
	return v.Max(min).Min(max)
}

// Max returns max of two points componentvise
func (v Point) Max(b Point) Point {
	return Point{Max(v.X, b.X), Max(v.Y, b.Y)}
}

// Min returns min of two points componentvise
func (v Point) Min(b Point) Point {
	return Point{Min(v.X, b.X), Min(v.Y, b.Y)}
}

// Max ...
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min ...
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
