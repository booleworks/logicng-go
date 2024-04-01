package handler

// A Handler can be used to abort computations.  The Started method is called
// when the computation is started.  The Aborted method returns whether the
// computation was aborted by the handler.
type Handler interface {
	Started()
	Aborted() bool
}

// Aborted returns true when the given handler is not nil and the computation
// was aborted by this handler.
func Aborted(handler Handler) bool {
	return handler != nil && handler.Aborted()
}

// Start starts the given handler it is not nil.
func Start(handler Handler) {
	if handler != nil {
		handler.Started()
	}
}

// Computation is a simple computation handler which can be embedded in more
// complex handlers.
type Computation struct {
	aborted bool
}

// Started indicates the handler that the computation was started.
func (c *Computation) Started() {
	c.aborted = false
}

// Aborted reports whether the computation was aborted by the handler.
func (c *Computation) Aborted() bool {
	return c.aborted
}

// SetAborted sets whether the computation was aborted.
func (c *Computation) SetAborted(aborted bool) {
	c.aborted = aborted
}
