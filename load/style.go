package load

import (
	"github.com/jakubDoka/mlok/mat"
	"github.com/jakubDoka/mlok/mat/rgba"

	"github.com/jakubDoka/goml/goss"
)

// Fill is constant important to ui, it is also redeclare here so check it out
const Fill float64 = 100000.767898765556788777667787666

// RawStyle is wrapper of goss.RawStyle and adds extra functionality
type RawStyle struct {
	goss.Style
}

// Sub retrieves Substyle under the key, or returns def
func (r RawStyle) Sub(key string, def RawStyle) (v RawStyle) {
	v = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	if val, ok := val[0].(goss.Style); ok {
		return RawStyle{val}
	}

	return
}

func (r RawStyle) Int(key string, def int) (v int) {
	v = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	switch v := val[0].(type) {
	case int:
		return v
	case float64:
		return int(v)
	}

	return
}

// Bool returns boolean value under the key, or def if retrieval fails
func (r RawStyle) Bool(key string, def bool) (v bool) {
	v = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	switch v := val[0].(type) {
	case bool:
		return v
	case int:
		return v == 1
	}

	return
}

// Float returns float under the key of default value if obtaining failed
func (r RawStyle) Float(key string, def float64) (v float64) {
	v = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	switch v := val[0].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}

	return v
}

// Ident returns string from RawStyle or def if not present
func (r RawStyle) Ident(key, def string) string {
	val, ok := r.Style.Ident(key)
	if !ok {
		return def
	}
	return val
}

// Vec returns vector under tha key, if parsing fails or vec is not present def is returned
func (r RawStyle) Vec(key string, def mat.Vec) (u mat.Vec) {
	u = def

	val, ok := r.Style[key]
	if !ok {
		return
	}

	components := [2]float64{}
	n := CollectFloats(val, components[:])

	switch n {
	case 1:
		return mat.V(components[0], components[0])
	case 2:
		return mat.V(components[0], components[1])
	}

	return
}

// Region returns a texture region, it can be hardcoded aabb or region name.
func (r RawStyle) Region(key string, regions map[string]mat.AABB, def mat.AABB) mat.AABB {
	m, ok := regions[r.Ident(key, "")]
	if !ok {
		return r.AABB(key, def)
	}
	return m
}

// AABB parser margin under the key, if parsing fails or margin is not present, default is returned
func (r RawStyle) AABB(key string, def mat.AABB) (m mat.AABB) {
	m = def

	val, ok := r.Style[key]
	if !ok {
		return
	}

	sides := [4]float64{}
	n := CollectFloats(val, sides[:])

	switch n {
	case 1:
		return mat.A(sides[0], sides[0], sides[0], sides[0])
	case 2:
		return mat.A(sides[0], sides[1], sides[0], sides[1])
	case 4:
		return mat.A(sides[0], sides[1], sides[2], sides[3])
	}

	return
}

// RGBA returns a color under the key, if color parsing fails or color is not present, def is returned
// this also accepts names mapped in mlok/mat/rgba.Colors
func (r RawStyle) RGBA(key string, def mat.RGBA) (c mat.RGBA) {
	c = def

	val, ok := r.Style[key]
	if !ok {
		return
	}

	res, n := ParseRGBA(val)
	if n != 0 {
		return res
	}

	return
}

// ParseRGBA parses slice if interfaces into color and returns how big the color was
//
// syntax:
//
//	name - simple name of color can be enough if its contained in map in mlok/mat/rgba package
//  a - specify only one float/int channel and mat.Alpha(alpha) will be returned
//  r, g, b - specify threes floats/ints and mat.RGB(r, g, b) will be returned
//  r, g, b, a - specify four floats/ints and mat.RGBA{r, g, b, a} will be returned
//  hex - hex color notation
func ParseRGBA(values []interface{}) (c mat.RGBA, ln int) {
	if v, ok := values[0].(string); ok {
		if v, ok := rgba.Colors[v]; ok {
			return v, 1
		}
		col, err := mat.HexToRGBA(v)
		if err != nil {
			return col, 1
		}
		return
	}

	channels := [4]float64{}
	ln = CollectFloats(values, channels[:])

	switch ln {
	case 1:
		return mat.Alpha(channels[0]), 1
	case 3:
		return mat.RGB(channels[0], channels[1], channels[2]), 3
	case 4:
		return mat.RGBA{
			R: channels[0],
			G: channels[1],
			B: channels[2],
			A: channels[3],
		}, 4
	}
	return
}

// CollectFloats collects consecutive floating point numbers and returns how match was collected
func CollectFloats(values []interface{}, buff []float64) int {
	for i := 0; i < len(values) && i < len(buff); i++ {
		switch v := values[i].(type) {
		case float64:
			buff[i] = v
		case int:
			buff[i] = float64(v)
		case string:
			if v != "fill" {
				return i
			}
			buff[i] = Fill
		default:
			return i
		}
	}
	return mat.Mini(len(buff), len(values))
}
