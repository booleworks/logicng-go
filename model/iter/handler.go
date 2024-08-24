package iter

import (
	"github.com/booleworks/logicng-go/event"
)

// A LimitHandler can be used to cancel a model iteration depending on the
// number of found models.
type LimitHandler struct {
	bound            int
	countCommitted   int
	countUncommitted int
}

// HandlerWithLimit generates a new handler which cancels a model iteration
// after the limit of found models is reached.
func HandlerWithLimit(limit int) *LimitHandler {
	return &LimitHandler{
		bound: limit,
	}
}

// ShouldResume processes the given event and returns true if the
// computation should be resumed and false if it should be cancelled.
func (h *LimitHandler) ShouldResume(e event.Event) bool {
	if e == event.ModelEnumerationStarted {
		h.countCommitted = 0
		h.countUncommitted = 0
	} else if e == event.ModelEnumerationCommit {
		h.countCommitted += h.countUncommitted
		h.countUncommitted = 0
	} else if e == event.ModelEnumerationRollback {
		h.countUncommitted = 0
	} else if efm, ok := e.(EventIteratorFoundModels); ok {
		h.countUncommitted += efm.NumberOfModels
	}
	return h.countUncommitted+h.countCommitted < h.bound
}
