package events

import "reflect"

// Handler allows callback registration and event fireing, its simply
// a way of comunication that comes handy in some cases, though is has
// performance cost, events are mapped to reflect.Type of an inputted
// anything, enum events are accessed by indexing into slice, if you
// don't have to pass arguments with event use Handler.Trigger and Handler.Run
// as it is cleaner and faster
type Handler struct {
	callbacks     map[reflect.Type][]func(interface{})
	enumCallbacks [][]func()
}

// NHandler initializes inner map
func NHandler() *Handler {
	return &Handler{
		callbacks: map[reflect.Type][]func(interface{}){},
	}
}

// Fire launches all event callbacks for given type
func (e *Handler) Fire(ev interface{}) {
	for _, e := range e.callbacks[reflect.TypeOf(ev)] {
		e(ev)
	}
}

// On registers event callback
func (e *Handler) On(ev interface{}, f func(interface{})) {
	e.callbacks[reflect.TypeOf(ev)] = append(e.callbacks[reflect.TypeOf(ev)], f)
}

// Trigger launches enum event callbacks
func (e *Handler) Trigger(en Enum) {
	e.ensureLength(en)
	for _, e := range e.enumCallbacks[en] {
		e()
	}
}

// Run register enum event callback
func (e *Handler) Run(en Enum, f func()) {
	e.ensureLength(en)
	e.enumCallbacks[en] = append(e.enumCallbacks[en], f)
}

// as we cannot know how match different EnEvents we have, its unreliable to hardcode a value
func (e *Handler) ensureLength(en Enum) {
	for int(en) >= len(e.enumCallbacks) {
		e.enumCallbacks = append(e.enumCallbacks, []func(){})
	}
}

// Enum is to distinguish event enumeration
type Enum int
