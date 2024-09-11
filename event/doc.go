// Package event provides datastructures for events.  Events are emitted within
// computations and handlers can then react to these events and potentially
// cancel computations.  Some computations define their own events with more
// information, e.g. an optimization can include the current bound in the event
// such that a handler can react to it.
package event
