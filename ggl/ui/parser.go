package ui

import (
	"strconv"
	"strings"

	"github.com/jakubDoka/gobatch/logic/event"
	"github.com/jakubDoka/goml"
	"github.com/jakubDoka/goml/goss"
	"github.com/jakubDoka/sterr"
)

// ParserRelated errors
var (
	ErrGoml           = sterr.New("error in goml")
	ErrUrlp           = sterr.New("error when loading attributes")
	ErrPath           = sterr.New("in element %v")
	ErrMissingFactory = sterr.New("missing factory with name %s")
)

// NoFactory is used for creating div
type NoFactory struct{}

// New implements ModuleFactory interface
func (*NoFactory) New() Module { return nil }

// Parser handles element parsing form goml, if you don't know goml syntax read
// the tutorial on github.com/jakubDoka/goml
type Parser struct {
	factories map[string]ModuleFactory
	GP        *goml.Parser
	GS        goss.Parser
}

// NParser creates ready-to-use Parser
func NParser() *Parser {
	p := &Parser{
		factories: map[string]ModuleFactory{},
	}

	p.GP = goml.NParser(&p.GS)

	p.AddFactory("div", &ModuleBase{})
	p.AddFactory("text", &Text{})
	p.AddFactory("scroll", &Scroll{})
	p.AddFactory("sprite", &Sprite{})
	p.AddFactory("patch", &Patch{})
	p.AddFactory("button", &Button{})
	p.AddFactory("area", &Area{})

	return p
}

// AddFactory adds a element factory to parser, every element with given name
// will receive module from the factory
func (p *Parser) AddFactory(name string, mf ModuleFactory) {
	p.factories[name] = mf
	p.GP.AddDefinitions(name)
}

// Parse parses goml source into list of elements
func (p *Parser) Parse(source []byte) ([]*Element, error) {
	div, err := p.GP.Parse(source)
	if err != nil {
		return nil, ErrGoml.Wrap(err)
	}
	elems := make([]*Element, len(div.Children))
	for i, e := range div.Children {
		ch, err := p.translateElement(i, e)
		if err != nil {
			return nil, err
		}
		elems[i] = ch
	}
	return elems, nil
}

func (p *Parser) translateElement(i int, elem goml.Element) (*Element, error) {
	val, ok := p.factories[elem.Name]
	if !ok {
		return nil, ErrMissingFactory.Args(elem.Name)
	}
	e := &Element{Module: val.New(), children: NChildren(), Raw: elem, Events: event.String{}}

	if val, ok := elem.Attributes["name"]; ok {
		e.name = val[0]
	} else {
		e.name = strconv.Itoa(i)
	}
	if _, ok := elem.Attributes["hidden"]; ok {
		e.hidden = true
	}
	if val, ok := elem.Attributes["id"]; ok {
		e.id = val[0]
	}
	if val, ok := elem.Attributes["group"]; ok {
		e.group = val[0]
	}
	if val, ok := elem.Attributes["styles"]; ok {
		if len(val) == 1 && strings.Contains(val[0], " ") { // make it more friendly, both list and string is valid
			e.Styles = strings.Split(val[0], " ")
		} else {
			e.Styles = val
		}
	}

	for i, ch := range elem.Children {
		el, err := p.translateElement(i, ch)
		if err != nil {
			return nil, ErrPath.Args(e.name).Wrap(err)
		}
		e.AddChild(el.name, el)
	}

	return e, nil
}

// ModuleFactory should be an producer of module instances for parser
// it gives you option to handle initialization your self, witch you can
// signalize by returning true
type ModuleFactory interface {
	New() Module
}
