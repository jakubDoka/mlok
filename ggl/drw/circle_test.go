package drw

import (
	"gobatch/ggl"
	"reflect"
	"testing"
)

func TestCircle(t *testing.T) {
	c := NCircle(10, 4)
	r := ggl.Indices{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 1}
	if !reflect.DeepEqual(c.Indices, r) {
		t.Errorf("\n%#v\n%#v", c.Indices, r)
	}
}
