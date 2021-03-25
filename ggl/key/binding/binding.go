// Package binding implements simple B system on top of window input
// methods
package binding

import (
	"fmt"

	"github.com/jakubDoka/mlok/ggl"
	"github.com/jakubDoka/mlok/ggl/key"
	"github.com/jakubDoka/mlok/logic/netw"
)

// S is collection of bindings, advantage over using just window is that you can send input over network easily
// and keep inputstate for each player. The way you should set up an inputstate is as follows:
//
//  const (
//		Forward binding.B = iota
//		Right
//     	Left
//      Backwards
//  )
//
//  var Bindings = binging.New(key.W, key.D, key.A, key.S) // its a slice and you index into it with constants
//
// Then each frame you can call Bindings.Update(win), though better approach is to Clone then and store with
// player for example. You can now create bindings with any key combination and listen to them with same
// constants. Other advantage is that you can Write and read bindings to netw.Buffer so stransporting input
// is lot easier.
type S []state

// New creates binding mapping to witch you can index with your constant
func New(keys ...key.Key) S {
	s := make(S, len(keys))
	for i := range s {
		s[i].Key = keys[i]
	}

	return s
}

// Write writes input to the buffer
func (s S) Write(b *netw.Buffer) {
	b.PutUint16(uint16(len(s)))
	for _, bid := range s {
		b.Data = append(b.Data, byte(bid.State))
	}
}

// Read reads S state from buffer
//
// panics if length does not match
func (s S) Read(b *netw.Buffer) {
	l := b.Uint16()
	if l != uint16(len(s)) {
		panic(fmt.Errorf("length of input state (len=%d) written in buffer does not match readers length(len=%d)", l, len(s)))
	}
	for i := range s {
		s[i].State = State(b.Byte())
	}
}

// Clone clones the S (independent copy)
func (s S) Clone() S {
	new := make(S, len(s))
	copy(new, s)
	return new
}

// Update updates S states and returns whether change happened
func (s S) Update(w *ggl.Window) (changed bool) {
	for i := range s {
		b := &s[i]
		old := b.State

		if w.JustPressed(b.Key) {
			b.State = JustPressed
		} else if w.Pressed(b.Key) {
			b.State = Pressed
		} else if w.JustReleased(b.Key) {
			b.State = JustReleased
		} else {
			b.State = Released
		}

		if b.State != old {
			changed = true
		}
	}

	return
}

// State returns state of binding
//
// panics if bid is indexing out of bounds, for example you registered les keys then there is iota constants
func (s S) State(bid B) State {
	if len(s) <= int(bid) {
		panic("constant you provided does not have any binding registered")
	}
	return s[bid].State
}

// Pressed returns whether binding is pressed
//
// panics at sam cases as State
func (s S) Pressed(bid B) bool {
	return s.State(bid) == Pressed
}

// Pressed returns whether binding was pressed in this frame
//
// panics at sam cases as State
func (s S) JustPressed(bid B) bool {
	return s.State(bid) == JustPressed
}

// Pressed returns whether binding was released in this frame
//
// panics at sam cases as State
func (s S) JustReleased(bid B) bool {
	return s.State(bid) == JustReleased
}

type state struct {
	key.Key
	State
}

// B is a arbitrary binding, should be declared with iota
type B uint16

// State stores Binding state
type State uint8

// State enum
const (
	Released State = iota
	JustReleased
	Pressed
	JustPressed
)
