package ui

import (
	"reflect"
	"testing"

	"github.com/jakubDoka/goml"
	"github.com/jakubDoka/sterr"
)

func TestEmptyScene(t *testing.T) {
	NEmptyScene()
}

func TestParser(t *testing.T) {
	p := NParser()

	testCases := []struct {
		desc   string
		input  string
		output []*Element
		init   func([]*Element)
		err    sterr.Err
	}{
		{
			desc:  "simple",
			input: `<div"/>`,
			output: []*Element{
				{
					Raw: goml.Element{
						Name:       "div",
						Attributes: goml.Attribs{},
					},
					children: NChildren(),
					name:     "0",
				},
			},
		},
		{
			desc:  "nested",
			input: `<div><div/></>`,
			output: []*Element{
				{
					Raw: goml.Element{
						Name:       "div",
						Attributes: goml.Attribs{},
						Children: []goml.Element{
							{
								Name:       "div",
								Attributes: goml.Attribs{},
							},
						},
					},
					children: NChildren(),
					name:     "0",
				},
			},
			init: func(e []*Element) {
				el := e[0]
				el.AddChild("0", &Element{
					Raw: goml.Element{
						Name:       "div",
						Attributes: goml.Attribs{},
					},
					children: NChildren(),
					name:     "0",
				})
			},
		},
		{
			desc:  "basic attributes",
			input: `<div name="s" id="b" styles=["k" "g" "h"]/>`,
			output: []*Element{
				{
					Raw: goml.Element{
						Name: "div",
						Attributes: goml.Attribs{
							"id":     {"b"},
							"name":   {"s"},
							"styles": {"k", "g", "h"},
						},
					},
					children: NChildren(),
					name:     "s",
					id:       "b",
					Styles:   []string{"k", "g", "h"},
				},
			},
		},
		{
			desc:  "text",
			input: `hello`,
			output: []*Element{
				{
					Raw: goml.Element{
						Name: "text",
						Attributes: goml.Attribs{
							"text": {"hello"},
						},
					},
					children: NChildren(),
					Module:   &Text{},
					name:     "0",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.init != nil {
				tC.init(tC.output)
			}
			out, err := p.Parse([]byte(tC.input))
			if !tC.err.SameSurface(err) {
				t.Error(err)
			}

			if err != nil {
				return
			}

			if !reflect.DeepEqual(out, tC.output) {
				for i := range out {
					if i < len(tC.output) {
						t.Errorf("\n%#v\n%#v", out[i], tC.output[i])
						val, _ := tC.output[i].Child("0")
						val2, _ := out[i].Child("0")
						t.Errorf("\n%#v\n%#v", val, val2)
					}
				}
			}
		})
	}
}

func TestGroup(t *testing.T) {
	s := NScene()
	ch := NElement()
	ch.SetGroup("a")

	if ch.Group() != "a" {
		t.Error(ch.Group(), ch.group)
	}

	s.Root.AddChild("ch", ch)

	g := s.Group("a")
	if len(g) != 1 || g[0] != ch {
		t.Errorf("%v %p", g, ch)
	}

	ch2 := NElement()
	s.Root.AddChild("ch2", ch2)
	ch2.SetGroup("a")

	g = s.Group("a")
	if len(g) != 2 || g[1] != ch2 {
		t.Errorf("%v %p", g, ch2)
	}

	ch.SetGroup("b")
	g = s.Group("a")
	if len(g) != 1 || g[0] != ch2 {
		t.Errorf("%v %p", g, ch2)
	}

	g = s.Group("b")
	if len(g) != 1 || g[0] != ch {
		t.Errorf("%v %p", g, ch)
	}
}

func TestID(t *testing.T) {
	s := NScene()
	ch := NElement()
	ch.SetID("a")

	if ch.ID() != "a" {
		t.Error(ch.ID(), ch.id)
	}

	s.Root.AddChild("ch", ch)

	id := s.ID("a")
	if id != ch {
		t.Error(id, ch)
	}

	ch.SetID("b")

	id = s.ID("a")
	if id != nil {
		t.Error(id, ch)
	}

	id = s.ID("b")
	if id != ch {
		t.Error(id, ch)
	}
}

func TestName(t *testing.T) {
	s := NScene()
	ch := NElement()

	ch.SetName("b")

	s.Root.AddChild("a", ch)

	if ch.Name() != "a" {
		t.Error(ch.Name(), "a")
	}

	ch.SetName("b")

	if ch.Name() != "b" {
		t.Error(ch.Name(), "b")
	}

	c, _ := s.Root.Child("b")
	if c != ch {
		t.Error(c, ch)
	}

	ok := s.Root.Rename("a", "b")
	if ok {
		t.Error("unexpected success")
	}

	ok = s.Root.Rename("b", "a")
	if !ok {
		t.Error("Unexpected fail")
	}
}

func TestIndex(t *testing.T) {
	s := NScene()
	ch := NElement()

	s.Root.AddChild("a", ch)

	ch2 := NElement()

	s.Root.AddChild("b", ch2)
	s.Root.AddChild("c", NElement())

	if ch.Index() != ch2.Index()-1 {
		t.Error(ch.Index(), ch2.Index())
	}

	s.Root.ReIndex(1, 0)

	if ch.Index()-1 != ch2.Index() {
		t.Error(ch.Index(), ch2.Index())
	}

	if s.Root.ChildAt(0) != ch2 {
		t.Error(s.Root.ChildAt(0), ch2)
	}

	if s.Root.ChildAt(1) != ch {
		t.Error(s.Root.ChildAt(1), ch)
	}
}

func TestFindChild(t *testing.T) {
	s := NScene()
	ch := NElement()

	s.Root.AddChild("a", ch)

	ch2 := NElement()

	ch.AddChild("a", ch2)

	ch3 := NElement()

	ch2.AddChild("b", ch3)

	ch.AddChild("f", NElement())

	var coll []*Element

	testCases := []struct {
		desc  string
		res   []*Element
		query string
		cap   int
	}{
		{
			desc:  "unlimited",
			res:   []*Element{ch, ch2},
			query: "a",
			cap:   -1,
		},
		{
			desc:  "limited",
			res:   []*Element{ch},
			query: "a",
			cap:   1,
		},
		{
			desc:  "deep",
			res:   []*Element{ch3},
			query: "b",
			cap:   0,
		},
		{
			desc:  "nothing",
			query: "c",
			cap:   -1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			coll := coll[:0]
			s.Root.FindChild(tC.query, tC.cap, &coll)
			if !reflect.DeepEqual(tC.res, coll) {
				t.Errorf("\n%#v\n%#v", tC.res, coll)
			}
		})
	}
}

func TestPath(t *testing.T) {
	s := NScene()
	ch := NElement()

	s.Root.AddChild("a", ch)

	ch2 := NElement()

	ch.AddChild("a", ch2)

	ch3 := NElement()

	ch2.AddChild("b", ch3)

	if ch3.Path() != "root.a.a.b" {
		t.Error(ch3.Path())
	}
}

func TestInsertChild(t *testing.T) {
	s := NScene()
	ch := NElement()

	s.Root.AddChild("a", ch)

	ch2 := NElement()

	s.Root.AddChild("b", ch2)

	ch3 := NElement()

	s.Root.AddChild("c", ch3)

	ch4 := NElement()

	s.Root.InsertChild("d", 1, ch4)

	ch5 := NElement()

	s.Root.InsertChild("e", 4, ch5)

	if ch.Index() != ch4.Index()-1 || ch2.Index() != ch3.Index()-1 || ch5.Index() != 4 {
		t.Error(ch.Index(), ch4.Index(), ch2.Index(), ch3.Index(), ch5.Index())
	}
}
