package drw


// Resize resizes the Base
func (v *Base) Resize(size int) {
	if cap(*v) >= size {
		*v = (*v)[:size]
	} else {
		ns := make(Base, size)
		copy(ns, *v)
		*v = ns
	}
}

