package handler

import "github.com/booleworks/logicng-go/event"

// A Handler can be used to cancel computations.  It reacts to different kinds
// of events.
type Handler interface {
	// ShouldResume processes the given event and returns true if the
	// computation should be resumed and false if it should be cancelled.
	ShouldResume(event.Event) bool
}

// A NopHandler never cancels the computation (equivalent to no handler).
var NopHandler = nopHandler{}

type nopHandler struct{}

func (nopHandler) ShouldResume(event.Event) bool {
	return true
}

// The State contains the information if a handler was cancelled and
// if so, which was the event which caused the cancellation.  If the handler
// was not cancelled, the cause is the "Nothing" event.
type State struct {
	Success     bool
	CancelCause event.Event
}

// Success generates a new successful handler state where the handler was
// not cancelled.
func Success() State {
	return State{true, event.Nothing}
}

// Cancellation generates a new handler state where the handler was cancelled
// with the given event as cause.
func Cancellation(cancelCause event.Event) State {
	return State{false, cancelCause}
}
