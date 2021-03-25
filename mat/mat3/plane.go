package mat3

type Plane struct {
	O, N Vec
}

// Formula returns 0 if point belongs to Plane
func (p Plane) Formula(pos Vec) float64 {
	/*
		a*X + b*Y + c*Z + d = 0
		p.N.X*X + p.N.Y*Y + p.N.Z*Z + d = 0

		p.N.X*p.O.X + p.N.Y*p.O.Y + p.N.Z*p.O.Z + d = 0
		- p.N.X*p.O.X - p.N.Y*p.O.Y - p.N.Z*p.O.Z = d

		p.N.X*pos.X + p.N.Y*pos.Y + p.N.Z*pos.Z - p.N.X*p.O.X - p.N.Y*p.O.Y - p.N.Z*p.O.Z = 0
	*/
	return p.N.X*pos.X + p.N.Y*pos.Y + p.N.Z*pos.Z - p.N.X*p.O.X - p.N.Y*p.O.Y - p.N.Z*p.O.Z
}
