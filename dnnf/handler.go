package dnnf

import "github.com/booleworks/logicng-go/handler"

// A Handler for a DNNF can abort the compilation of a DNNF.  The method
// ShannonExpansion is called after each performed Shannon expansion within the
// DNNF compiler.
type Handler interface {
	handler.Handler
	ShannonExpansion() bool
}

// A TimeoutHandler aborts the DNNF compilation dependent on the time it takes
// to generate the DNNF.
type TimeoutHandler struct {
	handler.Timeout
}

// HandlerWithTimeout returns a new DNNF TimoutHandler for the given timeout.
func HandlerWithTimeout(timeout handler.Timeout) *TimeoutHandler {
	return &TimeoutHandler{timeout}
}

// ShannonExpansion is called by the DNNF compiler everytime a Shannon
// expansion is performed.
func (t *TimeoutHandler) ShannonExpansion() bool {
	return !t.TimeLimitExceeded()
}
