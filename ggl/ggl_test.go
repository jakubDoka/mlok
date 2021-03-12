package ggl

import (
	"testing"

	"github.com/jakubDoka/gobatch/mat"
)

func TestNPSResize(t *testing.T) {
	n := NPatch(mat.A(0, 0, 50, 50), mat.A(10, 10, 10, 10))
	res := [3][3][4]mat.Vec{
		{
			{mat.V(-25, -25), mat.V(-25, -15), mat.V(-15, -15), mat.V(-15, -25)},
			{mat.V(-15, -25), mat.V(-15, -15), mat.V(15, -15), mat.V(15, -25)},
			{mat.V(15, -25), mat.V(15, -15), mat.V(25, -15), mat.V(25, -25)},
		},
		{
			{mat.V(-25, -15), mat.V(-25, 15), mat.V(-15, 15), mat.V(-15, -15)},
			{mat.V(-15, -15), mat.V(-15, 15), mat.V(15, 15), mat.V(15, -15)},
			{mat.V(15, -15), mat.V(15, 15), mat.V(25, 15), mat.V(25, -15)},
		},
		{
			{mat.V(-25, 15), mat.V(-25, 25), mat.V(-15, 25), mat.V(-15, 15)},
			{mat.V(-15, 15), mat.V(-15, 25), mat.V(15, 25), mat.V(15, 15)},
			{mat.V(15, 15), mat.V(15, 25), mat.V(25, 25), mat.V(25, 15)},
		},
	}
	got := [3][3][4]mat.Vec{}
	for y, v := range got {
		for x := range v {
			got[y][x] = n.s[y][x].tex
		}
	}

	if res != got {
		t.Errorf("\n%v\n%v", res, got)
	}
}
