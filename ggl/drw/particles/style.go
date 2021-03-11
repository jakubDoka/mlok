package particles

import (
	"gobatch/ggl"
	"gobatch/ggl/drw"
	"gobatch/load"
	"gobatch/mat"
	"gobatch/mat/lerp"
	"sync"

	"github.com/jakubDoka/goml/goss"
)

// Type defines behavior of a particles
type Type struct {
	EmissionShape

	base                           Drawer
	indices, vertexes, threadCount int
	dws                            []ggl.Drawer
	mut                            sync.Mutex

	Color lerp.Gradient

	ScaleMultiplier, Acceleration, TwerkAcceleration lerp.Tween

	Velocity, Livetime, Twerk, Rotation, Spread, Scale lerp.Tween

	Friction, OriginGravity float64

	Gravity mat.Vec

	RotationRelativeToVelocity bool
}

// SetDrawer sets the Data drawer
func (d *Type) SetDrawer(drw Drawer) {
	d.base = drw
	d.indices, d.vertexes = drw.Metrics()
}

// setThreads changes threadcount for witch Type is used
func (d *Type) setThreads(count int) {
	d.mut.Lock()
	d.threadCount = count
	d.dws = make([]ggl.Drawer, count)
	for i := range d.dws {
		d.dws[i] = d.base.Copy()
	}
	d.mut.Unlock()
}

// RawStyle extends load.Rawstile by particle specific functionality
type RawStyle struct {
	load.RawStyle
}

// WrapStyle wraps a goss.Style into Raw Style
func WrapStyle(s goss.Style) RawStyle {
	return RawStyle{load.RawStyle{Style: s}}
}

// Gradient parses gradient, if parsing fails def is returned
//
// syntax:
// 	color - single color syntax applies witch will result into lerp.ConstColor
//  color d color - two colors with no positions separated by 'd' will evaluate to lerp.LinearColor
//  color p float d color p float - 2 or more colors separated by 'd' that have positions specified
// separated by 'p' will result into lerp.ChainedColor fed with lerp.CP(color, float), 'p' is a separator
// but can be omitted in case color definition did not need to be terminated
// (white 0.2 d 1 1 1 1 1 is valid syntax)
func (r RawStyle) Gradient(key string, custom map[string]lerp.Gradient, def lerp.Gradient) (e lerp.Gradient) {
	e = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	var colors lerp.ChainedColor
	for i, n := 0, -1; ; i += n {
		var color mat.RGBA
		var position float64
		color, n = load.ParseRGBA(val[i:])
	bck:
		if n == 0 || i+n >= len(val) {
			colors = append(colors, lerp.CP(position, color))
			break
		}

		switch v := val[i+n].(type) {
		case float64:
			position = v
			i++
			goto bck
		case int:
			position = float64(v)
			i++
			goto bck
		case string:
			switch v {
			case "d":
				colors = append(colors, lerp.CP(position, color))
				i++
			case "p":
				i++
				goto bck
			default:
				return
			}
		default:
			return
		}
	}

	switch len(colors) {
	case 0:
		return
	case 1:
		return lerp.ConstColor(colors[0].Color)
	case 2:
		if colors[0].Position == 0 && colors[1].Position == 0 {
			return lerp.LC(colors[0].Color, colors[1].Color)
		}
		fallthrough
	default:
		return colors
	}
}

// Tween parses a tween, if parsing fails, def is returned
func (r RawStyle) Tween(key string, custom map[string]lerp.Tween, def lerp.Tween) (e lerp.Tween) {
	e = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	props := [4]float64{}
	load.CollectFloats(val, props[:])

	switch v := val[0].(type) {
	case int:
		return lerp.Const(v)
	case float64:
		return lerp.Const(v)
	case string:
		switch v {
		case "random":
			return lerp.R(props[0], props[1])
		case "bezier":
			return lerp.B(props[0], props[1], props[2], props[3])
		case "linear":
			return lerp.L(props[0], props[1])
		default:
			val, ok := custom[v]
			if ok {
				return val
			}
		}
	}

	return
}

// Drawer parses drawer under the key, if parsing fails, def is returned
func (r RawStyle) Drawer(key string, custom map[string]Drawer, def Drawer) (e Drawer) {
	e = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	name, ok := val[0].(string)
	if !ok {
		return
	}

	props := [2]float64{}
	n := load.CollectFloats(val, props[:])

	switch n {
	case 2:
		switch name {
		case "circle":
			return &Circle{drw.NCircle(props[0], int(props[1]))}
		case "rect":
			s := ggl.NSprite(mat.A(0, 0, props[0], props[1]))
			s.SetIntensity(0)
			return &Sprite{s}
		}
	case 0:
		val, ok := custom[name]
		if !ok {
			return
		}
		return val
	}

	return
}

// EmissionShape parses emission shape under given key, if parsing fails def is returned
func (r RawStyle) EmissionShape(key string, custom map[string]EmissionShape, def EmissionShape) (e EmissionShape) {
	e = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	name, ok := val[0].(string)
	if !ok {
		return
	}

	switch len(val) {
	case 3:
		props := [2]float64{}
		for i := 0; i < len(props); i++ {
			switch v := val[i+1].(type) {
			case float64:
				props[i] = v
			case int:
				props[i] = float64(v)
			default:
				return
			}
		}

		switch name {
		case "circle":
			return Circular{props[0], props[1]}
		case "rectangle":
			return Rectangle{props[0], props[1]}
		}
	case 1:
		val, ok := custom[name]
		if !ok {
			return
		}
		return val
	}
	return
}
