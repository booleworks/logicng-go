package sat

import (
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model"
)

// Handler is an interface for SAT handlers which can abort the computation
// based on the number of conflicts found during the computation.
type Handler interface {
	handler.Handler
	DetectedConflict() bool
	FinishedSolving()
}

func handlerFinishSolving(handler Handler) {
	if handler != nil {
		handler.FinishedSolving()
	}
}

// A TimeoutHandler can be used to abort a SAT computation depending on a
// timeout set beforehand.
type TimeoutHandler struct {
	handler.Timeout
}

// HandlerWithTimeout generates a new timeout handler with the given timeout.
func HandlerWithTimeout(timeout handler.Timeout) *TimeoutHandler {
	return &TimeoutHandler{timeout}
}

// DetectedConflict is calles by the solver each time a conflict is detected.
func (t *TimeoutHandler) DetectedConflict() bool {
	return !t.TimeLimitExceeded()
}

// FinishedSolving is called when the SAT solver has finished the solving
// process.
func (t *TimeoutHandler) FinishedSolving() {
	// do nothing
}

// OptimizationHandler is an interface for SAT-based optimizations which can
// abort the computation everytime a better solution is found during the
// optimization.
type OptimizationHandler interface {
	handler.Handler
	SatHandler() Handler
	FoundBetterBound(model *model.Model) bool
	SetModel(model *model.Model)
}

// A TimeoutOptimizationHandler can be used to abort a SAT optimization
// depending on a timeout set beforehand.
type TimeoutOptimizationHandler struct {
	handler.Timeout
	satHandler Handler
	lastModel  *model.Model
}

// OptimizationHandlerWithTimeout generates a new timeout handler with the
// given timeout.
func OptimizationHandlerWithTimeout(timeout handler.Timeout) *TimeoutOptimizationHandler {
	satHandler := HandlerWithTimeout(timeout)
	return &TimeoutOptimizationHandler{timeout, satHandler, nil}
}

// SatHandler returns the SAT handler of the optimization handler.
func (t *TimeoutOptimizationHandler) SatHandler() Handler {
	return t.satHandler
}

// FoundBetterBound is called everytime a better solution bound is found on
// the SAT solver.
func (t *TimeoutOptimizationHandler) FoundBetterBound(model *model.Model) bool {
	t.lastModel = model
	return !t.TimeLimitExceeded()
}

// SetModel is called by the solver with the current model everytime a better
// solutions bound is found.
func (t *TimeoutOptimizationHandler) SetModel(model *model.Model) {
	t.lastModel = model
}

// IntermediateResult returns the last found model of the solver.
func (t *TimeoutOptimizationHandler) IntermediateResult() *model.Model {
	return t.lastModel
}

func satHandler(optimizationHandler OptimizationHandler) Handler {
	if optimizationHandler == nil {
		return nil
	} else {
		return optimizationHandler.SatHandler()
	}
}
