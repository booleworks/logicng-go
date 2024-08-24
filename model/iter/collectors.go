package iter

import (
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/sat"
)

// Collector gathers functionality for model iteration collectors.
//
// An iteration collector gathers the found models given by AddModel.  Added
// models can potentially be discarded later via Rollback. To prevent models
// from being rolled back one can call Commit. With Result the result - the
// committed models - can be retrieved.  The generic type R is result type of
// the model iteration function.  It can be e.g. a model count, a list of
// models, or a BDD.
type Collector[R any] interface {
	// AddModel adds a model to the iteration collector.  Returns true if the
	// model was added successfully.
	AddModel(modelFromSolver []bool, solver *sat.Solver, relevantAllIndices []int32, hdl handler.Handler) handler.State

	// Commit confirms all models found since the last commit.  These cannot be
	// rolled back anymore.  Also calls the Commit method of the handler.
	// Returns true if the iteration should continue after the commit.
	Commit(hdl handler.Handler) handler.State

	// Rollback discards all models since the last commit.  Also calls the
	// Rollback method of the handler.  Returns true if the iteration should
	// continue after the commit.
	Rollback(hdl handler.Handler) handler.State

	// RollbackAndReturnModels discards all models since the last commit and
	// returns them.  Also calls the Rollback method of the handler.
	RollbackAndReturnModels(solver *sat.Solver, hdl handler.Handler) []*model.Model

	// Result returns the committed state of the collector .
	Result() R
}
