// Package mat3 is still in progress and will documented after its finished
package mat3

// Mat is 3D matrix
//
// layout:
//
//	| m.X.X m.Y.X m.Z.X m.C.X |
//	| m.X.Y m.Y.Y m.Z.Y m.C.Y |
//	| m.X.Z m.Y.Z m.Z.Z m.C.Z |
//  |   0     0     0     1   |
type Mat struct {
	X, Y, Z, C Vec
}

var IM = Mat{
	Vec{1, 0, 0}, // 0
	Vec{0, 1, 0}, // 0
	Vec{0, 0, 1}, // 0
	Vec{0, 0, 0}, // 1
}

// Mul performs matrix multiplication (m by o)
func (m Mat) Mul(o Mat) Mat {
	return Mat{
		Vec{
			m.X.X*o.X.X + m.X.Y*o.Y.X + m.X.Z*o.Z.X,
			m.X.X*o.X.Y + m.X.Y*o.Y.Y + m.X.Z*o.Z.Y,
			m.X.X*o.X.Z + m.X.Y*o.Y.Z + m.X.Z*o.Z.Z,
		},
		Vec{
			m.Y.X*o.X.X + m.Y.Y*o.Y.X + m.Y.Z*o.Z.X,
			m.Y.X*o.X.Y + m.Y.Y*o.Y.Y + m.Y.Z*o.Z.Y,
			m.Y.X*o.X.Z + m.Y.Y*o.Y.Z + m.Y.Z*o.Z.Z,
		},
		Vec{
			m.Z.X*o.X.X + m.Z.Y*o.Y.X + m.Z.Z*o.Z.X,
			m.Z.X*o.X.Y + m.Z.Y*o.Y.Y + m.Z.Z*o.Z.Y,
			m.Z.X*o.X.Z + m.Z.Y*o.Y.Z + m.Z.Z*o.Z.Z,
		},
		Vec{
			m.C.X*o.X.X + m.C.Y*o.Y.X + m.C.Z*o.Z.X + o.C.X,
			m.C.X*o.X.Y + m.C.Y*o.Y.Y + m.C.Z*o.Z.Y + o.C.Y,
			m.C.X*o.X.Z + m.C.Y*o.Y.Z + m.C.Z*o.Z.Z + o.C.Z,
		},
	}
}

// Project performs projection on point and returns a projected result
func (m *Mat) Project(v Vec) Vec {
	return Vec{
		v.X*m.X.X + v.Y*m.Y.X + v.Z*m.Z.X + m.C.X,
		v.X*m.X.Y + v.Y*m.Y.Y + v.Z*m.Z.Y + m.C.Y,
		v.X*m.X.Z + v.Y*m.Y.Z + v.Z*m.Z.Z + m.C.Z,
	}
}

func (m Mat) Moved(v Vec) Mat {
	m.C.AddE(v)

	return m
}

func (m Mat) ScaleComp(v Vec) Mat {
	m.X.MulE(v)
	m.Y.MulE(v)
	m.Z.MulE(v)
	m.C.MulE(v)

	return m
}

// Unproject performs reverse action of Project
func (m *Mat) Unproject(u Vec) (v Vec) {
	/*
		1) v.X*m.X.X + v.Y*m.Y.X + v.Z*m.Z.X + m.C.X = u.X
		2) v.X*m.X.Y + v.Y*m.Y.Y + v.Z*m.Z.Y + m.C.Y = u.Y
		3) v.X*m.X.Z + v.Y*m.Y.Z + v.Z*m.Z.Z + m.C.Z = u.Z

		get v.X from 1:
			v.X*m.X.X + v.Y*m.Y.X + v.Z*m.Z.X + m.C.X = u.X
			v.Y*m.Y.X + v.Z*m.Z.X + m.C.X - u.X = - v.X*m.X.X

			(- v.Y*m.Y.X - v.Z*m.Z.X - m.C.X + u.X)/m.X.X = v.X
			a := -m.C.X + u.X
			v.X = (- v.Y*m.Y.X - v.Z*m.Z.X + a)/m.X.X
		substitute to 2 and get v.Y:
			((- v.Y*m.Y.X - v.Z*m.Z.X + a)/m.X.X)*m.X.Y + v.Y*m.Y.Y + v.Z*m.Z.Y + m.C.Y = u.Y
			- v.Y*m.Y.X*m.X.Y - v.Z*m.Z.X*m.X.Y + a*m.X.Y + v.Y*m.Y.Y*m.X.X + v.Z*m.Z.Y*m.X.X + m.C.Y*m.X.X = u.Y*m.X.X
			- v.Z*m.Z.X*m.X.Y + a*m.X.Y + v.Z*m.Z.Y*m.X.X + m.C.Y*m.X.X - u.Y*m.X.X = v.Y*m.Y.X*m.X.Y - v.Y*m.Y.Y*m.X.X

			(- v.Z*m.Z.X*m.X.Y + a*m.X.Y + v.Z*m.Z.Y*m.X.X + m.C.Y*m.X.X - u.Y*m.X.X) / (m.Y.X*m.X.Y - m.Y.Y*m.X.X) =  v.Y
			b := a*m.X.Y + m.C.Y*m.X.X - u.Y*m.X.X
			c := m.Y.X*m.X.Y - m.Y.Y*m.X.X
			v.Y = (- v.Z*m.Z.X*m.X.Y + v.Z*m.Z.Y*m.X.X + b) / c
		substitute to 3 and get v.Y:
			((- v.Y*m.Y.X - v.Z*m.Z.X + a)/m.X.X)*m.X.Z + v.Y*m.Y.Z + v.Z*m.Z.Z + m.C.Z = u.Z
			- v.Y*m.Y.X*m.X.Z - v.Z*m.Z.X*m.X.Z + a*m.X.Z + v.Y*m.Y.Z*m.X.X + v.Z*m.Z.Z*m.X.X + m.C.Z*m.X.X = u.Z*m.X.X
			- v.Z*m.Z.X*m.X.Z + a*m.X.Z + v.Z*m.Z.Z*m.X.X + m.C.Z*m.X.X - u.Z*m.X.X = v.Y*m.Y.X*m.X.Z - v.Y*m.Y.Z*m.X.X

			(- v.Z*m.Z.X*m.X.Z + a*m.X.Z + v.Z*m.Z.Z*m.X.X + m.C.Z*m.X.X - u.Z*m.X.X) / (m.Y.X*m.X.Z - m.Y.Z*m.X.X) = v.Y
			d := a*m.X.Z + m.C.Z*m.X.X - u.Z*m.X.X
			e := m.Y.X*m.X.Z - m.Y.Z*m.X.X
			v.Y = (- v.Z*m.Z.X*m.X.Z + v.Z*m.Z.Z*m.X.X + d) / e
		put derived functions into equality and get v.Z
			(- v.Z*m.Z.X*m.X.Z + v.Z*m.Z.Z*m.X.X + d) / e = (- v.Z*m.Z.X*m.X.Y + v.Z*m.Z.Y*m.X.X + b) / c
			- v.Z*m.Z.X*m.X.Z*c + v.Z*m.Z.Z*m.X.X*c + d*c = - v.Z*m.Z.X*m.X.Y*e + v.Z*m.Z.Y*m.X.X*e + b*e
			- v.Z*m.Z.X*m.X.Z*c + v.Z*m.Z.Z*m.X.X*c + v.Z*m.Z.X*m.X.Y*e - v.Z*m.Z.Y*m.X.X*e = b*e - d*c
		thus the final expression
			v.Z = (b*e - d*c) / (- m.Z.X*m.X.Z*c + m.Z.Z*m.X.X*c + m.Z.X*m.X.Y*e - m.Z.Y*m.X.X*e)

	*/

	a := -m.C.X + u.X
	b := a*m.X.Y + m.C.Y*m.X.X - u.Y*m.X.X
	c := m.Y.X*m.X.Y - m.Y.Y*m.X.X
	d := a*m.X.Z + m.C.Z*m.X.X - u.Z*m.X.X
	e := m.Y.X*m.X.Z - m.Y.Z*m.X.X

	v.Z = (b*e - d*c) / (-m.Z.X*m.X.Z*c + m.Z.Z*m.X.X*c + m.Z.X*m.X.Y*e - m.Z.Y*m.X.X*e)

	if e != 0 {
		v.Y = (-v.Z*m.Z.X*m.X.Z + v.Z*m.Z.Z*m.X.X + d) / e
	} else {
		v.Y = (-v.Z*m.Z.X*m.X.Y + v.Z*m.Z.Y*m.X.X + b) / c
	}

	v.X = (-v.Y*m.Y.X - v.Z*m.Z.X + a) / m.X.X

	return
}

func (m Mat) Approx(o Mat, precision int) bool {
	return m.X.Approx(o.X, precision) && m.Y.Approx(o.Y, precision) && m.Z.Approx(o.Z, precision) && m.C.Approx(o.C, precision)
}
