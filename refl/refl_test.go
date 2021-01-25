package refl

import (
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

func TestOverwriteDefault(t *testing.T) {
	type Test1 struct {
		A, B int
	}

	type Test2 struct {
		T Test1
		C int
		G string
	}

	testCases := []struct {
		desc                     string
		target, defaults, result interface{}
	}{
		{
			"Simple",
			&Test1{1, 0},
			&Test1{0, 1},
			Test1{1, 1},
		},

		{
			"Nested",
			&Test2{Test1{1, 0}, 0, ""},
			&Test2{Test1{0, 1}, 1, "asd"},
			Test2{Test1{1, 1}, 1, "asd"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			OverwriteDefault(tC.target, tC.defaults)
			if reflect.DeepEqual(tC.target, tC.result) {
				t.Error(tC.target, tC.result)
			}
		})
	}
	t.Fail()
}

func BenchmarkOverwriteDefault(b *testing.B) {
	type Bench struct {
		A, B, C int
	}

	type BenchS struct {
		A, B, C    int
		F, G       Bench
		S, D, H, I int
	}

	value := BenchS{
		1, 1, 0, Bench{}, Bench{1, 1, 0}, 0, 0, 0, 0,
	}

	defa := BenchS{
		1, 0, 1, Bench{1, 1, 1}, Bench{1, 0, 1}, 1, 1, 1, 1,
	}

	for i := 0; i < b.N; i++ {
		OverwriteDefault(&value, &defa)
	}
}
