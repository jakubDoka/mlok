package particles


// Resize resizes the particles
func (v *particles) Resize(size int) {
	if cap(*v) >= size {
		*v = (*v)[:size]
	} else {
		ns := make(particles, size)
		copy(ns, *v)
		*v = ns
	}
}

