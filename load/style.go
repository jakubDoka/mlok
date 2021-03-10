package load

import (
	"gobatch/mat"
	"gobatch/mat/rgba"

	"github.com/jakubDoka/goml/goss"
)

// Fill is constant important to ui, it is also redeclare here so check it out
const Fill float64 = 100000.767898765556788777667787666

// RawStyle is wrapper of goss.RawStyle and adds extra functionality
type RawStyle struct {
	goss.Style
}

// Bool returns boolean value under the key of def, if retrieval fails
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
	for i := 0; i < len(components) && i < len(val); i++ {
		switch v := val[i].(type) {
		case float64:
			components[i] = v
		case int:
			components[i] = float64(v)
		case string:
			if v != "fill" {
				return
			}
			components[i] = Fill
		default:
			return
		}
	}

	switch len(val) {
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
	for i := 0; i < len(sides) && i < len(val); i++ {
		switch v := val[i].(type) {
		case float64:
			sides[i] = v
		case int:
			sides[i] = float64(v)
		case string:
			if v != "fill" {
				return
			}
			sides[i] = Fill
		default:
			return
		}

	}

	switch len(val) {
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
// this also accepts names mapped in gobatch/mat/rgba.Colors
func (r RawStyle) RGBA(key string, def mat.RGBA) (c mat.RGBA) {
	c = def

	val, ok := r.Style[key]
	if !ok {
		return
	}

	if v, ok := val[0].(string); ok {
		if v, ok := rgba.Colors[v]; ok {
			return v
		}
		return
	}

	channels := [4]float64{}
	for i := 0; i < len(channels) && i < len(val); i++ {

		switch vl := val[i].(type) {
		case float64:
			channels[i] = vl
		case int:
			channels[i] = float64(vl)
		default:
			return
		}
	}

	switch len(val) {
	case 1:
		return mat.Alpha(channels[0])
	case 3:
		return mat.RGB(channels[0], channels[1], channels[2])
	case 4:
		return mat.RGBA{
			R: channels[0],
			G: channels[1],
			B: channels[2],
			A: channels[3],
		}
	}

	return
}
