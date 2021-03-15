package mat3

import (
	"testing"
)

func BenchmarkCross(b *testing.B) {
	v := Vec{1, 0, 6}
	u := Vec{4, 8, 9}
	for i := 0; i < b.N; i++ {
		v.Cross(u)
	}
}

func BenchmarkCross2(b *testing.B) {
	v := Vec{1, 0, 6}
	u := Vec{4, 8, 9}
	for i := 0; i < b.N; i++ {
		v.Cross2(u)
	}
}
