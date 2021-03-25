package angle

import (
	"math"
	"testing"

	"github.com/jakubDoka/mlok/mat"
)

func TestNorm(t *testing.T) {
	test := func(name string, input, output float64) {
		r := Norm(input)
		if !mat.Approx(r, output, 8) {
			t.Error(name, r, output)
		}
	}

	test("nothing", 2, 2)
	test("nothing2", -1, -1)
	test("negative", -Pi2-1, -1)
	test("positive", Pi2+1, 1)
	test("big", 100*Pi2+1, 1)
}

func TestTo(t *testing.T) {
	test := func(name string, a, b, output float64) {
		r := To(a, b)
		if !mat.Approx(r, output, 8) {
			t.Error(name, r, output)
		}
	}

	test("nothing", 1, 2, 1)
	test("nothing2", -1, -2, -1)
	test("short step", math.Pi, -math.Pi, 0)
	test("edge step", math.Pi*.7, -math.Pi*.7, math.Pi*.6)
	test("edge step negative", -math.Pi*.7, math.Pi*.7, -math.Pi*.6)
}

func TestTurn(t *testing.T) {
	test := func(name string, start, end, vel, output float64) {
		r := Turn(start, end, vel)
		if !mat.Approx(r, output, 8) {
			t.Error(name, r, output)
		}
	}

	test("nothing", 0, 0, 0, 0)
	test("instant", 1, 2, 3, 2)
	test("negative", 3, 1, 1, 2)
	test("positive", 1, 3, 1, 2)
	test("shorter positive", 2, -3, 1, 3)
	test("shorter negative", -2, 3, 1, -3)
}
