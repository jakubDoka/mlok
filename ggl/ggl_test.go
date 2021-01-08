package ggl

import (
	"gobatch/mt"
	"reflect"
	"testing"
)

type Test struct {
	A, B, C, D float64
}

type Test2 struct {
	Test
	B  float64
	FF []float64
}

func TestAssertHomogeneity(t *testing.T) {
	tests := []struct {
		tested, supposed interface{}
		name             string
	}{
		{
			Test{},
			float64(0),
			"Simple struct",
		},
		{
			[]float64{},
			float64(0),
			"slice",
		},
		{
			Test2{},
			float64(0),
			"complex",
		},
	}
	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			if err := AssertHomogeneity(reflect.TypeOf(te.tested), reflect.TypeOf(te.supposed)); err != nil {
				t.Error(err)
			}
		})
	}

}

func BenchmarkSprite2D(b *testing.B) {
	s := NSprite2D(mt.AABB{})
	for i := 0; i < b.N; i++ {
		s.Update(mt.IM2, mt.RGBA{})
	}
}
