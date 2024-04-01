package maxsat

import (
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/sat"
)

// Handler is an interface for MAX-SAT handlers which can abort the computation
// based on the upper and lower bounds found during the computation.
type Handler interface {
	handler.Handler
	SatHandler() sat.Handler
	FoundLowerBound(lowerBound int, model *model.Model) bool
	FoundUpperBound(upperBound int, model *model.Model) bool
	FinishedSolving()
	LowerBoundApproximation() int
	UpperBoundApproximation() int
}

// A TimeoutHandler can be used to abort a MAX-SAT computation depending on a
// timeout set beforehand.
type TimeoutHandler struct {
	handler.Timeout
	satHandler sat.Handler
	currentLb  int
	currentUb  int
}

// HandlerWithTimeout generates a new timeout handler with the given timeout.
func HandlerWithTimeout(timeout handler.Timeout) *TimeoutHandler {
	satHandler := sat.HandlerWithTimeout(timeout)
	return &TimeoutHandler{timeout, satHandler, -1, -1}
}

// FoundLowerBound is called by the MAX-SAT solver each time a new lower bound
// is found.  The current model for this bound is recorded.  Returns whether
// the computation should be continued.
func (t *TimeoutHandler) FoundLowerBound(lowerBound int, _ *model.Model) bool {
	t.currentLb = lowerBound
	return !t.TimeLimitExceeded()
}

// FoundUpperBound is called by the MAX-SAT solver each time a new upper bound
// is found.  The current model for this bound is recorded.  Returns whether
// the computation should be continued.
func (t *TimeoutHandler) FoundUpperBound(upperBound int, _ *model.Model) bool {
	t.currentUb = upperBound
	return !t.TimeLimitExceeded()
}

// SatHandler returns the underlying SAT handler for the MAX-SAT handler.
func (t *TimeoutHandler) SatHandler() sat.Handler { return t.satHandler }

// FinishedSolving is called when the MAX-SAT solver has finished the solving
// process.
func (t *TimeoutHandler) FinishedSolving() {}

// LowerBoundApproximation returns the last found lower bound.
func (t *TimeoutHandler) LowerBoundApproximation() int { return t.currentLb }

// UpperBoundApproximation returns the last found upper bound.
func (t *TimeoutHandler) UpperBoundApproximation() int { return t.currentUb }
