package mat3

type Mat struct {
	X, Y, Z, C Vec
}

var IM = Mat{
	Vec{1, 0, 0}, // 0
	Vec{0, 1, 0}, // 0
	Vec{0, 0, 1}, // 0
	Vec{0, 0, 0}, // 1
}

func (m *Mat) Project(v Vec) Vec {
	return Vec{
		v.X*m.X.X + v.X*m.Y.X + v.X*m.Z.X + m.C.X,
		v.Y*m.X.Y + v.Y*m.Y.Y + v.Y*m.Z.Y + m.C.Y,
		v.Z*m.X.Z + v.Z*m.Y.Z + v.Z*m.Z.Z + m.C.Z,
	}
}
