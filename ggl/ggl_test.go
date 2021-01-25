package ggl

import (
	"gobatch/mat"
	"testing"
)

func BenchmarkSprite2D(b *testing.B) {
	s := NSprite2D(mat.AABB{})
	for i := 0; i < b.N; i++ {
		s.Update(mat.IM2, mat.RGBA{})
	}
}
