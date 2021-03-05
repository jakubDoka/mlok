package refl

import (
	"reflect"
	"strconv"
	"testing"
	"unicode/utf8"
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
			&Test1{1, 1},
		},

		{
			"Nested",
			&Test2{Test1{1, 0}, 0, ""},
			&Test2{Test1{0, 1}, 1, "asd"},
			&Test2{Test1{1, 1}, 1, "asd"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			Overwrite(tC.target, tC.defaults, true)
			if !reflect.DeepEqual(tC.target, tC.result) {
				t.Error(tC.target, tC.result)
			}
		})
	}
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
		Overwrite(&value, &defa, true)
	}
}

type Test3 string

func (t Test3) Dec() (int, error) {
	return strconv.Atoi(string(t))
}

type Test4 struct {
	A, B, C Test3
}

type Test5 struct {
	A, B, C int
}

func TestConvert(t *testing.T) {
	t1s := Test3("10")
	t1d := 0
	t1r := 10
	t2s := Test4{Test3("10"), Test3("100"), Test3("4")}
	t2d := Test5{}
	t2r := Test5{10, 100, 4}
	testCases := []struct {
		desc           string
		scr, dest, res interface{}
		fail           bool
	}{
		{
			desc: "string to int",
			scr:  &t1s,
			dest: &t1d,
			res:  &t1r,
		},
		{
			desc: "nested",
			scr:  &t2s,
			dest: &t2d,
			res:  &t2r,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := Convert(tC.scr, tC.dest, "Dec")
			if tC.fail && err != nil {
				t.Fail()
			}

			if tC.fail {
				return
			}

			if !reflect.DeepEqual(tC.dest, tC.res) {
				t.Error(reflect.ValueOf(tC.dest).Elem().Interface(), reflect.ValueOf(tC.res).Elem().Interface())
			}
		})
	}
}

func TestF(t *testing.T) {
	r, n := utf8.DecodeRuneInString("\n")
	t.Error(r, n, '\n')
}
