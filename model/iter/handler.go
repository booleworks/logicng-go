package iter

import (
	"booleworks.com/logicng/handler"
	"booleworks.com/logicng/sat"
)

// Handler describes the functionality of a model iteration handler.
type Handler interface {
	handler.Handler

	// SatHandler returns the embedded SAT handler.
	SatHandler() sat.Handler

	// FoundModels is called every time new models are found.  The found models
	// are in an uncommitted state until they are confirmed by calling Commit.
	// It is also possible to roll back the uncommitted models by calling
	// Rollback.  Returns true if the iteration should continue.
	FoundModels(numberOfModels int) bool

	// Commit is called every time the models are committed.  Returns true if
	// the iteration should continue.
	Commit() bool

	// Rollback is called everytime uncommitted models are rolled back.  Returns
	// true if the iteration should continue.
	Rollback() bool
}

// A TimeoutHandler can be used to abort a model iteration depending on a
// timeout set beforehand.
type TimeoutHandler struct {
	handler.Timeout
	satHandler sat.Handler
}

// HandlerWithTimeout generates a new timeout handler with the given timeout.
func HandlerWithTimeout(timeout handler.Timeout) *TimeoutHandler {
	satHandler := sat.HandlerWithTimeout(timeout)
	return &TimeoutHandler{timeout, satHandler}
}

// SatHandler returns the embedded SAT handler.
func (t *TimeoutHandler) SatHandler() sat.Handler {
	return t.satHandler
}

// FoundModels is called every time new models are found.  The found models
// are in an uncommitted state until they are confirmed by calling Commit.
// It is also possible to roll back the uncommitted models by calling
// Rollback.  Returns true if the iteration should continue.
func (t *TimeoutHandler) FoundModels(_ int) bool {
	return !t.TimeLimitExceeded()
}

// Commit is called every time the models are committed.  Returns true if
// the iteration should continue.
func (t *TimeoutHandler) Commit() bool {
	return !t.TimeLimitExceeded()
}

// Rollback is called everytime uncommitted models are rolled back.  Returns
// true if the iteration should continue.
func (t *TimeoutHandler) Rollback() bool {
	return !t.TimeLimitExceeded()
}

// A LimitHandler can be used to abort a model iteration depending on the
// number of found models.
type LimitHandler struct {
	handler.Computation
	bound            int
	countCommitted   int
	countUncommitted int
}

// HandlerWithLimit generates a new handler which aborts a model iteration
// after the limit of found models is reached.
func HandlerWithLimit(limit int) *LimitHandler {
	return &LimitHandler{
		Computation: handler.Computation{},
		bound:       limit,
	}
}

// SatHandler returns the embedded SAT handler.
func (l *LimitHandler) SatHandler() sat.Handler {
	return nil
}

// FoundModels is called every time new models are found.  The found models
// are in an uncommitted state until they are confirmed by calling Commit.
// It is also possible to roll back the uncommitted models by calling
// Rollback.  Returns true if the iteration should continue.
func (l *LimitHandler) FoundModels(numberOfModels int) bool {
	l.SetAborted(l.countUncommitted+l.countCommitted+numberOfModels > l.bound)
	if !l.Aborted() {
		l.countUncommitted += numberOfModels
		return true
	} else {
		return false
	}
}

// Commit is called every time the models are committed.  Returns true if
// the iteration should continue.
func (l *LimitHandler) Commit() bool {
	l.countCommitted += l.countUncommitted
	l.countUncommitted = 0
	return true
}

// Rollback is called everytime uncommitted models are rolled back.  Returns
// true if the iteration should continue.
func (l *LimitHandler) Rollback() bool {
	l.countUncommitted = 0
	return true
}
