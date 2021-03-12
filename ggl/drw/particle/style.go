package particle

import (
	"math"
	"sync"

	"github.com/jakubDoka/gobatch/ggl"
	"github.com/jakubDoka/gobatch/ggl/drw"
	"github.com/jakubDoka/gobatch/load"
	"github.com/jakubDoka/gobatch/mat"
	"github.com/jakubDoka/gobatch/mat/lerp"
	"github.com/jakubDoka/gobatch/mat/lerpc"

	"github.com/jakubDoka/goml/goss"
)

// Parser performs parsing of goss into particle Tapes
type Parser struct {
	Assets

	Parser goss.Parser
	Styles goss.Styles
}

// Construct constructs the Type under given name and returns it
// if there is no goss style defining this Type, nil is returned
func (p *Parser) Construct(name string) *Type {
	t := Type{}
	stl, ok := p.Styles[name]
	if !ok {
		return nil
	}

	t.Parse(WrapStyle(stl), &p.Assets)

	return &t
}

// AddGoss parses the goss source and adds all results to Styles
func (p *Parser) AddGoss(source []byte) error {
	stl, err := p.Parser.Parse(source)

	if err == nil {
		if p.Styles == nil {
			p.Styles = stl
		} else {
			p.Styles.Add(stl)
		}
	}

	return err
}

// Assets can contain custom Objects that Parser uses to load the Type
type Assets struct {
	ColorTweens     map[string]lerpc.Tween
	ColorGenerators map[string]lerpc.Generator
	FloatTweens     map[string]lerp.Tween
	FloatGenerators map[string]lerp.Generator
	Drawers         map[string]Drawer
	EmissionShapes  map[string]EmissionShape
}

// Type defines behavior of a particles
type Type struct {
	EmissionShape

	base                           Drawer
	indices, vertexes, threadCount int
	dws                            []ggl.Drawer
	mut                            sync.Mutex

	Color lerpc.Tween
	Mask  lerpc.Generator

	ScaleMultiplier, Acceleration, TwerkAcceleration lerp.Tween

	Velocity, Livetime, Twerk, Rotation, Spread, ScaleX, ScaleY lerp.Generator

	Friction, OriginGravity float64

	Gravity mat.Vec

	RotationRelativeToVelocity bool
}

// DefType contains default configuration ro Type
var DefType = Type{
	EmissionShape: Point{},
	base:          Square(10),

	Color: lerpc.Const(mat.White),
	Mask:  lerpc.Const(mat.White),

	ScaleMultiplier:   lerp.Const(1),
	Acceleration:      lerp.Const(0),
	TwerkAcceleration: lerp.Const(0),

	Velocity: lerp.Const(0),
	Livetime: lerp.Const(1),
	Rotation: lerp.Const(0),
	Spread:   lerp.Const(0),
	ScaleX:   lerp.Const(1),
	ScaleY:   lerp.Const(1),

	Gravity: mat.V(0, -300),
}

// Parse performs the Type parsing, DefType is used as template, if some field is missing or
// un-parsable Value from defaults is used
func (d *Type) Parse(r RawStyle, a *Assets) {
	d.EmissionShape = r.EmissionShape("emission_shape", a.EmissionShapes, DefType.EmissionShape)
	d.base = r.Drawer("drawer", a.Drawers, DefType.base)

	d.Color = r.ColorTween("color", a.ColorTweens, DefType.Color)
	d.Mask = r.ColorGen("mask", a.ColorGenerators, DefType.Mask)

	d.ScaleMultiplier = r.FloatTween("scale_multiplier", a.FloatTweens, DefType.ScaleMultiplier)
	d.Acceleration = r.FloatTween("acceleration", a.FloatTweens, DefType.Acceleration)
	d.TwerkAcceleration = r.FloatTween("twerk_acceleration", a.FloatTweens, DefType.TwerkAcceleration)

	d.Velocity = r.FloatGen("velocity", a.FloatGenerators, DefType.Velocity)
	d.Livetime = r.FloatGen("livetime", a.FloatGenerators, DefType.Livetime)
	d.Twerk = r.FloatGen("twerk", a.FloatGenerators, DefType.Twerk)
	d.Rotation = r.FloatGen("rotation", a.FloatGenerators, DefType.Rotation)
	d.Spread = r.FloatGen("spread", a.FloatGenerators, DefType.Spread)
	d.ScaleX = r.FloatGen("scale_x", a.FloatGenerators, DefType.ScaleX)
	d.ScaleY = r.FloatGen("scale_y", a.FloatGenerators, DefType.ScaleY)

	d.Gravity = r.Vec("gravity", mat.V(0, -200))
	d.Friction = r.Float("friction", 0)
	d.OriginGravity = r.Float("origin_gravity", 0)
	d.RotationRelativeToVelocity = r.Bool("rotation_relative_to_velocity", false)
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

// ColorTween parses a color tween, syntax documented on RawStyle.color method applies
func (r RawStyle) ColorTween(key string, custom map[string]lerpc.Tween, def lerpc.Tween) (e lerpc.Tween) {
	e = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	if name, ok := val[0].(string); ok {
		if val, ok := custom[name]; ok {
			return val
		}
	}

	res := r.color(val, false)
	if res == nil {
		return
	}

	if val, ok := res.(lerpc.Tween); ok {
		return val
	}

	return
}

// ColorGen parses a color generator, syntax documented on RawStyle.color method applies
func (r RawStyle) ColorGen(key string, custom map[string]lerpc.Generator, def lerpc.Generator) (e lerpc.Generator) {
	e = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	if name, ok := val[0].(string); ok {
		if val, ok := custom[name]; ok {
			return val
		}
	}

	res := r.color(val, true)
	if res == nil {
		return
	}

	if val, ok := res.(lerpc.Generator); ok {
		return val
	}

	return
}

// color parses a color tween, if parsing fails def is returned
//
// syntax:
// 	color - single color syntax applies witch will result into lerp.ConstColor
//  color d color - two colors with no positions separated by 'd' will evaluate to lerp.LinearColor
//  color p float d color p float - 2 or more colors separated by 'd' that have positions specified
// separated by 'p' will result into lerp.ChainedColor fed with lerp.CP(color, float), 'p' is a separator
// but can be omitted in case color definition did not need to be terminated
// (white 0.2 d 1 1 1 1 1 is valid syntax)
func (r RawStyle) color(val []interface{}, random bool) (e interface{}) {

	var colors lerpc.ChainedTween
	for i, n := 0, -1; ; i += n {
		var color mat.RGBA
		var position float64
		color, n = load.ParseRGBA(val[i:])
	bck:
		if n == 0 || i+n >= len(val) {
			colors = append(colors, lerpc.Point(position, color))
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
				colors = append(colors, lerpc.Point(position, color))
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
		return val
	case 1:
		return lerpc.Const(colors[0].Color)
	case 2:
		if colors[0].Position == 0 && colors[1].Position == 0 {
			if random {
				return lerpc.Random(colors[0].Color, colors[1].Color)
			}
			return lerpc.Linear(colors[0].Color, colors[1].Color)
		}
		fallthrough
	default:
		return colors
	}
}

// FloatGen returns float generator, returns def if process failed
//
// syntax:
//	number - returns lerp.Constant(number)
//  random number number - if you prefix field with 'random' and provide two numbers lerp.Random(number, number)
// is returned
//  name - if name is not equal random it will search for custom generator with this name
func (r RawStyle) FloatGen(key string, custom map[string]lerp.Generator, def lerp.Generator) (e lerp.Generator) {
	e = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	props := [2]float64{}
	load.CollectFloats(val[1:], props[:])

	switch v := val[0].(type) {
	case int:
		return lerp.Const(float64(v))
	case float64:
		return lerp.Const(v)
	case string:
		switch v {
		case "random":
			return lerp.Random(props[0], props[1])
		default:
			val, ok := custom[v]
			if ok {
				return val
			}
		}
	}

	return
}

// FloatTween returns float tween, returns def if process failed
//
// syntax:
//	number - returns lerp.Constant(number)
//  linear start end - if you prefix field with 'linear' and provide two numbers lerp.Linear(number, number)
// is returned
//  bezier start startHandle endHandle end - if you prefix field with 'bezier' and provide four numbers
// lerp.Bezier(start, startHandle, endHandle, end) is returned
//  name - if name is not equal 'random' of 'bezier' it will search for custom tween with this name
func (r RawStyle) FloatTween(key string, custom map[string]lerp.Tween, def lerp.Tween) (e lerp.Tween) {
	e = def
	val, ok := r.Style[key]
	if !ok {
		return
	}

	props := [4]float64{}
	load.CollectFloats(val[1:], props[:])

	switch v := val[0].(type) {
	case int:
		return lerp.Const(float64(v))
	case float64:
		return lerp.Const(v)
	case string:
		switch v {
		case "bezier":
			return lerp.Bezier(props[0], props[1], props[2], props[3])
		case "linear":
			return lerp.Linear(props[0], props[1])
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
//
// syntax:
// 	circle radius resolution - if you prefix property with 'circle' and provide two numbers
// circular drawer will be used
//  rect width height - if you prefix property with 'rect' and provide 2 values, rectangle drawer
// will be used
//	name - providing neather rect nor circle prefix, custom drawer will be searched
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
//
// syntax:
// 	circle radius spread - if you prefix property with 'circle' and provide two numbers
// circular shape will be used, if you omit second value, circle will fall back to default (math.Pi)
//  rect width height - if you prefix property with 'rect' and provide 2 values, rectangle shape
// will be used, if you omit height, width is used for both dimension
//	name - if you provide neither rect nor circle prefix, custom shape will be searched
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

	props := [2]float64{}
	n := load.CollectFloats(val[1:], props[:])

	switch n {
	case 2, 1:
		switch name {
		case "circle":
			if n == 1 {
				return Circular{props[0], math.Pi}
			}
			return Circular{props[0], props[1]}
		case "rect":
			if n == 1 {
				return Rectangle{props[0], props[0]}
			}
			return Rectangle{props[0], props[1]}
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
